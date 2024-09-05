package eip

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v1/{project_id}/bandwidths
func DataSourceBandWidth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBandWidthRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicips": publicIPListComputedSchema(),
		},
	}
}

func dataSourceBandWidthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bwClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating bandwidth v1 client: %s", err)
	}

	listOpts := bandwidths.ListOpts{
		ShareType:           "WHOLE",
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
	}

	allBWs, err := bandwidths.List(bwClient, listOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to list bandwidths: %s", err)
	}
	if len(allBWs) == 0 {
		return diag.Errorf("bandwidth was not found")
	}

	filter := map[string]interface{}{
		"Name": d.Get("name").(string),
	}
	if v, ok := d.GetOk("size"); ok {
		filter["Size"] = v
	}

	filterBWs, err := utils.FilterSliceWithField(allBWs, filter)
	if err != nil {
		return diag.Errorf("filter bandwidths failed: %s", err)
	}
	if len(filterBWs) == 0 {
		return diag.Errorf("bandwidth was not found by %+v", filter)
	}

	result := filterBWs[0].(bandwidths.BandWidth)
	log.Printf("[DEBUG] Retrieved bandwidth %s: %+v", result.ID, result)

	d.SetId(result.ID)
	mErr := multierror.Append(
		d.Set("name", result.Name),
		d.Set("size", result.Size),
		d.Set("enterprise_project_id", result.EnterpriseProjectID),

		d.Set("share_type", result.ShareType),
		d.Set("bandwidth_type", result.BandwidthType),
		d.Set("charge_mode", result.ChargeMode),
		d.Set("status", result.Status),
		d.Set("publicips", flattenPublicIPs(result)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting bandwidth fields: %s", err)
	}

	return nil
}
