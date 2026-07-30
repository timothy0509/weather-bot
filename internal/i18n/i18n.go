// Package i18n provides multi-language user-facing strings.
package i18n

import "strings"

// Language is a supported language mode.
type Language string

const (
	EN         Language = "en"
	TC         Language = "tc"
	SC         Language = "sc"
	Bilingual  Language = "bilingual"
)

// IsValid reports whether a language code is supported.
func IsValid(lang string) bool {
	switch Language(lang) {
	case EN, TC, SC, Bilingual:
		return true
	default:
		return false
	}
}

// Normalize normalizes a language string.
func Normalize(lang string) Language {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "en", "english":
		return EN
	case "tc", "traditional", "trad", "cht":
		return TC
	case "sc", "simplified", "simp", "chs":
		return SC
	case "bilingual", "bi", "both":
		return Bilingual
	default:
		return EN
	}
}

// StringSet holds translations for the four language modes.
type StringSet struct {
	EN, TC, SC string
}

// T returns the translation for the given language.
// For bilingual mode, English and Traditional Chinese are shown side-by-side.
// Falls back to English if the requested language has an empty translation.
func (s StringSet) T(lang Language) string {
	var result string
	switch lang {
	case EN:
		result = s.EN
	case TC:
		result = s.TC
	case SC:
		result = s.SC
	case Bilingual:
		en := s.EN
		tc := s.TC
		if tc == "" {
			tc = en
		}
		return en + " / " + tc
	default:
		result = s.EN
	}
	if result == "" {
		return s.EN
	}
	return result
}
