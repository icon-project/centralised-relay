package relayer

import (
	"context"
	"fmt"
	"time"

	"github.com/icon-project/centralized-relay/relayer/store"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	DefaultFlushInterval      = 5 * time.Minute
	listenerChannelBufferSize = 1000
	DefaultTxRetry            = 2
	SaveHeightMaxAfter        = 1000
	RouteDuration             = 1 * time.Second

	prefixMessageStore = "message"
	prefixBlockStore   = "block"
)

// main start loop
func Start(
	ctx context.Context,
	log *zap.Logger,
	chains map[string]*Chain,
	flushInterval time.Duration,
	fresh bool,
	db store.Store,
) (chan error, error) {
	errorChan := make(chan error, 1)
	relayer, err := NewRelayer(log, db, chains, fresh)
	if err != nil {
		return nil, fmt.Errorf("error creating new relayer %v", err)
	}

	// start all the chain listeners
	go relayer.StartChainListeners(ctx, errorChan)

	// start all the block processor
	go relayer.StartBlockProcessors(ctx, errorChan)

	// responsible to relaying  messages
	go relayer.StartRouter(ctx, flushInterval, fresh)

	return errorChan, nil
}

type Relayer struct {
	log          *zap.Logger
	chains       map[string]*ChainRuntime
	messageStore *store.MessageStore
	blockStore   *store.BlockStore
}

func NewRelayer(log *zap.Logger, db store.Store, chains map[string]*Chain, fresh bool) (*Relayer, error) {

	// if fresh clearing db
	if fresh {
		err := db.ClearStore()
		if err != nil {
			return nil, err
		}
	}

	// initializing message store
	messageStore := store.NewMessageStore(db, prefixMessageStore)

	// blockStore store
	blockStore := store.NewBlockStore(db, prefixBlockStore)

	chainRuntimes := make(map[string]*ChainRuntime, len(chains))
	for _, chain := range chains {
		chainRuntime, err := NewChainRuntime(log, chain)
		if err != nil {
			return nil, err
		}

		lastSavedHeight, err := blockStore.GetLastStoredBlock(chain.ChainID())
		if err == nil {
			// successfully fetched last savedBlock
			chainRuntime.LastSavedHeight = lastSavedHeight
		}
		chainRuntimes[chain.ChainID()] = chainRuntime

	}

	return &Relayer{
		log:          log,
		chains:       chainRuntimes,
		messageStore: messageStore,
		blockStore:   blockStore,
	}, nil
}

func (r *Relayer) StartChainListeners(
	ctx context.Context,
	errCh chan error,
) {

	var eg errgroup.Group
	for _, chainRuntime := range r.chains {
		chainRuntime := chainRuntime

		eg.Go(func() error {
			// listening to the block
			err := chainRuntime.Provider.Listener(ctx, chainRuntime.LastSavedHeight, chainRuntime.listenerChan)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		errCh <- err
	}
}

func (r *Relayer) StartBlockProcessors(ctx context.Context, errorChan chan error) {
	var eg errgroup.Group

	for _, chainRuntime := range r.chains {
		chainRuntime := chainRuntime
		listener := chainRuntime.listenerChan
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case blockInfo, ok := <-listener:
					if !ok {
						return nil
					}
					r.processBlockInfo(ctx, chainRuntime, blockInfo)
				}
			}
		})
	}

	if err := eg.Wait(); err != nil {
		errorChan <- err // Report the error to the error channel.
	}
}

func (r *Relayer) StartRouter(ctx context.Context, flushInterval time.Duration, fresh bool) {

	routeTimer := time.NewTicker(RouteDuration)

	// TODO: implement flush logic
	for {
		select {
		case <-routeTimer.C:
			r.processMessages(ctx)
		}
	}
}

func (r *Relayer) processMessages(ctx context.Context) {

	for _, srcChainRuntime := range r.chains {
		for _, routeMessage := range srcChainRuntime.MessageCache.Messages {
			dstChainRuntime, err := r.FindChainRuntime(routeMessage.Dst)
			if err != nil {
				r.log.Error("dst chain runtime not found ", zap.String("dst chain", routeMessage.Dst))
				continue
			}
			ok := dstChainRuntime.shouldSendMessage(ctx, routeMessage, srcChainRuntime)
			if !ok {
				continue
			}
			go r.RouteMessage(ctx, routeMessage, dstChainRuntime, srcChainRuntime)
		}
	}

}

