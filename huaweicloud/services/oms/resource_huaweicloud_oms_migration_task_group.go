package oms

import (
	"context"
	"errors"
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

func buildTaskGroupSrcNodeOpts(rawSrcNode []interface{}) map[string]interface{} {
	if len(rawSrcNode) != 1 {
		return nil
	}
	srcNode := rawSrcNode[0].(map[string]interface{})

	srcNodeOpts := map[string]interface{}{
		"cloud_type": utils.ValueIgnoreEmpty(srcNode["data_source"]),
		"region":     utils.ValueIgnoreEmpty(srcNode["region"]),
		"ak":         utils.ValueIgnoreEmpty(srcNode["access_key"]),
		"sk":         utils.ValueIgnoreEmpty(srcNode["secret_key"]),
		"app_id":     utils.ValueIgnoreEmpty(srcNode["app_id"]),
		"bucket":     utils.ValueIgnoreEmpty(srcNode["bucket"]),
	}

	if srcNode["list_file_bucket"].(string) != "" {
		srcNodeOpts["list_file"] = map[string]interface{}{
			"obs_bucket":    srcNode["list_file_bucket"].(string),
			"list_file_key": srcNode["list_file_key"].(string),
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
		srcNodeOpts["object_key"] = &s
	}

	return srcNodeOpts
}

func buildTaskGroupDstNodeOpts(conf *config.Config, rawDstNode []interface{}) (map[string]interface{}, error) {
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

	dstNodeOpts := map[string]interface{}{
		"region":      dstNode["region"],
		"ak":          ak,
		"sk":          sk,
		"cloud_type":  utils.ValueIgnoreEmpty(dstNode["data_source"]),
		"bucket":      dstNode["bucket"],
		"save_prefix": utils.ValueIgnoreEmpty(dstNode["save_prefix"]),
	}

	return dstNodeOpts, nil
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

func buildTaskGroupCreateOpts(conf *config.Config, d *schema.ResourceData) (map[string]interface{}, error) {
	dstNodeOpts, err := buildTaskGroupDstNodeOpts(conf, d.Get("destination_object").([]interface{}))
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

	createOpts := map[string]interface{}{
		"task_type":                      d.Get("type"),
		"src_node":                       buildTaskGroupSrcNodeOpts(d.Get("source_object").([]interface{})),
		"description":                    utils.ValueIgnoreEmpty(d.Get("description")),
		"dst_node":                       dstNodeOpts,
		"enable_failed_object_recording": d.Get("enable_failed_object_recording"),
		"enable_kms":                     d.Get("enable_kms").(bool),
		"bandwidth_policy":               buildBandwidthPolicyOpts(d.Get("bandwidth_policy").([]interface{})),
		"source_cdn":                     buildSourceCdnOpts(d.Get("source_cdn").([]interface{})),
		"migrate_since":                  migrateSinceOpt,
		"enable_requester_pays":          utils.ValueIgnoreEmpty(d.Get("enable_requester_pays")),
		"object_overwrite_mode":          utils.ValueIgnoreEmpty(d.Get("object_overwrite_mode")),
		"consistency_check":              utils.ValueIgnoreEmpty(d.Get("consistency_check")),
		"enable_metadata_migration":      d.Get("enable_metadata_migration").(bool),
		"dst_storage_policy":             utils.ValueIgnoreEmpty(d.Get("dst_storage_policy").(string)),
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createTaskGroupHttpUrl = "v2/{project_id}/taskgroups"
		createTaskGroupProduct = "oms"
	)
	createTaskGroupClient, err := cfg.NewServiceClient(createTaskGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	createTaskGroupPath := createTaskGroupClient.Endpoint + createTaskGroupHttpUrl
	createTaskGroupPath = strings.ReplaceAll(createTaskGroupPath, "{project_id}", createTaskGroupClient.ProjectID)

	createTaskGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts, err := buildTaskGroupCreateOpts(cfg, d)
	if err != nil {
		return nil
	}

	log.Printf("[DEBUG] Create Task Group options: %#v", createOpts)
	createTaskGroupOpt.JSONBody = utils.RemoveNil(createOpts)
	createTaskGroupResp, err := createTaskGroupClient.Request("POST", createTaskGroupPath, &createTaskGroupOpt)
	if err != nil {
		return diag.Errorf("error creating OMS migration task group: %s", err)
	}

	createTaskGroupRespBody, err := utils.FlattenResponse(createTaskGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("group_id", createTaskGroupRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating OMS migration task group: ID is not found in API response")
	}

	groupID := id.(string)

	d.SetId(groupID)

	err = waitForTaskGroupStartedOrCompleted(ctx, createTaskGroupClient, groupID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for task group (%s) started or completed: %s", groupID, err)
	}

	if action, ok := d.GetOk("action"); ok && action.(string) == "stop" {
		actionConfig := &TaskGroupActionConfig{
			Action:  action.(string),
			Ctx:     ctx,
			Conf:    cfg,
			GroupID: groupID,
		}
		if err := handleMigrationTaskGroupAction(actionConfig, createTaskGroupClient, d); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceMigrationTaskGroupRead(ctx, d, meta)
}

func handleMigrationTaskGroupAction(actionConfig *TaskGroupActionConfig, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
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

func retryMigrationTaskGroup(actionConfig *TaskGroupActionConfig, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
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

	retryTaskOpts := map[string]interface{}{
		"src_ak":                        utils.StringIgnoreEmpty(srcNode["access_key"].(string)),
		"src_sk":                        utils.StringIgnoreEmpty(srcNode["secret_key"].(string)),
		"dst_ak":                        &dstAk,
		"dst_sk":                        &dstSk,
		"source_cdn_authentication_key": sourceCdnAuthenticationKey,
	}

	retryTaskGroupHttpUrl := "v2/{project_id}/taskgroups/{group_id}/retry"
	retryTaskGroupPath := client.Endpoint + retryTaskGroupHttpUrl
	retryTaskGroupPath = strings.ReplaceAll(retryTaskGroupPath, "{project_id}", client.ProjectID)
	retryTaskGroupPath = strings.ReplaceAll(retryTaskGroupPath, "{group_id}", actionConfig.GroupID)

	retryTaskGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	retryTaskGroupOpt.JSONBody = utils.RemoveNil(retryTaskOpts)
	_, err = client.Request("PUT", retryTaskGroupPath, &retryTaskGroupOpt)
	if err != nil {
		return fmt.Errorf("error retrying OMS migration task group: %s", err)
	}

	err = waitForTaskGroupStartedOrCompleted(actionConfig.Ctx, client, actionConfig.GroupID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for task group (%s) started or completed: %s", actionConfig.GroupID, err)
	}
	return nil
}

func startMigrationTaskGroup(actionConfig *TaskGroupActionConfig, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
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

	startTaskOpts := map[string]interface{}{
		"src_ak":                        utils.StringIgnoreEmpty(srcNode["access_key"].(string)),
		"src_sk":                        utils.StringIgnoreEmpty(srcNode["secret_key"].(string)),
		"dst_ak":                        &dstAk,
		"dst_sk":                        &dstSk,
		"source_cdn_authentication_key": sourceCdnAuthenticationKey,
	}

	startTaskGroupHttpUrl := "v2/{project_id}/taskgroups/{group_id}/start"
	startTaskGroupPath := client.Endpoint + startTaskGroupHttpUrl
	startTaskGroupPath = strings.ReplaceAll(startTaskGroupPath, "{project_id}", client.ProjectID)
	startTaskGroupPath = strings.ReplaceAll(startTaskGroupPath, "{group_id}", actionConfig.GroupID)

	startTaskGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	startTaskGroupOpt.JSONBody = utils.RemoveNil(startTaskOpts)
	_, err = client.Request("PUT", startTaskGroupPath, &startTaskGroupOpt)
	if err != nil {
		return fmt.Errorf("error starting OMS migration task group: %s", err)
	}

	err = waitForTaskGroupStartedOrCompleted(actionConfig.Ctx, client, actionConfig.GroupID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for task group (%s) started or completed: %s", actionConfig.GroupID, err)
	}
	return nil
}

func stopMigrationTaskGroup(actionConfig *TaskGroupActionConfig, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stopTaskGroupHttpUrl := "v2/{project_id}/taskgroups/{group_id}/stop"
	stopTaskGroupPath := client.Endpoint + stopTaskGroupHttpUrl
	stopTaskGroupPath = strings.ReplaceAll(stopTaskGroupPath, "{project_id}", client.ProjectID)
	stopTaskGroupPath = strings.ReplaceAll(stopTaskGroupPath, "{group_id}", actionConfig.GroupID)

	stopTaskGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("PUT", stopTaskGroupPath, &stopTaskGroupOpt)
	if err != nil {
		return fmt.Errorf("error stopping OMS migration task group: %s", err)
	}

	err = waitForTaskGroupStopped(actionConfig.Ctx, client, actionConfig.GroupID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for task group (%s) stopped: %s", actionConfig.GroupID, err)
	}
	return nil
}

func resourceMigrationTaskGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getTaskGroupHttpUrl = "v2/{project_id}/taskgroups/{group_id}"
		getTaskGroupProduct = "oms"
	)
	getTaskGroupClient, err := cfg.NewServiceClient(getTaskGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	getTaskGroupPath := getTaskGroupClient.Endpoint + getTaskGroupHttpUrl
	getTaskGroupPath = strings.ReplaceAll(getTaskGroupPath, "{project_id}", getTaskGroupClient.ProjectID)
	getTaskGroupPath = strings.ReplaceAll(getTaskGroupPath, "{group_id}", d.Id())

	getTaskGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTaskGroupResp, err := getTaskGroupClient.Request("GET", getTaskGroupPath, &getTaskGroupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OMS migration task group")
	}

	getTaskGroupRespBody, err := utils.FlattenResponse(getTaskGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", utils.PathSearch("task_type", getTaskGroupRespBody, nil)),
		d.Set("enable_kms", utils.PathSearch("enable_kms", getTaskGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getTaskGroupRespBody, nil)),
		d.Set("object_overwrite_mode", utils.PathSearch("object_overwrite_mode", getTaskGroupRespBody, nil)),
		d.Set("consistency_check", utils.PathSearch("consistency_check", getTaskGroupRespBody, nil)),
		d.Set("enable_requester_pays", utils.PathSearch("enable_requester_pays", getTaskGroupRespBody, nil)),
		d.Set("enable_failed_object_recording", utils.PathSearch("enable_failed_object_recording", getTaskGroupRespBody, nil)),
		d.Set("bandwidth_policy", flattenBandwidthPolicy(getTaskGroupRespBody)),
		d.Set("source_cdn", flattenSourceCdn(getTaskGroupRespBody)),
		d.Set("status", utils.PathSearch("status", getTaskGroupRespBody, nil)),
		d.Set("total_time", utils.PathSearch("total_time", getTaskGroupRespBody, nil)),
		d.Set("total_num", utils.PathSearch("total_num", getTaskGroupRespBody, nil)),
		d.Set("success_num", utils.PathSearch("success_num", getTaskGroupRespBody, nil)),
		d.Set("fail_num", utils.PathSearch("fail_num", getTaskGroupRespBody, nil)),
		d.Set("total_size", utils.PathSearch("total_size", getTaskGroupRespBody, nil)),
		d.Set("complete_size", utils.PathSearch("complete_size", getTaskGroupRespBody, nil)),
	)

	if migrateSince := utils.PathSearch("migrate_since", getTaskGroupRespBody, float64(0)).(float64); migrateSince != 0 {
		mErr = multierror.Append(mErr,
			d.Set("migrate_since", utils.FormatTimeStampUTC(int64(migrateSince))),
		)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OMS migration task group fields: %s", err)
	}

	return nil
}

func resourceMigrationTaskGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateTaskGroupProduct = "oms"
	)
	updateTaskGroupClient, err := cfg.NewServiceClient(updateTaskGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	groupID := d.Id()
	if d.HasChange("bandwidth_policy") {
		updateBandwidthPolicyHttpUrl := "v2/{project_id}/taskgroups/{group_id}/update"
		updateBandwidthPolicyPath := updateTaskGroupClient.Endpoint + updateBandwidthPolicyHttpUrl
		updateBandwidthPolicyPath = strings.ReplaceAll(updateBandwidthPolicyPath, "{project_id}", updateTaskGroupClient.ProjectID)
		updateBandwidthPolicyPath = strings.ReplaceAll(updateBandwidthPolicyPath, "{group_id}", groupID)

		updateTaskGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		updateTaskGroupOpt.JSONBody = utils.RemoveNil(buildUpdateBandwidthPolicyBodyParams(d))
		_, err = updateTaskGroupClient.Request("PUT", updateBandwidthPolicyPath, &updateTaskGroupOpt)
		if err != nil {
			return diag.Errorf("error updating OMS migration task group: %s", err)
		}
	}

	if d.HasChange("action") {
		actionConfig := &TaskGroupActionConfig{
			Action:  d.Get("action").(string),
			Ctx:     ctx,
			Conf:    cfg,
			GroupID: groupID,
		}
		if err := handleMigrationTaskGroupAction(actionConfig, updateTaskGroupClient, d); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceMigrationTaskGroupRead(ctx, d, meta)
}

func resourceMigrationTaskGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	groupID := d.Id()

	var (
		deleteTaskGroupHttpUrl = "v2/{project_id}/taskgroups/{group_id}"
		deleteTaskGroupProduct = "oms"
	)
	deleteTaskGroupClient, err := cfg.NewServiceClient(deleteTaskGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	deleteTaskGroupPath := deleteTaskGroupClient.Endpoint + deleteTaskGroupHttpUrl
	deleteTaskGroupPath = strings.ReplaceAll(deleteTaskGroupPath, "{project_id}", deleteTaskGroupClient.ProjectID)
	deleteTaskGroupPath = strings.ReplaceAll(deleteTaskGroupPath, "{group_id}", d.Id())

	deleteTaskGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTaskGroupResp, err := deleteTaskGroupClient.Request("GET", deleteTaskGroupPath, &deleteTaskGroupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OMS migration task group")
	}

	getTaskGroupRespBody, err := utils.FlattenResponse(getTaskGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	status := utils.PathSearch("status", getTaskGroupRespBody, float64(0)).(float64)

	// the status is creating, which cannot be stopped or deleted
	if status == 1 {
		err := waitForTaskGroupStartedOrCompleted(ctx, deleteTaskGroupClient, groupID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for task group (%s) started or completed: %s", groupID, err)
		}
		return resourceMigrationTaskGroupDelete(ctx, d, meta)
	}

	// the status is monitoring, must stop the running task group before deleting it
	if status == 2 {
		actionConfig := &TaskGroupActionConfig{
			Action:  "stop",
			Ctx:     ctx,
			Conf:    cfg,
			GroupID: groupID,
		}
		err := stopMigrationTaskGroup(actionConfig, deleteTaskGroupClient, d)
		if err != nil {
			if !hasErrorCode(err, "OMS.0066") {
				return diag.Errorf("error stopping OMS migration task group: %s", err)
			}
		} else {
			err := waitForTaskGroupStopped(ctx, deleteTaskGroupClient, groupID, d.Timeout(schema.TimeoutDelete))
			if err != nil {
				return diag.Errorf("error waiting for task group (%s) stopped: %s", groupID, err)
			}
		}
	}

	if status == 7 {
		err := waitForTaskGroupStopped(ctx, deleteTaskGroupClient, groupID, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.Errorf("error waiting for task group (%s) stopped: %s", groupID, err)
		}
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := deleteTaskGroupClient.Request("DELETE", deleteTaskGroupPath, &deleteTaskGroupOpt)
		if err == nil {
			return nil
		}

		// ErrorCode "OMS.0063" means the task group is in progress. This ErrorCode is not accurate, we need retry it.
		if hasErrorCode(err, "OMS.0063") {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return diag.Errorf("error deleting OMS migration task group: %s", err)
	}

	// wait for delete
	err = waitForTaskGroupDeleted(ctx, deleteTaskGroupClient, groupID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getTaskGroupGroupStatus(client *golangsdk.ServiceClient, groupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getTaskGroupHttpUrl = "v2/{project_id}/taskgroups/{group_id}"
		)

		getTaskGroupPath := client.Endpoint + getTaskGroupHttpUrl
		getTaskGroupPath = strings.ReplaceAll(getTaskGroupPath, "{project_id}", client.ProjectID)
		getTaskGroupPath = strings.ReplaceAll(getTaskGroupPath, "{group_id}", groupID)

		getTaskGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getTaskGroupResp, err := client.Request("GET", getTaskGroupPath, &getTaskGroupOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return getTaskGroupResp, "DELETED", nil
			}
			return nil, "", err
		}

		getTaskGroupRespBody, err := utils.FlattenResponse(getTaskGroupResp)
		if err != nil {
			return nil, "", err
		}

		statusRaw := utils.PathSearch("status", getTaskGroupRespBody, nil)
		if statusRaw == nil {
			return nil, "", errors.New("unable to find status in the API response")
		}

		if statusRaw.(float64) == 4 {
			// the task group is failed, read the error reason field
			if reason := utils.PathSearch("error_reason", getTaskGroupRespBody, nil); reason != nil {
				err = fmt.Errorf("migration task group is failed,"+
					" error_code is: %s, error_msg is: %s", utils.PathSearch("error_reason.error_code", getTaskGroupRespBody, nil),
					utils.PathSearch("error_reason.error_msg", getTaskGroupRespBody, nil))
				return getTaskGroupRespBody, "4", err
			}
		}

		status := strconv.Itoa(int(statusRaw.(float64)))
		return getTaskGroupRespBody, status, nil
	}
}

func waitForTaskGroupStartedOrCompleted(ctx context.Context, client *golangsdk.ServiceClient, groupID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"0", "1"},
		Target:       []string{"2", "6"},
		Refresh:      getTaskGroupGroupStatus(client, groupID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForTaskGroupStopped(ctx context.Context, client *golangsdk.ServiceClient, groupID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"7"},
		Target:       []string{"3", "6"},
		Refresh:      getTaskGroupGroupStatus(client, groupID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForTaskGroupDeleted(ctx context.Context, client *golangsdk.ServiceClient, groupID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"8", "9"},
		Target:       []string{"DELETED"},
		Refresh:      getTaskGroupGroupStatus(client, groupID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
