// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GA
// ---------------------------------------------------------------

package ga

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
			"frozen_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The frozen details of cloud services or resources.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of a cloud service or resource.`,
						},
						"effect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the resource after being forzen.`,
						},
						"scene": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The service scenario.`,
						},
					},
				},
			},
		},
	}
}

func resourceHealthCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v1/health-checks"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateHealthCheckBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating GA health check: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	healthCheckId := utils.PathSearch("health_check.id", respBody, "").(string)
	if healthCheckId == "" {
		return diag.Errorf("error creating GA health check: ID is not found in API response")
	}
	d.SetId(healthCheckId)

	err = createHealthCheckWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the GA health check (%s) creation to complete: %s", d.Id(), err)
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
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/health-checks/{health_check_id}"
		product          = "ga"
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{health_check_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`health_check.status`, respBody, "").(string)
			if utils.StrSliceContains(targetStatus, status) {
				return respBody, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceHealthCheckRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/health-checks/{health_check_id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{health_check_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GA health check")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("created_at", utils.PathSearch("health_check.created_at", respBody, nil)),
		d.Set("enabled", utils.PathSearch("health_check.enabled", respBody, nil)),
		d.Set("endpoint_group_id", utils.PathSearch("health_check.endpoint_group_id", respBody, nil)),
		d.Set("interval", utils.PathSearch("health_check.interval", respBody, nil)),
		d.Set("max_retries", utils.PathSearch("health_check.max_retries", respBody, nil)),
		d.Set("port", utils.PathSearch("health_check.port", respBody, nil)),
		d.Set("protocol", utils.PathSearch("health_check.protocol", respBody, nil)),
		d.Set("status", utils.PathSearch("health_check.status", respBody, nil)),
		d.Set("timeout", utils.PathSearch("health_check.timeout", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("health_check.updated_at", respBody, nil)),
		d.Set("frozen_info", flattenHealthCheckFrozenInfo(utils.PathSearch("health_check.frozen_info", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHealthCheckFrozenInfo(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	frozenInfo := map[string]interface{}{
		"status": utils.PathSearch("status", resp, nil),
		"effect": utils.PathSearch("effect", resp, nil),
		"scene":  utils.PathSearch("scene", resp, []string{}),
	}

	return []map[string]interface{}{frozenInfo}
}

func resourceHealthCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	updateHealthCheckHasChanges := []string{
		"enabled",
		"interval",
		"max_retries",
		"port",
		"protocol",
		"timeout",
	}

	if d.HasChanges(updateHealthCheckHasChanges...) {
		var (
			httpUrl = "v1/health-checks/{health_check_id}"
			product = "ga"
		)
		client, err := conf.NewServiceClient(product, region)
		if err != nil {
			return diag.Errorf("error creating GA client: %s", err)
		}

		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{health_check_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateHealthCheckBodyParams(d)),
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating GA health check: %s", err)
		}

		err = updateHealthCheckWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the GA health check (%s) update to complete: %s", d.Id(), err)
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
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/health-checks/{health_check_id}"
		product          = "ga"
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{health_check_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`health_check.status`, respBody, "").(string)
			if utils.StrSliceContains(targetStatus, status) {
				return respBody, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceHealthCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v1/health-checks/{health_check_id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{health_check_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting GA health check: %s", err)
	}

	err = deleteHealthCheckWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the GA health check (%s) delete to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteHealthCheckWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/health-checks/{health_check_id}"
		product          = "ga"
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{health_check_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					// When the error code is `404`, the value of respBody is nil, and a non-null value is returned to
					// avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`health_check.status`, respBody, "").(string)
			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}
