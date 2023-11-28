package prediction_test

import (
	"testing"

	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/prediction"
	"github.com/stretchr/testify/assert"
)

func TestFilter_Validate(t *testing.T) {
	testCases := []struct {
		name          string
		filter        prediction.Filter
		expectedError error
	}{
		{
			name: "valid filter",
			filter: prediction.Filter{
				CreationDate: prediction.TimestampRange{
					StartDate: 100000,
					EndDate:   200000,
				},
			},
			expectedError: nil,
		},
		{
			name: "filter without start date is valid",
			filter: prediction.Filter{
				CreationDate: prediction.TimestampRange{
					EndDate: 200000,
				},
			},
			expectedError: nil,
		},
		{
			name: "filter without creation date range is not valid",
			filter: prediction.Filter{
				CreationDate: prediction.TimestampRange{},
			},
			expectedError: prediction.ErrDateRangeFilterRequired,
		},
		{
			name: "filter without end date is not valid",
			filter: prediction.Filter{
				CreationDate: prediction.TimestampRange{
					StartDate: 1000,
				},
			},
			expectedError: prediction.ErrDateRangeFilterRequired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, tc.filter.Validate(), tc.expectedError)
		})
	}
}
