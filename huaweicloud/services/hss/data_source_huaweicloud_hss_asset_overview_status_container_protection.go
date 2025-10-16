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

// @API HSS GET /v5/{project_id}/asset/overview/status/container/protection
func DataSourceAssetOverviewStatusContainerProtection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetOverviewStatusContainerProtectionRead,

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
			"no_risk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"risk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"no_protect": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildAssetOverviewStatusContainerProtectionQueryParams(epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return queryParams
}

func dataSourceAssetOverviewStatusContainerProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/asset/overview/status/container/protection"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAssetOverviewStatusContainerProtectionQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS asset overview status container protection: %s", err)
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
		d.Set("no_risk", utils.PathSearch("no_risk", respBody, nil)),
		d.Set("risk", utils.PathSearch("risk", respBody, nil)),
		d.Set("no_protect", utils.PathSearch("no_protect", respBody, nil)),
		d.Set("total_num", utils.PathSearch("total_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
