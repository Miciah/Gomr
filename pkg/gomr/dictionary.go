package gomr

import (
	"encoding/json"
	"errors"
	"regexp"
)

// This plugin requires a wordnik api key to be defined in configuration
//  Get an api key here: http://developer.wordnik.com/

type DictionaryPlugin struct {
	Nick          string
	WordnikAPIKey string
}

// Struct for json response from api when getting a word definition list
type DefineResp struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Text         string `json:"text"`
}

func (dp DictionaryPlugin) Register() (err error) {
	// If the wordnik api key was not supplied, this plugin is useless
	if dp.WordnikAPIKey == "" {
		return errors.New("Wordnik API key was not provided in configuration")
	}
	return nil
}

func (dp DictionaryPlugin) Parse(user, channel, input string, conn *Connection) (err error) {
	if Match(input, `(?i)`+dp.Nick+`[\S]* define[:]*\s+[\S]{2,}`) {
		wrgx, _ := regexp.Compile(`define[:]*\s+([\S]{2,})`)
		wordMatch := wrgx.FindStringSubmatch(input)

		if wordMatch != nil && len(wordMatch) > 1 {
			word := wordMatch[1]
			url := "http://api.wordnik.com:80/v4/word.json/" + word +
				"/definitions?limit=1&includeRelated=false&useCanonical=true&includeTags=false&api_key=" +
				dp.WordnikAPIKey

			var r []byte
			r, err = HttpGet(url)
			if err != nil {
				err = errors.New("ERROR unable to get " + url + " : " + err.Error())
				return
			}
			resp := []DefineResp{}
			json.Unmarshal(r, &resp)
			if len(resp) == 0 {
				conn.SendTo(channel, "No definition found for "+word+".")
				return
			}
			definition := resp[0].Text
			conn.SendTo(channel, word+": "+definition)
		} else {
			err = errors.New("Dictionary Error: unable to get word definition string from input:" + input)
			return
		}

	}

	return
}

func (dp DictionaryPlugin) Help() (texts []string) {
	texts = append(texts, dp.Nick+"[:] define[:] <word>")
	return texts
}
