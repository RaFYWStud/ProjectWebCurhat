package service

import (
	"encoding/json"
	"log"

	"projectwebcurhat/contract"
	"projectwebcurhat/database"
	"projectwebcurhat/dto"
)

type signalingService struct {
	roomService contract.RoomService
}

func NewSignalingService(roomService contract.RoomService) contract.SignalingService {
	return &signalingService{
		roomService: roomService,
	}
}

func (s *signalingService) HandleMessage(client *database.Client, data []byte) error {
	var msg dto.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return err
	}

	msg.From = client.ID

	switch msg.Type {
	case dto.MessageTypeJoin:
		s.handleJoin(client, &msg)
	case dto.MessageTypeOffer, dto.MessageTypeAnswer, dto.MessageTypeCandidate:
		s.relayMessage(client, &msg)
	case dto.MessageTypeLeave:
		s.handleLeave(client)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}

	return nil
}

func (s *signalingService) handleJoin(client *database.Client, msg *dto.Message) {
	if msg.Username != "" {
		client.Username = msg.Username
	}

	room := s.roomService.FindOrCreateRoom(client)

	readyMsg := dto.Message{
		Type:   dto.MessageTypeReady,
		RoomID: room.ID,
		From:   "server",
	}
	s.sendToClient(client, &readyMsg)

	if room.IsFull() {
		otherClient := room.GetOtherClient(client.ID)
		if otherClient != nil {
			peerJoinMsg := dto.Message{
				Type:     dto.MessageTypeJoin,
				From:     otherClient.ID,
				Username: otherClient.Username,
				RoomID:   room.ID,
			}
			s.sendToClient(client, &peerJoinMsg)

			peerJoinMsg.From = client.ID
			peerJoinMsg.Username = client.Username
			s.sendToClient(otherClient, &peerJoinMsg)

			log.Printf("Room %s is ready with clients %s and %s", room.ID, client.ID, otherClient.ID)
		}
	}
}

func (s *signalingService) handleLeave(client *database.Client) {
	room := s.roomService.GetRoom(client.RoomID)
	if room != nil {
		otherClient := room.GetOtherClient(client.ID)
		if otherClient != nil {
			leaveMsg := dto.Message{
				Type: dto.MessageTypeLeave,
				From: client.ID,
			}
			s.sendToClient(otherClient, &leaveMsg)
		}
	}

	s.roomService.RemoveClientFromRoom(client)
}

func (s *signalingService) relayMessage(client *database.Client, msg *dto.Message) {
	room := s.roomService.GetRoom(client.RoomID)
	if room == nil {
		log.Printf("Room not found for client %s", client.ID)
		return
	}

	otherClient := room.GetOtherClient(client.ID)
	if otherClient == nil {
		log.Printf("No other client found in room %s", room.ID)
		return
	}

	msg.To = otherClient.ID
	s.sendToClient(otherClient, msg)
}

func (s *signalingService) sendToClient(client *database.Client, msg *dto.Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case client.Send <- data:
	default:
		log.Printf("Client %s send channel full, closing connection", client.ID)
		close(client.Send)
	}
}

func (s *signalingService) DisconnectClient(client *database.Client) {
	s.handleLeave(client)
}
