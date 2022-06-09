package tags

import (
	"fmt"
	"strconv"
)

func ParseTags(tag string) (map[string]string, error) {
	originalTag := tag

	tagsMap := map[string]string{}

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 {
			return nil, fmt.Errorf("malformed tag: missing tag name on tag string: '%s'", originalTag)
		}
		if i+1 >= len(tag) {
			return nil, fmt.Errorf("malformed tag: expected tag value but got empty string on tag string: '%s'", originalTag)
		}
		if tag[i] != ':' {
			return nil, fmt.Errorf("malformed tag: invalid character detected on tag string: '%d'", tag[i])
		}
		if tag[i+1] != '"' {
			return nil, fmt.Errorf("malformed tag: missing quotes right after tag name in tag string: '%s'", originalTag)
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			return nil, fmt.Errorf("malformed tag: missing end quote on a tag value in tag string: '%s'", originalTag)
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			continue
		}

		tagsMap[name] = value
	}

	return tagsMap, nil
}
