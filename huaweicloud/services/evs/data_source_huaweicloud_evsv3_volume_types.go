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

// @API EVS GET /v3/{project_id}/types
func DataSourceV3VolumeTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3VolumeTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"volume_types": {
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
						"extra_specs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// This field in the API is `RESKEY:availability zones`,
									// the naming of fields in the schema cannot have ":",
									// so we did a mapping.
									"availability_zones": {
										Type:     schema.TypeString,
										Computed: true,
									},
									// This field in the API is `os-vendor-extended:sold_out_availability_zones`,
									// the naming of fields in the schema cannot have ":",
									// so we did a mapping.
									"sold_out_availability_zones": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceV3VolumeTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		mErr       *multierror.Error
		product    = "evs"
		getHttpUrl = "v3/{project_id}/types"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving EVS v3 volume types: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeTypesResp := utils.PathSearch("volume_types", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("volume_types", flattenV3VolumeTypes(volumeTypesResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV3VolumeTypes(volumeTypesResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(volumeTypesResp))
	for _, v := range volumeTypesResp {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"extra_specs": flattenExtraSpecs(utils.PathSearch("extra_specs", v, nil)),
			"description": utils.PathSearch("description", v, nil),
		})
	}

	return rst
}

func flattenExtraSpecs(extraSpecsResp interface{}) []map[string]interface{} {
	if extraSpecsResp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"availability_zones":          utils.PathSearch(`"RESKEY:availability_zones"`, extraSpecsResp, nil),
			"sold_out_availability_zones": utils.PathSearch(`"os-vendor-extended:sold_out_availability_zones"`, extraSpecsResp, nil),
		},
	}
}
