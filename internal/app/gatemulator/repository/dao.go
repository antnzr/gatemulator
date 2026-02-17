package repository

type DAO interface {
	Subscription() SubscriptionRepository
	Subscriber() SubscriberRepository
	MessageFile() MessageFileRepository
}

type dao struct {
	subscription SubscriptionRepository
	subscriber   SubscriberRepository
	messageFile  MessageFileRepository
}

func NewDAO(subscription SubscriptionRepository, subscriber SubscriberRepository, messageFile MessageFileRepository) DAO {
	return &dao{subscription, subscriber, messageFile}
}

func (d *dao) Subscription() SubscriptionRepository {
	return d.subscription
}

func (d *dao) Subscriber() SubscriberRepository {
	return d.subscriber
}

func (d *dao) MessageFile() MessageFileRepository {
	return d.messageFile
}
