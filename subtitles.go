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
	StartPosition int
	FinishPosition int
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
	Text string
	IsStartVisible bool

	start time.Time
}

var (
	htmlTags = []string {"b", "i", "u"}
	bom = "\xEF\xBB\xBF"
)

func (line *SubtitleLine) toJsonLine(index int) *JsonSubtitleLine {
	return &JsonSubtitleLine{Index:index, Start:line.Start.Format("15:04:05"), start:line.Start, Text:strings.Join(line.Text, "<br>")}
}

func (subs *Subtitles) toJsonSubtitles() *JsonSubtitles {
	result := &JsonSubtitles{FileName:filepath.Base(subs.FilePath), Lines:make([]*JsonSubtitleLine, len(subs.Lines))}
	var lastStart time.Time
	for i, line := range subs.Lines {
		result.Lines[i] = line.toJsonLine(i)
		if lastStart.IsZero() || result.Lines[i].start.Sub(lastStart) > (time.Second * 10) {
			result.Lines[i].IsStartVisible = true
			lastStart = result.Lines[i].start
		}
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
			currentLine.StartPosition = GetPosition(currentLine.Start)
			currentLine.FinishPosition = GetPosition(currentLine.Finish)
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

func GetPosition(moment time.Time) int {
	parts := strings.SplitN(moment.Format("15:04:05"), ":", 3)
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	seconds, _ := strconv.Atoi(parts[2])
	return hours*3600 + minutes*60 + seconds + 1
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