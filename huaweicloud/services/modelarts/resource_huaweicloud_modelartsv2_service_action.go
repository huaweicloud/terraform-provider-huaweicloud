package modelarts

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var v2ServiceActionNonUpdatableParams = []string{
	"service_id",
	"action",
}

// @API ModelArts POST /v2/{project_id}/services/{service_id}/start
// @API ModelArts POST /v2/{project_id}/services/{service_id}/stop
// @API ModelArts POST /v2/{project_id}/services/{service_id}/interrupt
func ResourceV2ServiceAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2ServiceActionCreate,
		ReadContext:   resourceV2ServiceActionRead,
		UpdateContext: resourceV2ServiceActionUpdate,
		DeleteContext: resourceV2ServiceActionDelete,

		CustomizeDiff: config.FlexibleForceNew(v2ServiceActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the service is located.`,
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The service ID to be operated.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action type.`,
			},
		},
	}
}

func resourceV2ServiceActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v2/{project_id}/services/{service_id}/{action}"
		serviceId = d.Get("service_id").(string)
		action    = d.Get("action").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{service_id}", serviceId)
	actionPath = strings.ReplaceAll(actionPath, "{action}", action)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("POST", actionPath, &opt)
	if err != nil {
		return diag.Errorf("unable to operate the service (%s), the action is: %s, the error is: %s",
			serviceId, action, err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	return nil
}

func resourceV2ServiceActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2ServiceActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2ServiceActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the service. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
