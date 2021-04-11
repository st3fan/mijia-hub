package main

// Mostly taken from https://github.com/pdf/golifx
// This is all way to complicated for what we need. We can simplify this.

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// DefaultTimeout is the default duration after which operations time out
	DefaultSubscriptionTimeout = 2 * time.Second

	// DefaultRetryInterval is the default interval at which operations are retried
	DefaultSubscriptionRetryInterval = 100 * time.Millisecond

	subscriptionChanSize = 16
)

var (
	// ErrSubscriptionClosed connection closed
	ErrSubscriptionClosed = errors.New(`Connection closed`)

	// ErrSubscriptionTimeout timed out
	ErrSubscriptionTimeout = errors.New(`Timed out`)

	// ErrSubscriptionNotFound
	ErrSubscriptionNotFound = errors.New(`Not found`)
)

type Subscription struct {
	sync.Mutex
	id       string
	events   chan interface{}
	quit     chan struct{}
	provider *SubscriptionProvider
}

func newSubscription(provider *SubscriptionProvider) *Subscription {
	return &Subscription{
		id:       uuid.NewString(),
		events:   make(chan interface{}, subscriptionChanSize),
		quit:     make(chan struct{}),
		provider: provider,
	}
}

func (s *Subscription) Events() <-chan interface{} {
	return s.events
}

func (s *Subscription) notify(event interface{}) error {
	timeout := time.After(DefaultSubscriptionTimeout)
	select {
	case <-s.quit:
		return ErrSubscriptionClosed
	case s.events <- event:
		return nil
	case <-timeout:
		return ErrSubscriptionTimeout
	}
}

func (s *Subscription) Close() error {
	s.Lock()
	defer s.Unlock()

	select {
	case <-s.quit:
		return ErrSubscriptionClosed
	default:
		close(s.quit)
		close(s.events)
	}

	return s.provider.unsubscribe(s)
}

//

type SubscriptionProvider struct {
	subscriptions map[string]*Subscription
	sync.RWMutex
}

// Notify sends the provided event to all subscribers
func (s *SubscriptionProvider) Notify(event interface{}) {
	s.RLock()
	defer s.RUnlock()

	for _, subscription := range s.subscriptions {
		if err := subscription.notify(event); err != nil {
			// TODO
		}
	}
}

func (s *SubscriptionProvider) Close() error {
	for _, subscription := range s.subscriptions {
		if err := subscription.Close(); err != nil {
			// TODO What is the best strategy here?
		}
	}
	return nil
}

func (s *SubscriptionProvider) unsubscribe(subscription *Subscription) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.subscriptions[subscription.id]; !ok {
		return ErrSubscriptionNotFound
	}

	return nil
}

// Subscribe returns a new Subscription for this provider
func (s *SubscriptionProvider) Subscribe() *Subscription {
	s.Lock()
	defer s.Unlock()
	if s.subscriptions == nil {
		s.subscriptions = make(map[string]*Subscription)
	}
	sub := newSubscription(s)
	s.subscriptions[sub.id] = sub

	return sub
}
