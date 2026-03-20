package service

import "hotel-project/model"

type BookingService interface {
	LockRoom(guestID string, typeID string, checkInDate string, checkOutDate string) (*model.Booking, error)
	ConfirmBooking(bookingID string) (*model.Booking, error)
	CancelBooking(bookingID string) (*model.Booking, error)
	GetBookingByID(bookingID string) (*model.Booking, error)
}