package actions

import (
    "os"
    "github.com/maxtoolbox/maxlog/internal/moby"
    "github.com/maxtoolbox/maxlog/internal/cmdln"
    "github.com/maxtoolbox/maxlog/internal/k8s"
)

// ActionLogs creates and initializes an Action for retrieving logs.
//
// Returns:
//   *Action - A pointer to the initialized Action instance.
//
// Behavior:
//   - Sets the name of the Action to "logs".
//   - Assigns the runLogs function to the Action's runAction field.
func ActionLogs() *Action {
    act := &Action{
        name: "logs",
    }
    act.runAction = runLogs
    return act
}

// runLogs retrieves logs based on the MAXLOG_MODE environment variable.
//
// Parameters:
//   act - A pointer to the Action instance.
//
// Behavior:
//   - Reads the MAXLOG_MODE environment variable to determine the mode of operation.
//   - If the mode is "k8s", retrieves logs for Kubernetes resources using the k8s.GetLog function.
//   - If the mode is "pod", retrieves logs for a specific container using the moby.GetLog function.
//   - Logs a fatal error if MAXLOG_MODE is not set or contains an invalid value.
func runLogs(act *Action) {
    tail := getEnv("MAXLOG_TAIL", "40")
    if act.tail != "" {
        tail = act.tail
    }

    if os.Getenv("MAXLOG_MODE") == "k8s" {
        selector := os.Getenv("MAXLOG_K8S_APPTYPE")
        namespace := os.Getenv("MAXLOG_K8S_NAMESPACE")
        /* tag := ""
        if len(os.Args) > 3 {
            if os.Args[2] == "tag" {
                tag = os.Args[3]
            }
        } */
        if namespace == "" || selector == "" {
            cmdln.Fatal("Please set MAXLOG_K8S_NAMESPACE and MAXLOG_K8S_APPTYPE environment variables.", nil)
        }
        k8s.GetLog(tail, act.follow, act.tag)
    } else if os.Getenv("MAXLOG_MODE") == "pod" {
        container := os.Getenv("MAXLOG_CONTAINER")
        if container == "" {
            cmdln.Fatal("Container name is not set. Please set MAXLOG_CONTAINER environment variable.", nil)
        }
        moby.GetLog(container, tail, act.follow, act.tag)
    } else {
        cmdln.Fatal("Unknown MAXLOG_MODE. Please set it to 'k8s' or 'pod'.", nil)
    }
}

// getEnv retrieves the value of an environment variable or returns a fallback value.
//
// Parameters:
//   key - A string representing the name of the environment variable.
//   fallback - A string representing the fallback value.
//
// Returns:
//   string - The value of the environment variable if set, otherwise the fallback value.
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}