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
	commands chan playerCommand

	Progress chan int
	position int
}

type playerCommand struct {
	command string
	responseLineCount int
}

func NewVlcPlayer(exePath string, tcpHost string, tcpPort int) *VlcPlayer {
	return &VlcPlayer{exePath:exePath, tcpHost:tcpHost, tcpPort:tcpPort, commands:make(chan playerCommand), Progress:make(chan int)}
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
	player.commands <- playerCommand{command:fmt.Sprintf("add %s", moviePath), responseLineCount:2}
}

func (player *VlcPlayer) Seek(position int) {
	player.commands <- playerCommand{command:fmt.Sprintf("seek %d", position), responseLineCount:1}
}

func (player *VlcPlayer) Play(position int) {
	player.commands <- playerCommand{command:fmt.Sprintf("seek %d", position), responseLineCount:1}
}

func (player *VlcPlayer) Run() {
	for {
		select {
		case command := <-player.commands:
			fmt.Println(command)
			player.execCommand(command)
		case <-time.After(time.Millisecond * 100):
			positionResponse := player.execCommand(playerCommand{command:"get_time", responseLineCount:1})
			positionResponse = strings.TrimSpace(positionResponse)
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

func (player *VlcPlayer) execCommand(command playerCommand) string {
	fmt.Fprintln(player.conn, command.command)
	output := ""
	count := 0
	for count < command.responseLineCount {
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
	return output
}