// processBlockInfo->
// save block height to database
// & merge message to src cache
func (r *Relayer) processBlockInfo(ctx context.Context, srcChainRuntime *ChainRuntime, blockInfo types.BlockInfo) {
	err := r.SaveBlockHeight(ctx, srcChainRuntime, blockInfo.Height, len(blockInfo.Messages))
	if err != nil {
		r.log.Error("unable to save height", zap.Error(err))
	}

	go srcChainRuntime.mergeMessages(ctx, blockInfo)
}

func (r *Relayer) SaveBlockHeight(ctx context.Context, chainRuntime *ChainRuntime, height uint64, messageCount int) error {
	r.log.Debug("saving height:", zap.String("srcChain", chainRuntime.Provider.ChainId()), zap.Uint64("height", height))

	if messageCount > 0 || (height-chainRuntime.LastSavedHeight) > uint64(SaveHeightMaxAfter) {
		chainRuntime.LastSavedHeight = height
		err := r.blockStore.StoreBlock(height, chainRuntime.Provider.ChainId())
		if err != nil {
			return fmt.Errorf("error while saving height of chain:%s %v", chainRuntime.Provider.ChainId(), err)
		}
	}
	return nil
}

func (r *Relayer) FindChainRuntime(chainId string) (*ChainRuntime, error) {
	var chainRuntime *ChainRuntime
	var ok bool

	if chainRuntime, ok = r.chains[chainId]; !ok {
		return nil, fmt.Errorf("chain runtime not found, chainId:%s ", chainId)
	}

	return chainRuntime, nil
}

func (r *Relayer) RouteMessage(ctx context.Context, m *types.RouteMessage, dst, src *ChainRuntime) {
	callback := func(response types.ExecuteMessageResponse) {
		src := src
		dst := dst

		if response.Code == types.Success {
			dst.log.Info("Successfully relayed message:",
				zap.String("src chain", src.Provider.ChainId()),
				zap.String("dst chain", dst.Provider.ChainId()),
				zap.Uint64("Sn number", response.Sn),
				zap.Any("Tx hash", response.TxHash),
			)
			if err := r.ClearMessages(ctx, []types.MessageKey{response.MessageKey}, src); err != nil {
				r.log.Error("error occured when clearing successful message", zap.Error(err))
			}
			return
		}

		routeMessage, ok := src.MessageCache.Messages[response.MessageKey]
		if !ok {
			r.log.Error("message not found for key", zap.Any("message key", response.MessageKey))
			return
		}

		// open to resend the message
		if routeMessage.GetRetry() >= uint64(DefaultTxRetry) {

			// save to db
			if err := r.messageStore.StoreMessage(routeMessage.Message); err != nil {
				r.log.Error("error occured when storing the message after max retry", zap.Error(err))
				return
			}
			// removed message from messageCache
			src.MessageCache.Remove(routeMessage.MessageKey())

			dst.log.Error("failed to send message saving to database",
				zap.String("src chain", routeMessage.Src),
				zap.String("dst chain", routeMessage.Dst),
				zap.Uint64("Sn number", routeMessage.Sn),
			)
			return
		}
		routeMessage.SetIsProcessing(false)
	}

	// setting before message is processed
	m.SetIsProcessing(true)
	m.IncrementRetry()

	err := dst.Provider.Route(ctx, m, callback)
	if err != nil {
		dst.log.Error("error occured during message route", zap.Error(err))
	}

}

func (r *Relayer) ClearMessages(ctx context.Context, msgs []types.MessageKey, srcChain *ChainRuntime) error {

	// clear from cache
	srcChain.clearMessageFromCache(msgs)

	for _, m := range msgs {
		err := r.messageStore.DeleteMessage(m)
		if err != nil {
			r.log.Error("error occured when deleting message from db ", zap.Error(err))
		}
	}
	return nil
}
