package thesaurus

type ThesaurusStatusResp struct {
	// Loading status
	// Loaded indicates that a custom word dictionary is successfully loaded.
	// Loading indicates that a custom word dictionary is being loaded.
	// Failed indicates that a custom word dictionary fails to be loaded.
	Status        string `json:"status"`
	Bucket        string `json:"bucket"`        // OBS bucket where word dictionary files are stored.
	MainObj       string `json:"mainObj"`       // Main word dictionary file object.
	StopObj       string `json:"stopObj"`       // Stop word dictionary file object.
	SynonymObj    string `json:"synonymObj"`    // Synonym word dictionary file object.
	UpdateTime    int    `json:"updateTime"`    // Last word dictionary update time.
	UpdateDetails string `json:"updateDetails"` // Update details.
	ClusterId     string `json:"clusterId"`
	OperateStatus string `json:"operateStatus"`
	Id            string `json:"id"` // ID of a word dictionary.
}
