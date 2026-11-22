//go:build !prod || full || igt

package triplediamond

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed triplediamond_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "IGT", Name: "Triple Diamond", LNum: 9, Date: game.Year(2005)}, // see: https://www.slotsmate.com/software/igt/triple-diamond
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlsel |
			game.GPmix |
			game.GPfgno |
			game.GPwild |
			game.GPwmult,
		SX: 3,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) {
		ai.RTP = game.MakeRtpList(ReelsMap)
		slot.MapSetPadding(ReelsMap, ai.SY)
	},
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["igt/triplediamond/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
