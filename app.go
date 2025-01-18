package main

import (
	"context"
	"csm/internal/chat"
	"encoding/base64"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// App struct
type App struct {
	ctx   context.Context
	chat  *chat.Chat
	db    *sqlx.DB
	index bleve.Index
}

func NewApp() *App {

	db := sqlx.MustConnect("sqlite3", "data.db")
	db.MustExec(schemaStmt)

	return &App{
		chat:  chat.NewChat(),
		db:    db,
		index: newIndex(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	buildIndex(a.index, a.db)
}

func (a *App) CreateImageDescriptor(path string) string {
	return a.chat.CreatePortfolioDescriptorForImage(a.ctx, path)
}

func (a *App) CreateAndSavePortfolioDescriptorForImage(id int, path string) {
	descriptor := a.chat.CreatePortfolioDescriptorForImage(a.ctx, path)
	savePortfolioDescriptor(a.index, a.db, id, descriptor)
}

func (a *App) GetImagePreview(path string) string {

	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(data)
}

func (a *App) LoadDemoSeller(id int) Profile {
	return getProfile(a.db, id)
}

func (a *App) ReturnTagsForImage(path string) []string {
	return a.chat.ReturnTagsForImage(a.ctx, path)
}

func (a *App) SaveNewSellerTag(id int, tag string) {
	saveSpecialtyTag(a.index, a.db, id, tag)
}

func (a *App) SearchFromImageDescriptor(path string) []Profile {
	descriptor := a.chat.CreatePortfolioDescriptorForImage(a.ctx, path)
	profiles := search(a.index, a.db, descriptor)

	return profiles
}

type BuyerChatResponse struct {
	Followup string    `json:"followup"`
	Profiles []Profile `json:"profiles"`
}

func (a *App) SendBuyerChatMessage(input string) BuyerChatResponse {
	proposedQuery, followup := a.chat.BuyerQuery(a.ctx, input)
	profiles := search(a.index, a.db, proposedQuery)

	return BuyerChatResponse{Followup: followup,
		Profiles: profiles}
}

func (a *App) SendSellerChatMessage(input string) []string {
	return a.chat.ReturnTagsForString(a.ctx, input)
}
