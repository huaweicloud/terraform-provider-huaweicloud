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

// @API DataArtsStudio POST /v1/{project_id}/service/apis/{api_id}/instances/{instance_id}/authorize
func ResourceDataServiceApiAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceApiAuthCreate,
		ReadContext:   resourceDataServiceApiAuthRead,
		DeleteContext: resourceDataServiceApiAuthDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The region where the API and APP(s) are located.`,
			},

			// Parameter in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The workspace ID to which the API and APP(s) belong.`,
			},

			// Arguments
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the published API.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The exclusive cluster ID to which the published API belongs.`,
			},
			"app_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of application IDs that used to authorize API.`,
			},
			"expired_at": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The expiration time of the APP authorize operation.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildApiAuthBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"app_ids": d.Get("app_ids"),
	}
	if expirationTime, ok := d.GetOk("expired_at"); ok {
		result["time"] = expirationTime.(string)
	} else {
		result["time"] = utils.CalculateNextWholeHourAfterFewTime(utils.GetCurrentTime(true), 48*time.Hour)
	}
	return result
}

func authorizeApi(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v1/{project_id}/service/apis/{api_id}/instances/{instance_id}/authorize"
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
		JSONBody: buildApiAuthBodyParams(d),
		OkCodes:  []int{204},
	}

	_, err := client.Request("POST", debugPath, &opt)
	return err
}

func resourceDataServiceApiAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = authorizeApi(client, d)
	if err != nil {
		return diag.Errorf("failed to authorize APP(s) to access API (%s): %s", d.Get("api_id").(string), err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceDataServiceApiAuthRead(ctx, d, meta)
}

func resourceDataServiceApiAuthRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is only a one-time action resource for authorizing the API.
	return nil
}

func resourceDataServiceApiAuthDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for authorizing the API. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
