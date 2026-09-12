package domain

import (
	"errors"
	"time"
)

var (
	ErrCapsuleNotFound      = errors.New("capsule not found")
	ErrCannotOpenYet        = errors.New("cannot open capsule before open_at time")
	ErrInvalidOpenAt        = errors.New("open_at time cannot be in the past")
	ErrCapsuleAlreadySealed = errors.New("capsule is already sealed")
	ErrValueEmpty           = errors.New("value cannot be empty")
	ErrNotEditable          = errors.New("cupsule not Editable")
	ErrIsNotSealed          = errors.New("capsule is not sealed")
)

type Cupsule struct {
	ID         string
	OwnerID    string
	Title      string
	Message    string
	Status     Status
	Visibility Visibility
	OpenAt     time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Status string

const (
	StatusUnspecified Status = "UNSPECIFIED"
	StatusDraft       Status = "DRAFT"
	StatusSealed      Status = "SEALED"
	StatusOpened      Status = "OPENED"
	StatusCancelled   Status = "CANCELLED"
)

type Visibility string

const (
	VisibilityUnspecified Visibility = "UNSPECIFIED"
	VisibilityPrivate     Visibility = "PRIVATE"
	VisibilityShared      Visibility = "SHARED"
	VisibilityPublic      Visibility = "PUBLIC"
)
