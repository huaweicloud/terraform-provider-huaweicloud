package dew

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

// @API DEW GET /v1/{project_id}/dew/cpcs/az
func DataSourceAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locales": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"en_us": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zh_cn": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/az"
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving availability zones: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	availabilityZones := utils.PathSearch("availability_zone", respBody, make([]interface{}, 0)).([]interface{})

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("availability_zone", flattenAvailabilityZones(availabilityZones)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAvailabilityZones(availabilityZones []interface{}) []interface{} {
	if len(availabilityZones) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(availabilityZones))
	for _, availabilityZone := range availabilityZones {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", availabilityZone, nil),
			"display_name": utils.PathSearch("display_name", availabilityZone, nil),
			"locales":      flattenZonesLocales(utils.PathSearch("locales", availabilityZone, nil)),
			"type":         utils.PathSearch("type", availabilityZone, nil),
			"region_id":    utils.PathSearch("region_id", availabilityZone, nil),
			"status":       utils.PathSearch("status", availabilityZone, nil),
		})
	}

	return result
}

func flattenZonesLocales(localesData interface{}) []map[string]interface{} {
	if localesData == nil {
		return nil
	}

	result := map[string]interface{}{
		"en_us": utils.PathSearch(`"en-us"`, localesData, nil),
		"zh_cn": utils.PathSearch(`"zh-cn"`, localesData, nil),
	}

	return []map[string]interface{}{result}
}
