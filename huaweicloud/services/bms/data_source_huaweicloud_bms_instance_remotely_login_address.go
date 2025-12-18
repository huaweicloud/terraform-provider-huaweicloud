package bms

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

// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/remote_console
func DataSourceInstanceRemotelyLoginAddress() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceRemotelyLoginAddressRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_console": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     remoteConsoleSchema(),
			},
		},
	}
}

func remoteConsoleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceInstanceRemotelyLoginAddressRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/remote_console"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_id}", d.Get("server_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOpt.JSONBody = buildGetInstanceRemotelyLoginAddressBodyParams(d)

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving BMS instance remotely login address: %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("remote_console", flattenRemoteConsole(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetInstanceRemotelyLoginAddressBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"remote_console": buildGetInstanceRemotelyLoginAddressRemoteConsoleBodyParams(d),
	}
	return bodyParams
}

func buildGetInstanceRemotelyLoginAddressRemoteConsoleBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("remote_console").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw, ok := rawParams[0].(map[string]interface{})
	if !ok {
		return nil
	}

	rst := map[string]interface{}{
		"protocol": raw["protocol"],
		"type":     raw["type"],
	}

	return rst
}

func flattenRemoteConsole(resp interface{}) []interface{} {
	curJson := utils.PathSearch("remote_console", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"type":     utils.PathSearch("type", curJson, nil),
			"protocol": utils.PathSearch("protocol", curJson, nil),
			"url":      utils.PathSearch("url", curJson, nil),
		},
	}
	return rst
}
