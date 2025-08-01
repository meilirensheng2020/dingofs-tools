package common

import (
	"context"

	"github.com/dingodb/dingofs-tools/pkg/base"
	"github.com/dingodb/dingofs-tools/pkg/output"
	"github.com/dingodb/dingofs-tools/proto/dingofs/proto/cachegroup"
	"github.com/dingodb/dingofs-tools/proto/dingofs/proto/mds"
	"github.com/dingodb/dingofs-tools/proto/dingofs/proto/metaserver"
	"github.com/dingodb/dingofs-tools/proto/dingofs/proto/topology"

	"google.golang.org/grpc"
)

// rpc services
type GetFsQuotaRpc struct {
	Info             *base.Rpc
	Request          *metaserver.GetFsQuotaRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type SetFsQuotaRpc struct {
	Info             *base.Rpc
	Request          *metaserver.SetFsQuotaRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type CheckQuotaRpc struct {
	Info             *base.Rpc
	Request          *metaserver.SetDirQuotaRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type DeleteQuotaRpc struct {
	Info             *base.Rpc
	Request          *metaserver.DeleteDirQuotaRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type GetQuotaRpc struct {
	Info             *base.Rpc
	Request          *metaserver.GetDirQuotaRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type ListQuotaRpc struct {
	Info             *base.Rpc
	Request          *metaserver.LoadDirQuotasRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type SetQuotaRpc struct {
	Info             *base.Rpc
	Request          *metaserver.SetDirQuotaRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type QueryCopysetRpc struct {
	Info           *base.Rpc
	Request        *topology.GetCopysetsInfoRequest
	topologyClient topology.TopologyServiceClient
}

type ListPartitionRpc struct {
	Info           *base.Rpc
	Request        *topology.ListPartitionRequest
	topologyClient topology.TopologyServiceClient
}

type CreateInodeRpc struct {
	Info             *base.Rpc
	Request          *metaserver.CreateInodeRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type UpdateInodeRpc struct {
	Info             *base.Rpc
	Request          *metaserver.UpdateInodeRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type DeleteInodeRpc struct {
	Info             *base.Rpc
	Request          *metaserver.DeleteInodeRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type GetInodeAttrRpc struct {
	Info             *base.Rpc
	Request          *metaserver.BatchGetInodeAttrRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type GetInodeRpc struct {
	Info             *base.Rpc
	Request          *metaserver.GetInodeRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type CreateDentryRpc struct {
	Info             *base.Rpc
	Request          *metaserver.CreateDentryRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type ListDentryRpc struct {
	Info             *base.Rpc
	Request          *metaserver.ListDentryRequest
	metaServerClient metaserver.MetaServerServiceClient
}
type GetDentryRpc struct {
	Info             *base.Rpc
	Request          *metaserver.GetDentryRequest
	metaServerClient metaserver.MetaServerServiceClient
}

type GetFsStatsRpc struct {
	Info      *base.Rpc
	Request   *mds.GetFsStatsRequest
	mdsClient mds.MdsServiceClient
}

type ListClusterFsRpc struct {
	Info      *base.Rpc
	Request   *mds.ListClusterFsInfoRequest
	mdsClient mds.MdsServiceClient
}

type ListCacheGroupRpc struct {
	Info             *base.Rpc
	Request          *cachegroup.LoadGroupsRequest
	cacheGroupClient cachegroup.CacheGroupMemberServiceClient
}

type ListCacheMemberRpc struct {
	Info             *base.Rpc
	Request          *cachegroup.LoadMembersRequest
	cacheGroupClient cachegroup.CacheGroupMemberServiceClient
}

type ReweightMemberRpc struct {
	Info             *base.Rpc
	Request          *cachegroup.ReweightMemberRequest
	cacheGroupClient cachegroup.CacheGroupMemberServiceClient
}

var _ base.RpcFunc = (*GetFsQuotaRpc)(nil)      // check interface
var _ base.RpcFunc = (*SetFsQuotaRpc)(nil)      // check interface
var _ base.RpcFunc = (*CheckQuotaRpc)(nil)      // check interface
var _ base.RpcFunc = (*DeleteQuotaRpc)(nil)     // check interface
var _ base.RpcFunc = (*GetQuotaRpc)(nil)        // check interface
var _ base.RpcFunc = (*ListQuotaRpc)(nil)       // check interface
var _ base.RpcFunc = (*SetQuotaRpc)(nil)        // check interface
var _ base.RpcFunc = (*QueryCopysetRpc)(nil)    // check interface
var _ base.RpcFunc = (*ListPartitionRpc)(nil)   // check interface
var _ base.RpcFunc = (*GetInodeAttrRpc)(nil)    // check interface
var _ base.RpcFunc = (*CreateInodeRpc)(nil)     // check interface
var _ base.RpcFunc = (*UpdateInodeRpc)(nil)     // check interface
var _ base.RpcFunc = (*DeleteInodeRpc)(nil)     // check interface
var _ base.RpcFunc = (*GetInodeRpc)(nil)        // check interface
var _ base.RpcFunc = (*ListDentryRpc)(nil)      // check interface
var _ base.RpcFunc = (*GetDentryRpc)(nil)       // check interface
var _ base.RpcFunc = (*GetFsStatsRpc)(nil)      // check interface
var _ base.RpcFunc = (*ListClusterFsRpc)(nil)   // check interface
var _ base.RpcFunc = (*ListCacheGroupRpc)(nil)  // check interface
var _ base.RpcFunc = (*ListCacheMemberRpc)(nil) // check interface
var _ base.RpcFunc = (*ReweightMemberRpc)(nil)  // check interface

func (getFsQuotaRpc *GetFsQuotaRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	getFsQuotaRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (getFsQuotaRpc *GetFsQuotaRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := getFsQuotaRpc.metaServerClient.GetFsQuota(ctx, getFsQuotaRpc.Request)
	output.ShowRpcData(getFsQuotaRpc.Request, response, getFsQuotaRpc.Info.RpcDataShow)
	return response, err
}

func (setFsQuotaRpc *SetFsQuotaRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	setFsQuotaRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (setFsQuotaRpc *SetFsQuotaRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := setFsQuotaRpc.metaServerClient.SetFsQuota(ctx, setFsQuotaRpc.Request)
	output.ShowRpcData(setFsQuotaRpc.Request, response, setFsQuotaRpc.Info.RpcDataShow)
	return response, err
}
func (checkQuotaRpc *CheckQuotaRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	checkQuotaRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (checkQuotaRpc *CheckQuotaRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := checkQuotaRpc.metaServerClient.SetDirQuota(ctx, checkQuotaRpc.Request)
	output.ShowRpcData(checkQuotaRpc.Request, response, checkQuotaRpc.Info.RpcDataShow)
	return response, err
}

func (deleteQuotaRpc *DeleteQuotaRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	deleteQuotaRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (deleteQuotaRpc *DeleteQuotaRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := deleteQuotaRpc.metaServerClient.DeleteDirQuota(ctx, deleteQuotaRpc.Request)
	output.ShowRpcData(deleteQuotaRpc.Request, response, deleteQuotaRpc.Info.RpcDataShow)
	return response, err
}

func (getQuotaRpc *GetQuotaRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	getQuotaRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (getQuotaRpc *GetQuotaRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := getQuotaRpc.metaServerClient.GetDirQuota(ctx, getQuotaRpc.Request)
	output.ShowRpcData(getQuotaRpc.Request, response, getQuotaRpc.Info.RpcDataShow)
	return response, err
}

func (listQuotaRpc *ListQuotaRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	listQuotaRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (listQuotaRpc *ListQuotaRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := listQuotaRpc.metaServerClient.LoadDirQuotas(ctx, listQuotaRpc.Request)
	output.ShowRpcData(listQuotaRpc.Request, response, listQuotaRpc.Info.RpcDataShow)
	return response, err
}

func (setQuotaRpc *SetQuotaRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	setQuotaRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (setQuotaRpc *SetQuotaRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := setQuotaRpc.metaServerClient.SetDirQuota(ctx, setQuotaRpc.Request)
	output.ShowRpcData(setQuotaRpc.Request, response, setQuotaRpc.Info.RpcDataShow)
	return response, err
}

func (qcRpc *QueryCopysetRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	qcRpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (qcRpc *QueryCopysetRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := qcRpc.topologyClient.GetCopysetsInfo(ctx, qcRpc.Request)
	output.ShowRpcData(qcRpc.Request, response, qcRpc.Info.RpcDataShow)
	return response, err
}

func (lpRp *ListPartitionRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	lpRp.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (lpRp *ListPartitionRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := lpRp.topologyClient.ListPartition(ctx, lpRp.Request)
	output.ShowRpcData(lpRp.Request, response, lpRp.Info.RpcDataShow)
	return response, err
}

func (getInodeRpc *GetInodeAttrRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	getInodeRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (getInodeRpc *GetInodeAttrRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := getInodeRpc.metaServerClient.BatchGetInodeAttr(ctx, getInodeRpc.Request)
	output.ShowRpcData(getInodeRpc.Request, response, getInodeRpc.Info.RpcDataShow)
	return response, err
}

func (createInodeRpc *CreateInodeRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	createInodeRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (createInodeRpc *CreateInodeRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := createInodeRpc.metaServerClient.CreateInode(ctx, createInodeRpc.Request)
	output.ShowRpcData(createInodeRpc.Request, response, createInodeRpc.Info.RpcDataShow)
	return response, err
}

func (updateInodeRpc *UpdateInodeRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	updateInodeRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (updateInodeRpc *UpdateInodeRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := updateInodeRpc.metaServerClient.UpdateInode(ctx, updateInodeRpc.Request)
	output.ShowRpcData(updateInodeRpc.Request, response, updateInodeRpc.Info.RpcDataShow)
	return response, err
}

func (deleteInodeRpc *DeleteInodeRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	deleteInodeRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (deleteInodeRpc *DeleteInodeRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := deleteInodeRpc.metaServerClient.DeleteInode(ctx, deleteInodeRpc.Request)
	output.ShowRpcData(deleteInodeRpc.Request, response, deleteInodeRpc.Info.RpcDataShow)
	return response, err
}

func (inodeRpc *GetInodeRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	inodeRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (inodeRpc *GetInodeRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := inodeRpc.metaServerClient.GetInode(ctx, inodeRpc.Request)
	output.ShowRpcData(inodeRpc.Request, response, inodeRpc.Info.RpcDataShow)
	return response, err
}

func (createDentryRpc *CreateDentryRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	createDentryRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (createDentryRpc *CreateDentryRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := createDentryRpc.metaServerClient.CreateDentry(ctx, createDentryRpc.Request)
	output.ShowRpcData(createDentryRpc.Request, response, createDentryRpc.Info.RpcDataShow)
	return response, err
}

func (listDentryRpc *ListDentryRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	listDentryRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (listDentryRpc *ListDentryRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := listDentryRpc.metaServerClient.ListDentry(ctx, listDentryRpc.Request)
	output.ShowRpcData(listDentryRpc.Request, response, listDentryRpc.Info.RpcDataShow)
	return response, err
}

func (getDentryRpc *GetDentryRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	getDentryRpc.metaServerClient = metaserver.NewMetaServerServiceClient(cc)
}

func (getDentryRpc *GetDentryRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := getDentryRpc.metaServerClient.GetDentry(ctx, getDentryRpc.Request)
	output.ShowRpcData(getDentryRpc.Request, response, getDentryRpc.Info.RpcDataShow)
	return response, err
}

func (getFsStatsRpc *GetFsStatsRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	getFsStatsRpc.mdsClient = mds.NewMdsServiceClient(cc)
}

func (getFsStatsRpc *GetFsStatsRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := getFsStatsRpc.mdsClient.GetFsStats(ctx, getFsStatsRpc.Request)
	output.ShowRpcData(getFsStatsRpc.Request, response, getFsStatsRpc.Info.RpcDataShow)
	return response, err
}

func (listFsRpc *ListClusterFsRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	listFsRpc.mdsClient = mds.NewMdsServiceClient(cc)
}

func (listFsRpc *ListClusterFsRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := listFsRpc.mdsClient.ListClusterFsInfo(ctx, listFsRpc.Request)
	output.ShowRpcData(listFsRpc.Request, response, listFsRpc.Info.RpcDataShow)
	return response, err
}

func (listCacheGroup *ListCacheGroupRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	listCacheGroup.cacheGroupClient = cachegroup.NewCacheGroupMemberServiceClient(cc)
}

func (listCacheGroup *ListCacheGroupRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := listCacheGroup.cacheGroupClient.LoadGroups(ctx, listCacheGroup.Request)
	output.ShowRpcData(listCacheGroup.Request, response, listCacheGroup.Info.RpcDataShow)
	return response, err
}

func (listCacheMember *ListCacheMemberRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	listCacheMember.cacheGroupClient = cachegroup.NewCacheGroupMemberServiceClient(cc)
}

func (listCacheMember *ListCacheMemberRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := listCacheMember.cacheGroupClient.LoadMembers(ctx, listCacheMember.Request)
	output.ShowRpcData(listCacheMember.Request, response, listCacheMember.Info.RpcDataShow)
	return response, err
}

func (rewightMember *ReweightMemberRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	rewightMember.cacheGroupClient = cachegroup.NewCacheGroupMemberServiceClient(cc)
}

func (rewightMember *ReweightMemberRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := rewightMember.cacheGroupClient.ReweightMember(ctx, rewightMember.Request)
	output.ShowRpcData(rewightMember.Request, response, rewightMember.Info.RpcDataShow)
	return response, err
}
