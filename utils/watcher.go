package utils

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/clipboard"
	"log"
	"time"
	"voxel/store"
)

type Utils struct {
	ctx context.Context
}

func NewUtils() *Utils {
	return &Utils{}
}

func (b *Utils) Startup(ctx context.Context) {
	b.ctx = ctx
}

func (u *Utils) WatchClipboard() {
	err := clipboard.Init()
	if err != nil {
		log.Fatalf("Failed to initialize clipboard: %v", err)
	}

	// Watch for text changes in the clipboard
	ch := clipboard.Watch(context.Background(), clipboard.FmtText)
	for data := range ch {
		if len(data) > 0 {
			text := string(data)
			fmt.Printf("New clipboard content: %s\n", text)
			u.saveToDatabase(text)
		}
	}
}

func (u *Utils) saveToDatabase(text string) {
	item := store.ClipboardItem{
		Content:   text,
		Timestamp: time.Now(),
		Type:      "text/plain",
	}
	categories := []string{"Default"}

	item.SetCategoriesFromArray(categories)

	store.AddClipboardItem(&item)
	runtime.EventsEmit(u.ctx, "newClipboardData")

	fmt.Printf("Saved to database: %s\n", text)
}
