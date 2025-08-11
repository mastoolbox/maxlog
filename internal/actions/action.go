package actions

import (
	"fmt"
	"strings"
)

// ActionFunc defines a function type that operates on an Action instance.
type ActionFunc func(*Action)

// Action represents an action with various attributes and a function to execute.
type Action struct {
	name      string      // The name of the action.
	tag       string      // The tag associated with the action.
	apptype   string      // The application type for the action.
	namespace string      // The namespace in which the action operates.
	tail      string      // The tail parameter for the action.
	follow    bool        // The follow parameter for the action.
	runAction ActionFunc  // The function to execute the action.
}

// ActionRunner defines an interface for initializing and running actions.
type ActionRunner interface {
	// Init initializes the action with the provided arguments.
	// Parameters:
	//   args - A slice of strings representing the arguments for initialization.
	// Returns:
	//   error - An error if initialization fails.
	Init([]string) error

	// Run executes the action.
	Run()

	// SetName sets the name of the action.
	// Parameters:
	//   name - A string representing the name to be assigned to the action.
	SetName(string)

	// GetName retrieves the name of the action.
	// Returns:
	//   string - The name of the action.
	GetName() string
}

// SetName sets the name of the Action.
// Parameters:
//   name - A string representing the name to be assigned to the Action.
func (act *Action) SetName(name string) {
	act.name = name
}

// GetName retrieves the name of the Action.
// Returns:
//   string - The name of the Action.
func (act *Action) GetName() string {
	return act.name
}

// trimSubCmd trims the leading dashes from a subcommand.
// Parameters:
//   subcmd - A string representing the subcommand to be trimmed.
// Returns:
//   string - The trimmed subcommand.
func trimSubCmd(subcmd string) string {
    if subcmd[0:2] == "--" {
		return subcmd[2:]
	} else if subcmd[0] == '-' {
		return subcmd[1:]
	}
	return subcmd
}

// splitSubCmd splits arguments into subcommands and their values.
// Parameters:
//   args - A slice of strings representing the arguments.
// Returns:
//   []string - A slice of strings containing subcommands and their values.
func splitSubCmd(args []string) []string {
	subCmds := []string{}
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			vals := strings.Split(arg, "=")
			subCmds = append(subCmds, vals[0])
			subCmds = append(subCmds, vals[1])
		} else {
			subCmds = append(subCmds, arg)
		}
	}
	return subCmds
}

// Init initializes the Action with the provided arguments.
// Parameters:
//   args - A slice of strings representing the arguments for initialization.
// Returns:
//   error - An error if the initialization fails due to missing parameters.
func (act *Action) Init(args []string) error {
    act.follow = true
    act.tag    = ""
	args = splitSubCmd(args)
	if len(args)%2 > 0 {
		return fmt.Errorf("Missing second parameter")
	}
	for i := 0; i < len(args); i += 2 {
		switch trimSubCmd(args[i]) {
		case "tag":
			act.tag = args[i+1]
		case "namespace":
			act.namespace = args[i+1]
		case "apptype":
			act.apptype = args[i+1]
		case "tail":
			act.tail = args[i+1]
		case "follow":
            act.follow = !(args[i+1] == "0" || args[i+1] == "no" || args[i+1] == "false")
		}
	}
	return nil
}

// Run executes the Action by invoking its associated function.
// Behavior:
//   - Calls the runAction function with the Action instance as a parameter.
func (act *Action) Run() {
	act.runAction(act)
}