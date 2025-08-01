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

// @API Workspace GET /v2/{project_id}/desktops/{desktop_id}/sysprep
func DataSourceDesktopSysprep() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDesktopSysprepRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the desktop sysprep is located.`,
			},
			"desktop_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the desktop to be queried.`,
			},
			"sysprep_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        sysprepInfoSchema(),
				Description: `The sysprep information of the desktop.`,
			},
		},
	}
}

func sysprepInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sysprep_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sysprep version of the desktop.`,
			},
			"support_create_image": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the desktop supports creating image.`,
			},
		},
	}
}

func queryDesktopSysprep(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl   = "v2/{project_id}/desktops/{desktop_id}/sysprep"
		desktopId = d.Get("desktop_id").(string)
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{desktop_id}", desktopId)

	requestOpts := &golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, requestOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func flattenSysprepInfo(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"sysprep_version":      utils.PathSearch("sysprep_version", respBody, nil),
			"support_create_image": utils.PathSearch("support_create_image", respBody, nil),
		},
	}
}

func dataSourceDesktopSysprepRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating WorkSpace client: %s", err)
	}

	respBody, err := queryDesktopSysprep(client, d)
	if err != nil {
		return diag.Errorf("error querying desktop sysprep information: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("sysprep_info", flattenSysprepInfo(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
