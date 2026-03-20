package repository

type InventoryRepository interface {
	CheckAvailableCount(typeID string, checkInDate string, checkOutDate string) (int, error)
	DecreaseAvailableCount(typeID string, checkInDate string, checkOutDate string) error
	IncreaseAvailableCount(typeID string, checkInDate string, checkOutDate string) error
}v