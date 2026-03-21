package featA

type AvailabilityService interface {
	CheckAvailability(
		roomTypeID string,
		checkInDate string,
		checkOutDate string,
	) (bool, int, float64, error)
}