package contract

import "projectwebcurhat/database"

type Service struct {
	Room      RoomService
	Signaling SignalingService
}

type RoomService interface {
	FindOrCreateRoom(client *database.Client) *database.Room
	GetRoom(roomID string) *database.Room
	RemoveClientFromRoom(client *database.Client)
	GetRoomCount() int
}

type SignalingService interface {
	HandleMessage(client *database.Client, data []byte) error
	DisconnectClient(client *database.Client)
}
