package dli

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v3/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func addTags(client *golangsdk.ServiceClient, resourceId, resourceType string, raw map[string]interface{}) error {
	opts := tags.UpdateTagsOpts{
		Action:       "create",
		ResourceType: resourceType,
		Tags:         utils.ExpandResourceTags(raw),
	}
	return tags.UpdateTagsToResource(client, resourceId, opts)
}

func removeTags(client *golangsdk.ServiceClient, resourceId, resourceType string, oldRaw map[string]interface{}) error {
	opts := tags.UpdateTagsOpts{
		Action:       "delete",
		ResourceType: resourceType,
		Tags:         utils.ExpandResourceTags(oldRaw),
	}
	return tags.UpdateTagsToResource(client, resourceId, opts)
}

func updateResourceTags(client *golangsdk.ServiceClient, resourceId, resourceType string, oldTags, newTags interface{}) error {
	// remove old tags and set new tags
	oldTagsRaw := oldTags.(map[string]interface{})
	if len(oldTagsRaw) > 0 {
		if err := removeTags(client, resourceId, resourceType, oldTagsRaw); err != nil {
			return err
		}
	}

	newTagsRaw := newTags.(map[string]interface{})
	if len(newTagsRaw) > 0 {
		if err := addTags(client, resourceId, resourceType, newTagsRaw); err != nil {
			return err
		}
	}

	return nil
}
