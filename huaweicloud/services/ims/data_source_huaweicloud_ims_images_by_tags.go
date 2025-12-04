package ims

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

// @API IMS POST /v2/{project_id}/images/resource_instances/action
func DataSourceIMSImagesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIMSImagesByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"without_any_tag": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsMatchesSchema(),
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsTagsSchema(),
			},
			"not_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsTagsSchema(),
			},
			"tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsTagsSchema(),
			},
			"not_tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourcesByTagsTagsSchema(),
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIMSImagesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("ims", region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	httpUrl := "v2/{project_id}/images/resource_instances/action"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	bodyParams := buildIMSImagesByTagsQueryParams(d)

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

func buildIMSImagesByTagsQueryParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":       d.Get("action"),
		"matches":      buildNIMSImagesByTagsMatchesBodyParams(d.Get("matches")),
		"tags":         buildIMSImagesByTagsTagsQueryParams(d.Get("tags")),
		"not_tags":     buildIMSImagesByTagsTagsQueryParams(d.Get("not_tags")),
		"tags_any":     buildIMSImagesByTagsTagsQueryParams(d.Get("tags_any")),
		"not_tags_any": buildIMSImagesByTagsTagsQueryParams(d.Get("not_tags_any")),
	}
	if v, ok := d.GetOk("without_any_tag"); ok {
		withoutAnyTag, _ := strconv.ParseBool(v.(string))
		bodyParams["without_any_tag"] = withoutAnyTag
	}

	return bodyParams
}

func buildNIMSImagesByTagsMatchesBodyParams(matchesRaw interface{}) []map[string]interface{} {
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

func buildIMSImagesByTagsTagsQueryParams(tagsRaw interface{}) []map[string]interface{} {
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

func flattenResourcesByTagsResourcesResponseBody(resp []interface{}) []interface{} {
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
		}
	}
	return resources
}
