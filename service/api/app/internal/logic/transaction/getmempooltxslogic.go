package transaction

import (
	"context"

	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/logic/errcode"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/block"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/mempool"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/svc"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/types"
	"github.com/zecrey-labs/zecrey-legend/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMempoolTxsLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	mempool mempool.Mempool
	block   block.Block
}

func NewGetMempoolTxsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMempoolTxsLogic {
	return &GetMempoolTxsLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		mempool: mempool.New(svcCtx.Config),
		block:   block.New(svcCtx.Config),
	}
}
func (l *GetMempoolTxsLogic) GetMempoolTxs(req *types.ReqGetMempoolTxs) (*types.RespGetMempoolTxs, error) {
	if utils.CheckTypeLimit(req.Limit) {
		logx.Errorf("[CheckTypeLimit] param:%v", req.Limit)
		return nil, errcode.ErrInvalidParam
	}
	if utils.CheckTypeOffset(req.Offset) {
		logx.Errorf("[CheckTypeOffset] param:%v", req.Offset)
		return nil, errcode.ErrInvalidParam
	}
	resp := &types.RespGetMempoolTxs{
		MempoolTxs: make([]*types.Tx, 0),
	}
	var err error
	resp.Total, err = l.mempool.GetMempoolTxsTotalCount()
	if err != nil {
		logx.Errorf("[GetMempoolTxsTotalCount] err:%v", err)
		return nil, err
	}
	mempoolTxs, err := l.mempool.GetMempoolTxs(int64(req.Limit), int64(req.Offset))
	if err != nil {
		logx.Errorf("[GetMempoolTxs] err:%v", err)
		return nil, err
	}
	for _, mempoolTx := range mempoolTxs {
		txDetails := make([]*types.TxDetail, 0)
		for _, txDetail := range mempoolTx.MempoolDetails {
			txDetails = append(txDetails, &types.TxDetail{
				AssetId:      uint32(txDetail.AssetId),
				AssetType:    uint32(txDetail.AssetType),
				AccountIndex: int32(txDetail.AccountIndex),
				AccountName:  txDetail.AccountName,
			})
		}
		if err != nil {
			logx.Error("[GetMempoolTxs] err:%v", err)
			return &types.RespGetMempoolTxs{}, err
		}
		blockInfo, err := l.block.GetBlockByBlockHeight(mempoolTx.L2BlockHeight)
		if err != nil {
			logx.Error("[transaction.GetTxsByAccountName] err:%v", err)
			return nil, err
		}
		resp.MempoolTxs = append(resp.MempoolTxs, &types.Tx{
			TxHash:        mempoolTx.TxHash,
			TxType:        uint32(mempoolTx.TxType),
			GasFeeAssetId: uint32(mempoolTx.GasFeeAssetId),
			GasFee:        mempoolTx.GasFee,
			NftIndex:      uint32(mempoolTx.NftIndex),
			PairIndex:     uint32(mempoolTx.PairIndex),
			AssetId:       uint32(mempoolTx.AssetId),
			TxAmount:      mempoolTx.TxAmount,
			NativeAddress: mempoolTx.NativeAddress,
			TxDetails:     txDetails,
			TxInfo:        mempoolTx.TxInfo,
			ExtraInfo:     mempoolTx.ExtraInfo,
			Memo:          mempoolTx.Memo,
			AccountIndex:  uint32(mempoolTx.AccountIndex),
			Nonce:         uint32(mempoolTx.Nonce),
			ExpiredAt:     uint32(mempoolTx.ExpiredAt),
			L2BlockHeight: uint32(mempoolTx.L2BlockHeight),
			Status:        uint32(mempoolTx.Status),
			CreatedAt:     uint32(mempoolTx.CreatedAt.Unix()),
			BlockID:       uint32(blockInfo.ID),
		})
	}
	return resp, nil
}
