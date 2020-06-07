package config

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	// Tags is set for other package
	Tags []string
)

func init() {
	makeTags()
}

// MakeTags for other package use
func makeTags() {
	tags := readConfigTags()
	scanner := bufio.NewScanner(strings.NewReader(tags))
	for scanner.Scan() {
		rawText := scanner.Text()
		tmp := strings.Split(rawText, ",")
		Tags = append(Tags, tmp...)
	}
}

func readConfigTags() string {
	var err error
	var data []byte
	configDir := os.Getenv("ConfigDir")
	tagsConfigPath := path.Join(configDir, "tagconfig.txt")
	_, err = os.Stat(tagsConfigPath)
	if os.IsNotExist(err) {
		f, _ := os.OpenFile(tagsConfigPath, os.O_CREATE|os.O_RDWR, 0644)
		defaultContent := "tag1,tag2"
		f.WriteString(defaultContent)
		f.Close()
		return defaultContent
	}

	data, _ = ioutil.ReadFile(tagsConfigPath)
	return string(data)
}
