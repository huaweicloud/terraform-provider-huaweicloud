package utils

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

const SysTagKeyEnterpriseProjectId = "_sys_enterprise_project_id"

// CreateResourceTags is a helper to create the tags for a resource.
// It expects the schema name must be "tags"
func CreateResourceTags(client *golangsdk.ServiceClient, d *schema.ResourceData, resourceType, id string) error {
	if tagRaw := d.Get("tags").(map[string]interface{}); len(tagRaw) > 0 {
		tagList := ExpandResourceTags(tagRaw)
		return tags.Create(client, resourceType, id, tagList).ExtractErr()
	}
	return nil
}

// UpdateResourceTags is a helper to update the tags for a resource.
// It expects the tags field to be named "tags"
func UpdateResourceTags(conn *golangsdk.ServiceClient, d *schema.ResourceData, resourceType, id string) error {
	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		// remove old tags
		if len(oMap) > 0 {
			taglist := ExpandResourceTags(oMap)
			err := tags.Delete(conn, resourceType, id, taglist).ExtractErr()
			if err != nil {
				return err
			}
		}

		// set new tags
		if len(nMap) > 0 {
			taglist := ExpandResourceTags(nMap)
			err := tags.Create(conn, resourceType, id, taglist).ExtractErr()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteResourceTagsWithKeys is a helper to delete the tags with tagKeys for a resource.
func DeleteResourceTagsWithKeys(client *golangsdk.ServiceClient, tagKeys []string, resourceType, id string) error {
	for _, key := range tagKeys {
		if err := tags.DeleteWithKey(client, resourceType, id, key).ExtractErr(); err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[WARN] The tag key (%s) of resource (%s) not exist.", key, id)
				continue
			}
			return err
		}
	}
	return nil
}

// SetResourceTagsToState is a helper to query tags of resource, then set to state.
// The schema argument name must be: tags
func SetResourceTagsToState(d *schema.ResourceData, client *golangsdk.ServiceClient, resourceType, id string) error {
	// set tags
	if resourceTags, err := tags.Get(client, resourceType, id).Extract(); err == nil {
		tagmap := TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmt.Errorf("error saving tags to state for %s (%s): %s", resourceType, id, err)
		}
	} else {
		log.Printf("[WARN] Error fetching tags of %s (%s): %s", resourceType, id, err)
	}
	return nil
}

// TagsToMap returns the list of tags into a map.
func TagsToMap(tags []tags.ResourceTag) map[string]string {
	result := make(map[string]string)
	for _, val := range tags {
		result[val.Key] = val.Value
	}

	// ignore system tags to keep the tags consistent with what the user set
	delete(result, "CCE-Cluster-ID")
	delete(result, "CCE-Dynamic-Provisioning-Node")

	return result
}

// FlattenTagsToMap returns the list of tags into a map.
func FlattenTagsToMap(tags interface{}) map[string]interface{} {
	if tagArray, ok := tags.([]interface{}); ok {
		result := make(map[string]interface{})
		for _, val := range tagArray {
			if t, ok := val.(map[string]interface{}); ok {
				result[t["key"].(string)] = t["value"]
			}
		}
		return result
	}

	return nil
}

// FlattenSameKeyTagsToMap using to flatten remote tag list and filter tag map with same key.
// Parameters:
// + d         : It is required that there is a `tags` field of map type in the parameter.
// + remoteTags: It is required that the parameter is a tag list, for example: [{"key":"key1","value":"value1"}].
func FlattenSameKeyTagsToMap(d *schema.ResourceData, remoteTags interface{}) map[string]interface{} {
	return FilterMapWithSameKey(d.Get("tags").(map[string]interface{}), FlattenTagsToMap(remoteTags))
}

// ExpandResourceTags returns the tags for the given map of data.
func ExpandResourceTags(tagmap map[string]interface{}) []tags.ResourceTag {
	var taglist []tags.ResourceTag

	for k, v := range tagmap {
		tag := tags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	return taglist
}

// ExpandResourceTagsMap returns the tags in format of list of maps for the given map of data.
// Parameters:
// + tagmap: tags input with the map structure format.
// + keepEmpty: whether to keep the empty list return instead of the nil.
func ExpandResourceTagsMap(tagmap map[string]interface{}, keepEmpty ...bool) []map[string]interface{} {
	// Empty list returned only keepEmpty is set and the value is 'true'.
	if len(tagmap) < 1 && (len(keepEmpty) < 1 || (len(keepEmpty) > 0 && !keepEmpty[0])) {
		// If not, returns nil and with the type, the value is '[]map[string]interface{}(nil)'.
		return nil
	}

	taglist := make([]map[string]interface{}, 0, len(tagmap))

	for k, v := range tagmap {
		tag := map[string]interface{}{
			"key":   k,
			"value": v,
		}
		taglist = append(taglist, tag)
	}

	return taglist
}

// GetDNSZoneTagType returns resource tag type of DNS zone by zoneType
func GetDNSZoneTagType(zoneType string) (string, error) {
	if zoneType == "public" {
		return "DNS-public_zone", nil
	} else if zoneType == "private" {
		return "DNS-private_zone", nil
	}
	return "", fmt.Errorf("invalid zone type: %s", zoneType)
}

// GetDNSRecordSetTagType returns resource tag type of DNS record set by zoneType
func GetDNSRecordSetTagType(zoneType string) (string, error) {
	if zoneType == "public" {
		return "DNS-public_recordset", nil
	} else if zoneType == "private" {
		return "DNS-private_recordset", nil
	}
	return "", fmt.Errorf("invalid zone type: %s", zoneType)
}

func ParseEnterpriseProjectIdFromSysTags(value []tags.ResourceTag) (enterpriseProjectId string) {
	if len(value) == 0 {
		return
	}

	for i := 0; i < len(value); i++ {
		item := value[i]
		if item.Key == SysTagKeyEnterpriseProjectId {
			return item.Value
		}
	}
	return
}

func BuildSysTags(enterpriseProjectID string) (enterpriseProjectTags []tags.ResourceTag) {
	if enterpriseProjectID != "" {
		t := tags.ResourceTag{
			Key:   SysTagKeyEnterpriseProjectId,
			Value: enterpriseProjectID,
		}
		enterpriseProjectTags = append(enterpriseProjectTags, t)
	}
	return
}
