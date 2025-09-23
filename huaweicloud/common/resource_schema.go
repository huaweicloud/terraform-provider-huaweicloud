package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

// TagsSchema returns the schema to use for tags.
func TagsSchema(description ...string) *schema.Schema {
	schemaObj := schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	if len(description) > 0 {
		schemaObj.Description = description[0]
	}
	return &schemaObj
}

// TagsForceNewSchema returns the schema to use for tags with ForceNew behavior.
func TagsForceNewSchema(description ...string) *schema.Schema {
	schemaObj := schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	if len(description) > 0 {
		schemaObj.Description = description[0]
	}
	return &schemaObj
}

// TagsComputedSchema returns the schema to use for tags as an attribute.
func TagsComputedSchema(description ...string) *schema.Schema {
	schemaObj := schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	if len(description) > 0 {
		schemaObj.Description = description[0]
	}
	return &schemaObj
}

func SchemaChargingMode(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		Computed: true,
		ValidateFunc: validation.StringInSlice([]string{
			"prePaid", "postPaid",
		}, false),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func SchemaPeriodUnit(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		RequiredWith: []string{"period"},
		ValidateFunc: validation.StringInSlice([]string{
			"month", "year",
		}, false),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func SchemaPeriod(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:          schema.TypeInt,
		Optional:      true,
		ForceNew:      true,
		RequiredWith:  []string{"period_unit"},
		ValidateFunc:  validation.IntBetween(1, 9),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func SchemaAutoRenew(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			"true", "false",
		}, false),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func SchemaAutoRenewUpdatable(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			"true", "false",
		}, false),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func SchemaAutoPay(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			"true", "false",
		}, false),
		ConflictsWith: conflicts,
		Deprecated:    "Deprecated",
	}

	return &resourceSchema
}

func ValidatePrePaidChargeInfo(d *schema.ResourceData) error {
	if _, ok := d.GetOk("period_unit"); !ok {
		return fmtp.Errorf("both of `period, period_unit` must be specified in prePaid charging mode")
	}
	return nil
}
