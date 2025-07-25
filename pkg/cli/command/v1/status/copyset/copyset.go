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
 * Created Date: 2022-06-25
 * Author: chengyi (Cyber-SiKu)
 */

package copyset

import (
	"fmt"
	"strings"

	cmderror "github.com/dingodb/dingofs-tools/internal/error"
	cobrautil "github.com/dingodb/dingofs-tools/internal/utils"
	basecmd "github.com/dingodb/dingofs-tools/pkg/cli/command"
	checkCopyset "github.com/dingodb/dingofs-tools/pkg/cli/command/v1/check/copyset"
	listCopyset "github.com/dingodb/dingofs-tools/pkg/cli/command/v1/list/copyset"
	"github.com/dingodb/dingofs-tools/pkg/config"
	"github.com/dingodb/dingofs-tools/pkg/output"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type CopysetCommand struct {
	basecmd.FinalDingoCmd
	health cobrautil.ClUSTER_HEALTH_STATUS
}

var _ basecmd.FinalDingoCmdFunc = (*CopysetCommand)(nil) // check interface

const (
	copysetExample = `$ dingo status copyset`
)

func NewCopysetCommand() *cobra.Command {
	return NewStatusCopysetCommand().Cmd
}

func (cCmd *CopysetCommand) AddFlags() {
	config.AddRpcRetryTimesFlag(cCmd.Cmd)
	config.AddRpcRetryDelayFlag(cCmd.Cmd)
	config.AddRpcTimeoutFlag(cCmd.Cmd)
	config.AddFsMdsAddrFlag(cCmd.Cmd)
	config.AddMarginOptionFlag(cCmd.Cmd)
}

func (cCmd *CopysetCommand) Init(cmd *cobra.Command, args []string) error {
	cCmd.health = cobrautil.HEALTH_ERROR
	response, err := listCopyset.GetCopysetsInfos(cCmd.Cmd)
	if err.TypeCode() != cmderror.CODE_SUCCESS {
		cCmd.Error = err
		return fmt.Errorf(err.Message)
	}
	var copysetIdVec []string
	var poolIdVec []string
	for _, v := range response.GetCopysetValues() {
		info := v.GetCopysetInfo()
		copysetIdVec = append(copysetIdVec, fmt.Sprintf("%d", info.GetCopysetId()))
		poolIdVec = append(poolIdVec, fmt.Sprintf("%d", info.GetPoolId()))
	}
	if len(copysetIdVec) == 0 {
		var err error
		cCmd.Error = cmderror.ErrSuccess()
		cCmd.Result = "No copyset found"
		cCmd.health = cobrautil.HEALTH_OK
		return err
	}
	copysetIds := strings.Join(copysetIdVec, ",")
	poolIds := strings.Join(poolIdVec, ",")
	result, table, errCheck, health := checkCopyset.GetCopysetsStatus(cCmd.Cmd, copysetIds, poolIds)
	cCmd.Result = result
	cCmd.TableNew = table
	cCmd.Error = errCheck
	cCmd.health = health
	return nil
}

func (cCmd *CopysetCommand) Print(cmd *cobra.Command, args []string) error {
	return output.FinalCmdOutput(&cCmd.FinalDingoCmd, cCmd)
}

func (cCmd *CopysetCommand) RunCommand(cmd *cobra.Command, args []string) error {
	return nil
}

func (cCmd *CopysetCommand) ResultPlainOutput() error {
	return output.FinalCmdOutputPlain(&cCmd.FinalDingoCmd)
}

func NewStatusCopysetCommand() *CopysetCommand {
	copysetCmd := &CopysetCommand{
		FinalDingoCmd: basecmd.FinalDingoCmd{
			Use:     "copyset",
			Short:   "status all copyset of the dingofs",
			Example: copysetExample,
		},
	}
	basecmd.NewFinalDingoCli(&copysetCmd.FinalDingoCmd, copysetCmd)
	return copysetCmd
}

func GetCopysetStatus(caller *cobra.Command) (*interface{}, *tablewriter.Table, *cmderror.CmdError, cobrautil.ClUSTER_HEALTH_STATUS) {
	copysetCmd := NewStatusCopysetCommand()
	copysetCmd.Cmd.SetArgs([]string{
		fmt.Sprintf("--%s", config.FORMAT), config.FORMAT_NOOUT,
	})
	config.AlignFlagsValue(caller, copysetCmd.Cmd, []string{
		config.RPCRETRYTIMES, config.RPCTIMEOUT, config.DINGOFS_MDSADDR,
	})
	copysetCmd.Cmd.SilenceErrors = true
	copysetCmd.Cmd.Execute()
	return &copysetCmd.Result, copysetCmd.TableNew, copysetCmd.Error, copysetCmd.health
}
