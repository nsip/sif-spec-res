package tool

import (
	"log"
	"regexp"
	"strings"

	v342 "github.com/nsip/sif-spec-res/3.4.2"
	v343 "github.com/nsip/sif-spec-res/3.4.3"
	v344 "github.com/nsip/sif-spec-res/3.4.4"
	v345 "github.com/nsip/sif-spec-res/3.4.5"
	v346 "github.com/nsip/sif-spec-res/3.4.6"
	v347 "github.com/nsip/sif-spec-res/3.4.7"
	v348 "github.com/nsip/sif-spec-res/3.4.8"
	v349 "github.com/nsip/sif-spec-res/3.4.9"

	goio "github.com/digisan/gotk/io"
)

// versions :
var versions = []string{
	"3.4.2",
	"3.4.3",
	"3.4.4",
	"3.4.5",
	"3.4.6",
	"3.4.7",
	"3.4.8",
	"3.4.9",
}

// GetAllVer :
func GetAllVer(prefix, suffix string) (vers []string) {
	if prefix == "" && suffix == "" {
		return versions
	}
	for _, v := range versions {
		vers = append(vers, prefix+v+suffix)
	}
	return
}

func GetAttrPaths(ver string) map[string]struct{} {

	mAttrPaths := make(map[string]struct{})

	var txt string
	switch ver {
	case "3.4.2":
		txt = string(v342.TXT["342"])
	case "3.4.3":
		txt = string(v343.TXT["343"])
	case "3.4.4":
		txt = string(v344.TXT["344"])
	case "3.4.5":
		txt = string(v345.TXT["345"])
	case "3.4.6":
		txt = string(v346.TXT["346"])
	case "3.4.7":
		txt = string(v347.TXT["347"])
	case "3.4.8":
		txt = string(v348.TXT["348"])
	case "3.4.9":
		txt = string(v349.TXT["349"])
	}

	r := regexp.MustCompile(`\s+([\w\d]+/)+@[\w\d]+\s+`)
	if _, err := goio.StrLineScan(txt, func(line string) (bool, string) {
		ss := r.FindAllString(line, -1)
		if len(ss) > 0 {
			path := strings.Trim(ss[0], " \t")
			mAttrPaths[path] = struct{}{}
		}
		return false, ""
	}, ""); err != nil {
		log.Fatalln(err)
	}

	return mAttrPaths
}
