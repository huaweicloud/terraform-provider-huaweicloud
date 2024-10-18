package dws

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

// @API DWS POST /v2/{project_id}/disaster-recoveries
// @API DWS GET /v2/{project_id}/disaster-recovery/{disaster_recovery_id}
// @API DWS DELETE /v2/{project_id}/disaster-recovery/{disaster_recovery_id}
// @API DWS PUT /v2/{project_id}/disaster-recovery/{disaster_recovery_id}
// @API DWS POST /v2/{project_id}/disaster-recovery/{disaster_recovery_id}/start
// @API DWS POST /v2/{project_id}/disaster-recovery/{disaster_recovery_id}/pause
// @API DWS POST /v2/{project_id}/disaster-recovery/{disaster_recovery_id}/switchover
func ResourceDwsDisasterRecoveryTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDwsDisasterRecoveryTaskCreate,
		UpdateContext: resourceDwsDisasterRecoveryTaskUpdate,
		ReadContext:   resourceDwsDisasterRecoveryTaskRead,
		DeleteContext: resourceDwsDisasterRecoveryTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dr_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"primary_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"standby_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dr_sync_period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"started_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_cluster": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     clusterSchema(),
			},
			"standby_cluster": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     clusterSchema(),
			},
		},
	}
}

func clusterSchema() *schema.Resource {
	nodeResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_az": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_success_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"obs_bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &nodeResource
}

func resourceDwsDisasterRecoveryTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/disaster-recoveries"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDisasterRecoveryTaskBodyParams(d)),
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DWS disaster recovery: %s", err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	recoveryId := utils.PathSearch("disaster_recovery.id", respBody, "").(string)
	if recoveryId == "" {
		return diag.Errorf("unable to find the DWS disaster recovery ID from the API response")
	}
	d.SetId(recoveryId)
	// When the disaster recovery successfully created, the status is unstart.
	err = waitingForActionDisasterRecvery(ctx, d, client, "unstart", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the creation of DWS disaster recovery (%s) to complete: %s", d.Id(), err)
	}

	action := d.Get("action").(string)
	if action == "start" {
		err = doActionDisasterRecoveryTask(client, d, action)
		if err != nil {
			return diag.Errorf("error starting DWS disaster recovery: %s", err)
		}
		// When the disaster recovery successfully started, the status is running.
		err = waitingForActionDisasterRecvery(ctx, d, client, "running", d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for the start of DWS disaster recovery (%s) to complete: %s", d.Id(), err)
		}
	}
	return resourceDwsDisasterRecoveryTaskRead(ctx, d, meta)
}

func DisasterRecoveryTaskStatusRefresh(client *golangsdk.ServiceClient, d *schema.ResourceData, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getDisasterRecoveryTask(client, d.Id())
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				obj := map[string]string{"code": "COMPLETED"}
				return obj, "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		status := utils.PathSearch("disaster_recovery.status", respBody, "").(string)
		unexpectedStatuses := []string{
			"create_failed", "start_failed", "stop_failed", "abnormal",
		}
		if utils.StrSliceContains(unexpectedStatuses, status) {
			return respBody, "ERROR", fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func waitingForActionDisasterRecvery(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	target string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      DisasterRecoveryTaskStatusRefresh(client, d, []string{target}),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func buildCreateDisasterRecoveryTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"disaster_recovery": map[string]interface{}{
			"name":               d.Get("name"),
			"dr_type":            d.Get("dr_type"),
			"primary_cluster_id": d.Get("primary_cluster_id"),
			"standby_cluster_id": d.Get("standby_cluster_id"),
			"dr_sync_period":     d.Get("dr_sync_period"),
		},
	}
}

func buildUpdateDisasterRecoveryTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"disaster_recovery": map[string]interface{}{
			"dr_sync_period": d.Get("dr_sync_period"),
		},
	}
}

func resourceDwsDisasterRecoveryTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	respBody, err := getDisasterRecoveryTask(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DWS disaster recovery")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("disaster_recovery.name", respBody, nil)),
		d.Set("dr_sync_period", utils.PathSearch("disaster_recovery.dr_sync_period", respBody, nil)),
		d.Set("dr_type", utils.PathSearch("disaster_recovery.dr_type", respBody, nil)),
		d.Set("status", utils.PathSearch("disaster_recovery.status", respBody, "")),
		d.Set("started_at", utils.PathSearch("disaster_recovery.start_time", respBody, nil)),
		d.Set("created_at", utils.PathSearch("disaster_recovery.create_time", respBody, nil)),
		d.Set("primary_cluster", flattenClusterInfo(utils.PathSearch("disaster_recovery.primary_cluster", respBody, nil))),
		d.Set("standby_cluster", flattenClusterInfo(utils.PathSearch("disaster_recovery.standby_cluster", respBody, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterInfo(cluster interface{}) []map[string]interface{} {
	if cluster == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":              utils.PathSearch("id", cluster, nil),
			"region":          utils.PathSearch("region", cluster, nil),
			"name":            utils.PathSearch("name", cluster, nil),
			"cluster_az":      utils.PathSearch("cluster_az", cluster, nil),
			"role":            utils.PathSearch("role", cluster, nil),
			"status":          utils.PathSearch("status", cluster, nil),
			"progress":        utils.PathSearch("progress", cluster, nil),
			"last_success_at": utils.PathSearch("last_success_time", cluster, nil),
			"obs_bucket_name": utils.PathSearch("obs_bucket_name", cluster, nil),
		},
	}
}

func getDisasterRecoveryTask(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	httpUrl := "v2/{project_id}/disaster-recovery/{disaster_recovery_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{disaster_recovery_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "DWS.10101")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func doActionDisasterRecoveryTask(client *golangsdk.ServiceClient, d *schema.ResourceData, action string) error {
	httpUrl := "v2/{project_id}/disaster-recovery/{disaster_recovery_id}/{action}"
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{action}", action)
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{disaster_recovery_id}", d.Id())
	actionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	resp, err := client.Request("POST", actionPath, &actionOpt)
	if err != nil {
		return common.ConvertExpected400ErrInto404Err(err, "error_code", "DWS.10101")
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return err
	}
	return nil
}

func updateDisasterRecoveryTask(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/disaster-recovery/{disaster_recovery_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{disaster_recovery_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateDisasterRecoveryTaskBodyParams(d)),
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	resp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "DWS.10101")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func resourceDwsDisasterRecoveryTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if d.HasChange("dr_sync_period") {
		_, err = updateDisasterRecoveryTask(client, d)
		if err != nil {
			return diag.Errorf("error updating DWS disaster recovery: %s", err)
		}
	}

	if d.HasChange("action") {
		action := d.Get("action").(string)
		err := doAction(ctx, action, client, d)
		if err != nil {
			return err
		}
	}

	return resourceDwsDisasterRecoveryTaskRead(ctx, d, meta)
}

func doAction(ctx context.Context, action string, client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	switch action {
	case "start":
		err := doActionDisasterRecoveryTask(client, d, action)
		if err != nil {
			return diag.Errorf("error starting DWS disaster recovery: %s", err)
		}
		// When the disaster recovery successfully started, the status is running.
		err = waitingForActionDisasterRecvery(ctx, d, client, "running", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the start of DWS disaster recovery (%s) to complete: %s", d.Id(), err)
		}
	case "pause":
		err := doActionDisasterRecoveryTask(client, d, action)
		if err != nil {
			return diag.Errorf("error stopping DWS disaster recovery: %s", err)
		}
		// When the disaster recovery successfully paused, the status is stopped.
		err = waitingForActionDisasterRecvery(ctx, d, client, "stopped", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the stop of DWS disaster recovery (%s) to complete: %s", d.Id(), err)
		}
	case "switchover":
		err := doActionDisasterRecoveryTask(client, d, action)
		if err != nil {
			return diag.Errorf("error switchovering DWS disaster recovery: %s", err)
		}
		// When the disaster recovery successfully swtichovered, the status is running.
		err = waitingForActionDisasterRecvery(ctx, d, client, "running", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the swtichover of DWS disaster recovery (%s) to complete: %s", d.Id(), err)
		}
	default:
		return diag.Errorf("not supported action: " + action)
	}

	return nil
}

func resourceDwsDisasterRecoveryTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/disaster-recovery/{disaster_recovery_id}"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	// When the disaster recovery is running, should stop it frist.
	respBody, err := getDisasterRecoveryTask(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting disaster recovery when check the status")
	}
	status := utils.PathSearch("disaster_recovery.status", respBody, "").(string)
	if status == "running" {
		err = doActionDisasterRecoveryTask(client, d, "pause")
		if err != nil {
			return diag.Errorf("error stopping DWS client: %s", err)
		}
		err = waitingForDisasterDeleteStopStatus(ctx, d, client)
		if err != nil {
			return diag.Errorf("error waiting for the stop of DWS disaster recovery (%s) to complete: %s", d.Id(), err)
		}
	}
	delPath := client.Endpoint + httpUrl
	delPath = strings.ReplaceAll(delPath, "{project_id}", client.ProjectID)
	delPath = strings.ReplaceAll(delPath, "{disaster_recovery_id}", d.Id())
	delStageOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	_, err = client.Request("DELETE", delPath, &delStageOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DWS.10101"),
			"error deleting disaster recovery")
	}

	err = waitingForDisasterDeleteStatus(ctx, d, client)
	if err != nil {
		return diag.Errorf("error waiting for the deleting of DWS disaster recovery (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func waitingForDisasterDeleteStatus(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      DisasterRecoveryTaskStatusRefresh(client, d, []string{}),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitingForDisasterDeleteStopStatus(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      DisasterRecoveryTaskStopStatusRefresh(client, d, []string{"stopped", "stop_failed"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func DisasterRecoveryTaskStopStatusRefresh(client *golangsdk.ServiceClient, d *schema.ResourceData, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getDisasterRecoveryTask(client, d.Id())
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				obj := map[string]string{"code": "COMPLETED"}
				return obj, "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		status := utils.PathSearch("disaster_recovery.status", respBody, "").(string)

		if utils.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}
