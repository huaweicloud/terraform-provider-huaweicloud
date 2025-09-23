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

// @API HSS GET /v5/{project_id}/cluster/asset/statistics
func DataSourceClusterAssetStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterAssetStatisticsRead,
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
			// Attributes
			"cluster_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"work_load_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pod_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceClusterAssetStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/cluster/asset/statistics"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		requestPath = fmt.Sprintf("%s?enterprise_project_id=%v", requestPath, epsId)
	}
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS cluster asset statistics: %s", err)
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

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cluster_num", utils.PathSearch("cluster_num", respBody, nil)),
		d.Set("work_load_num", utils.PathSearch("work_load_num", respBody, nil)),
		d.Set("service_num", utils.PathSearch("service_num", respBody, nil)),
		d.Set("pod_num", utils.PathSearch("pod_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
