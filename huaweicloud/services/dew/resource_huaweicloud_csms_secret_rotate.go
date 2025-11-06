package dew

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1/{project_id}/secrets/{secret_name}/rotate
// @API DEW GET /v1/{project_id}/csms/tasks
func ResourceSecretRotate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretRotateCreate,
		ReadContext:   resourceSecretRotateRead,
		UpdateContext: resourceSecretRotateUpdate,
		DeleteContext: resourceSecretRotateDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"secret_name",
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"secret_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the secret to rotate.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the rotation task.`,
			},
			"rotation_func_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URN of the rotation function.`,
			},
			"operate_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the operation.`,
			},
			"task_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The time when the rotation task was created, UNIX timestamp in milliseconds.`,
			},
			"attempt_nums": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of attempts to rotate the secret.`,
			},
		},
	}
}

func queryCsmsTargetTask(client *golangsdk.ServiceClient, taskId string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/csms/tasks"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?task_id=%s", taskId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	taskDetail := utils.PathSearch("tasks|[0]", respBody, nil)
	if taskDetail == nil {
		return nil, errors.New("unable to find the DEW CSMS secret rotation task detail from the API response")
	}

	return taskDetail, nil
}

func waitingForSecretRotateTaskSuccess(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			taskDetail, err := queryCsmsTargetTask(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			taskStatus := utils.PathSearch("task_status", taskDetail, "").(string)
			if taskStatus == "" {
				return nil, "ERROR", errors.New("unable to find the DEW CSMS secret rotation task status from the API response")
			}

			if taskStatus == "SUCCESS" {
				return taskDetail, "COMPLETED", nil
			}

			if taskStatus == "FAILED" {
				errCode := utils.PathSearch("task_error_code", taskDetail, "").(string)
				errMsg := utils.PathSearch("task_error_msg", taskDetail, "").(string)
				return taskDetail, "ERROR", fmt.Errorf("rotation task failed, error code: %s, error message: %s", errCode, errMsg)
			}

			return taskDetail, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func resourceSecretRotateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/secrets/{secret_name}/rotate"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{secret_name}", d.Get("secret_name").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error rotating DEW CSMS secret: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("rotation_task_id", respBody, "").(string)
	if taskId == "" {
		return diag.Errorf("unable to find the DEW CSMS secret rotation task ID from the API response")
	}

	d.SetId(taskId)

	taskDetail, err := waitingForSecretRotateTaskSuccess(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for DEW CSMS secret rotation task to be completed: %s", err)
	}

	mErr := multierror.Append(
		d.Set("task_id", taskId),
		d.Set("rotation_func_urn", utils.PathSearch("rotation_func_urn", taskDetail, nil)),
		d.Set("operate_type", utils.PathSearch("operate_type", taskDetail, nil)),
		d.Set("task_time", utils.PathSearch("task_time", taskDetail, nil)),
		d.Set("attempt_nums", utils.PathSearch("attempt_nums", taskDetail, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting DEW CSMS secret rotation task attributes: %s", err)
	}

	return resourceSecretRotateRead(ctx, d, meta)
}

func resourceSecretRotateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceSecretRotateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceSecretRotateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to rotate secret.
	Deleting this resource will not recover the rotated secret, but will only remove the resource information from
	the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
