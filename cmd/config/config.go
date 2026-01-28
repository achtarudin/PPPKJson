package config

import (
	"context"
	"cutbray/pppk-json/internal/ports"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type ConnectManager struct {
	Name    string
	Adapter ports.AdapterPort
}

func DisconnectAdapters(adapters ...ConnectManager) {
	var wgStop sync.WaitGroup

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, adapter := range adapters {
		wgStop.Add(1)
		go func(connect ConnectManager) {
			defer wgStop.Done()
			log.Printf("[Info %s] Stopping...", connect.Name)

			if err := connect.Adapter.Disconnect(timeoutCtx); err != nil {
				log.Printf("[Error] Failed to stop adapter: %v", err)
			} else {
				log.Printf("[Info %s] Stopped successfully.", connect.Name)
			}

		}(adapter)
	}
	wgStop.Wait()
	os.Exit(0)
}

func ConnectAdapters(ctx context.Context, adapters ...ConnectManager) error {
	var wg sync.WaitGroup

	// Provide a buffered channel based on the number of adapters
	errChan := make(chan error, len(adapters))

	// Timeout context for connection attempts
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for _, item := range adapters {
		wg.Add(1)

		go func(connect ConnectManager) {
			defer wg.Done()
			log.Printf("[%s] Connecting...", connect.Name)
			if err := connect.Adapter.Connect(timeoutCtx); err != nil {
				errChan <- fmt.Errorf("[%s] Error :%w", connect.Name, err)
			} else {
				log.Printf("[%s] Connected successfully.", connect.Name)
			}
		}(item)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
