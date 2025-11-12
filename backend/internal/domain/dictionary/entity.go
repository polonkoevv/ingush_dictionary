package dictionary

type Dictionary struct {
	DictID int    `json:"dict_id" db:"dict_id"`
	Abbr   string `json:"abbr" db:"abbr"`
	Name   string `json:"name" db:"name"`
	Author string `json:"author" db:"author"`
}
