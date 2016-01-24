package main

import (
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"strconv"
	"encoding/json"
	"strings"
)

var (
	upgrader = websocket.Upgrader{}
	subtitles *Subtitles
	player *VlcPlayer
)

func main() {
	main2()
	port := 3016
	http.HandleFunc("/socket", socket)
	http.HandleFunc("/getSubtitles", getSubtitles)
	http.HandleFunc("/play", play)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":" + strconv.Itoa(port), nil)
}

func main3() {
	subtitles, err := ParseSubtitlesFile(`D:\USERS\FALCON\_video\Frozen.srt`)
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range subtitles.Lines {
		fmt.Println(line.Index)
		fmt.Printf("%s --> %s\n", line.Start.Format("15:04:05"), line.Finish.Format("15:04:05"))
		for _, text := range line.Text {
			fmt.Println(text)
		}
		fmt.Println()
	}
}

func main2() {
	player = NewVlcPlayer(`C:\Program Files (x86)\VideoLAN\VLC\vlc.exe`, "localhost", 2016)
	err := player.Start()
	if err != nil {
		log.Fatalf("Failed to start VLC Player: ", err)
	}
	go player.Run()
	player.PlayMovie(`D:\USERS\FALCON\_Downloads\FRIENDS (eng+rus+subs)\Season 01\01x01 - The One Where Monica Gets A New Roomate.avi`)
}

var socketConn *websocket.Conn

func socket(w http.ResponseWriter, r *http.Request) {
	var err error
	socketConn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Failed to upgrade http connection to support websockects: ", err)
		return
	}
	go sendProgress()
}

func getSubtitles(w http.ResponseWriter, r *http.Request) {
	subtitles, _ = ParseSubtitlesFile(`D:\USERS\FALCON\_Downloads\FRIENDS (eng+rus+subs)\Season 01\01x01 - The One Where Monica Gets A New Roomate.srt`)
	content, _ := json.Marshal(subtitles.toJsonSubtitles())
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func play(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(r.FormValue("Index"))
	//shift, _ := strconv.Atoi(r.FormValue("Shift"))
	line := subtitles.Lines[index]
	parts := strings.SplitN(line.Start.Format("15:04:05"), ":", 3)
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	seconds, _ := strconv.Atoi(parts[2])
	position := hours*3600 + minutes*60 + seconds
	player.Seek(position)
}

var lastSendIndexes []int

func sendProgress() {
	for {
		select {
		case position := <-player.Progress:
			indexes := make([]int, 0, 10)
			nextIndex := -1
			for i, line := range subtitles.Lines {
				if line.StartPosition <= position && (position < line.FinishPosition || (line.StartPosition == position && position == line.FinishPosition))  {
					indexes = append(indexes, i)
				} else if position < line.StartPosition {
					nextIndex = i
					break
				}
			}
			if len(indexes) == 0 && nextIndex >= 0 {
				if subtitles.Lines[nextIndex].StartPosition > position + 1 {
					indexes = append(indexes, -nextIndex - 1)
				} else {
					indexes = append(indexes, nextIndex)
				}
			}
			if len(indexes) > 0 {
				if !isEqualSlices(indexes, lastSendIndexes) {
					lastSendIndexes = indexes
					socketConn.WriteJSON(indexes)
				}
			}
		}
	}
}

func isEqualSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, item := range a {
		if b[i] != item {
			return false
		}
	}
	return true
}