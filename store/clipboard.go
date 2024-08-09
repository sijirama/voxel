package store

import (
	"database/sql"
	"strings"
	"time"
)

// ClipboardItem represents a single clipboard entry.
type ClipboardItem struct {
	Content    string    `json:"content" db:"content"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	Type       string    `json:"type" db:"type"`
	Categories string    `json:"categories" db:"categories"`
}

// ClipboardItemDbRow represents a clipboard item with an ID from the database.
type ClipboardItemDbRow struct {
	ID int `json:"id" db:"id"`
	ClipboardItem
}

// GetClipboardItemById retrieves a clipboard item by its ID.
func GetClipboardItemById(id int) (ClipboardItemDbRow, error) {
	var ci ClipboardItemDbRow
	err := store.QueryRow("SELECT id, content, timestamp, type, categories FROM clipboard_items WHERE id = ?", id).Scan(&ci.ID, &ci.Content, &ci.Timestamp, &ci.Type, &ci.Categories)
	if err != nil {
		if err == sql.ErrNoRows {
			return ci, nil // Return an empty struct with no error if no rows were found
		}
		return ci, err
	}
	return ci, nil
}

// GetAllClipboardItems retrieves all clipboard items from the database.
func GetAllClipboardItems() ([]ClipboardItemDbRow, error) {
	rows, err := store.Query("SELECT id, content, timestamp, type, categories FROM clipboard_items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ClipboardItemDbRow, 0)

	for rows.Next() {
		var ci ClipboardItemDbRow
		err := rows.Scan(&ci.ID, &ci.Content, &ci.Timestamp, &ci.Type, &ci.Categories)
		if err != nil {
			return nil, err
		}
		items = append(items, ci)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// AddClipboardItem adds a new clipboard item to the database.
func AddClipboardItem(item *ClipboardItem) error {
	_, err := store.Exec("INSERT INTO clipboard_items (content, timestamp, type, categories) VALUES (?, ?, ?, ?)",
		item.Content, item.Timestamp, item.Type, item.Categories)
	if err != nil {
		return err
	}
	return nil
}

// DeleteClipboardItemById deletes a clipboard item by its ID.
func DeleteClipboardItemById(id int) error {
	_, err := store.Exec("DELETE FROM clipboard_items WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateClipboardItemById updates a clipboard item by its ID.
func UpdateClipboardItemById(id int, item *ClipboardItem) error {
	_, err := store.Exec("UPDATE clipboard_items SET content = ?, timestamp = ?, type = ?, categories = ? WHERE id = ?",
		item.Content, item.Timestamp, item.Type, item.Categories, id)
	if err != nil {
		return err
	}
	return nil
}

// ConvertCategoriesToArray splits the comma-separated categories into a slice of strings.
func (ci *ClipboardItemDbRow) ConvertCategoriesToArray() []string {
	return strings.Split(ci.Categories, ",")
}

// SetCategoriesFromArray converts a slice of categories into a comma-separated string.
func (ci *ClipboardItem) SetCategoriesFromArray(categories []string) {
	ci.Categories = strings.Join(categories, ",")
}
