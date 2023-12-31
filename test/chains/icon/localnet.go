package icon

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/icon-bridge/common/wallet"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	//chantypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	dockertypes "github.com/docker/docker/api/types"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/icon-project/centralized-relay/test/chains"
	icontypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type IconLocalnet struct {
	log           *zap.Logger
	testName      string
	cfg           ibc.ChainConfig
	numValidators int
	numFullNodes  int
	FullNodes     IconNodes
	findTxMu      sync.Mutex
	keystorePath  string
	scorePaths    map[string]string
	IBCAddresses  map[string]string     `json:"addresses"`
	Wallets       map[string]ibc.Wallet `json:"wallets"`
}

func (c *IconLocalnet) CreateKey(ctx context.Context, keyName string) error {
	//TODO implement me
	panic("implement me")
}

func NewIconLocalnet(testName string, log *zap.Logger, chainConfig ibc.ChainConfig, numValidators int, numFullNodes int, scorePaths map[string]string) chains.Chain {
	return &IconLocalnet{
		testName:      testName,
		cfg:           chainConfig,
		numValidators: numValidators,
		numFullNodes:  numFullNodes,
		log:           log,
		scorePaths:    scorePaths,
		Wallets:       map[string]ibc.Wallet{},
		IBCAddresses:  make(map[string]string),
	}
}

// Config fetches the chain configuration.
func (c *IconLocalnet) Config() ibc.ChainConfig {
	return c.cfg
}

func (c *IconLocalnet) OverrideConfig(key string, value any) {
	if value == nil {
		return
	}
	c.cfg.ConfigFileOverrides[key] = value
}

