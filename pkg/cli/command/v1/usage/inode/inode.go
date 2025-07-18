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

package inode

import (
	"fmt"
	"strconv"
	"strings"

	cmderror "github.com/dingodb/dingofs-tools/internal/error"
	cobrautil "github.com/dingodb/dingofs-tools/internal/utils"
	"github.com/dingodb/dingofs-tools/pkg/base"
	basecmd "github.com/dingodb/dingofs-tools/pkg/cli/command"

	config "github.com/dingodb/dingofs-tools/pkg/config"
	"github.com/dingodb/dingofs-tools/pkg/output"
	"github.com/spf13/cobra"
)

const (
	PREFIX = "topology_fs_id"
	SUFFIX = "inode_num"
)

type InodeNumCommand struct {
	basecmd.FinalDingoCmd
	FsId2Filetype2Metric map[string]map[string]*base.Metric
}

type Result struct {
	Result string
	Error  *cmderror.CmdError
	SubUri string
}

var _ basecmd.FinalDingoCmdFunc = (*InodeNumCommand)(nil) // check interface

const (
	inodeExample = `$ dingofs usage inode
$ dingofs usage inode --fsid 1,2,3`
)

func NewInodeNumCommand() *cobra.Command {
	inodeNumCmd := &InodeNumCommand{
		FinalDingoCmd: basecmd.FinalDingoCmd{
			Use:     "inode",
			Short:   "get the inode usage of dingofs",
			Example: inodeExample,
		},
	}
	basecmd.NewFinalDingoCli(&inodeNumCmd.FinalDingoCmd, inodeNumCmd)
	return inodeNumCmd.Cmd
}

func (iCmd *InodeNumCommand) AddFlags() {
	config.AddFsMdsAddrFlag(iCmd.Cmd)
	config.AddFsIdOptionDefaultAllFlag(iCmd.Cmd)
	config.AddHttpTimeoutFlag(iCmd.Cmd)
}

func (iCmd *InodeNumCommand) Init(cmd *cobra.Command, args []string) error {
	addrs, addrErr := config.GetFsMdsAddrSlice(iCmd.Cmd)
	if addrErr.TypeCode() != cmderror.CODE_SUCCESS {
		iCmd.Error = addrErr
		return fmt.Errorf(addrErr.Message)
	}

	iCmd.FsId2Filetype2Metric = make(map[string]map[string]*base.Metric)

	fsIds := config.GetFlagStringSliceDefaultAll(iCmd.Cmd, config.DINGOFS_FSID)
	if len(fsIds) == 0 {
		fsIds = []string{"*"}
	}
	for _, fsId := range fsIds {
		_, err := strconv.ParseUint(fsId, 10, 32)
		if err != nil && fsId != "*" {
			return fmt.Errorf("invalid fsId: %s", fsId)
		}
		if fsId == "*" {
			fsId = ""
		}
		subUri := fmt.Sprintf("/vars/"+PREFIX+"_%s*"+SUFFIX, fsId)
		timeout := config.GetHttpTimeout(cmd)
		metric := base.NewMetric(addrs, subUri, timeout)
		filetype2Metric := make(map[string]*base.Metric)
		filetype2Metric["inode_num"] = metric
		iCmd.FsId2Filetype2Metric[fsId] = filetype2Metric
	}
	header := []string{cobrautil.ROW_FS_ID, cobrautil.ROW_STORAGE_TYPE, cobrautil.ROW_NUM}
	iCmd.SetHeader(header)
	iCmd.TableNew.SetAutoMergeCellsByColumnIndex(cobrautil.GetIndexSlice(header, []string{cobrautil.ROW_FS_ID}))

	return nil
}

func (iCmd *InodeNumCommand) Print(cmd *cobra.Command, args []string) error {
	return output.FinalCmdOutput(&iCmd.FinalDingoCmd, iCmd)
}

func (iCmd *InodeNumCommand) RunCommand(cmd *cobra.Command, args []string) error {
	results := make(chan Result, config.MaxChannelSize())
	size := 0
	for fsId, filetype2Metric := range iCmd.FsId2Filetype2Metric {
		for filetype, metric := range filetype2Metric {
			size++
			go func(m *base.Metric, filetype string, id string) {
				result, err := base.QueryMetric(m)
				results <- Result{
					Result: result,
					Error:  err,
					SubUri: m.SubUri,
				}
			}(metric, filetype, fsId)
		}
	}
	count := 0
	rows := make([]map[string]string, 0)
	var errs []*cmderror.CmdError
	for res := range results {
		datas := strings.Split(res.Result, "\n")
		if res.Error.Code == cmderror.CODE_SUCCESS {
			for _, data := range datas {
				if data == "" {
					continue
				}
				data = cobrautil.RmWitespaceStr(data)
				resMap := strings.Split(data, ":")
				preMap := strings.Split(resMap[0], "_")
				if len(resMap) != 2 && len(preMap) < 4 {
					splitErr := cmderror.ErrDataNoExpected()
					splitErr.Format(data, "the length of the data does not meet the requirements")
					errs = append(errs, splitErr)
				} else {
					num, errNum := strconv.ParseInt(resMap[1], 10, 64)
					id, errId := strconv.ParseUint(preMap[3], 10, 32)
					prefix := fmt.Sprintf("%s_%s_", PREFIX, preMap[3])
					filetype := strings.Replace(resMap[0], prefix, "", 1)
					filetype = strings.Replace(filetype, SUFFIX, "", 1)
					if filetype == "" {
						filetype = "inode_num_"
					}
					filetype = filetype[0 : len(filetype)-1]
					if errNum == nil && errId == nil {
						row := make(map[string]string)
						row[cobrautil.ROW_FS_ID] = strconv.FormatUint(uint64(id), 10)
						row[cobrautil.ROW_STORAGE_TYPE] = filetype
						row[cobrautil.ROW_NUM] = strconv.FormatInt(num, 10)
						rows = append(rows, row)
					} else {
						toErr := cmderror.ErrDataNoExpected()
						toErr.Format(data)
						errs = append(errs, toErr)
					}
				}
			}
		} else {
			errs = append(errs, res.Error)
		}
		count++
		if count >= size {
			break
		}
	}
	iCmd.Error = cmderror.MostImportantCmdError(errs)

	mergeErr := cmderror.MergeCmdErrorExceptSuccess(errs)
	iCmd.Error = mergeErr

	if len(rows) > 0 {
		list := cobrautil.ListMap2ListSortByKeys(rows, iCmd.Header, []string{
			config.DINGOFS_FSID,
		})
		iCmd.TableNew.AppendBulk(list)
		iCmd.Result = rows
	}
	return mergeErr.ToError()
}

func (iCmd *InodeNumCommand) ResultPlainOutput() error {
	return output.FinalCmdOutputPlain(&iCmd.FinalDingoCmd)
}
