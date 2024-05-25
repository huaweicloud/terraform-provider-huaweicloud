package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/terminals"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v2/{project_id}/terminals/binding-desktops/batch-delete
// @API Workspace POST /v2/{project_id}/terminals/binding-desktops
// @API Workspace GET /v2/{project_id}/terminals/binding-desktops
// @API Workspace GET /v2/{project_id}/terminals/binding-desktops/config
// @API Workspace POST /v2/{project_id}/terminals/binding-desktops/config
func ResourceTerminalBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTerminalBindingCreate,
		ReadContext:   resourceTerminalBindingRead,
		UpdateContext: resourceTerminalBindingUpdate,
		DeleteContext: resourceTerminalBindingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the desktops (to be bound to the MAC address) are located.",
			},
			"bindings": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mac": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The MAC address.",
						},
						"desktop_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The desktop name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The binding description.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the binding policy.",
						},
					},
				},
				Description: "The managed resource configuration.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether bindings are available.",
			},
			"disabled_after_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether disabled the binding function before destroy resource.",
			},
		},
	}
}

func buildTerminalBindings(bindings *schema.Set) []terminals.TerminalBindingInfo {
	if bindings.Len() < 1 {
		return nil
	}
	result := make([]terminals.TerminalBindingInfo, bindings.Len())
	for i, val := range bindings.List() {
		binding := val.(map[string]interface{})
		result[i] = terminals.TerminalBindingInfo{
			Mac:         binding["mac"].(string),
			DesktopName: binding["desktop_name"].(string),
			Description: binding["description"].(string),
		}
	}
	return result
}

func isConfigFunctionEnabled(enabled bool) string {
	if enabled {
		return "ON"
	}
	return "OFF"
}

func resourceTerminalBindingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	configOpts := terminals.UpdateConfigOpts{
		TcBindSwitch: isConfigFunctionEnabled(d.Get("enabled").(bool)),
	}
	// No error will be reported when opening or closing repeatedly.
	err = terminals.UpdateConfig(client, configOpts)
	if err != nil {
		return diag.Errorf("error opening terminal binding configuration: %s", err)
	}

	createOpts := terminals.CreateOpts{
		BindList: buildTerminalBindings(d.Get("bindings").(*schema.Set)),
	}
	err = terminals.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error batch binding MAC addresses to desktops: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceTerminalBindingRead(ctx, d, meta)
}

func flattenTerminalBindings(bindings []terminals.TerminalBindingResp) []interface{} {
	if len(bindings) < 1 {
		return nil
	}

	result := make([]interface{}, len(bindings))
	for i, binding := range bindings {
		result[i] = map[string]interface{}{
			"mac":          binding.MAC,
			"desktop_name": binding.DesktopName,
			"description":  binding.Description,
			"id":           binding.ID,
		}
	}
	return result
}

func resourceTerminalBindingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	configStatus, err := terminals.GetConfig(client)
	if err != nil {
		diag.Errorf("error getting terminal binding configuration: %s", err)
	}

	opts := terminals.ListOpts{
		Offset: utils.Int(0), // Offset is the required query parameter.
		Limit:  1000,         // Limit is the required query parameter.
	}
	resp, err := terminals.List(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Terminal bindings")
	}
	if len(resp) < 1 {
		common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "Terminal bindings")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("enabled", configStatus == "ON"),
		d.Set("bindings", flattenTerminalBindings(resp)),
		d.Set("disabled_after_delete", d.Get("disabled_after_delete").(bool)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func deleteTerminalBindings(client *golangsdk.ServiceClient, bindings *schema.Set) error {
	if bindings.Len() < 1 {
		return nil
	}

	idList := make([]string, bindings.Len())
	for i, val := range bindings.List() {
		binding := val.(map[string]interface{})
		idList[i] = binding["id"].(string)
	}

	opts := terminals.DeleteOpts{
		IDs: idList,
	}
	_, err := terminals.Delete(client, opts)
	if err != nil {
		return fmt.Errorf("error deleting terminals and MAC addresses binding: %s", err)
	}
	return nil
}

func resourceTerminalBindingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	if d.HasChange("bindings") {
		oldRaw, newRaws := d.GetChange("bindings")

		err = deleteTerminalBindings(client, oldRaw.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}

		opts := terminals.CreateOpts{
			BindList: buildTerminalBindings(newRaws.(*schema.Set)),
		}
		err = terminals.Create(client, opts)
		if err != nil {
			return diag.Errorf("error batch binding MAC addresses to desktops: %s", err)
		}
	}

	if d.HasChange("enabled") {
		configOpts := terminals.UpdateConfigOpts{
			TcBindSwitch: isConfigFunctionEnabled(d.Get("enabled").(bool)),
		}
		err = terminals.UpdateConfig(client, configOpts)
		if err != nil {
			return diag.Errorf("error updating terminal binding configuration: %s", err)
		}
	}
	return resourceTerminalBindingRead(ctx, d, meta)
}

func resourceTerminalBindingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	bindings := d.Get("bindings").(*schema.Set)
	err = deleteTerminalBindings(client, bindings)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("disabled_after_delete").(bool) {
		opts := terminals.UpdateConfigOpts{
			TcBindSwitch: "OFF",
		}
		err = terminals.UpdateConfig(client, opts)
		if err != nil {
			return diag.Errorf("error closing terminal binding configuration: %s", err)
		}
	}
	return nil
}
