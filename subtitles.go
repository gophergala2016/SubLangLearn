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

var (
	htmlTags = []string {"b", "i", "u"}
	bom = "\xEF\xBB\xBF"
)

func ParseSubtitlesFile(filePath string) (*Subtitles, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	subtitles := &Subtitles{FilePath:filePath}
	var currentLine *SubtitleLine
	scanner := bufio.NewScanner(file)
	isFirstLine := true
	for scanner.Scan() {
		rawLine := strings.TrimSpace(scanner.Text())
		if isFirstLine {
			rawLine = strings.TrimPrefix(rawLine, bom)
			isFirstLine = false
		}

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
			currentLine.Text = append(currentLine.Text, removeHtml(rawLine))
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

func removeHtml(value string) string {
	if !strings.Contains(value, "<") {
		return value
	}

	for _, tag := range htmlTags {
		value = strings.Replace(value, "<" + tag + ">", "", -1)
		value = strings.Replace(value, "</" + tag + ">", "", -1)
		value = strings.Replace(value, "{" + tag + "}", "", -1)
		value = strings.Replace(value, "{/" + tag + "}", "", -1)
	}
	value = strings.Replace(value, "</font>", "", -1)
	fontIndex := strings.Index(value, "<font")
	if fontIndex >= 0 {
		endIndex := strings.Index(value[fontIndex:], ">")
		if endIndex >= 0 {
			value = value[:fontIndex] + value[fontIndex+endIndex+1:]
		}
	}
	return value
}