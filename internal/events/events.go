package events

import (
	"follower-maze/internal/notifications"
	"follower-maze/internal/types"
)

// ProcessBroadcast Adds the relationship of the users in the follow registry
func ProcessBroadcast(ctx *types.Context, evt types.Event) {
	notifications.SendEventToAllUsers(ctx, evt)
}

// ProcessDeadLetter adds the event to the dead letter queue
func ProcessDeadLetter(ctx *types.Context, payload string) {
	ctx.DeadLetterChannel <- payload
}

// ProcessFollow Adds the relationship of the users in the follow registry
func ProcessFollow(ctx *types.Context, evt types.Event) {
	// Fetch the follow list of the followed user and if it doesn't have one create it
	if _, ok := ctx.FollowRegistry[evt.ReceiverUserID]; !ok {
		ctx.FollowRegistry[evt.ReceiverUserID] = make(map[int]bool)
	}

	// Add the sender to the followers of the user
	followers := ctx.FollowRegistry[evt.ReceiverUserID]
	followers[evt.EmitterUserID] = true

	notifications.SendEventToUser(ctx, evt.ReceiverUserID, evt)
}

// ProcessPrivateMessage Adds the relationship of the users in the follow registry
func ProcessPrivateMessage(ctx *types.Context, evt types.Event) {
	notifications.SendEventToUser(ctx, evt.ReceiverUserID, evt)
}

// ProcessStatusUpdate Adds the relationship of the users in the follow registry
func ProcessStatusUpdate(ctx *types.Context, evt types.Event) {
	notifications.SendNotificationToAllFollowers(ctx, evt.EmitterUserID, evt)
}

// ProcessUnfollow Adds the relationship of the users in the follow registry
func ProcessUnfollow(ctx *types.Context, evt types.Event) {
	if followers, ok := ctx.FollowRegistry[evt.ReceiverUserID]; ok {
		delete(followers, evt.EmitterUserID)
	}
}
