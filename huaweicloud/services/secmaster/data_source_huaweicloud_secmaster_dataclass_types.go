package secmaster

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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/dataclasses/{dataclass_id}/types
func DataSourceDataclassTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataclassTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataclass_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sortby": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_built_in": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"layout_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"level": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataclassTypeDataSchema(),
			},
		},
	}
}

func dataclassTypeDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataclass_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"layout_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"layout_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_category_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_built_in": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sla": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataclass_business_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildDataclassTypesQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	if v, ok := d.GetOk("enabled"); ok {
		queryParams = fmt.Sprintf("%s&enabled=%v", queryParams, v)
	}
	if v, ok := d.GetOk("order"); ok {
		queryParams = fmt.Sprintf("%s&order=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sortby"); ok {
		queryParams = fmt.Sprintf("%s&sortby=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sub_category"); ok {
		queryParams = fmt.Sprintf("%s&sub_category=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("category_code"); ok {
		queryParams = fmt.Sprintf("%s&category_code=%v", queryParams, v)
	}
	if v, ok := d.GetOk("is_built_in"); ok {
		queryParams = fmt.Sprintf("%s&is_built_in=%v", queryParams, v)
	}
	if v, ok := d.GetOk("layout_name"); ok {
		queryParams = fmt.Sprintf("%s&layout_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("level"); ok {
		queryParams = fmt.Sprintf("%s&level=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDataclassTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/dataclasses/{dataclass_id}/types"
		dataclassId = d.Get("dataclass_id").(string)
		offset      = 0
		limit       = 1000
		result      = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{dataclass_id}", dataclassId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildDataclassTypesQueryParams(d, limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster dataclass types: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("dataclass_type_details", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, dataResp...)

		if len(dataResp) < limit {
			break
		}

		offset += len(dataResp)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenDataclassTypesData(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataclassTypesData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"id":                      utils.PathSearch("id", v, nil),
			"dataclass_id":            utils.PathSearch("dataclass_id", v, nil),
			"domain_id":               utils.PathSearch("domain_id", v, nil),
			"project_id":              utils.PathSearch("project_id", v, nil),
			"workspace_id":            utils.PathSearch("workspace_id", v, nil),
			"region_id":               utils.PathSearch("region_id", v, nil),
			"layout_id":               utils.PathSearch("layout_id", v, nil),
			"layout_name":             utils.PathSearch("layout_name", v, nil),
			"category":                utils.PathSearch("category", v, nil),
			"category_code":           utils.PathSearch("category_code", v, nil),
			"sub_category":            utils.PathSearch("sub_category", v, nil),
			"sub_category_code":       utils.PathSearch("sub_category_code", v, nil),
			"description":             utils.PathSearch("description", v, nil),
			"enabled":                 utils.PathSearch("enabled", v, nil),
			"level":                   utils.PathSearch("level", v, nil),
			"is_built_in":             utils.PathSearch("is_built_in", v, nil),
			"sla":                     utils.PathSearch("sla", v, nil),
			"creator_id":              utils.PathSearch("creator_id", v, nil),
			"creator_name":            utils.PathSearch("creator_name", v, nil),
			"modifier_id":             utils.PathSearch("modifier_id", v, nil),
			"modifier_name":           utils.PathSearch("modifier_name", v, nil),
			"create_time":             utils.PathSearch("create_time", v, nil),
			"update_time":             utils.PathSearch("update_time", v, nil),
			"dataclass_business_code": utils.PathSearch("dataclass_business_code", v, nil),
			"sub_count":               utils.PathSearch("sub_count", v, nil),
		})
	}

	return rst
}
