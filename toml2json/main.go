package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/embres"
	"github.com/cdutwhu/gotil/io"
	"github.com/cdutwhu/gotil/str"
	jt "github.com/cdutwhu/json-tool"
	sifspecres "github.com/nsip/sif-spec-res"
	"github.com/peterbourgon/mergemap"
)

var (
	fPf            = fmt.Printf
	fSp            = fmt.Sprint
	fSf            = fmt.Sprintf
	sCount         = strings.Count
	sReplaceAll    = strings.ReplaceAll
	sSplit         = strings.Split
	splitRev       = str.SplitRev
	mustWriteFile  = io.MustWriteFile
	failOnErr      = fn.FailOnErr
	createDirBytes = embres.CreateDirBytes
	printFileBytes = embres.PrintFileBytes
)

var (
	lsObjects        = []string{}
	mObjPaths        = map[string][]string{}
	mObjMaxLenOfPath = map[string]int{}

	clearBuf = func() {
		lsObjects = []string{}
		mObjPaths = map[string][]string{}
		mObjMaxLenOfPath = map[string]int{}
	}
)

// initGlobalMaps :
func initGlobalMaps(oneObjPathList interface{}, name, sep string) {
	// nameType := reflect.TypeOf(oneObjPathList).Name()
	value := reflect.ValueOf(oneObjPathList)
	nField := value.NumField()

	// for [****] version,
	// [nField] should be 1 as all paths have been wrapped into [****] Array
	for i := 0; i < nField; i++ {
		// [****] version
		lsPath := fSp(value.Field(i).Interface())
		lsPath = lsPath[1 : len(lsPath)-1]
		mObjPaths[name] = append(mObjPaths[name], sSplit(lsPath, " ")...)
		for _, path := range mObjPaths[name] {
			if n := sCount(path, sep) + 1; mObjMaxLenOfPath[name] < n {
				mObjMaxLenOfPath[name] = n
			}
		}
	}
	sort.SliceStable(mObjPaths[name], func(i, j int) bool {
		return sCount(mObjPaths[name][i], sep) < sCount(mObjPaths[name][j], sep)
	})
}

// InitCfgBuf :
func InitCfgBuf(cfg interface{}, sep string) {
	clearBuf()
	value := reflect.ValueOf(cfg)
	nField, valType := value.NumField(), value.Type()
	for i := 0; i < nField; i++ {
		fVal, fValTyp := value.Field(i), valType.Field(i)
		// nameType := reflect.TypeOf(fVal.Interface()).Name()
		// fPln(nameType)
		if fVal.Kind() == reflect.Struct {
			initGlobalMaps(fVal.Interface(), fValTyp.Name, sep)
			lsObjects = append(lsObjects, fValTyp.Name)
		}
	}
}

// GetLoadedObjects :
func GetLoadedObjects() []string {
	return append([]string{}, lsObjects...)
}

// GetAllFullPaths :
func GetAllFullPaths(obj, sep string) (paths []string) {
	for _, path := range mObjPaths[obj] {
		// fPln(path)
		paths = append(paths, obj+sep+path)
	}
	return
}

// GetLvlFullPaths :
func GetLvlFullPaths(obj, sep string, lvl int) (paths []string, valid bool) {
	if lvl > mObjMaxLenOfPath[obj] {
		return nil, false
	}
	for _, path := range mObjPaths[obj] {
		if lvl == sCount(path, sep)+1 {
			paths = append(paths, obj+sep+path)
		}
	}
	return paths, true
}

// -------------------------------------------------- //

// MakeBasicMap :
func MakeBasicMap(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{field: value}
}

// MakeOneMap :
func MakeOneMap(path, sep, valsymbol string) map[string]interface{} {
	var v interface{}
	for i, seg := range splitRev(path, sep) {
		if i == 0 {
			v = valsymbol
		}
		v = MakeBasicMap(seg, v)
	}
	return v.(map[string]interface{})
}

// MergeMaps :
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	var v map[string]interface{}
	for i, m := range maps {
		if i == 0 {
			v = m
		} else {
			v = mergemap.Merge(v, m)
		}
	}
	return v
}

