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
 * Created Date: 2022-06-16
 * Author: chengyi (Cyber-SiKu)
 */

package metaserver

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	cmderror "github.com/dingodb/dingofs-tools/internal/error"
	cobrautil "github.com/dingodb/dingofs-tools/internal/utils"
	"github.com/dingodb/dingofs-tools/pkg/base"
	basecmd "github.com/dingodb/dingofs-tools/pkg/cli/command"
	"github.com/dingodb/dingofs-tools/pkg/config"
	"github.com/dingodb/dingofs-tools/pkg/output"
	"github.com/dingodb/dingofs-tools/proto/dingofs/proto/topology"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	metaserverExample = `$ dingo query metaserver --metaserverid=1,2,3
$ dingo query metaserver --metaserveraddr=127.0.0.1:9700,127.0.0.1:9701,127.0.0.1:9702`
)

type QueryMetaserverRpc struct {
	Info           *base.Rpc
	Request        *topology.GetMetaServerInfoRequest
	topologyClient topology.TopologyServiceClient
}

var _ base.RpcFunc = (*QueryMetaserverRpc)(nil) // check interface

type MetaserverCommand struct {
	basecmd.FinalDingoCmd
	Rpc  []*QueryMetaserverRpc
	Rows []map[string]string
}

var _ basecmd.FinalDingoCmdFunc = (*MetaserverCommand)(nil) // check interface

