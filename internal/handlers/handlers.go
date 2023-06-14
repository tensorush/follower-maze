package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"follower-maze/internal/events"
	"follower-maze/internal/types"
)

// Subscribes a new user.
func NewSubscriptionHandler(ctx *types.Context) types.Handler {
	return func(conn types.Connection, message string) {
		subscriptionHandler(ctx, conn, message)
	}
}

func subscriptionHandler(ctx *types.Context, conn types.Connection, message string) {
	userID, err := strconv.Atoi(message)
	if err != nil {
		log.Fatal(err)
	}

	ctx.UsersPool[userID] = conn

	fmt.Printf("User connected: %d (%d total)\n", userID, len(ctx.UsersPool))
}

// Processes a new event
func NewEventsHandler(ctx *types.Context) types.Handler {
	return func(conn types.Connection, message string) {
		eventsHandler(ctx, conn, message)
	}
}

func eventsHandler(ctx *types.Context, conn types.Connection, message string) {
	var err error
	eventParts := strings.Split(message, "|")

	seqNum, err := strconv.Atoi(eventParts[0])
	if err != nil {
		events.ProcessDeadLetter(ctx, message)
		log.Printf("Couldn't process message %v", message)
		return
	}

	if len(eventParts) < 1 {
		events.ProcessDeadLetter(ctx, message)
		log.Printf("Couldn't process message %v", message)
		return
	}

	eventType := eventParts[1]

	emitterUserID := 0
	if len(eventParts) > 2 {
		emitterUserID, err = strconv.Atoi(eventParts[2])
		if err != nil {
			events.ProcessDeadLetter(ctx, message)
			log.Printf("Couldn't process message %v", err)
			return
		}
	}

	receiverUserID := 0
	if len(eventParts) > 3 {
		receiverUserID, err = strconv.Atoi(eventParts[3])
		if err != nil {
			events.ProcessDeadLetter(ctx, message)
			log.Printf("Couldn't process message %v", err)
			return
		}
	}

	evt := types.Event{
		Sequence:       seqNum,
		ReceiverUserID: receiverUserID,
		EmitterUserID:  emitterUserID,
		EventType:      eventType,
		Payload:        message + "\n",
	}

	ctx.EventChannel <- evt
}
