package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// tagsSchema returns the schema to use for tags.
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}
