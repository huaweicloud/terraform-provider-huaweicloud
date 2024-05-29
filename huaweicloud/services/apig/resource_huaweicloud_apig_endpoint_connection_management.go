package apig

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/vpc-endpoint/connections/action
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/vpc-endpoint/connections
func ResourceEndpointConnectionManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointConnectionManagementCreate,
		ReadContext:   resourceEndpointConnectionManagementRead,
		UpdateContext: resourceEndpointConnectionManagementUpdate,
		DeleteContext: resourceEndpointConnectionManagementDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},

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
				ForceNew:    true,
				Description: "Specifies the ID of the dedicated instance to which the endpoint connection belongs.",
			},
			"endpoint_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the endpoint connection.",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the operation type endpoint connection.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current ststus of the endpoint connection.",
			},
		},
	}
}

func resourceEndpointConnectionManagementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("apig", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG Client: %s", err)
	}

	err = approveEndpointConnection(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("endpoint_id").(string))
	return resourceEndpointConnectionManagementRead(ctx, d, meta)
}

func approveEndpointConnection(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout time.Duration) error {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/vpc-endpoint/connections/action"
		instanceId = d.Get("instance_id").(string)
		endpointId = d.Get("endpoint_id").(string)
		action     = d.Get("action").(string)
	)
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{instance_id}", instanceId)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"action":    action,
			"endpoints": []string{endpointId},
		},
	}
	_, err := client.Request("POST", path, &opts)
	if err != nil {
		return fmt.Errorf("failed to %s endpoint connection (%s) under dedicated instance (%s): %s", action, endpointId, instanceId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: waitForConnectionApproved(client, instanceId, endpointId),
		Timeout: timeout,
		// The approve operation has an intermediate status of creating, so it needs to wait for a short period of time.
		MinTimeout: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for endpoint connection operations completed: %s", err)
	}

	return nil
}

func analyseApproveAction(action string) string {
	actionStatus := map[string]string{
		"accepted": "receive",
		"rejected": "reject",
	}
	if v, ok := actionStatus[action]; ok {
		return v
	}
	return ""
}

func resourceEndpointConnectionManagementRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG Client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		endpointId = d.Get("endpoint_id").(string)
	)
	connection, err := GetEndpointConntionByEndpointId(client, instanceId, endpointId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting endpoint connection (%s) under dedicated instance (%s)",
			endpointId, instanceId))
	}

	status := utils.PathSearch("status", connection, "")
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("status", status),
		d.Set("action", analyseApproveAction(status.(string))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEndpointConnectionManagementUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("apig", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG Client: %s", err)
	}

	if d.HasChange("action") {
		err = approveEndpointConnection(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceEndpointConnectionManagementRead(ctx, d, meta)
}

// Destroying resources does not change the current action of the endpoint connection.
func resourceEndpointConnectionManagementDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

// GetEndpointConntionByEndpointId is a method that used to obtain details of connection by specified endpoint connection ID.
func GetEndpointConntionByEndpointId(client *golangsdk.ServiceClient, instanceId, endpointId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/vpc-endpoint/connections"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath += fmt.Sprintf("?id=%s", endpointId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	connection := utils.PathSearch("connections[0]", respBody, nil)
	if connection != nil {
		return connection, nil
	}

	return nil, golangsdk.ErrDefault404{}
}

func waitForConnectionApproved(client *golangsdk.ServiceClient, instanceId, endpointId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		connection, err := GetEndpointConntionByEndpointId(client, instanceId, endpointId)
		if err != nil {
			return connection, "ERROR", err
		}

		status := utils.PathSearch("status", connection, "").(string)
		if utils.StrSliceContains([]string{"failed"}, status) {
			return connection, "", fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains([]string{"accepted", "rejected"}, status) {
			return connection, "COMPLETED", nil
		}
		return connection, "PENDING", nil
	}
}
