package tools

import (
	"log"
	"math/rand"
	"sync"

	"github.com/gary23b/sprites/models"
)

type gridBlock struct {
	sprites map[int]models.NearMeInfo
	mutex   sync.RWMutex
}

type brokerPosInfo struct {
	state models.SpriteState
	yGrid int
	xGrid int
	mutex sync.RWMutex
}

type PositionBroker struct {
	sprites map[int]*brokerPosInfo
	mutex   sync.RWMutex

	grid [][]gridBlock
}

func NewPositionBroker() *PositionBroker {
	ret := &PositionBroker{
		sprites: make(map[int]*brokerPosInfo),
		grid:    make([][]gridBlock, 1000),
	}

	for y := range ret.grid {
		ret.grid[y] = make([]gridBlock, 1000)
		for x := range ret.grid[y] {
			g := &ret.grid[y][x]
			g.sprites = make(map[int]models.NearMeInfo)
		}
	}

	return ret
}

func (s *PositionBroker) AddSprite(id int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.sprites[id]
	if ok {
		log.Printf("id: %d already present\n", id)
		return
	}

	x := rand.Intn(1000)
	y := rand.Intn(1000)
	s.sprites[id] = &brokerPosInfo{
		xGrid: x,
		yGrid: y,
	}

	g := &s.grid[y][x]
	g.mutex.Lock()
	g.sprites[id] = models.NearMeInfo{
		SpriteID:   id,
		SpriteType: 0,
		X:          float64(x-500) * 20,
		Y:          float64(y-500) * 20,
	}
	g.mutex.Unlock()
}

func (s *PositionBroker) RemoveSprite(id int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	item, ok := s.sprites[id]
	if !ok {
		return
	}

	g := &s.grid[item.yGrid][item.xGrid]
	g.mutex.Lock()
	delete(g.sprites, id)
	g.mutex.Unlock()

	delete(s.sprites, id)
}

func (s *PositionBroker) UpdateSpriteInfo(id int, state models.SpriteState) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	item, ok := s.sprites[id]
	if !ok {
		log.Printf("id: %d was not found\n", id)
		return
	}
	item.mutex.Lock()
	defer item.mutex.Unlock()

	x := max(0, min(999, int(state.X/20+500)))
	y := max(0, min(999, int(state.Y/20+500)))

	if x != item.xGrid || y != item.yGrid {
		g := &s.grid[item.yGrid][item.xGrid]
		g.mutex.Lock()
		delete(g.sprites, id)
		g.mutex.Unlock()

		item.xGrid = x
		item.yGrid = y
	}

	g := &s.grid[y][x]
	g.mutex.Lock()
	g.sprites[id] = models.NearMeInfo{
		SpriteID:   id,
		SpriteType: state.SpriteType,
		X:          state.X,
		Y:          state.Y,
	}
	g.mutex.Unlock()

	item.state = state
}

func (s *PositionBroker) GetSpriteInfo(id int) models.SpriteState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	item, ok := s.sprites[id]
	if !ok {
		// log.Printf("id: %d was not found\n", id)
		return models.SpriteState{Deleted: true}
	}
	item.mutex.Lock()
	defer item.mutex.Unlock()
	return item.state
}

func (s *PositionBroker) GetSpritesNearMe(x, y, distance float64) []models.NearMeInfo {
	xMin := max(0, min(999, int((x-distance)/20+500)))
	yMin := max(0, min(999, int((y-distance)/20+500)))
	xMax := max(0, min(999, int((x+distance)/20+500)))
	yMax := max(0, min(999, int((y+distance)/20+500)))

	ret := []models.NearMeInfo{}

	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			g := &s.grid[y][x]
			g.mutex.RLock()
			for _, sprite := range g.sprites {
				ret = append(ret, sprite)
			}
			g.mutex.RUnlock()
		}
	}
	return ret
}
