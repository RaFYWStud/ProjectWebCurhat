package contract

import (
	"projectwebcurhat/database"
	"projectwebcurhat/dto"
)

type Service struct {
	Room      RoomService
	Signaling SignalingService
	Auth      AuthService
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

type AuthService interface {
	Register(payload *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(payload *dto.LoginRequest) (*dto.AuthResponse, error)
	GetProfile(userID int) (*dto.UserProfile, error)
}
