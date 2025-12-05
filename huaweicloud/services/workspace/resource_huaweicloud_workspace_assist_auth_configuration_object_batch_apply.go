package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v2/{project_id}/assist-auth-config/apply-objects
func ResourceAssistAuthConfigurationObjectBatchApply() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssistAuthConfigurationObjectBatchApplyCreate,
		ReadContext:   resourceAssistAuthConfigurationObjectBatchApplyRead,
		UpdateContext: resourceAssistAuthConfigurationObjectBatchApplyUpdate,
		DeleteContext: resourceAssistAuthConfigurationObjectBatchApplyDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the apply objects of assist auth configuration are located.`,
			},

			// Optional parameters.
			"add": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The list of objects to be added.`,
				Elem:        assistAuthConfigurationObjectSchema(),
			},
			"delete": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The list of objects to be removed.`,
				Elem:        assistAuthConfigurationObjectSchema(),
			},
		},
	}
}

func assistAuthConfigurationObjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the binding object.`,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the user or user group.`,
			},
			"object_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the user or user group.`,
			},
		},
	}
}

func buildAssistAuthConfigurationApplyObjects(objects []interface{}) []map[string]interface{} {
	if len(objects) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(objects))
	for _, obj := range objects {
		result = append(result, map[string]interface{}{
			"object_type": utils.PathSearch("object_type", obj, nil),
			"object_id":   utils.PathSearch("object_id", obj, nil),
			"object_name": utils.PathSearch("object_name", obj, nil),
		})
	}

	return result
}

func updateAssistAuthConfigurationApplyObjects(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/assist-auth-config/apply-objects"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"add":    utils.ValueIgnoreEmpty(buildAssistAuthConfigurationApplyObjects(d.Get("add").([]interface{}))),
			"delete": utils.ValueIgnoreEmpty(buildAssistAuthConfigurationApplyObjects(d.Get("delete").([]interface{}))),
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	return nil
}

func resourceAssistAuthConfigurationObjectBatchApplyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	err = updateAssistAuthConfigurationApplyObjects(client, d)
	if err != nil {
		return diag.Errorf("error updating assist auth configuration apply objects: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceAssistAuthConfigurationObjectBatchApplyRead(ctx, d, meta)
}

func resourceAssistAuthConfigurationObjectBatchApplyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAssistAuthConfigurationObjectBatchApplyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAssistAuthConfigurationObjectBatchApplyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to update assist auth configuration apply objects.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information
from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
