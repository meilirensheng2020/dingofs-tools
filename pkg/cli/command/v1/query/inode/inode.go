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
 * Created Date: 2022-06-28
 * Author: chengyi (Cyber-SiKu)
 */

package inode

import (
	"context"
	"fmt"

	cmderror "github.com/dingodb/dingofs-tools/internal/error"
	cobrautil "github.com/dingodb/dingofs-tools/internal/utils"
	"github.com/dingodb/dingofs-tools/pkg/base"
	basecmd "github.com/dingodb/dingofs-tools/pkg/cli/command"
	"github.com/dingodb/dingofs-tools/pkg/cli/command/v1/list/partition"
	"github.com/dingodb/dingofs-tools/pkg/cli/command/v1/query/copyset"
	"github.com/dingodb/dingofs-tools/pkg/config"
	"github.com/dingodb/dingofs-tools/pkg/output"
	"github.com/dingodb/dingofs-tools/proto/dingofs/proto/common"
	"github.com/dingodb/dingofs-tools/proto/dingofs/proto/metaserver"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
)

const (
	inodeExample = `$ dingo query inode --fsid 1 --inodeid 1`
)

type QueryInodeRpc struct {
	Info             *base.Rpc
	Request          *metaserver.GetInodeRequest
	metaserverClient metaserver.MetaServerServiceClient
}

var _ base.RpcFunc = (*QueryInodeRpc)(nil) // check interface

func (qiRpc *QueryInodeRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	qiRpc.metaserverClient = metaserver.NewMetaServerServiceClient(cc)
}

func (qiRpc *QueryInodeRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := qiRpc.metaserverClient.GetInode(ctx, qiRpc.Request)
	output.ShowRpcData(qiRpc.Request, response, qiRpc.Info.RpcDataShow)
	return response, err
}

type InodeCommand struct {
	basecmd.FinalDingoCmd
	QIRpc *QueryInodeRpc
}

var _ basecmd.FinalDingoCmdFunc = (*InodeCommand)(nil) // check interface

func NewInodeCommand() *cobra.Command {
	inodeCmd := &InodeCommand{
		FinalDingoCmd: basecmd.FinalDingoCmd{
			Use:     "inode",
			Short:   "query the inode of fs",
			Example: inodeExample,
		},
	}
	basecmd.NewFinalDingoCli(&inodeCmd.FinalDingoCmd, inodeCmd)
	return inodeCmd.Cmd
}

func (iCmd *InodeCommand) AddFlags() {
	config.AddRpcRetryTimesFlag(iCmd.Cmd)
	config.AddRpcRetryDelayFlag(iCmd.Cmd)
	config.AddRpcTimeoutFlag(iCmd.Cmd)
	config.AddFsMdsAddrFlag(iCmd.Cmd)
	config.AddFsIdRequiredFlag(iCmd.Cmd)
	config.AddInodeIdRequiredFlag(iCmd.Cmd)
}

func (iCmd *InodeCommand) Init(cmd *cobra.Command, args []string) error {
	header := []string{
		cobrautil.ROW_FS_ID, cobrautil.ROW_INODE_ID, cobrautil.ROW_LENGTH, cobrautil.ROW_TYPE, cobrautil.ROW_NLINK, cobrautil.ROW_PARENT,
	}
	iCmd.Header = header

	return nil
}

func (iCmd *InodeCommand) Print(cmd *cobra.Command, args []string) error {
	return output.FinalCmdOutput(&iCmd.FinalDingoCmd, iCmd)
}

func (iCmd *InodeCommand) Prepare() error {
	fsId := config.GetFlagUint32(iCmd.Cmd, config.DINGOFS_FSID)
	inodeId := config.GetFlagUint64(iCmd.Cmd, config.DINGOFS_INODEID)

	fsId2PartitionList, errGet := partition.GetFsPartition(iCmd.Cmd)
	if errGet.TypeCode() != cmderror.CODE_SUCCESS {
		iCmd.Error = errGet
		return errGet.ToError()
	}
	partitionInfoList := (*fsId2PartitionList)[fsId]
	if partitionInfoList == nil {
		return fmt.Errorf("inode[%d] is not found in fs[%d]", inodeId, fsId)
	}
	index := slices.IndexFunc(partitionInfoList,
		func(p *common.PartitionInfo) bool {
			return p.GetFsId() == fsId && p.GetStart() <= inodeId && p.GetEnd() >= inodeId
		})
	if index < 0 {
		return fmt.Errorf("inode[%d] is not on any partition of fs[%d]", inodeId, fsId)
	}
	partitionInfo := partitionInfoList[index]
	poolId := partitionInfo.GetPoolId()
	copyetId := partitionInfo.GetCopysetId()
	partitionId := partitionInfo.GetPartitionId()
	supportStream := false
	inodeRequest := &metaserver.GetInodeRequest{
		PoolId:           &poolId,
		CopysetId:        &copyetId,
		PartitionId:      &partitionId,
		FsId:             &fsId,
		InodeId:          &inodeId,
		SupportStreaming: &supportStream,
	}
	iCmd.QIRpc = &QueryInodeRpc{
		Request: inodeRequest,
	}
	// get addrs
	config.AddCopysetidSliceRequiredFlag(iCmd.Cmd)
	config.AddPoolidSliceRequiredFlag(iCmd.Cmd)
	iCmd.Cmd.ParseFlags([]string{
		fmt.Sprintf("--%s", config.DINGOFS_COPYSETID), fmt.Sprintf("%d", copyetId),
		fmt.Sprintf("--%s", config.DINGOFS_POOLID), fmt.Sprintf("%d", poolId),
	})
	key2Copyset, errQuery := copyset.QueryCopysetInfo(iCmd.Cmd)
	if errQuery.TypeCode() != cmderror.CODE_SUCCESS {
		return fmt.Errorf("query copyset info failed: %s", errQuery.Message)
	}
	if len(*key2Copyset) == 0 {
		return fmt.Errorf("no copysetinfo found")
	}
	key := cobrautil.GetCopysetKey(uint64(poolId), uint64(copyetId))
	leader := (*key2Copyset)[key].Info.GetLeaderPeer()
	addr, peerErr := cobrautil.PeertoAddr(leader)
	if peerErr.TypeCode() != cmderror.CODE_SUCCESS {
		return fmt.Errorf("pares leader peer[%s] failed: %s", leader, peerErr.Message)
	}
	addrs := []string{addr}

	timeout := config.GetRpcTimeout(iCmd.Cmd)
	retrytimes := config.GetRpcRetryTimes(iCmd.Cmd)
	retryDelay := config.GetRpcRetryDelay(iCmd.Cmd)
	verbose := config.GetFlagBool(iCmd.Cmd, config.VERBOSE)
	iCmd.QIRpc.Info = base.NewRpc(addrs, timeout, retrytimes, retryDelay, verbose, "GetInode")

	return nil
}

