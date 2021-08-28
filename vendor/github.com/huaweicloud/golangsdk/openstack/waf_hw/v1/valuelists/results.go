/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package valuelists

// WafValueList reference table detail
type WafValueList struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	CreationTime int64    `json:"timestamp"`
	Values       []string `json:"values"`
}

type ListValueListRst struct {
	Total int            `json:"total"`
	Items []WafValueList `json:"items"`
}
