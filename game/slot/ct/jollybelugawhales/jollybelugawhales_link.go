//go:build !prod || full || ct

package jollybelugawhales

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed jollybelugawhales_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Jolly Beluga Whales", LNum: 20, Date: game.Date(2012, 6, 15)}, // see: https://www.slotsmate.com/software/ct-interactive/jolly-beluga-whales
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPrwild,
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
	game.DataRouter["ctinteractive/jollybelugawhales/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
