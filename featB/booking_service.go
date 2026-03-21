package service

import (
	"context"
	"fmt"
	"time"
)

// ===== Request / Response =====

type LockBookingRequest struct {
	RoomTypeID string
	CheckIn    time.Time
	CheckOut   time.Time
	Quantity   int
	GuestID    string
}

type LockBookingResponse struct {
	BookingID string
	Status    string
	ExpiresAt time.Time
}

// ===== Model =====

type Booking struct {
	BookingID string
	Status    string
	CheckIn   time.Time
	CheckOut  time.Time
	GuestID   string
	ExpiresAt time.Time
}

// ===== Repository Interface =====

type BookingRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	TryLockInventory(ctx context.Context, tx Tx, roomTypeID string, date time.Time, qty int) (bool, error)
	CreateBooking(ctx context.Context, tx Tx, b Booking) error
}

type Tx interface {
	Commit() error
	Rollback() error
}

// ===== Service Interface =====

type BookingService interface {
	LockBooking(ctx context.Context, req LockBookingRequest) (*LockBookingResponse, error)
}

// ===== Implementation =====

type bookingService struct {
	repo BookingRepository
}

func NewBookingService(r BookingRepository) BookingService {
	return &bookingService{repo: r}
}

// ===== Core Logic: Lock + Create Booking =====

func (s *bookingService) LockBooking(ctx context.Context, req LockBookingRequest) (*LockBookingResponse, error) {

	// 1. Validate เบื้องต้น
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}
	if !req.CheckOut.After(req.CheckIn) {
		return nil, fmt.Errorf("invalid date range")
	}

	// 2. สร้างช่วงวัน
	dates := generateDateRange(req.CheckIn, req.CheckOut)

	// 3. เริ่ม transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Lock inventory ทุกวัน
	for _, d := range dates {
		ok, err := s.repo.TryLockInventory(ctx, tx, req.RoomTypeID, d, req.Quantity)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if !ok {
			tx.Rollback()
			return nil, fmt.Errorf("room not available on %s", d.Format("2006-01-02"))
		}
	}

	// 5. สร้าง booking
	booking := Booking{
		BookingID: generateBookingID(),
		Status:    "PENDING",
		CheckIn:   req.CheckIn,
		CheckOut:  req.CheckOut,
		GuestID:   req.GuestID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err = s.repo.CreateBooking(ctx, tx, booking)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 6. commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 7. response
	return &LockBookingResponse{
		BookingID: booking.BookingID,
		Status:    booking.Status,
		ExpiresAt: booking.ExpiresAt,
	}, nil
}

// ===== Helper =====

func generateDateRange(start, end time.Time) []time.Time {
	var dates []time.Time
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}
	return dates
}

func generateBookingID() string {
	return fmt.Sprintf("BKG-%d", time.Now().UnixNano())
}