package repository

import "hotel-project/model"

type RoomRepository interface {
	FindRoomTypeByID(typeID string) (*model.RoomType, error)
}