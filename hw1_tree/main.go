package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	var result string
	var prefix string
	err := readDirOrFiles(path, &result, prefix, printFiles)
	fmt.Println(result)
	out.Write([]byte(result))
	return err
}

func readDirOrFiles(path string, result *string, prefix string, printFiles bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if !printFiles {
		onlyDirs := files[:0]
		for _, f := range files {
			if f.IsDir() {
				onlyDirs = append(onlyDirs, f)
			} 
		}
		files = onlyDirs
	}
	for index, f := range files {
		var step string
		var prefixAdd string
		if index == len(files)-1 {
			step = "└───"
			prefixAdd = "\t"
		} else {
			step = "├───"
			prefixAdd = "│\t"
		}

		if f.IsDir() {
			*result += prefix + step + f.Name() + "\n"
			err = readDirOrFiles(path+"/"+f.Name(), result, prefix+prefixAdd, printFiles)
			if err != nil {
				return err
			}
		} else {
			if printFiles {
				if f.Size() == 0 {
					*result += fmt.Sprintf("%s%s%s (empty)\n", prefix, step, f.Name())
				} else {
					*result += fmt.Sprintf("%s%s%s (%db)\n", prefix, step, f.Name(), f.Size())
				}
			}
		}
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
