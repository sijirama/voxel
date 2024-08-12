package main
//https://video-downloads.googleusercontent.com/ADGPM2lKKGiTlV7rE_nUZ7Hqxmv41vx6FTWVFosFb4ZN5-YF3OpZfdkIDxhN0pRYEiExb9wvW5X6EHyzYcCQGYGyksu4jBuB2tZ-e-OMs6Sw1OKAho8KYXS05FauJ3IS3mP98l6o6-EDKybNiDnZAHUvHPNz2qfaW-l8WC5w2y8ZQcI1eHywe-mJSzQk8ZcGFs7lq0RZ3S3Nt8JWC89V9f6VJx4iCHg0Yqa2zvTu5HK2XxAUdddws70pywcwptTT_jVZFA_kzM_o_8dF59DqDAib16nGP2FnzZu5XCqu9s6S6gJgujhqCr7S_4PEpbcrs2AZyQNpJnmPOWO4i6Z0_a0-FjLfWqyZqwzlE8NMqBFt28YsmXcTrag_AWRAjfjCVuaIR2Jg2GpJx-FYhZ93bHEAfwY9gRfyMNOHJvAX5h_lX-5kuesCBWUGHVx8x34P8qEz2lJBBfPJh5X6l72tXTieNI5BBLlvFQGnT8p6UkhUfoSE6L4CmRKqKNUn7XjDJ-v6qQ3XLm5dSGZyPDf9f6D6zMx4ESKOh1-M-ElIbnpiqrIWTgOJZB7Ji5J5d8AeWIGt8ou0kbUjqQqkwQje4HkM-gpPTkkcSXJGb4a8CeH7eEBoJ_NfZWzXUGA5u-4U2pxWplYveTTQvBY0cmQ8Mrcf68rZKlTRtPU7tKwPj2BpfVUQHSBnnu33o2zHxAdmJazZqg-RB0UJcjk14Tazyqr-YhDBuH5R3QH9xcYM-qmIumt5_GGBCg_73N6VK5yy0HKb9Coi0dqaatMPp_vblQN6ZeMe9GIelZhgNc2s4nPOdswyZaL5QsaLql_JKJ63xPxzYQnSaCcj-qiSTbtPCInylQWN40nP9qm1E_mGZAYoGwXa0NfWsg-U4-D30pKwRFIE-yPu7IIYhQMu78H0P-Ez3xMECg-mEzV532xlyBuVQfVDM1WiRuwk2A3LQOsXMaTEPIOwxW9DtpfFiwliy_YSlt0v3Hz9EPwwLJhL0EvlRDrMVCcNZuz7Nv9G4AQk9Ku5zJzRa_VDkN7MciEbLL1EHoyD2K6BH--htn9MubnnjI_T6915BTZm04va9SMJP2fRvSnSwq28

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
