package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	iamPollTimeout      = 20 * time.Second
	iamPollInitialDelay = 1 * time.Second
	iamPollInterval     = 1 * time.Second
)

// PollRequest polls an IAM resource until the request function no longer returns 404.
func PollRequest(ctx context.Context, requestFunc func() (interface{}, error)) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(iamPollInitialDelay):
	}

	deadline := time.Now().Add(iamPollTimeout)

	for {
		resp, err := requestFunc()
		if err == nil {
			return resp, nil
		}

		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return nil, err
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for IAM resources to become expected status after %s", iamPollTimeout)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(iamPollInterval):
		}
	}
}

func buildV5ResourceTags(tagmap map[string]interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, 0, len(tagmap))
	for k, v := range tagmap {
		tags = append(tags, map[string]interface{}{
			"tag_key":   k,
			"tag_value": v,
		})
	}

	return tags
}

func addV5TagsToResource(client *golangsdk.ServiceClient, tags map[string]interface{}, resourceType, resourceId string) error {
	httpUrl := "v5/{resource_type}/{resource_id}/tags/create"
	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{resource_type}", resourceType)
	addPath = strings.ReplaceAll(addPath, "{resource_id}", resourceId)
	addOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"tags": buildV5ResourceTags(tags),
		},
	}
	_, err := client.Request("POST", addPath, &addOpt)
	return err
}

func expandV5TagsKeyToStringList(tagmap map[string]interface{}) []string {
	tagKeys := make([]string, 0, len(tagmap))
	for k := range tagmap {
		tagKeys = append(tagKeys, k)
	}

	return tagKeys
}

func removeV5TagsFromResource(client *golangsdk.ServiceClient, tags map[string]interface{}, resourceType, resourceId string) error {
	httpUrl := "v5/{resource_type}/{resource_id}/tags/delete"
	removePath := client.Endpoint + httpUrl
	removePath = strings.ReplaceAll(removePath, "{resource_type}", resourceType)
	removePath = strings.ReplaceAll(removePath, "{resource_id}", resourceId)
	removeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         expandV5TagsKeyToStringList(tags),
	}
	_, err := client.Request("DELETE", removePath, &removeOpt)
	return err
}

func getV5ResourceTags(client *golangsdk.ServiceClient, resourceType string, resourceId string) ([]interface{}, error) {
	getHttpUrl := "v5/{resource_type}/{resource_id}/tags"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{resource_type}", resourceType)
	getPath = strings.ReplaceAll(getPath, "{resource_id}", resourceId)
	getOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("tags", getRespBody, make([]interface{}, 0)).([]interface{}), nil
}
