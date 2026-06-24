package gaussdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var gaussdbPluginExtensionNonUpdatableParams = []string{"instance_id", "db_name", "plugin_name", "extension_name"}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/config-plugin-extensions
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/plugin-extensions
func ResourceGaussDbInstancePluginExtensionsConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginExtensionCreate,
		UpdateContext: resourcePluginExtensionUpdate,
		ReadContext:   resourcePluginExtensionRead,
		DeleteContext: resourcePluginExtensionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePluginExtensionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(gaussdbPluginExtensionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plugin_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"extension_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildPluginExtensionBodyParams(d *schema.ResourceData, action string) map[string]interface{} {
	return map[string]interface{}{
		"plugin_name":      d.Get("plugin_name"),
		"extension_name":   d.Get("extension_name"),
		"extension_action": action,
		"db_list":          []string{d.Get("db_name").(string)},
	}
}

func resourcePluginExtensionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/config-plugin-extensions"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildPluginExtensionBodyParams(d, "on")),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error configuring GaussDB instance plugin extension: %s", err)
	}

	pluginName := d.Get("plugin_name").(string)
	extensionName := d.Get("extension_name").(string)
	dbName := d.Get("db_name").(string)
	d.SetId(fmt.Sprintf("%s/%s/%s/%s", instanceID, dbName, pluginName, extensionName))

	return resourcePluginExtensionRead(ctx, d, meta)
}

func resourcePluginExtensionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl       = "v3/{project_id}/instances/{instance_id}/plugin-extensions"
		product       = "opengauss"
		instanceID    = d.Get("instance_id").(string)
		pluginName    = d.Get("plugin_name").(string)
		extensionName = d.Get("extension_name").(string)
		dbName        = d.Get("db_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	getPath += fmt.Sprintf("?plugin_name=%s&db_name=%s", pluginName, dbName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB plugin extensions")
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	extension := utils.PathSearch(
		fmt.Sprintf("[?extension_name=='%s'] | [0]", extensionName),
		respBody, nil)
	if extension == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
			fmt.Sprintf("GaussDB plugin extension %s is not available for plugin %s on database %s of instance %s",
				extensionName, pluginName, dbName, instanceID))
	}
	status := utils.PathSearch("status", extension, "").(string)
	if status == "" || status == "off" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
			fmt.Sprintf("GaussDB plugin extension %s is not available for plugin %s on database %s of instance %s",
				extensionName, pluginName, dbName, instanceID))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("db_name", dbName),
		d.Set("plugin_name", pluginName),
		d.Set("extension_name", extensionName),
		d.Set("status", status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePluginExtensionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePluginExtensionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/config-plugin-extensions"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)

	bodyParams := buildPluginExtensionBodyParams(d, "off")

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(bodyParams),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error disabling GaussDB instance plugin extension: %s", err)
	}

	return nil
}

func resourcePluginExtensionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 4 {
		return nil, errors.New("invalid format specified for import ID, " +
			"want '<instance_id>/<db_name>/<plugin_name>/<extension_name>'")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("db_name", parts[1]),
		d.Set("plugin_name", parts[2]),
		d.Set("extension_name", parts[3]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
