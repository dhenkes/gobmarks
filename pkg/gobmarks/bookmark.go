package gobmarks

import (
	"context"
)

// Bookmark constants.
const (
	MaxBookmarkTitleLen       = 255
	MaxBookmarkUrlLen         = 255
	MaxBookmarkDescriptionLen = 255
)

// Bookmark represents a bookmark in the system.
type Bookmark struct {
	ID          string `json:"id"`
	UserID      string `json:"users_id"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Html        string `json:"html"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	RemovedAt   int64  `json:"removed_at"`
}

// Validate returns an error if the bookmark contains invalid fields.
func (b *Bookmark) Validate() error {
	if b.UserID == "" {
		return NewError(EINVALID, "User ID required.")
	}

	if b.Title == "" {
		return NewError(EINVALID, "Title required.")
	}

	if len(b.Title) > MaxBookmarkTitleLen {
		return NewError(EINVALID, "Title must be less than %d characters.", MaxBookmarkTitleLen)
	}

	if b.Url == "" {
		return NewError(EINVALID, "Url required.")
	}

	if len(b.Url) > MaxBookmarkUrlLen {
		return NewError(EINVALID, "Url must be less than %d characters.", MaxBookmarkUrlLen)
	}

	if b.Description != "" && len(b.Description) > MaxBookmarkDescriptionLen {
		return NewError(EINVALID, "Description must be less than %d characters.", MaxBookmarkDescriptionLen)
	}

	return nil
}

// CanFindBookmark returns true if the current user can list bookmarks with
// the given filter.
func CanFindBookmark(ctx context.Context, filter BookmarkFilter) bool {
	id := UserIDFromContext(ctx)
	return id != "" && filter.UserID == &id
}

// CanUpdateBookmark returns true if the current user can update the bookmark.
func CanUpdateBookmark(ctx context.Context, bookmark *Bookmark) bool {
	if user := UserFromContext(ctx); user != nil && user.IsDemo {
		return false
	} else {
		id := UserIDFromContext(ctx)
		return id != "" && bookmark.UserID == id
	}
}

// BookmarkService represents a service for managing bookmarks. The functions
// should return ENOTFOUND if the bookmark could not be found and EUNAUTHORIZED
// if the user is not authorized to run the transaction.
type BookmarkService interface {
	FindBookmarkByID(ctx context.Context, id string) (*Bookmark, error)
	FindBookmarks(ctx context.Context, filter BookmarkFilter) ([]*Bookmark, int, error)
	CreateBookmark(ctx context.Context, bookmark *Bookmark) error
	UpdateBookmark(ctx context.Context, id string, update BookmarkUpdate) (*Bookmark, error)
	RemoveBookmark(ctx context.Context, id string) error
}

// BookmarkFilter represents a filter passed to FindBookmarks().
type BookmarkFilter struct {
	ID     *string `json:"id"`
	UserID *string `json:"users_id"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// BookmarkUpdate represents a set of fields to be updated via UpdateBookmark().
type BookmarkUpdate struct {
	Title       *string `json:"title"`
	Url         *string `json:"url"`
	Description *string `json:"description"`
	Html        *string `json:"html"`
}
