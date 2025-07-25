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
 * Created Date: 2022-05-25
 * Author: chengyi (Cyber-SiKu)
 */

package usage

import (
	basecmd "github.com/dingodb/dingofs-tools/pkg/cli/command"
	fsUsage "github.com/dingodb/dingofs-tools/pkg/cli/command/v1/usage/fs"
	inode "github.com/dingodb/dingofs-tools/pkg/cli/command/v1/usage/inode"
	"github.com/dingodb/dingofs-tools/pkg/cli/command/v1/usage/metadata"
	"github.com/spf13/cobra"
)

type UsageCommand struct {
	basecmd.MidDingoCmd
}

var _ basecmd.MidDingoCmdFunc = (*UsageCommand)(nil) // check interface

func (usageCmd *UsageCommand) AddSubCommands() {
	usageCmd.Cmd.AddCommand(
		inode.NewInodeNumCommand(),
		metadata.NewMetadataCommand(),
		fsUsage.NewFsUsageCommand(),
	)
}

func NewUsageCommand() *cobra.Command {
	usageCmd := &UsageCommand{
		basecmd.MidDingoCmd{
			Use:   "usage",
			Short: "get the usage info of dingofs",
		},
	}
	return basecmd.NewMidDingoCli(&usageCmd.MidDingoCmd, usageCmd)
}
