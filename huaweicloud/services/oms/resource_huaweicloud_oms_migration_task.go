package oms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

// @API OMS POST /v2/{project_id}/tasks/{task_id}/start
// @API OMS POST /v2/{project_id}/tasks/{task_id}/stop
// @API OMS GET /v2/{project_id}/tasks/{task_id}
// @API OMS DELETE /v2/{project_id}/tasks/{task_id}
// @API OMS POST /v2/{project_id}/tasks
// @API OMS PUT /v2/{project_id}/tasks/{task_id}/bandwidth-policy
func ResourceMigrationTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrationTaskCreate,
		ReadContext:   resourceMigrationTaskRead,
		UpdateContext: resourceMigrationTaskUpdate,
		DeleteContext: resourceMigrationTaskDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_object": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"bucket": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"access_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"secret_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
							ForceNew:  true,
							RequiredWith: []string{
								"source_object.0.access_key",
							},
						},
						"object": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							ExactlyOneOf: []string{
								"source_object.0.list_file_bucket",
							},
						},
						"security_token": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
							ForceNew:  true,
							RequiredWith: []string{
								"source_object.0.access_key",
							},
						},
						"data_source": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AWS", "Azure", "Aliyun", "Tencent", "HuaweiCloud", "QingCloud", "KingsoftCloud",
								"Baidu", "Qiniu", "URLSource", "UCloud",
							}, false),
							Default: "Aliyun",
						},
						"app_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"list_file_bucket": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							RequiredWith: []string{
								"source_object.0.list_file_key",
							},
						},
						"list_file_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							RequiredWith: []string{
								"source_object.0.list_file_bucket",
							},
						},
						"list_file_num": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"json_auth_file": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"destination_object": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"access_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"secret_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
							ForceNew:  true,
						},
						"security_token": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
							ForceNew:  true,
						},
						"save_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"list", "url_list", "object", "prefix",
				}, false),
			},
			"start_task": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enable_kms": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"migrate_since": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_restore": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_failed_object_recording": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"bandwidth_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_bandwidth": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"start": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"source_cdn": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http", "https",
							}, false),
						},
						"authentication_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "NONE",
						},
						"authentication_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"smn_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_urn": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"trigger_conditions": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"language": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"message_template_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"enable_metadata_migration": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"dst_storage_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"task_priority": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"consistency_check": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_requester_pays": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"object_overwrite_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildSrcNodeOpts(rawSrcNode []interface{}) map[string]interface{} {
	if len(rawSrcNode) != 1 {
		return nil
	}
	srcNode := rawSrcNode[0].(map[string]interface{})

	srcNodeOpts := map[string]interface{}{
		"cloud_type":     utils.ValueIgnoreEmpty(srcNode["data_source"]),
		"region":         utils.ValueIgnoreEmpty(srcNode["region"]),
		"ak":             utils.ValueIgnoreEmpty(srcNode["access_key"]),
		"sk":             utils.ValueIgnoreEmpty(srcNode["secret_key"]),
		"security_token": utils.ValueIgnoreEmpty(srcNode["security_token"]),
		"app_id":         utils.ValueIgnoreEmpty(srcNode["app_id"]),
		"bucket":         utils.ValueIgnoreEmpty(srcNode["bucket"]),
		"json_auth_file": utils.ValueIgnoreEmpty(srcNode["json_auth_file"]),
	}

	if srcNode["list_file_bucket"].(string) != "" {
		srcNodeOpts["list_file"] = map[string]interface{}{
			"obs_bucket":    srcNode["list_file_bucket"],
			"list_file_key": srcNode["list_file_key"],
			"list_file_num": utils.ValueIgnoreEmpty(srcNode["list_file_num"]),
		}
	}

	if len(srcNode["object"].([]interface{})) != 0 {
		s := make([]string, 0, len(srcNode["object"].([]interface{})))
		for _, val := range srcNode["object"].([]interface{}) {
			if strVal, ok := val.(string); ok {
				s = append(s, strVal)
			}
		}

		// "" will be ignored, add it here if the input is [ "" ]
		if len(s) == 0 {
			s = append(s, "")
		}

		srcNodeOpts["object_key"] = &s
	}

	return srcNodeOpts
}

