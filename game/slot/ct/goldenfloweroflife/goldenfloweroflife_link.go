//go:build !prod || full || ct

package goldenfloweroflife

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed goldenfloweroflife_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Golden Flower Of Life", LNum: 20, Date: game.Date(2013, 6, 15)}, // see: https://www.slotsmate.com/software/ct-interactive/golden-flower-of-life
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
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
	game.DataRouter["ctinteractive/goldenfloweroflife/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
