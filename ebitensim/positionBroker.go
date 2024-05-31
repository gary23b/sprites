package ebitensim

import (
	"log"
	"sync"

	"github.com/gary23b/sprites/models"
)

type brokerPosInfo struct {
	state models.SpriteState
	mutex sync.RWMutex
}

type positionBroker struct {
	sprites map[string]*brokerPosInfo
	mutex   sync.RWMutex
}

func newPositionBroker() *positionBroker {
	ret := &positionBroker{
		sprites: make(map[string]*brokerPosInfo),
	}

	return ret
}

func (s *positionBroker) addSprite(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.sprites[name]
	if ok {
		log.Printf("Name: %s already present\n", name)
		return
	}

	s.sprites[name] = &brokerPosInfo{}
}

func (s *positionBroker) removeSprite(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sprites, name)
}

func (s *positionBroker) updateSpriteInfo(name string, state models.SpriteState) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	item, ok := s.sprites[name]
	if !ok {
		log.Printf("Name: %s was not found\n", name)
		return
	}
	item.mutex.Lock()
	defer item.mutex.Unlock()
	item.state = state

}

func (s *positionBroker) getSpriteInfo(name string) models.SpriteState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	item, ok := s.sprites[name]
	if !ok {
		log.Printf("Name: %s was not found\n", name)
		return models.SpriteState{}
	}
	item.mutex.Lock()
	defer item.mutex.Unlock()
	return item.state
}
