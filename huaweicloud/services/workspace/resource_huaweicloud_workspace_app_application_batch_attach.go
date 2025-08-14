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

var applicationBatchAttachNonUpdateParams = []string{"server_id", "record_ids"}

// @API Workspace POST /v1/{project_id}/image-servers/{server_id}/actions/attach-app
func ResourceAppApplicationBatchAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppApplicationBatchAttachCreate,
		ReadContext:   resourceAppApplicationBatchAttachRead,
		UpdateContext: resourceAppApplicationBatchAttachUpdate,
		DeleteContext: resourceAppApplicationBatchAttachDelete,

		CustomizeDiff: config.FlexibleForceNew(applicationBatchAttachNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the Workspace APP is located.`,
			},
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the image server instance.`,
			},
			"record_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of application record IDs to be attach.`,
			},
			"uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the application attachment.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceAppApplicationBatchAttachCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		httpUrl  = "v1/{project_id}/image-servers/{server_id}/actions/attach-app"
		serverId = d.Get("server_id").(string)
	)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{server_id}", serverId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"items": utils.ExpandToStringList(d.Get("record_ids").([]interface{})),
		},
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("unable to attach applications to the server (%s): %s", serverId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	return diag.FromErr(d.Set("uri", utils.PathSearch("uri", respBody, nil)))
}

func resourceAppApplicationBatchAttachRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppApplicationBatchAttachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppApplicationBatchAttachDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to attach applications to specified image instance. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
