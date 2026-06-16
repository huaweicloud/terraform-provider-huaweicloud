package drs

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

// @API DRS POST /v5/{project_id}/{resource_type}/resource-instances/filter
func DataSourceDrsInstancesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceDrsInstancesByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
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
				Optional: true,
				Elem:     instancesByTagsTagsSchema(),
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     instancesByTagsMatchesSchema(),
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instancesByTagsResourcesSchema(),
			},
		},
	}
}

func instancesByTagsMatchesSchema() *schema.Resource {
	return &schema.Resource{
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
}

func instancesByTagsTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func instancesByTagsResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instancesByTagsResourcesTagsSchema(),
			},
		},
	}
}

func instancesByTagsResourcesTagsSchema() *schema.Resource {
	return &schema.Resource{
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
}

func buildInstancesByTagsTagsBodyParams(rawArray []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, raw := range rawArray {
		rawMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":    rawMap["key"],
			"values": rawMap["values"],
		})
	}

	return rst
}

func buildInstancesByTagsMatchesBodyParams(rawArray []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, raw := range rawArray {
		rawMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":   rawMap["key"],
			"value": rawMap["value"],
		})
	}

	return rst
}

func buildListResourceInstancesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"without_any_tag": d.Get("without_any_tag"),
		"tags":            buildInstancesByTagsTagsBodyParams(d.Get("tags").([]interface{})),
		"matches":         buildInstancesByTagsMatchesBodyParams(d.Get("matches").([]interface{})),
	}
}

func DataSourceDrsInstancesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/{resource_type}/resource-instances/filter"
		result  = make([]interface{}, 0)
		limit   = 1000
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{resource_type}", d.Get("resource_type").(string))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildListResourceInstancesBodyParams(d)),
	}

	for {
		queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
		currentListPath := listPath + queryParams

		listResp, err := client.Request("POST", currentListPath, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS resource instances: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		resources := utils.PathSearch("resources", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(resources) == 0 {
			break
		}

		result = append(result, resources...)

		offset += len(resources)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resources", flattenResources(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResources(resourcesResp []interface{}) []interface{} {
	if len(resourcesResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resourcesResp))
	for _, v := range resourcesResp {
		rst = append(rst, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_detail": utils.PathSearch("resource_detail", v, nil),
			"tags":            flattenResourceTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenResourceTags(tagsResp []interface{}) []interface{} {
	if len(tagsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(tagsResp))
	for _, v := range tagsResp {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
