package main

import (
	"time"
	"os"
	"bufio"
	"strings"
	"strconv"
	"path/filepath"
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

type JsonSubtitles struct {
	FileName string
	Lines []*JsonSubtitleLine
}

type JsonSubtitleLine struct {
	Index int
	Start string
	Finish string
	Text string
}

var (
	htmlTags = []string {"b", "i", "u"}
	bom = "\xEF\xBB\xBF"
)

func (line *SubtitleLine) toJsonLine(index int) *JsonSubtitleLine {
	return &JsonSubtitleLine{Index:index, Start:line.Start.Format("15:04:05"), Finish:line.Finish.Format("15:04:05"), Text:strings.Join(line.Text, "<br>")}
}

func (subs *Subtitles) toJsonSubtitles() *JsonSubtitles {
	result := &JsonSubtitles{FileName:filepath.Base(subs.FilePath), Lines:make([]*JsonSubtitleLine, len(subs.Lines))}
	for i, line := range subs.Lines {
		result.Lines[i] = line.toJsonLine(i)
	}
	return result
}

func ParseSubtitlesFile(filePath string) (*Subtitles, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	subtitles := &Subtitles{FilePath:filePath}
	var previousLine, currentLine *SubtitleLine
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
				if len(subtitles.Lines) == 0 || subtitles.Lines[len(subtitles.Lines)-1] != currentLine {
					subtitles.Lines = append(subtitles.Lines, currentLine)
				}
				previousLine = currentLine
				currentLine = nil
			}
			continue
		}
		if currentLine == nil {
			index, _ := strconv.Atoi(rawLine)
			currentLine = &SubtitleLine{Index:index}
		} else if strings.Contains(rawLine, "-->") {
			parts := strings.SplitN(rawLine, "-->", 2)
			start := parseTime(parts[0])
			finish := parseTime(parts[1])
			if previousLine != nil && start.Sub(previousLine.Finish) < time.Duration(50 * time.Millisecond) {
				currentLine = previousLine
				currentLine.Finish = finish
			} else {
				currentLine.Start = start
				currentLine.Finish = finish
			}
		} else {
			currentLine.Text = append(currentLine.Text, removeHtml(rawLine))
		}
	}
	if currentLine != nil {
		if subtitles.Lines[len(subtitles.Lines)-1] != currentLine {
			subtitles.Lines = append(subtitles.Lines, currentLine)
		}
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