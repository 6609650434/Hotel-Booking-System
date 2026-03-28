package repository

import "hotel-project/model"

type BookingRepository interface {
	Save(booking *model.Booking) error
	FindByID(bookingID string) (*model.Booking, error)
	UpdateStatus(bookingID string, status string) error
	UpdateLockTime(bookingID string, lockedAt string, expiresAt string) error
}