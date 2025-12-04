package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var appSharedFolderAssignNonUpdatableParams = []string{"storage_id", "storage_claim_id", "add_items", "del_items"}

// @API Workspace POST /v1/{project_id}/persistent-storages/{storage_id}/actions/assign-share-folder
func ResourceAppSharedFolderAssign() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSharedFolderAssignCreate,
		ReadContext:   resourceAppSharedFolderAssignRead,
		UpdateContext: resourceAppSharedFolderAssignUpdate,
		DeleteContext: resourceAppSharedFolderAssignDelete,

		CustomizeDiff: config.FlexibleForceNew(appSharedFolderAssignNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the WKS storage is located.`,
			},

			// Required parameters.
			"storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The WKS storage ID to which the shared folder belongs.`,
			},
			"storage_claim_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The WKS storage directory claim ID.`,
			},

			// Optional parameters.
			"add_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        addItemSchema(),
				Description: `The list of members to be added.`,
			},
			"del_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        delItemSchema(),
				Description: `The list of members to be removed.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func addItemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_statement_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The policy ID.`,
			},
			"attach": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The target.`,
			},
			"attach_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The associated object type.`,
			},
		},
	}
}

func delItemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attach": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The target.`,
			},
			"attach_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The associated object type.`,
			},
		},
	}
}

func buildAppSharedFolderAddItems(addItems []interface{}) []map[string]interface{} {
	if len(addItems) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(addItems))
	for _, item := range addItems {
		result = append(result, map[string]interface{}{
			"policy_statement_id": utils.PathSearch("policy_statement_id", item, ""),
			"attach":              utils.PathSearch("attach", item, ""),
			"attach_type":         utils.PathSearch("attach_type", item, ""),
		})
	}

	return result
}

func buildAppSharedFolderDelItems(delItems []interface{}) []map[string]interface{} {
	if len(delItems) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(delItems))
	for _, item := range delItems {
		result = append(result, map[string]interface{}{
			"attach":      utils.PathSearch("attach", item, ""),
			"attach_type": utils.PathSearch("attach_type", item, ""),
		})
	}

	return result
}

func buildAppSharedFolderAssignBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"storage_claim_id": d.Get("storage_claim_id").(string),
		"add_items":        buildAppSharedFolderAddItems(d.Get("add_items").([]interface{})),
		"del_items":        buildAppSharedFolderDelItems(d.Get("del_items").([]interface{})),
	}
}

func resourceAppSharedFolderAssignCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		httpUrl   = "v1/{project_id}/persistent-storages/{storage_id}/actions/assign-share-folder"
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		storageId = d.Get("storage_id").(string)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{storage_id}", storageId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAppSharedFolderAssignBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error assigning shared folder members: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceAppSharedFolderAssignRead(ctx, d, meta)
}

func resourceAppSharedFolderAssignRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppSharedFolderAssignUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppSharedFolderAssignDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for assignment of user access to a shared folder. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
