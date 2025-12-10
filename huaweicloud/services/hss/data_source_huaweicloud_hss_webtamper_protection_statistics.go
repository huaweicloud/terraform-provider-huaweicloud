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

// @API HSS GET /v5/{project_id}/wtp/statistics
func DataSourceWebtamperProtectionStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebtamperProtectionStatisticsRead,

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
			"protect_host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protect_success_host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protect_fail_host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"anti_tampering_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceWebtamperProtectionStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/wtp/statistics"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		getPath = fmt.Sprintf("%s?enterprise_project_id=%s", getPath, epsId)
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving protection data statistics: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("protect_host_num", utils.PathSearch("protect_host_num", respBody, nil)),
		d.Set("protect_success_host_num", utils.PathSearch("protect_success_host_num", respBody, nil)),
		d.Set("protect_fail_host_num", utils.PathSearch("protect_fail_host_num", respBody, nil)),
		d.Set("anti_tampering_num", utils.PathSearch("anti_tampering_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
