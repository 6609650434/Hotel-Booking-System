package service

import "hotel-project/model"

type RoomService interface {
	CheckAvailability(typeID string, checkInDate string, checkOutDate string) (int, error)
	GetRoomTypeByID(typeID string) (*model.RoomType, error)
}