package subutils

import (
	"bytes"
	"fmt"
	"log"
)

const metadata = `
[Script Info]
Title: MovieTitle
Original Script: simple-sub
ScriptType: v1.00
Collisions: Normal

[V4 Styles]
Format: Name,Fontname,Fontsize,PrimaryColour,SecondaryColour,TertiaryColour,BackColour,Bold,Italic,BorderStyle,Outline,Shadow,Alignment,MarginL,MarginR,MarginV,AlphaLevel,Encoding
Style: StyleA,Arial,20,&H00FFFFFF,&H00FFFFFF,0,0,0,0,1,2,0,6,30,30,10,0,0
Style: StyleB,Arial,20,&H00FFFFFF,&H00FFFFFF,0,0,0,0,1,2,0,2,30,30,10,0,0

[Events]
Format: Marked, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text`

const subUpLineTemplate = "Dialogue: Marked=0,%s,%s,StyleA,NTP,0000,0000,0000,!Effect,%s \n"
const subDownLineTemplate = "Dialogue: Marked=0,%s,%s,StyleB,NTP,0000,0000,0000,!Effect,%s \n"

// MergeSubtitles : merges subtitles and creates a new file.
func MergeSubtitles(c CommandArgs) {
	if hasAllSubMergeParams(c) {
		subTop, subBottom := createEntries(c)
		merged := Merge(subTop, subBottom)
		writeToFile(c.FileSubTop+".merged.ssa", merged)
	} else {
		log.Fatal("Missing merge params!")
	}
}

// Merge : merges two subtitles
func Merge(subUp []SubtitleEntry, subDown []SubtitleEntry) string {
	var buffer bytes.Buffer
	buffer.WriteString(metadata + "\n")
	for _, item := range subUp {
		buffer.WriteString(fmt.Sprintf(subUpLineTemplate, item.StartTime, item.EndTime, item.Text))
	}
	for _, item := range subDown {
		buffer.WriteString(fmt.Sprintf(subDownLineTemplate, item.StartTime, item.EndTime, item.Text))
	}
	return buffer.String()
}

func createEntries(c CommandArgs) ([]SubtitleEntry, []SubtitleEntry) {
	return createEntry(c.FileSubTop, c.EncSubTop), createEntry(c.FileSubBottom, c.EncSubBottom)
}

func createEntry(fileName string, enc string) []SubtitleEntry {
	if len(fileName) == 0 {
		return []SubtitleEntry{}
	}
	lines := ""
	if len(enc) == 0 {
		lines = simpleRead(fileName)
	} else {
		lines = readWithEncoding(fileName, getEncoding(enc))
	}
	return CreateSubEntries(lines)
}
