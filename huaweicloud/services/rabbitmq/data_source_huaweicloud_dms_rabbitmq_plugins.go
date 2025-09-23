package rabbitmq

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RabbitMQ GET /v2/{project_id}/instances/{instance_id}/rabbitmq/plugins
func DataSourceDmsRabbitmqPlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRabbitmqPluginsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RabbitMQ instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the plugin.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the plugin is enabled.`,
			},
			"running": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the plugin is running.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the version of the plugin.`,
			},
			"plugins": {
				Type:        schema.TypeList,
				Elem:        pluginsSchema(),
				Computed:    true,
				Description: `Indicates the list of the plugins.`,
			},
		},
	}
}

func pluginsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the plugin.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the plugin is enabled.`,
			},
			"running": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the plugin is running.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version of the plugin.`,
			},
		},
	}
	return &sc
}

func resourceDmsRabbitmqPluginsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getRabbitmqPlugins: query DMS RabbitMQ plugins
	var (
		getRabbitmqPluginsHttpUrl = "v2/{project_id}/instances/{instance_id}/rabbitmq/plugins"
		getRabbitmqPluginsProduct = "dms"
	)
	getRabbitmqPluginsClient, err := cfg.NewServiceClient(getRabbitmqPluginsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	getRabbitmqPluginsPath := getRabbitmqPluginsClient.Endpoint + getRabbitmqPluginsHttpUrl
	getRabbitmqPluginsPath = strings.ReplaceAll(getRabbitmqPluginsPath, "{project_id}",
		getRabbitmqPluginsClient.ProjectID)
	getRabbitmqPluginsPath = strings.ReplaceAll(getRabbitmqPluginsPath, "{instance_id}", instanceID)

	getRabbitmqPluginsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getRabbitmqPluginsResp, err := getRabbitmqPluginsClient.Request("GET", getRabbitmqPluginsPath,
		&getRabbitmqPluginsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS RabbitMQ plugins")
	}

	getRabbitmqPluginsRespBody, err := utils.FlattenResponse(getRabbitmqPluginsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("plugins", flattenPlugins(filterPlugins(d, getRabbitmqPluginsRespBody))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPlugins(pluginRespBody []interface{}) []interface{} {
	if pluginRespBody == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(pluginRespBody))
	for _, v := range pluginRespBody {
		rst = append(rst, map[string]interface{}{
			"name":    utils.PathSearch("name", v, nil),
			"enable":  utils.PathSearch("enable", v, nil),
			"running": utils.PathSearch("running", v, nil),
			"version": utils.PathSearch("version", v, nil),
		})
	}
	return rst
}

func filterPlugins(d *schema.ResourceData, resp interface{}) []interface{} {
	pluginJson := utils.PathSearch("plugins", resp, make([]interface{}, 0))
	pluginArray := pluginJson.([]interface{})
	if len(pluginArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(pluginArray))

	rawName, rawNameOK := d.GetOk("name")
	rawEnable := d.Get("enable")
	rawRunning := d.Get("running")
	rawVersion, rawVersionOK := d.GetOk("version")

	for _, plugin := range pluginArray {
		name := utils.PathSearch("name", plugin, nil)
		enable := utils.PathSearch("enable", plugin, false).(bool)
		running := utils.PathSearch("running", plugin, false).(bool)
		version := utils.PathSearch("version", plugin, nil)

		if rawNameOK && rawName != name {
			continue
		}
		if rawEnable != enable {
			continue
		}
		if rawRunning != running {
			continue
		}
		if rawVersionOK && rawVersion != version {
			continue
		}
		result = append(result, plugin)
	}

	return result
}
