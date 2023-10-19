package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/devoteamnl/opendora/api/models"
	"github.com/devoteamnl/opendora/api/sql_client"
	"github.com/devoteamnl/opendora/api/sql_client/sql_queries"
)

func TestDfService_ServeRequest(t *testing.T) {

	exampleWeeklyDataPoints := []models.DataPoint{{Key: "202338", Value: 0}, {Key: "202337", Value: 1}, {Key: "202336", Value: 2}}
	exampleMonthlyDataPoints := []models.DataPoint{{Key: "23/04", Value: 6}, {Key: "23/03", Value: 5}, {Key: "23/02", Value: 4}}
	exampleQuarterlyDataPoints := []models.DataPoint{{Key: "2024-01-01", Value: 6}, {Key: "2023-10-01", Value: 5}, {Key: "2023-07-01", Value: 4}}
	exampleAverageWeeklyDataPoints := []models.DataPoint{{Key: "202338", Value: 0.75}, {Key: "202337", Value: 1.5}, {Key: "202336", Value: 2}}
	exampleAverageMonthlyDataPoints := []models.DataPoint{{Key: "23/04", Value: 5.25}, {Key: "23/03", Value: 4.5}, {Key: "23/02", Value: 4}}
	exampleAverageQuarterlyDataPoints := []models.DataPoint{{Key: "2024-01-01", Value: 5.25}, {Key: "2023-10-01", Value: 4.5}, {Key: "2023-07-01", Value: 4}}

	dataMockMap := map[string]sql_client.MockDataReturn{
		sql_queries.WeeklyDeploymentSql + sql_queries.CountSql:      {Data: exampleWeeklyDataPoints},
		sql_queries.MonthlyDeploymentSql + sql_queries.CountSql:     {Data: exampleMonthlyDataPoints},
		sql_queries.QuarterlyDeploymentSql + sql_queries.CountSql:   {Data: exampleQuarterlyDataPoints},
		sql_queries.WeeklyDeploymentSql + sql_queries.AverageSql:    {Data: exampleAverageWeeklyDataPoints},
		sql_queries.MonthlyDeploymentSql + sql_queries.AverageSql:   {Data: exampleAverageMonthlyDataPoints},
		sql_queries.QuarterlyDeploymentSql + sql_queries.AverageSql: {Data: exampleAverageQuarterlyDataPoints},
	}

	errorMockMap := map[string]sql_client.MockDataReturn{
		sql_queries.WeeklyDeploymentSql + sql_queries.CountSql:      {Err: fmt.Errorf("error from weekly query")},
		sql_queries.MonthlyDeploymentSql + sql_queries.CountSql:     {Err: fmt.Errorf("error from monthly query")},
		sql_queries.QuarterlyDeploymentSql + sql_queries.CountSql:   {Err: fmt.Errorf("error from quarterly query")},
		sql_queries.WeeklyDeploymentSql + sql_queries.AverageSql:    {Err: fmt.Errorf("error from weekly average query")},
		sql_queries.MonthlyDeploymentSql + sql_queries.AverageSql:   {Err: fmt.Errorf("error from monthly average query")},
		sql_queries.QuarterlyDeploymentSql + sql_queries.AverageSql: {Err: fmt.Errorf("error from quarterly average query")},
	}

	tests := []struct {
		name           string
		params         ServiceParameters
		mockClient     sql_client.MockClient
		expectResponse models.Response
		expectError    string
	}{
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "df_count", Aggregation: "weekly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: errorMockMap},
			expectResponse: models.Response{Aggregation: "weekly", DataPoints: nil},
			expectError:    "error from weekly query",
		},
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "df_count", Aggregation: "monthly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: errorMockMap},
			expectResponse: models.Response{Aggregation: "monthly", DataPoints: nil},
			expectError:    "error from monthly query",
		},
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "df_count", Aggregation: "quarterly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: errorMockMap},
			expectResponse: models.Response{Aggregation: "quarterly", DataPoints: nil},
			expectError:    "error from quarterly query",
		},
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "df_average", Aggregation: "weekly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: errorMockMap},
			expectResponse: models.Response{Aggregation: "weekly", DataPoints: nil},
			expectError:    "error from weekly average query",
		},
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "df_average", Aggregation: "monthly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: errorMockMap},
			expectResponse: models.Response{Aggregation: "monthly", DataPoints: nil},
			expectError:    "error from monthly average query",
		},
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "df_average", Aggregation: "quarterly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: errorMockMap},
			expectResponse: models.Response{Aggregation: "quarterly", DataPoints: nil},
			expectError:    "error from quarterly average query",
		},
		{
			name:           "should return weekly data points from the database",
			params:         ServiceParameters{TypeQuery: "df_count", Aggregation: "weekly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: dataMockMap},
			expectResponse: models.Response{Aggregation: "weekly", DataPoints: exampleWeeklyDataPoints},
			expectError:    "",
		},
		{
			name:           "should return monthly data points from the database",
			params:         ServiceParameters{TypeQuery: "df_count", Aggregation: "monthly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: dataMockMap},
			expectResponse: models.Response{Aggregation: "monthly", DataPoints: exampleMonthlyDataPoints},
			expectError:    "",
		},
		{
			name:           "should return quarterly data points from the database",
			params:         ServiceParameters{TypeQuery: "df_count", Aggregation: "quarterly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: dataMockMap},
			expectResponse: models.Response{Aggregation: "quarterly", DataPoints: exampleQuarterlyDataPoints},
			expectError:    "",
		},
		{
			name:           "should return average weekly data points from the database",
			params:         ServiceParameters{TypeQuery: "df_average", Aggregation: "weekly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: dataMockMap},
			expectResponse: models.Response{Aggregation: "weekly", DataPoints: exampleAverageWeeklyDataPoints},
			expectError:    "",
		},
		{
			name:           "should return average monthly data points from the database",
			params:         ServiceParameters{TypeQuery: "df_average", Aggregation: "monthly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: dataMockMap},
			expectResponse: models.Response{Aggregation: "monthly", DataPoints: exampleAverageMonthlyDataPoints},
			expectError:    "",
		},
		{
			name:           "should return average quarterly data points from the database",
			params:         ServiceParameters{TypeQuery: "df_average", Aggregation: "quarterly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDataMap: dataMockMap},
			expectResponse: models.Response{Aggregation: "quarterly", DataPoints: exampleAverageQuarterlyDataPoints},
			expectError:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dfCountService := DfService{Client: tt.mockClient}
			got, err := dfCountService.ServeRequest(tt.params)

			if err == nil && tt.expectError != "" {
				t.Errorf("expected '%v' got no error", tt.expectError)
			}
			if err != nil && err.Error() != tt.expectError {
				t.Errorf("expected '%v' got '%v'", tt.expectError, err)
			}
			if !reflect.DeepEqual(got, tt.expectResponse) {
				t.Errorf("DfService.ServeRequest() = %v, want %v", got, tt.expectResponse)
			}
		})
	}
}
