package evm

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/chains/evm/types"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/pkg/errors"
)

const (
	BlockInterval              = 2 * time.Second
	BlockHeightPollInterval    = 60 * time.Second
	defaultReadTimeout         = 15 * time.Second
	monitorBlockMaxConcurrency = 1000 // number of concurrent requests to synchronize older blocks from source chain
)

type BnOptions struct {
	StartHeight uint64
	Concurrency uint64
}

func (r *EVMProvider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfoChan chan relayertypes.BlockInfo) error {

	startHeight, err := r.startFromHeight(ctx, lastSavedHeight)
	if err != nil {
		return err
	}

	concurrency := r.GetConcurrency(ctx)
	r.log.Info("Starting Evm listener from height ", zap.Uint64("start-height", startHeight))

	// block notification channel
	// (buffered: to avoid deadlock)
	// increase concurrency parameter for faster sync
	bnch := make(chan *types.BlockNotification, concurrency)

	heightTicker := time.NewTicker(BlockInterval)
	defer heightTicker.Stop()

	heightPoller := time.NewTicker(BlockHeightPollInterval)
	defer heightPoller.Stop()

	latestHeight := func() uint64 {
		height, err := r.client.eth.BlockNumber(context.TODO())
		if err != nil {
			// TODO:
			// r.Log.WithFields(log.Fields{"error": err}).Error("receiveLoop: failed to GetBlockNumber")
			return 0
		}
		return height
	}

	// Loop started
	next, latest := startHeight, latestHeight()
	// last unverified block notification
	var lbn *types.BlockNotification
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-heightTicker.C:
			latest++

		case <-heightPoller.C:
			if height := latestHeight(); height > latest {
				latest = height
				if next > latest {
					// TODO:
					// r.Log.Debugf("receiveLoop: skipping; latest=%d, next=%d", latest, next)
				}
			}

		case bn := <-bnch:
			// process all notifications
			for ; bn != nil; next++ {
				if lbn != nil {
					fmt.Println("block-notification received evm: ", lbn.Height)

					messages, err := r.FindMessages(ctx, lbn)
					if err != nil {
						return errors.Wrapf(err, "receiveLoop: callback: %v", err)
					}
					blockInfoChan <- relayertypes.BlockInfo{
						Height:   lbn.Height.Uint64(),
						Messages: messages,
					}
				}

				if lbn, bn = bn, nil; len(bnch) > 0 {
					bn = <-bnch
				}
			}
			// remove unprocessed notifications
			for len(bnch) > 0 {
				<-bnch
			}

		default:
			if next >= latest {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			type bnq struct {
				h     uint64
				v     *types.BlockNotification
				err   error
				retry int
			}

			qch := make(chan *bnq, cap(bnch))
			for i := next; i < latest &&
				len(qch) < cap(qch); i++ {
				qch <- &bnq{i, nil, nil, 3} // fill bch with requests
			}
			bns := make([]*types.BlockNotification, 0, len(qch))
			for q := range qch {
				switch {
				case q.err != nil:
					if q.retry > 0 {
						if !strings.HasSuffix(q.err.Error(), "requested block number greater than current block number") {
							q.retry--
							q.v, q.err = nil, nil
							qch <- q
							continue
						}
						if latest >= q.h {
							latest = q.h - 1
						}
					}
					//r.Log.Debugf("receiveLoop: bnq: h=%d:%v, %v", q.h, q.v.Header.Hash(), q.err)
					bns = append(bns, nil)
					if len(bns) == cap(bns) {
						close(qch)
					}

				case q.v != nil:
					bns = append(bns, q.v)
					if len(bns) == cap(bns) {
						close(qch)
					}
				default:
					go func(q *bnq) {
						defer func() {
							time.Sleep(500 * time.Millisecond)
							qch <- q
						}()
						if q.v == nil {
							q.v = &types.BlockNotification{}
						}
						q.v.Height = (&big.Int{}).SetUint64(q.h)
						q.v.Header, q.err = r.client.eth.HeaderByNumber(context.TODO(), q.v.Height)
						if q.err != nil {
							//q.err = errors.Wrapf(q.err, "GetHmyHeaderByHeight: %v", q.err)
							return
						}
						if q.v.Header.GasUsed > 0 {
							ht := big.NewInt(q.v.Height.Int64())
							r.BlockReq.FromBlock = ht
							r.BlockReq.ToBlock = ht
							q.v.Logs, q.err = r.client.eth.FilterLogs(context.TODO(), r.BlockReq)
							if q.err != nil {
								q.err = errors.Wrapf(q.err, "FilterLogs: %v", q.err)
								return
							}
						}
					}(q)
				}
			}
			// filter nil
			_bns_, bns := bns, bns[:0]
			for _, v := range _bns_ {
				if v != nil {
					bns = append(bns, v)
				}
			}
			// sort and forward notifications
			if len(bns) > 0 {
				sort.SliceStable(bns, func(i, j int) bool {
					return bns[i].Height.Uint64() < bns[j].Height.Uint64()
				})
				for i, v := range bns {
					if v.Height.Uint64() == next+uint64(i) {
						bnch <- v
					}
				}
			}
		}
	}
}

func (p *EVMProvider) FindMessages(ctx context.Context, lbn *types.BlockNotification) ([]relayertypes.Message, error) {

	return nil, nil

}

func (p *EVMProvider) GetConcurrency(ctx context.Context) int {

	// TODO: get concurrency from config
	// if opts.Concurrency < 1 || opts.Concurrency > monitorBlockMaxConcurrency {
	// 	concurrency := opts.Concurrency
	// 	if concurrency < 1 {
	// 		opts.Concurrency = 1
	// 	} else {
	// 		opts.Concurrency = monitorBlockMaxConcurrency
	// 	}
	// 	// r.Log.Warnf("receiveLoop: opts.Concurrency (%d): value out of range [%d, %d]: setting to default %d",
	// 	// concurrency, 1, monitorBlockMaxConcurrency, opts.Concurrency)
	// }
	return monitorBlockMaxConcurrency
}

func (p *EVMProvider) startFromHeight(ctx context.Context, lastSavedHeight uint64) (uint64, error) {
	latestHeight, err := p.QueryLatestHeight(ctx)
	if err != nil {
		return 0, err
	}

	if p.cfg.StartHeight > latestHeight {
		p.log.Error("start height provided on config cannot be greater than latest height",
			zap.Uint64("start-height", p.cfg.StartHeight),
			zap.Int64("latest-height", int64(latestHeight)),
		)
	}

	// priority1: startHeight from config
	if p.cfg.StartHeight != 0 && p.cfg.StartHeight < latestHeight {
		return p.cfg.StartHeight, nil
	}

	// priority2: lastsaveheight from db
	if lastSavedHeight != 0 && lastSavedHeight < latestHeight {
		return lastSavedHeight, nil
	}

	// priority3: latest height
	return latestHeight, nil
}