package consensus

import (
	"context"
	"fmt"
	"log"
	"time"
	timePkg "time"

	consensusTelemetry "github.com/pokt-network/pocket/consensus/telemetry"
	typesCons "github.com/pokt-network/pocket/consensus/types"
	"github.com/pokt-network/pocket/shared/messaging"
	"github.com/pokt-network/pocket/shared/modules"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	pacemakerModuleName = "pacemaker"
)

type Pacemaker interface {
	modules.Module
	PacemakerDebug

	// TODO(olshansky): Rather than exposing the underlying `ConsensusModule` struct,
	// we could create a `ConsensusModuleDebug` interface that'll expose setters/getters
	// for the height/round/step/etc, and interface with the module that way.
	SetConsensusModule(module *consensusModule)

	ValidateMessage(message *typesCons.HotstuffMessage) error
	RestartTimer()
	NewHeight()
	InterruptRound()
}

var (
	_ modules.Module             = &paceMaker{}
	_ modules.ConfigurableModule = &paceMaker{}
	_ PacemakerDebug             = &paceMaker{}
	_ modules.PacemakerConfig    = &typesCons.PacemakerConfig{}
)

type paceMaker struct {
	bus modules.Bus

	// TODO(olshansky): The reason `pacemaker_*` files are not a sub-package under consensus
	// due to it's dependency on the underlying implementation of `ConsensusModule`. Think
	// through a way to decouple these. This could be fixed with reflection but that's not
	// a great idea in production code.
	consensusMod *consensusModule

	pacemakerCfg modules.PacemakerConfig

	stepCancelFunc context.CancelFunc

	// Only used for development and debugging.
	paceMakerDebug
}

func CreatePacemaker(runtimeMgr modules.RuntimeMgr) (modules.Module, error) {
	var m paceMaker
	return m.Create(runtimeMgr)
}

func (m *paceMaker) Create(runtimeMgr modules.RuntimeMgr) (modules.Module, error) {
	cfg := runtimeMgr.GetConfig()
	if err := m.ValidateConfig(cfg); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	pacemakerCfg := cfg.GetConsensusConfig().(HasPacemakerConfig).GetPacemakerConfig()

	return &paceMaker{
		bus:          nil,
		consensusMod: nil,

		pacemakerCfg: pacemakerCfg,

		stepCancelFunc: nil, // Only set on restarts

		paceMakerDebug: paceMakerDebug{
			manualMode:                pacemakerCfg.GetManual(),
			debugTimeBetweenStepsMsec: pacemakerCfg.GetDebugTimeBetweenStepsMsec(),
			quorumCertificate:         nil,
		},
	}, nil
}

func (p *paceMaker) Start() error {
	p.RestartTimer()
	return nil
}
func (p *paceMaker) Stop() error {
	return nil
}

func (p *paceMaker) GetModuleName() string {
	return pacemakerModuleName
}

func (m *paceMaker) SetBus(pocketBus modules.Bus) {
	m.bus = pocketBus
}

func (m *paceMaker) GetBus() modules.Bus {
	if m.bus == nil {
		log.Fatalf("PocketBus is not initialized")
	}
	return m.bus
}

func (*paceMaker) ValidateConfig(cfg modules.Config) error {
	if _, ok := cfg.GetConsensusConfig().(HasPacemakerConfig); !ok {
		return fmt.Errorf("cannot cast to PacemakeredConsensus")
	}
	return nil
}

func (m *paceMaker) SetConsensusModule(c *consensusModule) {
	m.consensusMod = c
}

