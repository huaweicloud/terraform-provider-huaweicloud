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

// @API DataArtsStudio POST /v1/{project_id}/service/apis/authorize/action
func ResourceDataServiceApiAuthAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceApiAuthActionCreate,
		ReadContext:   resourceDataServiceApiAuthActionRead,
		DeleteContext: resourceDataServiceApiAuthActionDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The region where the published API and APP to be operated are located.`,
			},

			// Parameter in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The workspace ID to which the published API and APP to be operated belong.`,
			},

			// Arguments
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the published API that to be operated.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The exclusive cluster ID to which the published API belongs.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the APP to be operated.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The operation type of the authorization.`,
			},
			"expired_at": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The expiration time of the authorize operation.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildApiAuthActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"api_id":      d.Get("api_id"),
		"instance_id": d.Get("instance_id"),
		"app_id":      d.Get("app_id"),
		"apply_type":  d.Get("type"),
	}
	if expirationTime, ok := d.GetOk("expired_at"); ok {
		result["time"] = expirationTime.(string)
	} else {
		result["time"] = utils.CalculateNextWholeHourAfterFewTime(utils.GetCurrentTime(true), 48*time.Hour)
	}
	return result
}

func doApiAuthAction(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/service/apis/authorize/action"
	debugPath := client.Endpoint + httpUrl
	debugPath = strings.ReplaceAll(debugPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"Dlm-Type":     "EXCLUSIVE",
		},
		JSONBody: buildApiAuthActionBodyParams(d),
		OkCodes:  []int{204},
	}

	_, err := client.Request("POST", debugPath, &opt)
	return err
}

func resourceDataServiceApiAuthActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = doApiAuthAction(client, d)
	if err != nil {
		return diag.Errorf("failed to operate API (%s): %s", d.Get("api_id").(string), err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceDataServiceApiAuthActionRead(ctx, d, meta)
}

func resourceDataServiceApiAuthActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is only a one-time action resource for operating the (authorized) APP.
	return nil
}

func resourceDataServiceApiAuthActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the API. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
