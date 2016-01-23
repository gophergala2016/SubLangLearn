package main

import (
	"time"
	"os"
	"bufio"
	"strings"
	"strconv"
)

type SubtitleLine struct {
	Index int
	Start time.Time
	Finish time.Time
	Text []string
}

type Subtitles struct {
	FilePath string
	Lines []*SubtitleLine
}

func ParseSubtitlesFile(filePath string) (*Subtitles, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	subtitles := &Subtitles{FilePath:filePath}
	var currentLine *SubtitleLine
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rawLine := strings.TrimSpace(scanner.Text())
		if rawLine == "" {
			if currentLine != nil {
				subtitles.Lines = append(subtitles.Lines, currentLine)
				currentLine = nil
			}
			continue
		}
		if currentLine == nil {
			index, _ := strconv.Atoi(rawLine)
			currentLine = &SubtitleLine{Index:index}
		} else if strings.Contains(rawLine, "-->") {
			parts := strings.SplitN(rawLine, "-->", 2)
			currentLine.Start = parseTime(parts[0])
			currentLine.Finish = parseTime(parts[1])
		} else {
			currentLine.Text = append(currentLine.Text, rawLine)
		}
	}
	if currentLine != nil {
		subtitles.Lines = append(subtitles.Lines, currentLine)
	}
	return subtitles, nil
}

func parseTime(value string) time.Time {
	value = strings.TrimSpace(value)
	parts := strings.SplitN(value, ",", 2)
	result, _ := time.Parse("15:04:05", parts[0])
	if len(parts) > 1 {
		milliseconds, _ := strconv.Atoi(parts[1])
		result = result.Add(time.Duration(milliseconds) * time.Millisecond)
	}
	return result
}