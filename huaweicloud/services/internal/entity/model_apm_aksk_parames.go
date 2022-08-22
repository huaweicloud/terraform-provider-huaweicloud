package entity

import "encoding/json"

type CreateAkSkParam struct {
	Descp string `json:"descp,omitempty"`
}

type AkSkResultVO struct {
	Ak string `json:"ak,omitempty"`
	Sk string `json:"sk,omitempty"`
}

type GetAkSkListVO struct {
	AccessAkSkModels []AccessAkSkVO `json:"access_ak_sk_models,omitempty"`
}

type AccessAkSkVO struct {
	Id            int    `json:"id,omitempty"`
	GmtCreate     string `json:"gmt_create,omitempty"`
	GmtModify     string `json:"gmt_modify,omitempty"`
	InnerDomainId int    `json:"inner_domain_id,omitempty"`
	Ak            string `json:"ak,omitempty"`
	Sk            string `json:"sk,omitempty"`
	Status        string `json:"status,omitempty"`
	Descp         string `json:"descp,omitempty"`
}

func (c AccessAkSkVO) ToString() string {
	b, err := json.Marshal(&c)
	if err != nil {
		return "failed"
	}
	return string(b)
}

func (c AkSkResultVO) ToString() string {
	b, err := json.Marshal(&c)
	if err != nil {
		return "failed"
	}
	return string(b)
}

type ErrorResp struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}
