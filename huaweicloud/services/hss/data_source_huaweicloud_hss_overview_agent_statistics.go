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

// @API HSS GET /v5/{project_id}/overview/agent/statistics
func DataSourceOverviewAgentStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOverviewAgentStatisticsRead,

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
			"container_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"wait_upgrade_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"online_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"not_online_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"offline_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"incluster_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"not_incluster_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildOverviewAgentStatisticsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	queryParams = fmt.Sprintf("%s&container_type=%v", queryParams, d.Get("container_type").(int))

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceOverviewAgentStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/overview/agent/statistics"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildOverviewAgentStatisticsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS overview agent statistics: %s", err)
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
		d.Set("wait_upgrade_num", utils.PathSearch("wait_upgrade_num", respBody, nil)),
		d.Set("online_num", utils.PathSearch("online_num", respBody, nil)),
		d.Set("not_online_num", utils.PathSearch("not_online_num", respBody, nil)),
		d.Set("offline_num", utils.PathSearch("offline_num", respBody, nil)),
		d.Set("incluster_num", utils.PathSearch("incluster_num", respBody, nil)),
		d.Set("not_incluster_num", utils.PathSearch("not_incluster_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
