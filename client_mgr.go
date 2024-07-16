//Copyright The ZHIYUNCo.All rights reserved.
//Created by admin at2024/7/16.

package gogowebsocket

import "sync"

type clientsMgr struct {
	// 连接的clients.
	clients map[string]*Client
	rwlock  sync.RWMutex
}

func newClientsMgr() *clientsMgr {
	return &clientsMgr{
		clients: make(map[string]*Client),
	}
}

func (s *clientsMgr) addClient(c *Client) {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	s.clients[c.GetClientId()] = c
}

func (s *clientsMgr) delClient(clientId string) {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	delete(s.clients, clientId)
}

func (s *clientsMgr) delAll() {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	s.clients = make(map[string]*Client)
}

func (s *clientsMgr) getClient(clientId string) *Client {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()
	return s.clients[clientId]
}

func (s *clientsMgr) getClients() map[string]*Client {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()
	newclients := make(map[string]*Client)

	for key, value := range s.clients {
		newclients[key] = value
	}
	return newclients
}
