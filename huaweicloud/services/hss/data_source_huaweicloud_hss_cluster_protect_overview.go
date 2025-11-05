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

// @API HSS GET /v5/{project_id}/cluster-protect/overview
func DataSourceClusterProtectOverview() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterProtectOverviewRead,

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
			"cluster_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"under_protect_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"event_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"block_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildClusterProtectOverviewQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceClusterProtectOverviewRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/cluster-protect/overview"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildClusterProtectOverviewQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS cluster protect overview: %s", err)
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
		d.Set("cluster_num", utils.PathSearch("cluster_num", respBody, nil)),
		d.Set("under_protect_num", utils.PathSearch("under_protect_num", respBody, nil)),
		d.Set("policy_num", utils.PathSearch("policy_num", respBody, nil)),
		d.Set("event_num", utils.PathSearch("event_num", respBody, nil)),
		d.Set("block_num", utils.PathSearch("block_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
