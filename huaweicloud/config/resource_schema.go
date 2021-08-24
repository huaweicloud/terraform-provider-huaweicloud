package config

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}
