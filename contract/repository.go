package contract

import "projectwebcurhat/database"

type Repository struct {
	Room RoomRepository
	User UserRepository
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

type UserRepository interface {
	CreateUser(user *database.User) (*database.User, error)
	GetUserByEmail(email string) (*database.User, error)
	GetUserByID(id int) (*database.User, error)
	GetUserByUsername(username string) (*database.User, error)
	UpdateUser(user *database.User) (*database.User, error)
	SetOnlineStatus(userID int, online bool) error
}
