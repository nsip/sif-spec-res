package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cdutwhu/gonfig/strugen"
	"github.com/digisan/gotk/io"
	"github.com/digisan/gotk/slice/tsb"
	"github.com/digisan/logkit"
	"github.com/nsip/sif-spec-res/tool"
)

var (
	fPln          = fmt.Println
	fPf           = fmt.Printf
	fSf           = fmt.Sprintf
	sHasPrefix    = strings.HasPrefix
	sSplit        = strings.Split
	sReplaceAll   = strings.ReplaceAll
	sTrim         = strings.Trim
	sTrimPrefix   = strings.TrimPrefix
	sTrimSuffix   = strings.TrimSuffix
	sJoin         = strings.Join
	mustWriteFile = io.MustWriteFile
	failOnErr     = logkit.FailOnErr

	rmHeadToFirst = func(s, mark string) string {
		if segs := sSplit(s, mark); len(segs) > 1 {
			return sJoin(segs[1:], mark)
		}
		return s
	}
)

// Println :
func Println(num bool, slc ...string) {
	if num {
		for i, v := range slc {
			fPf("%d: %v\n", i, v)
		}
	} else {
		for _, v := range slc {
			fPln(v)
		}
	}
}

// ObjGrp :
func ObjGrp(sep string, listGrp ...string) []string {
	m := map[string]bool{}
	for _, lsPath := range listGrp {
		obj := sSplit(lsPath, sep)[0]
		if _, ok := m[obj]; !ok {
			m[obj] = true
		}
	}
	keys, _ := tsb.Map2KVs(m, nil, nil)
	return keys
}

// MapOfGrp :
func MapOfGrp(objs []string, sep string, xxxPathGrp ...string) map[string][]string {
	m := make(map[string][]string)
	for _, obj := range objs {
		prefix := obj + sep
		for _, lp := range xxxPathGrp {
			if sHasPrefix(lp, prefix) {
				lp = sSplit(lp, "\t")[0]
				lp = sSplit(lp, " ")[0]
				m[obj] = append(m[obj], rmHeadToFirst(lp, sep))
			} else {
				m[obj] = append(m[obj], obj)
				break
			}
		}
	}
	return m
}

// PrintGrp4Cfg :
func PrintGrp4Cfg(m map[string][]string, attr string) (toml string) {
	switch attr {
	case "LIST", "NUMERIC", "BOOLEAN", "ATTRIBUTE", "OBJECT":
		for obj, grp := range m {
			content := fSf("[%s]\n  %s = [", obj, attr)
			for _, path := range grp {
				content += fSf("\"%s\", ", path)
			}
			toml += content[:len(content)-2] + "]" + "\n\n"
		}
	}
	return
}

// GenTomlAndGoSrc :
func GenTomlAndGoSrc(specPath, outDir string) {

	const (
		SEP       = "/"
		VERSION   = "VERSION: "
		OBJECT    = "OBJECT: "
		LIST      = "LIST: "
		NUMERIC   = "NUMERIC: "
		BOOLEAN   = "BOOLEAN: "
		ATTRIBUTE = "SIMPLE ATTRIBUTE: "
	)

	var (
		objGrp      []string
		listPathGrp []string
		numPathGrp  []string
		boolPathGrp []string
		attrPathGrp []string
		SIFVer      string
	)

	bytes, err := ioutil.ReadFile(specPath)
	failOnErr("%v", err)

	for _, line := range sSplit(string(bytes), "\n") {
		switch {
		case sHasPrefix(line, VERSION):
			SIFVer = sTrim(line[len(VERSION):], " \t\r\n")
		case sHasPrefix(line, OBJECT):
			objGrp = append(objGrp, sTrim(line[len(OBJECT):], " \t\r\n"))
		case sHasPrefix(line, LIST):
			// listPathGrp = append(listPathGrp, rmTailFromLast(line[len(LIST):], "/")) // exclude last one
			listPathGrp = append(listPathGrp, sTrim(line[len(LIST):], " \t\r\n"))
		case sHasPrefix(line, NUMERIC):
			numPathGrp = append(numPathGrp, sTrim(line[len(NUMERIC):], " \t\r\n"))
		case sHasPrefix(line, BOOLEAN):
			boolPathGrp = append(boolPathGrp, sTrim(line[len(BOOLEAN):], " \t\r\n"))
		case sHasPrefix(line, ATTRIBUTE):
			attrPathGrp = append(attrPathGrp, sTrim(line[len(ATTRIBUTE):], " \t\r\n"))
		}
	}

	// Println(true, objGrp...)
	// fPln("-----------------------------")

	// Println(false, listPathGrp...)
	// fPln("-----------------------------")

	mListAttr := MapOfGrp(ObjGrp(SEP, listPathGrp...), SEP, listPathGrp...)
	mNumAttr := MapOfGrp(ObjGrp(SEP, numPathGrp...), SEP, numPathGrp...)
	mBoolAttr := MapOfGrp(ObjGrp(SEP, boolPathGrp...), SEP, boolPathGrp...)
	mObjAttr := MapOfGrp(ObjGrp(SEP, objGrp...), SEP, objGrp...)
	mAttr2 := MapOfGrp(ObjGrp(SEP, attrPathGrp...), SEP, attrPathGrp...)

	verln := fSf("Version = \"%s\"\n\n", SIFVer)
	toml4List := verln + PrintGrp4Cfg(mListAttr, "LIST")
	toml4Num := verln + PrintGrp4Cfg(mNumAttr, "NUMERIC")
	toml4Bool := verln + PrintGrp4Cfg(mBoolAttr, "BOOLEAN")
	toml4Obj := verln + PrintGrp4Cfg(mObjAttr, "OBJECT")
	toml4Attr := verln + PrintGrp4Cfg(mAttr2, "ATTRIBUTE")

	mustWriteFile(outDir+"toml/List2JSON.toml", []byte(toml4List))
	mustWriteFile(outDir+"toml/Num2JSON.toml", []byte(toml4Num))
	mustWriteFile(outDir+"toml/Bool2JSON.toml", []byte(toml4Bool))
	mustWriteFile(outDir+"toml/Obj2JSON.toml", []byte(toml4Obj))
	mustWriteFile(outDir+"toml/Attr2JSON.toml", []byte(toml4Attr))
}

func main() {

	cfgSrc, pkgName := "./toml2json/config.go", "main"
	os.Remove(cfgSrc)

	for _, spec := range tool.GetAllVer("./", ".txt") {
		ver := sTrimPrefix(spec, "./")
		ver = sTrimSuffix(ver, ".txt")
		outdir := "./" + ver + "/"
		GenTomlAndGoSrc(spec, outdir)
		tomlPath := outdir + "toml/"
		v := sReplaceAll(ver, ".", "")

		CfgL2J, CfgB2J, CfgN2J, CfgO2J, CfgA2J := "CfgL2J"+v, "CfgB2J"+v, "CfgN2J"+v, "CfgO2J"+v, "CfgA2J"+v
		strugen.GenStruct(tomlPath+"List2JSON.toml", CfgL2J, pkgName, cfgSrc)
		strugen.GenStruct(tomlPath+"Bool2JSON.toml", CfgB2J, pkgName, cfgSrc)
		strugen.GenStruct(tomlPath+"Num2JSON.toml", CfgN2J, pkgName, cfgSrc)
		strugen.GenStruct(tomlPath+"Obj2JSON.toml", CfgO2J, pkgName, cfgSrc)
		strugen.GenStruct(tomlPath+"Attr2JSON.toml", CfgA2J, pkgName, cfgSrc)
	}

	strugen.GenNewCfg(cfgSrc)
}
