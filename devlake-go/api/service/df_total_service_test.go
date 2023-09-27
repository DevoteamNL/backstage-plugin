package service

import (
	"devlake-go/group-sync/api/models"
	"devlake-go/group-sync/api/sql_client"
	"fmt"
	"reflect"
	"testing"
)

func TestDfTotalService_ServeRequest(t *testing.T) {

	exampleWeeklyDataPoints := []models.DataPoint{{Key: "202338", Value: 0}, {Key: "202337", Value: 1}, {Key: "202336", Value: 2}}
	exampleMonthlyDataPoints := []models.DataPoint{{Key: "23/04", Value: 6}, {Key: "23/03", Value: 5}, {Key: "23/02", Value: 4}}

	tests := []struct {
		name           string
		params         ServiceParameters
		mockClient     sql_client.MockClient
		expectResponse models.Response
		expectError    string
	}{
		{
			name:   "should return an error with an unexpected error from the database",
			params: ServiceParameters{TypeQuery: "df_total", Aggregation: "weekly", Project: "", To: 0, From: 0},
			mockClient: sql_client.MockClient{
				WeeklyErrorToReturn:  fmt.Errorf("error from weekly query"),
				MonthlyErrorToReturn: fmt.Errorf("error from monthly query"),
			},
			expectResponse: models.Response{Aggregation: "weekly", DataPoints: nil},
			expectError:    "error from weekly query",
		},
		{
			name:   "should return an error with an unexpected error from the database",
			params: ServiceParameters{TypeQuery: "df_total", Aggregation: "monthly", Project: "", To: 0, From: 0},
			mockClient: sql_client.MockClient{
				WeeklyErrorToReturn:  fmt.Errorf("error from monthly query"),
				MonthlyErrorToReturn: fmt.Errorf("error from monthly query"),
			},
			expectResponse: models.Response{Aggregation: "monthly", DataPoints: nil},
			expectError:    "error from monthly query",
		},
		{
			name:   "should return weekly data points from the database",
			params: ServiceParameters{TypeQuery: "df_total", Aggregation: "weekly", Project: "", To: 0, From: 0},
			mockClient: sql_client.MockClient{
				WeeklyDataPointsToReturn:  exampleWeeklyDataPoints,
				MonthlyDataPointsToReturn: exampleMonthlyDataPoints,
			},
			expectResponse: models.Response{Aggregation: "weekly", DataPoints: exampleWeeklyDataPoints},
			expectError:    "",
		},
		{
			name:   "should return monthly data points from the database",
			params: ServiceParameters{TypeQuery: "df_total", Aggregation: "monthly", Project: "", To: 0, From: 0},
			mockClient: sql_client.MockClient{
				WeeklyDataPointsToReturn:  exampleWeeklyDataPoints,
				MonthlyDataPointsToReturn: exampleMonthlyDataPoints,
			},
			expectResponse: models.Response{Aggregation: "monthly", DataPoints: exampleMonthlyDataPoints},
			expectError:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dfTotalService := DfTotalService{Client: tt.mockClient}
			got, err := dfTotalService.ServeRequest(tt.params)

			if err == nil && tt.expectError != "" {
				t.Errorf("expected '%v' got no error", tt.expectError)
			}
			if err != nil && err.Error() != tt.expectError {
				t.Errorf("expected '%v' got '%v'", tt.expectError, err)
			}
			if !reflect.DeepEqual(got, tt.expectResponse) {
				t.Errorf("DfTotalService.ServeRequest() = %v, want %v", got, tt.expectResponse)
			}
		})
	}
}
