package cts

import (
	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"
)

func flattenNotificationUsers(users []cts.NotificationUsers) []map[string]interface{} {
	ret := make([]map[string]interface{}, len(users))
	for i, v := range users {
		ret[i] = map[string]interface{}{
			"group": v.UserGroup,
			"users": v.UserList,
		}
	}

	return ret
}

func flattenNotificationOperations(ops []cts.Operations) []map[string]interface{} {
	ret := make([]map[string]interface{}, len(ops))
	for i, v := range ops {
		ret[i] = map[string]interface{}{
			"service":     v.ServiceType,
			"resource":    v.ResourceType,
			"trace_names": v.TraceNames,
		}
	}

	return ret
}

func flattenNotificationFilter(filter *cts.Filter) []map[string]interface{} {
	if filter == nil {
		return nil
	}
	result := map[string]interface{}{
		"condition": filter.Condition.Value(),
		"rule":      filter.Rule,
	}

	return []map[string]interface{}{result}
}
