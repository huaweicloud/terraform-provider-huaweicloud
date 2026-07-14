package dsc

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

// @API DSC GET /v1/{project_id}/scan-security-levels
func DataSourceDscScanSecurityLevels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscSecurityLevelsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_deleted": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_levels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscSecurityLevelsSchema(),
			},
		},
	}
}

func dscSecurityLevelsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"level_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_level_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"color_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_level_desc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sort_weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildDscSecurityLevelsQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%v", queryParams, v)
	}
	if v, ok := d.GetOk("is_deleted"); ok {
		queryParams = fmt.Sprintf("%s&is_deleted=%v", queryParams, v)
	}
	return queryParams
}

func dataSourceDscSecurityLevelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/scan-security-levels"
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	for {
		currentPath := requestPath + buildDscSecurityLevelsQueryParams(d, limit, offset)
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC security levels: %s", err)
		}

		requestRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		securityLevelsList := utils.PathSearch("security_levels_list", requestRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, securityLevelsList...)
		if len(securityLevelsList) < limit {
			break
		}
		offset += len(securityLevelsList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("security_levels", flattenDscSecurityLevels(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscSecurityLevels(securityLevelsList []interface{}) []interface{} {
	if len(securityLevelsList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(securityLevelsList))
	for _, v := range securityLevelsList {
		rst = append(rst, map[string]interface{}{
			"level_id":            utils.PathSearch("level_id", v, nil),
			"project_id":          utils.PathSearch("project_id", v, nil),
			"security_level_name": utils.PathSearch("security_level_name", v, nil),
			"color_number":        utils.PathSearch("color_number", v, nil),
			"security_level_desc": utils.PathSearch("security_level_desc", v, nil),
			"used_count":          utils.PathSearch("used_count", v, nil),
			"category":            utils.PathSearch("category", v, nil),
			"create_time":         utils.PathSearch("create_time", v, nil),
			"update_time":         utils.PathSearch("update_time", v, nil),
			"sort_weight":         utils.PathSearch("sort_weight", v, nil),
			"is_deleted":          utils.PathSearch("is_deleted", v, nil),
		})
	}

	return rst
}
