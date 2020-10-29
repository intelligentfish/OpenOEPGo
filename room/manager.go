package room

import "sync"

var (
	managerOnce sync.Once // manager once
	managerInst *Manager  // manager instance
)

// Manager
type Manager struct {
	sync.RWMutex
	rooms map[int]*Room
}

// NewManager
func NewManager() *Manager {
	return &Manager{rooms: make(map[int]*Room)}
}

// ManagerInst
func ManagerInst() *Manager {
	managerOnce.Do(func() {
		managerInst = NewManager()
	})
	return managerInst
}

// AddRoom
func (object *Manager) AddRoom(room *Room) {
	object.Lock()
	defer object.Unlock()
	object.rooms[room.Id] = room
}

// RemoveRoom
func (object *Manager) RemoveRoom(room *Room) {
	object.Lock()
	defer object.Unlock()
	delete(object.rooms, room.Id)
}
