package swrenterprise

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR POST /v2/{project_id}/{resource_type}/{resource_id}/{sub_resource_type}/resource-instances/filter
func DataSourceSwrEnterpriseSubResourcesFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseSubResourcesFilterRead,

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
				Description: `The type of the resource`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the resource`,
			},
			"sub_resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the sub resource`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The resource tags used to filter the target resources.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The key of the resource tag used to filter the target resources.`,
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The values corresponding to the current key used to filter the target resources.`,
						},
					},
				},
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of target resources that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource.`,
						},
						"resource_detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The detailed information of the resource, in JSON format.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource.`,
						},
						"tags":     common.TagsComputedSchema("The key/value tag pairs to associate with the resource."),
						"sys_tags": common.TagsComputedSchema("The key/value system tag pairs to associate with the resource."),
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total count of the resources.`,
			},
		},
	}
}

func dataSourceSwrEnterpriseSubResourcesFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	httpUrl := "v2/{project_id}/{resource_type}/{resource_id}/{sub_resource_type}/resource-instances/filter?limit=200"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{resource_type}", d.Get("resource_type").(string))
	listPath = strings.ReplaceAll(listPath, "{resource_id}", d.Get("resource_id").(string))
	listPath = strings.ReplaceAll(listPath, "{sub_resource_type}", d.Get("sub_resource_type").(string))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSwrEnterpriseSubResourcesFilterQueryParams(d)),
	}

	rst := make([]interface{}, 0)
	limit := 200
	offset := 0
	totalCount := float64(0)
	for {
		currentPath := listPath + fmt.Sprintf("&offset=%d", offset)
		listResp, err := client.Request("POST", currentPath, &listOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = utils.PathSearch("total_count", listRespBody, float64(0)).(float64)
		resources := utils.PathSearch("resources", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(resources) == 0 {
			break
		}

		rst = append(rst, resources...)

		offset += limit
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("total_count", totalCount),
		d.Set("resources", flattenSwrEnterpriseSubResourcesFilterResponseBody(rst)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildSwrEnterpriseSubResourcesFilterQueryParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"tags": buildSwrEnterpriseSubResourcesFilterQueryParamsTags(d.Get("tags")),
	}

	return bodyParams
}

func buildSwrEnterpriseSubResourcesFilterQueryParamsTags(tagsRaw interface{}) []map[string]interface{} {
	tags := tagsRaw.([]interface{})
	if len(tags) == 0 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		bodyParams[i] = map[string]interface{}{
			"key":    utils.PathSearch("key", tag, nil),
			"values": utils.PathSearch("values", tag, nil),
		}
	}

	return bodyParams
}

func flattenSwrEnterpriseSubResourcesFilterResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	resources := make([]interface{}, len(resp))
	for i, v := range resp {
		resources[i] = map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_detail": utils.JsonToString(utils.PathSearch("resource_detail", v, nil)),
			"tags":            utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"sys_tags":        utils.FlattenTagsToMap(utils.PathSearch("sys_tags", v, nil)),
		}
	}
	return resources
}
