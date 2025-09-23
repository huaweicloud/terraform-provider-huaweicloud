package dataarts

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

// @API DataArtsStudio POST /v1/{project_id}/service/messages
func ResourceDataServiceMessageApprove() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceMessageApproveCreate,
		ReadContext:   resourceDataServiceMessageApproveRead,
		DeleteContext: resourceDataServiceMessageApproveDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the message (to be approved) is located.`,
			},

			// Parameter in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The workspace ID of the exclusive API to which the message (to be approved) belongs.`,
			},

			// Arguments
			"message_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the message (to be approved).`,
			},
			"action": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `The approve action performed by the message.`,
			},
			"time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The execution time of the message action.`,
			},
		},
	}
}

func buildMessageApproveBodyParams(d *schema.ResourceData) map[string]interface{} {
	action := d.Get("action").(int)
	result := map[string]interface{}{
		"message_id": utils.ValueIgnoreEmpty(d.Get("message_id")),
		"action":     action,
	}

	if action == 1 {
		result["time"] = d.Get("time")
	}
	return result
}

func doMessageApprove(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v1/{project_id}/service/messages"
		workspaceId = d.Get("workspace_id").(string)
	)
	debugPath := client.Endpoint + httpUrl
	debugPath = strings.ReplaceAll(debugPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     "EXCLUSIVE",
		},
		JSONBody: buildMessageApproveBodyParams(d),
		OkCodes:  []int{204},
	}

	_, err := client.Request("POST", debugPath, &opt)
	return err
}

func resourceDataServiceMessageApproveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = doMessageApprove(client, d)
	if err != nil {
		return diag.Errorf("error approving message: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceDataServiceMessageApproveRead(ctx, d, meta)
}

func resourceDataServiceMessageApproveRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is a one-time action resource used only for approval messages.
	return nil
}

func resourceDataServiceMessageApproveDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used only for approval messages. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
