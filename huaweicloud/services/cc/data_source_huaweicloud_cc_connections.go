package cc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC GET /v3/{domain_id}/ccaas/cloud-connections
// @API CC POST /v3/{domain_id}/ccaas/cloud-connections/filter
func DataSourceCloudConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Cloud connection ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Cloud connection name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Cloud connection description.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Cloud connection status.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Enterprise project ID.`,
			},
			"used_scene": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Application scenario.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Cloud connection tags.`,
			},
			"connections": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Cloud connection list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Cloud connection ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Cloud connection name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Cloud connection description.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `ID of the account that the instance belongs to.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `ID of the enterprise project that the cloud connection belongs to.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Time when the cloud connection was created.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Time when the cloud connection was updated.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Cloud connection tags.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Cloud connection status.`,
						},
						"used_scene": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application scenarios of the cloud connection.`,
						},
						"network_instance_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Number of the network instances loaded to the cloud connection.`,
						},
						"bandwidth_package_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Number of the bandwidth packages bound to the cloud connection.`,
						},
						"inter_region_bandwidth_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Number of the inter-region bandwidths configured for the cloud connection.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)

	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	result, err := getCloudConnection(client, cfg, d)
	if err != nil {
		return diag.Errorf("error retrieving cloud connections: %s", err)
	}

	if tags, ok := d.GetOk("tags"); ok {
		resourceIDs, err := filterCloudConnectionsByTags(tags.(map[string]interface{}), client, cfg.DomainID)
		if err != nil {
			return diag.Errorf("error filtering global connection bandwidths by tags: %s", err)
		}
		result = filter(result, resourceIDs)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("connections", flattenListCloudConnectionsResponseBody(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getCloudConnection(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v3/{domain_id}/ccaas/cloud-connections"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)

	params := buildCloudConnectionsQueryParams(d, cfg)
	path += params

	resp, err := pagination.ListAllItems(
		client,
		"marker",
		path,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return nil, err
	}

	curJson := utils.PathSearch("cloud_connections", respBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	return curArray, nil
}

func buildCloudConnectionsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)

	if v, ok := d.GetOk("connection_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("used_scene"); ok {
		res = fmt.Sprintf("%s&used_scene=%v", res, v)
	}
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsID)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func filterCloudConnectionsByTags(tags map[string]interface{}, client *golangsdk.ServiceClient, domainID string) ([]string, error) {
	var resourceIDs []string
	httpUrl := "v3/{domain_id}/ccaas/cloud-connections/filter"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", domainID)

	if len(tags) < 1 {
		return nil, nil
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildFilterCloudConnectionByTagsOpts(tags),
	}

	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	curJson := utils.PathSearch("cloud_connections[*].id", respBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	for _, item := range curArray {
		resourceIDs = append(resourceIDs, item.(string))
	}

	return resourceIDs, nil
}

func buildFilterCloudConnectionByTagsOpts(tagmap map[string]interface{}) map[string]interface{} {
	taglist := make([]interface{}, 0, len(tagmap))

	for k, v := range tagmap {
		taglist = append(taglist, map[string]interface{}{
			"key":    k,
			"values": []string{v.(string)},
		})
	}

	return map[string]interface{}{
		"tags": taglist,
	}
}

func flattenListCloudConnectionsResponseBody(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(items))
	for _, v := range items {
		rst = append(rst, map[string]interface{}{
			"id":                            utils.PathSearch("id", v, nil),
			"name":                          utils.PathSearch("name", v, nil),
			"description":                   utils.PathSearch("description", v, nil),
			"domain_id":                     utils.PathSearch("domain_id", v, nil),
			"enterprise_project_id":         utils.PathSearch("enterprise_project_id", v, nil),
			"created_at":                    utils.PathSearch("created_at", v, nil),
			"updated_at":                    utils.PathSearch("updated_at", v, nil),
			"tags":                          utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"status":                        utils.PathSearch("status", v, nil),
			"used_scene":                    utils.PathSearch("used_scene", v, nil),
			"network_instance_number":       utils.PathSearch("network_instance_number", v, nil),
			"bandwidth_package_number":      utils.PathSearch("bandwidth_package_number", v, nil),
			"inter_region_bandwidth_number": utils.PathSearch("inter_region_bandwidth_number", v, nil),
		})
	}

	return rst
}
