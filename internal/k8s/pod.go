package k8s

import (
	"os"
	"context"
	"bufio"
	"fmt"
	"strconv"
	"io"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
    "github.com/maxtoolbox/maxlog/internal/cmdln"
)

// GetClientSet creates and returns a Kubernetes clientset.
//
// Returns:
//   *kubernetes.Clientset - A clientset for interacting with the Kubernetes API.
//   error - An error if the clientset creation fails.
//
// Behavior:
//   - Loads the default kubeconfig file using clientcmd.
//   - Applies configuration overrides if necessary.
//   - Creates a clientset using the loaded kubeconfig.
//   - Terminates the program with a fatal error if the kubeconfig cannot be loaded.
func GetClientSet() (*kubernetes.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		cmdln.Fatal("Error loading kubeconfig:", err)
	}
	return kubernetes.NewForConfig(config)
}

// GetNSPods retrieves the PodInterface for the specified Kubernetes namespace.
//
// Returns:
//   typedv1.PodInterface - An interface for interacting with pods in the namespace.
//
// Behavior:
//   - Creates a Kubernetes clientset using GetClientSet.
//   - Retrieves the namespace from the environment variable MAXLOG_K8S_NAMESPACE.
//   - Returns the PodInterface for the specified namespace.
//   - Terminates the program with a fatal error if the clientset creation fails.
func GetNSPods() typedv1.PodInterface {
	clientset, err := GetClientSet()
	if err != nil {
		cmdln.Fatal("Error creating clientset:", err)
	}
	project := os.Getenv("MAXLOG_K8S_NAMESPACE")
	return clientset.CoreV1().Pods(project)
}

// GetPods retrieves a list of pods matching the specified label selector.
//
// Returns:
//   *corev1.PodList - A list of pods matching the label selector.
//   error - An error if the pod retrieval fails.
//
// Behavior:
//   - Constructs a label selector using the environment variable MAXLOG_K8S_APPTYPE.
//   - Uses the PodInterface from GetNSPods to list pods in the namespace.
//   - Returns the list of pods and any error encountered during retrieval.
func GetPods() (*corev1.PodList, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: "mas.ibm.com/appTypeName=" + os.Getenv("MAXLOG_K8S_APPTYPE"),
	}
    return GetNSPods().List(context.TODO(), listOptions)
}

// GetLog retrieves logs for Kubernetes pods and processes them.
//
// Parameter:
//   tail - A string representing the number of lines to tail from the logs.
//   follow - A boolean indicating whether to follow the log stream.
//   tag - A string representing a tag to be applied to the log lines.
//
// Behavior:
//   - Fetches the list of pods using the GetPods function.
//   - Handles errors that occur during pod retrieval by terminating the program.
//   - Passes the retrieved pods and tail parameter to the getPodLogs function for log processing.
func GetLog(tail string, follow bool, tag string) {
    pods, err := GetPods()
    if err != nil {
        cmdln.Fatal("Error getting pods:", err)
    }
    getPodLogs(pods, tail, follow, tag)
}

// getPodLogs retrieves and processes logs for a list of Kubernetes pods.
//
// Parameters:
//   pods - A pointer to a corev1.PodList containing the pods to retrieve logs from.
//   tail - A string representing the number of lines to tail from the logs.
//   follow - A boolean indicating whether to follow the log stream.
//   tag - A string representing a tag to be applied to the log lines.
//
// Behavior:
//   - Parses the `tail` parameter into an integer value.
//   - Configures pod log options, including tailing the specified number of lines and following the logs.
//   - Iterates through the list of pods and retrieves their logs using the Kubernetes client.
//   - Starts a goroutine for each pod to process its logs using the `writeLogs` function.
//   - Waits for all log processing goroutines to complete before returning.
func getPodLogs(pods *corev1.PodList, tail string, follow bool, tag string) {
    tailnum, err := strconv.ParseInt(tail, 10, 64)
    if err != nil {
        cmdln.Fatal("Error parsing tail number:", err)
    }

    podLogOpts := corev1.PodLogOptions{
        Follow:    follow,
        TailLines: &tailnum,
        Container: os.Getenv("MAXLOG_K8S_APPTYPE"),
    }

    ctx := context.TODO()
    ch := make(chan bool)

    for _, pod := range pods.Items {
        podLogs, err := GetNSPods().GetLogs(pod.Name, &podLogOpts).Stream(ctx)
        if err != nil {
            cmdln.Fatal("Error getting pod logs:", err)
        }
        go writeLogs(bufio.NewReader(podLogs), ch, tag)
    }

    <-ch
}

// writeLogs reads log lines from a buffered reader and processes them.
//
// Parameters:
//   buffer - A pointer to a bufio.Reader that provides the log lines to read.
//   ch - A channel used to signal when the log processing is complete.
//   tag - A string representing a tag to be applied to the log lines.
//
// Behavior:
//   - Continuously reads lines from the buffer until EOF is reached.
//   - Each line is processed using the SetLabels function to apply formatting.
//   - Outputs the formatted log lines to the standard output.
//   - Signals completion by sending a value to the provided channel.
func writeLogs(buffer *bufio.Reader, ch chan bool, tag string) {
    defer func() { ch <- true }()

    for {
        line, err := buffer.ReadString('\n')
        if err == io.EOF {
            break
        }
        fmt.Print(cmdln.SetLabels(line, tag))
    }
}