package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yanzay/bstools"
	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
)

var template = `Passive Income:
Upgrade %s
%s
Income increase: %d gold/min
Time to upgrade: %s

Rush Barracks:
Upgrade %s
%s
Time to upgrade: %s

Balanced Battle:
Upgrade %s
%s
Time to upgrade: %s
`

var bStore *BuildStore

var (
	dbFile = flag.String("data", "bsupgrade.db", "Database file name")
)

func main() {
	flag.Parse()
	bStore = NewBuildStore(*dbFile)
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Handle("/start", "Forward your 🏘 Buildings and ⚒ Workshop here")
	bot.HandleDefault(parserHandler)
	bot.ListenAndServe()
}

func parserHandler(m *tbot.Message) {
	log.Infof("Message from [%d] %s %s (%s)\n%s", m.From.ID, m.From.FirstName, m.From.LastName, m.From.UserName, m.Text())
	state, err := bstools.ParseState(m.Text())
	if err != nil {
		m.Reply(err.Error())
		return
	}
	savedState := bStore.GetBuildings(fmt.Sprint(m.From.ID))
	if err != nil {
		m.Reply(err.Error())
		return
	}
	state.Merge(savedState)
	log.Infof(fmt.Sprint(state))
	bStore.SaveBuildings(fmt.Sprint(m.From.ID), state)
	err = state.Valid()
	if err != nil {
		m.Reply(err.Error())
		return
	}
	balUp := state.BalancedUpgrade()
	rushUp := state.RushUpgrade()
	batUp := state.BattleUpgrade()

	reply := fmt.Sprintf(template,
		balUp.Type, printPrice(balUp), state.IncomeDelta(balUp.Type), balUp.Duration,
		rushUp.Type, printPrice(rushUp), rushUp.Duration,
		batUp.Type, printPrice(batUp), batUp.Duration)
	m.Reply("```\n"+reply+"```", tbot.WithMarkdown)
	log.Infof("Recommendation: %s, %s, %s", balUp.Type, rushUp.Type, batUp.Type)
}

func printPrice(up *bstools.Upgrade) string {
	return fmt.Sprintf("💰 %s\n🌲 %s\n⛏ %s", comma(up.Gold), comma(up.Wood), comma(up.Stone))
}

func comma(n int) string {
	runes := []rune(fmt.Sprint(n))
	commaRunes := make([]rune, 0)
	for i, r := range runes {
		commaRunes = append(commaRunes, r)
		pos := len(runes) - i - 1
		if pos != 0 && pos%3 == 0 {
			commaRunes = append(commaRunes, ',')
		}
	}
	return string(commaRunes)
}
