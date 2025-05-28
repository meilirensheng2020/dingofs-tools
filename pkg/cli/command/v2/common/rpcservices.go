package common

import (
	"context"

	basecmd "github.com/dingodb/dingofs-tools/pkg/cli/command"
	"github.com/dingodb/dingofs-tools/pkg/output"
	pbmdsv2 "github.com/dingodb/dingofs-tools/proto/dingofs/proto/mdsv2"
	"google.golang.org/grpc"
)

// rpc services
type CreateFsRpc struct {
	Info      *basecmd.Rpc
	Request   *pbmdsv2.CreateFsRequest
	mdsClient pbmdsv2.MDSServiceClient
}

type DeleteFsRpc struct {
	Info      *basecmd.Rpc
	Request   *pbmdsv2.DeleteFsRequest
	mdsClient pbmdsv2.MDSServiceClient
}

type ListFsRpc struct {
	Info      *basecmd.Rpc
	Request   *pbmdsv2.ListFsInfoRequest
	mdsClient pbmdsv2.MDSServiceClient
}

type GetFsRpc struct {
	Info      *basecmd.Rpc
	Request   *pbmdsv2.GetFsInfoRequest
	mdsClient pbmdsv2.MDSServiceClient
}

type GetMdsRpc struct {
	Info      *basecmd.Rpc
	Request   *pbmdsv2.GetMDSListRequest
	mdsClient pbmdsv2.MDSServiceClient
}

type SetFsQuotaRpc struct {
	Info      *basecmd.Rpc
	Request   *pbmdsv2.SetFsQuotaRequest
	mdsClient pbmdsv2.MDSServiceClient
}

type GetFsQuotaRpc struct {
	Info      *basecmd.Rpc
	Request   *pbmdsv2.GetFsQuotaRequest
	mdsClient pbmdsv2.MDSServiceClient
}

type GetInodeRpc struct {
	Info      *basecmd.Rpc
	Request   *pbmdsv2.GetInodeRequest
	mdsClient pbmdsv2.MDSServiceClient
}

type ListDentryRpc struct {
	Info      *basecmd.Rpc
	Request   *pbmdsv2.ListDentryRequest
	mdsClient pbmdsv2.MDSServiceClient
}

var _ basecmd.RpcFunc = (*CreateFsRpc)(nil)   // check interface
var _ basecmd.RpcFunc = (*DeleteFsRpc)(nil)   // check interface
var _ basecmd.RpcFunc = (*ListFsRpc)(nil)     // check interface
var _ basecmd.RpcFunc = (*GetFsRpc)(nil)      // check interface
var _ basecmd.RpcFunc = (*GetMdsRpc)(nil)     // check interface
var _ basecmd.RpcFunc = (*SetFsQuotaRpc)(nil) // check interface
var _ basecmd.RpcFunc = (*GetFsQuotaRpc)(nil) // check interface
var _ basecmd.RpcFunc = (*GetInodeRpc)(nil)   // check interface
var _ basecmd.RpcFunc = (*ListDentryRpc)(nil) // check interface

func (createFs *CreateFsRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	createFs.mdsClient = pbmdsv2.NewMDSServiceClient(cc)
}

func (createFs *CreateFsRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := createFs.mdsClient.CreateFs(ctx, createFs.Request)
	output.ShowRpcData(createFs.Request, response, createFs.Info.RpcDataShow)
	return response, err
}

func (deleteFs *DeleteFsRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	deleteFs.mdsClient = pbmdsv2.NewMDSServiceClient(cc)
}

func (deleteFs *DeleteFsRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := deleteFs.mdsClient.DeleteFs(ctx, deleteFs.Request)
	output.ShowRpcData(deleteFs.Request, response, deleteFs.Info.RpcDataShow)
	return response, err
}

func (listFs *ListFsRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	listFs.mdsClient = pbmdsv2.NewMDSServiceClient(cc)
}

func (listFs *ListFsRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := listFs.mdsClient.ListFsInfo(ctx, listFs.Request)
	output.ShowRpcData(listFs.Request, response, listFs.Info.RpcDataShow)
	return response, err
}

func (getFs *GetFsRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	getFs.mdsClient = pbmdsv2.NewMDSServiceClient(cc)
}

func (getFs *GetFsRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := getFs.mdsClient.GetFsInfo(ctx, getFs.Request)
	output.ShowRpcData(getFs.Request, response, getFs.Info.RpcDataShow)
	return response, err
}

func (getMds *GetMdsRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	getMds.mdsClient = pbmdsv2.NewMDSServiceClient(cc)
}

func (getMds *GetMdsRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := getMds.mdsClient.GetMDSList(ctx, getMds.Request)
	output.ShowRpcData(getMds.Request, response, getMds.Info.RpcDataShow)
	return response, err
}

func (setFsQuota *SetFsQuotaRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	setFsQuota.mdsClient = pbmdsv2.NewMDSServiceClient(cc)
}

func (setFsQuota *SetFsQuotaRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := setFsQuota.mdsClient.SetFsQuota(ctx, setFsQuota.Request)
	output.ShowRpcData(setFsQuota.Request, response, setFsQuota.Info.RpcDataShow)
	return response, err
}

func (getFsQuota *GetFsQuotaRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	getFsQuota.mdsClient = pbmdsv2.NewMDSServiceClient(cc)
}

func (getFsQuota *GetFsQuotaRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := getFsQuota.mdsClient.GetFsQuota(ctx, getFsQuota.Request)
	output.ShowRpcData(getFsQuota.Request, response, getFsQuota.Info.RpcDataShow)
	return response, err
}

func (getInode *GetInodeRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	getInode.mdsClient = pbmdsv2.NewMDSServiceClient(cc)
}

func (getInode *GetInodeRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := getInode.mdsClient.GetInode(ctx, getInode.Request)
	output.ShowRpcData(getInode.Request, response, getInode.Info.RpcDataShow)
	return response, err
}

func (listDentry *ListDentryRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	listDentry.mdsClient = pbmdsv2.NewMDSServiceClient(cc)
}

func (listDentry *ListDentryRpc) Stub_Func(ctx context.Context) (interface{}, error) {
	response, err := listDentry.mdsClient.ListDentry(ctx, listDentry.Request)
	output.ShowRpcData(listDentry.Request, response, listDentry.Info.RpcDataShow)
	return response, err
}
