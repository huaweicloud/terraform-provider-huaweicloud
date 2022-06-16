package oms

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	v2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2"
	oms "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceMigrationTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrationTaskCreate,
		ReadContext:   resourceMigrationTaskRead,
		UpdateContext: resourceMigrationTaskUpdate,
		DeleteContext: resourceMigrationTaskDelete,

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
				ForceNew: true,
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
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 200),
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
							ValidateFunc: validation.StringInSlice([]string{
								"NONE", "QINIU_PRIVATE_AUTHENTICATION", "ALIYUN_OSS_A", "ALIYUN_OSS_B", "ALIYUN_OSS_C",
								"KSYUN_PRIVATE_AUTHENTICATION",
							}, false),
							Default: "NONE",
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
					},
				},
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

func buildSrcNodeOpts(rawSrcNode []interface{}) *oms.SrcNodeReq {
	if len(rawSrcNode) != 1 {
		return nil
	}
	srcNode := rawSrcNode[0].(map[string]interface{})

	srcNodeOpts := oms.SrcNodeReq{
		CloudType:     utils.StringIgnoreEmpty(srcNode["data_source"].(string)),
		Region:        utils.StringIgnoreEmpty(srcNode["region"].(string)),
		Ak:            utils.StringIgnoreEmpty(srcNode["access_key"].(string)),
		Sk:            utils.StringIgnoreEmpty(srcNode["secret_key"].(string)),
		SecurityToken: utils.StringIgnoreEmpty(srcNode["security_token"].(string)),
		AppId:         utils.StringIgnoreEmpty(srcNode["app_id"].(string)),
		Bucket:        utils.StringIgnoreEmpty(srcNode["bucket"].(string)),
	}

	if srcNode["list_file_bucket"].(string) != "" {
		srcNodeOpts.ListFile = &oms.ListFile{
			ObsBucket:   srcNode["list_file_bucket"].(string),
			ListFileKey: srcNode["list_file_key"].(string),
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

		srcNodeOpts.ObjectKey = &s
	}

	return &srcNodeOpts
}

func buildDstNodeOpts(config *config.Config, rawDstNode []interface{}) (*oms.DstNodeReq, error) {
	if len(rawDstNode) != 1 {
		return nil, nil
	}
	dstNode := rawDstNode[0].(map[string]interface{})

	dstNodeOpts := oms.DstNodeReq{
		Region:     dstNode["region"].(string),
		Ak:         dstNode["access_key"].(string),
		Sk:         dstNode["secret_key"].(string),
		Bucket:     dstNode["bucket"].(string),
		SavePrefix: utils.StringIgnoreEmpty(dstNode["save_prefix"].(string)),
	}

	if dstNode["access_key"].(string) != "" {
		dstNodeOpts.Ak = dstNode["access_key"].(string)
	} else if config.AccessKey != "" {
		dstNodeOpts.Ak = config.AccessKey
	} else {
		return nil, fmt.Errorf("unable to find access_key")
	}

	if dstNode["secret_key"].(string) != "" {
		dstNodeOpts.Sk = dstNode["secret_key"].(string)
	} else if config.SecretKey != "" {
		dstNodeOpts.Sk = config.SecretKey
	} else {
		return nil, fmt.Errorf("unable to find secret_key")
	}

	if dstNode["security_token"].(string) != "" {
		dstNodeOpts.SecurityToken = utils.String(dstNode["security_token"].(string))
	} else if config.SecurityToken != "" {
		dstNodeOpts.SecurityToken = utils.String(config.SecurityToken)
	}

	return &dstNodeOpts, nil
}

func buildBandwidthPolicyOpts(rawBandwidthPolicy []interface{}) *[]oms.BandwidthPolicyDto {
	if len(rawBandwidthPolicy) < 1 {
		return nil
	}

	bandwidthPolicyOpts := make([]oms.BandwidthPolicyDto, len(rawBandwidthPolicy))
	for i, rawPolicy := range rawBandwidthPolicy {
		policy := rawPolicy.(map[string]interface{})
		bandwidthPolicyOpts[i] = oms.BandwidthPolicyDto{
			MaxBandwidth: int64(policy["max_bandwidth"].(int) * 1024 * 1024),
			Start:        policy["start"].(string),
			End:          policy["end"].(string),
		}
	}

	return &bandwidthPolicyOpts
}

func buildSourceCdnOpts(rawSourceCdn []interface{}) (*oms.SourceCdnReq, error) {
	if len(rawSourceCdn) != 1 {
		return nil, nil
	}
	sourceCdn := rawSourceCdn[0].(map[string]interface{})

	sourceCdnOpts := oms.SourceCdnReq{
		Domain:            sourceCdn["domain"].(string),
		AuthenticationKey: utils.String(sourceCdn["authentication_key"].(string)),
	}

	if sourceCdn["protocol"].(string) == "http" {
		sourceCdnOpts.Protocol = oms.GetSourceCdnReqProtocolEnum().HTTP
	} else {
		sourceCdnOpts.Protocol = oms.GetSourceCdnReqProtocolEnum().HTTPS
	}

	var authenticationType oms.SourceCdnReqAuthenticationType
	if err := authenticationType.UnmarshalJSON([]byte(sourceCdn["authentication_type"].(string))); err != nil {
		return nil, fmt.Errorf("error parsing the argument authentication_type: %s", err)
	}

	return &sourceCdnOpts, nil
}

func buildSmnConfigOpts(rawSmnConfig []interface{}) *oms.SmnConfig {
	if len(rawSmnConfig) != 1 {
		return nil
	}
	smnInfo := rawSmnConfig[0].(map[string]interface{})

	smnInfoOpts := oms.SmnConfig{
		TopicUrn:          smnInfo["topic_urn"].(string),
		TriggerConditions: utils.ExpandToStringList(smnInfo["trigger_conditions"].([]interface{})),
	}
	var language oms.SmnConfigLanguage
	if smnInfo["language"].(string) == "zh-cn" {
		language = oms.GetSmnConfigLanguageEnum().ZH_CN
	} else {
		language = oms.GetSmnConfigLanguageEnum().EN_US
	}
	smnInfoOpts.Language = &language

	return &smnInfoOpts
}

func resourceMigrationTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcOmsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	var taskType oms.CreateTaskReqTaskType
	if err := taskType.UnmarshalJSON([]byte(d.Get("type").(string))); err != nil {
		return diag.Errorf("error parsing the argument type: %s", err)
	}

	createOpts := oms.CreateTaskReq{
		TaskType:                    &taskType,
		SrcNode:                     buildSrcNodeOpts(d.Get("source_object").([]interface{})),
		EnableKms:                   utils.Bool(d.Get("enable_kms").(bool)),
		Description:                 utils.StringIgnoreEmpty(d.Get("description").(string)),
		BandwidthPolicy:             buildBandwidthPolicyOpts(d.Get("bandwidth_policy").([]interface{})),
		SmnConfig:                   buildSmnConfigOpts(d.Get("smn_config").([]interface{})),
		EnableRestore:               utils.Bool(d.Get("enable_restore").(bool)),
		EnableFailedObjectRecording: utils.Bool(d.Get("enable_failed_object_recording").(bool)),
	}

	if d.Get("migrate_since").(string) != "" {
		migrateSince, err := utils.FormatUTCTimeStamp(d.Get("migrate_since").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.MigrateSince = &migrateSince
	}

	sourceCdn, err := buildSourceCdnOpts(d.Get("source_cdn").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts.SourceCdn = sourceCdn

	dstNodeOpts, err := buildDstNodeOpts(config, d.Get("destination_object").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts.DstNode = dstNodeOpts

	log.Printf("[DEBUG] Create Task options: %#v", createOpts)

	resp, err := client.CreateTask(&oms.CreateTaskRequest{Body: &createOpts})
	if err != nil {
		return diag.Errorf("error creating OMS migration task: %s", err)
	}

	if resp.Id == nil {
		return diag.Errorf("unable to find the task ID")
	}

	taskID := strconv.FormatInt(*resp.Id, 10)
	d.SetId(taskID)

	err = waitForTaskStartedORCompleted(ctx, client, *resp.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for task (%s) started: %s", taskID, err)
	}

	if !d.Get("start_task").(bool) {
		_, err = client.StopTask(&oms.StopTaskRequest{TaskId: *resp.Id})
		if err != nil {
			return diag.Errorf("error stopping OMS migration task: %s", err)
		}
	}

	return resourceMigrationTaskRead(ctx, d, meta)
}

func resourceMigrationTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcOmsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	taskID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("the task ID must be integer: %s", err)
	}

	resp, err := client.ShowTask(&oms.ShowTaskRequest{TaskId: taskID})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OMS migration task")
	}
	log.Printf("[DEBUG] Retrieved Task %s: %#v", d.Id(), resp)

	mErr := multierror.Append(nil,
		d.Set("type", resp.TaskType.Value()),
		d.Set("enable_kms", resp.EnableKms),
		d.Set("description", resp.Description),
		d.Set("bandwidth_policy", flattenBandwidthPolicy(resp.BandwidthPolicy)),
		d.Set("source_cdn", flattenSourceCdn(resp.SourceCdn)),
		d.Set("enable_restore", resp.EnableRestore),
		d.Set("enable_failed_object_recording", resp.EnableFailedObjectRecording),
		d.Set("name", resp.Name),
		d.Set("status", resp.Status),
	)

	if resp.MigrateSince != nil {
		mErr = multierror.Append(mErr,
			d.Set("migrate_since", utils.FormatTimeStampUTC(*resp.MigrateSince)),
		)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OMS migration task fields: %s", err)
	}

	return nil
}

func resourceMigrationTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcOmsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	taskID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("the task ID must be integer: %s", err)
	}

	if d.HasChange("bandwidth_policy") {
		updateBandwidthPolicyOpts := buildBandwidthPolicyOpts(d.Get("bandwidth_policy").([]interface{}))
		updateBandwidthPolicyReq := oms.UpdateBandwidthPolicyRequest{
			TaskId: taskID,
			Body: &oms.UpdateBandwidthPolicyReq{
				BandwidthPolicy: *updateBandwidthPolicyOpts,
			},
		}
		_, err = client.UpdateBandwidthPolicy(&updateBandwidthPolicyReq)
		if err != nil {
			diag.Errorf("error retrieving OMS migration task: %s", err)
		}
	}

	if d.HasChange("start_task") {
		if d.Get("start_task").(bool) {
			_, err := client.StartTask(&oms.StartTaskRequest{TaskId: taskID})
			if err != nil {
				diag.Errorf("error starting OMS migration task: %s", err)
			}
		} else {
			_, err := client.StopTask(&oms.StopTaskRequest{TaskId: taskID})
			if err != nil {
				diag.Errorf("error stopping OMS migration task: %s", err)
			}
		}
	}

	return resourceMigrationTaskRead(ctx, d, meta)
}

func resourceMigrationTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcOmsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	taskID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("the task ID must be integer: %s", err)
	}

	// must stop the running task before deleting it
	resp, err := client.ShowTask(&oms.ShowTaskRequest{TaskId: taskID})
	if err != nil {
		return diag.Errorf("error retrieving OMS migration task: %s", err)
	}

	if resp.Status == nil {
		diag.Errorf("unable to find the status OMS migration task: %d", taskID)
	}

	if *resp.Status == 2 {
		_, err = client.StopTask(&oms.StopTaskRequest{TaskId: taskID})
		if err != nil {
			// ErrorCode "OMS.0066" means the task is not running, don't need top stop it before deleting
			if responseErr, ok := err.(*sdkerr.ServiceResponseError); !ok || responseErr.ErrorCode != "OMS.0066" {
				return diag.Errorf("error stopping OMS migration task: %s", err)
			}
		}
	}

	_, err = client.DeleteTask(&oms.DeleteTaskRequest{TaskId: taskID})
	if err != nil {
		return diag.Errorf("error deleting OMS migration task: %s", err)
	}

	return nil
}

