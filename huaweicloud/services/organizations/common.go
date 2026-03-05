package organizations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	rootType     = "organizations:roots"
	unitType     = "organizations:ous"
	accountsType = "organizations:accounts"
	policiesType = "organizations:policies"
)

func getRoot(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v1/organizations/roots"
	getRootPath := client.Endpoint + httpUrl
	getRootOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getRootResp, err := client.Request("GET", getRootPath, &getRootOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getRootResp)
}

func addTags(client *golangsdk.ServiceClient, resourceType, resourceId string, tagList []tags.ResourceTag) error {
	var (
		addTagsToHttpUrl = "v1/organizations/{resource_type}/{resource_id}/tags/create"
	)

	addTagsToPath := client.Endpoint + addTagsToHttpUrl
	addTagsToPath = strings.ReplaceAll(addTagsToPath, "{resource_type}", resourceType)
	addTagsToPath = strings.ReplaceAll(addTagsToPath, "{resource_id}", resourceId)

	addTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildTagsBodyParams(tagList)),
	}
	_, err := client.Request("POST", addTagsToPath, &addTagsOpt)
	if err != nil {
		return fmt.Errorf("error creating tags of resourceType (%s): %s", resourceType, err)
	}

	return nil
}

func deleteTags(client *golangsdk.ServiceClient, resourceType, resourceId string, tagList []tags.ResourceTag) error {
	var (
		addTagsHttpUrl = "v1/organizations/{resource_type}/{resource_id}/tags/delete"
	)

	addTagsPath := client.Endpoint + addTagsHttpUrl
	addTagsPath = strings.ReplaceAll(addTagsPath, "{resource_type}", resourceType)
	addTagsPath = strings.ReplaceAll(addTagsPath, "{resource_id}", resourceId)

	addTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildTagsBodyParams(tagList)),
	}
	_, err := client.Request("POST", addTagsPath, &addTagsOpt)
	if err != nil {
		return fmt.Errorf("error deleting tags of resourceType (%s): %s", resourceType, err)
	}

	return nil
}

func buildTagsBodyParams(tagList []tags.ResourceTag) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"tags": tagList,
	}
	return bodyParams
}

func updateTags(d *schema.ResourceData, client *golangsdk.ServiceClient, resourceType, resourceId, tagsName string) error {
	oRaw, nRaw := d.GetChange(tagsName)
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})
	if len(oMap) > 0 {
		tagList := utils.ExpandResourceTags(oMap)
		err := deleteTags(client, resourceType, resourceId, tagList)
		if err != nil {
			return err
		}
	}
	if len(nMap) > 0 {
		tagList := utils.ExpandResourceTags(nMap)
		err := addTags(client, resourceType, resourceId, tagList)
		if err != nil {
			return err
		}
	}

	return nil
}

func getTags(client *golangsdk.ServiceClient, resourceType, resourceId string) (map[string]interface{}, error) {
	var (
		httpUrl = "v1/organizations/{resource_type}/{resource_id}/tags"
		limit   = 200
		marker  = ""
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{resource_type}", resourceType)
	listPath = strings.ReplaceAll(listPath, "{resource_id}", resourceId)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		getPathWithMarker := listPath
		if marker != "" {
			getPathWithMarker = fmt.Sprintf("%s&marker=%s", getPathWithMarker, marker)
		}

		resp, err := client.Request("GET", getPathWithMarker, &opt)
		if err != nil {
			return nil, fmt.Errorf("error getting tags of the resource (%s/%s): %s", resourceType, resourceId, err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		tagList := utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tagList...)
		if len(tagList) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return utils.FlattenTagsToMap(result), nil
}
