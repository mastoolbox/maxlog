/*
Copyright 2025 maxlog authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/maxtoolbox/maxlog/internal/actions"
)

func validateArgs(args []string) error {
	if len(args) < 2 {
		args = append(args, "logs")
	}

	command := args[1]
	offset := 2
	if strings.HasPrefix(command, "tag=") || strings.HasPrefix(command, "focus=") {
		command = "logs"
		offset = 1
	}

	cmds := []actions.ActionRunner{
		actions.ActionLogs(),
		actions.ActionInspect(),
		actions.ActionVersion(),
		actions.ActionHelp(),
	}

	for _, cmd := range cmds {
		if cmd.GetName() == command {
			cmd.Init(args[offset:])
			cmd.Run()
			return nil
		}
	}

	return fmt.Errorf("Unknown command: %s", command)
}

func main() {
	if err := validateArgs(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