// Initialize initializes node structs so that things like initializing keys can be done before starting the chain
func (c *IconLocalnet) Initialize(ctx context.Context, testName string, cli *client.Client, networkID string) error {
	chainCfg := c.Config()
	// c.pullImages(ctx, cli)
	image := chainCfg.Images[0]

	newFullNodes := make(IconNodes, c.numFullNodes)
	copy(newFullNodes, c.FullNodes)

	eg, egCtx := errgroup.WithContext(ctx)
	for i := len(c.FullNodes); i < c.numFullNodes; i++ {
		i := i
		eg.Go(func() error {
			fn, err := c.NewChainNode(egCtx, testName, cli, networkID, image, false)
			if err != nil {
				return err
			}
			fn.Index = i
			newFullNodes[i] = fn
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	c.findTxMu.Lock()
	defer c.findTxMu.Unlock()
	c.FullNodes = newFullNodes
	return nil
}

func (c *IconLocalnet) pullImages(ctx context.Context, cli *client.Client) {
	for _, image := range c.Config().Images {
		rc, err := cli.ImagePull(
			ctx,
			image.Repository+":"+image.Version,
			dockertypes.ImagePullOptions{},
		)
		if err != nil {
			c.log.Error("Failed to pull image",
				zap.Error(err),
				zap.String("repository", image.Repository),
				zap.String("tag", image.Version),
			)
		} else {
			_, _ = io.Copy(io.Discard, rc)
			_ = rc.Close()
		}
	}
}

func (c *IconLocalnet) NewChainNode(
	ctx context.Context,
	testName string,
	cli *client.Client,
	networkID string,
	image ibc.DockerImage,
	validator bool,
) (*IconNode, error) {
	// Construct the ChainNode first so we can access its name.
	// The ChainNode's VolumeName cannot be set until after we create the volume.
	in := &IconNode{
		log:          c.log,
		Chain:        c,
		DockerClient: cli,
		NetworkID:    networkID,
		TestName:     testName,
		Image:        image,
	}

	v, err := cli.VolumeCreate(ctx, volumetypes.VolumeCreateBody{
		Labels: map[string]string{
			dockerutil.CleanupLabel: testName,

			dockerutil.NodeOwnerLabel: in.Name(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating volume for chain node: %w", err)
	}
	in.VolumeName = v.Name

	if err := dockerutil.SetVolumeOwner(ctx, dockerutil.VolumeOwnerOptions{
		Log: c.log,

		Client: cli,

		VolumeName: v.Name,
		ImageRef:   image.Ref(),
		TestName:   testName,
		UidGid:     image.UidGid,
	}); err != nil {
		return nil, fmt.Errorf("set volume owner: %w", err)
	}
	return in, nil
}

// Start sets up everything needed (validators, gentx, fullnodes, peering, additional accounts) for chain to start from genesis.
func (c *IconLocalnet) Start(testName string, ctx context.Context, additionalGenesisWallets ...ibc.WalletAmount) error {
	c.findTxMu.Lock()
	defer c.findTxMu.Unlock()
	eg, egCtx := errgroup.WithContext(ctx)
	for _, n := range c.FullNodes {
		n := n
		eg.Go(func() error {
			if err := n.CreateNodeContainer(egCtx, additionalGenesisWallets...); err != nil {
				return err
			}
			// All (validators, gentx, fullnodes, peering, additional accounts) are included in the image itself.
			return n.StartContainer(ctx)
		})
	}
	return eg.Wait()
}

// Exec runs an arbitrary command using Chain's docker environment.
// Whether the invoked command is run in a one-off container or execing into an already running container
// is up to the chain implementation.
//
// "env" are environment variables in the format "MY_ENV_VAR=value"
func (c *IconLocalnet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	return c.getFullNode().Exec(ctx, cmd, env)
}

// ExportState exports the chain state at specific height.
func (c *IconLocalnet) ExportState(ctx context.Context, height int64) (string, error) {
	block, err := c.getFullNode().GetBlockByHeight(ctx, height)
	return block, err
}

// GetRPCAddress retrieves the rpc address that can be reached by other containers in the docker network.
func (c *IconLocalnet) GetRPCAddress() string {
	return fmt.Sprintf("http://%s:9080/api/v3/", c.getFullNode().HostName())
}

func (c *IconLocalnet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	return c.FullNodes[0].GetChainConfig(ctx, rlyHome, keyName)
}

// GetGRPCAddress retrieves the grpc address that can be reached by other containers in the docker network.
// Not Applicable for Icon
func (c *IconLocalnet) GetGRPCAddress() string {
	return ""
}

// GetHostRPCAddress returns the rpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
func (c *IconLocalnet) GetHostRPCAddress() string {
	return "http://" + c.getFullNode().HostRPCPort + "/api/v3"
}

// GetHostGRPCAddress returns the grpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
// Not applicable for Icon
func (c *IconLocalnet) GetHostGRPCAddress() string {
	return ""
}

// HomeDir is the home directory of a node running in a docker container. Therefore, this maps to
// the container's filesystem (not the host).
func (c *IconLocalnet) HomeDir() string {
	return c.getFullNode().HomeDir()
}

func (c *IconLocalnet) createKeystore(ctx context.Context, keyName string) (string, string, error) {
	w := wallet.New()
	ks, err := wallet.KeyStoreFromWallet(w, []byte(keyName))
	if err != nil {
		return "", "", err
	}

	err = c.getFullNode().RestoreKeystore(ctx, ks, keyName)
	if err != nil {
		c.log.Error("fail to restore keystore", zap.Error(err))
		return "", "", err
	}
	ksd, err := wallet.NewKeyStoreData(ks)
	if err != nil {
		return "", "", err
	}
	key, err := wallet.DecryptICONKeyStore(ksd, []byte(keyName))
	if err != nil {
		return "", "", err
	}
	return w.Address(), hex.EncodeToString(key.Bytes()), nil
}

// RecoverKey recovers an existing user from a given mnemonic.
func (c *IconLocalnet) RecoverKey(ctx context.Context, name string, mnemonic string) error {
	panic("not implemented") // TODO: Implement
}

// GetAddress fetches the bech32 address for a test key on the "user" node (either the first fullnode or the first validator if no fullnodes).
func (c *IconLocalnet) GetAddress(ctx context.Context, keyName string) ([]byte, error) {
	addrInByte, err := json.Marshal(keyName)
	if err != nil {
		return nil, err
	}
	return addrInByte, nil
}

// SendFunds sends funds to a wallet from a user account.
func (c *IconLocalnet) SendFunds(ctx context.Context, keyName string, amount ibc.WalletAmount) error {
	c.CheckForKeyStore(ctx, keyName)

	cmd := c.getFullNode().NodeCommand("rpc", "sendtx", "transfer", "--key_store", c.keystorePath, "--key_password", keyName,
		"--to", amount.Address, "--value", fmt.Sprint(amount.Amount)+"000000000000000000", "--step_limit", "10000000000000")
	_, _, err := c.getFullNode().Exec(ctx, cmd, nil)
	return err
}

// Height returns the current block height or an error if unable to get current height.
func (c *IconLocalnet) Height(ctx context.Context) (uint64, error) {
	return c.getFullNode().Height(ctx)
}

// GetGasFeesInNativeDenom gets the fees in native denom for an amount of spent gas.
func (c *IconLocalnet) GetGasFeesInNativeDenom(gasPaid int64) int64 {
	gasPrice, _ := strconv.ParseFloat(strings.Replace(c.cfg.GasPrices, c.cfg.Denom, "", 1), 64)
	fees := float64(gasPaid) * gasPrice
	return int64(fees)
}

// BuildRelayerWallet will return a chain-specific wallet populated with the mnemonic so that the wallet can
// be restored in the relayer node using the mnemonic. After it is built, that address is included in
// genesis with some funds.
func (c *IconLocalnet) BuildRelayerWallet(ctx context.Context, keyName string) (ibc.Wallet, error) {
	return c.BuildWallet(ctx, keyName, "")
}

func (c *IconLocalnet) BuildWallet(ctx context.Context, keyName string, mnemonic string) (ibc.Wallet, error) {
	address, privateKey, err := c.createKeystore(ctx, keyName)
	if err != nil {
		return nil, fmt.Errorf("failed to create key with name %q on chain %s: %w", keyName, c.cfg.Name, err)

	}

	w := NewWallet(keyName, []byte(address), privateKey, c.cfg)
	c.Wallets[keyName] = w
	return w, nil
}

func (c *IconLocalnet) getFullNode() *IconNode {
	c.findTxMu.Lock()
	defer c.findTxMu.Unlock()
	if len(c.FullNodes) > 0 {
		// use first full node
		return c.FullNodes[0]
	}
	return c.FullNodes[0]
}

func (c *IconLocalnet) FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error) {
	fn := c.getFullNode()
	return fn.FindTxs(ctx, height)
}

// GetBalance fetches the current balance for a specific account address and denom.
func (c *IconLocalnet) GetBalance(ctx context.Context, address string, denom string) (int64, error) {
	return c.getFullNode().GetBalance(ctx, address)
}

func (c *IconLocalnet) SetupConnection(ctx context.Context, keyName string, target chains.Chain) error {
	//testcase := ctx.Value("testcase").(string)
	xcall := c.IBCAddresses["xcall"]
	_ = c.CheckForKeyStore(ctx, keyName)
	relayerKey := fmt.Sprintf("relayer-%s", c.Config().Name)
	relayerAddress := c.Wallets[relayerKey].FormattedAddress()

	connection, err := c.getFullNode().DeployContract(ctx, c.scorePaths["connection"], c.keystorePath, `{"_xCall":"`+xcall+`","_relayer":"`+relayerAddress+`"}`)
	if err != nil {
		return err
	}

	params := `{"networkId":"` + target.Config().ChainID + `", "messageFee":"0x0", "responseFee":"0x0"}`
	ctx, err = c.executeContract(context.Background(), connection, relayerKey, "setFee", params)
	if err != nil {
		return err
	}
	c.IBCAddresses["connection"] = connection
	return nil
}

func (c *IconLocalnet) SetupXCall(ctx context.Context, keyName string) error {
	//testcase := ctx.Value("testcase").(string)
	nid := c.cfg.ChainID
	//ibcAddress := c.IBCAddresses["ibc"]
	_ = c.CheckForKeyStore(ctx, keyName)
	xcall, err := c.getFullNode().DeployContract(ctx, c.scorePaths["xcall"], c.keystorePath, `{"networkId":"`+nid+`"}`)
	if err != nil {
		return err
	}
	//

	//c.IBCAddresses[fmt.Sprintf("xcall-%s", testcase)] = xcall
	//c.IBCAddresses[fmt.Sprintf("connection-%s", testcase)] = connection
	c.IBCAddresses["xcall"] = xcall

	return nil
}

func (c *IconLocalnet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	testcase := ctx.Value("testcase").(string)
	c.CheckForKeyStore(ctx, keyName)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)

	xCall := c.IBCAddresses["xcall"]
	params := `{"_callService":"` + xCall + `"}`
	dapp, err := c.getFullNode().DeployContract(ctx, c.scorePaths["dapp"], c.keystorePath, params)
	if err != nil {
		return err
	}
	c.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dapp

	for _, connection := range connections {
		//connectionKey := fmt.Sprintf("%s-%s", connection.Connection, testcase)
		params = `{"nid":"` + connection.Nid + `", "source":"` + c.IBCAddresses[connection.Connection] + `", "destination":"` + connection.Destination + `"}`
		ctx, err = c.executeContract(context.Background(), dapp, keyName, "addConnection", params)
		if err != nil {
			c.log.Error("Unable to add connection",
				zap.Error(err),
				zap.String("nid", connection.Nid),
				zap.String("source", c.IBCAddresses[connection.Connection]),
				zap.String("destination", connection.Destination),
			)
		}
	}

	return nil
}

func (c *IconLocalnet) GetContractAddress(key string) string {
	value, exist := c.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}
	return value
}

func (c *IconLocalnet) BackupConfig() ([]byte, error) {
	wallets := make(map[string]interface{})
	for key, value := range c.Wallets {
		wallets[key] = map[string]string{
			"mnemonic":         value.Mnemonic(),
			"address":          hex.EncodeToString(value.Address()),
			"formattedAddress": value.FormattedAddress(),
		}
	}
	backup := map[string]interface{}{
		"addresses": c.IBCAddresses,
		"wallets":   wallets,
	}
	return json.MarshalIndent(backup, "", "\t")
}

func (c *IconLocalnet) RestoreConfig(backup []byte) error {
	result := make(map[string]interface{})
	err := json.Unmarshal(backup, &result)
	if err != nil {
		return err
	}
	c.IBCAddresses = result["addresses"].(map[string]string)
	wallets := make(map[string]ibc.Wallet)

	for key, value := range result["wallets"].(map[string]interface{}) {
		_value := value.(map[string]string)
		mnemonic := _value["mnemonic"]
		address, _ := hex.DecodeString(_value["address"])
		wallets[key] = NewWallet(key, address, mnemonic, c.Config())
	}
	c.Wallets = wallets
	return nil
}

func (c *IconLocalnet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	// TODO: send fees
	var params = `{"_to":"` + _to + `", "_data":"` + hex.EncodeToString(data) + `"}`
	if rollback != nil {
		params = `{"_to":"` + _to + `", "_data":"` + hex.EncodeToString(data) + `", "_rollback":"` + hex.EncodeToString(rollback) + `"}`
	}
	ctx, err := c.executeContract(ctx, c.IBCAddresses[dappKey], keyName, "sendMessage", params)
	if err != nil {
		return nil, err
	}
	txn := ctx.Value("txResult").(*icontypes.TransactionResult)

	return context.WithValue(ctx, "sn", getSn(txn)), nil
}

// HasPacketReceipt returns the receipt of the packet sent to the target chain
func (c *IconLocalnet) IsPacketReceived(ctx context.Context, params map[string]interface{}, order ibc.Order) bool {
	if order == ibc.Ordered {
		sequence := params["sequence"].(uint64) //2
		ctx, err := c.QueryContract(ctx, c.IBCAddresses["ibc"], chains.GetNextSequenceReceive, params)
		if err != nil {
			fmt.Printf("Error--%v\n", err)
			return false
		}
		response, err := formatHexNumberFromResponse(ctx.Value("query-result").([]byte))

		if err != nil {
			fmt.Printf("Error--%v\n", err)
			return false
		}
		fmt.Printf("response[\"data\"]----%v", response)
		return sequence < response
	}
	ctx, _ = c.QueryContract(ctx, c.IBCAddresses["ibc"], chains.HasPacketReceipt, params)

	response, err := formatHexNumberFromResponse(ctx.Value("query-result").([]byte))
	if err != nil {
		fmt.Printf("Error--%v\n", err)
		return false
	}
	return response == 1
}

func formatHexNumberFromResponse(value []byte) (uint64, error) {
	pattern := `0x[0-9a-fA-F]+`
	regex := regexp.MustCompile(pattern)
	result := regex.FindString(string(value))
	if result == "" {
		return 0, fmt.Errorf("number not found")

	}

	response, err := strconv.ParseInt(result, 0, 64)
	if err != nil {
		return 0, err
	}
	return uint64(response), nil
}

// FindTargetXCallMessage returns the request id and the data of the message sent to the target chain
func (c *IconLocalnet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	sn := ctx.Value("sn").(string)
	reqId, destData, err := target.FindCallMessage(ctx, height, c.cfg.ChainID+"/"+c.IBCAddresses[dappKey], to, sn)
	return &chains.XCallResponse{SerialNo: sn, RequestID: reqId, Data: destData}, err
}

func (c *IconLocalnet) XCall(ctx context.Context, targetChain chains.Chain, keyName, to string, data, rollback []byte) (*chains.XCallResponse, error) {
	height, err := targetChain.(ibc.Chain).Height(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: send fees
	ctx, err = c.SendPacketXCall(ctx, keyName, to, data, rollback)
	if err != nil {
		return nil, err
	}
	return c.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func getSn(tx *icontypes.TransactionResult) string {
	for _, log := range tx.EventLogs {
		if string(log.Indexed[0]) == "CallMessageSent(Address,str,int)" {
			sn, _ := strconv.ParseInt(log.Indexed[3], 0, 64)
			return strconv.FormatInt(sn, 10)
		}
	}
	return ""
}

func (c *IconLocalnet) ExecuteCall(ctx context.Context, reqId, data string) (context.Context, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	return c.executeContract(ctx, c.IBCAddresses["xcall"], interchaintest.UserAccount, "executeCall", `{"_reqId":"`+reqId+`","_data":"`+data+`"}`)
}

func (c *IconLocalnet) ExecuteRollback(ctx context.Context, sn string) (context.Context, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	ctx, err := c.executeContract(ctx, c.IBCAddresses["xcall"], interchaintest.UserAccount, "executeRollback", `{"_sn":"`+sn+`"}`)
	if err != nil {
		return nil, err
	}
	txn := ctx.Value("txResult").(*icontypes.TransactionResult)
	sequence, err := icontypes.HexInt(txn.EventLogs[0].Indexed[1]).Int()
	return context.WithValue(ctx, "IsRollbackEventFound", fmt.Sprintf("%d", sequence) == sn), nil

}

func (c *IconLocalnet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	index := []*string{&from, &to, &sn}
	event, err := c.FindEvent(ctx, startHeight, "xcall", "CallMessage(str,str,int,int,bytes)", index)
	if err != nil {
		return "", "", err
	}

	intHeight, _ := event.Height.Int()
	block, _ := c.getFullNode().Client.GetBlockByHeight(&icontypes.BlockHeightParam{Height: icontypes.NewHexInt(int64(intHeight - 1))})
	i, _ := event.Index.Int()
	tx := block.NormalTransactions[i]
	trResult, _ := c.getFullNode().TransactionResult(ctx, string(tx.TxHash))
	eventIndex, _ := event.Events[0].Int()
	reqId := trResult.EventLogs[eventIndex].Data[0]
	data := trResult.EventLogs[eventIndex].Data[1]
	return reqId, data, nil
}

func (c *IconLocalnet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	index := []*string{&sn}
	event, err := c.FindEvent(ctx, startHeight, "xcall", "ResponseMessage(int,int)", index)
	if err != nil {
		return "", err
	}

	intHeight, _ := event.Height.Int()
	block, _ := c.getFullNode().Client.GetBlockByHeight(&icontypes.BlockHeightParam{Height: icontypes.NewHexInt(int64(intHeight - 1))})
	i, _ := event.Index.Int()
	tx := block.NormalTransactions[i]
	trResult, _ := c.getFullNode().TransactionResult(ctx, string(tx.TxHash))
	eventIndex, _ := event.Events[0].Int()
	code, _ := strconv.ParseInt(trResult.EventLogs[eventIndex].Data[0], 0, 64)

	return strconv.FormatInt(code, 10), nil
}

func (c *IconLocalnet) FindEvent(ctx context.Context, startHeight uint64, contract, signature string, index []*string) (*icontypes.EventNotification, error) {
	filter := icontypes.EventFilter{
		Addr:      icontypes.Address(c.IBCAddresses[contract]),
		Signature: signature,
		Indexed:   index,
	}

	// Create a context with a timeout of 16 seconds.
	_ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create an event request with the given filter and start height.
	req := &icontypes.EventRequest{
		EventFilter: filter,
		Height:      icontypes.NewHexInt(int64(startHeight)),
	}
	channel := make(chan *icontypes.EventNotification)
	response := func(_ *websocket.Conn, v *icontypes.EventNotification) error {
		channel <- v
		return nil
	}
	errRespose := func(conn *websocket.Conn, err error) {}
	go func(ctx context.Context, req *icontypes.EventRequest, response func(*websocket.Conn, *icontypes.EventNotification) error, errRespose func(*websocket.Conn, error)) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered: %v", err)
			}
		}()
		if err := c.getFullNode().Client.MonitorEvent(ctx, req, response, errRespose); err != nil {
			log.Printf("MonitorEvent error: %v", err)
		}
	}(ctx, req, response, errRespose)

	select {
	case v := <-channel:
		return v, nil
	case <-_ctx.Done():
		return nil, errors.New(fmt.Sprintf("timeout : Event %s not found after %d block", signature, startHeight))
	}
}

// DeployContract implements chains.Chain
func (c *IconLocalnet) DeployContract(ctx context.Context, keyName string) (context.Context, error) {
	// Get contract Name from context
	ctxValue := ctx.Value(chains.ContractName{}).(chains.ContractName)
	contractName := ctxValue.ContractName

	// Get Init Message from context
	ctxVal := ctx.Value(chains.InitMessageKey("init-msg")).(chains.InitMessage)

	initMessage := c.getInitParams(ctx, contractName, ctxVal.Message)

	var contracts chains.ContractKey

	// Check if keystore is alreadry available for given keyName
	ownerAddr := c.CheckForKeyStore(ctx, keyName)
	if ownerAddr != nil {
		contracts.ContractOwner = map[string]string{
			keyName: ownerAddr.FormattedAddress(),
		}
	}

	// Get ScoreAddress
	scoreAddress, err := c.getFullNode().DeployContract(ctx, c.scorePaths[contractName], c.keystorePath, initMessage)

	contracts.ContractAddress = map[string]string{
		contractName: scoreAddress,
	}

	testcase := ctx.Value("testcase").(string)
	contract := fmt.Sprintf("%s-%s", contractName, testcase)
	c.IBCAddresses[contract] = scoreAddress
	return context.WithValue(ctx, chains.Mykey("contract Names"), chains.ContractKey{
		ContractAddress: contracts.ContractAddress,
		ContractOwner:   contracts.ContractOwner,
	}), err
}

// executeContract implements chains.Chain
func (c *IconLocalnet) executeContract(ctx context.Context, contractAddress, keyName, methodName, params string) (context.Context, error) {
	c.CheckForKeyStore(ctx, keyName)

	hash, err := c.getFullNode().ExecuteContract(ctx, contractAddress, methodName, c.keystorePath, params)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Transaction Hash: %s\n", hash)

	txHashByte, err := hex.DecodeString(strings.TrimPrefix(hash, "0x"))
	if err != nil {
		return nil, fmt.Errorf("error when executing contract %v ", err)
	}
	_, res, err := c.getFullNode().Client.WaitForResults(ctx, &icontypes.TransactionHashParam{Hash: icontypes.NewHexBytes(txHashByte)})
	if err != nil {
		return nil, err
	}
	if res.Status == "0x1" {
		return context.WithValue(ctx, "txResult", res), nil
	}
	//TODO add debug flag to print trace
	trace, err := c.getFullNode().GetDebugTrace(ctx, icontypes.NewHexBytes(txHashByte))
	if err == nil {
		logs, _ := json.Marshal(trace.Logs)
		fmt.Printf("---------debug trace start-----------\n%s\n---------debug trace end-----------\n", string(logs))
	}
	return ctx, fmt.Errorf("%s", res.Failure.MessageValue)
}

func (c *IconLocalnet) ExecuteContract(ctx context.Context, contractAddress, keyName, methodName string, params map[string]interface{}) (context.Context, error) {
	execMethodName, execParams := c.getExecuteParam(ctx, methodName, params)
	return c.executeContract(ctx, contractAddress, keyName, execMethodName, execParams)
}

// GetBlockByHeight implements chains.Chain
func (c *IconLocalnet) GetBlockByHeight(ctx context.Context) (context.Context, error) {
	panic("unimplemented")
}

// GetLastBlock implements chains.Chain
func (c *IconLocalnet) GetLastBlock(ctx context.Context) (context.Context, error) {
	h, err := c.getFullNode().Height(ctx)
	return context.WithValue(ctx, chains.LastBlock{}, h), err
}

func (c *IconLocalnet) InitEventListener(ctx context.Context, contract string) chains.EventListener {
	listener := NewIconEventListener(c, contract)
	return listener
}

// QueryContract implements chains.Chain
func (c *IconLocalnet) QueryContract(ctx context.Context, contractAddress, methodName string, params map[string]interface{}) (context.Context, error) {
	time.Sleep(2 * time.Second)

	// get query msg
	query := c.GetQueryParam(methodName, params)
	_params, _ := json.Marshal(query.Value)
	output, err := c.getFullNode().QueryContract(ctx, contractAddress, query.MethodName, string(_params))

	chains.Response = output
	fmt.Printf("Response is : %s \n", output)
	return context.WithValue(ctx, "query-result", chains.Response), err

}

func (c *IconLocalnet) BuildWallets(ctx context.Context, keyName string) (ibc.Wallet, error) {
	w := c.CheckForKeyStore(ctx, keyName)
	if w == nil {
		return nil, fmt.Errorf("error keyName already exists")
	}

	amount := ibc.WalletAmount{
		Address: w.FormattedAddress(),
		Amount:  10000,
	}
	var err error

	err = c.SendFunds(ctx, interchaintest.FaucetAccountKeyName, amount)
	return w, err
}

// PauseNode pauses the node
func (c *IconLocalnet) PauseNode(ctx context.Context) error {
	return c.getFullNode().DockerClient.ContainerPause(ctx, c.getFullNode().ContainerID)
}

// UnpauseNode starts the paused node
func (c *IconLocalnet) UnpauseNode(ctx context.Context) error {
	return c.getFullNode().DockerClient.ContainerUnpause(ctx, c.getFullNode().ContainerID)
}
