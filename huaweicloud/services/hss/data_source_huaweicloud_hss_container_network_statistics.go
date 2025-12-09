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

// @API HSS GET /v5/{project_id}/container-network/network-statistics
func DataSourceContainerNetworkStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerNetworkStatisticsRead,

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
			"protected_cluster_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"namespace_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"network_policy_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildContainerNetworkStatisticsQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func dataSourceContainerNetworkStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/container-network/network-statistics"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerNetworkStatisticsQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS container network statistics: %s", err)
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
		d.Set("protected_cluster_total_num", utils.PathSearch("protected_cluster_total_num", respBody, nil)),
		d.Set("cluster_total_num", utils.PathSearch("cluster_total_num", respBody, nil)),
		d.Set("namespace_total_num", utils.PathSearch("namespace_total_num", respBody, nil)),
		d.Set("network_policy_total_num", utils.PathSearch("network_policy_total_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
