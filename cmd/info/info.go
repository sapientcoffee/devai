/*
Copyright 2024 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package info

import (
	"fmt"

	"github.com/spf13/cobra"
)

// infoCmd represents the echo command
var InfoCmd = &cobra.Command{
	Aliases: []string{"echo"},
	Use:     "info",
	Short:   "Info command and test",
	Long: `A useful command to view details and information about the CLI.
	
	It is alos a useful subcommand to test the CLI is working`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current implementation of info command is limited.")
	},
}

func init() {
}
