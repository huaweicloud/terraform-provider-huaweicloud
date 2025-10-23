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

// @API HSS GET /v5/{project_id}/container/cluster/statistics
func DataSourceContainerClusterStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerClusterStatisticsRead,

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
			"risk_cluster_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"app_vul_cluster_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unscan_cluster_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"all_cluster_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildContainerClusterStatisticsQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceContainerClusterStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/container/cluster/statistics"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerClusterStatisticsQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS container cluster statistics: %s", err)
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
		d.Set("risk_cluster_num", utils.PathSearch("risk_cluster_num", respBody, nil)),
		d.Set("app_vul_cluster_num", utils.PathSearch("app_vul_cluster_num", respBody, nil)),
		d.Set("unscan_cluster_num", utils.PathSearch("unscan_cluster_num", respBody, nil)),
		d.Set("all_cluster_num", utils.PathSearch("all_cluster_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
