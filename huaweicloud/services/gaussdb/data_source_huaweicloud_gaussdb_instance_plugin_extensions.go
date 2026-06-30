package gaussdb

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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/plugin-extensions
func DataSourceInstancePluginExtensions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancePluginExtensionsRead,

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
			"plugin_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"extensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     pluginExtensionSchema(),
			},
		},
	}
}

func pluginExtensionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"extension_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceInstancePluginExtensionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/plugin-extensions?plugin_name={plugin_name}&db_name={db_name}"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{plugin_name}", d.Get("plugin_name").(string))
	listPath = strings.ReplaceAll(listPath, "{db_name}", d.Get("db_name").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error querying GaussDB instance plugin extensions: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("extensions", flattenPluginExtensions(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPluginExtensions(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curArray := utils.PathSearch("[*]", resp, make([]interface{}, 0)).([]interface{})
	res := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"extension_name": utils.PathSearch("extension_name", v, nil),
			"status":         utils.PathSearch("status", v, nil),
		})
	}
	return res
}
