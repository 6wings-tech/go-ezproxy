package cfg

type config struct {
	ReposRoot      string `json:"reposRoot"`
	Domain         string `json:"domain"`
	SshUser        string `json:"sshUser"`
	SshPort        int    `json:"sshPort"`
	CloneUrlPrefix string `json:"cloneUrlPrefix"`
}
