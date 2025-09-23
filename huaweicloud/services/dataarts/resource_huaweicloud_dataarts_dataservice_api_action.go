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

// @API DataArtsStudio POST /v1/{project_id}/service/apis/{api_id}/instances/{instance_id}/action
func ResourceDataServiceApiAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceApiActionCreate,
		ReadContext:   resourceDataServiceApiActionRead,
		DeleteContext: resourceDataServiceApiActionDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the API to be operated is located.`,
			},

			// Parameter in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the workspace to which the API belongs.`,
			},

			// Arguments
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The exclusive API ID, which in published status.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The exclusive cluster ID to which the published API belongs on Data Service side.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The action type.`,
			},
		},
	}
}

func doApiAction(client *golangsdk.ServiceClient, workspaceId, apiId, instanceId, actionType string) error {
	httpUrl := "v1/{project_id}/service/apis/{api_id}/instances/{instance_id}/action"
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{api_id}", apiId)
	actionPath = strings.ReplaceAll(actionPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     "EXCLUSIVE",
		},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			// The unpublish/stop function requires sufficient preparation time for authorized applications,
			// and the request must be initiated at least 2 days in advance.
			// The time must be on the hour.
			"time":   utils.CalculateNextWholeHourAfterFewTime(utils.GetCurrentTime(true), 48*time.Hour),
			"action": actionType,
		}),
		OkCodes: []int{204},
	}

	_, err := client.Request("POST", actionPath, &opt)
	if err != nil {
		return err
	}
	return nil
}

func resourceDataServiceApiActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		workspaceId = d.Get("workspace_id").(string)
		apiId       = d.Get("api_id").(string)
		instanceId  = d.Get("instance_id").(string)
		actionType  = d.Get("type").(string)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = doApiAction(client, workspaceId, apiId, instanceId, actionType)
	if err != nil {
		return diag.Errorf("failed to %s API (%s): %s", strings.ToLower(actionType), apiId, err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceDataServiceApiActionRead(ctx, d, meta)
}

func resourceDataServiceApiActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is only a one-time action resource for unpublishing/stopping/recovering the API.
	// There is no API for the provider to query the history of this API action.
	return nil
}

func resourceDataServiceApiActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for unpublishing/stopping/recovering the API.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
