package subutils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"simple-sub/converter"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// CommandArgs : commandline args
type CommandArgs struct {
	Mode     string
	FileName string
	Encoding string
}

var validModes = map[string]func(c CommandArgs){
	"--remove-accent": removeAccentLetters,
}

var validEncodings = map[string]*charmap.Charmap{
	"pl": charmap.Windows1250,
	"tr": charmap.Windows1254,
}

// GetValidModes : returns the valid modes
func GetValidModes() []string {
	m := make([]string, 1)
	for k := range validModes {
		m = append(m, k)
	}
	return m
}

// Run : runs in accordance with commandline arguments
func (c *CommandArgs) Run() {
	if fn, ok := validModes[c.Mode]; ok {
		fn(*c)
	} else {
		fmt.Println("provide a valid mode. Valid modes : ", GetValidModes())
	}
}

func (c CommandArgs) String() string {
	return fmt.Sprintf("Mode: %s\nFileName: %s\nEncoding: %s\n", c.Mode, c.FileName, c.Encoding)
}

func removeAccentLetters(commandArgs CommandArgs) {
	if len(commandArgs.FileName) > 0 && len(commandArgs.Encoding) > 0 {
		txt := readWithEncoding(commandArgs.FileName, getEncoding(commandArgs.Encoding))
		writeToFile(commandArgs.FileName+".new", txt)
	}
}

func getEncoding(cmdStr string) *charmap.Charmap {
	if enc, ok := validEncodings[cmdStr]; ok {
		return enc
	}
	return charmap.ISO8859_1
}

func readWithEncoding(filename string, charmap *charmap.Charmap) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := transform.NewReader(f, charmap.NewDecoder())

	sc := bufio.NewScanner(r)
	allText := ""
	for sc.Scan() {
		allText += sc.Text() + "\n"
	}

	return getConvertAccentText(allText)
}

func writeToFile(fileName string, text string) {
	perm := os.FileMode(0644)
	ioutil.WriteFile(fileName, []byte(text), perm)
}

func getConvertAccentText(text string) string {
	newText := ""
	for _, runeValue := range text {
		newText += converter.Convert2NonAccent(string(runeValue))
	}
	return newText
}
