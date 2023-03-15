package blok_cli

type BlokConfig struct {
	Src         string           `json:"src"`
	Dst         string           `json:"dst"`
	Site        string           `json:"site"`
	DstProvider string           `json:"dst_provider"`
	Build       *BlokBuildConfig `json:"build"`
}

type BlokBuildConfig struct {
	BuildPath string `json:"build_path"`
	Output    string `json:"output"`
	Command   string `json:"command"`
}
