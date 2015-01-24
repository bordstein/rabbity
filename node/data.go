package main

type Repo struct {
	Name          string
	LatestVersion int64
}

type RepoVersion struct {
	Name      string
	CommitMsg string
	Files     map[string]FileInfo
	Version   int64
}

type FileInfo struct {
	Name string
}
