package socket

import (
	"net"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/store"
	"github.com/icon-project/centralized-relay/relayer/types"
)

type Event string

type Message struct {
	Event Event
	Data  []byte
}

type Server struct {
	listener net.Listener
	rly      *relayer.Relayer
}

type ReqMessageList struct {
	Chain      string
	Pagination *store.Pagination
}

type ReqGetBlock struct {
	Chain string
	All   bool
}

type ReqRelayMessage struct {
	Chain  string
	Sn     uint64
	Height uint64
}

type ReqMessageRemove struct {
	Chain string
	Sn    uint64
}

type ResMessageRemove struct {
	Sn     uint64
	Chain  string
	Dst    string
	Height uint64
	Event  string
}

type ResMessageList struct {
	Messages []*types.RouteMessage
	Total    int
}

type ResGetBlock struct {
	Chain  string
	Height uint64
}

type ResRelayMessage struct {
	*types.RouteMessage
}

type ReqPruneDB struct {
	Chain string
}

type ResPruneDB struct {
	Status string
}
