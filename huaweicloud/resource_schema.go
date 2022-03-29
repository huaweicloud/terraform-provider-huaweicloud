package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

// tagsSchema returns the schema to use for tags.
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

// tagsForceNewSchema returns the schema to use for tags with ForceNew
func tagsForceNewSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func schemaChargingMode(conflicts []string) *schema.Schema {
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

func schemaPeriodUnit(conflicts []string) *schema.Schema {
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

func schemaPeriod(conflicts []string) *schema.Schema {
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

func schemaAutoRenew(conflicts []string) *schema.Schema {
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

func schemaAutoPay(conflicts []string) *schema.Schema {
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

func validatePrePaidChargeInfo(d *schema.ResourceData) error {
	if _, ok := d.GetOk("period_unit"); !ok {
		return fmtp.Errorf("both of `period, period_unit` must be specified in prePaid charging mode")
	}
	return nil
}
