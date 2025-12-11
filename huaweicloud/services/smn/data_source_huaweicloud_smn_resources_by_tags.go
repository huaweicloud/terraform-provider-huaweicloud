package smn

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SMN POST /v2/{project_id}/{resource_type}/resource_instances/action
func DataSourceResourcesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcesByTagsRead,

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
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsTagsSchema(),
			},
			"tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsTagsSchema(),
			},
			"not_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsTagsSchema(),
			},
			"not_tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsTagsSchema(),
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsMatchesSchema(),
			},
			"without_any_tag": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourcesByTagsResourcesSchema(),
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcesByTagsTagsSchema() *schema.Resource {
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

func resourcesByTagsMatchesSchema() *schema.Resource {
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

func resourcesByTagsResourcesSchema() *schema.Resource {
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
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourcesByTagsResourceDetailSchema(),
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourcesByTagsResourceTagSchema(),
			},
		},
	}
}

func resourcesByTagsResourceDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"detail_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcesByTagsResourceTagSchema() *schema.Resource {
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

func dataSourceResourcesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v2/{project_id}/{resource_type}/resource_instances/action"
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{resource_type}", d.Get("resource_type").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	bodyParams := buildResourcesByTagsQueryParams(d)

	action := d.Get("action").(string)
	rst := make([]interface{}, 0)
	limit := 1000
	offset := 0
	totalCount := float64(0)
	for {
		if action == "filter" {
			bodyParams["limit"] = strconv.Itoa(limit)
			bodyParams["offset"] = strconv.Itoa(offset)
		}

		listOpt.JSONBody = utils.RemoveNil(bodyParams)
		listResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		if action == "count" {
			totalCount = utils.PathSearch("total_count", listRespBody, float64(0)).(float64)
			break
		}

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
		d.Set("resources", flattenResourcesByTagsResourcesResponseBody(rst)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildResourcesByTagsQueryParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":       d.Get("action"),
		"tags":         buildResourcesByTagsTagsQueryParams(d.Get("tags")),
		"tags_any":     buildResourcesByTagsTagsQueryParams(d.Get("tags_any")),
		"not_tags":     buildResourcesByTagsTagsQueryParams(d.Get("not_tags")),
		"not_tags_any": buildResourcesByTagsTagsQueryParams(d.Get("not_tags_any")),
		"matches":      buildResourcesByTagsMatchesBodyParams(d.Get("matches")),
	}
	if v, ok := d.GetOk("without_any_tag"); ok {
		withoutAnyTag, _ := strconv.ParseBool(v.(string))
		bodyParams["without_any_tag"] = withoutAnyTag
	}

	return bodyParams
}

func buildResourcesByTagsTagsQueryParams(tagsRaw interface{}) []map[string]interface{} {
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

func buildResourcesByTagsMatchesBodyParams(matchesRaw interface{}) []map[string]interface{} {
	matches := matchesRaw.([]interface{})
	if len(matches) == 0 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(matches))

	for i, match := range matches {
		bodyParams[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", match, nil),
			"value": utils.PathSearch("value", match, nil),
		}
	}

	return bodyParams
}

func flattenResourcesByTagsResourcesResponseBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_detail": flattenResourcesByTagsResourceDetailResponseBody(v),
			"tags":            flattenResourcesByTagsResourceTagResponseBody(v),
		})
	}
	return rst
}

func flattenResourcesByTagsResourceDetailResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("resource_detail", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"detail_id":             utils.PathSearch("detailId", curJson, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", curJson, nil),
			"topic_urn":             utils.PathSearch("topic_urn", curJson, nil),
			"display_name":          utils.PathSearch("display_name", curJson, nil),
		},
	}
	return rst
}

func flattenResourcesByTagsResourceTagResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("tags", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		},
		)
	}
	return rst
}
