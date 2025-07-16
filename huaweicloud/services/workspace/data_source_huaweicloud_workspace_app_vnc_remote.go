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

// @API Workspace GET /v1/{project_id}/app-servers/{server_id}/actions/vnc
func DataSourceAppVncRemote() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppVncRemoteRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the APP server is located.`,
			},
			// Required parameters.
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the APP server to get VNC remote information.`,
			},
			// Attributes.
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The remote login console address.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The login type.`,
			},
			// Internal parameters.
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The login protocol.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func getAppVncRemote(client *golangsdk.ServiceClient, serverId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/app-servers/{server_id}/actions/vnc"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_id}", serverId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func dataSourceAppVncRemoteRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	serverId := d.Get("server_id").(string)
	resp, err := getAppVncRemote(client, serverId)
	if err != nil {
		return diag.Errorf("error getting VNC remote information: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("url", utils.PathSearch("url", resp, nil)),
		d.Set("type", utils.PathSearch("type", resp, nil)),
		d.Set("protocol", utils.PathSearch("protocol", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
