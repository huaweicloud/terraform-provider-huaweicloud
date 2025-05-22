package dew

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var keypairAssociateNonUpdatableParams = []string{
	"keypair_name",
	"server",
	"server.*.id",
	"server.*.auth",
	"server.*.port",
	"server.*.disable_password",
	"server.*.auth.*.type",
	"server.*.auth.*.key",
}

// @API DEW POST /v3/{project_id}/keypairs/associate
// @API DEW GET /v3/{project_id}/tasks/{task_id}
func ResourceKmsKeypairAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsKeypairAssociateCreate,
		ReadContext:   resourceKmsKeypairAssociateRead,
		UpdateContext: resourceKmsKeypairAssociateUpdate,
		DeleteContext: resourceKmsKeypairAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(keypairAssociateNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"keypair_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of SSH keypair.`,
			},
			"server": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        resourceKeypairAssociateServerSchema(),
				Description: `Specifies the ECS information.`,
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

func resourceKeypairAssociateServerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies ID of the ECS.`,
			},
			"auth": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        resourceKeypairAssociateServerAuthSchema(),
				Description: `Specifies the authentication information.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the SSH listening port.`,
			},
			"disable_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the password is disabled.`,
			},
		},
	}
}

func resourceKeypairAssociateServerAuthSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the value of the authentication type.`,
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the value of the key.`,
			},
		},
	}
}

func buildKeypairAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"keypair_name": d.Get("keypair_name"),
		"server":       utils.RemoveNil(buildKeypairAssociateServerBodyParams(d.Get("server").([]interface{}))),
	}
}

func buildKeypairAssociateServerAuthBodyParams(auths []interface{}) map[string]interface{} {
	if len(auths) == 0 || auths[0] == nil {
		return nil
	}

	auth, ok := auths[0].(map[string]interface{})
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"type": utils.ValueIgnoreEmpty(auth["type"]),
		"key":  utils.ValueIgnoreEmpty(auth["key"]),
	}

	return bodyParams
}

func buildKeypairAssociateServerBodyParams(servers []interface{}) map[string]interface{} {
	if len(servers) == 0 || servers[0] == nil {
		return nil
	}

	bodyParams := make(map[string]interface{}, 0)
	for _, v := range servers {
		item := v.(map[string]interface{})
		bodyParams = map[string]interface{}{
			"id":               item["id"],
			"auth":             buildKeypairAssociateServerAuthBodyParams(item["auth"].([]interface{})),
			"port":             utils.ValueIgnoreEmpty(item["port"]),
			"disable_password": item["disable_password"],
		}
	}

	return bodyParams
}

func getKeypairAssociateTask(client *golangsdk.ServiceClient, taskId string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/tasks/{task_id}"
		getOpt  = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return getResp, fmt.Errorf("error retrieving KMS keypair associate task: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func keypairAssociateTaskStatusRefreshFunc(taskId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getTaskBody, err := getKeypairAssociateTask(client, taskId)
		if err != nil {
			return getTaskBody, "ERROR", err
		}

		taskStatus := utils.PathSearch("task_status", getTaskBody, "").(string)
		if utils.StrSliceContains([]string{"FAILED_BIND", "FAILED_RESET", "FAILED_REPLACE"}, taskStatus) {
			return getTaskBody, "ERROR", fmt.Errorf("unexpect status (%s)", taskStatus)
		}

		if utils.StrSliceContains([]string{"SUCCESS_BIND", "SUCCESS_RESET", "SUCCESS_REPLACE"}, taskStatus) {
			return getTaskBody, "COMPLETED", nil
		}

		return getTaskBody, "PENDING", nil
	}
}

func waitForKeypairAssociateStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, taskId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      keypairAssociateTaskStatusRefreshFunc(taskId, client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func resourceKmsKeypairAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/keypairs/associate"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildKeypairAssociateBodyParams(d)),
	}
	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error associating the KMS keypair to the ECS: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("task_id", respBody, "").(string)
	if taskId == "" {
		return diag.Errorf("unable to find KMS task ID from API response")
	}
	d.SetId(taskId)

	err = waitForKeypairAssociateStatusCompleted(ctx, client, taskId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the task (%s) completed: %s", taskId, err)
	}

	return nil
}

func resourceKmsKeypairAssociateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsKeypairAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsKeypairAssociateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to associate a SSH keypair to a specified ECS.
Deleting this resource will not change the current SSH key pair, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
