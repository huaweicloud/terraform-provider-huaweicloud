package dataarts

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DataArtsStudio POST /v2/{project_id}/design/approvals/batch-offline
func ResourceArchitectureBatchUnpublish() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureBatchUnpublishCreate,
		ReadContext:   resourceArchitectureBatchUnpublishRead,
		DeleteContext: resourceArchitectureBatchUnpublishDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of DataArts Studio workspace.",
			},
			"biz_infos": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the list of objects to be unpublished.",
				Elem:        bizInfoSchema(),
			},
			"approver_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the user ID of the architecture reviewer.",
			},
			"approver_user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the user name of the architecture reviewer.",
			},
			"fast_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies whether to automatically review.",
			},
		},
	}
}

func resourceArchitectureBatchUnpublishCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if err := batchOfflineResource(client, d); err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)
	return nil
}

func resourceArchitectureBatchUnpublishRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceArchitectureBatchUnpublishDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for unpublishing resources. Deleting this resource will not clear
	the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
