package lts

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
)

// The tag field information.
type TagsMap struct {
	// The key of the tag.
	Key string `json:"key" required:"true"`
	// The value of the tag.
	// The value can be an empty string.
	Value string `json:"value"`
	// Whether to apply to the log stream.
	TagsToStreamsEnable bool `json:"tags_to_streams_enable"`
}

func expandTagsToList(tagRaw map[string]interface{}) []TagsMap {
	var taglist []TagsMap

	for k, v := range tagRaw {
		tag := TagsMap{
			Key:                 k,
			Value:               v.(string),
			TagsToStreamsEnable: false,
		}
		taglist = append(taglist, tag)
	}

	return taglist
}

func updateTags(client *golangsdk.ServiceClient, resourceType, resourceId string, d *schema.ResourceData) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	httpUrl := "v1/{project_id}/{resource_type}/{resource_id}/tags/action"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{resource_type}", resourceType)
	path = strings.ReplaceAll(path, "{resource_id}", resourceId)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		OkCodes: []int{
			// For log stream, the status code of deleting tags is 201.
			200, 201,
		},
	}
	// remove old tags
	if len(oMap) > 0 {
		opts.JSONBody = map[string]interface{}{
			"action": "delete",
			"tags":   expandTagsToList(oMap),
		}
		_, err := client.Request("POST", path, &opts)
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		opts.JSONBody = map[string]interface{}{
			"action": "create",
			"tags":   expandTagsToList(nMap),
		}

		_, err := client.Request("POST", path, &opts)
		if err != nil {
			return err
		}
	}
	return nil
}
