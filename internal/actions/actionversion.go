package actions

import (
	"fmt"
)

// ActionVersion creates and initializes an Action for displaying version information.
//
// Returns:
//
//	*Action - A pointer to the initialized Action instance.
//
// Behavior:
//   - Sets the name of the Action to "version".
//   - Assigns the runVersion function to the Action's runAction field.
func ActionVersion() *Action {
	act := &Action{
		name: "version",
	}
	act.runAction = runVersion
	return act
}

// runVersion displays the version information of the application.
//
// Parameters:
//
//	act - A pointer to the Action instance.
//
// Behavior:
//   - Prints the current version of the application to the console.
func runVersion(act *Action) {
	fmt.Println("maxlog version: 0.0.4")
}
