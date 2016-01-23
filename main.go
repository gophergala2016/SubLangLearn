package main

import (
	"log"
	"time"
	"fmt"
	"math/rand"
)

func main() {
	player := NewVlcPlayer(`C:\Program Files (x86)\VideoLAN\VLC\vlc.exe`, "localhost", 2016)
	err := player.Start()
	if err != nil {
		log.Fatalf("Failed to start VLC Player: ", err)
	}
	player.PlayMovie(`D:\USERS\FALCON\_video\Sample.avi`)
	for {
		time.Sleep(3000*time.Millisecond)
		fmt.Println(player.CurrentPosition())
		player.Seek(rand.Int() % 3600)
	}
	player.Wait()
}
