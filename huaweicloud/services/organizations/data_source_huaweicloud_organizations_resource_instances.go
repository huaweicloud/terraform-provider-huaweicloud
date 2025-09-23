package organizations

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

// @API Organizations POST /v1/organizations/{resource_type}/resource-instances/filter
func DataSourceOrganizationsResourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationsResourceInstancesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Elem:     organizationsResourceInstancesTagSchema(),
				Optional: true,
			},
			"matches": {
				Type:     schema.TypeList,
				Elem:     organizationsResourceInstancesMatchSchema(),
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Elem:     organizationsResourceInstancesResourceSchema(),
				Computed: true,
			},
		},
	}
}

func organizationsResourceInstancesTagSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
	return &sc
}

func organizationsResourceInstancesMatchSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func organizationsResourceInstancesResourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Elem:     organizationsResourceInstancesResourceTagSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func organizationsResourceInstancesResourceTagSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceOrganizationsResourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/organizations/{resource_type}/resource-instances/filter"
		product = "organizations"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	listBasePath := client.Endpoint + httpUrl
	listBasePath = strings.ReplaceAll(listBasePath, "{resource_type}", d.Get("resource_type").(string))
	listsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listsOpt.JSONBody = utils.RemoveNil(buildGetResourceInstancesBodyParams(d))

	var res []interface{}
	limit := 100
	offset := 0

	for {
		listPath := listBasePath + buildListResourceInstancesQueryParams(limit, offset)
		listResp, err := client.Request("POST", listPath, &listsOpt)
		if err != nil {
			return diag.Errorf("error retrieving Organizations resource instances: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceInstances := flattenResourceInstancesResp(listRespBody)
		res = append(res, resourceInstances...)
		offset += len(resourceInstances)
		totalCount := utils.PathSearch("total_count", listRespBody, float64(0)).(float64)
		if offset >= int(totalCount) {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
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
