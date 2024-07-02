package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/plugins"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourcePlugin defines the provider resource of the APIG plugin.
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instanceId}/plugins/{plugin_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instanceId}/plugins/{plugin_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instanceId}/plugins/{plugin_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instanceId}/plugins
func ResourcePlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginCreate,
		ReadContext:   resourcePluginRead,
		UpdateContext: resourcePluginUpdate,
		DeleteContext: resourcePluginDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePluginImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the plugin is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the plugin belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plugin name.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The plugin type.",
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return utils.JSONStringsEqual(old, new)
				},
				Description: "The configuration details for plugin.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The plugin description.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the plugin.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the plugin.",
			},
		},
	}
}

func resourcePluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := plugins.CreateOpts{
		InstanceId:  d.Get("instance_id").(string),
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Scope:       "global",
		Content:     d.Get("content").(string),
		Description: d.Get("description").(string),
	}
	resp, err := plugins.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating the plugin: %s", err)
	}
	d.SetId(resp.ID)

	return resourcePluginRead(ctx, d, meta)
}

func resourcePluginRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	pluginId := d.Id()
	resp, err := plugins.Get(client, instanceId, pluginId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "plugin policy")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("type", resp.Type),
		d.Set("content", resp.Content),
		d.Set("description", resp.Description),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving plugin fields: %s", err)
	}
	return nil
}

func resourcePluginUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		pluginId = d.Id()
		opts     = plugins.UpdateOpts{
			InstanceId:  d.Get("instance_id").(string),
			ID:          pluginId,
			Name:        d.Get("name").(string),
			Type:        d.Get("type").(string),
			Scope:       "global",
			Content:     d.Get("content").(string),
			Description: utils.String(d.Get("description").(string)),
		}
	)

	_, err = plugins.Update(client, opts)
	if err != nil {
		return diag.Errorf("error updating the plugin (%s): %s", pluginId, err)
	}

	return resourcePluginRead(ctx, d, meta)
}

func resourcePluginDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		pluginId   = d.Id()
	)
	err = plugins.Delete(client, instanceId, pluginId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting the plugin (%s)", pluginId))
	}
	return nil
}

func resourcePluginImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])
	err := d.Set("instance_id", parts[0])
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving instance ID field: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
