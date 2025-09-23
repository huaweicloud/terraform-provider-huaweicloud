package er

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var attachmentAccepterNonUpdatableParams = []string{"instance_id", "attachment_id", "action"}

// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/attachments/{attachment_id}/accept
// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/attachments/{attachment_id}/reject
func ResourceAttachmentAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAttachmentAccepterCreate,
		UpdateContext: resourceAttachmentAccepterUpdate,
		ReadContext:   resourceAttachmentAccepterRead,
		DeleteContext: resourceAttachmentAccepterDelete,

		CustomizeDiff: config.FlexibleForceNew(attachmentAccepterNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the ER instance.`,
			},
			"attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the attachment to be action.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action type.`,
			},
		},
	}
}

func resourceAttachmentAccepterCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("er", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	var (
		httpUrl      = "v3/{project_id}/enterprise-router/{er_id}/attachments/{attachment_id}/{action}"
		attachmentId = d.Get("attachment_id").(string)
		action       = d.Get("action").(string)
	)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{er_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{attachment_id}", attachmentId)
	createPath = strings.ReplaceAll(createPath, "{action}", action)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	reap, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("unable to %s the attachment(%s): %s", action, attachmentId, err)
	}

	respBody, err := utils.FlattenResponse(reap)
	if err != nil {
		return diag.FromErr(err)
	}

	requestId := utils.PathSearch("request_id", respBody, "").(string)
	if requestId == "" {
		return diag.Errorf("unable to find the request ID from the API response")
	}

	d.SetId(requestId)

	return nil
}

func resourceAttachmentAccepterRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAttachmentAccepterUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAttachmentAccepterDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the attachment. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
