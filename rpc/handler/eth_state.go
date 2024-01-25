package handler

import (
	"context"
	"math/big"

	"github.com/Conflux-Chain/confura/node"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go/types"
)

// EthStateHandler handles evm space state RPC method by redirecting requests to another
// full state node if state is not available on normal full node.
type EthStateHandler struct {
	cp *node.EthClientProvider
}

func NewEthStateHandler(cp *node.EthClientProvider) *EthStateHandler {
	return &EthStateHandler{cp: cp}
}

func (h *EthStateHandler) Balance(
	ctx context.Context,
	w3c *node.Web3goClient,
	addr common.Address,
	block *types.BlockNumberOrHash,
) (*big.Int, error) {
	bal, err := h.doRequest(ctx, w3c, func(w3c *node.Web3goClient) (interface{}, error) {
		return w3c.Eth.Balance(addr, block)
	})

	if err != nil {
		return nil, err
	}

	return bal.(*big.Int), err
}

func (h *EthStateHandler) TransactionCount(
	ctx context.Context,
	w3c *node.Web3goClient,
	addr common.Address,
	blockNum *types.BlockNumberOrHash,
) (*big.Int, error) {
	txnCnt, err := h.doRequest(ctx, w3c, func(w3c *node.Web3goClient) (interface{}, error) {
		return w3c.Eth.TransactionCount(addr, blockNum)
	})

	if err != nil {
		return nil, err
	}

	return txnCnt.(*big.Int), err
}

func (h *EthStateHandler) StorageAt(
	ctx context.Context,
	w3c *node.Web3goClient,
	addr common.Address,
	location *big.Int,
	block *types.BlockNumberOrHash,
) (common.Hash, error) {
	storage, err := h.doRequest(ctx, w3c, func(w3c *node.Web3goClient) (interface{}, error) {
		return w3c.Eth.StorageAt(addr, location, block)
	})

	if err != nil {
		return common.Hash{}, err
	}

	return storage.(common.Hash), err
}

func (h *EthStateHandler) CodeAt(
	ctx context.Context,
	w3c *node.Web3goClient,
	addr common.Address,
	blockNum *types.BlockNumberOrHash,
) ([]byte, error) {
	code, err := h.doRequest(ctx, w3c, func(w3c *node.Web3goClient) (interface{}, error) {
		return w3c.Eth.CodeAt(addr, blockNum)
	})

	if err != nil {
		return []byte{}, err
	}

	return code.([]byte), err
}

func (h *EthStateHandler) Call(
	ctx context.Context,
	w3c *node.Web3goClient,
	callRequest types.CallRequest,
	blockNum *types.BlockNumberOrHash,
) ([]byte, error) {
	result, err := h.doRequest(ctx, w3c, func(w3c *node.Web3goClient) (interface{}, error) {
		return w3c.Eth.Call(callRequest, blockNum)
	})

	if err != nil {
		return []byte{}, err
	}

	return result.([]byte), err
}

func (h *EthStateHandler) EstimateGas(
	ctx context.Context,
	w3c *node.Web3goClient,
	callRequest types.CallRequest,
	blockNum *types.BlockNumberOrHash,
) (*big.Int, error) {
	est, err := h.doRequest(ctx, w3c, func(w3c *node.Web3goClient) (interface{}, error) {
		return w3c.Eth.EstimateGas(callRequest, blockNum)
	})

	if err != nil {
		return nil, err
	}

	return est.(*big.Int), err
}

func (h *EthStateHandler) doRequest(
	ctx context.Context,
	initW3c *node.Web3goClient,
	clientFunc func(w3c *node.Web3goClient) (interface{}, error),
) (interface{}, error) {
	result, err := clientFunc(initW3c)
	if err == nil || !isStateNotAvailable(err) {
		return result, err
	}

	fsW3c, cperr := h.cp.GetClientByIP(ctx, node.GroupEthFullState)
	if cperr != nil {
		return result, err
	}

	return clientFunc(fsW3c)
}
