package users

import (
	"time"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/query"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/resource"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/status"
)

const (
	// TableName is the database table for this resource
	TableName = "users"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "name asc, id desc"
)

// AllowedParams returns an array of acceptable params in update
func AllowedParams() []string {
	return []string{"name", "summary", "email", "text", "title", "password_hash", "notification", "subscription", "keywords"}
}

// AllowedParamsAdmin returns the cols editable by admins
func AllowedParamsAdmin() []string {
	return []string{"status", "name", "summary", "email", "text", "title", "password_hash", "notification", "subscription", "keywords"}
}

// NewWithColumns creates a new user instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *User {

	user := New()
	user.ID = resource.ValidateInt(cols["id"])
	user.CreatedAt = resource.ValidateTime(cols["created_at"])
	user.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	user.Status = resource.ValidateInt(cols["status"])
	user.Email = resource.ValidateString(cols["email"])
	user.Name = resource.ValidateString(cols["name"])
	user.PasswordHash = resource.ValidateString(cols["password_hash"])
	user.PasswordResetAt = resource.ValidateTime(cols["password_reset_at"])
	user.Role = resource.ValidateInt(cols["role"])
	user.Text = resource.ValidateString(cols["text"])
	user.Title = resource.ValidateString(cols["title"])
	user.Notification = resource.ValidateBoolean(cols["notification"])

	user.ApprovedEmail = resource.ValidateString(cols["approved_email"])
	user.Keywords = resource.ValidateStringArray(cols["keywords"])
	user.PersonalEmail = resource.ValidateString(cols["personal_email"])
	user.Subscription = resource.ValidateBoolean(cols["subscription"])
	user.Plan = resource.ValidateString(cols["plan"])
	user.TrialEnd = resource.ValidateTime(cols["trial_end"])

	user.TwitterConnected = resource.ValidateBoolean(cols["twitter_connected"])
	user.TwitterId = resource.ValidateString(cols["twitter_id"])
	user.TwitterUsername = resource.ValidateString(cols["twitter_username"])
	user.TwitterAccessToken = resource.ValidateString(cols["twitter_access_token"])
	user.TwitterRefreshToken = resource.ValidateString(cols["twitter_refresh_token"])
	user.TwitterTokenExpiryTime = resource.ValidateTime(cols["twitter_token_expiry_time"])
	user.TwitterListID = resource.ValidateString(cols["twitter_list_id"])
	user.TwitterFollowers = resource.ValidateStringArray(cols["twitter_followers"])
	user.TwitterListCreationTime = resource.ValidateTime(cols["twitter_list_creation_time"])

	user.AutoLike = resource.ValidateBoolean(cols["auto_like"])

	return user
}

// New creates and initialises a new user instance.
func New() *User {
	user := &User{}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.TableName = TableName
	user.KeyName = KeyName
	user.Status = status.Draft
	return user
}

// FindFirst fetches a single user record from the database using
// a where query with the format and args provided.
func FindFirst(format string, args ...interface{}) (*User, error) {
	result, err := Query().Where(format, args...).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Find fetches a single user record from the database by id.
func Find(id int64) (*User, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all user records matching this query from the database.
func FindAll(q *query.Query) ([]*User, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of users constructed from the results
	var users []*User
	for _, cols := range results {
		p := NewWithColumns(cols)
		users = append(users, p)
	}

	return users, nil
}

// Query returns a new query for users with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for users with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}
