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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC GET /v3/{domain_id}/gcb/gcbandwidths
// @API CC POST /v3/gcb/resource-instances/filter
func DataSourceCcGlobalConnectionBandwidths() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCGlobalConnectionBandwidthsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Resource name.`,
			},
			"gcb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Resource ID.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Range of a global connection bandwidth, in Mbit/s.`,
			},
			"admin_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Status of the global connection bandwidth.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Type of a global connection bandwidth.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Bound instance ID.`,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Instance type.`,
			},
			"binding_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Binding service.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Billing option.`,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Enterprise project ID.`,
			},
			"globalconnection_bandwidths": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Response body for querying the global connection bandwidth list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Resource ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Resource name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Resource description.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `ID of the account that the resource belongs to.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Type of a global connection bandwidth.`,
						},
						"bordercross": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the global connection bandwidth is used for cross-border communications.`,
						},
						"binding_service": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Binding service.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `ID of the enterprise project that the global connection bandwidth belongs to.`,
						},
						"charge_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Billing option.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Range of a global connection bandwidth, in Mbit/s.`,
						},
						"sla_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Class of a global connection bandwidth.`,
						},
						"local_site_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Code of the local access point.`,
						},
						"remote_site_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Code of the remote access point.`,
						},
						"admin_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Global connection bandwidth status.`,
						},
						"remote_area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Name of a remote access point.`,
						},
						"local_area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Name of a local access point.`,
						},
						"frozen": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether a global connection bandwidth is frozen.`,
						},
						"spec_code_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `UUID of a line specification code.`,
						},
						"tags": common.TagsComputedSchema(),
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Time when the resource was created.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Time when the resource was updated.`,
						},
						"enable_share": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether a global connection bandwidth can be used by multiple instances.`,
						},
						"instances": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of instances that the global connection bandwidth is bound to.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Bound instance ID.`,
									},
									"project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Project ID of the bound instance.`,
									},
									"region_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Region of the bound instance.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Bound instance type.`,
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

func dataSourceCCGlobalConnectionBandwidthsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)

	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	result, err := getGCB(client, cfg, d)
	if err != nil {
		return diag.Errorf("error retrieving global connection bandwidths: %s", err)
	}

	if tags, ok := d.GetOk("tags"); ok {
		resourceIDs, err := filterGCBByTags(tags.(map[string]interface{}), client)
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
		d.Set("globalconnection_bandwidths", flattenListGCBResponseBody(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getGCB(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v3/{domain_id}/gcb/gcbandwidths"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)

	params := buildGlobalConnectionBandwidthsQueryParams(d, cfg)
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

	curJson := utils.PathSearch("globalconnection_bandwidths", respBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))

	if sizeFilter, ok := d.GetOk("size"); ok {
		for _, item := range curArray {
			itemMap := item.(map[string]interface{})
			if int(itemMap["size"].(float64)) != sizeFilter {
				continue
			}
			rst = append(rst, item)
		}
	} else {
		rst = append(rst, curArray...)
	}

	return rst, nil
}

func buildGlobalConnectionBandwidthsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("gcb_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("admin_state"); ok {
		res = fmt.Sprintf("%s&admin_state=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("instance_type"); ok {
		res = fmt.Sprintf("%s&instance_type=%v", res, v)
	}
	if v, ok := d.GetOk("binding_service"); ok {
		res = fmt.Sprintf("%s&binding_service=%v", res, v)
	}
	if v, ok := d.GetOk("charge_mode"); ok {
		res = fmt.Sprintf("%s&charge_mode=%v", res, v)
	}
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsID)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func filterGCBByTags(tagsFilter map[string]interface{}, client *golangsdk.ServiceClient) ([]string, error) {
	var resourceIDs []string
	httpUrl := "v3/gcb/resource-instances/filter"
	basePath := client.Endpoint + httpUrl

	if len(tagsFilter) < 1 {
		return nil, nil
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildGCBByTagsOpts(tagsFilter),
	}

	offset := 0
	for {
		path := fmt.Sprintf("%s?limit=100&offset=%d", basePath, offset)
		resp, err := client.Request("POST", path, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		curJson := utils.PathSearch("resources[*].resource_id", respBody, make([]interface{}, 0))
		curArray := curJson.([]interface{})

		if len(curArray) == 0 {
			break
		}

		for _, item := range curArray {
			resourceIDs = append(resourceIDs, item.(string))
		}

		offset += 100
	}

	return resourceIDs, nil
}

func buildGCBByTagsOpts(tagmap map[string]interface{}) map[string]interface{} {
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

func filter(items []interface{}, resourceIDs []string) []interface{} {
	if len(resourceIDs) < 1 || len(items) < 1 {
		return nil
	}

	set := make(map[string]struct{})
	for _, id := range resourceIDs {
		set[id] = struct{}{}
	}

	var result []interface{}
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		if _, ok := set[itemMap["id"].(string)]; ok {
			result = append(result, item)
		}
	}

	return result
}

func flattenListGCBResponseBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"instances":             flattenGCBBindingInstances(v),
			"enable_share":          utils.PathSearch("enable_share", v, nil),
			"updated_at":            utils.PathSearch("updated_at", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"spec_code_id":          utils.PathSearch("spec_code_id", v, nil),
			"frozen":                utils.PathSearch("frozen", v, nil),
			"local_area":            utils.PathSearch("local_area", v, nil),
			"remote_area":           utils.PathSearch("remote_area", v, nil),
			"admin_state":           utils.PathSearch("admin_state", v, nil),
			"remote_site_code":      utils.PathSearch("remote_site_code", v, nil),
			"local_site_code":       utils.PathSearch("local_site_code", v, nil),
			"sla_level":             utils.PathSearch("sla_level", v, nil),
			"size":                  utils.PathSearch("size", v, nil),
			"charge_mode":           utils.PathSearch("charge_mode", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"binding_service":       utils.PathSearch("binding_service", v, nil),
			"type":                  utils.PathSearch("type", v, nil),
			"bordercross":           utils.PathSearch("bordercross", v, nil),
			"domain_id":             utils.PathSearch("domain_id", v, nil),
			"description":           utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func flattenGCBBindingInstances(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instances", resp, nil)

	if curJson == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	if curArray, ok := curJson.([]interface{}); ok {
		for _, item := range curArray {
			rst = append(rst, map[string]interface{}{
				"project_id": utils.PathSearch("project_id", item, nil),
				"region_id":  utils.PathSearch("region_id", item, nil),
				"type":       utils.PathSearch("type", item, nil),
				"id":         utils.PathSearch("id", item, nil),
			})
		}
	}
	return rst
}
