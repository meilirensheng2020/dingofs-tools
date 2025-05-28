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

package common

import (
	"context"
	"fmt"
	cmderror "github.com/dingodb/dingofs-tools/internal/error"
	basecmd "github.com/dingodb/dingofs-tools/pkg/cli/command"
	cmdcommon "github.com/dingodb/dingofs-tools/pkg/cli/command/common"
	"github.com/dingodb/dingofs-tools/pkg/config"
	pbmdsv2error "github.com/dingodb/dingofs-tools/proto/dingofs/proto/error"
	pbmdsv2 "github.com/dingodb/dingofs-tools/proto/dingofs/proto/mdsv2"
	"github.com/spf13/cobra"
	"log"
	"sync"
	"sync/atomic"
)

//public functions

// create new mds rpc
func CreateNewMdsRpc(cmd *cobra.Command, serviceName string) (*basecmd.Rpc, error) {
	// get mds address
	addrs, getAddrErr := config.GetFsMdsAddrSlice(cmd)
	if getAddrErr.TypeCode() != cmderror.CODE_SUCCESS {
		return nil, fmt.Errorf(getAddrErr.Message)
	}
	// new rpc
	timeout := config.GetRpcTimeout(cmd)
	retryTimes := config.GetRpcRetryTimes(cmd)
	mdsRpc := basecmd.NewRpc(addrs, timeout, retryTimes, serviceName)
	mdsRpc.RpcDataShow = config.GetFlagBool(cmd, "verbose")

	return mdsRpc, nil
}

// get fsinfo by fsid or fsname
func GetFsInfo(cmd *cobra.Command, fsId uint32, fsName string) (*pbmdsv2.FsInfo, error) {
	// new prc
	mdsRpc, err := CreateNewMdsRpc(cmd, "GetFsInfo")
	if err != nil {
		return nil, err
	}
	// set request info
	var getFsRpc *GetFsRpc
	if fsId > 0 {
		getFsRpc = &GetFsRpc{Info: mdsRpc, Request: &pbmdsv2.GetFsInfoRequest{FsId: fsId}}
	} else {
		getFsRpc = &GetFsRpc{Info: mdsRpc, Request: &pbmdsv2.GetFsInfoRequest{FsName: fsName}}
	}
	// get rpc result
	response, errCmd := basecmd.GetRpcResponse(getFsRpc.Info, getFsRpc)
	if errCmd.TypeCode() != cmderror.CODE_SUCCESS {
		return nil, fmt.Errorf(errCmd.Message)
	}
	result := response.(*pbmdsv2.GetFsInfoResponse)
	if mdsErr := result.GetError(); mdsErr.GetErrcode() != pbmdsv2error.Errno_OK {
		return nil, cmderror.MDSV2Error(mdsErr).ToError()
	}
	return result.GetFsInfo(), nil
}

// retrieve fsid from command-line parameters,if not set, get by GetFsInfo via fsname
func GetFsId(cmd *cobra.Command) (uint32, error) {
	fsId, fsName, fsErr := cmdcommon.CheckAndGetFsIdOrFsNameValue(cmd)
	if fsErr != nil {
		return 0, fsErr
	}
	// fsId is not set,need to get fsId by fsName (fsName -> fsId)
	if fsId == 0 {
		fsInfo, fsErr := GetFsInfo(cmd, 0, fsName)
		if fsErr != nil {
			return 0, fsErr
		}
		fsId = fsInfo.GetFsId()
		if fsId == 0 {
			return 0, fmt.Errorf("fsid is invalid")
		}
	}
	return fsId, nil
}

// retrieve fsid from command-line parameters,if not set, get by GetFsInfo via fsid
func GetFsName(cmd *cobra.Command) (string, error) {
	fsId, fsName, fsErr := cmdcommon.CheckAndGetFsIdOrFsNameValue(cmd)
	if fsErr != nil {
		return "", fsErr
	}
	if len(fsName) == 0 { // fsName is not set,need to get fsName by fsId (fsId->fsName)
		fsInfo, fsErr := GetFsInfo(cmd, fsId, "")
		if fsErr != nil {
			return "", fsErr
		}
		fsName = fsInfo.GetFsName()
		if len(fsName) == 0 {
			return "", fmt.Errorf("fsName is invalid")
		}
	}
	return fsName, nil
}

