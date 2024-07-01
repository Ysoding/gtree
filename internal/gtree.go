package internal

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Counter struct {
	dirs  int
	files int
}

func (c *Counter) count(isDir bool) {
	if isDir {
		c.dirs++
	} else {
		c.files++
	}
}

func (c Counter) String() string {
	return fmt.Sprintf("\n%d directories, %d files\n", c.dirs, c.files)
}

func Run(cmd *cobra.Command, args []string) {
	var targetPath string

	if len(os.Args) != 2 {
		targetPath = "."
	} else {
		targetPath = os.Args[1]
	}

	counter := new(Counter)
	info, err := os.Stat(targetPath)

	defer fmt.Print(counter)

	if err != nil || !info.IsDir() {
		if os.IsNotExist(err) {
			fmt.Println("Error: target dir path not exists.")
		} else {
			fmt.Printf("%s\t[error opening dir]\n", targetPath)
		}
		return
	}

	fmt.Println(targetPath)
	printTree(counter, targetPath, "")
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

func printTree(counter *Counter, dirPath string, prefix string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println(prefix + dirPath + " (Permission denied)")
		return
	}

	entries = sortEntryByName(entries)

	for i, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		isDir := entry.IsDir()

		fi, err := entry.Info()
		if err != nil || isIgnoreName(fi.Name()) {
			continue
		}

		counter.count(isDir)
		isLink := fi.Mode()&os.ModeSymlink != 0

		var entryStr string

		if isLink {
			target, _ := os.Readlink(entryPath)
			entryStr = fmt.Sprintf("%s -> %s", entry.Name(), target)
		} else {
			entryStr = entry.Name()
		}

		isLast := i == len(entries)-1

		if isLast {
			fmt.Print(prefix + "└── ")
		} else {
			fmt.Print(prefix + "├── ")
		}

		randomColorPrintln(entryStr)

		if isDir {
			if isLast {
				printTree(counter, entryPath, prefix+"    ")
			} else {
				printTree(counter, entryPath, prefix+"│   ")
			}
		}
	}
}

func randomColorPrintln(a ...any) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

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
