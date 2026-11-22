//go:build !prod || full || agt

package merrychristmas

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed merrychristmas_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Merry Christmas", LNum: 5}, // see: https://agtsoftware.com/games/agt/christmas
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPfgseq |
			game.GPscat |
			game.GPwild,
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
	game.DataRouter["agt/merrychristmas/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
