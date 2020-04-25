package ghost

import "time"

// String returns a pointer to the string.
func String(s string) *string {
	return &s
}

// Bool returns a pointer to the bool.
func Bool(b bool) *bool {
	return &b
}

// Int returns a pointer to the int.
func Int(i int) *int {
	return &i
}

// Time creates a timestamp from the RFC3339 string and returns a pointer,
// ignoring any errors that occur during construction.
func Time(s string) *time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return &t
}
