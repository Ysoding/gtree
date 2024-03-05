package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func sortEntryByName(entries []fs.DirEntry) []fs.DirEntry {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	return entries
}

func printTree(dirPath string, level int, prefix string) {
	indent := strings.Repeat(" ", level*3)
	if level > 0 {
		prefix = prefix + indent + "|"
	}

	tmpPrefix := prefix

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println(prefix + "─ ─ " + dirPath + " (Permission denied)")
		return
	}

	if level == 0 {
		fmt.Println(dirPath)
	}

	entries = sortEntryByName(entries)

	for i, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		isDir := entry.IsDir()
		fi, err := entry.Info()
		if err != nil {
			continue
		}

		isLink := fi.Mode()&os.ModeSymlink != 0

		var entryStr string
		if isLink {
			target, _ := os.Readlink(entryPath)
			entryStr = fmt.Sprintf("%s -> %s", entry.Name(), target)
		} else {
			entryStr = entry.Name()
		}

		isLast := i == len(entries)-1
		if isLast { // last element
			prefix = prefix[:len(prefix)-1] + "`"
		}

		if isDir {
			fmt.Println(prefix + "─ ─ " + entryStr)
			if isLast {
				printTree(entryPath, level+1, tmpPrefix[:len(prefix)-1])
			} else {
				printTree(entryPath, level+1, tmpPrefix)
			}
		} else {

			fmt.Println(prefix + "─ ─ " + entryStr)
		}
	}
}

func main() {
	var targetPath string

	if len(os.Args) != 2 {
		targetPath = "."
	} else {
		targetPath = os.Args[1]
	}

	_, err := os.Stat(targetPath)

	if os.IsNotExist(err) {
		fmt.Println("Error: target dir path not exists.")
		os.Exit(0)
	}

	printTree(targetPath, 0, "|")
}
