package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/plugins"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourcePluginAssociate defines the provider resource of the APIG plugin binding.
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/detach
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attached-apis
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attach
func ResourcePluginAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginAssociateCreate,
		ReadContext:   resourcePluginAssociateRead,
		UpdateContext: resourcePluginAssociateUpdate,
		DeleteContext: resourcePluginAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePluginAssociateImportState,
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
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The plugin ID.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The environment ID where the API was published.",
			},
			"api_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The APIs bound by the plugin.",
			},
		},
	}
}

func resourcePluginAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		pluginId   = d.Get("plugin_id").(string)
		envId      = d.Get("env_id").(string)
		opts       = plugins.BindOpts{
			InstanceId: instanceId,
			PluginId:   pluginId,
			EnvId:      envId,
			ApiIds:     utils.ExpandToStringListBySet(d.Get("api_ids").(*schema.Set)),
		}
	)

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err = plugins.Bind(client, opts)
		retryable, err := handleMultiOperationsError(err)
		if retryable {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error binding the plugin to the APIs: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, pluginId, envId))

	return resourcePluginAssociateRead(ctx, d, meta)
}

func flattenBindApis(bindList []plugins.BindApiInfo) []string {
	if len(bindList) < 1 {
		return nil
	}

	result := make([]string, len(bindList))
	for i, v := range bindList {
		result[i] = v.ApiId
	}
	return result
}

func resourcePluginAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		pluginId   = d.Get("plugin_id").(string)
		listOpts   = plugins.ListBindOpts{
			InstanceId: instanceId,
			PluginId:   pluginId,
			EnvId:      d.Get("env_id").(string),
		}
	)

	resp, err := plugins.ListBind(client, listOpts)
	if err != nil {
		// A 404 error will returned if the instance is not exist.
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error querying the bind details of the plugin")
	}

	bindApiIds := flattenBindApis(resp)
	if len(bindApiIds) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "plugin associate")
	}
	if err = d.Set("api_ids", bindApiIds); err != nil {
		return diag.Errorf("error saving API IDs field: %s", err)
	}
	return nil
}

func resourcePluginAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId     = d.Get("instance_id").(string)
		pluginId       = d.Get("plugin_id").(string)
		envId          = d.Get("env_id").(string)
		oldVal, newVal = d.GetChange("api_ids")
		rmRaw          = oldVal.(*schema.Set).Difference(newVal.(*schema.Set))
		addRaw         = newVal.(*schema.Set).Difference(oldVal.(*schema.Set))
	)

	// Step 1, unbind the APIs.
	if rmRaw.Len() > 0 {
		unbindOpts := plugins.UnbindOpts{
			InstanceId: instanceId,
			PluginId:   pluginId,
			EnvId:      envId,
			ApiIds:     utils.ExpandToStringListBySet(rmRaw),
		}
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			err = plugins.Unbind(client, unbindOpts)
			retryable, err := handleMultiOperationsError(err)
			if retryable {
				return resource.RetryableError(err)
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return diag.Errorf("an error occurred while updating plugin bindings, some APIs could not be unbound from "+
				"the plugin (%s): %s", pluginId, err)
		}
	}

	// Step 2, bind the APIs.
	if addRaw.Len() > 0 {
		bindOpts := plugins.BindOpts{
			InstanceId: instanceId,
			PluginId:   pluginId,
			EnvId:      envId,
			ApiIds:     utils.ExpandToStringListBySet(addRaw),
		}
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			_, err = plugins.Bind(client, bindOpts)
			retryable, err := handleMultiOperationsError(err)
			if retryable {
				return resource.RetryableError(err)
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return diag.Errorf("an error occurred while updating plugin bindings, some APIs could not be bound to the "+
				"plugin (%s): %s", pluginId, err)
		}
	}

	return resourcePluginAssociateRead(ctx, d, meta)
}

func resourcePluginAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		pluginId   = d.Get("plugin_id").(string)
		envId      = d.Get("env_id").(string)
		listOpts   = plugins.ListBindOpts{
			InstanceId: instanceId,
			PluginId:   pluginId,
			EnvId:      d.Get("env_id").(string),
		}
	)
	resp, err := plugins.ListBind(client, listOpts)
	if err != nil {
		return diag.Errorf("error querying the bind details of the plugin: %s", err)
	}

	if apiIds := flattenBindApis(resp); len(apiIds) > 0 {
		unbindOpts := plugins.UnbindOpts{
			InstanceId: instanceId,
			PluginId:   pluginId,
			EnvId:      envId,
			ApiIds:     apiIds,
		}
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			err = plugins.Unbind(client, unbindOpts)
			retryable, err := handleMultiOperationsError(err)
			if retryable {
				return resource.RetryableError(err)
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return diag.Errorf("an error occurred while deleting plugin bindings, some APIs could not be unbound from "+
				"the plugin (%s): %s", pluginId, err)
		}
	}
	return nil
}

func resourcePluginAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<plugin_id>/<env_id>', "+
			"but got '%s'", importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("plugin_id", parts[1]),
		d.Set("env_id", parts[2]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving associate resource field: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
