package workspace

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

// @API Workspace GET /v1/{project_id}/volume-type
func DataSourceAppAvailableVolumes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppAvailableVolumesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource.`,
			},
			"volume_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of available volume types.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource specification code.`,
						},
						"volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The volume type.`,
						},
						"volume_product_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The volume product type.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type.`,
						},
						"cloud_service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cloud service type code.`,
						},
						"name": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The volume type name in different languages.`,
							Elem: &schema.Schema{
								Type: schema.TypeMap,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
						"volume_type_extra_specs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The extra specifications of volume type.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The availability zone for this volume type.`,
									},
									"sold_out_availability_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The sold out availability zone for this volume type.`,
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

func getAppAvailableVolumes(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/volume-type"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("volume_types", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceAppAvailableVolumesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace app client: %s", err)
	}

	volumes, err := getAppAvailableVolumes(client)
	if err != nil {
		return diag.Errorf("error querying volume types: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("volume_types", flattenVolumeTypes(volumes)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVolumeTypes(volumeTypes []interface{}) []map[string]interface{} {
	if len(volumeTypes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(volumeTypes))
	for _, volumeType := range volumeTypes {
		vt := volumeType.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"resource_spec_code":  utils.PathSearch("resource_spec_code", vt, nil),
			"volume_type":         utils.PathSearch("volume_type", vt, nil),
			"volume_product_type": utils.PathSearch("volume_product_type", vt, nil),
			"resource_type":       utils.PathSearch("resource_type", vt, nil),
			"cloud_service_type":  utils.PathSearch("cloud_service_type", vt, nil),
			"name":                utils.PathSearch("name", vt, nil),
			"volume_type_extra_specs": flattenVolumeTypeExtraSpecs(utils.PathSearch("volume_type_extra_specs", vt,
				make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return result
}

func flattenVolumeTypeExtraSpecs(specs map[string]interface{}) []map[string]interface{} {
	if len(specs) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"availability_zone":          utils.PathSearch("availability_zone", specs, nil),
			"sold_out_availability_zone": utils.PathSearch("sold_out_availability_zone", specs, nil),
		},
	}
}
