package vpn

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

// @API VPN POST /v5/{project_id}/{resource_type}/resource-instances/filter
func DataSourceVpnInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource type.`,
			},
			"without_any_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to filter instances.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tag list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the tag key.`,
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the value list of the tag.`,
						},
					},
				},
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the search field, including a key and a value.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the match key.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the match value.`,
						},
					},
				},
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the resource object list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the resource ID.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the resource name.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the tag list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the tag key.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the tag value.`,
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

func dataSourceVpnInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	httpUrl := "v5/{project_id}/{resource_type}/resource-instances/filter"
	listBasePath := client.Endpoint + httpUrl
	listBasePath = strings.ReplaceAll(listBasePath, "{project_id}", client.ProjectID)
	listBasePath = strings.ReplaceAll(listBasePath, "{resource_type}", d.Get("resource_type").(string))
	listsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildGetResourceInstancesBodyParams(d)),
	}

	var res []interface{}
	limit := 100
	offset := 0
	for {
		listPath := listBasePath + buildListResourceInstancesQueryParams(limit, offset)
		listResp, err := client.Request("POST", listPath, &listsOpt)
		if err != nil {
			return diag.Errorf("error retrieving VPN resource instances: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceInstances := flattenResourceInstancesResp(listRespBody)
		if len(resourceInstances) == 0 {
			break
		}
		res = append(res, resourceInstances...)
		offset += len(resourceInstances)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListResourceInstancesQueryParams(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func buildGetResourceInstancesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"without_any_tag": utils.ValueIgnoreEmpty(d.Get("without_any_tag")),
		"tags":            buildGetResourceInstanceTagsBodyParams(d),
		"matches":         buildGetResourceInstanceMatchesBodyParams(d),
	}
	return bodyParams
}

func buildGetResourceInstanceTagsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	tagsRawParams := d.Get("tags").([]interface{})
	if len(tagsRawParams) == 0 {
		return nil
	}
	res := make([]map[string]interface{}, 0, len(tagsRawParams))
	for _, v := range tagsRawParams {
		raw := v.(map[string]interface{})
		res = append(res, map[string]interface{}{
			"key":    raw["key"],
			"values": raw["values"],
		})
	}
	return res
}

func buildGetResourceInstanceMatchesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	matchesRawParams := d.Get("matches").([]interface{})
	if len(matchesRawParams) == 0 {
		return nil
	}
	res := make([]map[string]interface{}, 0, len(matchesRawParams))
	for _, v := range matchesRawParams {
		raw := v.(map[string]interface{})
		res = append(res, map[string]interface{}{
			"key":   raw["key"],
			"value": raw["value"],
		})
	}
	return res
}

func flattenResourceInstancesResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"tags":          flattenResourceInstanceTagsResp(v),
		})
	}
	return rst
}

func flattenResourceInstanceTagsResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("tags", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
