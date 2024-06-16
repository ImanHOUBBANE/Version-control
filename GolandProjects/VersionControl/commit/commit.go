package commit

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func HandleCommitCommand() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 3 {
		fmt.Println("Message was not passed.")
		return
	}

	err := os.MkdirAll("./vcs/commits", os.ModePerm)
	if err != nil {
		fmt.Println("Cannot create directory \"./vcs/commits\"")
		return
	}

	sliceFile := GetLineIndex("./vcs/index.txt")
	currentHashes := make(map[string]string)
	for _, filePath := range sliceFile {
		currentHashes[filePath] = HashFile(filePath)
	}

	lastCommitHashes := make(map[string]string)
	if fileInfo, err := os.Stat("./vcs/log.txt"); err == nil && fileInfo.Size() > 0 {
		lastCommit := GetLastCommit("./vcs/log.txt")
		lastCommitFiles := GetFileLastCommitdir(filepath.Join("./vcs/commits", lastCommit[7:71]))
		for _, file := range lastCommitFiles {
			lastCommitHashes[file] = HashFile(filepath.Join("./vcs/commits", lastCommit[7:71], file))
		}
	}

	modified := false
	for file, currentHash := range currentHashes {
		if lastHash, exists := lastCommitHashes[file]; !exists || lastHash != currentHash {
			modified = true
			break
		}
	}

	if modified {
		commitName := HashString(strconv.Itoa(rand.Int()))
		pathCommit := filepath.Join("./vcs/commits", commitName)
		err := os.MkdirAll(pathCommit, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating commit directory:", err)
			return
		}

		for file, _ := range currentHashes {
			src := file
			dst := filepath.Join(pathCommit, filepath.Base(file))
			if err := CopyFile(src, dst); err != nil {
				fmt.Println("Error copying file for commit:", err)
				return
			}
		}

		AddInfoLog(commitName)
		fmt.Println("Changes are committed.")
	} else {
		fmt.Println("Nothing to commit.")
	}
}
func GetLineIndex(filePath string) []string {
	var lines []string
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return lines
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	firstLine := true
	for scanner.Scan() {
		if firstLine {
			firstLine = false
			continue
		}
		lines = append(lines, scanner.Text())
	}
	return lines
}

func HashFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		fmt.Println("Error hashing file:", err)
		return ""
	}
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func GetLastCommit(filePath string) string {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	commits := strings.Split(strings.TrimSpace(string(data)), "\n\n")
	if len(commits) > 0 {
		return commits[len(commits)-1]
	}
	return ""
}

func GetFileLastCommitdir(commitPath string) []string {
	var files []string
	entries, err := ioutil.ReadDir(commitPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return files
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files
}

func AddInfoLog(commitName string) {
	logPath := "./vcs/log.txt"
	author, _ := os.ReadFile("./vcs/config.txt")
	logEntry := fmt.Sprintf("commit %s\nAuthor: %s\n%s", commitName, strings.TrimSpace(string(author)), os.Args[2])
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(logEntry + "\n\n"); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

func HashString(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("could not open source file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("copying failed: %v", err)
	}

	return nil
}
