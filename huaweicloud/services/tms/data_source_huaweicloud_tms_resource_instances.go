package tms

import (
	"context"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TMS POST /v1.0/resource-instances/filter
func DataSourceResourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsResourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"resource_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     tmsTagsSchema(),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"without_any_tag": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tmsResourcesSchema(),
			},
			"errors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tmsErrorsSchema(),
			},
		},
	}
}

func tmsTagsSchema() *schema.Resource {
	sc := schema.Resource{
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
	return &sc
}

func tmsResourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tmsResourcesTagsSchema(),
			},
		},
	}
	return &sc
}

func tmsResourcesTagsSchema() *schema.Resource {
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

func tmsErrorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceTmsResourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1.0/resource-instances/filter"
		product = "tms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	limit := 200
	offset := 0
	resourcesRes := make([]interface{}, 0)
	errorsRes := make([]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildGetTmsResourceInstancesBodyParams(d, limit, offset))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving TMS resource instances: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		resources := flattenTmsResourceInstancesResources(getRespBody)
		errors := flattenTmsResourceInstancesErrors(getRespBody)
		resourcesRes = append(resourcesRes, resources...)
		errorsRes = append(errorsRes, errors...)
		if len(resources) == 0 && len(errors) == 0 {
			break
		}
		offset += limit
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("resources", resourcesRes),
		d.Set("errors", errorsRes),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetTmsResourceInstancesBodyParams(d *schema.ResourceData, limit, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_id":     utils.ValueIgnoreEmpty(d.Get("project_id")),
		"resource_types": d.Get("resource_types"),
		"tags":           buildGetTmsResourceInstancesTagsBody(d),
		"limit":          limit,
		"offset":         offset,
	}
	if v, ok := d.GetOk("without_any_tag"); ok {
		withoutAnyTag, _ := strconv.ParseBool(v.(string))
		bodyParams["without_any_tag"] = withoutAnyTag
	}
	return bodyParams
}

func buildGetTmsResourceInstancesTagsBody(d *schema.ResourceData) []map[string]interface{} {
	rawTags := d.Get("tags").([]interface{})

	rst := make([]map[string]interface{}, 0, len(rawTags))
	for _, rawTag := range rawTags {
		if v, ok := rawTag.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"key":    v["key"],
				"values": v["values"],
			})
		}
	}
	return rst
}

func flattenTmsResourceInstancesResources(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"project_id":    utils.PathSearch("project_id", v, nil),
			"project_name":  utils.PathSearch("project_name", v, nil),
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
			"tags":          flattenTmsResourceInstancesResourceTags(v),
		})
	}
	return rst
}

func flattenTmsResourceInstancesResourceTags(resp interface{}) []interface{} {
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

func flattenTmsResourceInstancesErrors(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("errors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"error_code":    utils.PathSearch("error_code", v, nil),
			"error_msg":     utils.PathSearch("error_msg", v, nil),
			"project_id":    utils.PathSearch("project_id", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
		})
	}
	return rst
}
