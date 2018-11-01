package main

import (
	"github.com/deckarep/gosx-notifier"
	"github.com/mholt/archiver"
	"github.com/radovskyb/watcher"
	log2 "log"
	"os"
	"strconv"
	"time"
)

var verbose = false
var noNotifications = false
var projectFile = ""
var sourcesFolder = ""

func main() {
	// Get program arguments.
	args := os.Args[1:]

	// Todo: Add better flag parsing.
	// Todo: Add .framerx file creation.

	// Check if no arguments exist or if the only argument is '-help'.
	if (len(args) < 2 || (len(args) == 1 && (args[0] == "-help") || args[0] == "❓")) {
		// Print the help message.
		printHelpMessage()
		return
	}

	// Enable Verbose mode
	if (stringInArray("-v", args) || stringInArray("📢", args)) {
		verbose = true
		log("Verose Mode enabled.")
	}

	// Disable notifications
	if (stringInArray("-nn", args) || stringInArray("🔕", args)) {
		noNotifications = false
		log("Notifications disabled")
	}

	// Set project file
	log("Setting projectFile variable to " + args[0])
	projectFile = args[0]

	// Set sources folder
	log("Setting projectFile variable to " + args[1])
	sourcesFolder = args[1]

	notify(1)

	// Create a new file watcher
	w := watcher.New()

	// Make the file watcher react to write events
	w.FilterOps(watcher.Write)

	// Define what happens when an event fires.
	go func() {
		for {
			select {
				case event := <- w.Event: {
					str := strconv.FormatInt(event.Size(), 10)
					log("Update detected! New file size: " + str)
					notify(0)
					extractSources()
				}
			}
		}
	}()

	if err := w.Add(args[0]); err != nil {
		error(err.Error())
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		error(err.Error())
	}
}

// Checks if a string is in an Array.
func stringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Pushes a notification with a specific type.
// Types:
//   0: Update notification
//   1: Start notification
func notify(notifType int) {
	log("Preparing to push notification..")

	switch notifType {
		case 0:
			if (!noNotifications) {
				note := gosxnotifier.NewNotification("Sources have been updated 👌")
				note.Sender = "design.lucasritter.framerx.watcher"
				note.Group = "design.lucasritter.framerx.watcher"
				note.Title = "Framer X Watcher 👀"

				err := note.Push()

				if err != nil {
					error("Failed to push notification..")
				}
			}
			break;

		case 1:
			if (!noNotifications) {
				note := gosxnotifier.NewNotification("Started watching " + projectFile + " 😎")
				note.Sender = "design.lucasritter.framerx.watcher"
				note.Group = "design.lucasritter.framerx.watcher"
				note.Title = "Framer X Watcher 👀"

				err := note.Push()

				if err != nil {
					error("Failed to push notification..")
				}
			}
			break;
	}
}

// Logs a message.
func log(message string) {
	if (verbose) {
		println(message)
	}
}

// Also logs a message.
func error(message string) {
	log2.Fatal(message)
}

// Extracts the sources from the Framer X project file.
func extractSources() {
	// Checks if the path exists.
	if (!pathExists(sourcesFolder)) {
		log("Path doesn't exist yet. Creating..")
		os.Mkdir(sourcesFolder, 0700)
	}
	archiver.Zip.Open(projectFile, sourcesFolder)
}

// Returns a value indicating whether a path exists or not.
func pathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}