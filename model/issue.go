package model

type Issue struct {
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
}

type Fields struct {
	Worklog     WorklogField `json:"worklog"`
	Summary     string       `json:"summary"`
	BudgetOwner BudgetOwner  `json:"customfield_11501"`
}

type BudgetOwner struct {
	Name string `json:"value"`
}

type WorklogField struct {
	Total    int       `json:"total"`
	Worklogs []Worklog `json:"worklogs"`
}

type Worklog struct {
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	Created          string `json:"created"`
	Started          string `json:"started"`
	Comment          string `json:"comment"`
	Author           Author `json:"author"`
	Key              string
}

type Worklogs []*Worklog

func (w Worklogs) Len() int {
	return len(w)
}

func (w Worklogs) Less(i, j int) bool {
	return w[i].Started < w[j].Started
}

func (w Worklogs) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

type Author struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
}
