package main

import (
	"fmt"
	"os"

	"github.com/yanzay/bstools"
	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
)

var template = `Upgrade %s
💰 %s
🌲 %s
⛏ %s

Income increase: %d gold/min
Time to upgrade: %s`

func main() {
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Handle("/start", "Forward your 🏘 Buildings here")
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
	log.Infof(fmt.Sprint(state))
	up := state.BalancedUpgrade()
	reply := fmt.Sprintf(template, up.Type, comma(up.Gold), comma(up.Wood), comma(up.Stone), state.IncomeDelta(up.Type), up.Duration)
	m.Reply("```\n"+reply+"```", tbot.WithMarkdown)
	log.Infof("Recommendation: %s", up.Type)
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
