package repoisitory

import "github.com/Soyaib10/comfort-cocoon/internal/models"

type DatabaseRepo interface{
	AllUsers() bool

	InsertReservation(res models.Reservation) error
} 