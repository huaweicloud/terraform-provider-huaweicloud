package entity

type CreateLogtankRequestBody struct {
	Logtank *CreateLogTankOption `json:"logtank"`
}
type CreateLogTankOption struct {
	LoadBalancerId string `json:"loadbalancer_id"`

	LogGroupId string `json:"log_group_id"`

	LogTopicId string `json:"log_topic_id"`
}
type CreateLogtankResponse struct {
	Logtank *Logtank `json:"logtank,omitempty"`

	RequestId *string `json:"request_id,omitempty"`
}

type Logtank struct {
	ID             string `json:"id"`
	ProjectID      string `json:"project_id"`
	LoadBalancerID string `json:"loadbalancer_id"`
	LogGroupID     string `json:"log_group_id"`
	LogTopicID     string `json:"log_topic_id"`
}
