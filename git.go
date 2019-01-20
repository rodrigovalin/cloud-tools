package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func findRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	root, err := findInDir(cwd)
	if err != nil {
		return "", err
	}

	return root, nil
}

func findInDir(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() && f.Name() == ".git" {
			return dir, nil
		}
	}

	if path.Dir(dir) == "" || path.Dir(dir) == "/" {
		return "", fmt.Errorf("Could not find Git root directory.")
	}
	return findInDir(path.Dir(dir))
}
