package controller

import "github.com/antnzr/gatemulator/internal/app/gatemulator/service"

type Controller struct {
	Subscription SubscriptionController
	Subscriber   SubscriberController
	Message      MessageController
	Store        StoreController
}

func NewController(
	subscriptionService service.SubscriptionService,
	subscriberService service.SubscriberService,
	messageService service.MessageService,
	storeService service.StoreService,
) *Controller {
	subscriptionController := NewSubscriptionController(subscriptionService)
	subscriberController := NewSubscriberController(subscriberService)
	messageController := NewMessageController(messageService)
	storeController := NewStoreController(storeService)

	return &Controller{
		Subscription: *subscriptionController,
		Subscriber:   *subscriberController,
		Message:      *messageController,
		Store:        *storeController,
	}
}
