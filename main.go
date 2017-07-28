package main

import (
	"fmt"
	"os"

	"github.com/yanzay/bstools"
	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
)

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
	state, err := bstools.ParseState(m.Text())
	if err != nil {
		m.Reply(err.Error())
		return
	}
	log.Infof(fmt.Sprint(state))
	up := state.BalancedUpgrade()
	reply := fmt.Sprintf("%s\n💰 %s\n🌲 %s\n⛏ %s", up.Type, comma(up.Gold), comma(up.Wood), comma(up.Stone))
	m.Reply("```\n"+reply+"```", tbot.WithMarkdown)
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
