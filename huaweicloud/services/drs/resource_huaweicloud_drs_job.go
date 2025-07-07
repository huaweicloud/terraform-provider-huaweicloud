package drs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
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

var notFoundErrCode = []string{
	"DRS.M00289", // non exist
	"DRS.M05004", // deleted
}

// @API DRS POST /v3/{project_id}/jobs/batch-status
// @API DRS POST /v3/{project_id}/jobs
// @API DRS POST /v3/{project_id}/jobs/batch-connection
// @API DRS POST /v3/{project_id}/jobs/cluster/batch-connection
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
// @API DRS POST /v3/{project_id}/jobs/batch-sync-policy
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
			"node_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "high",
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
			},
			"migrate_definer": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"limit_speed": {
				Type:     schema.TypeSet,
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
			"policy_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_ddl_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"conflict_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"index_trans": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},

						// Kafka
						"topic_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"partition_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"kafka_data_format": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"topic_name_format": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"partitions_num": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"replication_factor": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						// PostgreSQL
						"is_fill_materialized_view": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"export_snapshot": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},

						// GaussDB Primary/Standby to Kafka primary and standby task
						"slot_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						// incre
						"file_and_position": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"gtid_set": {
							Type:     schema.TypeString,
							Optional: true,
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
			"public_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"master_az": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"slave_az"},
			},
			"slave_az": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"master_az"},
			},

			// charge info: charging_mode, period_unit, period, auto_renew
			// once start the job, the bill will be auto paid
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
				Description: "schema: Internal",
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period"},
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
				Description: "schema: Internal",
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_unit"},
				Description:  "schema: Internal",
			},
			"auto_renew": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
				Description: "schema: Internal",
			},

			"alarm_notify": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_urn": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"delay_time": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"rpo_delay": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"rto_delay": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"is_open_fast_clean": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slave_job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
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
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"original_job_direction": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dbInfoSchemaResource() *schema.Resource {
	nodeResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
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
			"kafka_security_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     dbInfoKafkaSecurityConfigSchemaResource(),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &nodeResource
}

func dbInfoKafkaSecurityConfigSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sasl_mechanism": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"trust_store_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"trust_store_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"trust_store_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"endpoint_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delegation_tokens": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_key_store": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"key_store_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_store_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_store_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"set_private_key_password": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"key_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
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

	if _, ok := d.GetOk("master_az"); ok {
		valid := testConnectionsForDualAZ(client, jobId, opts.Jobs[0])
		if !valid {
			return diag.Errorf("test db connection of job: %s failed", jobId)
		}
	} else {
		valid := testConnections(client, jobId, opts.Jobs[0])
		if !valid {
			return diag.Errorf("test db connection of job: %s failed", jobId)
		}
	}

	err = reUpdateJob(client, jobId, opts.Jobs[0], d.Get("migrate_definer").(bool))
	if err != nil {
		return diag.FromErr(err)
	}

	// Configure the transmission speed for the job.
	if v, ok := d.GetOk("limit_speed"); ok {
		configRaw := v.(*schema.Set).List()
		speedLimits := make([]jobs.SpeedLimitInfo, len(configRaw))
		for i, v := range configRaw {
			tmp := v.(map[string]interface{})
			speedLimits[i] = jobs.SpeedLimitInfo{
				Speed: tmp["speed"].(string),
				Begin: tmp["start_time"].(string),
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

	if _, ok := d.GetOk("policy_config"); ok {
		err = updateJobPolicyConfig(client, d)
		if err != nil {
			return diag.FromErr(err)
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

	if _, ok := d.GetOk("alarm_notify"); ok {
		err = updateJobConfig(clientV5, buildUpdateJobConfigBodyParams(d, "notify"), "notify", d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	startTime := d.Get("start_time").(string)
	startMode := "start"
	if startTime != "" && startTime != "0" {
		startMode = "start_later"
	}

	startReq := jobs.StartJobReq{
		Jobs: []jobs.StartInfo{
			{
				JobId:     jobId,
				StartTime: startTime,
			},
		},
	}
	_, err = jobs.Start(client, startReq)
	if err != nil {
		return diag.Errorf("start DRS job failed,error: %s", err)
	}

	err = waitingforJobStatus(ctx, client, jobId, startMode, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceJobRead(ctx, d, meta)
}

func buildUpdateJobConfigBodyParams(d *schema.ResourceData, updateType string) map[string]interface{} {
	switch updateType {
	case "db_object":
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
	case "notify":
		return map[string]interface{}{
			"alarm_notify": buildAlarmNotify(d.Get("alarm_notify").([]interface{})),
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

func buildAlarmNotify(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}
	raw := rawArray[0].(map[string]interface{})
	rst := map[string]interface{}{
		"alarm_to_user": true,
		"topic_urn":     raw["topic_urn"],
		"delay_time":    utils.ValueIgnoreEmpty(raw["delay_time"]),
		"rpo_delay":     utils.ValueIgnoreEmpty(raw["rpo_delay"]),
		"rto_delay":     utils.ValueIgnoreEmpty(raw["rto_delay"]),
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

func buildJobPolicyConfigRequestBody(rawArray []interface{}, id string) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}
	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}
	rst := map[string]interface{}{
		"jobs": []map[string]interface{}{
			{
				"job_id":                    id,
				"filter_ddl_policy":         utils.ValueIgnoreEmpty(raw["filter_ddl_policy"]),
				"conflict_policy":           utils.ValueIgnoreEmpty(raw["conflict_policy"]),
				"index_trans":               utils.ValueIgnoreEmpty(raw["index_trans"]),
				"topic_policy":              utils.ValueIgnoreEmpty(raw["topic_policy"]),
				"topic":                     utils.ValueIgnoreEmpty(raw["topic"]),
				"partition_policy":          utils.ValueIgnoreEmpty(raw["partition_policy"]),
				"kafka_data_format":         utils.ValueIgnoreEmpty(raw["kafka_data_format"]),
				"topic_name_format":         utils.ValueIgnoreEmpty(raw["topic_name_format"]),
				"partitions_num":            utils.ValueIgnoreEmpty(raw["partitions_num"]),
				"replication_factor":        utils.ValueIgnoreEmpty(raw["replication_factor"]),
				"is_fill_materialized_view": utils.ValueIgnoreEmpty(raw["is_fill_materialized_view"]),
				"export_snapshot":           utils.ValueIgnoreEmpty(raw["export_snapshot"]),
				"slot_name":                 utils.ValueIgnoreEmpty(raw["slot_name"]),
				"file_and_position":         utils.ValueIgnoreEmpty(raw["file_and_position"]),
				"gtid_set":                  utils.ValueIgnoreEmpty(raw["gtid_set"]),
			},
		},
	}
	return rst
}

func updateJobPolicyConfig(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updatePolicyConfigHttpUrl := "v3/{project_id}/jobs/batch-sync-policy"
	updatePolicyConfigPath := client.Endpoint + updatePolicyConfigHttpUrl
	updatePolicyConfigPath = strings.ReplaceAll(updatePolicyConfigPath, "{project_id}", client.ProjectID)
	updatePolicyConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildJobPolicyConfigRequestBody(d.Get("policy_config").([]interface{}), d.Id())),
	}
	_, err := client.Request("POST", updatePolicyConfigPath, &updatePolicyConfigOpt)
	if err != nil {
		return fmt.Errorf("error updating policy config: %s", err)
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
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "results[0].error_code", notFoundErrCode...),
			"error retrieving DRS job")
	}
	if len(detailResp.Results) == 0 || detailResp.Results[0].Status == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DRS job")
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

	progressResp, err := jobs.Progress(client, jobs.QueryJobReq{Jobs: []string{d.Id()}})
	if err != nil {
		return diag.Errorf("error getting job progress: %s", err)
	}

	// get topicUrn, it's not in api response
	topicUrn := ""
	if v, ok := d.GetOk("alarm_notify"); ok {
		alarmNotify := v.([]interface{})[0].(map[string]interface{})
		topicUrn = alarmNotify["topic_urn"].(string)
	}

	createdAt, _ := strconv.ParseInt(detail.CreateTime, 10, 64)
	updatedAt, _ := strconv.ParseInt(detail.UpdateTime, 10, 64)
	// engine_type input mongodb will return mongodb-to-dds
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.Name),
		d.Set("type", detail.DbUseType),
		d.Set("direction", detail.JobDirection),
		d.Set("net_type", listResp.Jobs[0].NetType),
		d.Set("public_ip", detail.InstInfo.PublicIp),
		d.Set("private_ip", detail.InstInfo.Ip),
		d.Set("node_type", detail.InstInfo.InstType),
		d.Set("destination_db_readnoly", detail.IsTargetReadonly),
		d.Set("migration_type", detail.TaskType),
		d.Set("description", detail.Description),
		d.Set("multi_write", detail.MultiWrite),
		d.Set("created_at", utils.FormatTimeStampRFC3339(createdAt/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(updatedAt/1000, false)),
		d.Set("status", detail.Status),
		d.Set("progress", progressResp.Results[0].Progress),
		d.Set("tags", utils.TagsToMap(detail.Tags)),
		d.Set("alarm_notify", flattenAlarmNotify(detail.AlarmNotify, topicUrn)),
		d.Set("limit_speed", flattenLimitSpeed(detail.SpeedLimit)),
		d.Set("master_az", detail.MasterAz),
		d.Set("master_job_id", detail.MasterJobId),
		d.Set("slave_az", detail.SlaveAz),
		d.Set("slave_job_id", getSlaveJobID(listResp.Jobs[0].Children, detail.MasterJobId)),
		d.Set("vpc_id", detail.VpcId),
		d.Set("subnet_id", detail.SubnetId),
		d.Set("security_group_id", detail.SecurityGroupId),
		setDbInfoToState(d, detail.SourceEndpoint, "source_db"),
		setDbInfoToState(d, detail.TargetEndpoint, "destination_db"),
		d.Set("is_open_fast_clean", detail.IsOpenFastClean),
		d.Set("original_job_direction", detail.OriginalJobDirection),
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

func flattenAlarmNotify(alarmNotify jobs.AlarmNotifyInfo, topicUrn string) []interface{} {
	rst := make([]interface{}, 0)
	v := map[string]interface{}{
		"topic_urn":  topicUrn,
		"delay_time": alarmNotify.DelayTime,
		"rpo_delay":  alarmNotify.RpoDelay,
		"rto_delay":  alarmNotify.RtoDelay,
	}
	rst = append(rst, v)
	return rst
}

func flattenLimitSpeed(speedLimit []jobs.SpeedLimitInfo) []interface{} {
	rst := make([]interface{}, 0)
	for _, limit := range speedLimit {
		v := map[string]interface{}{
			"speed":      limit.Speed,
			"start_time": limit.Begin,
			"end_time":   limit.End,
		}
		rst = append(rst, v)
	}
	return rst
}

func getSlaveJobID(children []jobs.ChildrenJobInfo, masterJobId string) string {
	if len(children) != 2 {
		return ""
	}
	if children[0].Id == masterJobId {
		return children[1].Id
	}
	return children[0].Id
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
			return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "results[0].error_code", notFoundErrCode...),
				"error retrieving DRS job")
		}
		if len(detailResp.Results) == 0 || detailResp.Results[0].Status == "DELETED" {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DRS job")
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

	if d.HasChange("start_time") {
		if v := d.Get("start_time").(string); v != "0" && v != "" {
			err = executeJobAction(clientV5, buildExecuteJobActionBodyParams(d, "start_later"), "start", d.Id())
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("action") {
		if action, ok := d.GetOk("action"); ok && utils.StrSliceContains([]string{"stop", "restart", "reset", "start"}, action.(string)) {
			// precheck status
			resp, err := jobs.Status(client, jobs.QueryJobReq{Jobs: []string{d.Id()}})
			if err != nil {
				return diag.Errorf("error retrieving job status: %s", err)
			}
			if resp.Count == 0 || resp.Results[0].ErrorCode != "" {
				return diag.Errorf("error retrieving job status, %s: %s", resp.Results[0].ErrorCode, resp.Results[0].ErrorMessage)
			}
			err = preCheckStatus(action.(string), resp.Results[0].Status)
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
			err = executeJobAction(clientV5, buildExecuteJobActionBodyParams(d, action.(string)), action.(string), d.Id())
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

	// wait 10 seconds before getting the job, to avoid delay for getting children info
	// lintignore:R018
	time.Sleep(10 * time.Second)

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
	if len(listResp.Jobs) == 0 || len(listResp.Jobs[0].Children) == 0 {
		return fmt.Errorf("error getting children job info")
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

func preCheckStatus(action, status string) error {
	switch action {
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
	case "start":
		if status != "WAITING_FOR_START" {
			return fmt.Errorf("error starting job for status(%s)", status)
		}
	}

	return nil
}

func buildExecuteJobActionBodyParams(d *schema.ResourceData, action string) map[string]interface{} {
	switch action {
	case "stop":
		return map[string]interface{}{
			"pause_mode": utils.ValueIgnoreEmpty(d.Get("pause_mode")),
		}
	case "restart":
		return map[string]interface{}{
			"is_sync_re_edit": utils.ValueIgnoreEmpty(d.Get("is_sync_re_edit")),
		}
	case "reset":
		return map[string]interface{}{}
	case "start":
		return map[string]interface{}{}
	case "start_later":
		return map[string]interface{}{
			"start_time": utils.ValueIgnoreEmpty(d.Get("start_time")),
		}
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

func waitForOrderDetail(ctx context.Context, bssV2Client *golangsdk.ServiceClient, timeout time.Duration, orderId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			resourceIDs, err := common.GetResourceIDsByOrder(bssV2Client, orderId, 1)
			if err != nil {
				if strings.Contains(err.Error(), "response empty") {
					return resourceIDs, "pending", nil
				}
				return nil, "error", err
			}

			return resourceIDs, "complete", nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error getting resource_id: %s", err)
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
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "results[0].error_code", notFoundErrCode...),
			"error retrieving DRS job")
	}
	if len(detailResp.Results) == 0 || detailResp.Results[0].Status == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DRS job")
	}
	orderId := detailResp.Results[0].PeriodOrder.OrderId

	if d.Get("charging_mode").(string) == "prePaid" && strings.TrimSpace(orderId) != "" {
		// unsubscribe the order
		// resource_id is different from job_id
		// searching order has delay
		err := waitForOrderDetail(ctx, bssV2Client, d.Timeout(schema.TimeoutDelete), orderId)
		if err != nil {
			return diag.FromErr(err)
		}
		resourceIDs, err := common.GetResourceIDsByOrder(bssV2Client, orderId, 1)
		if err != nil {
			return diag.Errorf("error getting resource IDs: %s", err)
		}
		if len(resourceIDs) != 1 {
			return diag.Errorf("error getting resource IDs, more than 1 resources are get by order (%s)", orderId)
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
		pending = []string{"STARTJOBING", "WAITING_FOR_START", "CONFIGURATION"}
		target = []string{"FULL_TRANSFER_STARTED", "FULL_TRANSFER_COMPLETE", "INCRE_TRANSFER_STARTED"}
	case "start_later":
		pending = []string{"CONFIGURATION"}
		target = []string{"WAITING_FOR_START"}
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
		pending = []string{"STARTJOBING", "WAITING_FOR_START", "CHILD_TRANSFER_STARTING", "CHILD_TRANSFER_STARTED",
			"CHILD_TRANSFER_COMPLETE", "RELEASE_CHILD_TRANSFER_STARTED"}
		target = []string{"FULL_TRANSFER_STARTED", "FULL_TRANSFER_COMPLETE", "INCRE_TRANSFER_STARTED",
			"RELEASE_CHILD_TRANSFER_COMPLETE", "DELETED"}
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
		NodeType:         d.Get("node_type").(string),
		SourceEndpoint:   *sourceDb,
		TargetEndpoint:   *targetDb,
		SubnetId:         subnetId,
		Tags:             utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		SysTags:          utils.BuildSysTags(enterpriseProjectID),
		MasterAz:         d.Get("master_az").(string),
		SlaveAz:          d.Get("slave_az").(string),
		PublciIpList:     buildPublicIpListParam(d.Get("public_ip_list").([]interface{})),
		IsOpenFastClean:  d.Get("is_open_fast_clean").(bool),
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

func buildPublicIpListParam(publicIpList []interface{}) []jobs.PublciIpList {
	if len(publicIpList) == 0 {
		return nil
	}

	publicIps := make([]jobs.PublciIpList, 0, len(publicIpList))
	for _, v := range publicIpList {
		tmp := v.(map[string]interface{})
		publicIps = append(publicIps, jobs.PublciIpList{
			Id:       tmp["id"].(string),
			PublicIp: tmp["public_ip"].(string),
			Type:     tmp["type"].(string),
		})
	}

	return publicIps
}

func buildDbConfigParamter(d *schema.ResourceData, dbType, projectId string) (*jobs.Endpoint, error) {
	configRaw := d.Get(dbType).([]interface{})[0].(map[string]interface{})
	configs := jobs.Endpoint{
		DbType:              configRaw["engine_type"].(string),
		Ip:                  configRaw["ip"].(string),
		DbName:              configRaw["name"].(string),
		DbUser:              configRaw["user"].(string),
		DbPassword:          configRaw["password"].(string),
		DbPort:              golangsdk.IntToPointer(configRaw["port"].(int)),
		InstanceId:          configRaw["instance_id"].(string),
		Region:              configRaw["region"].(string),
		VpcId:               configRaw["vpc_id"].(string),
		SubnetId:            configRaw["subnet_id"].(string),
		ProjectId:           projectId,
		SslCertPassword:     configRaw["ssl_cert_password"].(string),
		SslCertCheckSum:     configRaw["ssl_cert_check_sum"].(string),
		SslCertKey:          configRaw["ssl_cert_key"].(string),
		SslCertName:         configRaw["ssl_cert_name"].(string),
		SslLink:             utils.Bool(configRaw["ssl_enabled"].(bool)),
		KafkaSecurityConfig: buildKafkaSecurityConfigParamter(configRaw["kafka_security_config"].([]interface{})),
	}
	return &configs, nil
}

func buildKafkaSecurityConfigParamter(kafkaSecurityConfig []interface{}) *jobs.KafkaSecurityConfig {
	if len(kafkaSecurityConfig) == 0 {
		return nil
	}
	params, ok := kafkaSecurityConfig[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return &jobs.KafkaSecurityConfig{
		Type:                  params["type"].(string),
		SaslMechanism:         params["sasl_mechanism"].(string),
		TrustStoreKeyName:     params["trust_store_key_name"].(string),
		TrustStoreKey:         params["trust_store_key"].(string),
		TrustStorePassword:    params["trust_store_password"].(string),
		EndpointAlgorithm:     params["endpoint_algorithm"].(string),
		DelegationTokens:      params["delegation_tokens"].(bool),
		EnableKeyStore:        params["enable_key_store"].(bool),
		KeyStoreKeyName:       params["key_store_key_name"].(string),
		KeyStoreKey:           params["key_store_key"].(string),
		KeyStorePassword:      params["key_store_password"].(string),
		SetPrivateKeyPassword: params["set_private_key_password"].(bool),
		KeyPassword:           params["key_password"].(string),
	}
}

func setDbInfoToState(d *schema.ResourceData, endpoint jobs.Endpoint, fieldName string) error {
	result := make([]interface{}, 1)
	// IP sometimes will not same as input, if input 1, will return 2
	item := map[string]interface{}{
		"engine_type":           endpoint.DbType,
		"ip":                    d.Get(fieldName + ".0.ip"),
		"port":                  endpoint.DbPort,
		"password":              endpoint.DbPassword,
		"user":                  endpoint.DbUser,
		"instance_id":           endpoint.InstanceId,
		"name":                  endpoint.DbName,
		"region":                endpoint.Region,
		"vpc_id":                endpoint.VpcId,
		"subnet_id":             endpoint.SubnetId,
		"ssl_cert_password":     endpoint.SslCertPassword,
		"ssl_cert_check_sum":    endpoint.SslCertCheckSum,
		"ssl_cert_key":          endpoint.SslCertKey,
		"ssl_cert_name":         endpoint.SslCertName,
		"ssl_enabled":           endpoint.SslLink,
		"security_group_id":     endpoint.SecurityGroupId,
		"kafka_security_config": flattenKafkaSecurityConfig(endpoint.KafkaSecurityConfig),
	}
	result[0] = item
	// lintignore:R001
	return d.Set(fieldName, result)
}

func flattenKafkaSecurityConfig(kafkaSecurityConfig *jobs.KafkaSecurityConfig) interface{} {
	if kafkaSecurityConfig == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"type":                     kafkaSecurityConfig.Type,
			"sasl_mechanism":           kafkaSecurityConfig.SaslMechanism,
			"trust_store_key_name":     kafkaSecurityConfig.TrustStoreKeyName,
			"trust_store_key":          kafkaSecurityConfig.TrustStoreKey,
			"trust_store_password":     kafkaSecurityConfig.TrustStorePassword,
			"endpoint_algorithm":       kafkaSecurityConfig.EndpointAlgorithm,
			"delegation_tokens":        kafkaSecurityConfig.DelegationTokens,
			"enable_key_store":         kafkaSecurityConfig.EnableKeyStore,
			"key_store_key_name":       kafkaSecurityConfig.KeyStoreKeyName,
			"key_store_key":            kafkaSecurityConfig.KeyStoreKey,
			"key_store_password":       kafkaSecurityConfig.KeyStorePassword,
			"set_private_key_password": kafkaSecurityConfig.SetPrivateKeyPassword,
			"key_password":             kafkaSecurityConfig.KeyPassword,
		},
	}
}

// only mongodb have to set 0 for port, kafka will return error if input 0
func processPort(dbType string, port int) *int {
	if port == 0 && dbType != "mongodb" {
		return nil
	}
	return &port
}

func testConnections(client *golangsdk.ServiceClient, jobId string, opts jobs.CreateJobReq) (valid bool) {
	reqParams := jobs.TestConnectionsReq{
		Jobs: []jobs.TestEndPoint{
			{
				JobId:               jobId,
				NetType:             opts.NetType,
				EndPointType:        "so",
				ProjectId:           client.ProjectID,
				Region:              opts.SourceEndpoint.Region,
				VpcId:               opts.SourceEndpoint.VpcId,
				SubnetId:            opts.SourceEndpoint.SubnetId,
				DbType:              opts.SourceEndpoint.DbType,
				DbName:              opts.SourceEndpoint.DbName,
				Ip:                  opts.SourceEndpoint.Ip,
				DbUser:              opts.SourceEndpoint.DbUser,
				DbPassword:          opts.SourceEndpoint.DbPassword,
				DbPort:              processPort(opts.SourceEndpoint.DbType, *opts.SourceEndpoint.DbPort),
				SslLink:             opts.SourceEndpoint.SslLink,
				SslCertKey:          opts.SourceEndpoint.SslCertKey,
				SslCertName:         opts.SourceEndpoint.SslCertName,
				SslCertCheckSum:     opts.SourceEndpoint.SslCertCheckSum,
				SslCertPassword:     opts.SourceEndpoint.SslCertPassword,
				InstId:              opts.SourceEndpoint.InstanceId,
				KafkaSecurityConfig: opts.SourceEndpoint.KafkaSecurityConfig,
			},
			{
				JobId:               jobId,
				NetType:             opts.NetType,
				EndPointType:        "ta",
				ProjectId:           client.ProjectID,
				Region:              opts.TargetEndpoint.Region,
				VpcId:               opts.TargetEndpoint.VpcId,
				SubnetId:            opts.TargetEndpoint.SubnetId,
				DbType:              opts.TargetEndpoint.DbType,
				DbName:              opts.TargetEndpoint.DbName,
				Ip:                  opts.TargetEndpoint.Ip,
				DbUser:              opts.TargetEndpoint.DbUser,
				DbPassword:          opts.TargetEndpoint.DbPassword,
				DbPort:              processPort(opts.TargetEndpoint.DbType, *opts.TargetEndpoint.DbPort),
				SslLink:             opts.TargetEndpoint.SslLink,
				SslCertKey:          opts.TargetEndpoint.SslCertKey,
				SslCertName:         opts.TargetEndpoint.SslCertName,
				SslCertCheckSum:     opts.TargetEndpoint.SslCertCheckSum,
				SslCertPassword:     opts.TargetEndpoint.SslCertPassword,
				InstId:              opts.TargetEndpoint.InstanceId,
				KafkaSecurityConfig: opts.TargetEndpoint.KafkaSecurityConfig,
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

func processIpAndPort(ip, port string) string {
	if strings.Contains(ip, ",") || strings.Contains(ip, ":") {
		return ip
	}
	return ip + ":" + port
}

func testConnectionsForDualAZ(client *golangsdk.ServiceClient, jobId string, opts jobs.CreateJobReq) (valid bool) {
	sourceEndpoint := []jobs.PropertyParam{
		{
			DbName:              opts.SourceEndpoint.DbName,
			DbType:              opts.SourceEndpoint.DbType,
			NetType:             opts.NetType,
			EndPointType:        "so",
			Ip:                  processIpAndPort(opts.SourceEndpoint.Ip, strconv.Itoa(*opts.SourceEndpoint.DbPort)),
			DbUser:              opts.SourceEndpoint.DbUser,
			DbPassword:          opts.SourceEndpoint.DbPassword,
			ProjectId:           client.ProjectID,
			Region:              opts.SourceEndpoint.Region,
			VpcId:               opts.SourceEndpoint.VpcId,
			SubnetId:            opts.SourceEndpoint.SubnetId,
			InstId:              opts.SourceEndpoint.InstanceId,
			SslLink:             opts.SourceEndpoint.SslLink,
			SslCertKey:          opts.SourceEndpoint.SslCertKey,
			SslCertName:         opts.SourceEndpoint.SslCertName,
			SslCertCheckSum:     opts.SourceEndpoint.SslCertCheckSum,
			KafkaSecurityConfig: opts.SourceEndpoint.KafkaSecurityConfig,
		},
	}
	targetEndpoint := []jobs.PropertyParam{
		{
			DbName:              opts.TargetEndpoint.DbName,
			DbType:              opts.TargetEndpoint.DbType,
			NetType:             opts.NetType,
			EndPointType:        "ta",
			Ip:                  processIpAndPort(opts.TargetEndpoint.Ip, strconv.Itoa(*opts.TargetEndpoint.DbPort)),
			DbUser:              opts.TargetEndpoint.DbUser,
			DbPassword:          opts.TargetEndpoint.DbPassword,
			ProjectId:           client.ProjectID,
			Region:              opts.TargetEndpoint.Region,
			VpcId:               opts.TargetEndpoint.VpcId,
			SubnetId:            opts.TargetEndpoint.SubnetId,
			InstId:              opts.TargetEndpoint.InstanceId,
			SslLink:             opts.TargetEndpoint.SslLink,
			SslCertKey:          opts.TargetEndpoint.SslCertKey,
			SslCertName:         opts.TargetEndpoint.SslCertName,
			SslCertCheckSum:     opts.TargetEndpoint.SslCertCheckSum,
			KafkaSecurityConfig: opts.TargetEndpoint.KafkaSecurityConfig,
		},
	}

	sourceEndpointJson, _ := json.Marshal(sourceEndpoint)
	targetEndpointJson, _ := json.Marshal(targetEndpoint)
	reqParams := jobs.TestClusterConnectionsReq{
		Jobs: []jobs.TestJob{
			{
				Action:   "testConnection",
				JobId:    jobId,
				Property: string(sourceEndpointJson),
			},
			{
				Action:   "testConnection",
				JobId:    jobId,
				Property: string(targetEndpointJson),
			},
		},
	}
	rsp, err := jobs.TestClusterConnections(client, reqParams)
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
		return fmt.Errorf("error waiting for DRS job(%s) precheck mode(%s) to be completed: %s", jobId, precheckMode, err)
	}
	return nil
}
