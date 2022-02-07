package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
)

type FileItem struct {
	fullPath    string
	sizeAndHash string
	modTime     string
}

func (fi *FileItem) String() string {
	return fmt.Sprintf("%v %v %v", fi.fullPath, fi.sizeAndHash, fi.modTime)
}

func getWorkDir() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	}

	dir, _ := os.Getwd()
	return dir
}

func calcQuickHash100k(fullPath string) (string, error) {
	f, err := os.Open(fullPath)

	if err != nil {
		return "", err
	}
	defer f.Close()

	//
	buf := make([]byte, 100*1024)
	n, err := f.Read(buf)

	if err != nil {
		return "", err
	}

	buf = buf[0:n]

	//
	value := md5.Sum(buf)
	return hex.EncodeToString(value[:]), nil
}

func newFileItem(fullPath string, size int64, modTime string) *FileItem {
	sizeAndHash, err := calcQuickHash100k(fullPath)

	if err != nil {
		return nil
	}

	return &FileItem{
		fullPath:    fullPath,
		sizeAndHash: sizeAndHash,
		modTime:     modTime,
	}
}

func enumFiles(dir string) []*FileItem {
	var files []*FileItem
	var enum func(dir string)

	enum = func(dir string) {
		if dirs, err := os.ReadDir(dir); err == nil {
			for _, d := range dirs {
				if d.IsDir() {
					enum(path.Join(dir, d.Name()))
				} else if info, err := d.Info(); err == nil {
					fi := newFileItem(path.Join(dir, info.Name()), info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
					if fi != nil {
						files = append(files, fi)
					}
				}
			}
		}
	}

	// 枚举
	enum(dir)

	// deletes
	var toDels []string

	// map
	m := make(map[string]*FileItem)

	for _, fi := range files {
		fiOld := m[fi.sizeAndHash]

		if fiOld == nil {
			m[fi.sizeAndHash] = fi
			continue
		}

		if fi.modTime < fiOld.modTime {
			toDels = append(toDels, fiOld.fullPath)
			m[fi.sizeAndHash] = fi
		} else {
			toDels = append(toDels, fi.fullPath)
		}
	}

	for _, f := range toDels {
		fmt.Println(f)
		os.Remove(f)
	}

	return files
}

func main() {
	dir := getWorkDir()

	enumFiles(dir)
}
