# SubLangLearn

Watch a movie with foreign audio track, play with subtitles - and learn new language.

It is a simple golang webserver that controls VLC Player via Remote Control. A uses reads subtitles in browser and can jump to previous or next phrases easy.

![screenshot](https://raw.github.com/gophergala2016/SubLangLearn/master/static/SubLangLearn_sceenshot.png)

## Prerequisites

VLC Player.

Tested only on Windows 10 x64.

Optional "config.ini" file (sample file is in sources):

path to VLC player:

`vlc_path=C:\Program Files (x86)\VideoLAN\VLC\vlc.exe`

tcp port of VLC Player Remote Control:

`vlc_port=3016`

tcp port of golang webserver:

`web_port=2016`

## Start

Usage:
`SubLangLearn.exe <movie_path> <subtitles_path>`

Example:
`SubLangLearn.exe "D:\FRIENDS\Season 02\02x01 - The One With Ross' New Girlfriend.avi" "D:\FRIENDS\Season 02\02x01 - The One With Ross' New Girlfriend.srt"`

!!! Sometimes start is crashed - repeat please.

Subtitles should be in .srt format and utf-8 encoded.

Open page "localhost:2016" (`web_port` from config.ini)

## Play

Seek phrases, repeat complex sentences, listen to difficult words, ...

## Stop

Close VLC Player.
Stop SubLangLearn.exe (Ctrl-C).