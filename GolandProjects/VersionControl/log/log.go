package log

import (
	"fmt"
	"os"
	"strings"
)

func HandleLogCommand() {
	logFilePath := "./vcs/log.txt"
	if !fileExists(logFilePath) {
		fmt.Println("No commits yet.")
		return
	}
	data, err := os.ReadFile(logFilePath)
	if err != nil {
		fmt.Println("Cannot read log file:", err)
		return
	}

	// Split the log entries and print them in reverse order
	logEntries := strings.Split(strings.TrimSpace(string(data)), "\n\n")
	for i := len(logEntries) - 1; i >= 0; i-- {
		fmt.Println(logEntries[i])
		if i > 0 {
			fmt.Println()
		}
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
