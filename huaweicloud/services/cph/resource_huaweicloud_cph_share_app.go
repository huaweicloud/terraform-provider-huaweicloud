package cph

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

var shareAppNonUpdatableParams = []string{
	"server_id",
	"package_name",
	"bucket_name",
	"object_path",
	"pre_install_app",
}

// @API CPH POST /v1/{project_id}/cloud-phone/phones/share-apps
// @API CPH DELETE /v1/{project_id}/cloud-phone/phones/share-apps
// @API CPH GET /v1/{project_id}/cloud-phone/jobs/{job_id}
func ResourceShareApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceShareAppCreate,
		UpdateContext: resourceShareAppUpdate,
		ReadContext:   resourceShareAppRead,
		DeleteContext: resourceShareAppDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(shareAppNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CPH server ID.`,
			},
			"package_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the package name.`,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the OBS bucket name.`,
			},
			"object_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the OBS object path.`,
			},
			"pre_install_app": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies whether to pre-install the application.`,
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

func resourceShareAppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("cph", region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	// createShareApp: create CPH share app
	createShareAppHttpUrl := "v1/{project_id}/cloud-phone/phones/share-apps"
	createShareAppPath := client.Endpoint + createShareAppHttpUrl
	createShareAppPath = strings.ReplaceAll(createShareAppPath, "{project_id}", client.ProjectID)

	createShareAppOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createShareAppOpt.JSONBody = buildcreateShareAppBodyParams(d)
	createAdbShareAppResp, err := client.Request("POST", createShareAppPath, &createShareAppOpt)
	if err != nil {
		return diag.Errorf("error creating CPH share app: %s", err)
	}

	createShareAppRespBody, err := utils.FlattenResponse(createAdbShareAppResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("jobs|[0].job_id", createShareAppRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the pushing share app job ID from the API response")
	}
	d.SetId(jobId)

	err = checkShareAppJobStatus(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for pushing CPH share app to be completed: %s", err)
	}

	return nil
}

func buildcreateShareAppBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"package_name":    d.Get("package_name"),
		"pre_install_app": d.Get("pre_install_app"),
		"bucket_name":     d.Get("bucket_name"),
		"object_path":     d.Get("object_path"),
		"server_ids":      []string{d.Get("server_id").(string)},
	}

	return bodyParams
}

func resourceShareAppRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceShareAppUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceShareAppDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("cph", region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	// deleteShareApp: delete CPH share app
	deleteShareAppHttpUrl := "v1/{project_id}/cloud-phone/phones/share-apps"
	deleteShareAppPath := client.Endpoint + deleteShareAppHttpUrl
	deleteShareAppPath = strings.ReplaceAll(deleteShareAppPath, "{project_id}", client.ProjectID)

	deleteShareAppOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteShareAppOpt.JSONBody = map[string]interface{}{
		"package_name": d.Get("package_name"),
		"server_ids":   []string{d.Get("server_id").(string)},
	}
	deleteAdbShareAppResp, err := client.Request("DELETE", deleteShareAppPath, &deleteShareAppOpt)
	if err != nil {
		return diag.Errorf("error deleting CPH share app: %s", err)
	}

	deleteShareAppRespBody, err := utils.FlattenResponse(deleteAdbShareAppResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("jobs|[0].job_id", deleteShareAppRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the deleting share app job ID from the API response")
	}

	err = checkShareAppJobStatus(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for deleting CPH share app to be completed: %s", err)
	}

	return nil
}

func checkShareAppJobStatus(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      jobStatusRefreshFunc(client, id),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func jobStatusRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getJobsHttpUrl := "v1/{project_id}/cloud-phone/jobs/{job_id}"
		getJobsPath := client.Endpoint + getJobsHttpUrl
		getJobsPath = strings.ReplaceAll(getJobsPath, "{project_id}", client.ProjectID)
		getJobsPath = strings.ReplaceAll(getJobsPath, "{job_id}", id)

		getJobsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getJobsResp, err := client.Request("GET", getJobsPath, &getJobsOpt)
		if err != nil {
			return nil, "ERROR", err
		}

		getJobsRespBody, err := utils.FlattenResponse(getJobsResp)
		if err != nil {
			return nil, "ERROR", err
		}

		// status is 1, indicates the job is Running
		// status is -1, indicates the job is failed
		status := utils.PathSearch("status", getJobsRespBody, float64(0)).(float64)
		if int(status) == 1 {
			return getJobsRespBody, "PENDING", nil
		}
		if int(status) == -1 {
			return getJobsRespBody, "ERROR", fmt.Errorf("error_msg: %v", utils.PathSearch("error_msg", getJobsRespBody, "").(string))
		}
		return getJobsRespBody, "COMPLETED", nil
	}
}
