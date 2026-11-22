//go:build !prod || full || agt

package egypt

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed egypt_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Egypt", LNum: 15, Date: game.Year(2025)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPrgmult |
			game.GPscat |
			game.GPwild,
		SX: 5,
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
	game.DataRouter["agt/egypt/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
