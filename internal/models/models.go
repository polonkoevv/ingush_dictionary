package models

type WordResponse struct {
	CurrentPage  int    `json:"current_page"`
	Data         []Data `json:"data"`
	FirstPageURL string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	LastPageURL  string `json:"last_page_url"`
	NextPageURL  string `json:"next_page_url"`
	Path         string `json:"path"`
	PerPage      string `json:"per_page"`
	PrevPageURL  string `json:"prev_page_url"`
	To           int    `json:"to"`
	Total        int    `json:"total"`
}

type Data struct {
	ID         int    `json:"id"`
	Word       string `json:"word"`
	Language   string `json:"language"`
	WordID     int    `json:"word_id"`
	Words      []Word `json:"words"`
	Translates []Word `json:"translates"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Audio      Audio  `json:"audio"`
}

type Word struct {
	ID              int           `json:"id"`
	Word            string        `json:"word"`
	Examples        []interface{} `json:"examples"`
	PartSpeech      int           `json:"part_speech"`
	Class           int           `json:"class"`
	Decline         string        `json:"decline"`
	CreatedAt       string        `json:"created_at"`
	UpdatedAt       string        `json:"updated_at"`
	ViewerFavorited bool          `json:"viewerFavorited"`
}

type Audio struct {
	ID        int    `json:"id"`
	Filename  string `json:"filename"`
	Thumb     string `json:"thumb"`
	Ext       string `json:"ext"`
	Size      int    `json:"size"`
	Accepted  int    `json:"accepted"`
	ParentID  int    `json:"parent_id"`
	UserID    int64  `json:"user_id"`
	UserIP    string `json:"user_ip"`
	CreatedAt string `json:"created_at"`
	Type      string `json:"type"`
}
