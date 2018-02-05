package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/peti2001/jira-time-log/model"
	"sync"
	"log"
)

type jiraApi struct {
	url     string
	cookies []*http.Cookie
}

func JiraApifactory(url string, cookiesText string) *jiraApi {
	cookies := strings.Split(cookiesText, ";")
	httpCookies := make([]*http.Cookie, len(cookies))
	for _, cookie := range cookies {
		splittedCookie := strings.Split(strings.Trim(cookie, " "), "=")
		if len(splittedCookie) == 2 {
			httpCookies = append(httpCookies, &http.Cookie{Name: splittedCookie[0], Value: splittedCookie[1]})
		}
	}

	return &jiraApi{
		url + "/rest/api/2",
		httpCookies,
	}
}

func (j *jiraApi) makeRequest(url string) ([]byte, error) {
	log.Println("Make request: " + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for _, c := range j.cookies {
		if c != nil {
			req.AddCookie(c)
		}
	}
	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (j *jiraApi) getFilterResult(url string) (*model.FilterResult, error) {
	ret, err := j.makeRequest(url)
	if err != nil {
		return nil, err
	}
	var filterResult model.FilterResult
	err = json.Unmarshal(ret, &filterResult)
	if err != nil {
		return nil, err
	}

	return &filterResult, nil
}

func (j *jiraApi) GetIssuesByFilter(filterId string) ([]*model.Issue, error) {
	ret, err := j.makeRequest(j.url + "/filter/" + filterId)
	if err != nil {
		return nil, err
	}

	var filter model.Filter
	err = json.Unmarshal(ret, &filter)
	if err != nil {
		fmt.Println(string(ret))
		return nil, err
	}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	maxResult := 50
	startsAt := 0
	filterResult, err := j.getFilterResult(fmt.Sprintf("%s&startAt=0&maxResults=%d", filter.SearchUrl, maxResult))
	issues := make([]*model.Issue, filterResult.Total)
	i := 0
	for startsAt < filterResult.Total {
		filterResult, err = j.getFilterResult(fmt.Sprintf("%s&startAt=%d&maxResults=%d", filter.SearchUrl, startsAt, maxResult))
		issueQueue := make(chan string, 10)

		for _, issue := range filterResult.Issues {
			if issue.Key == "" {
				continue
			}
			wg.Add(1)
			issueQueue <- issue.Key
			go func(k string) {
				issue, err := j.GetIssue(k)
				if err != nil {
					log.Println(err)
				}
				mutex.Lock()
				issues[i] = issue
				i++
				mutex.Unlock()
				wg.Done()
				<-issueQueue
			}(issue.Key)
		}
		wg.Wait()
		startsAt += maxResult
	}

	return issues, nil
}

func (j *jiraApi) GetIssue(key string) (*model.Issue, error) {
	ret, err := j.makeRequest(j.url + "/issue/" + key)

	var jiraIssue model.Issue
	err = json.Unmarshal(ret, &jiraIssue)
	if err != nil {
		fmt.Println(string(ret))
		return nil, err
	}

	return &jiraIssue, nil
}