func buildDstNodeOpts(cfg *config.Config, rawDstNode []interface{}) (map[string]interface{}, error) {
	if len(rawDstNode) != 1 {
		return nil, nil
	}
	dstNode := rawDstNode[0].(map[string]interface{})
	ak, err := getDstAccessKey(cfg, dstNode)
	if err != nil {
		return nil, err
	}
	sk, err := getDstSecretKey(cfg, dstNode)
	if err != nil {
		return nil, err
	}
	securityToken := getDstSecurityToken(cfg, dstNode)

	dstNodeOpts := map[string]interface{}{
		"region":         dstNode["region"].(string),
		"ak":             ak,
		"sk":             sk,
		"security_token": securityToken,
		"bucket":         dstNode["bucket"].(string),
		"save_prefix":    utils.ValueIgnoreEmpty(dstNode["save_prefix"].(string)),
	}

	return dstNodeOpts, nil
}

func getDstAccessKey(cfg *config.Config, dstNode map[string]interface{}) (string, error) {
	if ak := dstNode["access_key"].(string); ak != "" {
		return ak, nil
	}
	if ak := cfg.AccessKey; ak != "" {
		return ak, nil
	}
	return "", fmt.Errorf("unable to find access_key")
}

func getDstSecretKey(cfg *config.Config, dstNode map[string]interface{}) (string, error) {
	if sk := dstNode["secret_key"].(string); sk != "" {
		return sk, nil
	}
	if sk := cfg.SecretKey; sk != "" {
		return sk, nil
	}
	return "", fmt.Errorf("unable to find secret_key")
}

func getDstSecurityToken(cfg *config.Config, dstNode map[string]interface{}) *string {
	if securityToken := dstNode["security_token"].(string); securityToken != "" {
		return utils.String(securityToken)
	}
	if securityToken := cfg.SecurityToken; securityToken != "" {
		return utils.String(securityToken)
	}
	return nil
}

func buildBandwidthPolicyOpts(rawBandwidthPolicy []interface{}) []map[string]interface{} {
	if len(rawBandwidthPolicy) < 1 {
		return nil
	}

	bandwidthPolicyOpts := make([]map[string]interface{}, len(rawBandwidthPolicy))
	for i, rawPolicy := range rawBandwidthPolicy {
		policy := rawPolicy.(map[string]interface{})
		bandwidthPolicyOpts[i] = map[string]interface{}{
			"max_bandwidth": float64(policy["max_bandwidth"].(int) * 1024 * 1024),
			"start":         policy["start"].(string),
			"end":           policy["end"].(string),
		}
	}

	return bandwidthPolicyOpts
}

func buildSourceCdnOpts(rawSourceCdn []interface{}) map[string]interface{} {
	if len(rawSourceCdn) != 1 {
		return nil
	}
	sourceCdn := rawSourceCdn[0].(map[string]interface{})

	sourceCdnOpts := map[string]interface{}{
		"domain":              sourceCdn["domain"].(string),
		"authentication_key":  utils.ValueIgnoreEmpty(sourceCdn["authentication_key"].(string)),
		"protocol":            sourceCdn["protocol"].(string),
		"authentication_type": sourceCdn["authentication_type"].(string),
	}

	return sourceCdnOpts
}

func buildSmnConfigOpts(rawSmnConfig []interface{}) map[string]interface{} {
	if len(rawSmnConfig) != 1 {
		return nil
	}
	smnInfo := rawSmnConfig[0].(map[string]interface{})

	smnInfoOpts := map[string]interface{}{
		"topic_urn":          smnInfo["topic_urn"].(string),
		"trigger_conditions": utils.ExpandToStringList(smnInfo["trigger_conditions"].([]interface{})),
		"language":           utils.ValueIgnoreEmpty(smnInfo["language"].(string)),
	}

	return smnInfoOpts
}

func buildcreateTaskBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	bodyParams := map[string]interface{}{
		"task_type":                      d.Get("type"),
		"src_node":                       buildSrcNodeOpts(d.Get("source_object").([]interface{})),
		"enable_kms":                     d.Get("enable_kms").(bool),
		"description":                    utils.ValueIgnoreEmpty(d.Get("description").(string)),
		"bandwidth_policy":               buildBandwidthPolicyOpts(d.Get("bandwidth_policy").([]interface{})),
		"smn_config":                     buildSmnConfigOpts(d.Get("smn_config").([]interface{})),
		"enable_restore":                 d.Get("enable_restore"),
		"enable_failed_object_recording": d.Get("enable_failed_object_recording"),
		"source_cdn":                     buildSourceCdnOpts(d.Get("source_cdn").([]interface{})),
		"enable_metadata_migration":      d.Get("enable_metadata_migration").(bool),
		"enable_requester_pays":          d.Get("enable_requester_pays").(bool),
		"task_priority":                  utils.ValueIgnoreEmpty(d.Get("task_priority").(string)),
		"consistency_check":              utils.ValueIgnoreEmpty(d.Get("consistency_check").(string)),
		"object_overwrite_mode":          utils.ValueIgnoreEmpty(d.Get("object_overwrite_mode").(string)),
		"dst_storage_policy":             utils.ValueIgnoreEmpty(d.Get("dst_storage_policy").(string)),
	}

	dstNodeOpts, err := buildDstNodeOpts(cfg, d.Get("destination_object").([]interface{}))
	if err != nil {
		return nil, err
	}

	bodyParams["dst_node"] = dstNodeOpts

	if d.Get("migrate_since").(string) != "" {
		migrateSince, err := utils.FormatUTCTimeStamp(d.Get("migrate_since").(string))
		if err != nil {
			return nil, err
		}
		bodyParams["migrate_since"] = migrateSince
	}

	return bodyParams, nil
}

func resourceMigrationTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createTaskHttpUrl = "v2/{project_id}/tasks"
		createTaskProduct = "oms"
	)
	createTaskClient, err := cfg.NewServiceClient(createTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	createTaskPath := createTaskClient.Endpoint + createTaskHttpUrl
	createTaskPath = strings.ReplaceAll(createTaskPath, "{project_id}", createTaskClient.ProjectID)

	createTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts, err := buildcreateTaskBodyParams(d, cfg)
	if err != nil {
		return nil
	}

	log.Printf("[DEBUG] Create Task options: %#v", createOpts)
	createTaskOpt.JSONBody = utils.RemoveNil(createOpts)
	createTaskResp, err := createTaskClient.Request("POST", createTaskPath, &createTaskOpt)
	if err != nil {
		return diag.Errorf("error creating OMS migration task: %s", err)
	}

	createTaskRespBody, err := utils.FlattenResponse(createTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createTaskRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating OMS migration task: ID is not found in API response")
	}

	taskID := strconv.FormatInt(int64(id.(float64)), 10)
	d.SetId(taskID)

	err = waitForTaskStartedORCompleted(ctx, createTaskClient, taskID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for task (%s) started: %s", taskID, err)
	}

	if !d.Get("start_task").(bool) {
		err := stopTask(createTaskClient, d)
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForTaskStopped(ctx, createTaskClient, taskID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for task (%s) stopped: %s", taskID, err)
		}
	}

	return resourceMigrationTaskRead(ctx, d, meta)
}

func resourceMigrationTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getTaskHttpUrl = "v2/{project_id}/tasks/{task_id}"
		getTaskProduct = "oms"
	)
	getTaskClient, err := cfg.NewServiceClient(getTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	getTaskPath := getTaskClient.Endpoint + getTaskHttpUrl
	getTaskPath = strings.ReplaceAll(getTaskPath, "{project_id}", getTaskClient.ProjectID)
	getTaskPath = strings.ReplaceAll(getTaskPath, "{task_id}", d.Id())

	getTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTaskResp, err := getTaskClient.Request("GET", getTaskPath, &getTaskOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "OMS.1009"),
			"error retrieving OMS migration task")
	}

	getTaskRespBody, err := utils.FlattenResponse(getTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("type", utils.PathSearch("task_type", getTaskRespBody, nil)),
		d.Set("enable_kms", utils.PathSearch("enable_kms", getTaskRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getTaskRespBody, nil)),
		d.Set("bandwidth_policy", flattenBandwidthPolicy(getTaskRespBody)),
		d.Set("source_cdn", flattenSourceCdn(getTaskRespBody)),
		d.Set("enable_restore", utils.PathSearch("enable_restore", getTaskRespBody, nil)),
		d.Set("enable_failed_object_recording", utils.PathSearch("enable_failed_object_recording", getTaskRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getTaskRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getTaskRespBody, nil)),
	)

	if migrateSince := utils.PathSearch("migrate_since", getTaskRespBody, float64(0)).(float64); migrateSince != 0 {
		mErr = multierror.Append(mErr,
			d.Set("migrate_since", utils.FormatTimeStampUTC(int64(migrateSince))),
		)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OMS migration task fields: %s", err)
	}

	return nil
}

func hasErrorCode(err error, expectCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode := utils.PathSearch("error_code", response, nil)
			if errorCode == nil {
				log.Printf("[WARN] failed to parse error_code from response body")
			}

			if errorCode == expectCode {
				return true
			}
		}
	}

	return false
}

func resourceMigrationTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateTaskProduct = "oms"
	)
	updateTaskClient, err := cfg.NewServiceClient(updateTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	taskID := d.Id()

	if d.HasChange("bandwidth_policy") {
		err := updateBandwidthPolicy(updateTaskClient, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("start_task") {
		if d.Get("start_task").(bool) {
			err := startTask(updateTaskClient, cfg, d)
			if err != nil {
				return diag.FromErr(err)
			}

			err = waitForTaskStartedORCompleted(ctx, updateTaskClient, taskID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.Errorf("error waiting for task (%s) started: %s", taskID, err)
			}
		} else {
			err := stopTask(updateTaskClient, d)
			if err != nil {
				return diag.FromErr(err)
			}

			err = waitForTaskStopped(ctx, updateTaskClient, taskID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.Errorf("error waiting for task (%s) stopped: %s", taskID, err)
			}
		}
	}

	return resourceMigrationTaskRead(ctx, d, meta)
}

func buildUpdateBandwidthPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_policy": buildBandwidthPolicyOpts(d.Get("bandwidth_policy").([]interface{})),
	}

	return bodyParams
}

func updateBandwidthPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateBandwidthPolicyHttpUrl := "v2/{project_id}/tasks/{task_id}/bandwidth-policy"
	updateBandwidthPolicyPath := client.Endpoint + updateBandwidthPolicyHttpUrl
	updateBandwidthPolicyPath = strings.ReplaceAll(updateBandwidthPolicyPath, "{project_id}", client.ProjectID)
	updateBandwidthPolicyPath = strings.ReplaceAll(updateBandwidthPolicyPath, "{task_id}", d.Id())

	updateBandwidthPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateBandwidthPolicyOpt.JSONBody = utils.RemoveNil(buildUpdateBandwidthPolicyBodyParams(d))
	_, err := client.Request("PUT", updateBandwidthPolicyPath, &updateBandwidthPolicyOpt)
	if err != nil {
		return fmt.Errorf("error updating bandwidth policy OMS migration task: %s", err)
	}

	return nil
}

