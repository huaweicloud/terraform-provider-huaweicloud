package dws

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

// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/db-manager/om-user/status
func DataSourceOmAccountConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOmAccountConfigurationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DWS cluster ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the OM account.`,
			},
			"om_user_expires_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time of the OM account, in RFC3339 format.`,
			},
		},
	}
}

func dataSourceOmAccountConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/clusters/{cluster_id}/db-manager/om-user/status"
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DWS OM account configuration: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	expirationTime := utils.PathSearch("om_user_info.om_user_expires_time", getRespBody, float64(0))
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("status", utils.PathSearch("om_user_info.om_user_status", getRespBody, nil)),
		d.Set("om_user_expires_time", utils.FormatTimeStampRFC3339(int64(expirationTime.(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
