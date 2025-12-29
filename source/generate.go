package main

import (
	"os"
	"path"
	"time"
)

var Data struct {
	ProjectName string
	Version     string
	Date        string
	Files       []string
}

func MakeHome(out string) {
	f, err := os.Create(path.Join(out, "index.html"))
	check(err)
	defer f.Close()

	err = Templates.Index.Execute(f, Data)
	check(err)
}

func Generate(files []ParsedFile, out string) {
	err := os.Mkdir(out, 0755)
	if !os.IsExist(err) {
		check(err)
	}

	Data.Date = time.Now().Format("2006-01-02")

	for _, v := range files {
		// Individual files

		Data.Files = append(Data.Files, v.path)

		for _, w := range v.domments {
			// Per file domments

			for _, d := range w.Tags {
				switch d.Attribute {
				case "project", "prj", "p":
					Data.ProjectName = d.Content
				case "version", "ver", "v":
					Data.Version = d.Content
				}
			}
		}
	}

	MakeHome(out)
}