func startTask(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) error {
	startTaskHttpUrl := "v2/{project_id}/tasks/{task_id}/start"
	startTaskPath := client.Endpoint + startTaskHttpUrl
	startTaskPath = strings.ReplaceAll(startTaskPath, "{project_id}", client.ProjectID)
	startTaskPath = strings.ReplaceAll(startTaskPath, "{task_id}", d.Id())

	startTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	opts, err := buildStartTaskBodyParams(cfg, d)
	if err != nil {
		return err
	}

	startTaskOpt.JSONBody = utils.RemoveNil(opts)
	_, err = client.Request("POST", startTaskPath, &startTaskOpt)
	if err != nil {
		return fmt.Errorf("error starting OMS migration task: %s", err)
	}

	return nil
}

func stopTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stopTaskHttpUrl := "v2/{project_id}/tasks/{task_id}/stop"
	stopTaskPath := client.Endpoint + stopTaskHttpUrl
	stopTaskPath = strings.ReplaceAll(stopTaskPath, "{project_id}", client.ProjectID)
	stopTaskPath = strings.ReplaceAll(stopTaskPath, "{task_id}", d.Id())

	stopTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("POST", stopTaskPath, &stopTaskOpt)
	if err != nil {
		return fmt.Errorf("error stopping OMS migration task: %s", err)
	}

	return nil
}

func waitForTaskStopped(ctx context.Context, client *golangsdk.ServiceClient, taskID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"0", "7"},
		Target:     []string{"3"},
		Refresh:    getTaskStatus(client, taskID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func buildStartTaskBodyParams(cfg *config.Config, d *schema.ResourceData) (map[string]interface{}, error) {
	srcNode := make(map[string]interface{})
	dstNode := make(map[string]interface{})
	if sourceObjects := d.Get("source_object").([]interface{}); len(sourceObjects) > 0 {
		srcNode = sourceObjects[0].(map[string]interface{})
	}
	if destinationObjects := d.Get("destination_object").([]interface{}); len(destinationObjects) > 0 {
		dstNode = destinationObjects[0].(map[string]interface{})
	}

	dstAk, err := getDstAccessKey(cfg, dstNode)
	if err != nil {
		return nil, err
	}
	dstSk, err := getDstSecretKey(cfg, dstNode)
	if err != nil {
		return nil, err
	}
	dstSecurityToken := getDstSecurityToken(cfg, dstNode)

	startTaskReqOpt := map[string]interface{}{
		"src_ak":             utils.ValueIgnoreEmpty(srcNode["access_key"].(string)),
		"src_sk":             utils.ValueIgnoreEmpty(srcNode["secret_key"].(string)),
		"src_security_token": utils.ValueIgnoreEmpty(srcNode["security_token"].(string)),
		"dst_ak":             dstAk,
		"dst_sk":             dstSk,
		"dst_security_token": dstSecurityToken,
	}
	if sourceCDNs := d.Get("source_cdn").([]interface{}); len(sourceCDNs) > 0 {
		sourceCdn := sourceCDNs[0].(map[string]interface{})
		startTaskReqOpt["source_cdn_authentication_key"] = utils.String(sourceCdn["authentication_key"].(string))
	}
	return startTaskReqOpt, nil
}

func resourceMigrationTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	taskID := d.Id()

	var (
		deleteTaskHttpUrl = "v2/{project_id}/tasks/{task_id}"
		deleteTaskProduct = "oms"
	)
	client, err := cfg.NewServiceClient(deleteTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	deleteTaskPath := client.Endpoint + deleteTaskHttpUrl
	deleteTaskPath = strings.ReplaceAll(deleteTaskPath, "{project_id}", client.ProjectID)
	deleteTaskPath = strings.ReplaceAll(deleteTaskPath, "{task_id}", taskID)

	deleteTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// get the status of task before deleting it
	getTaskResp, err := client.Request("GET", deleteTaskPath, &deleteTaskOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "OMS.1009"),
			"error retrieving OMS migration task")
	}

	getTaskRespBody, err := utils.FlattenResponse(getTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}
	status := utils.PathSearch("status", getTaskRespBody, float64(0)).(float64)

	// must stop the running task before deleting it
	if status == 2 {
		err = stopTask(client, d)
		if err != nil {
			// ErrorCode "OMS.0066" means the task is not running, don't need to stop it before deleting
			if !hasErrorCode(err, "OMS.0066") {
				return diag.Errorf("error stopping OMS migration task: %s", err)
			}
		} else {
			err := waitForTaskStopped(ctx, client, taskID, d.Timeout(schema.TimeoutDelete))
			if err != nil {
				return diag.Errorf("error waiting for task (%s) stopped: %s", taskID, err)
			}
		}
	}

	if status == 7 {
		err := waitForTaskStopped(ctx, client, taskID, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.Errorf("error waiting for task (%s) stopped: %s", taskID, err)
		}
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err = client.Request("DELETE", deleteTaskPath, &deleteTaskOpt)
		if err == nil {
			return nil
		}

		// ErrorCode "OMS.0063" means the task is in progress. This ErrorCode is not accurate, we need retry it.
		if hasErrorCode(err, "OMS.0063") {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return diag.Errorf("error deleting OMS migration task: %s", err)
	}

	return nil
}

