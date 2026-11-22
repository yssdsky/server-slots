//go:build !prod || full || ct

package nordicsong

import (
	_ "embed"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
)

//go:embed nordicsong_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Nordic Song", LNum: 25, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/nordic-song
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
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
		ReelsBon.SetPadding(ai.SY)
	},
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/nordicsong/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/nordicsong/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
