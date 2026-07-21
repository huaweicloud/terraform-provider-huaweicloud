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

// @API DSC GET /v1/{project_id}/scan-jobs/{job_id}/security-levels
func DataSourceDscSecurityLevels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscJobSecurityLevelsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the scan job ID.",
			},
			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the keyword for fuzzy search on object names.",
			},
			"asset_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the asset type for filtering.",
			},
			"asset_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the asset ID for filtering.",
			},
			"security_level_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the security level IDs for filtering.",
			},
			"security_level_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The security level distribution list.",
				Elem:        dscSecurityLevelInfoSchema(),
			},
		},
	}
}

func dscSecurityLevelInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"security_level_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security level ID.",
			},
			"security_level_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security level name.",
			},
			"security_level_color": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The security level color.",
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of risk objects under this security level.",
			},
			"percent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The percentage of this security level among all risks.",
			},
		},
	}
}

func buildDscJobSecurityLevelsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("keyword"); ok {
		queryParams = fmt.Sprintf("%s&keyword=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_type"); ok {
		queryParams = fmt.Sprintf("%s&asset_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_id"); ok {
		queryParams = fmt.Sprintf("%s&asset_id=%v", queryParams, v)
	}
	securityLevelIds := d.Get("security_level_ids").([]interface{})
	for _, id := range securityLevelIds {
		queryParams = fmt.Sprintf("%s&security_level_ids=%v", queryParams, id)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceDscJobSecurityLevelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		jobId   = d.Get("job_id").(string)
		httpUrl = "v1/{project_id}/scan-jobs/{job_id}/security-levels"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	currentPath := requestPath + buildDscJobSecurityLevelsQueryParams(d)

	resp, err := client.Request("GET", currentPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC security levels: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("security_level_list", flattenDscJobSecurityLevels(
			utils.PathSearch("security_level_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscJobSecurityLevels(levels []interface{}) []interface{} {
	if len(levels) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(levels))
	for _, v := range levels {
		rst = append(rst, map[string]interface{}{
			"security_level_id":    utils.PathSearch("security_level_id", v, nil),
			"security_level_name":  utils.PathSearch("security_level_name", v, nil),
			"security_level_color": utils.PathSearch("security_level_color", v, nil),
			"count":                utils.PathSearch("count", v, nil),
			"percent":              utils.PathSearch("percent", v, nil),
		})
	}

	return rst
}
