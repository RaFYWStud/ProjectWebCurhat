package contract

import "projectwebcurhat/database"

type Repository struct {
	Room RoomRepository
}

type RoomRepository interface {
	CreateRoom(id string) *database.Room
	GetRoom(roomID string) *database.Room
	DeleteRoom(roomID string)
	GetRoomCount() int
	GetWaitingRoom() *database.Room
	SetWaitingRoom(room *database.Room)
	StoreRoom(room *database.Room)
}
