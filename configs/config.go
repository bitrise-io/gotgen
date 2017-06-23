package configs

// DelimiterModel ...
type DelimiterModel struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

// Model ...
type Model struct {
	Inventory map[string]interface{} `json:"inventory"`
	Delimiter DelimiterModel         `json:"delimiter"`
}
