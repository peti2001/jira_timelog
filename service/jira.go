package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/peti2001/jira-time-log/model"
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
		url + "/rest/api/latest/issue/",
		httpCookies,
	}
}

func (j *jiraApi) GetIssue(key string) (*model.Issue, error) {
	req, err := http.NewRequest("GET", j.url+key, nil)
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
	var jiraIssue model.Issue
	err = json.Unmarshal(ret, &jiraIssue)
	if err != nil {
		fmt.Println(string(ret))
		return nil, err
	}

	return &jiraIssue, nil
}
