/*
 *  Copyright (c) 2022 NetEase Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

/*
 * Project: DingoCli
 * Created Date: 2022-06-20
 * Author: chengyi (Cyber-SiKu)
 */

package create

import (
	basecmd "github.com/dingodb/dingofs-tools/pkg/cli/command"
	"github.com/dingodb/dingofs-tools/pkg/cli/command/v1/create/fs"
	"github.com/dingodb/dingofs-tools/pkg/cli/command/v1/create/subpath"

	"github.com/dingodb/dingofs-tools/pkg/cli/command/v1/create/topology"
	"github.com/spf13/cobra"
)

type CreateCommand struct {
	basecmd.MidDingoCmd
}

var _ basecmd.MidDingoCmdFunc = (*CreateCommand)(nil) // check interface

func (createCmd *CreateCommand) AddSubCommands() {
	createCmd.Cmd.AddCommand(
		fs.NewFsCommand(),
		topology.NewTopologyCommand(),
		subpath.NewSubPathCommand(),
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
