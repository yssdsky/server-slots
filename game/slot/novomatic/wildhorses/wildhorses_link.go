//go:build !prod || full || novomatic

package wildhorses

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed wildhorses_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Wild Horses", LNum: 50, Date: game.Date(2016, 12, 4)}, // see: https://www.slotsmate.com/software/novomatic/wild-horses
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPpick |
			game.GPfgseq |
			game.GPfgreel |
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
		ReelsBon.SetPadding(ai.SY)
	},
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["novomatic/wildhorses/bon"] = &ReelsBon
	game.DataRouter["novomatic/wildhorses/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
