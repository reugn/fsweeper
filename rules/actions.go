package rules

import (
	"log"
	"os"
	"path"
	"time"
)

// Action to perform on a file
type Action struct {
	Action  string `yaml:"action"`
	Payload string `yaml:"payload"`
}

func (action *Action) echoAction(filePath string, vars *Vars) {
	msg := vars.Process(action.Payload, filePath)
	log.Printf("[%s] %s\n", filePath, msg)
}

func (action *Action) touchAction(filePath string) {
	now := time.Now().Local()
	err := os.Chtimes(filePath, now, now)
	if err != nil {
		log.Fatalf("Failed to touch file: %s", filePath)
	}
}

func (action *Action) moveAction(filePath string, vars *Vars) {
	newPath := vars.Process(action.Payload, filePath)
	err := os.Rename(filePath, newPath)
	if err != nil {
		log.Fatalf("Failed to move file: %s to: %s", filePath, newPath)
	}
}

func (action *Action) renameAction(filePath string, vars *Vars) {
	newFileName := vars.Process(action.Payload, filePath)
	err := os.Rename(filePath, validatePath(path.Dir(filePath))+newFileName)
	if err != nil {
		log.Fatalf("Failed to rename file: %s to: %s", filePath, newFileName)
	}
}

func (action *Action) deleteAction(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Fatalf("Failed to remove file: %s", filePath)
	}
}
