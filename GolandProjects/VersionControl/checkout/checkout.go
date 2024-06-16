package checkout

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CopyFilecheckout(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %v", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("could not copy content: %v", err)
	}

	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("could not flush file to storage: %v", err)
	}

	return nil
}

func HandleCheckoutCommand() {
	if len(os.Args) != 3 {
		fmt.Println("Commit id was not passed.")
		return
	}

	commitID := os.Args[2]
	commitIDCorrect := false

	entries, err := os.ReadDir("./vcs/commits")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == commitID {
			commitIDCorrect = true
			break
		}
	}

	if !commitIDCorrect {
		fmt.Println("Commit does not exist.")
		return
	}

	commitPath := filepath.Join("./vcs/commits", commitID)
	commitEntries, err := os.ReadDir(commitPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range commitEntries {
		src := filepath.Join(commitPath, entry.Name())
		dst := filepath.Join(".", entry.Name())
		if err := CopyFilecheckout(src, dst); err != nil {
			fmt.Println("Error copying file for commit:", err)
			return
		}
	}

	fmt.Printf("Switched to commit %s.", commitID)
}
