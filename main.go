package main

import (
	"fmt"
	"io/ioutil"

	"github.com/peti2001/jira-time-log/service"
)

func main() {

	cookies, err := ioutil.ReadFile("./cookies.txt")
	if err != nil {
		panic(err)
	}

	jiraApi := service.JiraApifactory(
		"https://innonic.atlassian.net",
		string(cookies),
	)

	issues, err := jiraApi.GetIssuesByFilter("16100")
	if err != nil {
		panic(err)
	}

	sheet := make([]string, 0)
	for _, issue := range issues {
		if err != nil {
			panic(err)
		}
		for _, worklog := range issue.Fields.Worklog.Worklogs {
			sheet = append(
				sheet,
				fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%d\n", worklog.Author.EmailAddress, worklog.Started, issue.Fields.BudgetOwner.Name, issue.Key, issue.Fields.Summary, worklog.TimeSpentSeconds),
			)
		}
	}
	fmt.Println("Email\tStarted at\tBudget Owner\tIssue key\tIssue summary\tTime Spent (sec)")
	for _, row := range sheet {
		fmt.Print(row)
	}

}
