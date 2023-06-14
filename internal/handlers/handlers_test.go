package handlers

import (
	"fmt"
	"testing"

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

func TestEventsHandlers(t *testing.T) {
	ctx := resetContext()

	tests := []struct {
		name         string
		events       []types.Event
		wantErr      bool
		emptyResults bool
	}{
		{
			"TestBadPayload",
			[]types.Event{
				{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "F",
					Payload:        "1\n",
				},
				{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "F",
					Payload:        "\n",
				},
				{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "F",
					Payload:        "1|1|B|2\n",
				},
				{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "F",
					Payload:        "1|2|3|B\n",
				},
			},
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, evt := range tt.events {
				eventsHandler(&ctx, nil, evt.Payload)
				x := <-ctx.DeadLetterChannel
				fmt.Printf("ASSERT %v vs %v", evt.Payload, x)
				assert.Equal(t, evt.Payload, x)
			}

			close(ctx.DeadLetterChannel)
			close(ctx.EventChannel)
			ctx = resetContext()
		})
	}
}
