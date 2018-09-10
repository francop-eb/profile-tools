package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

const (
	list  = "├── "
	box   = "│"
	arrow = "├──▶ "
)

func main() {
	s := make([]string, 0, 128)
	s = getBaseImage(s, os.Args[1])

	c := color.New(color.FgCyan, color.Bold)
	y := color.New(color.FgGreen)
	b := color.New(color.Bold)

	c.Printf("\nBase image chain for %s:\n", os.Args[1])
	for i := len(s) - 1; i >= 0; i-- {
		y.Print(box)
		b.Print(s[i])
		if i != 0 {
			y.Print(arrow)
		} else {
			y.Print(box)

		}
	}
	fmt.Print("\n\n")

	c.Printf("Images dependent on %s:\n", os.Args[1])
	searchDeps(os.Args[1])

}

func getBaseImage(s []string, folder string) []string {
	s = append(s, folder)
	file, err := os.Open(path.Join(os.Getenv("HOME"), "eventbrite/docker-dev", folder, "Dockerfile"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	next := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "FROM") || strings.HasPrefix(scanner.Text(), "from") {
			next = strings.SplitAfter(scanner.Text(), " ")[1]
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	newFolder := strings.Split(next, "/")
	if len(newFolder) > 1 {
		return getBaseImage(s, SpaceMap(string(strings.Split(newFolder[1], ":")[0])))

	}
	s = append(s, string(next))
	return s
}

func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func searchDeps(image string) {
	config.Load(file.NewSource(
		file.WithPath("conf.yaml"),
	))
	dir := config.Get("dir").String("")

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	g := color.New(color.FgGreen)
	b := color.New(color.Bold)

	for _, f := range files {
		if f.IsDir() {
			if _, err := os.Stat(filepath.Join(dir, f.Name(), "Dockerfile")); !os.IsNotExist(err) {
				file, err := os.Open(path.Join(filepath.Join(dir, f.Name(), "Dockerfile")))
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					if strings.HasPrefix(scanner.Text(), "FROM") || strings.HasPrefix(scanner.Text(), "from") {
						if strings.Index(scanner.Text(), image) >= 0 {
							g.Print(list)
							b.Println(f.Name())
						}
						break
					}
				}
				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
