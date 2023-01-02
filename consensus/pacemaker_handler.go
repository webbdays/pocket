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
		//log.Printf("\n pacemakerMsg SetHeight Called Via Bus, Height is: %d", pacemakerMsg.GetHeight().Height)
		m.height = pacemakerMsg.GetHeight().Height
		//log.Printf(" pacemaker_handler New Height is: %d \n\n\n", m.height)
	case typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_SET_ROUND:
		//log.Printf("\n pacemaker_handler SetRound Called Via Bus, Round is: %d", pacemakerMsg.GetRound().Round)
		m.round = pacemakerMsg.GetRound().Round
		//log.Printf("In pacemaker_handler New Round is: %d \n\n\n", m.round)
	case typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_RESET_FOR_NEW_HEIGHT:
		//log.Printf("\n pacemaker_handler Reset for new round Called Via Bus")
		m.resetForNewHeight()
	default:
		log.Printf("\n\n Unexpected case, message Action: %s \n\n", &pacemakerMsg.Action)
	}
	return nil
}
