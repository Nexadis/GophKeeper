package models

import (
	"time"
)

// TimeProvider позволяет переопределить источник времени для объектов в пакетах
type TimeProvider interface {
	Now() time.Time
}
