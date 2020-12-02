package sifspecres

// versions :
var versions = []string{
	"3.4.6",
	"3.4.7",
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
