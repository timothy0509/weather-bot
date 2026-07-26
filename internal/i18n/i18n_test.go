package i18n

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		input    string
		expected Language
	}{
		{"en", EN},
		{"English", EN},
		{"tc", TC},
		{"traditional", TC},
		{"sc", SC},
		{"simplified", SC},
		{"bilingual", Bilingual},
		{"bi", Bilingual},
		{"unknown", EN},
		{"", EN},
	}
	for _, c := range cases {
		got := Normalize(c.input)
		if got != c.expected {
			t.Errorf("Normalize(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestT(t *testing.T) {
	if got, want := T("temperature", EN), "Temperature"; got != want {
		t.Errorf("T(temperature, EN) = %q, want %q", got, want)
	}
	if got, want := T("temperature", TC), "氣溫"; got != want {
		t.Errorf("T(temperature, TC) = %q, want %q", got, want)
	}
	if got, want := T("temperature", Bilingual), "Temperature / 氣溫"; got != want {
		t.Errorf("T(temperature, Bilingual) = %q, want %q", got, want)
	}
}

func TestIsValid(t *testing.T) {
	if !IsValid("en") {
		t.Error("expected en to be valid")
	}
	if IsValid("fr") {
		t.Error("expected fr to be invalid")
	}
}
