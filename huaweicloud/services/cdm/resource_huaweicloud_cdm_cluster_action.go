// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDM
// ---------------------------------------------------------------

package cdm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDM POST /v1.1/{project_id}/clusters/{cluster_id}/action
// @API CDM GET /v1.1/{project_id}/clusters/{clusterId}
func ResourceClusterAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterActionCreate,
		ReadContext:   resourceClusterActionRead,
		DeleteContext: resourceClusterActionDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `ID of CDM cluster.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Action type.`,
			},
			"restart": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ClusterActionRestartSchema(),
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func ClusterActionRestartSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Restart level.`,
				ForceNew:    true,
			},
			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Restart mode.`,
				ForceNew:    true,
			},
			"delay_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Restart delay, in seconds.`,
				ForceNew:    true,
			},
		},
	}
	return &sc
}

func resourceClusterActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createClusterAction: create a cdm cluster action.
	var (
		createClusterActionHttpUrl = "v1.1/{project_id}/clusters/{cluster_id}/action"
		createClusterActionProduct = "cdm"
	)
	createClusterActionClient, err := cfg.NewServiceClient(createClusterActionProduct, region)
	if err != nil {
		return diag.Errorf("error creating CDM client: %s", err)
	}

	createClusterActionPath := createClusterActionClient.Endpoint + createClusterActionHttpUrl
	createClusterActionPath = strings.ReplaceAll(createClusterActionPath, "{project_id}", createClusterActionClient.ProjectID)
	createClusterActionPath = strings.ReplaceAll(createClusterActionPath, "{cluster_id}", d.Get("cluster_id").(string))

	createClusterActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createClusterActionOpt.JSONBody = buildCreateClusterActionBodyParams(d)
	createClusterActionResp, err := createClusterActionClient.Request("POST", createClusterActionPath, &createClusterActionOpt)

	if err != nil {
		return diag.Errorf("error creating CDM cluster action: %s", err)
	}

	createClusterActionRespBody, err := utils.FlattenResponse(createClusterActionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("jobId[0]", createClusterActionRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating cluster action: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = createClusterActionWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the cluster action (%s) creation to complete: %s", d.Id(), err)
	}
	return nil
}

func buildCreateClusterActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	if v := d.Get("type").(string); v == "start" {
		return map[string]interface{}{
			"start": map[string]interface{}{},
		}
	}

	return map[string]interface{}{
		"restart": utils.RemoveNil(buildCreateClusterActionRequestBodyRestart(d.Get("restart"))),
	}
}

func buildCreateClusterActionRequestBodyRestart(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"restartDelayTime": utils.ValueIgnoreEmpty(raw["delay_time"]),
			"restartMode":      utils.ValueIgnoreEmpty(raw["mode"]),
			"restartLevel":     utils.ValueIgnoreEmpty(raw["level"]),
			"type":             "cdm",
			"instance":         "",
			"group":            "",
		}
		return params
	}
	return nil
}

func createClusterActionWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				createClusterActionWaitingHttpUrl = "v1.1/{project_id}/clusters/{cluster_id}"
				createClusterActionWaitingProduct = "cdm"
			)
			createClusterActionWaitingClient, err := cfg.NewServiceClient(createClusterActionWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CDM client: %s", err)
			}

			createClusterActionWaitingPath := createClusterActionWaitingClient.Endpoint + createClusterActionWaitingHttpUrl
			createClusterActionWaitingPath = strings.ReplaceAll(createClusterActionWaitingPath, "{project_id}",
				createClusterActionWaitingClient.ProjectID)
			createClusterActionWaitingPath = strings.ReplaceAll(createClusterActionWaitingPath,
				"{cluster_id}", d.Get("cluster_id").(string))

			createClusterActionWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json"},
			}

			createClusterActionWaitingResp, err := createClusterActionWaitingClient.Request("GET",
				createClusterActionWaitingPath, &createClusterActionWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createClusterActionWaitingRespBody, err := utils.FlattenResponse(createClusterActionWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`status`, createClusterActionWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"200",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createClusterActionWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"300", "303",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createClusterActionWaitingRespBody, status, nil
			}

			return createClusterActionWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceClusterActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterActionDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
