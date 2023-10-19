package nds

import (
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type NDSFile struct {
	Path              string
	Name              string
	Developer         string
	Size              string
	Base64EncodedIcon string
}

func CleanText(text string) string {
	cleanedTitle := strings.Map(func(r rune) rune {
		if r == 0x0A || r == 0x0D {
			return ' '
		}
		return r
	}, text)

	return RemoveNonUnicodeChars(cleanedTitle)
}

func ISO2UTF8(text string) string {
	iso88591 := charmap.ISO8859_1.NewDecoder()
	utf8Bytes, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(text), iso88591))
	if err != nil {
		return ""
	}

	return string(utf8Bytes)
}

func RemoveNonUnicodeChars(input string) string {
	reg := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	return reg.ReplaceAllStringFunc(input, func(match string) string {
		if !unicode.IsLetter([]rune(match)[0]) && !unicode.IsNumber([]rune(match)[0]) && !unicode.IsSpace([]rune(match)[0]) {
			return ""
		}

		return match
	})
}
