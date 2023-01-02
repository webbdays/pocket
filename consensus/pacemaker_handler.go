package consensus

import (
	"encoding/binary"
	"fmt"
	"log"

	typesCons "github.com/pokt-network/pocket/consensus/types"
	"github.com/pokt-network/pocket/shared/codec"
	"google.golang.org/protobuf/types/known/anypb"
)

func (m *consensusModule) HandlePacemakerMessage(pacemakerMessage *anypb.Any) error {
	m.m.Lock()
	defer m.m.Unlock()

	switch pacemakerMessage.MessageName() {
	case PacemakerAccessContentType:

		msg, err := codec.GetCodec().FromAny(pacemakerMessage)
		if err != nil {
			return err
		}

		paceMakerMessage, ok := msg.(*typesCons.PacemakerMessage)
		if !ok {
			return fmt.Errorf("failed to cast message to HotstuffMessage")
		}

		if err := m.handlePaceMakerMessage(paceMakerMessage); err != nil {
			return err
		}

	default:
		return typesCons.ErrUnknownPacemakerMessageType(pacemakerMessage.MessageName())
	}

	return nil
}

func (m *consensusModule) handlePaceMakerMessage(pacemakerMsg *typesCons.PacemakerMessage) error {
	switch pacemakerMsg.Action {
	case typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_SET_HEIGHT:
		log.Printf("\n\nPacemaker SetHeight Called Via Bus\n\n")
		m.SetHeight(binary.BigEndian.Uint64(pacemakerMsg.Message.Value))

	// case messaging.DebugMessageAction_DEBUG_CONSENSUS_PRINT_NODE_STATE:
	// 	m.printNodeState(debugMessage)
	default:
		log.Println("Debug message: \n", &pacemakerMsg.Message)
	}
	return nil
}
