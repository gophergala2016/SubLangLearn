package main

import (
	"os/exec"
	"fmt"
	"net"
	"bufio"
	"strings"
	"strconv"
)

type VlcPlayer struct {
	exePath string
	tcpHost string
	tcpPort int

	cmd *exec.Cmd
	conn net.Conn
	connReader *bufio.Reader
}

func NewVlcPlayer(exePath string, tcpHost string, tcpPort int) *VlcPlayer {
	return &VlcPlayer{exePath:exePath, tcpHost:tcpHost, tcpPort:tcpPort}
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

func (player *VlcPlayer) Wait() {
	player.cmd.Wait()
}

func (player *VlcPlayer) PlayMovie(moviePath string) {
	player.execCommand(fmt.Sprintf("add %s", moviePath), 2)
}

func (player *VlcPlayer) Seek(position int) {
	player.execCommand(fmt.Sprintf("seek %d", position), 1)
}

func (player *VlcPlayer) CurrentPosition() int {
	positionResponse := player.execCommand("get_time", 1)
	positionResponse = strings.TrimSpace(positionResponse)
	position, _ := strconv.Atoi(positionResponse)
	return position
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
	return output
}