func getTaskStatus(client *golangsdk.ServiceClient, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getTaskHttpUrl = "v2/{project_id}/tasks/{task_id}"
		)

		getTaskPath := client.Endpoint + getTaskHttpUrl
		getTaskPath = strings.ReplaceAll(getTaskPath, "{project_id}", client.ProjectID)
		getTaskPath = strings.ReplaceAll(getTaskPath, "{task_id}", taskId)

		getTaskOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getTaskResp, err := client.Request("GET", getTaskPath, &getTaskOpt)
		if err != nil {
			return nil, "", err
		}

		getTaskRespBody, err := utils.FlattenResponse(getTaskResp)
		if err != nil {
			return nil, "", err
		}

		status := strconv.Itoa(int(utils.PathSearch("status", getTaskRespBody, float64(0)).(float64)))
		return getTaskRespBody, status, nil
	}
}

func waitForTaskStartedORCompleted(ctx context.Context, client *golangsdk.ServiceClient, taskID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"0", "1"},
		Target:     []string{"2", "5"},
		Refresh:    getTaskStatus(client, taskID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func flattenBandwidthPolicy(getTaskRespBody interface{}) []map[string]interface{} {
	bandwidthPolicy := utils.PathSearch("bandwidth_policy", getTaskRespBody, nil)
	if bandwidthPolicy == nil {
		return nil
	}

	bandwidthPolicyResult := make([]map[string]interface{}, len(bandwidthPolicy.([]interface{})))
	for i, policy := range bandwidthPolicy.([]interface{}) {
		bandwidthPolicyResult[i] = map[string]interface{}{
			"max_bandwidth": utils.PathSearch("max_bandwidth", policy, float64(0)).(float64) / (1024 * 1024),
			"start":         utils.PathSearch("start", policy, nil),
			"end":           utils.PathSearch("end", policy, nil),
		}
	}
	return bandwidthPolicyResult
}

func flattenSourceCdn(getTaskRespBody interface{}) []map[string]interface{} {
	sourceCdn := utils.PathSearch("source_cdn", getTaskRespBody, nil)
	if sourceCdn == nil {
		return nil
	}

	sourceCdnResult := map[string]interface{}{
		"domain":              utils.PathSearch("domain", sourceCdn, nil),
		"protocol":            utils.PathSearch("protocol", sourceCdn, nil),
		"authentication_type": utils.PathSearch("authentication_type", sourceCdn, nil),
	}

	return []map[string]interface{}{sourceCdnResult}
}
