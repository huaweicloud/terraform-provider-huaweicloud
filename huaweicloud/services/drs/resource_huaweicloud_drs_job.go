package drs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/drs/v3/jobs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v3/{project_id}/jobs/batch-status
// @API DRS POST /v3/{project_id}/jobs
// @API DRS POST /v3/{project_id}/jobs/batch-connection
// @API DRS DELETE /v3/{project_id}/jobs/batch-jobs
// @API DRS PUT /v3/{project_id}/jobs/batch-limit-speed
// @API DRS POST /v3/{project_id}/jobs/batch-precheck-result
// @API DRS POST /v3/{project_id}/jobs/batch-precheck
// @API DRS POST /v3/{project_id}/jobs/batch-starting
// @API DRS POST /v3/{project_id}/jobs/batch-creation
// @API DRS POST /v3/{project_id}/jobs/batch-detail
// @API DRS PUT /v3/{project_id}/jobs/batch-modification
// @API DRS POST /v5/{project_id}/jobs/{resource_type}/{job_id}/tags/action
// @API DRS POST /v5/{project_id}/jobs/{job_id}/action
// @API DRS PUT /v5/{project_id}/jobs/{job_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
func ResourceDrsJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobCreate,
		ReadContext:   resourceJobRead,
		UpdateContext: resourceJobUpdate,
		DeleteContext: resourceJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^([A-Za-z][A-Za-z0-9-_\.]*)$`),
						"The name consists of 4 to 50 characters, starting with a letter. "+
							"Only letters, digits, underscores (_) and hyphens (-) are allowed."),
					validation.StringLenBetween(4, 50),
				),
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_db": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     dbInfoSchemaResource(),
			},

			"destination_db": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     dbInfoSchemaResource(),
			},

			"destination_db_readnoly": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"net_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "eip",
			},

			"migration_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "FULL_INCR_TRANS",
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^!<>&'"\\]*$`),
						"The 'description' has special character"),
					validation.StringLenBetween(1, 256),
				),
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"multi_write": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"expired_days": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  14,
			},

			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"migrate_definer": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"limit_speed": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 3,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"speed": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"start_time": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"end_time": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"tags": common.TagsSchema(),

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"is_sync_re_edit": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"action"},
			},

			"pause_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"action"},
			},

			"databases": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"tables"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},

			"tables": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"databases"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:     schema.TypeString,
							Required: true,
						},

						"table_names": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			// charge info: charging_mode, period_unit, period, auto_renew
			// once start the job, the bill will be auto paid
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func dbInfoSchemaResource() *schema.Resource {
	nodeResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"mysql", "mongodb", "gaussdbv5"}, false),
			},

			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"ssl_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"ssl_cert_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ssl_cert_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ssl_cert_check_sum": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ssl_cert_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}

	return &nodeResource
}

func resourceJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v3 client, error: %s", err)
	}
	clientV5, err := conf.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v5 client, error: %s", err)
	}

	opts, err := buildCreateParamter(d, client.ProjectID, conf.GetEnterpriseProjectID(d))
	if err != nil {
		return diag.FromErr(err)
	}

	rst, err := jobs.Create(client, *opts)
	if err != nil {
		return diag.Errorf("error creating DRS job: %s", err)
	}

	jobId := rst.Results[0].Id
	d.SetId(jobId)

	err = waitingforJobStatus(ctx, client, jobId, "create", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	valid := testConnections(client, jobId, opts.Jobs[0])
	if !valid {
		return diag.Errorf("test db connection of job: %s failed", jobId)
	}

	err = reUpdateJob(client, jobId, opts.Jobs[0], d.Get("migrate_definer").(bool))
	if err != nil {
		return diag.FromErr(err)
	}

	// Configure the transmission speed for the job.
	if v, ok := d.GetOk("limit_speed"); ok {
		configRaw := v.([]interface{})
		speedLimits := make([]jobs.SpeedLimitInfo, len(configRaw))
		for i, v := range configRaw {
			tmp := v.(map[string]interface{})
			speedLimits[i] = jobs.SpeedLimitInfo{
				Speed: tmp["speed"].(string),
				Begin: tmp["begin_time"].(string),
				End:   tmp["end_time"].(string),
			}
		}
		_, err = jobs.LimitSpeed(client, jobs.BatchLimitSpeedReq{
			SpeedLimits: []jobs.LimitSpeedReq{
				{
					JobId:      jobId,
					SpeedLimit: speedLimits,
				},
			},
		})

		if err != nil {
			return diag.Errorf("limit speed of job: %s failed, error: %s", jobId, err)
		}
	}

	// Only support migration or synchronization job to select objects, this limitation is stated in docs.
	_, ok1 := d.GetOk("databases")
	_, ok2 := d.GetOk("tables")
	if ok1 || ok2 {
		err = updateJobConfig(clientV5, buildUpdateJobConfigBodyParams(d, "db_object"), "db_object", d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = preCheck(ctx, client, jobId, d.Timeout(schema.TimeoutCreate), "forStartJob")
	if err != nil {
		return diag.FromErr(err)
	}

	startReq := jobs.StartJobReq{
		Jobs: []jobs.StartInfo{
			{
				JobId:     jobId,
				StartTime: d.Get("start_time").(string),
			},
		},
	}
	_, err = jobs.Start(client, startReq)

	if err != nil {
		return diag.Errorf("start DRS job failed,error: %s", err)
	}

	err = waitingforJobStatus(ctx, client, jobId, "start", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceJobRead(ctx, d, meta)
}

func buildUpdateJobConfigBodyParams(d *schema.ResourceData, updateType string) map[string]interface{} {
	// in next update for policy, it will change to switch/case
	if updateType == "db_object" {
		if _, ok1 := d.GetOk("databases"); ok1 {
			return map[string]interface{}{
				"db_object": map[string]interface{}{
					"object_scope": "database",
					"object_info":  buildDatabaseInfos(d.Get("databases").(*schema.Set).List()),
				},
			}
		}
		return map[string]interface{}{
			"db_object": map[string]interface{}{
				"object_scope": "table",
				"object_info":  buildTables(d.Get("tables").(*schema.Set).List()),
			},
		}
	}
	return nil
}

func buildDatabaseInfos(list []interface{}) map[string]interface{} {
	rst := make(map[string]interface{})
	for _, val := range list {
		if v, ok := val.(string); ok {
			m := map[string]interface{}{
				"name": v,
				"all":  true,
			}
			rst[v] = m
		}
	}
	return rst
}

func buildTables(tables []interface{}) map[string]interface{} {
	rst := make(map[string]interface{})
	for _, val := range tables {
		v := val.(map[string]interface{})
		database := v["database"].(string)
		tableNames := v["table_names"].(*schema.Set).List()
		rst[database] = map[string]interface{}{
			"name":   database,
			"tables": buildTableInfos(tableNames),
		}
	}
	return rst
}

func buildTableInfos(list []interface{}) map[string]interface{} {
	rst := make(map[string]interface{})
	for _, val := range list {
		if v, ok := val.(string); ok {
			m := map[string]interface{}{
				"name": v,
				"all":  true,
				"type": "table",
			}
			rst[v] = m
		}
	}
	return rst
}

func updateJobConfig(client *golangsdk.ServiceClient, jsonBody map[string]interface{}, updateType, id string) error {
	updateJobConfigHttpUrl := "v5/{project_id}/jobs/{job_id}"
	updateJobConfigPath := client.Endpoint + updateJobConfigHttpUrl
	updateJobConfigPath = strings.ReplaceAll(updateJobConfigPath, "{project_id}", client.ProjectID)
	updateJobConfigPath = strings.ReplaceAll(updateJobConfigPath, "{job_id}", id)
	updateJobConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"job": map[string]interface{}{
				"type":   updateType,
				"params": jsonBody,
			},
		}),
	}
	_, err := client.Request("PUT", updateJobConfigPath, &updateJobConfigOpt)
	if err != nil {
		return fmt.Errorf("error updating job for %s: %s", updateType, err)
	}
	return nil
}

func resourceJobRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v3 client, error: %s", err)
	}

	detailResp, err := jobs.Get(client, jobs.QueryJobReq{Jobs: []string{d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, parseDrsJobErrorToError404(err), "error retrieving DRS job")
	}
	detail := detailResp.Results[0]

	// Net_type is not in detail, so query by list.
	listResp, err := jobs.List(client, jobs.ListJobsReq{
		CurPage:   1,
		PerPage:   1,
		Name:      d.Id(),
		DbUseType: detail.DbUseType,
	})

	if err != nil {
		return diag.Errorf("query the job list by jobId: %s, error: %s", d.Id(), err)
	}

	createdAt, _ := strconv.ParseInt(detail.CreateTime, 10, 64)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.Name),
		d.Set("type", detail.DbUseType),
		d.Set("engine_type", detail.InstInfo.EngineType),
		d.Set("direction", detail.JobDirection),
		d.Set("net_type", listResp.Jobs[0].NetType),
		d.Set("public_ip", detail.InstInfo.PublicIp),
		d.Set("private_ip", detail.InstInfo.Ip),
		d.Set("destination_db_readnoly", detail.IsTargetReadonly),
		d.Set("migration_type", detail.TaskType),
		d.Set("description", detail.Description),
		d.Set("multi_write", detail.MultiWrite),
		d.Set("created_at", utils.FormatTimeStampRFC3339(createdAt/1000, false)),
		d.Set("status", detail.Status),
		d.Set("tags", utils.TagsToMap(detail.Tags)),
		setDbInfoToState(d, detail.SourceEndpoint, "source_db"),
		setDbInfoToState(d, detail.TargetEndpoint, "destination_db"),
	)

	// set objects
	if detail.ObjectSwitch {
		objectName := flattenObjectName(detail.ObjectInfos, detail.SyncDatabase)
		if detail.SyncDatabase {
			mErr = multierror.Append(mErr,
				d.Set("databases", objectName),
			)
		} else {
			mErr = multierror.Append(mErr,
				d.Set("tables", objectName),
			)
		}
	}

	// set charging info
	if !reflect.DeepEqual(detail.PeriodOrder, jobs.OrderInfo{}) {
		// `2` menas `month`, `3` means `year`
		periodType := "month"
		if detail.PeriodOrder.PeriodType == 3 {
			periodType = "year"
		}
		// after updating auto_renew from CBC, auto_renew doesn't change in DRS returns
		mErr = multierror.Append(mErr,
			d.Set("charging_mode", "prePaid"),
			d.Set("period_unit", periodType),
			d.Set("period", detail.PeriodOrder.PeriodNum),
			d.Set("order_id", detail.PeriodOrder.OrderId),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("charging_mode", "postPaid"),
		)
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting DRS job fields: %s", mErr)
	}

	return nil
}

func flattenObjectName(objectInfos []jobs.ObjectInfo, isDateBase bool) []interface{} {
	if len(objectInfos) == 0 {
		return nil
	}
	rst := make([]interface{}, 0)

	if isDateBase {
		for _, objectInfo := range objectInfos {
			rst = append(rst, objectInfo.Name)
		}
	} else {
		databases := make(map[string][]string)
		for _, objectInfo := range objectInfos {
			if objectInfo.Type == "table" {
				databases[objectInfo.ParentId] = append(databases[objectInfo.ParentId], objectInfo.Name)
			}
		}
		if len(databases) == 0 {
			return nil
		}
		for database, tableNames := range databases {
			rst = append(rst, map[string]interface{}{
				"database":    database,
				"table_names": tableNames,
			})
		}
	}

	return rst
}

func resourceJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v3 client, error: %s", err)
	}
	clientV5, err := conf.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v5 client, error: %s", err)
	}
	bssClient, err := conf.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS V2 client: %s", err)
	}

	// update name and description
	if d.HasChanges("name", "description") {
		detailResp, err := jobs.Get(client, jobs.QueryJobReq{Jobs: []string{d.Id()}})
		if err != nil {
			return common.CheckDeletedDiag(d, parseDrsJobErrorToError404(err), "error retrieving DRS job")
		}
		detail := detailResp.Results[0]

		if utils.StrSliceContains(
			[]string{"RELEASE_RESOURCE_COMPLETE", "RELEASE_RESOURCE_STARTED", "RELEASE_RESOURCE_FAILED"}, detail.Status) {
			return nil
		}

		updateParams := jobs.UpdateReq{
			Jobs: []jobs.UpdateJobReq{
				{
					JobId:       d.Id(),
					Name:        d.Get("name").(string),
					Description: d.Get("description").(string),
				},
			},
		}

		_, err = jobs.Update(client, updateParams)
		if err != nil {
			return diag.Errorf("update job: %s failed,error: %s", d.Id(), err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(clientV5, d, "jobs/"+d.Get("type").(string), d.Id())
		if tagErr != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("failed to update tags for DRS job(%s): %s", d.Id(), tagErr),
				},
			}
		}
	}

	if d.HasChange("action") {
		if action, ok := d.GetOk("action"); ok && utils.StrSliceContains([]string{"stop", "restart", "reset"}, action.(string)) {
			// precheck status
			resp, err := jobs.Status(client, jobs.QueryJobReq{Jobs: []string{d.Id()}})
			if err != nil {
				return diag.Errorf("error retrieving job status: %s", err)
			}
			if resp.Count == 0 || resp.Results[0].ErrorCode != "" {
				return diag.Errorf("error retrieving job status, %s: %s", resp.Results[0].ErrorCode, resp.Results[0].ErrorMessage)
			}
			err = preCheckStatus(resp.Results[0].Status)
			if err != nil {
				return diag.FromErr(err)
			}

			// preCheck job before reset
			if action.(string) == "reset" {
				err = preCheck(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate), "forRetryJob")
				if err != nil {
					return diag.FromErr(err)
				}
			}

			// execute action
			err = executeJobAction(clientV5, buildExecuteJobActionBodyParams(d), action.(string), d.Id())
			if err != nil {
				return diag.FromErr(err)
			}

			// wait for action complete
			err = waitingforJobStatus(ctx, client, d.Id(), action.(string), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChanges("databases", "tables") {
		err := updateObjectsSelection(ctx, d, client, clientV5)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		// resource_id is different from job_id
		resourceIDs, err := common.GetResourceIDsByOrder(bssClient, d.Get("order_id").(string), 1)
		if err != nil || len(resourceIDs) == 0 {
			return diag.Errorf("error getting resource IDs: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), resourceIDs[0]); err != nil {
			return diag.Errorf("error updating the auto-renew of DRS job (%s): %s", d.Id(), err)
		}
	}

	return resourceJobRead(ctx, d, meta)
}

func updateObjectsSelection(ctx context.Context, d *schema.ResourceData, client, clientV5 *golangsdk.ServiceClient) error {
	// preCheck type and status
	if d.Get("type").(string) != "sync" {
		return fmt.Errorf("only synchronization job supports updating object selection")
	}
	if !utils.StrSliceContains([]string{"INCRE_TRANSFER_STARTED", "INCRE_TRANSFER_FAILED"}, d.Get("status").(string)) {
		return fmt.Errorf("error updating synchronization object for status(%s)", d.Get("status").(string))
	}

	// update object
	_, ok1 := d.GetOk("databases")
	_, ok2 := d.GetOk("tables")
	if ok1 || ok2 {
		err := updateJobConfig(clientV5, buildUpdateJobConfigBodyParams(d, "db_object"), "db_object", d.Id())
		if err != nil {
			return err
		}
	}

	// preCheck job
	err := preCheck(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate), "forRetryJob")
	if err != nil {
		return err
	}

	// restart re-edit sync job
	reEditRequestBody := map[string]interface{}{
		"is_sync_re_edit": true,
	}
	err = executeJobAction(clientV5, reEditRequestBody, "restart", d.Id())
	if err != nil {
		return err
	}

	// wait for children transfer job started
	listResp, err := jobs.List(client, jobs.ListJobsReq{
		CurPage:   1,
		PerPage:   1,
		Name:      d.Id(),
		DbUseType: "sync",
	})
	if err != nil {
		return err
	}
	if listResp.Jobs[0].Children[0].Id == "" {
		return fmt.Errorf("error updating synchronization object: children synchronization job ID not found")
	}
	err = waitingforJobStatus(ctx, client, listResp.Jobs[0].Children[0].Id, "restart", d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func preCheckStatus(status string) error {
	switch status {
	case "stop":
		if !utils.StrSliceContains(
			[]string{"FULL_TRANSFER_STARTED", "FULL_TRANSFER_COMPLETE", "INCRE_TRANSFER_STARTED"}, status) {
			return fmt.Errorf("error pausing job for status(%s)", status)
		}
	case "restart":
		if status != "PAUSING" {
			return fmt.Errorf("error restarting job for status(%s)", status)
		}
	case "reset":
		if !utils.StrSliceContains([]string{"FULL_TRANSFER_FAILED", "INCRE_TRANSFER_FAILED"}, status) {
			return fmt.Errorf("error reseting job for status(%s)", status)
		}
	}
	return nil
}

func buildExecuteJobActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	action := d.Get("action").(string)
	switch action {
	case "stop":
		return map[string]interface{}{
			"pause_mode": utils.ValueIngoreEmpty(d.Get("pause_mode")),
		}
	case "restart":
		return map[string]interface{}{
			"is_sync_re_edit": utils.ValueIngoreEmpty(d.Get("is_sync_re_edit")),
		}
	case "reset":
		return map[string]interface{}{}
	}
	return nil
}

func executeJobAction(client *golangsdk.ServiceClient, jsonBody map[string]interface{}, action, id string) error {
	executeJobActionHttpUrl := "v5/{project_id}/jobs/{job_id}/action"
	executeJobActionPath := client.Endpoint + executeJobActionHttpUrl
	executeJobActionPath = strings.ReplaceAll(executeJobActionPath, "{project_id}", client.ProjectID)
	executeJobActionPath = strings.ReplaceAll(executeJobActionPath, "{job_id}", id)
	executeJobActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"job": map[string]interface{}{
				"action_name":   action,
				"action_params": jsonBody,
			},
		}),
	}
	executeJobActionResp, err := client.Request("POST", executeJobActionPath, &executeJobActionOpt)
	if err != nil {
		return fmt.Errorf("error executing job action: %s", err)
	}
	executeJobActionRespBody, err := utils.FlattenResponse(executeJobActionResp)
	if err != nil {
		return fmt.Errorf("error flattening job action response: %s", err)
	}
	status := utils.PathSearch("status", executeJobActionRespBody, nil)
	if status == nil {
		return fmt.Errorf("error getting job action status")
	} else if status.(string) != "success" {
		return fmt.Errorf("error executing job action: status(%s)", status)
	}

	return nil
}

func resourceJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v3 client, error: %s", err)
	}
	bssV2Client, err := conf.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	detailResp, err := jobs.Get(client, jobs.QueryJobReq{Jobs: []string{d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, parseDrsJobErrorToError404(err), "error retrieving DRS job")
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		// unsubscribe the order
		// resource_id is different from job_id
		resourceIDs, err := common.GetResourceIDsByOrder(bssV2Client, detailResp.Results[0].PeriodOrder.OrderId, 1)
		if err != nil {
			return diag.Errorf("error getting resource IDs: %s", err)
		}
		err = common.UnsubscribePrePaidResource(d, conf, resourceIDs)
		if err != nil {
			return diag.Errorf("error unsubscribing DRS job: %s", err)
		}
		err = waitingforJobStatus(ctx, client, d.Id(), "terminate", d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if !utils.StrSliceContains([]string{"CREATE_FAILED", "RELEASE_RESOURCE_COMPLETE", "RELEASE_CHILD_TRANSFER_COMPLETE"},
		detailResp.Results[0].Status) {
		// force terminate
		if !d.Get("force_destroy").(bool) {
			return diag.Errorf("the job: %s cannot be deleted when it is running. "+
				"If you want to forcibly delete the job please set force_destroy to True", d.Id())
		}

		dErr := jobs.Delete(client, jobs.BatchDeleteJobReq{
			Jobs: []jobs.DeleteJobReq{
				{
					DeleteType: jobs.DeleteTypeForceTerminate,
					JobId:      d.Id(),
				},
			},
		})

		if dErr.Err != nil {
			return diag.Errorf("terminate DRS job failed. %q: %s", d.Id(), dErr)
		}

		err = waitingforJobStatus(ctx, client, d.Id(), "terminate", d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	dErr := jobs.Delete(client, jobs.BatchDeleteJobReq{
		Jobs: []jobs.DeleteJobReq{
			{
				DeleteType: jobs.DeleteTypeDelete,
				JobId:      d.Id(),
			},
		},
	})
	if dErr.Err != nil {
		return diag.Errorf("delete DRS job failed. %q: %s", d.Id(), dErr)
	}

	return nil
}

func waitingforJobStatus(ctx context.Context, client *golangsdk.ServiceClient, id, statusType string,
	timeout time.Duration) error {
	var pending []string
	var target []string

	switch statusType {
	case "create":
		pending = []string{"CREATING"}
		target = []string{"CONFIGURATION"}
	case "start":
		pending = []string{"STARTJOBING", "WAITING_FOR_START"}
		target = []string{"FULL_TRANSFER_STARTED", "FULL_TRANSFER_COMPLETE", "INCRE_TRANSFER_STARTED"}
	case "terminate":
		pending = []string{"PENDING"}
		target = []string{"RELEASE_RESOURCE_COMPLETE"}
	case "stop":
		pending = []string{"FULL_TRANSFER_STARTED", "FULL_TRANSFER_COMPLETE", "INCRE_TRANSFER_STARTED"}
		target = []string{"PAUSING"}
	case "reset":
		pending = []string{"STARTJOBING", "WAITING_FOR_START"}
		target = []string{"FULL_TRANSFER_STARTED", "FULL_TRANSFER_COMPLETE", "INCRE_TRANSFER_STARTED"}
	case "restart":
		pending = []string{"STARTJOBING", "WAITING_FOR_START", "CHILD_TRANSFER_STARTING"}
		target = []string{"FULL_TRANSFER_STARTED", "FULL_TRANSFER_COMPLETE", "INCRE_TRANSFER_STARTED",
			"CHILD_TRANSFER_STARTED", "CHILD_TRANSFER_COMPLETE"}
	}

	stateConf := &resource.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			resp, err := jobs.Status(client, jobs.QueryJobReq{Jobs: []string{id}})
			if err != nil {
				return nil, "", err
			}
			if resp.Count == 0 || resp.Results[0].ErrorCode != "" {
				return resp, "failed", fmt.Errorf("%s: %s", resp.Results[0].ErrorCode, resp.Results[0].ErrorMessage)
			}

			if utils.StrSliceContains([]string{"START_JOB_FAILED", "CREATE_FAILED", "RELEASE_RESOURCE_FAILED",
				"CHILD_TRANSFER_FAILED"}, resp.Results[0].Status) {
				return resp, "failed", fmt.Errorf("%s", resp.Results[0].Status)
			}

			// if job is prepaid mode, need to unsubscribe from CBC, status delay changes into RELEASE_RESOURCE_STARTED
			if statusType == "terminate" && resp.Results[0].Status != "RELEASE_RESOURCE_COMPLETE" {
				return resp, "PENDING", nil
			}

			return resp, resp.Results[0].Status, nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DRS job: %s to be %s: %s", id, statusType, err)
	}
	return nil
}

func buildCreateParamter(d *schema.ResourceData, projectId, enterpriseProjectID string) (*jobs.BatchCreateJobReq, error) {
	jobDirection := d.Get("direction").(string)

	sourceDb, err := buildDbConfigParamter(d, "source_db", projectId)
	if err != nil {
		return nil, err
	}

	targetDb, err := buildDbConfigParamter(d, "destination_db", projectId)
	if err != nil {
		return nil, err
	}

	var subnetId string
	if jobDirection == "up" {
		if targetDb.InstanceId == "" {
			return nil, fmt.Errorf("destination_db.0.instance_id is required When diretion is down")
		}
		subnetId = targetDb.SubnetId
	} else {
		if sourceDb.InstanceId == "" {
			return nil, fmt.Errorf("source_db.0.instance_id is required When diretion is down")
		}
		subnetId = sourceDb.SubnetId
	}

	var bindEip bool
	if d.Get("net_type").(string) == "eip" {
		bindEip = true
	}

	job := jobs.CreateJobReq{
		Name:             d.Get("name").(string),
		DbUseType:        d.Get("type").(string),
		EngineType:       d.Get("engine_type").(string),
		JobDirection:     jobDirection,
		NetType:          d.Get("net_type").(string),
		BindEip:          utils.Bool(bindEip),
		IsTargetReadonly: utils.Bool(d.Get("destination_db_readnoly").(bool)),
		TaskType:         d.Get("migration_type").(string),
		Description:      d.Get("description").(string),
		MultiWrite:       utils.Bool(d.Get("multi_write").(bool)),
		ExpiredDays:      fmt.Sprint(d.Get("expired_days").(int)),
		NodeType:         "high",
		SourceEndpoint:   *sourceDb,
		TargetEndpoint:   *targetDb,
		SubnetId:         subnetId,
		Tags:             utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		SysTags:          utils.BuildSysTags(enterpriseProjectID),
	}

	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode.(string) == "prePaid" {
		// the prepaid mode only take effect when `type` is `sync` or `cloudDataGuard`,
		// this limitation has been stated in docs
		if err = common.ValidatePrePaidChargeInfo(d); err != nil {
			return nil, err
		}

		// `2` menas `month`, `3` means `year`
		periodType := 2
		if d.Get("period_unit").(string) == "year" {
			periodType = 3
		}

		autoRenew := 0
		if d.Get("auto_renew").(string) == "true" {
			autoRenew = 1
		}
		job.ChargingMode = "period"
		job.PeriodOrder = &jobs.PeriodOrder{
			PeriodType:  periodType,
			PeriodNum:   d.Get("period").(int),
			IsAutoRenew: autoRenew,
		}
	}

	return &jobs.BatchCreateJobReq{Jobs: []jobs.CreateJobReq{job}}, nil
}

func buildDbConfigParamter(d *schema.ResourceData, dbType, projectId string) (*jobs.Endpoint, error) {
	configRaw := d.Get(dbType).([]interface{})[0].(map[string]interface{})
	configs := jobs.Endpoint{
		DbType:          configRaw["engine_type"].(string),
		Ip:              configRaw["ip"].(string),
		DbName:          configRaw["name"].(string),
		DbUser:          configRaw["user"].(string),
		DbPassword:      configRaw["password"].(string),
		DbPort:          golangsdk.IntToPointer(configRaw["port"].(int)),
		InstanceId:      configRaw["instance_id"].(string),
		Region:          configRaw["region"].(string),
		VpcId:           configRaw["vpc_id"].(string),
		SubnetId:        configRaw["subnet_id"].(string),
		ProjectId:       projectId,
		SslCertPassword: configRaw["ssl_cert_password"].(string),
		SslCertCheckSum: configRaw["ssl_cert_check_sum"].(string),
		SslCertKey:      configRaw["ssl_cert_key"].(string),
		SslCertName:     configRaw["ssl_cert_name"].(string),
		SslLink:         utils.Bool(configRaw["ssl_enabled"].(bool)),
	}
	return &configs, nil
}

func parseDrsJobErrorToError404(respErr error) error {
	var apiError jobs.JobDetailResp

	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiError)
		if pErr == nil &&
			(apiError.Results[0].ErrorCode == "DRS.M00289" || apiError.Results[0].ErrorCode == "DRS.M05004") {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}

func setDbInfoToState(d *schema.ResourceData, endpoint jobs.Endpoint, fieldName string) error {
	result := make([]interface{}, 1)
	item := map[string]interface{}{
		"engine_type":        endpoint.DbType,
		"ip":                 endpoint.Ip,
		"port":               endpoint.DbPort,
		"password":           endpoint.DbPassword,
		"user":               endpoint.DbUser,
		"instance_id":        endpoint.InstanceId,
		"name":               endpoint.InstanceName,
		"region":             endpoint.Region,
		"vpc_id":             endpoint.VpcId,
		"subnet_id":          endpoint.SubnetId,
		"ssl_cert_password":  endpoint.SslCertPassword,
		"ssl_cert_check_sum": endpoint.SslCertCheckSum,
		"ssl_cert_key":       endpoint.SslCertKey,
		"ssl_cert_name":      endpoint.SslCertName,
		"ssl_enabled":        endpoint.SslLink,
	}
	result[0] = item
	// lintignore:R001
	return d.Set(fieldName, result)
}

func testConnections(client *golangsdk.ServiceClient, jobId string, opts jobs.CreateJobReq) (valid bool) {
	reqParams := jobs.TestConnectionsReq{
		Jobs: []jobs.TestEndPoint{
			{
				JobId:        jobId,
				NetType:      opts.NetType,
				EndPointType: "so",
				ProjectId:    client.ProjectID,
				Region:       opts.SourceEndpoint.Region,
				VpcId:        opts.SourceEndpoint.VpcId,
				SubnetId:     opts.SourceEndpoint.SubnetId,
				DbType:       opts.SourceEndpoint.DbType,
				Ip:           opts.SourceEndpoint.Ip,
				DbUser:       opts.SourceEndpoint.DbUser,
				DbPassword:   opts.SourceEndpoint.DbPassword,
				DbPort:       opts.SourceEndpoint.DbPort,
				SslLink:      opts.SourceEndpoint.SslLink,
				InstId:       opts.SourceEndpoint.InstanceId,
			},
			{
				JobId:        jobId,
				NetType:      opts.NetType,
				EndPointType: "ta",
				ProjectId:    client.ProjectID,
				Region:       opts.TargetEndpoint.Region,
				VpcId:        opts.TargetEndpoint.VpcId,
				SubnetId:     opts.TargetEndpoint.SubnetId,
				DbType:       opts.TargetEndpoint.DbType,
				Ip:           opts.TargetEndpoint.Ip,
				DbUser:       opts.TargetEndpoint.DbUser,
				DbPassword:   opts.TargetEndpoint.DbPassword,
				DbPort:       opts.TargetEndpoint.DbPort,
				SslLink:      opts.TargetEndpoint.SslLink,
				InstId:       opts.TargetEndpoint.InstanceId,
			},
		},
	}
	rsp, err := jobs.TestConnections(client, reqParams)
	if err != nil || rsp.Count != 2 {
		log.Printf("[ERROR] test connections of job: %s failed,error: %s", jobId, err)
		return false
	}

	valid = rsp.Results[0].Success && rsp.Results[1].Success
	return
}

func reUpdateJob(client *golangsdk.ServiceClient, jobId string, opts jobs.CreateJobReq, migrateDefiner bool) error {
	reqParams := jobs.UpdateReq{
		Jobs: []jobs.UpdateJobReq{
			{
				JobId:            jobId,
				Name:             opts.Name,
				NetType:          opts.NetType,
				EngineType:       opts.EngineType,
				NodeType:         opts.NodeType,
				StoreDbInfo:      true,
				IsRecreate:       utils.Bool(false),
				DbUseType:        opts.DbUseType,
				Description:      opts.Description,
				TaskType:         opts.TaskType,
				JobDirection:     opts.JobDirection,
				IsTargetReadonly: opts.IsTargetReadonly,
				ReplaceDefiner:   &migrateDefiner,
				SourceEndpoint:   &opts.SourceEndpoint,
				TargetEndpoint:   &opts.TargetEndpoint,
			},
		},
	}

	_, err := jobs.Update(client, reqParams)
	if err != nil {
		return fmt.Errorf("update job failed,error: %s", err)
	}

	return nil
}

func preCheck(ctx context.Context, client *golangsdk.ServiceClient, jobId string, timeout time.Duration, precheckMode string) error {
	_, err := jobs.PreCheckJobs(client, jobs.BatchPrecheckReq{
		Jobs: []jobs.PreCheckInfo{
			{
				JobId:        jobId,
				PrecheckMode: precheckMode,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("start job: %s preCheck failed,error: %s", jobId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			resp, err := jobs.CheckResults(client, jobs.QueryPrecheckResultReq{
				Jobs: []string{jobId},
			})
			if err != nil {
				return nil, "", err
			}
			if resp.Count == 0 || resp.Results[0].ErrorCode != "" {
				return resp, "failed", fmt.Errorf("%s: %s", resp.Results[0].ErrorCode, resp.Results[0].ErrorMsg)
			}

			if resp.Results[0].Process != "100%" {
				return resp, "pending", nil
			}

			if resp.Results[0].TotalPassedRate == "100%" {
				return resp, "complete", nil
			}

			return resp, "failed", fmt.Errorf("some preCheck item failed: %v", resp)
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DRS job: %s to be terminate: %s", jobId, err)
	}
	return nil
}