func getTaskStatus(client *v2.OmsClient, taskId int64) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		taskGet, err := client.ShowTask(&oms.ShowTaskRequest{TaskId: taskId})
		if err != nil {
			return nil, "", err
		}

		status := strconv.Itoa(int(*taskGet.Status))
		return taskGet, status, nil
	}
}

func waitForTaskStartedORCompleted(ctx context.Context, client *v2.OmsClient, taskID int64, timeout time.Duration) error {
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

func flattenBandwidthPolicy(bandwidthPolicy *[]oms.BandwidthPolicyDto) []map[string]interface{} {
	if bandwidthPolicy == nil {
		return nil
	}

	bandwidthPolicyResult := make([]map[string]interface{}, len(*bandwidthPolicy))
	for i, policy := range *bandwidthPolicy {
		bandwidthPolicyResult[i] = map[string]interface{}{
			"max_bandwidth": policy.MaxBandwidth / (1024 * 1024),
			"start":         policy.Start,
			"end":           policy.End,
		}
	}
	return bandwidthPolicyResult
}

func flattenSourceCdn(sourceCdn *oms.SourceCdnResp) []map[string]interface{} {
	if sourceCdn == nil {
		return nil
	}

	sourceCdnResult := map[string]interface{}{
		"domain":              sourceCdn.Domain,
		"protocol":            sourceCdn.Protocol,
		"authentication_type": sourceCdn.AuthenticationType.Value(),
	}

	return []map[string]interface{}{sourceCdnResult}
}
