package jobs

// Job is the structure that represents the detail of the process job.
type Job struct {
	// Job ID.
	ID string `json:"id"`
	// Job name.
	Name string `json:"name"`
	// Status info.
	// + Running
	// + Completed
	// + Failed
	Status string `json:"status"`
	// Creation time, the format is "yyyy-mm-ddThh:mm:ssZ".
	Created string `json:"created"`
	// End time, the format is "yyyy-mm-ddThh:mm:ssZ".
	Ended string `json:"ended"`
	// The execution progress of the job.
	Progress string `json:"progress"`
	// The DDS instance info to which the job belongs.
	Instance Instance `json:"instance"`
	// Error information generated when the job fails to be executed.
	FailReason string `json:"fail_reason"`
}

// Instance is the structure that represents the detail of the DDS instnace.
type Instance struct {
	// Instance ID.
	ID string `json:"id"`
	// Instance name.
	Name string `json:"name"`
}
