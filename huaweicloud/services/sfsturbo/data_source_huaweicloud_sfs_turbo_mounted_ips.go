package sfsturbo

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

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{share_id}/action
func DataSourceSfsTurboMountedIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSfsTurboMountedIpsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ips": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ips_attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildMountedIpsBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"ips": utils.ValueIgnoreEmpty(d.Get("ips")),
	}
	return map[string]interface{}{"get_client_ips": body}
}

func dataSourceSfsTurboMountedIpsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/action"
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{share_id}", d.Get("share_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildMountedIpsBodyParams(d)),
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving SFS Turbo mounted IP addresses: %s", err)
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ips_attribute", utils.PathSearch("ips", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
