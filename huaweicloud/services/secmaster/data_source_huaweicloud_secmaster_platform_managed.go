package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/siem/cloud-logs/managers
func DataSourcePlatformManaged() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlatformManagedRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dw_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform_managed_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publish_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_managed_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"whitelist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePlatformManagedRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/siem/cloud-logs/managers"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster platform managed information: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("create_time", utils.PathSearch("result[0].create_time", respBody, nil)),
		d.Set("dw_region", utils.PathSearch("result[0].dw_region", respBody, nil)),
		d.Set("platform_managed_domain_id", utils.PathSearch(
			"result[0].platform_managed_domain_id", respBody, nil)),
		d.Set("publish_status", utils.PathSearch("result[0].publish_status", respBody, nil)),
		d.Set("tenant_managed_domain_id", utils.PathSearch(
			"result[0].tenant_managed_domain_id", respBody, nil)),
		d.Set("update_time", utils.PathSearch("result[0].update_time", respBody, nil)),
		d.Set("whitelist", utils.PathSearch("result[0].whitelist", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
