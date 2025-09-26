package apig

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

var apigChannelMemberBatchActionNonUpdatableParams = []string{"instance_id", "vpc_channel_id", "action", "member_ids"}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members/batch-enable
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members/batch-disable
func ResourceChannelMemberBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelMemberBatchActionCreate,
		ReadContext:   resourceChannelMemberBatchActionRead,
		UpdateContext: resourceChannelMemberBatchActionUpdate,
		DeleteContext: resourceChannelMemberBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(apigChannelMemberBatchActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the dedicated instance to which the VPC channel belongs is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the VPC channel belongs.`,
			},
			"vpc_channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPC channel.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action to perform on the VPC channel members.`,
			},

			// Optional parameters.
			"member_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of member IDs to perform the action on.`,
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

func buildChannelMemberBatchActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"member_ids": d.Get("member_ids"),
	}
}

func resourceChannelMemberBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		httpUrl      = "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members/batch-{action}"
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		instanceId   = d.Get("instance_id").(string)
		vpcChannelId = d.Get("vpc_channel_id").(string)
		action       = d.Get("action").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{vpc_channel_id}", vpcChannelId)
	createPath = strings.ReplaceAll(createPath, "{action}", action)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildChannelMemberBatchActionBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error performing batch action on VPC channel members: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceChannelMemberBatchActionRead(ctx, d, meta)
}

func resourceChannelMemberBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceChannelMemberBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceChannelMemberBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for performing an operation with the VPC channel member list.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
