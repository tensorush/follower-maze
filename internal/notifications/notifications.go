package notifications

import (
	"fmt"

	"follower-maze/internal/types"
)

// Sends an event payload to a connected user, otherwise - to the dead letter queue.
func SendEventToUser(ctx *types.Context, userID int, evt types.Event) {
	var err error = nil
	clientConn, ok := ctx.UsersPool[userID]
	if ok {
		_, err = fmt.Fprint(clientConn, evt.Payload)
	}

	if !ok {
		ctx.DeadLetterChannel <- evt.Payload
	}

	if err != nil {
		fmt.Printf("Could not send message to user: %d error: %v", userID, err)
		ctx.DeadLetterChannel <- evt.Payload
	}
}

// Sends an event payload to all connected users.
func SendEventToAllUsers(ctx *types.Context, evt types.Event) {
	for userID, clientConn := range ctx.UsersPool {
		_, err := fmt.Fprint(clientConn, evt.Payload)

		if err != nil {
			fmt.Printf("Could not send message to user: %d error: %v", userID, err)
			ctx.DeadLetterChannel <- evt.Payload
		}
	}
}

// Sends an event payload to all followers of a user.
func SendNotificationToAllFollowers(ctx *types.Context, userID int, evt types.Event) {
	if followers, ok := ctx.FollowRegistry[userID]; ok {
		for follower := range followers {
			SendEventToUser(ctx, follower, evt)
		}
	}
}
