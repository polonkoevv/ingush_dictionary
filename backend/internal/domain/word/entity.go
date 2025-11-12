package word

type Word struct {
	WordID      int    `json:"word_id" db:"word_id"`
	Word        string `json:"word" db:"word"`
	SpeechPart  string `json:"speech_part" db:"speech_part"`
	Translation string `json:"translation" db:"translation"`
	Topic       string `json:"topic" db:"topic"`
	DictID      int    `json:"dict_id" db:"dict_id"`
	DictAbbr    string `json:"abbr" db:"abbr"`
}
