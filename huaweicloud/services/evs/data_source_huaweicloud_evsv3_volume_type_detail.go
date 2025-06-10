package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS GET /v3/{project_id}/types/{type_id}
func DataSourceEvsv3VolumeTypeDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsv3VolumeTypeDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_type": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extra_specs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reskey_availability_zones": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_vendor_extended_sold_out_availability_zones": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceEvsv3VolumeTypeDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		product = "evs"
		httpUrl = "v3/{project_id}/types/{type_id}"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS v3 client: %s", err)
	}
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{type_id}", d.Get("type_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving EVS volume type detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	volumeType := utils.PathSearch("volume_type", respBody, make(map[string]interface{})).(map[string]interface{})

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("volume_type", flattenVolumeTypeDetail(volumeType)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVolumeTypeDetail(volumeTypeDetail map[string]interface{}) []map[string]interface{} {
	if len(volumeTypeDetail) == 0 {
		return nil
	}
	return []map[string]interface{}{{
		"id":          utils.PathSearch("id", volumeTypeDetail, nil),
		"name":        utils.PathSearch("name", volumeTypeDetail, nil),
		"description": utils.PathSearch("description", volumeTypeDetail, nil),
		"extra_specs": flattenVolumeTypeDetailExtraSpecs(utils.PathSearch("extra_specs", volumeTypeDetail, nil)),
	}}
}

func flattenVolumeTypeDetailExtraSpecs(extraSpecs interface{}) interface{} {
	if extraSpecs == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"reskey_availability_zones":                      utils.PathSearch("\"RESKEY:availability_zones\"", extraSpecs, nil),
			"os_vendor_extended_sold_out_availability_zones": utils.PathSearch("\"os-vendor-extended:sold_out_availability_zones\"", extraSpecs, nil),
		},
	}
}
