package rc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var resourceRemoveNonUpdatableParams = []string{
	"group_id",
	"resource_id",
}

// @API RC DELETE /v1/resource-center/groups/{group_id}/resources/{resource_id}
func ResourceResourceRemove() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceRemoveCreate,
		ReadContext:   resourceResourceRemoveRead,
		UpdateContext: resourceResourceRemoveUpdate,
		DeleteContext: resourceResourceRemoveDelete,

		CustomizeDiff: config.FlexibleForceNew(resourceRemoveNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource to be removed is located.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the resource group.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the resource to be removed.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceResourceRemoveCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		groupId    = d.Get("group_id").(string)
		resourceId = d.Get("resource_id").(string)
		httpUrl    = "v1/resource-center/groups/{group_id}/resources/{resource_id}"
	)

	client, err := cfg.NewServiceClient("rc", region)
	if err != nil {
		return diag.Errorf("error creating RC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{group_id}", groupId)
	createPath = strings.ReplaceAll(createPath, "{resource_id}", resourceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", createPath, &opt)
	if err != nil {
		return diag.Errorf("error removing resource from resource group: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	return nil
}

func resourceResourceRemoveRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourceRemoveUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourceRemoveDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to remove resource. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
