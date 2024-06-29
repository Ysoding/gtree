package internal

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var directoryCnt = 0
var fileCnt = 0

func Run(cmd *cobra.Command, args []string) {
	var targetPath string

	if len(os.Args) != 2 {
		targetPath = "."
	} else {
		targetPath = os.Args[1]
	}

	info, err := os.Stat(targetPath)

	defer printCountInfo()

	if err != nil || !info.IsDir() {
		if os.IsNotExist(err) {
			fmt.Println("Error: target dir path not exists.")
		} else {
			fmt.Printf("%s\t[error opening dir]\n", targetPath)
		}
		return
	}

	printTree(targetPath, 0, "|")
}

func printCountInfo() {
	fmt.Printf("\n%d directories, %d files\n", directoryCnt, fileCnt)
}

func sortEntryByName(entries []fs.DirEntry) []fs.DirEntry {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	return entries
}

func isIgnoreName(name string) bool {
	return name[0] == '.'
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
		if err != nil || isIgnoreName(fi.Name()) {
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
			directoryCnt++
			randomColorPrintln(prefix + "─ ─ " + entryStr)
			if isLast {
				printTree(entryPath, level+1, tmpPrefix[:len(prefix)-1])
			} else {
				printTree(entryPath, level+1, tmpPrefix)
			}
		} else {
			randomColorPrintln(prefix + "─ ─ " + entryStr)
			fileCnt++
		}
	}
}

func randomColorPrintln(a ...any) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number between 0 and 99
	randomNumber := r.Intn(100)

	if randomNumber < 80 {
		color.Unset()
		fmt.Println(a...)
	} else {
		randomNumber = r.Intn(5)

		var c *color.Color
		switch randomNumber {
		case 0:
			c = color.New(color.FgRed)
		case 1:
			c = color.New(color.FgBlue)
		case 2:
			c = color.New(color.FgGreen)
		case 3:
			c = color.New(color.FgYellow)
		case 4:
			c = color.New(color.FgMagenta)
		}
		c.Println(a...)
	}
}
