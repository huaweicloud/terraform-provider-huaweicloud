package mrs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API MRS GET /v1.1/{region_id}/available-zones
func DataSourceAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailabilityZonesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the availability zones are located.`,
			},
			// Optional parameters.
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The scope of the availability zone.`,
			},
			// Attributes.
			"default_az_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default availability zone code.`,
			},
			"support_physical_az_group": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the physical availability zone grouping is supported.`,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone code.`,
						},
						"az_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone ID.`,
						},
						"az_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone code.`,
						},
						"az_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone name.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone status.`,
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The region ID.`,
						},
						"az_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone group ID.`,
						},
						"az_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone type.`,
						},
						"az_category": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The availability zone category.`,
						},
						"charge_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The charge policy of the availability zone.`,
						},
						"az_tags": {
							Type:        schema.TypeList,
							Elem:        availabilityZoneTagsSchema(),
							Computed:    true,
							Description: `The availability zone tags.`,
						},
					},
				},
				Description: `The availability zone list that matched the filter parameters.`,
			},
		},
	}
}

func availabilityZoneTagsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The mode of the availability zone.`,
			},
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alias of the availability zone.`,
			},
			"public_border_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public border group to which the availability zone belongs.`,
			},
		},
	}
	return &sc
}

func dataSourceAvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	respBody, err := listAvailabilityZones(client, d, region)
	if err != nil {
		return diag.Errorf("error retrieving availability zones: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("default_az_code", utils.PathSearch("default_az_code", respBody, nil)),
		d.Set("support_physical_az_group", utils.PathSearch("support_physical_az_group", respBody, nil)),
		d.Set("available_zones", flattenAvailabilityZones(utils.PathSearch("available_zones",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listAvailabilityZones(client *golangsdk.ServiceClient, d *schema.ResourceData, region string) (interface{}, error) {
	httpUrl := "v1.1/{region_id}/available-zones"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{region_id}", region)

	if scope, ok := d.GetOk("scope"); ok {
		listPath += fmt.Sprintf("?scope=%v", scope)
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func flattenAvailabilityZones(availabilityZones []interface{}) []map[string]interface{} {
	if len(availabilityZones) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(availabilityZones))
	for _, v := range availabilityZones {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"az_code":       utils.PathSearch("az_code", v, nil),
			"az_name":       utils.PathSearch("az_name", v, nil),
			"az_id":         utils.PathSearch("az_id", v, nil),
			"status":        utils.PathSearch("status", v, nil),
			"region_id":     utils.PathSearch("region_id", v, nil),
			"az_group_id":   utils.PathSearch("az_group_id", v, nil),
			"az_type":       utils.PathSearch("az_type", v, nil),
			"az_category":   utils.PathSearch("az_category", v, nil),
			"charge_policy": utils.PathSearch("charge_policy", v, nil),
			"az_tags":       flattenAvailabilityZoneTags(utils.PathSearch("az_tags", v, nil)),
		})
	}
	return result
}

func flattenAvailabilityZoneTags(tag interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"mode":                utils.PathSearch("mode", tag, nil),
			"alias":               utils.PathSearch("alias", tag, nil),
			"public_border_group": utils.PathSearch("public_border_group", tag, nil),
		},
	}
}
