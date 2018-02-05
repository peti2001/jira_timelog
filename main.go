package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/peti2001/jira-time-log/service"
)

type Result []WorklogResult

type WorklogResult struct {
	Email            string `json:"email"`
	Started          string `json:"started"`
	BudgetOwner      string `json:"budgetOwner"`
	Key              string `json:"key"`
	Summary          string `json:"summary"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}

func getWorkLogByFilter(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filterId, _ := r.Form["filterId"]
	cookies, err := ioutil.ReadFile("./cookies.txt")
	if err != nil {
		panic(err)
	}

	jiraApi := service.JiraApifactory(
		"https://innonic.atlassian.net",
		string(cookies),
	)

	issues, err := jiraApi.GetIssuesByFilter(filterId[0])
	if err != nil {
		panic(err)
	}

	sheet := make(Result, 0)
	for i, issue := range issues {
		fmt.Printf("%d - %s\n", i, issue.Key)
		for _, worklog := range issue.Fields.Worklog.Worklogs {
			sheet = append(
				sheet,
				WorklogResult{
					worklog.Author.EmailAddress,
					worklog.Started,
					issue.Fields.BudgetOwner.Name,
					issue.Key,
					issue.Fields.Summary,
					worklog.TimeSpentSeconds,
				},
			)
		}
	}
	jsonResp, err := json.Marshal(sheet)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonResp)
}

func main() {
	http.HandleFunc("/getWorkLogByFilter", getWorkLogByFilter)
	fmt.Println("Start server")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
