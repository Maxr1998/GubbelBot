package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strconv"
	"strings"
)

const (
	bRune         = 0x1F171
	b             = string(bRune)
	mNameNormal   = "Normal"
	mNameAdvanced = "Advanced"
	mNameB        = b + b + b
)

// Mode specifies the amount of B you want in your messages
type Mode int

const (
	mNormal Mode = iota
	mAdvanced
	mB
)

var (
	normalReplacements   = []rune{'b', 'p'}
	advancedReplacements = []rune{'b', 'p', 'g', 'n', 'm', 'd'}

	normalReplacer   = getReplacer(normalReplacements)
	advancedReplacer = getReplacer(advancedReplacements)
	bReplacer        = func(s string) string {
		runes := []rune(s)
		for i, r := range runes {
			if r != ' ' {
				runes[i] = bRune
			}
		}
		return string(runes)
	}
)

func getReplacer(replacements []rune) *strings.Replacer {
	params := make([]string, 0, len(replacements)*4)

	for _, replace := range replacements {
		params = append(params, string(replace), b, strings.ToUpper(string(replace)), b)
	}

	return strings.NewReplacer(params...)
}

func main() {
	// Register API
	bot, err := tgbotapi.NewBotAPI(APIKey)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create and listen to update channel
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		query := update.InlineQuery
		if query == nil {
			continue
		}
		if q := query.Query; q != "" {
			go func() {
				resultNormal := generateResult(query, mNormal)
				resultAdvanced := generateResult(query, mAdvanced)
				resultB := generateResult(query, mB)
				responseConfig := tgbotapi.InlineConfig{
					InlineQueryID: query.ID,
					IsPersonal:    true,
					CacheTime:     0,
					Results:       []interface{}{resultNormal, resultAdvanced, resultB},
				}

				if _, err := bot.AnswerInlineQuery(responseConfig); err != nil {
					log.Println(err)
				}
			}()
		}
	}
}

func generateResult(query *tgbotapi.InlineQuery, mode Mode) tgbotapi.InlineQueryResultArticle {
	var resultModeName, resultText string
	switch mode {
	case mNormal:
		resultModeName = mNameNormal
		resultText = normalReplacer.Replace(query.Query)
	case mAdvanced:
		resultModeName = mNameAdvanced
		resultText = advancedReplacer.Replace(query.Query)
	case mB:
		resultModeName = mNameB
		resultText = bReplacer(query.Query)
	}

	result := tgbotapi.NewInlineQueryResultArticle(query.ID+strconv.Itoa(int(mode)), resultModeName, resultText)
	result.Description = resultText
	return result
}
