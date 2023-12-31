package chains

import (
	"fmt"
	"strings"
)

const (
	BASE_PATH = "BASE_PATH"
)

type ContractKey struct {
	ContractAddress map[string]string
	ContractOwner   map[string]string
}

type XCallConnection struct {
	Connection    string
	Nid           string
	Destination   string
	TimeoutHeight int `default:"100"`
}

type XCallResponse struct {
	SerialNo  string
	RequestID string
	Data      string
}

type LastBlock struct{}

type ContractName struct {
	ContractName string
}

type InitMessage struct {
	Message map[string]interface{}
}

type InitMessageKey string

type Param struct {
	Data string
}

type Query struct {
	Method string
	Param  string
}

type Mykey string

type Admins struct {
	Admin map[string]string
}

type AdminKey string

var Response interface{}

type MinimumGasPriceEntity struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Event map[string][]string
type Filter map[string]interface{}

type EventListener interface {
	Start()
	Stop()
	FindEvent(filters Filter) (Event, error)
}

type TimeoutResponse struct {
	HasTimeout        bool
	IsPacketFound     bool
	HasRollbackCalled bool
}

type PacketTransferResponse struct {
	IsPacketSent              bool
	IsPacketReceiptEventFound bool
}

type BufferArray []byte

func (u BufferArray) MarshalJSON() ([]byte, error) {
	if u == nil {
		return []byte{}, nil
	}

	return []byte(strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")), nil
}
