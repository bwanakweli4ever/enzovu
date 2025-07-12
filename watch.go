//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	fmt.Println("🔄 Simple Go file watcher for Enzovu")

	var cmd *exec.Cmd

	// Handle Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
		os.Exit(0)
	}()

	lastMod := time.Now()

	for {
		// Check for Go file modifications
		modified := false
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if filepath.Ext(path) == ".go" && info.ModTime().After(lastMod) {
				modified = true
				lastMod = time.Now()
			}
			return nil
		})

		if modified {
			// Kill existing process if running!
			if cmd != nil && cmd.Process != nil {
				cmd.Process.Kill()
				cmd.Wait()
			}

			fmt.Println("📝 Changes detected, restarting...")

			// Start new process
			cmd = exec.Command("go", "run", "main.go")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		}

		time.Sleep(500 * time.Millisecond)
	}
}
