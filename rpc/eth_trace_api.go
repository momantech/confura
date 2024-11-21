package rpc

import (
	"context"

	"github.com/Conflux-Chain/confura/rpc/handler"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go/types"
)

// ethTraceAPI provides evm space trace RPC proxy API.
type ethTraceAPI struct {
	stateHandler *handler.EthStateHandler
}

func (api *ethTraceAPI) Block(ctx context.Context, blockNumOrHash types.BlockNumberOrHash) ([]types.LocalizedTrace, error) {
	w3c := GetEthClientFromContext(ctx)
	return api.stateHandler.TraceBlock(ctx, w3c, blockNumOrHash)
}

func (api *ethTraceAPI) Filter(ctx context.Context, filter types.TraceFilter) ([]types.LocalizedTrace, error) {
	w3c := GetEthClientFromContext(ctx)
	return api.stateHandler.TraceFilter(ctx, w3c, filter)
}

func (api *ethTraceAPI) Transaction(ctx context.Context, txHash common.Hash) ([]types.LocalizedTrace, error) {
	w3c := GetEthClientFromContext(ctx)
	return api.stateHandler.TraceTransaction(ctx, w3c, txHash)
}
