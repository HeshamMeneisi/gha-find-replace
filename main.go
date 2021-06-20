package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"encoding/json"
)

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

func listFiles(include string, exclude string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if doesFileMatch(path, include, exclude) {
			fileList = append(fileList, path)
		}
		return nil
	})
	return fileList, err
}

func doesFileMatch(path string, include string, exclude string) bool {
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		includeRe := regexp.MustCompile(include)
		excludeRe := regexp.MustCompile(exclude)
		return includeRe.Match([]byte(path)) && !excludeRe.Match([]byte(path))
	}
	return false
}

func findAndReplace(path string, find string, replace string) (bool, error) {
	if find != replace {
		read, readErr := ioutil.ReadFile(path)
		check(readErr)

		re := regexp.MustCompile(find)
		newContents := re.ReplaceAllString(string(read), replace)

		if newContents != string(read) {
			writeErr := ioutil.WriteFile(path, []byte(newContents), 0)
			check(writeErr)
			fmt.Println(fmt.Sprintf(`Replaced %s with %s in %s`, find, replace, path))
			return true, nil
		}
		fmt.Println(fmt.Sprintf(`Failed to replace %s with %s in %s`, find, replace, path))
	}

	return false, nil
}

func replaceSimple(files []string, find string, replace string) int {
	modifiedCount := 0
	fmt.Println(fmt.Sprintf(`Replacing %s with %s`, find, replace))
	for _, path := range files {
		modified, findAndReplaceErr := findAndReplace(path, find, replace)
		check(findAndReplaceErr)

		if modified {
			modifiedCount += 1
		}
	}
	return modifiedCount
}

func replaceMapping(files []string, mapping map[string]string) int {
	modifiedCount := 0
	for key, value := range mapping {
		modifiedCount += replaceSimple(files, key, value)
	}
	return modifiedCount
}

func main() {
	include := os.Getenv("INPUT_INCLUDE")
	exclude := os.Getenv("INPUT_EXCLUDE")
	find := os.Getenv("INPUT_FIND")
	replace := os.Getenv("INPUT_REPLACE")
	mapping_json := os.Getenv("INPUT_MAPPING")

	files, filesErr := listFiles(include, exclude)
	check(filesErr)

	modifiedCount := 0
	if mapping_json != "" {
		var mapping map[string]string
		json.Unmarshal([]byte(mapping_json), &mapping)
		fmt.Println(fmt.Sprintf(`Replacing according to mapping %s`, mapping_json))
		modifiedCount = replaceMapping(files, mapping)
	} else {
		modifiedCount = replaceSimple(files, find, replace)
	}

	fmt.Println(fmt.Sprintf(`::set-output name=modifiedFiles::%d`, modifiedCount))
}
