package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"voxel/store"
	"voxel/utils"
	"github.com/atotto/clipboard"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	err := store.InitDatabase("./clipboard.db")
	if err != nil {
		log.Fatal(err)
	}

	a.ctx = ctx

	utils := utils.NewUtils()

	utils.Startup(ctx)

	go utils.WatchClipboard()
	select {}
}

func (b *App) shutdown(ctx context.Context) {
	store.ShutDownDatabase()
}

func (a *App) Greet(name string) string {
	fmt.Println(name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

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

func (a *App) GetClipBoardItemById(id int) (store.ClipboardItemDbRow, error) {
	item, err := store.GetClipboardItemById(id)
	if err != nil {
		log.Printf("%v", err)
		return store.ClipboardItemDbRow{}, fmt.Errorf("failed to get clipboard item: %v", err)
	}
	return item, nil
}

func (a *App) GetAllClipBoardItems() ([]store.ClipboardItemDbRow, error) {
	items, err := store.GetAllClipboardItems()
	if err != nil {
		log.Printf("%v", err)
		return nil, fmt.Errorf("failed to get all clipboard items: %v", err)
	}
	return items, nil
}


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
	currentItem, err := clipboard.ReadAll()

	if err != nil {
		fmt.Println(err)
	}

	if currentItem != "" {
		return currentItem
	}
	//NOTE: this is fucking retarded siji, wtf are you even doing lmao

	return ""
}
