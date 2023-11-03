package icon

import (
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func parseMessagesFromEventlogs(log *zap.Logger, eventlogs []types.EventLog, height uint64) []providerTypes.Message {
	msgs := make([]providerTypes.Message, 0)
	for _, el := range eventlogs {
		message, ok := parseMessageFromEvent(log, el, height)
		if ok {
			msgs = append(msgs, message)
		}
	}
	return msgs
}

func parseMessageFromEvent(
	log *zap.Logger,
	event types.EventLog,
	height uint64,
) (providerTypes.Message, bool) {
	eventName := string(event.Indexed[0][:])
	eventType := EventNameToType[eventName]

	switch eventName {
	case EmitMessage:
		// TODO: fetch message from eventlog
		return providerTypes.Message{
			MessageHeight: height,
			EventType:     eventType,
		}, true
	}
	return providerTypes.Message{}, false
}