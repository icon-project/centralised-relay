// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// AbiMetaData contains all meta data concerning the Abi contract.
var AbiMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"targetNetwork\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"Message\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"connSn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"response\",\"type\":\"bool\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"srcNetwork\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_connSn\",\"type\":\"uint256\"}],\"name\":\"getReceipt\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_relayer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_xCall\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"srcNetwork\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_connSn\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"recvMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"sn\",\"type\":\"uint256\"}],\"name\":\"revertMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"svc\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"sn\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"sendMessage\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"networkId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"messageFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"responseFee\",\"type\":\"uint256\"}],\"name\":\"setFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AbiABI is the input ABI used to generate the binding from.
// Deprecated: Use AbiMetaData.ABI instead.
var AbiABI = AbiMetaData.ABI

// Abi is an auto generated Go binding around an Ethereum contract.
type Abi struct {
	AbiCaller     // Read-only binding to the contract
	AbiTransactor // Write-only binding to the contract
	AbiFilterer   // Log filterer for contract events
}

// AbiCaller is an auto generated read-only Go binding around an Ethereum contract.
type AbiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AbiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AbiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AbiSession struct {
	Contract     *Abi              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AbiCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AbiCallerSession struct {
	Contract *AbiCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AbiTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AbiTransactorSession struct {
	Contract     *AbiTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AbiRaw is an auto generated low-level Go binding around an Ethereum contract.
type AbiRaw struct {
	Contract *Abi // Generic contract binding to access the raw methods on
}

// AbiCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AbiCallerRaw struct {
	Contract *AbiCaller // Generic read-only contract binding to access the raw methods on
}

// AbiTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AbiTransactorRaw struct {
	Contract *AbiTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAbi creates a new instance of Abi, bound to a specific deployed contract.
func NewAbi(address common.Address, backend bind.ContractBackend) (*Abi, error) {
	contract, err := bindAbi(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Abi{AbiCaller: AbiCaller{contract: contract}, AbiTransactor: AbiTransactor{contract: contract}, AbiFilterer: AbiFilterer{contract: contract}}, nil
}

// NewAbiCaller creates a new read-only instance of Abi, bound to a specific deployed contract.
func NewAbiCaller(address common.Address, caller bind.ContractCaller) (*AbiCaller, error) {
	contract, err := bindAbi(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AbiCaller{contract: contract}, nil
}

// NewAbiTransactor creates a new write-only instance of Abi, bound to a specific deployed contract.
func NewAbiTransactor(address common.Address, transactor bind.ContractTransactor) (*AbiTransactor, error) {
	contract, err := bindAbi(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AbiTransactor{contract: contract}, nil
}

// NewAbiFilterer creates a new log filterer instance of Abi, bound to a specific deployed contract.
func NewAbiFilterer(address common.Address, filterer bind.ContractFilterer) (*AbiFilterer, error) {
	contract, err := bindAbi(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AbiFilterer{contract: contract}, nil
}

// bindAbi binds a generic wrapper to an already deployed contract.
func bindAbi(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AbiMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Abi *AbiRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Abi.Contract.AbiCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Abi *AbiRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abi.Contract.AbiTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Abi *AbiRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Abi.Contract.AbiTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Abi *AbiCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Abi.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Abi *AbiTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abi.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Abi *AbiTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Abi.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Abi *AbiCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Abi *AbiSession) Admin() (common.Address, error) {
	return _Abi.Contract.Admin(&_Abi.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Abi *AbiCallerSession) Admin() (common.Address, error) {
	return _Abi.Contract.Admin(&_Abi.CallOpts)
}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_Abi *AbiCaller) ConnSn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "connSn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_Abi *AbiSession) ConnSn() (*big.Int, error) {
	return _Abi.Contract.ConnSn(&_Abi.CallOpts)
}

// ConnSn is a free data retrieval call binding the contract method 0x99f1fca7.
//
// Solidity: function connSn() view returns(uint256)
func (_Abi *AbiCallerSession) ConnSn() (*big.Int, error) {
	return _Abi.Contract.ConnSn(&_Abi.CallOpts)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_Abi *AbiCaller) GetFee(opts *bind.CallOpts, to string, response bool) (*big.Int, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "getFee", to, response)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_Abi *AbiSession) GetFee(to string, response bool) (*big.Int, error) {
	return _Abi.Contract.GetFee(&_Abi.CallOpts, to, response)
}

// GetFee is a free data retrieval call binding the contract method 0x7d4c4f4a.
//
// Solidity: function getFee(string to, bool response) view returns(uint256 fee)
func (_Abi *AbiCallerSession) GetFee(to string, response bool) (*big.Int, error) {
	return _Abi.Contract.GetFee(&_Abi.CallOpts, to, response)
}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_Abi *AbiCaller) GetReceipt(opts *bind.CallOpts, srcNetwork string, _connSn *big.Int) (bool, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "getReceipt", srcNetwork, _connSn)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_Abi *AbiSession) GetReceipt(srcNetwork string, _connSn *big.Int) (bool, error) {
	return _Abi.Contract.GetReceipt(&_Abi.CallOpts, srcNetwork, _connSn)
}

// GetReceipt is a free data retrieval call binding the contract method 0x9664da0e.
//
// Solidity: function getReceipt(string srcNetwork, uint256 _connSn) view returns(bool)
func (_Abi *AbiCallerSession) GetReceipt(srcNetwork string, _connSn *big.Int) (bool, error) {
	return _Abi.Contract.GetReceipt(&_Abi.CallOpts, srcNetwork, _connSn)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_Abi *AbiTransactor) ClaimFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "claimFees")
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_Abi *AbiSession) ClaimFees() (*types.Transaction, error) {
	return _Abi.Contract.ClaimFees(&_Abi.TransactOpts)
}

