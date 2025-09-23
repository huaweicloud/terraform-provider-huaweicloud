package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/tenant/profile
func DataSourceAppService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppServiceRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the APP service is located.`,
			},
			"service_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the APP service.`,
			},
			"open_with_ad": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the APP service is connected to AD.`,
			},
			"tenant_domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain ID to which the APP service belongs.`,
			},
			"tenant_domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name to which the APP service belongs.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the Workspace service, in RFC3339 format.`,
			},
		},
	}
}

func dataSourceAppServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/tenant/profile"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return diag.FromErr(err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.Errorf("error getting APP service information: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("service_status", utils.PathSearch("service_status", respBody, nil)),
		d.Set("open_with_ad", utils.PathSearch("open_with_ad", respBody, nil)),
		d.Set("tenant_domain_id", utils.PathSearch("tenant_domain_id", respBody, nil)),
		d.Set("tenant_domain_name", utils.PathSearch("tenant_domain_name", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
			respBody, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
