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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/workbenches
func DataSourceWorkbenches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkbenchesRead,

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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"global_search_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     workbenchDataSchema(),
			},
		},
	}
}

func workbenchDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url_openwith_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"icon": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"basic_properties": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_id": {
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
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_favorite": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildWorkbenchesQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("creator_type"); ok {
		queryParams = fmt.Sprintf("%s&creator_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("global_search_text"); ok {
		queryParams = fmt.Sprintf("%s&global_search_text=%v", queryParams, v)
	}
	if v, ok := d.GetOk("tags"); ok {
		queryParams = fmt.Sprintf("%s&tags=%v", queryParams, v)
	}
	if v, ok := d.GetOk("from_date"); ok {
		queryParams = fmt.Sprintf("%s&from_date=%v", queryParams, v)
	}
	if v, ok := d.GetOk("to_date"); ok {
		queryParams = fmt.Sprintf("%s&to_date=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceWorkbenchesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/workbenches"
		offset  = 0
		limit   = 200
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildWorkbenchesQueryParams(d, limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster workbenches: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
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
		d.Set("data", flattenWorkbenchesData(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWorkbenchesData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"type":              utils.PathSearch("type", v, nil),
			"url":               utils.PathSearch("url", v, nil),
			"url_openwith_type": utils.PathSearch("url_openwith_type", v, nil),
			"tags":              utils.PathSearch("tags", v, nil),
			"description":       utils.PathSearch("description", v, nil),
			"icon":              utils.PathSearch("icon", v, nil),
			"basic_properties":  utils.PathSearch("basic_properties", v, nil),
			"domain_id":         utils.PathSearch("domain_id", v, nil),
			"region_id":         utils.PathSearch("region_id", v, nil),
			"workspace_id":      utils.PathSearch("workspace_id", v, nil),
			"create_time":       utils.PathSearch("create_time", v, nil),
			"update_time":       utils.PathSearch("update_time", v, nil),
			"creator_id":        utils.PathSearch("creator_id", v, nil),
			"creator_name":      utils.PathSearch("creator_name", v, nil),
			"modifier_id":       utils.PathSearch("modifier_id", v, nil),
			"modifier_name":     utils.PathSearch("modifier_name", v, nil),
			"is_deleted":        utils.PathSearch("is_deleted", v, nil),
			"is_favorite":       utils.PathSearch("is_favorite", v, nil),
			"status":            utils.PathSearch("status", v, nil),
		})
	}

	return rst
}
