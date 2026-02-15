package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mon-go/internal/config"
	"mon-go/internal/handler"
	"mon-go/internal/server"
	"mon-go/internal/store"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	db, err := store.NewMongoDB(ctx, cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Printf("mongo: %v (server starting without database)", err)
		db = nil
	}
	defer func() {
		if db != nil {
			if err := db.Close(context.Background()); err != nil {
				log.Printf("mongo close: %v", err)
			}
		}
	}()

	var itemStore *store.ItemStore
	if db != nil {
		itemStore = store.NewItemStore(db)
	}
	itemHandler := &handler.ItemHandler{Store: itemStore}

	srv := server.New(cfg.ServerPort, db, itemHandler)

	go func() {
		log.Printf("server listening on :%d", cfg.ServerPort)
		if err := srv.Run(); err != nil && err != context.Canceled {
			log.Printf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
