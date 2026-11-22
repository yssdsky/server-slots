//go:build !prod || full || agt

package gems50

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed gems50_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "50 Gems", LNum: 50}, // see: https://agtsoftware.com/games/agt/gems50
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 4,
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
	game.DataRouter["agt/50gems/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
