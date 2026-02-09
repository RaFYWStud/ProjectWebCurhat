package repository

import (
	"projectwebcurhat/database"
	"sync"
)

type roomRepository struct {
	rooms       map[string]*database.Room
	waitingRoom *database.Room
	mutex       sync.RWMutex
}

func NewRoomRepository() *roomRepository {
	return &roomRepository{
		rooms: make(map[string]*database.Room),
	}
}

func (r *roomRepository) CreateRoom(id string) *database.Room {
	room := database.NewRoom(id)
	r.mutex.Lock()
	r.rooms[id] = room
	r.mutex.Unlock()
	return room
}

func (r *roomRepository) GetRoom(roomID string) *database.Room {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.rooms[roomID]
}

func (r *roomRepository) DeleteRoom(roomID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.rooms, roomID)
	if r.waitingRoom != nil && r.waitingRoom.ID == roomID {
		r.waitingRoom = nil
	}
}

func (r *roomRepository) GetRoomCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.rooms)
}

func (r *roomRepository) GetWaitingRoom() *database.Room {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.waitingRoom
}

func (r *roomRepository) SetWaitingRoom(room *database.Room) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.waitingRoom = room
}

func (r *roomRepository) StoreRoom(room *database.Room) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.rooms[room.ID] = room
}
