package pair

import (
	"context"

	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/l2asset"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/liquidity"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/svc"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvailablePairsLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	liquidity liquidity.LiquidityModel
	l2asset   l2asset.L2asset
}

func NewGetAvailablePairsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvailablePairsLogic {
	return &GetAvailablePairsLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		liquidity: liquidity.New(svcCtx),
		l2asset:   l2asset.New(svcCtx),
	}
}

func (l *GetAvailablePairsLogic) GetAvailablePairs(_ *types.ReqGetAvailablePairs) (*types.RespGetAvailablePairs, error) {
	liquidityAssets, err := l.liquidity.GetAllLiquidityAssets()
	if err != nil {
		logx.Error("[GetAllLiquidityAssets] err:%v", err)
		return nil, err
	}
	resp := &types.RespGetAvailablePairs{}
	for _, asset := range liquidityAssets {
		assetA, err := l.l2asset.GetSimpleL2AssetInfoByAssetId(l.ctx, uint32(asset.AssetAId))
		if err != nil {
			logx.Error("[GetSimpleL2AssetInfoByAssetId] err:%v", err)
			return nil, err
		}
		assetB, err := l.l2asset.GetSimpleL2AssetInfoByAssetId(l.ctx, uint32(asset.AssetBId))
		if err != nil {
			logx.Error("[GetSimpleL2AssetInfoByAssetId] err:%v", err)
			return nil, err
		}
		resp.Pairs = append(resp.Pairs, &types.Pair{
			PairIndex:    uint32(asset.PairIndex),
			AssetAId:     uint32(asset.AssetAId),
			AssetAName:   assetA.AssetName,
			AssetAAmount: asset.AssetA,
			AssetBId:     uint32(asset.AssetBId),
			AssetBName:   assetB.AssetName,
			AssetBAmount: asset.AssetB,
			FeeRate:      asset.FeeRate,
			TreasuryRate: asset.TreasuryRate,
		})
	}
	return resp, nil
}
