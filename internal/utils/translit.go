package utils

import (
	"bytes"
	"strings"
	"unicode"
)

var rusASCII = map[string]string{
	"а": "a",
	"б": "b",
	"в": "v",
	"г": "g",
	"д": "d",
	"е": "e",
	"ё": "yo",
	"ж": "zh",
	"з": "z",
	"и": "i",
	"й": "j",
	"к": "k",
	"л": "l",
	"м": "m",
	"н": "n",
	"о": "o",
	"п": "p",
	"р": "r",
	"с": "s",
	"т": "t",
	"у": "u",
	"ф": "f",
	"х": "h",
	"ц": "c",
	"ч": "ch",
	"ш": "sh",
	"щ": "sch",
	"ъ": "'",
	"ы": "y",
	"ь": "",
	"э": "e",
	"ю": "ju",
	"я": "ja",
}

// func seo(s string, collection map[int]string) (string, error) {
// 	runes := []rune(s)
// 	result := make([]string, 0, len(s))

// 	for symbol, i := range runes {
// 		if val, ok := collection[unicode.ToLower(i)]; ok {
// 			fmt.Println(i)
// 		} else {
// 		}
// 	}

// 	return strings.Join(result, ""), nil
// }

func rusToLatin(text string, table map[string]string, proc specProc) string {
	if text == "" {
		return ""
	}

	var input = bytes.NewBufferString(text)
	var output = bytes.NewBuffer(nil)

	// Previous, next letter for special processor
	var p, n rune
	var rr string
	var ok bool

	for {
		r, _, err := input.ReadRune()

		if err != nil {
			break
		}

		lowerRune := unicode.ToLower(r)

		if !isRussianChar(lowerRune) {
			output.WriteRune(lowerRune)
			p = lowerRune
			continue
		}

		if proc != nil {
			n, _, _ = input.ReadRune()

			err = input.UnreadRune()
			if err != nil {
				break
			}

			rr, ok = proc(p, lowerRune, n, table)

			if ok {
				output.WriteString(rr)
				continue
			}
		}

		p = lowerRune

		rr, ok = rusASCII[string(lowerRune)]

		if ok {
			output.WriteString(rr)
			continue
		}

		rr, ok = table[string(lowerRune)]

		if ok {
			output.WriteString(rr)
		}
	}

	return output.String()
}

type specProc func(p, c, n rune, table map[string]string) (string, bool)

func isRussianChar(r rune) bool {
	switch {
	case r >= 1040 && r <= 1103,
		r == 1105, r == 1025:
		return true
	}

	return false
}

func EncodeRus(text string) string {
	encodeText := rusToLatin(text, rusASCII, nil)
	encodeText = strings.ReplaceAll(encodeText, " ", "-")
	return strings.ToLower(strings.Trim(encodeText, "-"))
}
