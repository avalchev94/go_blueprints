package thesaurus

import (
	"encoding/json"
	"errors"
	"net/http"
)

type BigHugh struct {
	APIKEY string
}

type words struct {
	Syn []string `json:"syn"`
}

type synonyms struct {
	Noun      *words `json:"noun"`
	Verb      *words `json:"verb"`
	Adverb    *words `json:"adverb"`
	Adjective *words `json:"adjective"`
}

func (s *synonyms) getSynonyms() (syns []string) {
	if s.Noun != nil && len(s.Noun.Syn) > 0 {
		syns = append(syns, s.Noun.Syn...)
	}

	if s.Verb != nil && len(s.Verb.Syn) > 0 {
		syns = append(syns, s.Verb.Syn...)
	}

	if s.Adverb != nil && len(s.Adverb.Syn) > 0 {
		syns = append(syns, s.Adverb.Syn...)
	}

	if s.Adjective != nil && len(s.Adjective.Syn) > 0 {
		syns = append(syns, s.Adjective.Syn...)
	}

	return syns
}

// Synonyms only thing about this package. Just call it with the words whoose
// synonyms are needed.
func (b *BigHugh) Synonyms(term string) ([]string, error) {
	var syns []string
	response, err := http.Get("http://words.bighugelabs.com/api/2/" +
		b.APIKEY + "/" + term + "/json")
	if err != nil {
		return syns, errors.New("bighugh: Failed when looking for synonyms for \"" +
			term + "\"" + err.Error())
	}

	var data synonyms
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return syns, err
	}

	syns = data.getSynonyms()
	return syns, nil
}
