package bereal

type Memory struct {
	ID        string `json:"id"`
	Thumbnail struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnail"`
	Primary struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"primary"`
	Secondary struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"secondary"`
	IsLate    bool   `json:"isLate"`
	MemoryDay string `json:"memoryDay"`
}
