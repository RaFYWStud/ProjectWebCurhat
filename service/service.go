package service

import "projectwebcurhat/contract"

func New(repo *contract.Repository) *contract.Service {
	roomSvc := NewRoomService(repo)
	return &contract.Service{
		Room:      roomSvc,
		Signaling: NewSignalingService(roomSvc),
	}
}
