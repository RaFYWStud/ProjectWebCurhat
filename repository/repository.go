package repository

import (
	"projectwebcurhat/contract"

	"gorm.io/gorm"
)

func New(db *gorm.DB) *contract.Repository {
	return &contract.Repository{
		Room: NewRoomRepository(),
		User: NewUserRepository(db),
	}
}
