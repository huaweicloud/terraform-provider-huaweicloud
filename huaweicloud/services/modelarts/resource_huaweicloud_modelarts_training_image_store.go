package modelarts

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var trainingImageStoreNonUpdatableParams = []string{
	"training_job_id",
	"task_id",
	"name",
	"namespace",
	"tag",
	"description",
}

// @API ModelArts POST /v2/{project_id}/training-jobs/{training_job_id}/tasks/{task_id}/save-image-job
// @API ModelArts GET /v2/{project_id}/training-jobs/{training_job_id}/tasks/{task_id}/save-image-job
// @API SWR DELETE /v2/manage/namespaces/{namespace}/repos/{repository}/tags/{tag}
func ResourceTrainingImageStore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrainingImageStoreCreate,
		ReadContext:   resourceTrainingImageStoreRead,
		UpdateContext: resourceTrainingImageStoreUpdate,
		DeleteContext: resourceTrainingImageStoreDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(trainingImageStoreNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the training job to be stored as image is located.`,
			},

			// Required parameters.
			"training_job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the training job.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The task name of the training job.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the SWR repository to which the image is stored.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the SWR organization to which the image is stored.`,
			},
			"tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The tag of the image.`,
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the image.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
						Required: true,
					}),
			},
		},
	}
}

func buildTrainingImageStoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"namespace":   d.Get("namespace"),
		"tag":         d.Get("tag"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func createTrainingImageStore(client *golangsdk.ServiceClient, trainingJobId, taskId string, d *schema.ResourceData) error {
	httpURL := "v2/{project_id}/training-jobs/{training_job_id}/tasks/{task_id}/save-image-job"
	createPath := client.Endpoint + httpURL
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{training_job_id}", trainingJobId)
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildTrainingImageStoreBodyParams(d)),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func getTrainingImageStoreById(client *golangsdk.ServiceClient, trainingJobId, taskId string) (interface{}, error) {
	httpURL := "v2/{project_id}/training-jobs/{training_job_id}/tasks/{task_id}/save-image-job"
	getPath := client.Endpoint + httpURL
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{training_job_id}", trainingJobId)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)

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

func refreshTrainingImageStoreStatusFunc(client *golangsdk.ServiceClient, trainingJobId, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getTrainingImageStoreById(client, trainingJobId, taskId)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		if status == "CREATE_FAILED" {
			return respBody, "ERROR", fmt.Errorf("unexpected status (%s)", status)
		}

		if status == "ACTIVE" {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func waitForTrainingImageStoreCompleted(ctx context.Context, client *golangsdk.ServiceClient, trainingJobId, taskId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshTrainingImageStoreStatusFunc(client, trainingJobId, taskId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceTrainingImageStoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		trainingJobId = d.Get("training_job_id").(string)
		taskId        = d.Get("task_id").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = createTrainingImageStore(client, trainingJobId, taskId, d)
	if err != nil {
		return diag.Errorf("error saving image from training job (%s): %s", trainingJobId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", trainingJobId, taskId, d.Get("tag").(string)))

	if err = waitForTrainingImageStoreCompleted(ctx, client, trainingJobId, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for image saved from training job (%s) to complete: %s", trainingJobId, err)
	}

	return resourceTrainingImageStoreRead(ctx, d, meta)
}

func resourceTrainingImageStoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTrainingImageStoreUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func deleteStoredImageRepositoryTag(client *golangsdk.ServiceClient, namespace, repositoryName, tag string) error {
	httpURL := "v2/manage/namespaces/{namespace}/repos/{repository}/tags/{tag}"
	deletePath := client.Endpoint + httpURL
	deletePath = strings.ReplaceAll(deletePath, "{namespace}", namespace)
	deletePath = strings.ReplaceAll(deletePath, "{repository}", repositoryName)
	deletePath = strings.ReplaceAll(deletePath, "{tag}", tag)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceTrainingImageStoreDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		namespace      = d.Get("namespace").(string)
		repositoryName = d.Get("name").(string)
		tag            = d.Get("tag").(string)
	)

	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	err = deleteStoredImageRepositoryTag(client, namespace, repositoryName, tag)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error deleting stored image tag (%s/%s:%s) from training job (%s)", namespace, repositoryName, tag,
				d.Get("training_job_id").(string)))
	}

	return nil
}
