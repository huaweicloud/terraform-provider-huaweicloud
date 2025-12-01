package hss

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

// @API HSS GET /v5/{project_id}/security-check/config
func DataSourceSecurityCheckConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityCheckConfigRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"check_period_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"day_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"week_period": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hour": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host_id_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildSecurityCheckConfigQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceSecurityCheckConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/security-check/config"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildSecurityCheckConfigQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS security check config: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("task_id", utils.PathSearch("task_id", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("check_period_type", utils.PathSearch("check_period_type", respBody, nil)),
		d.Set("day_period", utils.PathSearch("day_period", respBody, nil)),
		d.Set("week_period", utils.ExpandToStringList(
			utils.PathSearch("week_period", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("hour", utils.PathSearch("hour", respBody, nil)),
		d.Set("content", utils.ExpandToStringList(
			utils.PathSearch("content", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("host_id_list", utils.ExpandToStringList(
			utils.PathSearch("host_id_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
