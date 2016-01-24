package main

import (
	"os/exec"
	"fmt"
	"net"
	"bufio"
	"strings"
	"strconv"
	"time"
)

type VlcPlayer struct {
	exePath string
	tcpHost string
	tcpPort int

	cmd *exec.Cmd
	conn net.Conn
	connReader *bufio.Reader
	commands chan PlayerCommand

	Progress chan int
	position int
}

func NewVlcPlayer(exePath string, tcpHost string, tcpPort int) *VlcPlayer {
	return &VlcPlayer{exePath:exePath, tcpHost:tcpHost, tcpPort:tcpPort, commands:make(chan PlayerCommand), Progress:make(chan int)}
}

func (player *VlcPlayer) Start() error {
	host := player.tcpHost
	port := player.tcpPort
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 2016
	}
	player.cmd = exec.Command(player.exePath, "--extraintf=rc", fmt.Sprintf("--rc-host=%s:%d", host, port), "--one-instance") //"--rc-quiet",
	err := player.cmd.Start()
	if err != nil {
		return err
	}
	player.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	player.connReader = bufio.NewReader(player.conn)
	return nil
}

func (player *VlcPlayer) PlayMovie(moviePath string) {
	player.commands <- simpleCommand{command:fmt.Sprintf(`add "%s"`, moviePath), responseLineCount:2}
	player.commands <- simpleCommand{command:"strack -1", responseLineCount:1}
}

func (player *VlcPlayer) Seek(position int) {
	player.commands <- seekCommand{position:position}
}

func (player *VlcPlayer) SlowSpeed() {
	player.commands <- simpleCommand{command:fmt.Sprintf(`slower`), responseLineCount:1}
}

func (player *VlcPlayer) NormalSpeed() {
	player.commands <- simpleCommand{command:fmt.Sprintf(`normal`), responseLineCount:1}
}

func (player *VlcPlayer) Run() {
	for {
		select {
		case command := <-player.commands:
			command.Execute(player)
		case <-time.After(time.Millisecond * 100):
			positionResponse := simpleCommand{command:"get_time", responseLineCount:1}.Execute(player)
			position, _ := strconv.Atoi(positionResponse)
			if position != player.position {
				select {
				case player.Progress <- position:
					player.position = position
				default:
				}
			}
		}
	}
}

func (player *VlcPlayer) execCommand(command string, responseLineCount int) string {
	fmt.Fprintln(player.conn, command)
	output := ""
	count := 0
	for count < responseLineCount {
		line, err := player.connReader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "status change:") {
			continue
		}
		output += line
		count++
	}
	return strings.TrimSpace(output)
}