package nod

import "time"

// TimeFilter specifies a time range for filtering queries.
type TimeFilter struct {
	From *time.Time
	To   *time.Time
}

// StringFilter specifies string matching criteria for filtering queries.
type StringFilter struct {
	Equals     *string
	Contains   *string
	StartsWith *string
	EndsWith   *string
}

// StringEquals creates a StringFilter matching an exact value.
func StringEquals(value string) *StringFilter {
	return &StringFilter{Equals: &value}
}

// StringContains creates a StringFilter matching a substring.
func StringContains(value string) *StringFilter {
	return &StringFilter{Contains: &value}
}

// StringStartsWith creates a StringFilter matching a prefix.
func StringStartsWith(value string) *StringFilter {
	return &StringFilter{StartsWith: &value}
}

// StringEndsWith creates a StringFilter matching a suffix.
func StringEndsWith(value string) *StringFilter {
	return &StringFilter{EndsWith: &value}
}

// TimeFrom creates a TimeFilter with a starting time.
func TimeFrom(value time.Time) *TimeFilter {
	return &TimeFilter{From: &value}
}

// TimeTo creates a TimeFilter with an ending time.
func TimeTo(value time.Time) *TimeFilter {
	return &TimeFilter{To: &value}
}

// TimeBetween creates a TimeFilter with a range between two times.
func TimeBetween(from, to time.Time) *TimeFilter {
	return &TimeFilter{From: &from, To: &to}
}

// NewStringFilter creates a new StringFilter with the specified criteria.
func NewStringFilter(equals, contains, startsWith, endsWith *string) *StringFilter {
	return &StringFilter{
		Equals:     equals,
		Contains:   contains,
		StartsWith: startsWith,
		EndsWith:   endsWith,
	}
}

// NewTimeFilter creates a new TimeFilter with the specified range.
func NewTimeFilter(from, to *time.Time) *TimeFilter {
	return &TimeFilter{
		From: from,
		To:   to,
	}
}
