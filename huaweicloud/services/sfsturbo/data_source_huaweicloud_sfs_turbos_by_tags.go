package sfsturbo

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

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/resource_instances/action
func DataSourceSfsTurbosByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSfsTurbosByTagsRead,

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
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
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
				},
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
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
				},
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
							Elem: &schema.Resource{
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
							},
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSfsTurbosByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/sfs-turbo/resource_instances/action"
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	reqBodyParams := utils.RemoveNil(buildSfsTurbosBodyParams(d))
	action := d.Get("action").(string)
	resources := make([]interface{}, 0)
	offset := 0
	totalCount := 0

	for {
		if action == "filter" {
			reqBodyParams["offset"] = fmt.Sprintf("%v", offset)
		}

		listOpt.JSONBody = reqBodyParams
		listResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.FromErr(err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		if action == "count" {
			totalCount = int(utils.PathSearch("total_count", listRespBody, float64(0)).(float64))
			break
		}

		turbos := utils.PathSearch("resources", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(turbos) == 0 {
			totalCount = len(resources)
			break
		}

		resources = append(resources, turbos...)

		offset += len(turbos)
	}

	datasourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(datasourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resources", flattenSfsTurbosRespByTags(resources)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildSfsTurbosTagsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	v, ok := d.GetOk("tags")
	if !ok {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(v.([]interface{})))

	for i, tag := range v.([]interface{}) {
		bodyParams[i] = map[string]interface{}{
			"key":    utils.PathSearch("key", tag, nil),
			"values": utils.PathSearch("values", tag, nil),
		}
	}

	return bodyParams
}

func buildSfsTurbosMatchesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	v, ok := d.GetOk("matches")
	if !ok {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(v.([]interface{})))

	for i, match := range v.([]interface{}) {
		bodyParams[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", match, nil),
			"value": utils.PathSearch("value", match, nil),
		}
	}

	return bodyParams
}

func buildSfsTurbosBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":  d.Get("action"),
		"tags":    buildSfsTurbosTagsBodyParams(d),
		"matches": buildSfsTurbosMatchesBodyParams(d),
	}

	if d.Get("without_any_tag").(bool) {
		bodyParams["without_any_tag"] = true
	}

	return bodyParams
}

func flattenSfsTurbosRespByTags(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, len(resp))
	for i, v := range resp {
		rst[i] = map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_detail": utils.PathSearch("resource_detail", v, nil),
			"tags":            flattenResourceTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}

	return rst
}

func flattenResourceTags(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	tags := make([]interface{}, len(resp))
	for i, v := range resp {
		tags[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}

	return tags
}