func (qmRpc *QueryMetaserverRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	qmRpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (qmRpc *QueryMetaserverRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := qmRpc.topologyClient.GetMetaServer(ctx, qmRpc.Request)
	output.ShowRpcData(qmRpc.Request, response, qmRpc.Info.RpcDataShow)
	return response, err
}

func NewMetaserverCommand() *cobra.Command {
	metaserverCmd := &MetaserverCommand{
		FinalDingoCmd: basecmd.FinalDingoCmd{
			Use:     "metaserver",
			Short:   "query metaserver in dingofs by metaserverid or metaserveraddr",
			Long:    "when both metaserverid and metaserveraddr exist, query only by metaserverid",
			Example: metaserverExample,
		},
	}
	basecmd.NewFinalDingoCli(&metaserverCmd.FinalDingoCmd, metaserverCmd)
	return metaserverCmd.Cmd
}

func (mCmd *MetaserverCommand) AddFlags() {
	config.AddRpcRetryTimesFlag(mCmd.Cmd)
	config.AddRpcRetryDelayFlag(mCmd.Cmd)
	config.AddRpcTimeoutFlag(mCmd.Cmd)
	config.AddFsMdsAddrFlag(mCmd.Cmd)
	config.AddMetaserverAddrOptionFlag(mCmd.Cmd)
	config.AddMetaserverIdOptionFlag(mCmd.Cmd)
}

func (mCmd *MetaserverCommand) Init(cmd *cobra.Command, args []string) error {
	addrs, addrErr := config.GetFsMdsAddrSlice(mCmd.Cmd)
	if addrErr.TypeCode() != cmderror.CODE_SUCCESS {
		mCmd.Error = addrErr
		return fmt.Errorf(addrErr.Message)
	}

	var metaserverAddrs []string
	var metaserverIds []string
	if viper.IsSet(config.VIPER_DINGOFS_METASERVERADDR) && !viper.IsSet(config.VIPER_DINGOFS_METASERVERID) {
		// metaserveraddr is set, but metaserverid is not set
		metaserverAddrs = viper.GetStringSlice(config.VIPER_DINGOFS_METASERVERADDR)
	} else {
		metaserverIds = viper.GetStringSlice(config.VIPER_DINGOFS_METASERVERID)
	}

	if len(metaserverAddrs) == 0 && len(metaserverIds) == 0 {
		return fmt.Errorf("%s or %s is required", config.DINGOFS_METASERVERADDR, config.DINGOFS_METASERVERID)
	}

	header := []string{cobrautil.ROW_ID, cobrautil.ROW_HOSTNAME, cobrautil.ROW_INTERNAL_ADDR, cobrautil.ROW_EXTERNAL_ADDR, cobrautil.ROW_ONLINE_STATE}
	mCmd.SetHeader(header)

	mCmd.Rows = make([]map[string]string, 0)
	timeout := config.GetRpcTimeout(cmd)
	retrytimes := config.GetRpcRetryTimes(cmd)
	retryDelay := config.GetRpcRetryDelay(cmd)
	verbose := config.GetFlagBool(cmd, config.VERBOSE)
	for i := range metaserverAddrs {
		addr := strings.Split(metaserverAddrs[i], ":")
		if len(addr) != 2 {
			return fmt.Errorf("unrecognized metaserver addr: %s", metaserverAddrs[i])
		}
		port, err := strconv.ParseUint(addr[1], 10, 32)
		if err != nil {
			return fmt.Errorf("unrecognized metaserver port: %s", metaserverAddrs[i])
		}
		port32 := uint32(port)
		request := &topology.GetMetaServerInfoRequest{
			HostIp: &addr[0],
			Port:   &port32,
		}
		rpc := &QueryMetaserverRpc{
			Request: request,
		}
		rpc.Info = base.NewRpc(addrs, timeout, retrytimes, retryDelay, verbose, "GetMetaServerInfo")
		mCmd.Rpc = append(mCmd.Rpc, rpc)
		row := make(map[string]string)
		row[cobrautil.ROW_ID] = cobrautil.ROW_VALUE_DNE
		row[cobrautil.ROW_HOSTNAME] = cobrautil.ROW_VALUE_DNE
		row[cobrautil.ROW_INTERNAL_ADDR] = cobrautil.ROW_VALUE_DNE
		row[cobrautil.ROW_EXTERNAL_ADDR] = metaserverAddrs[i]
		row[cobrautil.ROW_ONLINE_STATE] = cobrautil.ROW_VALUE_DNE
		mCmd.Rows = append(mCmd.Rows, row)
	}

	for i := range metaserverIds {
		id, err := strconv.ParseUint(metaserverIds[i], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid %s: %s", config.DINGOFS_METASERVERID, metaserverIds[i])
		}
		id32 := uint32(id)
		request := &topology.GetMetaServerInfoRequest{
			MetaServerID: &id32,
		}
		rpc := &QueryMetaserverRpc{
			Request: request,
		}
		rpc.Info = base.NewRpc(addrs, timeout, retrytimes, retryDelay, verbose, "GetMetaServerInfo")
		rpc.Info.RpcDataShow = config.GetFlagBool(mCmd.Cmd, config.VERBOSE)
		mCmd.Rpc = append(mCmd.Rpc, rpc)
		row := make(map[string]string)
		row[cobrautil.ROW_ID] = metaserverIds[i]
		row[cobrautil.ROW_EXTERNAL_ADDR] = ""
		mCmd.Rows = append(mCmd.Rows, row)
	}

	return nil
}

func (mCmd *MetaserverCommand) Print(cmd *cobra.Command, args []string) error {
	return output.FinalCmdOutput(&mCmd.FinalDingoCmd, mCmd)
}

func (mCmd *MetaserverCommand) RunCommand(cmd *cobra.Command, args []string) error {
	var infos []*base.Rpc
	var funcs []base.RpcFunc
	for _, rpc := range mCmd.Rpc {
		infos = append(infos, rpc.Info)
		funcs = append(funcs, rpc)
	}

	results, errs := base.GetRpcListResponse(infos, funcs)
	if len(errs) == len(infos) {
		mergeErr := cmderror.MergeCmdErrorExceptSuccess(errs)
		return mergeErr.ToError()
	}
	var resList []interface{}
	for _, result := range results {
		if result == nil {
			continue
		}
		response := result.(*topology.GetMetaServerInfoResponse)
		res, err := output.MarshalProtoJson(response)
		if err != nil {
			errMar := cmderror.ErrMarShalProtoJson()
			errMar.Format(err.Error())
			errs = append(errs, errMar)
		}
		resList = append(resList, res)
		if response.GetStatusCode() != topology.TopoStatusCode_TOPO_OK {
			code := response.GetStatusCode()
			err := cmderror.ErrGetMetaserverInfo(int(code))
			err.Format(topology.TopoStatusCode_name[int32(response.GetStatusCode())])
			errs = append(errs, err)
			continue
		}
		metaserverInfo := response.GetMetaServerInfo()
		for _, row := range mCmd.Rows {
			id := strconv.FormatUint(uint64(metaserverInfo.GetMetaServerID()), 10)
			externalAddr := fmt.Sprintf("%s:%d", metaserverInfo.GetExternalIp(), metaserverInfo.GetExternalPort())
			if row[cobrautil.ROW_ID] == id || row[cobrautil.ROW_EXTERNAL_ADDR] == externalAddr {
				row[cobrautil.ROW_ID] = id
				row[cobrautil.ROW_HOSTNAME] = metaserverInfo.GetHostname()
				internalAddr := fmt.Sprintf("%s:%d", metaserverInfo.GetInternalIp(), metaserverInfo.GetInternalPort())
				row[cobrautil.ROW_INTERNAL_ADDR] = internalAddr
				row[cobrautil.ROW_EXTERNAL_ADDR] = externalAddr
				row[cobrautil.ROW_ONLINE_STATE] = metaserverInfo.GetOnlineState().String()
			}
		}
	}

	list := cobrautil.ListMap2ListSortByKeys(mCmd.Rows, mCmd.Header, []string{cobrautil.ROW_ID})
	mCmd.TableNew.AppendBulk(list)
	mCmd.Result = resList
	mCmd.Error = cmderror.MostImportantCmdError(errs)

	return nil
}

func (mCmd *MetaserverCommand) ResultPlainOutput() error {
	return output.FinalCmdOutputPlain(&mCmd.FinalDingoCmd)
}
