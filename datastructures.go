package testi

type CollectionVersion struct {
	Name      string
	CommitMsg string
	Version   int
	Files     []string
}

type FileMetaData struct {
	Name    string
	Sha3sum string
}
