package users

import (
	"time"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/resource"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/status"
)

// User handles saving and retrieving users from the database
type User struct {
	// resource.Base defines behavior and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behavior
	status.ResourceStatus

	Organs []Organ

	Email        string
	Name         string
	Role         int64
	Text         string
	Title        string
	Notification bool

	PasswordHash    string
	PasswordResetAt time.Time

	ApprovedEmail string
	Keywords      []string

	PersonalEmail string
	Plan          string
	Subscription  bool
	TrialEnd      time.Time

	TwitterConnected        bool
	TwitterId               string
	TwitterUsername         string
	TwitterAccessToken      string
	TwitterRefreshToken     string
	TwitterTokenExpiryTime  time.Time
	TwitterListID           string
	TwitterFollowers        []string
	TwitterListCreationTime time.Time

	AutoLike bool
}
