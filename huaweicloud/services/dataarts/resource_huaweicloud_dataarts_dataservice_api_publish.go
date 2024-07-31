package dataarts

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v1/{project_id}/service/apis/{api_id}/instances/{instance_id}/publish
func ResourceDataServiceApiPublish() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceApiPublishCreate,
		ReadContext:   resourceDataServiceApiPublishRead,
		DeleteContext: resourceDataServiceApiPublishDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the published API is located.`,
			},

			// Parameter in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The workspace ID to which the published API belongs.`,
			},

			// Arguments
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the API to be published.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The exclusive cluster ID to which the published API belongs in Data Service side.`,
			},
			"apig_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The type of the APIG object.`,
			},
			"apig_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The APIG instance ID to which the API is published simultaneously in APIG service.`,
			},
			"apig_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The APIG group ID to which the published API belongs.`,
			},
			"roma_app_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The application ID for ROMA APIC.`,
			},
		},
	}
}

func buildApiPublishBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"apig_type":        utils.ValueIgnoreEmpty(d.Get("apig_type")),
		"apig_instance_id": utils.ValueIgnoreEmpty(d.Get("apig_instance_id")),
		"group_id_in_apig": utils.ValueIgnoreEmpty(d.Get("apig_group_id")),
		"roma_app_id":      utils.ValueIgnoreEmpty(d.Get("roma_app_id")),
	}
}

func publishApi(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v1/{project_id}/service/apis/{api_id}/instances/{instance_id}/publish"
		workspaceId = d.Get("workspace_id").(string)
		apiId       = d.Get("api_id").(string)
		instanceId  = d.Get("instance_id").(string)
	)
	debugPath := client.Endpoint + httpUrl
	debugPath = strings.ReplaceAll(debugPath, "{project_id}", client.ProjectID)
	debugPath = strings.ReplaceAll(debugPath, "{api_id}", apiId)
	debugPath = strings.ReplaceAll(debugPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     "EXCLUSIVE",
		},
		JSONBody: utils.RemoveNil(buildApiPublishBodyParams(d)),
		OkCodes:  []int{204},
	}

	// Only one publishing task can be executed at a time.
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)
	_, err := client.Request("POST", debugPath, &opt)
	if err != nil {
		return err
	}
	return nil
}

func resourceDataServiceApiPublishCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = publishApi(client, d)
	if err != nil {
		return diag.Errorf("error publishing API: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceDataServiceApiPublishRead(ctx, d, meta)
}

func resourceDataServiceApiPublishRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is only a one-time action resource for publishing the API.
	return nil
}

func resourceDataServiceApiPublishDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for publishing the API. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
