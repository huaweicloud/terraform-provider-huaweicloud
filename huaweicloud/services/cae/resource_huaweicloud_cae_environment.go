package cae

import (
	"context"
	"fmt"
	"log"
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

var envResourceNotFoundCodes = []string{"CAE.01500404"}

// @API CAE POST /v1/{project_id}/cae/environments
// @API CAE GET /v1/{project_id}/cae/jobs/{job_id}
// @API CAE POST /v1/{project_id}/cae/jobs/{job_id}
// @API CAE GET /v1/{project_id}/cae/environments
// @API CAE DELETE /v1/{project_id}/cae/environments/{environment_id}
func ResourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the environment is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the environment.",
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					oldVal, newVal := d.GetChange("annotations")
					for key, value := range newVal.(map[string]interface{}) {
						if mapValue, exists := oldVal.(map[string]interface{})[key]; exists && mapValue == value {
							continue
						}
						return false
					}
					return true
				},
				Description: utils.SchemaDesc(
					"The additional attributes of the environment.",
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the enterprise project to which the environment belongs.",
			},
			// Sometimes, debug requests may be affected by network fluctuations and time out.
			// Retry as appropriate to resolve the issue.
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum retry number in the create or delete operation.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the environment.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the environment, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the environment, in RFC3339 format.",
			},
		},
	}
}

func buildCreateEnvironmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "Environment",
		"metadata": map[string]interface{}{
			"annotations": utils.ValueIgnoreEmpty(d.Get("annotations")),
			"name":        d.Get("name"),
		},
	}
}

func buildEnvRequestMoreHeaders(epsId string) map[string]string {
	moreHeaders := map[string]string{
		"Content-Type": "application/json",
	}
	if epsId != "" {
		moreHeaders["X-Enterprise-Project-ID"] = epsId
	}
	return moreHeaders
}

func getJobById(client *golangsdk.ServiceClient, epsId, jobId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/jobs/{job_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildEnvRequestMoreHeaders(epsId),
	}
	createResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error querying job by its ID (%s): %s", jobId, err)
	}
	return utils.FlattenResponse(createResp)
}

