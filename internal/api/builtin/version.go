package builtin

var (
	Version VersionInfo = VersionInfo{
		Quickbeam: "0.1",
		API: "0.1",
		BBB: "0.1",
	}
)

type VersionInfo struct {
	Quickbeam string `json:"quickbeam"`
	API string `json:"api"`
	BBB string `json:"bigbluebutton"`
}

func GetVersion() VersionInfo {
	return Version
}
