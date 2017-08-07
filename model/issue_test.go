package model_test

import (
	"sort"
	"testing"

	"github.com/peti2001/jira-time-log/model"
	"github.com/stretchr/testify/assert"
)

func TestSortWorklogs(t *testing.T) {
	//Arrange
	issues := model.Worklogs{
		&model.Worklog{
			100,
			"2017-07-20T14:45:00.000+0200",
			"2017-07-20T14:45:00.000+0200",
			"",
			model.Author{
				"test.peter",
				"test.peter@innonic.com",
			},
			"AB-001",
		},
		&model.Worklog{
			100,
			"2017-12-20T14:45:00.000+0200",
			"2017-12-20T14:45:00.000+0200",
			"",
			model.Author{
				"test.peter",
				"test.peter@innonic.com",
			},
			"AB-002",
		},
		&model.Worklog{
			100,
			"2017-12-20T14:45:01.000+0200",
			"2017-12-20T14:45:01.000+0200",
			"",
			model.Author{
				"test.peter",
				"test.peter@innonic.com",
			},
			"AB-003",
		},
		&model.Worklog{
			100,
			"2017-06-20T14:45:00.000+0200",
			"2017-06-20T14:45:00.000+0200",
			"",
			model.Author{
				"test.peter",
				"test.peter@innonic.com",
			},
			"AB-004",
		},
		&model.Worklog{
			100,
			"2010-07-20T13:45:00.000+0200",
			"2010-07-20T13:45:00.000+0200",
			"",
			model.Author{
				"test.peter",
				"test.peter@innonic.com",
			},
			"AB-005",
		},
	}

	//Act
	sort.Sort(issues)

	//Assert
	assert.Equal(t, "AB-005", issues[0].Key)
	assert.Equal(t, "AB-002", issues[3].Key)
	assert.Equal(t, "AB-003", issues[4].Key)
}