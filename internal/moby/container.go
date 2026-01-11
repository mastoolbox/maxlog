package moby

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/maxtoolbox/maxlog/internal/cmdln"

	// Never mind. We use the Moby client for Podman. We can swap it out later.
	// Unfortunately, I'm having some problems with Windows right now.
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// GetLog retrieves and processes logs for a specific container.
//
// Parameters:
//
//	name - A string representing the name of the container whose logs are to be retrieved.
//	tail - A string specifying the number of lines to tail from the logs.
//	follow - A boolean indicating whether to follow the log stream.
//	tag - A string representing a tag to be applied to the log lines.
//
// Behavior:
//   - Resolves the container ID using the GetCID function based on the provided container name.
//   - Passes the container ID and tail parameter to the getContainerLogs function for log retrieval and processing.
func GetLog(name, tail string, follow bool, tag string) {
	getContainerLogs(GetCID(name), tail, follow, tag)
}

// GetCID retrieves the container ID for a given container name.
//
// Parameters:
//
//	name - A string representing the name of the container to search for.
//
// Returns:
//
//	string - The ID of the container matching the given name.
//
// Behavior:
//   - Creates a Moby client to interact with the container runtime.
//   - Retrieves a list of all containers using the Moby client.
//   - Iterates through the list of containers to find one whose name matches the provided name.
//   - Returns the container ID if a match is found.
//   - Logs a fatal error and terminates the program if no matching container is found.
func GetCID(name string) string {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if slices.Contains(container.Names, "/"+name) {
			return container.ID
		}
	}

	log.Fatal("The search for a container has not yielded any results.")
	return ""
}

// getContainerLogs retrieves and processes logs for a specific container.
//
// Parameters:
//
//	cid - A string representing the container ID whose logs are to be retrieved.
//	tail - A string specifying the number of lines to tail from the logs.
//	follow - A boolean indicating whether to follow the log stream.
//	tag - A string representing a tag to be applied to the log lines.
//
// Behavior:
//   - Creates a Moby client to interact with the container runtime.
//   - Configures log options, including stdout, stderr, timestamps, and tailing.
//   - Retrieves the container logs using the specified options.
//   - Processes the log stream by reading headers and data chunks.
//   - Extracts and formats log lines using the SetLabels function.
//   - Handles errors during log retrieval and processing, terminating the program if necessary.
func getContainerLogs(cid, tail string, follow bool, tag string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      "",
		Until:      "",
		Timestamps: true,
		Follow:     follow,
		Tail:       tail,
		Details:    false,
	}

	reader, err := cli.ContainerLogs(context.Background(), cid, options)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer reader.Close()

	hdr := make([]byte, 8)
	for {
		_, err := reader.Read(hdr)
		if err != nil {
			if err == io.EOF {
				return
			}

			log.Fatal(err)
		}

		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = reader.Read(dat)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		// time, line, found
		_, line, found := strings.Cut(string(dat), " ")
		if found {
			text := cmdln.SetLabels(line, tag)
			if text != "" {
				fmt.Print(text)
			}
		}
	}
}
