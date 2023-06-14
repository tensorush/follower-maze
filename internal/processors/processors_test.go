package processors

import (
	"testing"
	"time"

	"follower-maze/internal/types"

	"github.com/stretchr/testify/assert"
)

func resetContext() types.Context {
	return types.Context{
		EventQueue:        make(map[int]types.Event),
		FollowRegistry:    make(map[int]map[int]bool),
		UsersPool:         make(map[int]types.Connection),
		EventsPort:        9090,
		SubscriptionPort:  9099,
		EventChannel:      make(chan types.Event, 100),
		DeadLetterChannel: make(chan string, 100),
	}
}

func TestProcessEvent(t *testing.T) {
	ctx := resetContext()

	tests := []struct {
		name         string
		events       []types.Event
		wantErr      bool
		emptyResults bool
	}{
		{
			"TestFollow",
			[]types.Event{
				{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "F",
					Payload:        "1|F|1|2\n",
				},
				{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "F",
					Payload:        "1|F|1|2\n",
				},
			},
			false,
			false,
		},
		{
			"TestNonExistent",
			[]types.Event{
				{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "X",
					Payload:        "1|F|1|2\n",
				},
				{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "Y",
					Payload:        "1|U|1|2\n",
				},
			},
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, evt := range tt.events {
				processEvent(&ctx, evt)

				switch evt.EventType {
				case "F":
					assert.Equal(t, true, ctx.FollowRegistry[evt.ReceiverUserID][evt.EmitterUserID])
				case "X", "Y":
					assert.Equal(t, evt.Payload, <-ctx.DeadLetterChannel)
				}

			}

			close(ctx.DeadLetterChannel)
			close(ctx.EventChannel)
			ctx = resetContext()
		})
	}
}

func TestNewDeadLettersProcessor(t *testing.T) {
	ctx := resetContext()

	tests := []struct {
		name         string
		events       []types.Event
		wantErr      bool
		emptyResults bool
	}{
		{
			"TestFollow",
			[]types.Event{
				{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "F",
					Payload:        "1|F|1|2\n",
				},
			},
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := []string{}
			p := deadLettersProcessor{
				ctx: &ctx,
				printDL: func(v string) {
					results = append(results, v)
				},
			}

			for _, evt := range tt.events {
				ctx.DeadLetterChannel <- evt.Payload

				go p.Start()
				time.Sleep(300 * time.Millisecond)
				assert.Equal(t, evt.Payload, results[0])
				close(ctx.DeadLetterChannel)
				close(ctx.EventChannel)
			}

			ctx = resetContext()
		})
	}
}
