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

// @API Workspace GET /v1/{project_id}/availability-zone/summary
func DataSourceIesAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIesAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the availability zones are located.",
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of availability zones.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the availability zone.`,
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the availability zone.",
						},
						"i18n": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The internationalization information of the availability zone.",
						},
						"sold_out": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The sold out information for the availability zone.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"products": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The list of sold out product IDs.",
									},
								},
							},
						},
						"product_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of custom supported product IDs for the availability zone.",
						},
						"visible": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the availability zone is visible.",
						},
						"default_availability_zone": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is the default availability zone.",
						},
					},
				},
			},
		},
	}
}

func getIesAvailabilityZones(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/availability-zone/summary"
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

	return utils.PathSearch("azs.IES", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceIesAvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	allZones, err := getIesAvailabilityZones(client)
	if err != nil {
		return diag.Errorf("error querying availability zones: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	// The current structure returned by IES is the same as the center, so it remains consistent at present.
	// If there are any changes to the edge cloud in the future, flattened availability zones need to be modified.
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("availability_zones", flattenAvailabilityZones(allZones)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
