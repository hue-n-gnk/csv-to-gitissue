package pointer

import "time"

// String returns pointer of given string.
func String(s string) *string {
	return &s
}

// SafeString ensures to return empty string even pointer is nil.
func SafeString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// Bool returns pointer of given bool.
func Bool(b bool) *bool {
	return &b
}

// SafeBool ensures to return false even pointer is nil.
func SafeBool(b *bool) bool {
	if b == nil {
		return false
	}

	return *b
}

// Int returns pointer of given int.
func Int(s int) *int {
	return &s
}

// SafeInt ensures to return default int value even pointer is nil.
func SafeInt(s *int) int {
	if s == nil {
		return 0
	}

	return *s
}

// Float32 returns pointer of given float32.
func Float32(f float32) *float32 {
	return &f
}

// SafeFloat32 ensures to return default float32 value even pointer is nil.
func SafeFloat32(f *float32) float32 {
	if f == nil {
		return 0
	}

	return *f
}

func SafeTime(f *time.Time) time.Time {
	if f == nil {
		return time.Time{}
	}

	return *f
}
