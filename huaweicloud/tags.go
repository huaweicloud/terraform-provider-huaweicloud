package huaweicloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// tagsSchema returns the schema to use for tags.
//
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}
