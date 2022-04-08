package repositories

// Authorization is the structure that represents the detail of the repository authorization.
type Authorization struct {
	// Specified the authorization name.
	Name string `json:"name"`
	// Specified the repository type
	// The valid values are: devcloud, github, gitlab, gitee and bitbucket。
	RepoType string `json:"repo_type"`
	// Specified the repository address.
	RepoHost string `json:"repo_host"`
	// Specified the repository homepage.
	RepoHome string `json:"repo_home"`
	// Specified the repository username.
	RepoUser string `json:"repo_user"`
	// Specified the repository avatar.
	Avartar string `json:"avartar"`
	// Specified the authorization mode.
	TokenType string `json:"token_type"`
	// Specified the creation time.
	CreateTime int `json:"create_time"`
	// Specified the update time.
	UpdateTime int `json:"update_time"`
	// Specified the authorization status.
	Status int `json:"status"`
}
