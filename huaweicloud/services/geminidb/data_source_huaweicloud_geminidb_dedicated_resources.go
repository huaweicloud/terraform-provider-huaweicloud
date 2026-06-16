package geminidb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3/{project_id}/dedicated-resources
func DataSourceDedicatedResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedResourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vcpus": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ram": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
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

func dataSourceDedicatedResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dedicated-resources?limit=100"
		// The `offset` value is `0` in API document, Actually, the `offset` value is `1`.
		offset = 1
		result = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the dedicated resources: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dedicatedResources := utils.PathSearch("resources", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dedicatedResources) == 0 {
			break
		}

		result = append(result, dedicatedResources...)
		offset += len(dedicatedResources)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("resources", flattenDedicatedResources(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDedicatedResources(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"resource_name":     utils.PathSearch("resource_name", v, nil),
			"engine_name":       utils.PathSearch("engine_name", v, nil),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
			"architecture":      utils.PathSearch("architecture", v, nil),
			"capacity":          flattenDedicatedResourceCapacity(utils.PathSearch("capacity", v, nil)),
			"status":            utils.PathSearch("status", v, nil),
		})
	}

	return result
}

func flattenDedicatedResourceCapacity(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"vcpus":  utils.PathSearch("vcpus", resp, nil),
		"ram":    utils.PathSearch("ram", resp, nil),
		"volume": utils.PathSearch("volume", resp, nil),
	}

	return []map[string]interface{}{result}
}
