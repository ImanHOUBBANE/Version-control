package add

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func HandleAddCommand() {
	err := os.MkdirAll("../Version Control/", os.ModePerm)
	if err != nil {
		fmt.Println("Cannot create directory \"./vcs\"")
		return
	}
	filePath := "../Version Control/index.txt"
	if len(os.Args) == 3 {
		if os.Args[2] == "." {
			addAllFiles(filePath)
		} else {
			addSingleFile(filePath, os.Args[2])
		}
	} else if len(os.Args) == 2 {
		listTrackedFiles(filePath)
	}
}

func addAllFiles(indexPath string) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			addSingleFile(indexPath, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error adding files:", err)
	}
}

func addSingleFile(indexPath, filePath string) {
	file, err := os.OpenFile(indexPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return
	}
	defer file.Close()

	if fileExists(filePath) {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Println("Cannot resolve absolute path:", err)
			return
		}

		trackedFiles := getTrackedFiles(indexPath)
		if _, exists := trackedFiles[absPath]; exists {
			fmt.Printf("The file '%s' is already tracked.\n", filePath)
			return
		}

		if len(trackedFiles) == 0 {
			_, err := fmt.Fprintf(file, "Tracked files:\n")
			if err != nil {
				fmt.Println("Cannot write to file:", err)
				return
			}
		}

		_, err = fmt.Fprintln(file, absPath)
		if err != nil {
			fmt.Println("Cannot write to file:", err)
			return
		}
		fmt.Printf("The file '%s' is tracked.\n", filePath)
	} else {
		fmt.Printf("Can't find '%s'.\n", filePath)
	}
}

func listTrackedFiles(indexPath string) {
	file, err := os.OpenFile(indexPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return
	}
	defer file.Close()

	if info, _ := file.Stat(); info.Size() == 0 {
		err := os.WriteFile(indexPath, []byte("Tracked files:\n"), 0666)
		if err != nil {
			fmt.Println("Cannot write to file:", err)
			return
		}
		fmt.Println("Add a file to the index.")
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}
	}
}

func getTrackedFiles(indexPath string) map[string]struct{} {
	file, err := os.Open(indexPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	trackedFiles := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "Tracked files:") {
			trackedFiles[line] = struct{}{}
		}
	}
	return trackedFiles
}
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