func (p *paceMaker) ValidateMessage(m *typesCons.HotstuffMessage) error {
	currentHeight := p.consensusMod.height
	currentRound := p.consensusMod.round
	// Consensus message is from the past
	if m.Height < currentHeight {
		return typesCons.ErrPacemakerUnexpectedMessageHeight(typesCons.ErrOlderMessage, currentHeight, m.Height)
	}

	// Current node is out of sync
	if m.Height > currentHeight {
		// TODO(design): Need to restart state sync
		return typesCons.ErrPacemakerUnexpectedMessageHeight(typesCons.ErrFutureMessage, currentHeight, m.Height)
	}

	// Do not handle messages if it is a self proposal
	if p.consensusMod.isLeader() && m.Type == Propose && m.Step != NewRound {
		// TODO(olshansky): This code branch is a result of the optimization in the leader
		// handlers. Since the leader also acts as a replica but doesn't use the replica's
		// handlers given the current implementation, it is safe to drop proposal that the leader made to itself.
		return typesCons.ErrSelfProposal
	}

	// Message is from the past
	if m.Round < currentRound || (m.Round == currentRound && m.Step < p.consensusMod.step) {
		return typesCons.ErrPacemakerUnexpectedMessageStepRound(typesCons.ErrOlderStepRound, p.consensusMod.step, currentRound, m)
	}

	// Everything checks out!
	if m.Height == currentHeight && m.Step == p.consensusMod.step && m.Round == currentRound {
		return nil
	}

	// Pacemaker catch up! Node is synched to the right height, but on a previous step/round so we just jump to the latest state.
	if m.Round > currentRound || (m.Round == currentRound && m.Step > p.consensusMod.step) {
		p.consensusMod.nodeLog(typesCons.PacemakerCatchup(currentHeight, uint64(p.consensusMod.step), currentRound, m.Height, uint64(m.Step), m.Round))
		p.consensusMod.step = m.Step
		p.consensusMod.round = m.Round

		// TODO(olshansky): Add tests for this. When we catch up to a later step, the leader is still the same.
		// However, when we catch up to a later round, the leader at the same height will be different.
		if currentRound != m.Round || p.consensusMod.leaderId == nil {
			p.consensusMod.electNextLeader(m)
		}

		return nil
	}

	return typesCons.ErrUnexpectedPacemakerCase
}

func (p *paceMaker) RestartTimer() {
	if p.stepCancelFunc != nil {
		p.stepCancelFunc()
	}

	// NOTE: Not defering a cancel call because this function is asynchronous.

	stepTimeout := p.getStepTimeout(p.consensusMod.round)

	clock := p.GetBus().GetRuntimeMgr().GetClock()

	ctx, cancel := clock.WithTimeout(context.TODO(), stepTimeout)
	p.stepCancelFunc = cancel

	go func() {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				p.consensusMod.nodeLog(typesCons.PacemakerTimeout(p.consensusMod.CurrentHeight(), p.consensusMod.step, p.consensusMod.round))
				p.InterruptRound()
			}
		case <-clock.After(stepTimeout + 30*timePkg.Millisecond): // Adding 30ms to the context timeout to avoid race condition.
			return
		}
	}()
}

func (p *paceMaker) InterruptRound() {
	currentRound := p.GetBus().GetConsensusModule().CurrentRound()
	p.consensusMod.nodeLog(typesCons.PacemakerInterrupt(p.GetBus().GetConsensusModule().CurrentHeight(), typesCons.HotstuffStep(p.GetBus().GetConsensusModule().CurrentStep()), currentRound))
	//log.Printf("\n\n\n InterruptRound pre-call Current Round is: %d", currentRound)
	currentRound++

	pacemakerMesage := &typesCons.PacemakerMessage{
		Action: typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_SET_ROUND,
		Function: &typesCons.PacemakerMessage_Round{
			Round: &typesCons.SetRound{
				Round: currentRound,
			},
		},
	}

	anyProto, err := anypb.New(pacemakerMesage)
	if err != nil {
		log.Println("[WARN] Failed to interrupt round: ", err)
		return
	}

	e := &messaging.PocketEnvelope{Content: anyProto}
	p.GetBus().PublishEventToBus(e)

	//TODO! Remove sleep statement. This is currently highly fragile, not suitable for production.
	//! If bus is busy, currently 5 microseconds sleep period will not be enough and tests will fail.
	//! We must confirm the event is consumed in the bust, and then continue.
	time.Sleep(5 * time.Microsecond)

	//log.Printf("\n\n InterruptRound POST-call Current Round is: %d", p.consensusMod.round)
	//log.Printf("\n\n Called Current Round is: %d", p.GetBus().GetConsensusModule().CurrentRound())
	p.startNextView(p.consensusMod.highPrepareQC, false)
}

