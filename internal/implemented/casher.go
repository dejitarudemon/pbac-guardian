package implemented

import (
	"context"
	"sync"
)

type Storage map[string]any

type SessionMap map[string]Storage

type DefaultCasher struct {
	mx      *sync.RWMutex
	storage SessionMap
}

func NewDefaultCasher() *DefaultCasher {
	return &DefaultCasher{
		mx:      &sync.RWMutex{},
		storage: make(SessionMap),
	}
}

func (c *DefaultCasher) Set(ctx context.Context, sessionID, key string, value any) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.storage[sessionID]; !ok {
		c.storage[sessionID] = make(Storage)
	}

	c.storage[sessionID][key] = value

	return nil
}

func (c *DefaultCasher) Get(ctx context.Context, sessionID, key string) (any, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if session, sok := c.storage[sessionID]; sok {
		if value, vok := session[key]; vok {
			return value, nil
		}
	}

	return nil, nil
}

func (c *DefaultCasher) Clear(ctx context.Context, sessionID string) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	delete(c.storage, sessionID)

	return nil
}