func environmentJobRefreshFunc(client *golangsdk.ServiceClient, epsId, jobId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getJobById(client, epsId, jobId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "not_found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		status := utils.PathSearch("spec.status", respBody, "").(string)
		if utils.StrSliceContains([]string{"failed", "timeout"}, status) {
			return respBody, "ERROR", fmt.Errorf("unexpect job status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func retryJob(client *golangsdk.ServiceClient, epsId, jobId string) error {
	httpUrl := "v1/{project_id}/cae/jobs/{job_id}"

	retryPath := client.Endpoint + httpUrl
	retryPath = strings.ReplaceAll(retryPath, "{project_id}", client.ProjectID)
	retryPath = strings.ReplaceAll(retryPath, "{job_id}", jobId)

	retryOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildEnvRequestMoreHeaders(epsId),
	}
	_, err := client.Request("POST", retryPath, &retryOpts)
	if err != nil {
		return fmt.Errorf("error retring failed job (%s): %s", jobId, err)
	}
	return nil
}

func buildCurrentRetryCount(num int) string {
	if num < 1 {
		return "Retry number is invalid"
	}

	var result string
	switch num {
	case 1:
		result = "1st"
	case 2:
		result = "2nd"
	default:
		result = fmt.Sprintf("%dth", num)
	}
	return result
}

func waitForEnvironmentJobComplete(ctx context.Context, client *golangsdk.ServiceClient,
	epsId, jobId string, timeout time.Duration, maxRetries int) error {
	var totalTryNum = 0

	for {
		totalTryNum++
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      environmentJobRefreshFunc(client, epsId, jobId, []string{"success"}),
			Timeout:      timeout,
			Delay:        5 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err := stateConf.WaitForStateContext(ctx)
		if err == nil {
			break
		}
		if totalTryNum > maxRetries {
			return fmt.Errorf("error waiting for the job status to become expect value: %s", err)
		}
		log.Printf("[DEBUG][%s retry] Prepare to retry the failed job (%s)", buildCurrentRetryCount(totalTryNum), jobId)
		err = retryJob(client, epsId, jobId)
		if err != nil {
			log.Printf("[DEBUG][%s retry] An error occurred while retring failed job (%s): %s",
				buildCurrentRetryCount(totalTryNum), jobId, err)
		}
	}

	return nil
}

func createEnvironment(client *golangsdk.ServiceClient, d *schema.ResourceData, epsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/environments"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildEnvRequestMoreHeaders(epsId),
		JSONBody:         utils.RemoveNil(buildCreateEnvironmentBodyParams(d)),
	}
	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return nil, fmt.Errorf("error creating environment: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func getEnvironmentByName(client *golangsdk.ServiceClient, epsId, envName string) (interface{}, error) {
	envList, err := getEnvironments(client, epsId)
	if err != nil {
		return nil, err
	}

	env := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", envName), envList, nil)
	if env == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/cae/environments",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the environment (%s) does not exist", envName)),
			},
		}
	}
	return env, nil
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		envName = d.Get("name").(string)
		epsId   = cfg.GetEnterpriseProjectID(d)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	respBody, err := createEnvironment(client, d, epsId)
	if err != nil {
		return diag.FromErr(err)
	}
	// Use the environment name as the temporary ID during resource creation.
	d.SetId(envName)

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}
	err = waitForEnvironmentJobComplete(ctx, client, epsId, jobId, d.Timeout(schema.TimeoutCreate), d.Get("max_retries").(int))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func getEnvironments(client *golangsdk.ServiceClient, epsId string) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/cae/environments"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildEnvRequestMoreHeaders(epsId),
	}
	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func GetEnvironmentById(client *golangsdk.ServiceClient, epsId, resourceId string) (interface{}, error) {
	envList, err := getEnvironments(client, epsId)
	if err != nil {
		return nil, err
	}

	env := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", resourceId), envList, nil)
	if env == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/cae/environments",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the environment (%s) does not exist", resourceId)),
			},
		}
	}
	return env, nil
}

func resourceEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		resourceId = d.Id()
		env        interface{}
		epsId      = cfg.GetEnterpriseProjectID(d)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	if utils.IsUUID(resourceId) {
		env, err = GetEnvironmentById(client, epsId, resourceId)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error querying environment")
		}
	} else {
		// Since the creation interface does not return the environment ID, the temporary ID needs to be refreshed to
		// the real environment ID regardless of whether the creation is completed.
		env, err = getEnvironmentByName(client, epsId, resourceId)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error querying environment")
		}
		envId := utils.PathSearch("id", env, "").(string)
		if envId == "" {
			return diag.Errorf("the environment do not have the ID attribute, please re-import this resource")
		}
		d.SetId(envId)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", env, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("annotations.enterprise_project_id", env, nil)),
		d.Set("annotations", flattenAnnotations(utils.PathSearch("annotations",
			env, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("status", utils.PathSearch("status", env, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at",
			env, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_at",
			env, "").(string))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAnnotations(annotations map[string]interface{}) map[string]interface{} {
	if len(annotations) < 1 {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range annotations {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Only refresh the number of max retries.
	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/cae/environments/{environment_id}"
		epsId   = cfg.GetEnterpriseProjectID(d)
		envId   = d.Id()
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{environment_id}", envId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildEnvRequestMoreHeaders(epsId),
	}
	requestResp, err := client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", envResourceNotFoundCodes...),
			fmt.Sprintf("error deleting environment (%s)", envId))
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}
	err = waitForEnvironmentJobComplete(ctx, client, epsId, jobId, d.Timeout(schema.TimeoutDelete), d.Get("max_retries").(int))
	return diag.FromErr(err)
}
