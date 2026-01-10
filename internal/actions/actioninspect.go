package actions

import (
	"fmt"
	"os"

	"github.com/maxtoolbox/maxlog/internal/cmdln"
	"github.com/maxtoolbox/maxlog/internal/k8s"
	"github.com/maxtoolbox/maxlog/internal/moby"
)

// ActionInspect creates and initializes an Action for inspecting logs or resources.
//
// Returns:
//
//	*Action - A pointer to the initialized Action instance.
//
// Behavior:
//   - Sets the name of the Action to "inspect".
//   - Assigns the runInspect function to the Action's runAction field.
func ActionInspect() *Action {
	act := &Action{
		name: "inspect",
	}
	act.runAction = runInspect
	return act
}

// runInspect executes the inspection logic based on the MAXLOG_MODE environment variable.
//
// Parameters:
//
//	act - A pointer to the Action instance.
//
// Behavior:
//   - Reads the MAXLOG_MODE environment variable to determine the mode of operation.
//   - Calls inspectK8s if the mode is "k8s".
//   - Calls inspectPod if the mode is "pod".
//   - Logs a fatal error if MAXLOG_MODE is not set or contains an invalid value.
func runInspect(act *Action) {
	mode := os.Getenv("MAXLOG_MODE")
	if mode == "" {
		cmdln.Fatal(" MAXLOG_MODE is not set. Please set it to 'k8s' for Kubernetes mode or 'pod' for podman.", nil)
	} else if mode == "k8s" {
		inspectK8s()
	} else if mode == "pod" {
		inspectPod()
	} else {
		cmdln.Fatal("Unknown mode: '"+mode+"'. Please set MAXLOG_MODE to 'k8s' or 'pod'.", nil)
	}
}

// inspectK8s retrieves and displays information about Kubernetes resources.
//
// Behavior:
//   - Reads the MAXLOG_K8S_NAMESPACE and MAXLOG_K8S_APPTYPE environment variables.
//   - Logs a fatal error if the required environment variables are not set.
//   - Retrieves the list of pods using the k8s.GetPods function.
//   - Displays the namespace, application type, tail parameter, and the number of selected pods.
func inspectK8s() {
	ns := os.Getenv("MAXLOG_K8S_NAMESPACE")
	apptype := cmdln.GetEnv("MAXLOG_K8S_APPTYPE", cmdln.DefaultLabels)
	fmt.Println("Namespace    :", ns)
	fmt.Println("AppType      :", apptype)
	fmt.Println("Tail         :", os.Getenv("MAXLOG_TAIL"))
	if len(ns) == 0 {
		cmdln.Fatal("Please set MAXLOG_K8S_NAMESPACE environment variables.", nil)
	}
	pods, err := k8s.GetPods()
	if err != nil {
		cmdln.Fatal(" Error getting pods: ", err)
	}
	fmt.Println("Selected Pods:", len(pods.Items))
}

// inspectPod retrieves and displays information about a specific container.
//
// Behavior:
//   - Reads the MAXLOG_CONTAINER and MAXLOG_TAIL environment variables.
//   - Retrieves the container ID using the moby.GetCID function.
//   - Displays the container name, container ID, and tail parameter.
func inspectPod() {
	fmt.Println("Container    :", os.Getenv("MAXLOG_CONTAINER"))
	fmt.Println("CID          :", moby.GetCID(os.Getenv("MAXLOG_CONTAINER")))
	fmt.Println("Tail         :", os.Getenv("MAXLOG_TAIL"))
}
