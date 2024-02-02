package ecs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/remote_console
func DataSourceComputeInstanceRemoteConsole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeInstanceRemoteConsoleRead,

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

			// attributes
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeInstanceRemoteConsoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Get("instance_id").(string)
	product := "ecs"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"remote_console": utils.RemoveNil(map[string]interface{}{
				"protocol": "vnc",
				"type":     "novnc",
			}),
		},
	}

	getVNCRemoteAddressHttpUrl := "v1/{project_id}/cloudservers/{server_id}/remote_console"
	getVNCRemoteAddressPath := client.Endpoint + getVNCRemoteAddressHttpUrl
	getVNCRemoteAddressPath = strings.ReplaceAll(getVNCRemoteAddressPath, "{project_id}", client.ProjectID)
	getVNCRemoteAddressPath = strings.ReplaceAll(getVNCRemoteAddressPath, "{server_id}", instanceId)
	getVNCRemoteAddressResp, err := client.Request("POST", getVNCRemoteAddressPath, &requestOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VNC Remote Address")
	}
	getVNCRemoteAddressBody, err := utils.FlattenResponse(getVNCRemoteAddressResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceId)
	remoteConsole := utils.PathSearch("remote_console", getVNCRemoteAddressBody, nil)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("protocol", utils.PathSearch("url", remoteConsole, nil)),
		d.Set("type", utils.PathSearch("url", remoteConsole, nil)),
		d.Set("url", utils.PathSearch("url", remoteConsole, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