func (iCmd *InodeCommand) RunCommand(cmd *cobra.Command, args []string) error {
	preErr := iCmd.Prepare()
	if preErr != nil {
		return preErr
	}

	inodeResult, err := base.GetRpcResponse(iCmd.QIRpc.Info, iCmd.QIRpc)
	if err.TypeCode() != cmderror.CODE_SUCCESS {
		return fmt.Errorf("get inode failed: %s", err.Message)
	}
	getInodeResponse := inodeResult.(*metaserver.GetInodeResponse)
	if getInodeResponse.GetStatusCode() != metaserver.MetaStatusCode_OK {
		return fmt.Errorf("get inode failed: %s", getInodeResponse.GetStatusCode().String())
	}
	inode := getInodeResponse.GetInode()
	tableRows := make([]map[string]string, 0)
	if len(inode.S3ChunkInfoMap) == 0 {
		row := make(map[string]string)
		row[cobrautil.ROW_FS_ID] = fmt.Sprintf("%d", inode.GetFsId())
		row[cobrautil.ROW_INODE_ID] = fmt.Sprintf("%d", inode.GetInodeId())
		row[cobrautil.ROW_LENGTH] = fmt.Sprintf("%d", inode.GetLength())
		row[cobrautil.ROW_TYPE] = inode.GetType().String()
		row[cobrautil.ROW_NLINK] = fmt.Sprintf("%d", inode.GetNlink())
		row[cobrautil.ROW_PARENT] = fmt.Sprintf("%d", inode.GetParent())
		tableRows = append(tableRows, row)
	} else {
		rows := make([]map[string]string, 0)
		infoMap := inode.GetS3ChunkInfoMap()
		iCmd.Header = append(iCmd.Header, cobrautil.ROW_S3CHUNKINFO_CHUNKID)
		iCmd.Header = append(iCmd.Header, cobrautil.ROW_S3CHUNKINFO_OFFSET)
		iCmd.Header = append(iCmd.Header, cobrautil.ROW_S3CHUNKINFO_LENGTH)
		iCmd.Header = append(iCmd.Header, cobrautil.ROW_S3CHUNKINFO_SIZE)
		for _, infoList := range infoMap {
			for _, info := range infoList.GetS3Chunks() {
				row := make(map[string]string)
				row[cobrautil.ROW_FS_ID] = fmt.Sprintf("%d", inode.GetFsId())
				row[cobrautil.ROW_INODE_ID] = fmt.Sprintf("%d", inode.GetInodeId())
				row[cobrautil.ROW_LENGTH] = fmt.Sprintf("%d", inode.GetLength())
				row[cobrautil.ROW_TYPE] = inode.GetType().String()
				row[cobrautil.ROW_NLINK] = fmt.Sprintf("%d", inode.GetNlink())
				row[cobrautil.ROW_PARENT] = fmt.Sprintf("%d", inode.GetParent())
				row[cobrautil.ROW_S3CHUNKINFO_CHUNKID] = fmt.Sprintf("%d", info.GetChunkId())
				row[cobrautil.ROW_S3CHUNKINFO_OFFSET] = fmt.Sprintf("%d", info.GetOffset())
				row[cobrautil.ROW_S3CHUNKINFO_LENGTH] = fmt.Sprintf("%d", info.GetLen())
				row[cobrautil.ROW_S3CHUNKINFO_SIZE] = fmt.Sprintf("%d", info.GetSize())
				rows = append(rows, row)
			}
		}
		tableRows = append(tableRows, rows...)
	}

	iCmd.SetHeader(iCmd.Header)
	list := cobrautil.ListMap2ListSortByKeys(tableRows, iCmd.Header, []string{})
	iCmd.TableNew.AppendBulk(list)
	iCmd.Result = getInodeResponse
	iCmd.Error = cmderror.ErrSuccess()
	return nil
}

func (iCmd *InodeCommand) ResultPlainOutput() error {
	return output.FinalCmdOutputPlain(&iCmd.FinalDingoCmd)
}
