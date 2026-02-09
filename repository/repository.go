package repository

import "projectwebcurhat/contract"

func New() *contract.Repository {
	return &contract.Repository{
		Room: NewRoomRepository(),
	}
}