// MakeMap :
func MakeMap(paths []string, sep, valsymbol string) map[string]interface{} {
	maps := []map[string]interface{}{}
	for _, path := range paths {
		maps = append(maps, MakeOneMap(path, sep, valsymbol))
	}
	return MergeMaps(maps...)
}

// MakeJSON :
func MakeJSON(m map[string]interface{}) string {
	jsonbytes, e := json.Marshal(m)
	failOnErr("MakeJSON Fatal: %v", e)
	return string(jsonbytes)
}

// ------------------------------------------------------------------------------- //

// YieldJSON4OneCfg :
func YieldJSON4OneCfg(obj, sep, outDir, jsonVal string, levelized, extContent bool) {
	if outDir[len(outDir)-1] != '/' {
		outDir += "/"
	}
	path := outDir + obj + "/"

	// delete all obsolete json files when new config-json files are coming
	failOnErr("%v", os.RemoveAll(path))
	fPf("%s is removed\n", path)
	failOnErr("%v", os.MkdirAll(path, 0700))
	fPf("%s is created\n", path)

	if levelized {
		for lvl := 1; lvl < 100; lvl++ {
			if paths, valid := GetLvlFullPaths(obj, sep, lvl); valid {
				mm := MakeMap(paths, sep, jsonVal)
				if len(mm) == 0 {
					continue
				}
				jsonstr := MakeJSON(mm)
				jsonfmt := jt.Fmt(jsonstr, "  ")
				mustWriteFile(fSf("%s%d.json", path, lvl), []byte(jsonfmt))
			} else {
				break
			}
		}
	} else {
		paths := GetAllFullPaths(obj, sep)
		mm := MakeMap(paths, sep, jsonVal)
		jsonstr := MakeJSON(mm)
		jsonfmt := jt.Fmt(jsonstr, "  ")
		mustWriteFile(fSf("%s0.json", path), []byte(jsonfmt))

		if extContent {
			// extend jsonstr, such as xml->json '#content', "30" => { "#content": "30" }
			jsonext := sReplaceAll(jsonstr, fSf(`"%s"`, jsonVal), fSf(`{"#content": "%s"}`, jsonVal))
			jsonextfmt := jt.Fmt(jsonext, "  ")
			mustWriteFile(fSf("%s1.json", path), []byte(jsonextfmt))
		}
	}
}

// YieldJSONBySIFList :
func YieldJSONBySIFList(cfgPath, ver string) {
	JSONCfgOutDir := "./" + ver + "/json/LIST/"
	switch ver {
	case "3.4.2":
		InitCfgBuf(*NewCfg("CfgL2J342", nil, cfgPath).(*CfgL2J342), "/") // Init Global Maps
	case "3.4.3":
		InitCfgBuf(*NewCfg("CfgL2J343", nil, cfgPath).(*CfgL2J343), "/")
	case "3.4.4":
		InitCfgBuf(*NewCfg("CfgL2J344", nil, cfgPath).(*CfgL2J344), "/")
	case "3.4.5":
		InitCfgBuf(*NewCfg("CfgL2J345", nil, cfgPath).(*CfgL2J345), "/")
	case "3.4.6":
		InitCfgBuf(*NewCfg("CfgL2J346", nil, cfgPath).(*CfgL2J346), "/")
	case "3.4.7":
		InitCfgBuf(*NewCfg("CfgL2J347", nil, cfgPath).(*CfgL2J347), "/")
	case "3.4.8":
		InitCfgBuf(*NewCfg("CfgL2J348", nil, cfgPath).(*CfgL2J348), "/")
	default:
		panic("unsupported version: " + ver)
	}

	for _, obj := range GetLoadedObjects() {
		YieldJSON4OneCfg(obj, "/", JSONCfgOutDir, "[]", true, false)
	}
}