func (p *paceMaker) NewHeight() {
	currentHeight := p.GetBus().GetConsensusModule().CurrentHeight()
	currentHeight++
	p.consensusMod.nodeLog(typesCons.PacemakerNewHeight(currentHeight))

	p.consensusMod.height++
	pacemakerMesage := &typesCons.PacemakerMessage{
		Action: typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_SET_HEIGHT,
		Function: &typesCons.PacemakerMessage_Height{
			Height: &typesCons.SetHeight{
				Height: currentHeight,
			},
		},
	}

	anyProto, err := anypb.New(pacemakerMesage)
	if err != nil {
		log.Println("[WARN] Failed to interrupt round: ", err)
		return
	}
	//log.Printf("\n\n InterruptRound PRE-call Current Height is: %d", p.consensusMod.height)
	e := &messaging.PocketEnvelope{Content: anyProto}
	p.GetBus().PublishEventToBus(e)

	// //TODO! Remove sleep statement. This is currently highly fragile, not suitable for production.
	// //! If bus is busy, currently 5 microseconds sleep period will not be enough and tests will fail.
	// //! We must confirm the event is consumed in the bust, and then continue.
	time.Sleep(5 * time.Microsecond)

	// log.Printf("InterruptRound POST-call Current Height is: %d ", p.consensusMod.height)
	// log.Printf("InterruptRound p.mod. POST-call Current Height is: %d \n\n\n", p.GetBus().GetConsensusModule().CurrentRound())

	//p.consensusMod.resetForNewHeight()
	pacemakerMesage = &typesCons.PacemakerMessage{
		Action: typesCons.PacemakerMessageType_PACEMAKER_MESSAGE_RESET_FOR_NEW_HEIGHT,
	}

	anyProto, err = anypb.New(pacemakerMesage)
	if err != nil {
		log.Println("[WARN] Failed to interrupt round: ", err)
		return
	}

	e = &messaging.PocketEnvelope{Content: anyProto}
	p.GetBus().PublishEventToBus(e)

	//TODO! Remove sleep statement. This is currently highly fragile, not suitable for production.
	//! If bus is busy, currently 5 microseconds sleep period will not be enough and tests will fail.
	//! We must confirm the event is consumed in the bust, and then continue.
	time.Sleep(5 * time.Microsecond)

	p.startNextView(nil, false) // TODO(design): We are omitting CommitQC and TimeoutQC here.

	p.
		GetBus().
		GetTelemetryModule().
		GetTimeSeriesAgent().
		CounterIncrement(
			consensusTelemetry.CONSENSUS_BLOCKCHAIN_HEIGHT_COUNTER_NAME,
		)
}

func (p *paceMaker) startNextView(qc *typesCons.QuorumCertificate, forceNextView bool) {
	// DISCUSS: Should we lock the consensus module here?

	p.consensusMod.step = NewRound
	p.consensusMod.clearLeader()
	p.consensusMod.clearMessagesPool()
	// TECHDEBT: This should be avoidable altogether
	if p.consensusMod.utilityContext != nil {
		if err := p.consensusMod.utilityContext.Release(); err != nil {
			log.Println("[WARN] Failed to release utility context: ", err)
		}
		p.consensusMod.utilityContext = nil
	}

	// TECHDEBT: This if structure for debug purposes only; think of a way to externalize it from the main consensus flow...
	if p.manualMode && !forceNextView {
		p.quorumCertificate = qc
		return
	}

	hotstuffMessage := &typesCons.HotstuffMessage{
		Type:          Propose,
		Height:        p.consensusMod.height,
		Step:          NewRound,
		Round:         p.consensusMod.round,
		Block:         nil,
		Justification: nil, // Set below if qc is not nil
	}

	if qc != nil {
		hotstuffMessage.Justification = &typesCons.HotstuffMessage_QuorumCertificate{
			QuorumCertificate: qc,
		}
	}

	p.RestartTimer()
	p.consensusMod.broadcastToNodes(hotstuffMessage)
}

// TODO(olshansky): Increase timeout using exponential backoff.
func (p *paceMaker) getStepTimeout(round uint64) timePkg.Duration {
	baseTimeout := timePkg.Duration(int64(timePkg.Millisecond) * int64(p.pacemakerCfg.GetTimeoutMsec()))
	return baseTimeout
}
