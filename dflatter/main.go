package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"unicode"
)

func main() {
	s := make([]string, 0, 128)
	s = getBaseImage(s, os.Args[1])

	for i := len(s) - 1; i >= 0; i-- {
		log.Printf("Merging %s artifacts\n", s[i])
		if s[i] != os.Args[1] {
			err := cpyDir(path.Join(os.Getenv("HOME"), "eventbrite/docker-dev", s[i]), path.Join("./bundle/", os.Args[1], "/base"))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := cpyDir(path.Join(os.Getenv("HOME"), "eventbrite/docker-dev", s[i]), path.Join("./bundle/", os.Args[1]))
			if err != nil {
				fmt.Println(err)
			}
		}
		appendFiles(path.Join(os.Getenv("HOME"), "eventbrite/docker-dev", s[i], "Dockerfile"), path.Join("./bundle/", os.Args[1], "Dockerfile"))

	}
	fmt.Print("\n\n")
	cleanDockerfile(os.Args[1])
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
		return getBaseImage(s, spaceMap(string(strings.Split(newFolder[1], ":")[0])))

	}
	//s = append(s, string(next))
	return s
}

func spaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func cpyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = cpyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if fd.Name() != "Dockerfile" {
				if err = File(srcfp, dstfp); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return nil
}

//File copies a single file from src to dst
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

//appendFiles appends the content of inName at the end of outName
func appendFiles(inName, outName string) {

	in, err := os.Open(inName)
	if err != nil {
		log.Fatalln("failed to open second file for reading:", err)
	}
	defer in.Close()

	out, err := os.OpenFile(outName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("failed to dest  file for writing:", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Fatalln("failed to append second file to first:", err)
	}
	log.Printf("Appended %s to the end of %s\n", inName, outName)

	in.Close()
	out.Close()
}

func cleanDockerfile(folder string) {
	file, err := os.Open(path.Join("./bundle/", folder, "Dockerfile"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "FROM") || strings.HasPrefix(scanner.Text(), "from") || strings.HasPrefix(scanner.Text(), "MAINTAINER") {
			fmt.Println(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
