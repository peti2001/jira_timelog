package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/peti2001/jira-time-log/service"
)

func main() {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	cookies, err := ioutil.ReadFile("./cookies.txt")
	if err != nil {
		panic(err)
	}
	issues, err := ioutil.ReadFile("./issues.txt")
	if err != nil {
		panic(err)
	}

	jiraApi := service.JiraApifactory(
		"https://innonic.atlassian.net",
		string(cookies),
	)
	issuesSlice := strings.Split(string(issues), "\n")
	sheet := make([]string, 0)
	issueQueue := make(chan string, 15)
	for _, key := range issuesSlice {
		if key == "" {
			continue
		}
		wg.Add(1)
		issueQueue <- key
		go func(k string) {
			issue, err := jiraApi.GetIssue(k)
			if err != nil {
				panic(err)
			}
			for _, worklog := range issue.Fields.Worklog.Worklogs {
				worklog.Comment = strings.Replace(worklog.Comment, "\n", " ", -1)
				mutex.Lock()
				sheet = append(
					sheet,
					fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%d\n", worklog.Author.EmailAddress, worklog.Comment, worklog.Created, issue.Fields.BudgetOwner.Name, k, issue.Fields.Summary, worklog.TimeSpentSeconds),
				)
				mutex.Unlock()
			}
			wg.Done()
			<-issueQueue
		}(key)
	}
	wg.Wait()
	fmt.Println("Email\tCommit\tCreated at\tBudget Owner\tIssue key\tIssue summary\tTime Spent (sec)")
	for _, row := range sheet {
		fmt.Print(row)
	}

}
