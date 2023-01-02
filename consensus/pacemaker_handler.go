package consensus

import (
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
		m.height = pacemakerMsg.GetHeight().Height
		log.Printf("handlePaceMakerMessage Height is: %d", m.height)
	case typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_SET_ROUND:
		m.round = pacemakerMsg.GetRound().Round
	case typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_SET_STEP:
		m.step = typesCons.HotstuffStep(pacemakerMsg.GetStep().Step)
	case typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_RESET_FOR_NEW_HEIGHT:
		m.resetForNewHeight()
	case typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_CLEAR_LEADER_MESSAGE_POOL:
		m.clearLeader()
		m.clearMessagesPool()
	case typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_RELEASE_UTILITY_CONTEXT:
		if m.utilityContext != nil {
			if err := m.utilityContext.Release(); err != nil {
				log.Println("[WARN] Failed to release utility context: ", err)
				return err
			}
			m.utilityContext = nil
		}
	case typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_BROADCAST_HOTSTUFF_MESSAGE_TO_NODES:
		m.broadcastToNodes(pacemakerMsg.GetMessage())
	default:
		log.Printf("\n\n Unexpected case, message Action: %s \n\n", &pacemakerMsg.Action)
	}
	return nil
}
