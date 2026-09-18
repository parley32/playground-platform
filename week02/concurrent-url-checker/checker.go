package main

import (
	"net/url"
	"os"
	"strings"
)

func getURLList(fileName string) ([]string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	newContent := strings.Split(string(content), "\n")
	for index, url := range newContent {
		newContent[index] = strings.TrimSpace(url)
	}
	return purgeEmptyLines(newContent), nil
}

func purgeEmptyLines(list []string) []string {
	newList := []string{}
	for _, item := range list {
		if item != "" {
			newList = append(newList, item)
		}
	}
	return newList
}

func normalizeURL(url string) string {
	if !strings.Contains(url, "://") {
		url = "https://" + url
	}
	return url
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") &&
		parsedURL.Hostname() != ""
}
