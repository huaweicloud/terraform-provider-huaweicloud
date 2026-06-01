package geminidb

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

// @API GeminiDB GET /v3/{project_id}/ip-num-requirement
func DataSourceIpNumRequirement() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIpNumRequirementRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildIpNumRequirementQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?node_num=%v", d.Get("node_num"))

	if v, ok := d.GetOk("engine_name"); ok {
		queryParams = fmt.Sprintf("%s&engine_name=%v", queryParams, v)
	}

	if v, ok := d.GetOk("instance_mode"); ok {
		queryParams = fmt.Sprintf("%s&instance_mode=%v", queryParams, v)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		queryParams = fmt.Sprintf("%s&instance_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceIpNumRequirementRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/ip-num-requirement"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildIpNumRequirementQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the number of IP addresses: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ip_address_count", utils.PathSearch("count", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
