package policygroups

// PolicyGroup is the structure that represents the policy group detail.
type PolicyGroup struct {
	// Policy group name.
	PolicyGroupName string `json:"policy_group_name" required:"true"`
	// Policy group ID.
	PolicyGroupId string `json:"policy_group_id,omitempty"`
	// Priority.
	Priority int `json:"priority,omitempty"`
	// Policy group description.
	Description string `json:"description,omitempty"`
	// The update time of the policy group.
	UpdateTime string `json:"update_time,omitempty"`
	// List of target objects.
	Targets []Target `json:"targets,omitempty"`
	// List of policy.
	Policy Policy `json:"policies,omitempty"`
}
