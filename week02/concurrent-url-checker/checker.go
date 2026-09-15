package main

import (
	"net/url"
	"os"
	"strings"
)

func GetURLList(fileName string) ([]string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		//err reading file
		return nil, err
	}
	newContent := strings.Split(string(content), "\n")
	for index, url := range newContent {
		url = strings.TrimSpace(url)
		if !strings.Contains(url, "://") {
			url = "https://" + url
		}
		newContent[index] = url
	}
	return newContent, nil
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") &&
		parsedURL.Hostname() != ""
}
