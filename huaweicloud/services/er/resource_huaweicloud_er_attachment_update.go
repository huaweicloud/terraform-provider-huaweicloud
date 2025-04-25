package er

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var attachmenUpdateNonUpdatableParams = []string{
	"instance_id",
	"attachment_id",
}

// @API ER PUT /v3/{project_id}/enterprise-router/{er_id}/attachments/{attachment_id}
// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/attachments/{attachment_id}
func ResourceAttachmentUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAttachmentUpdateCreate,
		UpdateContext: resourceAttachmentUpdateUpdate,
		ReadContext:   resourceAttachmentUpdateRead,
		DeleteContext: resourceAttachmentUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(attachmenUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the attachment is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the ER instance.`,
			},
			"attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the attachment to be updated.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The new name of the attachment.`,
				AtLeastOneOf: []string{
					"description",
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The new description of the attachment.`,
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

func doAttachmentUpdate(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl      = "v3/{project_id}/enterprise-router/{er_id}/attachments/{attachment_id}"
		instanceId   = d.Get("instance_id").(string)
		attachmentId = d.Get("attachment_id").(string)
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{er_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{attachment_id}", attachmentId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"attachment": map[string]interface{}{
				"name":        utils.ValueIgnoreEmpty(d.Get("name")),
				"description": utils.ValueIgnoreEmpty(d.Get("description")),
			},
		},
	}
	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceAttachmentUpdateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	err = doAttachmentUpdate(client, d)
	if err != nil {
		return diag.Errorf("error updating attachment information: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceAttachmentUpdateRead(ctx, d, meta)
}

func GetAttachmentById(client *golangsdk.ServiceClient, instanceId, attachmentId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/enterprise-router/{er_id}/attachments/{attachment_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{er_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{attachment_id}", attachmentId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceAttachmentUpdateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		instanceId   = d.Get("instance_id").(string)
		attachmentId = d.Get("attachment_id").(string)
	)
	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	attachment, err := GetAttachmentById(client, instanceId, attachmentId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting attachment information")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("attachment.name", attachment, nil)),
		d.Set("description", utils.PathSearch("attachment.description", attachment, nil)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the ER resource tags: %s", mErr)
	}
	return nil
}

func resourceAttachmentUpdateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	err = doAttachmentUpdate(client, d)
	if err != nil {
		return diag.Errorf("error updating attachment information: %s", err)
	}

	return resourceAttachmentUpdateRead(ctx, d, meta)
}

func resourceAttachmentUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for updating the attachment. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
