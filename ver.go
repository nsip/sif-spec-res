package sifspecres

// versions :
var versions = []string{
	// "3.4.2",
	// "3.4.3",
	// "3.4.4",
	// "3.4.5",
	"3.4.6",
	"3.4.7",
	"3.4.8.draft",
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
