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
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
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
	player.PlayMovie(`D:\USERS\FALCON\_video\Frozen.avi`)
	//player.Wait()
}

func socket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Failed to upgrade http connection to support websockects: ", err)
		return
	}
	defer c.Close()
	c.WriteMessage(websocket.TextMessage, []byte("Hello"))
//	for {
//		messageType, message, err := c.ReadMessage()
//		log.Printf("recv: %v %s", messageType, message, err)
//	}
}

func getSubtitles(w http.ResponseWriter, r *http.Request) {
	subtitles, _ = ParseSubtitlesFile(`D:\USERS\FALCON\_video\Frozen 2013 720p EN.srt`)
	content, _ := json.Marshal(subtitles.toJsonSubtitles())
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func play(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(r.FormValue("Index"))
	line := subtitles.Lines[index]
	parts := strings.SplitN(line.Start.Format("15:04:05"), ":", 3)
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	seconds, _ := strconv.Atoi(parts[2])
	position := hours*3600 + minutes*60 + seconds + 1
	player.Seek(position)
}