// ClaimFees is a paid mutator transaction binding the contract method 0xd294f093.
//
// Solidity: function claimFees() returns()
func (_Abi *AbiTransactorSession) ClaimFees() (*types.Transaction, error) {
	return _Abi.Contract.ClaimFees(&_Abi.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_Abi *AbiTransactor) Initialize(opts *bind.TransactOpts, _relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "initialize", _relayer, _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_Abi *AbiSession) Initialize(_relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _Abi.Contract.Initialize(&_Abi.TransactOpts, _relayer, _xCall)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _relayer, address _xCall) returns()
func (_Abi *AbiTransactorSession) Initialize(_relayer common.Address, _xCall common.Address) (*types.Transaction, error) {
	return _Abi.Contract.Initialize(&_Abi.TransactOpts, _relayer, _xCall)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_Abi *AbiTransactor) RecvMessage(opts *bind.TransactOpts, srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "recvMessage", srcNetwork, _connSn, _msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_Abi *AbiSession) RecvMessage(srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Abi.Contract.RecvMessage(&_Abi.TransactOpts, srcNetwork, _connSn, _msg)
}

// RecvMessage is a paid mutator transaction binding the contract method 0xb58b4cec.
//
// Solidity: function recvMessage(string srcNetwork, uint256 _connSn, bytes _msg) returns()
func (_Abi *AbiTransactorSession) RecvMessage(srcNetwork string, _connSn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Abi.Contract.RecvMessage(&_Abi.TransactOpts, srcNetwork, _connSn, _msg)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_Abi *AbiTransactor) RevertMessage(opts *bind.TransactOpts, sn *big.Int) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "revertMessage", sn)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_Abi *AbiSession) RevertMessage(sn *big.Int) (*types.Transaction, error) {
	return _Abi.Contract.RevertMessage(&_Abi.TransactOpts, sn)
}

// RevertMessage is a paid mutator transaction binding the contract method 0x2d3fb823.
//
// Solidity: function revertMessage(uint256 sn) returns()
func (_Abi *AbiTransactorSession) RevertMessage(sn *big.Int) (*types.Transaction, error) {
	return _Abi.Contract.RevertMessage(&_Abi.TransactOpts, sn)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_Abi *AbiTransactor) SendMessage(opts *bind.TransactOpts, to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "sendMessage", to, svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_Abi *AbiSession) SendMessage(to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Abi.Contract.SendMessage(&_Abi.TransactOpts, to, svc, sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0x522a901e.
//
// Solidity: function sendMessage(string to, string svc, int256 sn, bytes _msg) payable returns()
func (_Abi *AbiTransactorSession) SendMessage(to string, svc string, sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Abi.Contract.SendMessage(&_Abi.TransactOpts, to, svc, sn, _msg)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Abi *AbiTransactor) SetAdmin(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "setAdmin", _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Abi *AbiSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _Abi.Contract.SetAdmin(&_Abi.TransactOpts, _address)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _address) returns()
func (_Abi *AbiTransactorSession) SetAdmin(_address common.Address) (*types.Transaction, error) {
	return _Abi.Contract.SetAdmin(&_Abi.TransactOpts, _address)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_Abi *AbiTransactor) SetFee(opts *bind.TransactOpts, networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "setFee", networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_Abi *AbiSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _Abi.Contract.SetFee(&_Abi.TransactOpts, networkId, messageFee, responseFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x43f08a89.
//
// Solidity: function setFee(string networkId, uint256 messageFee, uint256 responseFee) returns()
func (_Abi *AbiTransactorSession) SetFee(networkId string, messageFee *big.Int, responseFee *big.Int) (*types.Transaction, error) {
	return _Abi.Contract.SetFee(&_Abi.TransactOpts, networkId, messageFee, responseFee)
}

// AbiInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Abi contract.
type AbiInitializedIterator struct {
	Event *AbiInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiInitialized represents a Initialized event raised by the Abi contract.
type AbiInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Abi *AbiFilterer) FilterInitialized(opts *bind.FilterOpts) (*AbiInitializedIterator, error) {

	logs, sub, err := _Abi.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AbiInitializedIterator{contract: _Abi.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Abi *AbiFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AbiInitialized) (event.Subscription, error) {

	logs, sub, err := _Abi.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiInitialized)
				if err := _Abi.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Abi *AbiFilterer) ParseInitialized(log types.Log) (*AbiInitialized, error) {
	event := new(AbiInitialized)
	if err := _Abi.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiMessageIterator is returned from FilterMessage and is used to iterate over the raw logs and unpacked data for Message events raised by the Abi contract.
type AbiMessageIterator struct {
	Event *AbiMessage // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiMessage)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiMessage)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiMessage represents a Message event raised by the Abi contract.
type AbiMessage struct {
	TargetNetwork string
	Sn            *big.Int
	Msg           []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMessage is a free log retrieval operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_Abi *AbiFilterer) FilterMessage(opts *bind.FilterOpts) (*AbiMessageIterator, error) {

	logs, sub, err := _Abi.contract.FilterLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return &AbiMessageIterator{contract: _Abi.contract, event: "Message", logs: logs, sub: sub}, nil
}

// WatchMessage is a free log subscription operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_Abi *AbiFilterer) WatchMessage(opts *bind.WatchOpts, sink chan<- *AbiMessage) (event.Subscription, error) {

	logs, sub, err := _Abi.contract.WatchLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiMessage)
				if err := _Abi.contract.UnpackLog(event, "Message", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMessage is a log parse operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string targetNetwork, uint256 sn, bytes _msg)
func (_Abi *AbiFilterer) ParseMessage(log types.Log) (*AbiMessage, error) {
	event := new(AbiMessage)
	if err := _Abi.contract.UnpackLog(event, "Message", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
