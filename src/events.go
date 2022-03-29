package src

type PullRequest struct {
	Action          string `json:"action"`
	Number          int32  `json:"number"`
	PullRequestData PRData `json:"pull_request"`
	Commits         int64  `json:"commits"`
	Additions       int64  `json:"additions"`
	Deletions       int64  `json:"deletions"`
	ChangedFiles    int64  `json:"changed_files"`
}
type PRData struct {
	Url            string `json:"url"`
	ID             uint64 `json:"id"`
	HtmlUrl        string `json:"html_url"`
	DiffUrl        string `json:"diff_url"`
	PatchUrl       string `json:"patch_url"`
	IssueUrl       string `json:"issue_url"`
	State          string `json:"state"`
	Title          string `json:"title"`
	User           User   `json:"user"`
	CommitUrl      string `json:"commit_url"`
	MergeCommitSha string `json:"merge_commit_sha"`
	Head           Head   `json:"head"`
	Base           Base   `json:"base"`
}
type User struct {
	Login   string `json:"login"`
	ID      uint64 `json:"id"`
	Url     string `json:"url"`
	HtmlUrl string `json:"html_url"`
	Type    string `json:"type"`
}
type Head struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	Sha   string `json:"sha"`
	User  User   `json:"user"`
	Repo  Repo   `json:"repo"`
}
type Base struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	Sha   string `json:"sha"`
	User  User   `json:"user"`
	Repo  Repo   `json:"repo"`
}

type Repo struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	Owner       User   `json:"owner"`
	HtmlUrl     string `json:"html_url"`
	Description string `json:"descirption"`
	Url         string `json:"url"`
	GitUrl      string `json:"git_url"`
	SshUrl      string `json:"ssh_url"`
	CloneUrl    string `json:"clone_url"`
	Visibility  string `json:"visibility"`
}
