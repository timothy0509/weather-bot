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
func (s StringSet) T(lang Language) string {
	switch lang {
	case EN:
		return s.EN
	case TC:
		return s.TC
	case SC:
		return s.SC
	case Bilingual:
		return s.EN + " / " + s.TC
	default:
		return s.EN
	}
}

// Format bilingual returns a bilingual label/value string.
func Format(label, valueEn, valueTC string, lang Language) string {
	switch lang {
	case EN:
		return label + ": " + valueEn
	case TC:
		return label + "：" + valueTC
	case SC:
		return label + "：" + valueTC
	case Bilingual:
		return label + " / " + label + "：" + valueEn + " / " + valueTC
	default:
		return label + ": " + valueEn
	}
}
