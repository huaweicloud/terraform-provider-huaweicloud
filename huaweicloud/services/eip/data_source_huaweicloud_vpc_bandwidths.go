package eip

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v1/{project_id}/bandwidths
func DataSourceBandWidths() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBandWidthsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidths": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
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
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"publicips": publicIPListComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBandWidthsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	filter := map[string]interface{}{}
	if v, ok := d.GetOk("name"); ok {
		filter["Name"] = v
	}
	if v, ok := d.GetOk("bandwidth_id"); ok {
		filter["ID"] = v
	}
	if v, ok := d.GetOk("size"); ok {
		filter["Size"] = v
	}
	if v, ok := d.GetOk("charge_mode"); ok {
		filter["ChargeMode"] = v
	}

	filterBWs, err := utils.FilterSliceWithField(allBWs, filter)
	if err != nil {
		return diag.Errorf("filter bandwidths failed: %s", err)
	}

	rst := make([]map[string]interface{}, len(filterBWs))
	for i, v := range filterBWs {
		item := v.(bandwidths.BandWidth)
		bandwidth := map[string]interface{}{
			"id":                    item.ID,
			"name":                  item.Name,
			"size":                  item.Size,
			"enterprise_project_id": item.EnterpriseProjectID,
			"share_type":            item.ShareType,
			"bandwidth_type":        item.BandwidthType,
			"charge_mode":           item.ChargeMode,
			"status":                item.Status,
			"created_at":            item.CreatedAt,
			"updated_at":            item.UpdatedAt,
			"publicips":             flattenPublicIPs(item),
		}
		rst[i] = bandwidth
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("bandwidths", rst),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting bandwidth fields: %s", err)
	}

	return nil
}
