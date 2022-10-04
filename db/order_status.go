package db

// OrderStatus contains the different types of Order status.
type OrderStatus int

const (
	New OrderStatus = iota
	InProgress
	Completed
	Rejected
)

func (o OrderStatus) String() string {
	return [...]string{"NEW", "IN_PROGRESS", "COMPLETED", "REJECTED"}[o]
}