// YieldJSONBySIFNum :
func YieldJSONBySIFNum(cfgPath, ver string) {
	JSONCfgOutDir := "./" + ver + "/json/NUMERIC/"
	switch ver {
	case "3.4.2":
		InitCfgBuf(*NewCfg("CfgN2J342", nil, cfgPath).(*CfgN2J342), "/") // Init Global Maps
	case "3.4.3":
		InitCfgBuf(*NewCfg("CfgN2J343", nil, cfgPath).(*CfgN2J343), "/")
	case "3.4.4":
		InitCfgBuf(*NewCfg("CfgN2J344", nil, cfgPath).(*CfgN2J344), "/")
	case "3.4.5":
		InitCfgBuf(*NewCfg("CfgN2J345", nil, cfgPath).(*CfgN2J345), "/")
	case "3.4.6":
		InitCfgBuf(*NewCfg("CfgN2J346", nil, cfgPath).(*CfgN2J346), "/")
	case "3.4.7":
		InitCfgBuf(*NewCfg("CfgN2J347", nil, cfgPath).(*CfgN2J347), "/")
	case "3.4.8":
		InitCfgBuf(*NewCfg("CfgN2J348", nil, cfgPath).(*CfgN2J348), "/")
	default:
		panic("unsupported version: " + ver)
	}
	for _, obj := range GetLoadedObjects() {
		YieldJSON4OneCfg(obj, "/", JSONCfgOutDir, "(N)", false, true)
	}
}

// YieldJSONBySIFBool :
func YieldJSONBySIFBool(cfgPath, ver string) {
	JSONCfgOutDir := "./" + ver + "/json/BOOLEAN/"
	switch ver {
	case "3.4.2":
		InitCfgBuf(*NewCfg("CfgB2J342", nil, cfgPath).(*CfgB2J342), "/") // Init Global Maps
	case "3.4.3":
		InitCfgBuf(*NewCfg("CfgB2J343", nil, cfgPath).(*CfgB2J343), "/")
	case "3.4.4":
		InitCfgBuf(*NewCfg("CfgB2J344", nil, cfgPath).(*CfgB2J344), "/")
	case "3.4.5":
		InitCfgBuf(*NewCfg("CfgB2J345", nil, cfgPath).(*CfgB2J345), "/")
	case "3.4.6":
		InitCfgBuf(*NewCfg("CfgB2J346", nil, cfgPath).(*CfgB2J346), "/")
	case "3.4.7":
		InitCfgBuf(*NewCfg("CfgB2J347", nil, cfgPath).(*CfgB2J347), "/")
	case "3.4.8":
		InitCfgBuf(*NewCfg("CfgB2J348", nil, cfgPath).(*CfgB2J348), "/")
	default:
		panic("unsupported version: " + ver)
	}
	for _, obj := range GetLoadedObjects() {
		YieldJSON4OneCfg(obj, "/", JSONCfgOutDir, "(B)", false, true)
	}
}

// YieldJSONBySIF :
func YieldJSONBySIF(listCfg, numCfg, boolCfg, ver string) {
	YieldJSONBySIFList(listCfg, ver)
	YieldJSONBySIFNum(numCfg, ver)
	YieldJSONBySIFBool(boolCfg, ver)
}

func main() {
	for _, ver := range sifspecres.GetAllVer("", "") {
		v := sReplaceAll(ver, ".", "")
		YieldJSONBySIF(
			"./"+ver+"/toml/List2JSON.toml",
			"./"+ver+"/toml/Num2JSON.toml",
			"./"+ver+"/toml/Bool2JSON.toml",
			ver,
		)
		pkg := "sif" + v
		printFileBytes(pkg, "TXT", "./"+ver+"/txt.go", false, "./"+ver+".txt")
		createDirBytes(pkg, "JSON_BOOL", "./"+ver+"/json/BOOLEAN/", "./"+ver+"/json_bool.go", false, v, "json", "BOOLEAN")
		createDirBytes(pkg, "JSON_LIST", "./"+ver+"/json/LIST/", "./"+ver+"/json_list.go", false, v, "json", "LIST")
		createDirBytes(pkg, "JSON_NUM", "./"+ver+"/json/NUMERIC/", "./"+ver+"/json_num.go", false, v, "json", "NUMERIC")
	}
}
