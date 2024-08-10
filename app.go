package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"voxel/store"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	err := store.InitDatabase("./clipboard.db")
	if err != nil {
		log.Fatal(err)
	}
	a.ctx = ctx
}

func (b *App) shutdown(ctx context.Context) {
	store.ShutDownDatabase()
}

// Greet returns a greeting message
func (a *App) Greet(name string) string {
	fmt.Println(name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// AddClipBoardItem adds a new clipboard item to the database
func (a *App) AddClipBoardItem(content string, categories []string, contentType string) error {
	item := store.ClipboardItem{
		Content:   content,
		Timestamp: time.Now(),
		Type:      contentType,
	}
	item.SetCategoriesFromArray(categories)
	err := store.AddClipboardItem(&item)

	log.Println("New clipboard  item is getting added to database")

	if err != nil {
		log.Printf("%v", err)
		return fmt.Errorf("failed to add clipboard item: %v", err)
	}
	log.Println("New clipboard item added to database")
	return nil
}

// GetClipBoardItemById retrieves a clipboard item by its ID
func (a *App) GetClipBoardItemById(id int) (store.ClipboardItemDbRow, error) {
	item, err := store.GetClipboardItemById(id)
	if err != nil {
		log.Printf("%v", err)
		return store.ClipboardItemDbRow{}, fmt.Errorf("failed to get clipboard item: %v", err)
	}
	return item, nil
}

// GetAllClipBoardItems retrieves all clipboard items from the database
func (a *App) GetAllClipBoardItems() ([]store.ClipboardItemDbRow, error) {
	items, err := store.GetAllClipboardItems()
	if err != nil {
		log.Printf("%v", err)
		return nil, fmt.Errorf("failed to get all clipboard items: %v", err)
	}

	for _, item := range items {
		fmt.Printf("Clipboard things %s", item.Content)
	}

	return items, nil
}

// DeleteClipBoardItemById deletes a clipboard item by its ID

func (a *App) DeleteClipBoardItem(id int) error {
	return a.DeleteClipBoardItemById(id)
}

func (a *App) DeleteClipBoardItemById(id int) error {
	err := store.DeleteClipboardItemById(id)
	if err != nil {
		log.Printf("%v", err)
		return fmt.Errorf("failed to delete clipboard item: %v", err)
	}
	return nil
}

// UpdateClipBoardItemById updates a clipboard item by its ID
func (a *App) UpdateClipBoardItemById(id int, content string, categories []string, contentType string) error {
	item := store.ClipboardItem{
		Content:   content,
		Timestamp: time.Now(),
		Type:      contentType,
	}
	item.SetCategoriesFromArray(categories)

	err := store.UpdateClipboardItemById(id, &item)
	if err != nil {
		log.Printf("%v", err)
		return fmt.Errorf("failed to update clipboard item: %v", err)
	}
	return nil
}

func (a *App) GetClipboardContent() string {
	return "This is a clipboard item from the backend"
}
