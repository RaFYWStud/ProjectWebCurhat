package service

import (
	"log"

	"github.com/google/uuid"
	"projectwebcurhat/contract"
	"projectwebcurhat/database"
)

type roomService struct {
	repo *contract.Repository
}

func NewRoomService(repo *contract.Repository) contract.RoomService {
	return &roomService{repo: repo}
}

func (s *roomService) FindOrCreateRoom(client *database.Client) *database.Room {
	waitingRoom := s.repo.Room.GetWaitingRoom()

	if waitingRoom != nil && !waitingRoom.IsFull() {
		if waitingRoom.AddClient(client) {
			log.Printf("Client %s joined existing room %s", client.ID, waitingRoom.ID)
			s.repo.Room.SetWaitingRoom(nil)
			return waitingRoom
		}
	}

	roomID := uuid.New().String()
	room := s.repo.Room.CreateRoom(roomID)
	room.AddClient(client)
	s.repo.Room.SetWaitingRoom(room)

	log.Printf("Created new room %s for client %s", roomID, client.ID)
	return room
}

func (s *roomService) GetRoom(roomID string) *database.Room {
	return s.repo.Room.GetRoom(roomID)
}

func (s *roomService) RemoveClientFromRoom(client *database.Client) {
	if client.RoomID == "" {
		return
	}

	room := s.repo.Room.GetRoom(client.RoomID)
	if room == nil {
		return
	}

	room.RemoveClient(client.ID)
	log.Printf("Client %s removed from room %s", client.ID, room.ID)

	if room.IsEmpty() {
		s.repo.Room.DeleteRoom(room.ID)
		log.Printf("Room %s deleted (empty)", room.ID)
	}
}

func (s *roomService) GetRoomCount() int {
	return s.repo.Room.GetRoomCount()
}
