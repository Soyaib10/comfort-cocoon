package repoisitory

import (
	"time"

	"github.com/Soyaib10/comfort-cocoon/internal/models"
)

type DatabaseRepo interface{
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error)
} 