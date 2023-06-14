package main

import (
	"sync"

	"follower-maze/internal/handlers"
	"follower-maze/internal/processors"
	"follower-maze/internal/servers"
	"follower-maze/internal/types"
)

func main() {
	ctx := types.Context{
		EventQueue:        make(map[int]types.Event),
		FollowRegistry:    make(map[int]map[int]bool),
		UsersPool:         make(map[int]types.Connection),
		EventsPort:        9090,
		SubscriptionPort:  9099,
		EventChannel:      make(chan types.Event, 100),
		DeadLetterChannel: make(chan string, 100),
	}

	eventsServer := servers.NewEventsServer(handlers.NewEventsHandler(&ctx))
	subscriptionServer := servers.NewSubscriptionServer(handlers.NewSubscriptionHandler(&ctx))

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		eventsServer.Listen(ctx.EventsPort)
	}()

	go func() {
		defer wg.Done()
		subscriptionServer.Listen(ctx.SubscriptionPort)
	}()

	go func() {
		defer wg.Done()
		processors.NewEventsProcessor(&ctx).Start()
	}()

	go func() {
		defer wg.Done()
		processors.NewDeadLettersProcessor(&ctx).Start()
	}()

	wg.Wait()
}
