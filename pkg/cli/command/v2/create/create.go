// Copyright (c) 2025 dingodb.com, Inc. All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package create

import (
	basecmd "github.com/dingodb/dingofs-tools/pkg/cli/command"
	"github.com/spf13/cobra"
)

type CreateCommand struct {
	basecmd.MidDingoCmd
}

var _ basecmd.MidDingoCmdFunc = (*CreateCommand)(nil) // check interface

func (createCmd *CreateCommand) AddSubCommands() {
	createCmd.Cmd.AddCommand(
		NewCreateFsCommand(),
		NewSubPathCommand(),
	)
}

func NewCreateCommand() *cobra.Command {
	createCmd := &CreateCommand{
		basecmd.MidDingoCmd{
			Use:   "create",
			Short: "create resources in the dingofs",
		},
	}
	return basecmd.NewMidDingoCli(&createCmd.MidDingoCmd, createCmd)
}
