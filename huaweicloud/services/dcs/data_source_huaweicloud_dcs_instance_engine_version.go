package dcs

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

// @API DCS GET /v2/{project_id}/instances/{instance_id}/version
func DataSourceDcsInstanceEngineVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsInstanceEngineVersionRead,

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
			"engine_minor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_engine_minor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_minor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_proxy_minor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_minor_version_upgradable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"proxy_minor_version_upgradable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsInstanceEngineVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	httpUrl := "v2/{project_id}/instances/{instance_id}/version"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS instance engine version: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("engine_minor_version", utils.PathSearch("engine_minor_version", getRespBody, nil)),
		d.Set("latest_engine_minor_version", utils.PathSearch("latest_engine_minor_version", getRespBody, nil)),
		d.Set("proxy_minor_version", utils.PathSearch("proxy_minor_version", getRespBody, nil)),
		d.Set("latest_proxy_minor_version", utils.PathSearch("latest_proxy_minor_version", getRespBody, nil)),
		d.Set("engine_minor_version_upgradable", utils.PathSearch("engine_minor_version_upgradable", getRespBody, nil)),
		d.Set("proxy_minor_version_upgradable", utils.PathSearch("proxy_minor_version_upgradable", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
