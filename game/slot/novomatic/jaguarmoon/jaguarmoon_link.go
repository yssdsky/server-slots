//go:build !prod || full || novomatic

package jaguarmoon

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed jaguarmoon_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Jaguar Moon", Date: game.Date(2018, 7, 11)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: sn,
		WN: 243,
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) {
		ai.RTP = game.MakeRtpList(ReelsMap)
		slot.MapSetPadding(ReelsMap, ai.SY)
		ReelsBon.SetPadding(ai.SY)
	},
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["novomatic/jaguarmoon/bon"] = &ReelsBon
	game.DataRouter["novomatic/jaguarmoon/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
