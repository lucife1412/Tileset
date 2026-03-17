package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

// OpenUrl opens the specified URL in the default browser of the user.
func OpenUrl(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	port := "8080"
	addr := ":" + port
	url := "http://localhost" + addr

	// Serve static files from the current directory
	fs := http.FileServer(http.Dir(dir))
	
	// Add a handler to add CORS headers and cache control
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow loading resources (like tileset.json) from other origins if needed
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Set Cache-Control to avoid caching issues during development
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		
		fs.ServeHTTP(w, r)
	}))

	fmt.Println("Starting server at " + url)
	fmt.Println("Serving files from: " + dir)
	fmt.Println("Press Ctrl+C to stop the server.")

	// Try to open the browser automatically
	go func() {
		// Wait a bit to ensure server is up (simple sleep)
		// In a robust app we might check if port is listening
		OpenUrl(url)
	}()

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
