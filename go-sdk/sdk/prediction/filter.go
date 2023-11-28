package prediction

import "errors"

var (
	ErrDateRangeFilterRequired = errors.New("creation date range is required")
)

type Filter struct {
	RequestID    string
	Workflow     string
	WorkflowType string
	Process      string
	Version      string
	CreationDate TimestampRange
}

type TimestampRange struct {
	StartDate int64
	EndDate   int64
}

func (f Filter) Validate() error {
	if f.CreationDate.EndDate == 0 {
		return ErrDateRangeFilterRequired
	}

	return nil
}
