package main

import (
	"bytes"
	"errors"
	"os"
	"unicode"
)

type Tag struct {
	Attribute string
	Content   string
	Line      int
}

type Signature struct {
	Type       string
	Matches    []string
	Definition string
	Line       int
}

type Domment struct {
	Tags []Tag
	Sig  Signature
	Line int
}

/*!
 * @? Attempts to parse the signature of a domment
 */
func ParseSignature(signature string, line int) (Signature, error) {
	for _, v := range r.Expressions {
		if v.Expression.MatchString(signature) {
			matches := v.Expression.FindAllStringSubmatch(signature, -1)[0][1:]

			return Signature{Type: v.Name, Matches: matches, Definition: signature, Line: line}, nil
		}
	}

	return Signature{}, errors.New("Failed to identify domment signature")
}

func ParseDomment(contents []byte, startPtr int, startLine int) (Domment, int, int, error) {
	var dmnt Domment
	dmnt.Line = startLine
	ptr := startPtr
	currentLine := startLine

	for ptr < len(contents) {
		// End of block detection
		if bytes.HasPrefix(contents[ptr:], []byte("*/")) {
			ptr += 2 // Move past "*/"

			// Count trailing newlines to look for signature
			newlineCount := 0
			for ptr < len(contents) && unicode.IsSpace(rune(contents[ptr])) {
				if contents[ptr] == '\n' {
					newlineCount++
					currentLine++
				}
				ptr++
			}

			// Signature must be on the immediate next line (exactly 1 newline)
			if newlineCount == 1 && ptr < len(contents) {
				ptrStart := ptr
				sigLine := currentLine
				for ptr < len(contents) && contents[ptr] != '\n' {
					ptr++
				}
				signature := string(contents[ptrStart:ptr])

				// Optional: currentLine++ if we consumed the signature's newline
				if ptr < len(contents) && contents[ptr] == '\n' {
					currentLine++
					ptr++
				}

				sig, _ := ParseSignature(signature, sigLine)
				dmnt.Sig = sig
				return dmnt, ptr, currentLine, nil
			}

			return dmnt, ptr, currentLine, nil
		}

		// Tag detection
		if contents[ptr] == '@' {
			tagLine := currentLine
			ptr++ // consume @

			// read attribute name
			ptrStart := ptr
			for ptr < len(contents) && !unicode.IsSpace(rune(contents[ptr])) {
				ptr++
			}
			tag := string(contents[ptrStart:ptr])

			// read content
			for ptr < len(contents) && (contents[ptr] == ' ' || contents[ptr] == '\t') {
				ptr++
			}

			ptrStart = ptr
			ptrEnd := ptr
			for ptr < len(contents) && contents[ptr] != '@' && !bytes.HasPrefix(contents[ptr:], []byte("*/")) {
				if contents[ptr] == '\n' {
					currentLine++
				}
				if !unicode.IsSpace(rune(contents[ptr])) && contents[ptr] != '*' {
					ptrEnd = ptr
				}
				ptr++
			}

			value := contents[ptrStart : ptrEnd+1]
			var valueCleaned []byte

			for i := 0; i < len(value); i++ {
				if value[i] == '\n' {
					valueCleaned = append(valueCleaned, '\n')
					i++
					for i < len(value) && (value[i] == ' ' || value[i] == '\t' || value[i] == '*') {
						i++
					}
					i--
				} else {
					valueCleaned = append(valueCleaned, value[i])
				}
			}

			dmnt.Tags = append(dmnt.Tags, Tag{
				Attribute: tag,
				Content:   string(bytes.TrimSpace(valueCleaned)),
				Line:      tagLine,
			})
			continue
		}

		if contents[ptr] == '\n' {
			currentLine++
		}
		ptr++
	}

	return dmnt, ptr, currentLine, nil
}

func Document(path string, file_name string) []Domment {
	contents, err := os.ReadFile(path + "/" + file_name)
	if err != nil {
		return nil
	}

	var dmntSlc []Domment
	bytePtr := 0
	currentLine := 1

	for bytePtr < len(contents) {
		// Match the start of a block
		if bytes.HasPrefix(contents[bytePtr:], []byte("/*!")) {
			startLine := currentLine
			bytePtr += 3 // Move past "/*!"

			dmnt, endPtr, lastLine, err := ParseDomment(contents, bytePtr, startLine)
			if err == nil {
				dmntSlc = append(dmntSlc, dmnt)
			}

			bytePtr = endPtr
			currentLine = lastLine
			continue
		}

		if contents[bytePtr] == '\n' {
			currentLine++
		}
		bytePtr++
	}

	return dmntSlc
}
