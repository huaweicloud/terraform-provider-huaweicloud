package organizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations POST /v1/organizations/{resource_type}/resource-instances/filter
func DataSourceResourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceInstancesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource type.`,
			},
			"without_any_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to only get the resources without tags.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Elem:        resourceInstancesTagSchema(),
				Optional:    true,
				Description: `The list of tags to be queried.`,
			},
			"matches": {
				Type:        schema.TypeList,
				Elem:        resourceInstancesMatchSchema(),
				Optional:    true,
				Description: `The fields to be queried.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Elem:        resourceInstancesResourceSchema(),
				Computed:    true,
				Description: `The list of resources that match the filter parameters.`,
			},
		},
	}
}

func resourceInstancesTagSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The key of the tag.`,
			},
			"values": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `The list of values of the tag.`,
			},
		},
	}
	return &sc
}

func resourceInstancesMatchSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The field name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The value corresponding to the field name.`,
			},
		},
	}
	return &sc
}

func resourceInstancesResourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource ID.`,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource name.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Elem:        resourceInstancesResourceTagSchema(),
				Computed:    true,
				Description: `The list of resource tags.`,
			},
		},
	}
	return &sc
}

func resourceInstancesResourceTagSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key of the tag.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the tag.`,
			},
		},
	}
	return &sc
}

func listResourceInstances(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/{resource_type}/resource-instances/filter"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{resource_type}", d.Get("resource_type").(string))
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildGetResourceInstancesBodyParams(d)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		listResp, err := client.Request("POST", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		resourceInstances := utils.PathSearch("resources", listRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, resourceInstances...)
		if len(resourceInstances) < limit {
			break
		}
		offset += len(resourceInstances)
	}

	return result, nil
}

func dataSourceResourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	resources, err := listResourceInstances(client, d)
	if err != nil {
		return diag.Errorf("error retrieving resource instances: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	return diag.FromErr(d.Set("resources", flattenResourceInstancesResp(resources)))
}

func buildGetResourceInstancesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"without_any_tag": utils.ValueIgnoreEmpty(d.Get("without_any_tag")),
		"tags":            buildGetResourceInstanceTagsBodyParams(d.Get("tags").([]interface{})),
		"matches":         parseResourceInstanceTags(d.Get("matches").([]interface{})),
	}
	return bodyParams
}

func buildGetResourceInstanceTagsBodyParams(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}
	res := make([]map[string]interface{}, 0, len(tags))
	for _, v := range tags {
		res = append(res, map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"values": utils.PathSearch("values", v, nil),
		})
	}
	return res
}

func flattenResourceInstancesResp(resources []interface{}) []interface{} {
	if len(resources) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		rst = append(rst, map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"tags":          parseResourceInstanceTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func parseResourceInstanceTags(tags []interface{}) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(tags))
	for _, v := range tags {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
