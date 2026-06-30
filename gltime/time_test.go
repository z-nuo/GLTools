package gltime

import (
	"testing"
	"time"
)

func TestFormatAndParse(t *testing.T) {
	base := time.Date(2026, 6, 30, 15, 4, 5, 123000000, time.Local)
	layout := "2006-01-02 15:04:05"

	formatted := Format(base, layout)
	if formatted != "2026-06-30 15:04:05" {
		t.Fatalf("Format() = %q, want %q", formatted, "2026-06-30 15:04:05")
	}

	parsed, err := Parse(layout, formatted)
	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if parsed.Format(layout) != formatted {
		t.Fatalf("Parse() = %q, want %q", parsed.Format(layout), formatted)
	}
}

func TestParseInvalidValue(t *testing.T) {
	if _, err := Parse("2006-01-02", "bad"); err == nil {
		t.Fatal("Parse() error = nil, want error")
	}
}

func TestUnixMilliRoundTrip(t *testing.T) {
	base := time.Date(2026, 6, 30, 15, 4, 5, 123000000, time.Local)

	ms := UnixMilli(base)
	got := FromUnixMilli(ms)
	if !got.Equal(base) {
		t.Fatalf("FromUnixMilli(UnixMilli()) = %v, want instant %v", got, base)
	}
}

func TestDayStart(t *testing.T) {
	base := time.Date(2026, 6, 30, 15, 4, 5, 123000000, time.Local)
	want := time.Date(2026, 6, 30, 0, 0, 0, 0, time.Local)

	if got := DayStart(base); !got.Equal(want) {
		t.Fatalf("DayStart() = %v, want %v", got, want)
	}
}

func TestDayEnd(t *testing.T) {
	base := time.Date(2026, 6, 30, 15, 4, 5, 123000000, time.Local)
	want := time.Date(2026, 6, 30, 23, 59, 59, int(time.Second-time.Nanosecond), time.Local)

	if got := DayEnd(base); !got.Equal(want) {
		t.Fatalf("DayEnd() = %v, want %v", got, want)
	}
}
