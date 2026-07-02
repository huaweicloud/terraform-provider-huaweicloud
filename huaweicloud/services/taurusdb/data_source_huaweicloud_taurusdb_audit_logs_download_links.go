package taurusdb

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

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/audit-log-link
func DataSourceTaurusDBAuditLogsDownloadLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBAuditLogsDownloadLinksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTaurusDBAuditLogsDownloadLinksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/audit-log-link"
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	postPath := client.Endpoint + httpUrl
	postPath = strings.ReplaceAll(postPath, "{project_id}", client.ProjectID)
	postPath = strings.ReplaceAll(postPath, "{instance_id}", d.Get("instance_id").(string))

	postOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"ids": utils.ExpandToStringList(d.Get("ids").([]interface{})),
		},
	}

	postResp, err := client.Request("POST", postPath, &postOpt)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB audit logs download links: %s", err)
	}

	postRespBody, err := utils.FlattenResponse(postResp)
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
		d.Set("links", utils.ExpandToStringList(
			utils.PathSearch("links", postRespBody, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
