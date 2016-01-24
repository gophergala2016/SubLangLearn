package main
import (
	"fmt"
	"strings"
)

type PlayerCommand interface {
	Execute(player *VlcPlayer) string
}

type simpleCommand struct {
	command string
	responseLineCount int
}

func (command simpleCommand) Execute(player *VlcPlayer) string {
	return player.execCommand(command.command, command.responseLineCount)
}

type seekCommand struct {
	position int
}

func (command seekCommand) Execute(player *VlcPlayer) string {
	output := player.execCommand(fmt.Sprintf("seek %d", command.position), 1)
	if !strings.HasPrefix(output, "seek:") {
		player.execCommand("pause", 2)
		output = player.execCommand(fmt.Sprintf("seek %d", command.position), 1)
	}
	return output
}
