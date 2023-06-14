package processors

import (
	"fmt"

	"follower-maze/internal/events"
	"follower-maze/internal/types"
)

type processor struct {
	ctx *types.Context
}

type deadLettersProcessor struct {
	ctx     *types.Context
	printDL func(string)
}

// Returns a new dead letter event processor
func NewDeadLettersProcessor(ctx *types.Context) types.Processor {
	return deadLettersProcessor{
		ctx: ctx,
		printDL: func(v string) {
			fmt.Printf("Dead letter event %v", v)
		},
	}
}

// Processes dead letter events
func (p deadLettersProcessor) Start() {
	for newDeadLetter := range p.ctx.DeadLetterChannel {
		p.printDL(newDeadLetter)
	}
}

// Returns a new events processor
func NewEventsProcessor(ctx *types.Context) types.Processor {
	return processor{
		ctx: ctx,
	}
}

// Processes events
func (p processor) Start() {
	lastSeqNo := 0

	for newEvent := range p.ctx.EventChannel {
		p.ctx.EventQueue[newEvent.Sequence] = newEvent
		lastSeqNo = processEventQueue(p.ctx, lastSeqNo)
	}
}

// Processes events in the queue
func processEventQueue(ctx *types.Context, lastSeqNo int) int {
	for {
		nextEvent, ok := ctx.EventQueue[lastSeqNo+1]
		if !ok {
			break
		}

		delete(ctx.EventQueue, lastSeqNo+1)
		processEvent(ctx, nextEvent)
		lastSeqNo++
	}

	return lastSeqNo
}

// Multiplexes across different processors
func processEvent(ctx *types.Context, evt types.Event) {
	switch evt.EventType {
	case "F":
		events.ProcessFollow(ctx, evt)
	case "U":
		events.ProcessUnfollow(ctx, evt)
	case "P":
		events.ProcessPrivateMessage(ctx, evt)
	case "B":
		events.ProcessBroadcast(ctx, evt)
	case "S":
		events.ProcessStatusUpdate(ctx, evt)
	default:
		events.ProcessDeadLetter(ctx, evt.Payload)
	}
}
