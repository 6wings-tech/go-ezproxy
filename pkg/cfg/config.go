package cfg

type config struct {
	GitDirs        []string `json:"gitdirs"`
	Domain         string   `json:"domain"`
	CloneUrlPrefix string   `json:"cloneUrlPrefix"`
}
