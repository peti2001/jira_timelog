package model

type Filter struct {
	SearchUrl string `json:"searchUrl"`
}

type FilterResult struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}
