package tests

import (
	"context"
	"errors"
	"fmt"
	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/testsuite"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

type XCallTestSuite struct {
	*testsuite.E2ETestSuite
	T *testing.T
}

func (x *XCallTestSuite) TextXCall() {
	testcase := "xcall"
	portId := "transfer"
	ctx := context.WithValue(context.TODO(), "testcase", testcase)
	//x.Require().NoError(x.SetupXCall(ctx), "fail to setup xcall")
	x.Require().NoError(x.DeployXCallMockApp(ctx, portId), "fail to deploy xcall demo dapp")
	chainA, chainB := x.GetChains()
	x.T.Run("xcall one way message chainA-chainB", func(t *testing.T) {
		err := x.testOneWayMessage(ctx, t, chainA, chainB)
		assert.NoErrorf(t, err, "fail xCall one way message chainA-chainB ::%v\n ", err)
	})

	x.T.Run("xcall one way message chainB-chainA", func(t *testing.T) {
		err := x.testOneWayMessage(ctx, t, chainB, chainA)
		assert.NoErrorf(t, err, "fail xCall one way message chainB-chainA ::%v\n ", err)
	})

	x.T.Run("xcall test rollback chainA-chainB", func(t *testing.T) {
		err := x.testRollback(ctx, t, chainA, chainB)
		assert.NoErrorf(t, err, "fail xCall rollback message chainB-chainA ::%v\n ", err)

	})

	x.T.Run("xcall test rollback chainB-chainA", func(t *testing.T) {
		err := x.testRollback(ctx, t, chainB, chainA)
		assert.NoErrorf(t, err, "fail xcCll rollback message chainB-chainA ::%v\n ", err)

	})

	x.T.Run("xcall test send maxSize Data: 2048 bytes", func(t *testing.T) {
		x.testOneWayMessageWithSize(ctx, t, 1300, chainA, chainB)
		x.testOneWayMessageWithSize(ctx, t, 1300, chainB, chainA)
	})

	x.T.Run("xcall test send maxSize Data: 2049bytes", func(t *testing.T) {
		x.testOneWayMessageWithSizeExpectingError(ctx, t, 2000, chainB, chainA)
		x.testOneWayMessageWithSizeExpectingError(ctx, t, 2100, chainA, chainB)
	})

}

func (x *XCallTestSuite) TestXCallPacketDrop() {
	testcase := "packet-drop"
	portId := "transfer-1"
	ctx := context.WithValue(context.TODO(), "testcase", testcase)
	//x.Require().NoError(x.SetupXCall(ctx), "fail to setup xcall")
	x.Require().NoError(x.DeployXCallMockApp(ctx, portId), "fail to deploy xcall dapp")
	chainA, chainB := x.GetChains()
	x.T.Run("xcall packet drop chainA-chainB", func(t *testing.T) {
		x.testPacketDrop(ctx, t, chainA, chainB)
	})

	x.T.Run("xcall packet drop chainB-chainA", func(t *testing.T) {
		x.testPacketDrop(ctx, t, chainB, chainA)
	})
}

func (x *XCallTestSuite) testPacketDrop(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "drop-msg"
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	//height, _internal := chainA.(ibc.Chain).Height(ctx)
	listener := chainA.InitEventListener(ctx, "ibc")
	defer listener.Stop()
	res, err := chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, []byte(msg), []byte("rollback-data"))
	assert.Errorf(t, err, "failed to find eventlog - %w", err)

	sn := res.SerialNo
	snInt, _ := strconv.Atoi(res.SerialNo)
	params := map[string]interface{}{
		"port_id":    "transfer-1",
		"channel_id": "channel-0",
		"sequence":   uint64(snInt),
	}

	ctx, err = chainA.CheckForTimeout(ctx, chainB, params, listener)
	response := ctx.Value("timeout-response").(*chains.TimeoutResponse)
	assert.Truef(t, response.HasTimeout, "timeout event not found - %s", sn)
	assert.Falsef(t, response.IsPacketFound, "packet found on target chain - %s", sn)
	assert.Truef(t, response.HasRollbackCalled, "failed to call rollback  - %s", sn)
}

func (x *XCallTestSuite) testOneWayMessage(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "MessageTransferTestingWithoutRollback"
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)

	res, err := chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, []byte(msg), nil)
	result := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !result {
		return err
	}
	ctx, err = chainB.ExecuteCall(ctx, res.RequestID, res.Data)
	result = assert.NoErrorf(t, err, "error on execute call packet- %v", err)
	if !result {
		return err
	}
	//x.Require().NoErrorf(err, "error on execute call packet- %v", err)
	dataOutput := x.ConvertToPlainString(res.Data)
	//x.Require().NoErrorf(err, "error on converting res data as msg- %v", err)
	result = assert.NoErrorf(t, err, "error on converting res data as msg- %v", err)
	if !result {
		return err
	}
	result = assert.Equal(t, msg, dataOutput)
	if !result {
		return err
	}
	fmt.Println("Data Transfer Testing Without Rollback from " + chainA.(ibc.Chain).Config().ChainID + " to " + chainB.(ibc.Chain).Config().ChainID + " with data " + msg + " and Received:" + dataOutput + " PASSED")
	return nil
}

func (x *XCallTestSuite) testRollback(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "rollback"
	rollback := "RollbackDataTesting"
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	res, err := chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, []byte(msg), []byte(rollback))
	isSuccess := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !isSuccess {
		return err
	}
	height, err := chainA.(ibc.Chain).Height(ctx)
	_, err = chainB.ExecuteCall(ctx, res.RequestID, res.Data)
	code, err := chainA.FindCallResponse(ctx, height, res.SerialNo)
	isSuccess = assert.NoErrorf(t, err, "no call response found %v", err)
	isSuccess = assert.Equal(t, "0", code)
	if !isSuccess {
		return err
	}
	ctx, err = chainA.ExecuteRollback(ctx, res.SerialNo)
	assert.NoErrorf(t, err, "error on excute rollback- %w", err)
	return err
}

func (x *XCallTestSuite) testOneWayMessageWithSize(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	res, err := chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, _msg, nil)
	assert.NoError(t, err)

	_, err = chainB.ExecuteCall(ctx, res.RequestID, res.Data)
	assert.NoError(t, err)
}

func (x *XCallTestSuite) testOneWayMessageWithSizeExpectingError(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	_, err := chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, _msg, nil)
	result := assert.Errorf(t, err, "large data transfer should failed")
	if result {
		result = false
		if strings.Contains(err.Error(), "submessages:") {
			subStart := strings.Index(err.Error(), "submessages:") + len("submessages:")
			subEnd := strings.Index(err.Error(), ": execute")
			subMsg := err.Error()[subStart:subEnd]
			result = assert.ObjectsAreEqual(strings.TrimSpace(subMsg), "MaxDataSizeExceeded")
		} else if strings.Contains(err.Error(), "MaxDataSizeExceeded") {
			result = true
		} else {
			result = assert.ObjectsAreEqual(errors.New("UnknownFailure"), err)
		}
		if result {
			t.Logf("Test passed: %v", err)
		} else {
			t.Errorf("Test failed: %v", err)
		}
	}

}
