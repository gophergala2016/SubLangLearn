package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"strconv"
	"encoding/json"
	"strings"
	"path/filepath"
	"fmt"
	"os"
)

var (
	upgrader = websocket.Upgrader{}
	subtitles *Subtitles
	player *VlcPlayer
	moviePath string
	subtitlesPath string
	socketConn *websocket.Conn
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <movie_path> <subtitles_path>\n", filepath.Base(os.Args[0]))
		return
	}
	if strings.HasSuffix(strings.ToLower(filepath.Base(os.Args[1])), ".srt") {
		moviePath, subtitlesPath = os.Args[2], os.Args[1]
	} else {
		moviePath, subtitlesPath = os.Args[1], os.Args[2]
	}

	config := LoadConfig(filepath.Join(".", "config.ini"))
	launchVlcPlayer(config.VlcPlayerPath, config.VlcPlayerTcpPort)
	player.PlayMovie(moviePath)
	subtitles, _ = ParseSubtitlesFile(subtitlesPath)
	port := config.WebServerHttpPort
	http.HandleFunc("/socket", socket)
	http.HandleFunc("/getSubtitles", getSubtitles)
	http.HandleFunc("/play", play)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":" + strconv.Itoa(port), nil)
}

func launchVlcPlayer(vlcPath string, vlcPort int) {
	player = NewVlcPlayer(vlcPath, "localhost", vlcPort)
	err := player.Start()
	if err != nil {
		log.Fatal("Failed to start VLC Player: ", err)
	}
	go player.Run()
}

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
	content, _ := json.Marshal(subtitles.toJsonSubtitles())
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

var lastIndex = -1
var slowSpeed = false

func play(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(r.FormValue("Index"))
	if subtitles == nil || index >= len(subtitles.Lines) {
		return
	}

	if index == lastIndex {
		lastIndex = -1
		slowSpeed = true
		player.SlowSpeed()
	} else {
		lastIndex = index
		if slowSpeed {
			slowSpeed = false
			player.NormalSpeed()
		}
	}
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