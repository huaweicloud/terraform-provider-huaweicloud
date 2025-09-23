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

// @API Workspace GET /v2/{project_id}/desktops/{desktop_id}/remote-consoles
func DataSourceDesktopRemoteConsole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDesktopRemoteConsoleRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the desktop is located.`,
			},
			"desktop_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the desktop.`,
			},
			"remote_console": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopRemoteConsoleSchema(),
				Description: `The remote console information.`,
			},
		},
	}
}

func desktopRemoteConsoleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The login type of console.`,
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The remote login URL of console.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The login protocol of console.`,
			},
		},
	}
}

func queryDesktopRemoteConsole(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/desktops/{desktop_id}/remote-consoles"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{desktop_id}", d.Get("desktop_id").(string))

	requestOpts := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, requestOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenRemoteConsole(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	return []map[string]interface{}{{
		"type":     utils.PathSearch("type", resp, nil),
		"url":      utils.PathSearch("url", resp, nil),
		"protocol": utils.PathSearch("protocol", resp, nil),
	}}
}

func dataSourceDesktopRemoteConsoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	remoteConsole, err := queryDesktopRemoteConsole(client, d)
	if err != nil {
		return diag.Errorf("error getting desktop remote console: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("remote_console", flattenRemoteConsole(remoteConsole)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
