package oms

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	oms "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2"
	omsmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OMS PUT /v2/{project_id}/taskgroups/{group_id}/retry
// @API OMS PUT /v2/{project_id}/taskgroups/{group_id}/start
// @API OMS PUT /v2/{project_id}/taskgroups/{group_id}/stop
// @API OMS PUT /v2/{project_id}/taskgroups/{group_id}/update
// @API OMS DELETE /v2/{project_id}/taskgroups/{group_id}
// @API OMS GET /v2/{project_id}/taskgroups/{group_id}
// @API OMS POST /v2/{project_id}/taskgroups
func ResourceMigrationTaskGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrationTaskGroupCreate,
		ReadContext:   resourceMigrationTaskGroupRead,
		UpdateContext: resourceMigrationTaskGroupUpdate,
		DeleteContext: resourceMigrationTaskGroupDelete,

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
						"data_source": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "Aliyun",
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
						"data_source": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "HEC",
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
			},
			"enable_kms": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"migrate_since": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"object_overwrite_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"consistency_check": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_requester_pays": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"retry", "start", "stop",
				}, false),
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
				Computed: true,
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
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"success_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fail_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"complete_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildTaskGroupSrcNodeOpts(rawSrcNode []interface{}) *omsmodel.TaskGroupSrcNode {
	if len(rawSrcNode) != 1 {
		return nil
	}
	srcNode := rawSrcNode[0].(map[string]interface{})

	srcNodeOpts := omsmodel.TaskGroupSrcNode{
		CloudType: utils.StringIgnoreEmpty(srcNode["data_source"].(string)),
		Region:    utils.StringIgnoreEmpty(srcNode["region"].(string)),
		Ak:        utils.StringIgnoreEmpty(srcNode["access_key"].(string)),
		Sk:        utils.StringIgnoreEmpty(srcNode["secret_key"].(string)),
		AppId:     utils.StringIgnoreEmpty(srcNode["app_id"].(string)),
		Bucket:    utils.StringIgnoreEmpty(srcNode["bucket"].(string)),
	}

	if srcNode["list_file_bucket"].(string) != "" {
		srcNodeOpts.ListFile = &omsmodel.ListFile{
			ObsBucket:   srcNode["list_file_bucket"].(string),
			ListFileKey: srcNode["list_file_key"].(string),
		}
	} else {
		srcNodeObjects := srcNode["object"].([]interface{})
		s := make([]string, 0, len(srcNodeObjects))
		for _, val := range srcNodeObjects {
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

func buildTaskGroupDstNodeOpts(conf *config.Config, rawDstNode []interface{}) (*omsmodel.TaskGroupDstNode, error) {
	if len(rawDstNode) != 1 {
		return nil, nil
	}
	dstNode := rawDstNode[0].(map[string]interface{})
	ak, err := getTaskGroupDstAccessKey(conf, dstNode)
	if err != nil {
		return nil, err
	}
	sk, err := getTaskGroupDstSecretKey(conf, dstNode)
	if err != nil {
		return nil, err
	}

	dstNodeOpts := omsmodel.TaskGroupDstNode{
		Region:     dstNode["region"].(string),
		Ak:         ak,
		Sk:         sk,
		CloudType:  utils.StringIgnoreEmpty(dstNode["data_source"].(string)),
		Bucket:     dstNode["bucket"].(string),
		SavePrefix: utils.StringIgnoreEmpty(dstNode["save_prefix"].(string)),
	}

	return &dstNodeOpts, nil
}

func getTaskGroupDstAccessKey(conf *config.Config, dstNode map[string]interface{}) (string, error) {
	if ak := dstNode["access_key"].(string); ak != "" {
		return ak, nil
	}
	if ak := conf.AccessKey; ak != "" {
		return ak, nil
	}
	return "", fmt.Errorf("unable to find access_key")
}

func getTaskGroupDstSecretKey(conf *config.Config, dstNode map[string]interface{}) (string, error) {
	if sk := dstNode["secret_key"].(string); sk != "" {
		return sk, nil
	}
	if sk := conf.SecretKey; sk != "" {
		return sk, nil
	}
	return "", fmt.Errorf("unable to find secret_key")
}

func buildTaskGroupCreateOpts(conf *config.Config, d *schema.ResourceData) (*omsmodel.CreateTaskGroupReq, error) {
	var taskType omsmodel.CreateTaskGroupReqTaskType
	if err := taskType.UnmarshalJSON([]byte(d.Get("type").(string))); err != nil {
		return nil, fmt.Errorf("error parsing the argument type: %s", err)
	}

	dstNodeOpts, err := buildTaskGroupDstNodeOpts(conf, d.Get("destination_object").([]interface{}))
	if err != nil {
		return nil, err
	}

	sourceCdn, err := buildSourceCdnOpts(d.Get("source_cdn").([]interface{}))
	if err != nil {
		return nil, err
	}

	var migrateSinceOpt *int64
	if v, ok := d.GetOk("migrate_since"); ok {
		migrateSince, err := utils.FormatUTCTimeStamp(v.(string))
		if err != nil {
			return nil, err
		}
		migrateSinceOpt = &migrateSince
	}

	createOpts := &omsmodel.CreateTaskGroupReq{
		SrcNode:                     buildTaskGroupSrcNodeOpts(d.Get("source_object").([]interface{})),
		Description:                 utils.StringIgnoreEmpty(d.Get("description").(string)),
		DstNode:                     dstNodeOpts,
		EnableFailedObjectRecording: utils.Bool(d.Get("enable_failed_object_recording").(bool)),
		EnableKms:                   d.Get("enable_kms").(bool),
		TaskType:                    &taskType,
		BandwidthPolicy:             buildBandwidthPolicyOpts(d.Get("bandwidth_policy").([]interface{})),
		SourceCdn:                   sourceCdn,
		MigrateSince:                migrateSinceOpt,
		EnableRequesterPays:         utils.Bool(d.Get("enable_requester_pays").(bool)),
	}

	if v, ok := d.GetOk("object_overwrite_mode"); ok {
		var objectOverwriteMode omsmodel.CreateTaskGroupReqObjectOverwriteMode
		if err := objectOverwriteMode.UnmarshalJSON([]byte(v.(string))); err != nil {
			return nil, fmt.Errorf("error parsing the argument object_overwrite_mode: %s", err)
		}
		createOpts.ObjectOverwriteMode = &objectOverwriteMode
	}

	if v, ok := d.GetOk("consistency_check"); ok {
		var consistencyCheck omsmodel.CreateTaskGroupReqConsistencyCheck
		if err := consistencyCheck.UnmarshalJSON([]byte(v.(string))); err != nil {
			return nil, fmt.Errorf("error parsing the argument consistency_check: %s", err)
		}
		createOpts.ConsistencyCheck = &consistencyCheck
	}

	return createOpts, nil
}

type TaskGroupActionConfig struct {
	Action  string
	Ctx     context.Context
	Conf    *config.Config
	GroupID string
}

func resourceMigrationTaskGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.HcOmsV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}
	createOpts, err := buildTaskGroupCreateOpts(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateTaskGroup(&omsmodel.CreateTaskGroupRequest{Body: createOpts})
	if err != nil {
		return diag.Errorf("error creating OMS migration task group: %s", err)
	}
	if resp.GroupId == nil {
		return diag.Errorf("unable to find the task group ID")
	}
	groupID := *resp.GroupId
	d.SetId(groupID)

	err = waitForTaskGroupStartedOrCompleted(ctx, client, groupID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for task group (%s) started or completed: %s", groupID, err)
	}

	if action, ok := d.GetOk("action"); ok && action.(string) == "stop" {
		actionConfig := &TaskGroupActionConfig{
			Action:  action.(string),
			Ctx:     ctx,
			Conf:    conf,
			GroupID: groupID,
		}
		if err := handleMigrationTaskGroupAction(actionConfig, client, d); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceMigrationTaskGroupRead(ctx, d, meta)
}

func handleMigrationTaskGroupAction(actionConfig *TaskGroupActionConfig, client *oms.OmsClient, d *schema.ResourceData) error {
	var err error
	switch actionConfig.Action {
	case "retry":
		err = retryMigrationTaskGroup(actionConfig, client, d)
	case "start":
		err = startMigrationTaskGroup(actionConfig, client, d)
	case "stop":
		err = stopMigrationTaskGroup(actionConfig, client, d)
	default:
		err = fmt.Errorf("invalid argument action(%s), valid values are retry, start and stop", actionConfig.Action)
	}
	return err
}

func retryMigrationTaskGroup(actionConfig *TaskGroupActionConfig, client *oms.OmsClient, d *schema.ResourceData) error {
	srcNode := make(map[string]interface{})
	dstNode := make(map[string]interface{})
	if sourceObjects := d.Get("source_object").([]interface{}); len(sourceObjects) > 0 {
		srcNode = sourceObjects[0].(map[string]interface{})
	}
	if destinationObjects := d.Get("destination_object").([]interface{}); len(destinationObjects) > 0 {
		dstNode = destinationObjects[0].(map[string]interface{})
	}

	dstAk, err := getTaskGroupDstAccessKey(actionConfig.Conf, dstNode)
	if err != nil {
		return err
	}

	dstSk, err := getTaskGroupDstSecretKey(actionConfig.Conf, dstNode)
	if err != nil {
		return err
	}

	var sourceCdnAuthenticationKey *string
	if sourceCDNs := d.Get("source_cdn").([]interface{}); len(sourceCDNs) > 0 {
		sourceCdn := sourceCDNs[0].(map[string]interface{})
		sourceCdnAuthenticationKey = utils.String(sourceCdn["authentication_key"].(string))
	}

	retryTaskGroupReq := &omsmodel.RetryTaskGroupReq{
		SrcAk:                      utils.StringIgnoreEmpty(srcNode["access_key"].(string)),
		SrcSk:                      utils.StringIgnoreEmpty(srcNode["secret_key"].(string)),
		DstAk:                      &dstAk,
		DstSk:                      &dstSk,
		SourceCdnAuthenticationKey: sourceCdnAuthenticationKey,
	}

	retryTaskGroupRequest := &omsmodel.RetryTaskGroupRequest{
		GroupId: actionConfig.GroupID,
		Body:    retryTaskGroupReq,
	}
	_, err = client.RetryTaskGroup(retryTaskGroupRequest)
	if err != nil {
		return fmt.Errorf("error retry OMS migration task group: %s", err)
	}

	err = waitForTaskGroupStartedOrCompleted(actionConfig.Ctx, client, actionConfig.GroupID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for task group (%s) started or completed: %s", actionConfig.GroupID, err)
	}
	return nil
}

func startMigrationTaskGroup(actionConfig *TaskGroupActionConfig, client *oms.OmsClient, d *schema.ResourceData) error {
	srcNode := make(map[string]interface{})
	dstNode := make(map[string]interface{})
	if sourceObjects := d.Get("source_object").([]interface{}); len(sourceObjects) > 0 {
		srcNode = sourceObjects[0].(map[string]interface{})
	}
	if destinationObjects := d.Get("destination_object").([]interface{}); len(destinationObjects) > 0 {
		dstNode = destinationObjects[0].(map[string]interface{})
	}

	dstAk, err := getTaskGroupDstAccessKey(actionConfig.Conf, dstNode)
	if err != nil {
		return err
	}

	dstSk, err := getTaskGroupDstSecretKey(actionConfig.Conf, dstNode)
	if err != nil {
		return err
	}

	var sourceCdnAuthenticationKey *string
	if sourceCDNs := d.Get("source_cdn").([]interface{}); len(sourceCDNs) > 0 {
		sourceCdn := sourceCDNs[0].(map[string]interface{})
		if sourceCdn["authentication_key"].(string) != "" {
			sourceCdnAuthenticationKey = utils.String(sourceCdn["authentication_key"].(string))
		}
	}

	startTaskGroupReq := &omsmodel.StartTaskGroupReq{
		SrcAk:                      utils.StringIgnoreEmpty(srcNode["access_key"].(string)),
		SrcSk:                      utils.StringIgnoreEmpty(srcNode["secret_key"].(string)),
		DstAk:                      dstAk,
		DstSk:                      dstSk,
		SourceCdnAuthenticationKey: sourceCdnAuthenticationKey,
	}

	startTaskGroupRequest := &omsmodel.StartTaskGroupRequest{
		GroupId: actionConfig.GroupID,
		Body:    startTaskGroupReq,
	}
	_, err = client.StartTaskGroup(startTaskGroupRequest)
	if err != nil {
		return fmt.Errorf("error start OMS migration task group: %s", err)
	}

	err = waitForTaskGroupStartedOrCompleted(actionConfig.Ctx, client, actionConfig.GroupID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for task group (%s) started or completed: %s", actionConfig.GroupID, err)
	}
	return nil
}

func stopMigrationTaskGroup(actionConfig *TaskGroupActionConfig, client *oms.OmsClient, d *schema.ResourceData) error {
	stopTaskGroupRequest := &omsmodel.StopTaskGroupRequest{GroupId: actionConfig.GroupID}
	_, err := client.StopTaskGroup(stopTaskGroupRequest)
	if err != nil {
		return fmt.Errorf("error stop OMS migration task group: %s", err)
	}

	err = waitForTaskGroupStopped(actionConfig.Ctx, client, actionConfig.GroupID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for task group (%s) stopped: %s", actionConfig.GroupID, err)
	}
	return nil
}

func resourceMigrationTaskGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.HcOmsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	groupID := d.Id()
	resp, err := client.ShowTaskGroup(&omsmodel.ShowTaskGroupRequest{GroupId: groupID})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OMS migration task group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", resp.TaskType.Value()),
		d.Set("enable_kms", resp.EnableKms),
		d.Set("description", resp.Description),
		d.Set("object_overwrite_mode", resp.ObjectOverwriteMode.Value()),
		d.Set("consistency_check", resp.ConsistencyCheck.Value()),
		d.Set("enable_requester_pays", resp.EnableRequesterPays),
		d.Set("enable_failed_object_recording", resp.EnableFailedObjectRecording),
		d.Set("bandwidth_policy", flattenBandwidthPolicy(resp.BandwidthPolicy)),
		d.Set("source_cdn", flattenSourceCdn(resp.SourceCdn)),
		d.Set("status", resp.Status),
		d.Set("total_time", resp.TotalTime),
		d.Set("total_num", resp.TotalNum),
		d.Set("success_num", resp.SuccessNum),
		d.Set("fail_num", resp.FailNum),
		d.Set("total_size", resp.TotalSize),
		d.Set("complete_size", resp.CompleteSize),
	)

	if resp.MigrateSince != nil {
		mErr = multierror.Append(mErr,
			d.Set("migrate_since", utils.FormatTimeStampUTC(*resp.MigrateSince)),
		)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OMS migration task group fields: %s", err)
	}

	return nil
}

func resourceMigrationTaskGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.HcOmsV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	groupID := d.Id()
	if d.HasChange("bandwidth_policy") {
		updateBandwidthPolicyOpts := buildBandwidthPolicyOpts(d.Get("bandwidth_policy").([]interface{}))
		updateTaskGroupRequest := &omsmodel.UpdateTaskGroupRequest{
			GroupId: groupID,
			Body: &omsmodel.UpdateBandwidthPolicyReq{
				BandwidthPolicy: *updateBandwidthPolicyOpts,
			},
		}
		_, err := client.UpdateTaskGroup(updateTaskGroupRequest)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("action") {
		actionConfig := &TaskGroupActionConfig{
			Action:  d.Get("action").(string),
			Ctx:     ctx,
			Conf:    conf,
			GroupID: groupID,
		}
		if err := handleMigrationTaskGroupAction(actionConfig, client, d); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceMigrationTaskGroupRead(ctx, d, meta)
}

func resourceMigrationTaskGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.HcOmsV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	groupID := d.Id()
	resp, err := client.ShowTaskGroup(&omsmodel.ShowTaskGroupRequest{GroupId: groupID})
	if err != nil {
		return diag.Errorf("error retrieving OMS migration task group: %s", err)
	}

	if resp.Status == nil {
		return diag.Errorf("unable to find the status OMS migration task group: %s", groupID)
	}
	status := *resp.Status

	// the status is creating, which cannot be stopped or deleted
	if status == 1 {
		err := waitForTaskGroupStartedOrCompleted(ctx, client, groupID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for task group (%s) started or completed: %s", groupID, err)
		}
		return resourceMigrationTaskGroupDelete(ctx, d, meta)
	}

	// the status is monitoring, must stop the running task group before deleting it
	if status == 2 {
		_, err := client.StopTaskGroup(&omsmodel.StopTaskGroupRequest{GroupId: groupID})
		if err != nil {
			if responseErr, ok := err.(*sdkerr.ServiceResponseError); !ok || responseErr.ErrorCode != "OMS.0066" {
				return diag.Errorf("error stopping OMS migration task group: %s", err)
			}
		} else {
			err := waitForTaskGroupStopped(ctx, client, groupID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.Errorf("error waiting for task group (%s) stopped: %s", groupID, err)
			}
		}
	}

	if status == 7 {
		err := waitForTaskGroupStopped(ctx, client, groupID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for task group (%s) stopped: %s", groupID, err)
		}
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err = client.DeleteTaskGroup(&omsmodel.DeleteTaskGroupRequest{GroupId: groupID})
		if err == nil {
			return nil
		}
		responseErr, ok := err.(*sdkerr.ServiceResponseError)
		if !ok {
			return resource.NonRetryableError(err)
		}

		// ErrorCode "OMS.0063" means the task group is in progress. This ErrorCode is not accurate, we need retry it.
		if responseErr.ErrorCode == "OMS.0063" {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return diag.Errorf("error deleting OMS migration task group: %s", err)
	}

	// wait for delete
	err = waitForTaskGroupDeleted(ctx, client, groupID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getTaskGroupStatus(client *oms.OmsClient, groupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		groupGet, err := client.ShowTaskGroup(&omsmodel.ShowTaskGroupRequest{GroupId: groupID})
		if err != nil {
			responseErr, ok := err.(*sdkerr.ServiceResponseError)
			if ok && responseErr.StatusCode == 404 {
				return groupGet, "DELETED", nil
			}
			return nil, "", err
		}

		if *groupGet.Status == 4 {
			// the task group is failed, read the error reason field
			if reason := groupGet.ErrorReason; reason != nil {
				err = fmt.Errorf("migration task group is failed,"+
					" error_code is: %s, error_msg is: %s", *reason.ErrorCode, *reason.ErrorMsg)
				return groupGet, "4", err
			}
		}

		status := strconv.Itoa(int(*groupGet.Status))
		return groupGet, status, nil
	}
}

func waitForTaskGroupStartedOrCompleted(ctx context.Context, client *oms.OmsClient, groupID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"0", "1"},
		Target:       []string{"2", "6"},
		Refresh:      getTaskGroupStatus(client, groupID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForTaskGroupStopped(ctx context.Context, client *oms.OmsClient, groupID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"7"},
		Target:       []string{"3", "6"},
		Refresh:      getTaskGroupStatus(client, groupID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForTaskGroupDeleted(ctx context.Context, client *oms.OmsClient, groupID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"8", "9"},
		Target:       []string{"DELETED"},
		Refresh:      getTaskGroupStatus(client, groupID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
