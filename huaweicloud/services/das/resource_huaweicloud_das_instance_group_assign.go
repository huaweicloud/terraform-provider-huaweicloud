package das

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

var instanceGroupAssignNonUpdatableParams = []string{
	"group_id",
	"instance_ids",
}

// @API DAS POST /v3/{project_id}/batch-inspection/add-instance-to-group
func ResourceInstanceGroupAssign() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceGroupAssignCreate,
		ReadContext:   resourceInstanceGroupAssignRead,
		UpdateContext: resourceInstanceGroupAssignUpdate,
		DeleteContext: resourceInstanceGroupAssignDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceGroupAssignNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the DAS instance group is located.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The instance group ID.`,
			},
			"instance_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of instance IDs to be assigned to the group.`,
			},

			// Internal
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func buildInstanceGroupAssignBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"group_id":     d.Get("group_id").(string),
		"instance_ids": utils.ExpandToStringList(d.Get("instance_ids").([]interface{})),
	}
}

func resourceInstanceGroupAssignCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/batch-inspection/add-instance-to-group"
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildInstanceGroupAssignBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error assigning instances to DAS instance group: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	return resourceInstanceGroupAssignRead(ctx, d, meta)
}

func resourceInstanceGroupAssignRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceGroupAssignUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceGroupAssignDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for assigning instances to a DAS instance group.
Deleting this resource will not clear the corresponding request record, but will only remove the resource
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
