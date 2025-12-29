/*!
 * @p Domment
 * @v 0.0.1β
 */

package main

import (
	"os"
)

type ParsedFile struct {
	path     string
	domments []Domment
}

func main() {
	InitialiseRegex()
	InitialiseTemplates()

	path := "./"
	out := "../docs/"

	entries, err := os.ReadDir(path)

	check(err)

	var files []ParsedFile
	for _, v := range entries {
		if v.IsDir() {
			continue
		}

		files = append(files, ParsedFile{
			path:     v.Name(),
			domments: Document(path, v.Name()),
		})
	}

	Generate(files, out)
}