// get inode
func GetInode(cmd *cobra.Command, fsId uint32, inodeId uint64) (*pbmdsv2.Inode, error) {
	// new prc
	mdsRpc, err := CreateNewMdsRpc(cmd, "GetInode")
	if err != nil {
		return nil, err
	}
	// set request info
	getInodeRpc := &GetInodeRpc{Info: mdsRpc, Request: &pbmdsv2.GetInodeRequest{FsId: fsId, Ino: inodeId}}
	// get rpc result
	response, errCmd := basecmd.GetRpcResponse(getInodeRpc.Info, getInodeRpc)
	if errCmd.TypeCode() != cmderror.CODE_SUCCESS {
		return nil, fmt.Errorf(errCmd.Message)
	}
	result := response.(*pbmdsv2.GetInodeResponse)
	if mdsErr := result.GetError(); mdsErr.GetErrcode() != pbmdsv2error.Errno_OK {
		return nil, cmderror.MDSV2Error(mdsErr).ToError()
	}
	return result.GetInode(), nil
}

// list dentry
func ListDentry(cmd *cobra.Command, fsId uint32, inodeId uint64) ([]*pbmdsv2.Dentry, error) {
	// new prc
	mdsRpc, err := CreateNewMdsRpc(cmd, "ListDentry")
	if err != nil {
		return nil, err
	}
	// set request info
	listDentryRpc := &ListDentryRpc{Info: mdsRpc, Request: &pbmdsv2.ListDentryRequest{FsId: fsId, Parent: inodeId, Limit: 100}}
	// get rpc result
	response, errCmd := basecmd.GetRpcResponse(listDentryRpc.Info, listDentryRpc)
	if errCmd.TypeCode() != cmderror.CODE_SUCCESS {
		return nil, fmt.Errorf(errCmd.Message)
	}
	result := response.(*pbmdsv2.ListDentryResponse)
	if mdsErr := result.GetError(); mdsErr.GetErrcode() != pbmdsv2error.Errno_OK {
		return nil, cmderror.MDSV2Error(mdsErr).ToError()
	}
	return result.GetDentries(), nil
}

// get directory size and inodes by inode
func GetDirSummarySize(cmd *cobra.Command, fsId uint32, inode uint64, summary *cmdcommon.Summary, concurrent chan struct{},
	ctx context.Context, cancel context.CancelFunc, isFsCheck bool, inodeMap *sync.Map) error {
	var err error
	entries, entErr := ListDentry(cmd, fsId, inode)
	if entErr != nil {
		return entErr
	}
	var wg sync.WaitGroup
	var errCh = make(chan error, 1)
	for _, entry := range entries {
		if entry.GetType() == pbmdsv2.FileType_FILE {
			inodeAttr, err := GetInode(cmd, fsId, entry.GetIno())
			if err != nil {
				return err
			}
			if isFsCheck && inodeAttr.GetNlink() >= 2 { //filesystem check, hardlink is ignored
				if _, ok := inodeMap.LoadOrStore(inodeAttr.GetIno(), struct{}{}); ok {
					continue
				}
			}
			atomic.AddUint64(&summary.Length, inodeAttr.GetLength())
		}
		atomic.AddUint64(&summary.Inodes, 1)
		if entry.GetType() != pbmdsv2.FileType_DIRECTORY {
			continue
		}
		select {
		case err := <-errCh:
			cancel()
			return err
		case <-ctx.Done():
			return fmt.Errorf("cancel scan directory for other goroutine error")
		case concurrent <- struct{}{}:
			wg.Add(1)
			go func(e *pbmdsv2.Dentry) {
				defer wg.Done()
				sumErr := GetDirSummarySize(cmd, fsId, e.GetIno(), summary, concurrent, ctx, cancel, isFsCheck, inodeMap)
				<-concurrent
				if sumErr != nil {
					select {
					case errCh <- sumErr:
					default:
					}
				}
			}(entry)
		default:
			if sumErr := GetDirSummarySize(cmd, fsId, entry.GetIno(), summary, concurrent, ctx, cancel, isFsCheck, inodeMap); sumErr != nil {
				return sumErr
			}
		}
	}
	wg.Wait()
	select {
	case err = <-errCh:
	default:
	}
	return err
}

// get directory size and inodes by path name
func GetDirectorySizeAndInodes(cmd *cobra.Command, fsId uint32, dirInode uint64, isFsCheck bool) (int64, int64, error) {
	log.Printf("start to summary directory statistics, inode[%d]", dirInode)
	summary := &cmdcommon.Summary{0, 0}
	concurrent := make(chan struct{}, 50)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var inodeMap *sync.Map = &sync.Map{}
	sumErr := GetDirSummarySize(cmd, fsId, dirInode, summary, concurrent, ctx, cancel, isFsCheck, inodeMap)
	log.Printf("end summary directory statistics, inode[%d],inodes[%d],size[%d]", dirInode, summary.Inodes, summary.Length)
	if sumErr != nil {
		return 0, 0, sumErr
	}
	return int64(summary.Length), int64(summary.Inodes), nil
}
