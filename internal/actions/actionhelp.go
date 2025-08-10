package actions

import (
	"fmt"
)

// ActionHelp creates and initializes an Action for displaying help information.
//
// Returns:
//   *Action - A pointer to the initialized Action instance.
//
// Behavior:
//   - Sets the name of the Action to "help".
//   - Assigns the runHelp function to the Action's runAction field.
func ActionHelp() *Action {
	act := &Action{
		name: "help",
	}
	act.runAction = runHelp
	return act
}

// runHelp displays usage instructions and available actions.
//
// Parameters:
//   act - A pointer to the Action instance.
//
// Behavior:
//   - Prints usage instructions for the `maxlog` command.
//   - Lists available actions and their descriptions.
//   - Provides an example usage of the `maxlog` command.
//   - Mentions relevant environment variables for configuration.
func runHelp(act *Action) {
    fmt.Println("Usage: maxlog [action] [options]")
    fmt.Println("Available actions:")
    fmt.Println("  logs       - Show logs of containers")
    fmt.Println("  inspect    - Inspect pods or containers")
    fmt.Println("  version    - Show version information")
    fmt.Println("  help       - Show this help message")
    fmt.Println("Example: maxlog logs --tag=mytag --tail=100")
    fmt.Println("If no action is set, logs will be used.")
    fmt.Println("Environment variables:")
    fmt.Println("  MAXLOG_MODE        - Set to 'k8s' for Kubernetes mode or 'pod' for podman mode")
    fmt.Println("  Podman mode")
    fmt.Println("  MAXLOG_CONTAINER   - Specify the container name in podman mode")
    fmt.Println("  K8s mode")
}