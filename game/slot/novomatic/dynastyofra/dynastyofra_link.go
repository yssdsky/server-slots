//go:build !prod || full || novomatic

package dynastyofra

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed dynastyofra_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Dynasty of Ra", LNum: 10, Date: game.Date(2014, 1, 29)}, // see: https://www.slotsmate.com/software/novomatic/dynasty-of-ra
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
			game.GPwsc,
		SX: 5,
		SY: 3,
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
	game.DataRouter["novomatic/dynastyofra/bon"] = &ReelsBon
	game.DataRouter["novomatic/dynastyofra/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
