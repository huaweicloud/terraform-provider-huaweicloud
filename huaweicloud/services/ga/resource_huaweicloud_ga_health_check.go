// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GA
// ---------------------------------------------------------------

package ga

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA POST /v1/health-checks
// @API GA GET /v1/health-checks/{health_check_id}
// @API GA PUT /v1/health-checks/{health_check_id}
// @API GA DELETE /v1/health-checks/{health_check_id}
func ResourceHealthCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHealthCheckCreate,
		UpdateContext: resourceHealthCheckUpdate,
		ReadContext:   resourceHealthCheckRead,
		DeleteContext: resourceHealthCheckDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether to enable health check.`,
			},
			"endpoint_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the endpoint group ID.`,
			},
			"interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the health check interval, in seconds.`,
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the maximum number of retries.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the port used for the health check.`,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the timeout duration of the health check, in seconds.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TCP",
				Description: `Specifies the health check protocol.`,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP",
				}, false),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the provisioning status.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates when the health check was configured.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates when the health check was updated. `,
			},
		},
	}
}

func resourceHealthCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// createHealthCheck: Create a GA Health Check.
	var (
		createHealthCheckHttpUrl = "v1/health-checks"
		createHealthCheckProduct = "ga"
	)
	createHealthCheckClient, err := conf.NewServiceClient(createHealthCheckProduct, region)
	if err != nil {
		return diag.Errorf("error creating HealthCheck Client: %s", err)
	}

	createHealthCheckPath := createHealthCheckClient.Endpoint + createHealthCheckHttpUrl

	createHealthCheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createHealthCheckOpt.JSONBody = utils.RemoveNil(buildCreateHealthCheckBodyParams(d))
	createHealthCheckResp, err := createHealthCheckClient.Request("POST", createHealthCheckPath, &createHealthCheckOpt)
	if err != nil {
		return diag.Errorf("error creating HealthCheck: %s", err)
	}

	createHealthCheckRespBody, err := utils.FlattenResponse(createHealthCheckResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("health_check.id", createHealthCheckRespBody)
	if err != nil {
		return diag.Errorf("error creating HealthCheck: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = createHealthCheckWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of HealthCheck (%s) to complete: %s", d.Id(), err)
	}
	return resourceHealthCheckRead(ctx, d, meta)
}

func buildCreateHealthCheckBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"health_check": map[string]interface{}{
			"enabled":           utils.ValueIgnoreEmpty(d.Get("enabled")),
			"endpoint_group_id": utils.ValueIgnoreEmpty(d.Get("endpoint_group_id")),
			"interval":          utils.ValueIgnoreEmpty(d.Get("interval")),
			"max_retries":       utils.ValueIgnoreEmpty(d.Get("max_retries")),
			"port":              utils.ValueIgnoreEmpty(d.Get("port")),
			"protocol":          utils.ValueIgnoreEmpty(d.Get("protocol")),
			"timeout":           utils.ValueIgnoreEmpty(d.Get("timeout")),
		},
	}
	return bodyParams
}

func createHealthCheckWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// createHealthCheckWaiting: missing operation notes
			var (
				createHealthCheckWaitingHttpUrl = "v1/health-checks/{health_check_id}"
				createHealthCheckWaitingProduct = "ga"
			)
			createHealthCheckWaitingClient, err := config.NewServiceClient(createHealthCheckWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating HealthCheck Client: %s", err)
			}

			createHealthCheckWaitingPath := createHealthCheckWaitingClient.Endpoint + createHealthCheckWaitingHttpUrl
			createHealthCheckWaitingPath = strings.ReplaceAll(createHealthCheckWaitingPath, "{health_check_id}", d.Id())

			createHealthCheckWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createHealthCheckWaitingResp, err := createHealthCheckWaitingClient.Request("GET",
				createHealthCheckWaitingPath, &createHealthCheckWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createHealthCheckWaitingRespBody, err := utils.FlattenResponse(createHealthCheckWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`health_check.status`, createHealthCheckWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `health_check.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createHealthCheckWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createHealthCheckWaitingRespBody, status, nil
			}

			return createHealthCheckWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceHealthCheckRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// getHealthCheck: Query the GA Health Check detail
	var (
		getHealthCheckHttpUrl = "v1/health-checks/{health_check_id}"
		getHealthCheckProduct = "ga"
	)
	getHealthCheckClient, err := conf.NewServiceClient(getHealthCheckProduct, region)
	if err != nil {
		return diag.Errorf("error creating HealthCheck Client: %s", err)
	}

	getHealthCheckPath := getHealthCheckClient.Endpoint + getHealthCheckHttpUrl
	getHealthCheckPath = strings.ReplaceAll(getHealthCheckPath, "{health_check_id}", d.Id())

	getHealthCheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getHealthCheckResp, err := getHealthCheckClient.Request("GET", getHealthCheckPath, &getHealthCheckOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving HealthCheck")
	}

	getHealthCheckRespBody, err := utils.FlattenResponse(getHealthCheckResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("created_at", utils.PathSearch("health_check.created_at", getHealthCheckRespBody, nil)),
		d.Set("enabled", utils.PathSearch("health_check.enabled", getHealthCheckRespBody, nil)),
		d.Set("endpoint_group_id", utils.PathSearch("health_check.endpoint_group_id", getHealthCheckRespBody, nil)),
		d.Set("interval", utils.PathSearch("health_check.interval", getHealthCheckRespBody, nil)),
		d.Set("max_retries", utils.PathSearch("health_check.max_retries", getHealthCheckRespBody, nil)),
		d.Set("port", utils.PathSearch("health_check.port", getHealthCheckRespBody, nil)),
		d.Set("protocol", utils.PathSearch("health_check.protocol", getHealthCheckRespBody, nil)),
		d.Set("status", utils.PathSearch("health_check.status", getHealthCheckRespBody, nil)),
		d.Set("timeout", utils.PathSearch("health_check.timeout", getHealthCheckRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("health_check.updated_at", getHealthCheckRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceHealthCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	updateHealthCheckhasChanges := []string{
		"enabled",
		"interval",
		"max_retries",
		"port",
		"protocol",
		"timeout",
	}

	if d.HasChanges(updateHealthCheckhasChanges...) {
		// updateHealthCheck: Update the configuration of GA Health Check
		var (
			updateHealthCheckHttpUrl = "v1/health-checks/{health_check_id}"
			updateHealthCheckProduct = "ga"
		)
		updateHealthCheckClient, err := conf.NewServiceClient(updateHealthCheckProduct, region)
		if err != nil {
			return diag.Errorf("error creating HealthCheck Client: %s", err)
		}

		updateHealthCheckPath := updateHealthCheckClient.Endpoint + updateHealthCheckHttpUrl
		updateHealthCheckPath = strings.ReplaceAll(updateHealthCheckPath, "{health_check_id}", d.Id())

		updateHealthCheckOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateHealthCheckOpt.JSONBody = utils.RemoveNil(buildUpdateHealthCheckBodyParams(d))
		_, err = updateHealthCheckClient.Request("PUT", updateHealthCheckPath, &updateHealthCheckOpt)
		if err != nil {
			return diag.Errorf("error updating HealthCheck: %s", err)
		}
		err = updateHealthCheckWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of HealthCheck (%s) to complete: %s", d.Id(), err)
		}
	}
	return resourceHealthCheckRead(ctx, d, meta)
}

func buildUpdateHealthCheckBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"health_check": map[string]interface{}{
			"enabled":     d.Get("enabled"),
			"interval":    utils.ValueIgnoreEmpty(d.Get("interval")),
			"max_retries": utils.ValueIgnoreEmpty(d.Get("max_retries")),
			"port":        utils.ValueIgnoreEmpty(d.Get("port")),
			"protocol":    utils.ValueIgnoreEmpty(d.Get("protocol")),
			"timeout":     utils.ValueIgnoreEmpty(d.Get("timeout")),
		},
	}
	return bodyParams
}

func updateHealthCheckWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// updateHealthCheckWaiting: missing operation notes
			var (
				updateHealthCheckWaitingHttpUrl = "v1/health-checks/{health_check_id}"
				updateHealthCheckWaitingProduct = "ga"
			)
			updateHealthCheckWaitingClient, err := config.NewServiceClient(updateHealthCheckWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating HealthCheck Client: %s", err)
			}

			updateHealthCheckWaitingPath := updateHealthCheckWaitingClient.Endpoint + updateHealthCheckWaitingHttpUrl
			updateHealthCheckWaitingPath = strings.ReplaceAll(updateHealthCheckWaitingPath, "{health_check_id}", d.Id())

			updateHealthCheckWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			updateHealthCheckWaitingResp, err := updateHealthCheckWaitingClient.Request("GET",
				updateHealthCheckWaitingPath, &updateHealthCheckWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			updateHealthCheckWaitingRespBody, err := utils.FlattenResponse(updateHealthCheckWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`health_check.status`, updateHealthCheckWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `health_check.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return updateHealthCheckWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return updateHealthCheckWaitingRespBody, status, nil
			}

			return updateHealthCheckWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceHealthCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// deleteHealthCheck: Delete an existing GA Health Check
	var (
		deleteHealthCheckHttpUrl = "v1/health-checks/{health_check_id}"
		deleteHealthCheckProduct = "ga"
	)
	deleteHealthCheckClient, err := conf.NewServiceClient(deleteHealthCheckProduct, region)
	if err != nil {
		return diag.Errorf("error creating HealthCheck Client: %s", err)
	}

	deleteHealthCheckPath := deleteHealthCheckClient.Endpoint + deleteHealthCheckHttpUrl
	deleteHealthCheckPath = strings.ReplaceAll(deleteHealthCheckPath, "{health_check_id}", d.Id())

	deleteHealthCheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteHealthCheckClient.Request("DELETE", deleteHealthCheckPath, &deleteHealthCheckOpt)
	if err != nil {
		return diag.Errorf("error deleting HealthCheck: %s", err)
	}

	err = deleteHealthCheckWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Delete of HealthCheck (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteHealthCheckWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// deleteHealthCheckWaiting: missing operation notes
			var (
				deleteHealthCheckWaitingHttpUrl = "v1/health-checks/{health_check_id}"
				deleteHealthCheckWaitingProduct = "ga"
			)
			deleteHealthCheckWaitingClient, err := config.NewServiceClient(deleteHealthCheckWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating HealthCheck Client: %s", err)
			}

			deleteHealthCheckWaitingPath := deleteHealthCheckWaitingClient.Endpoint + deleteHealthCheckWaitingHttpUrl
			deleteHealthCheckWaitingPath = strings.ReplaceAll(deleteHealthCheckWaitingPath, "{health_check_id}", d.Id())

			deleteHealthCheckWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteHealthCheckWaitingResp, err := deleteHealthCheckWaitingClient.Request("GET",
				deleteHealthCheckWaitingPath, &deleteHealthCheckWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteHealthCheckWaitingRespBody, err := utils.FlattenResponse(deleteHealthCheckWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`health_check.status`, deleteHealthCheckWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `health_check.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return deleteHealthCheckWaitingRespBody, status, nil
			}

			return deleteHealthCheckWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
