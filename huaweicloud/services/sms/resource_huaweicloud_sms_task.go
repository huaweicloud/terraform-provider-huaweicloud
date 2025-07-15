package sms

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/sms/v3/tasks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var taskNonUpdatableParams = []string{"auto_start", "use_ipv6", "start_network_check", "migrate_speed_limit",
	"over_speed_threshold", "is_need_consistency_check", "need_migration_test"}

// ResourceMigrateTask is the impl of huaweicloud_sms_task
// @API SMS GET /v3/sources/{id}
// @API SMS POST /v3/tasks
// @API SMS POST /v3/tasks/{id}/action
// @API SMS GET /v3/tasks/{id}
// @API SMS DELETE /v3/tasks/{id}
// @API ECS GET /v1/{project_id}/cloudservers/{server_id}
// @API SMS POST /v3/tasks/{task_id}/speed-limit
// @API SMS GET /v3/tasks/{task_id}/speed-limit
// @API SMS GET /v3/tasks/{task_id}/passphrase
// @API SMS POST /v3/tasks/{task_id}/configuration-setting
// @API SMS GET /v3/tasks/{task_id}/configuration-setting
func ResourceMigrateTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrateTaskCreate,
		ReadContext:   resourceMigrateTaskRead,
		UpdateContext: resourceMigrateTaskUpdate,
		DeleteContext: resourceMigrateTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(taskNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MIGRATE_FILE", "MIGRATE_BLOCK",
				}, true),
			},
			"os_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LINUX", "WINDOWS",
				}, true),
			},
			"source_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"project_id"},
				Computed:     true,
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"region"},
				Computed:     true,
			},
			"target_server_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"vm_template_id"},
				RequiredWith: []string{"migration_ip"},
			},
			"vm_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"target_server_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"device_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "schema: Required",
						},
						"used_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"physical_volumes": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"device_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"file_system": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"mount_point": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"index": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"used_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"migration_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"start_target_server": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"use_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"syncing": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start", "stop", "restart",
				}, false),
			},
			"auto_start": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"use_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"start_network_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"migrate_speed_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"over_speed_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"is_need_consistency_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"need_migration_test": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"speed_limit": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end": {
							Type:     schema.TypeString,
							Required: true,
						},
						"speed": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"over_speed_threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
			"configurations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config_status": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"target_server_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"migrate_speed": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"passphrase": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"migrate_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getProjectID(d *schema.ResourceData, cfg *config.Config, region string) string {
	var projectID string

	if v, ok := d.GetOk("project_id"); ok {
		projectID = v.(string)
	} else {
		// get project ID from config
		projectID = cfg.RegionProjectIDMap[region]
	}

	return projectID
}

func getProjectName(d *schema.ResourceData, cfg *config.Config) string {
	// get project name from config
	projectName := cfg.TenantName

	if v, ok := d.GetOk("region"); ok {
		region := v.(string)
		if region != cfg.Region {
			// seem the region specified in resource as the project name
			projectName = region
		}
	}
	return projectName
}

func buildCreateTaskBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	region := cfg.GetRegion(d)

	target, err := buildCreateTaskTargetServerBodyParams(d, cfg)
	if err != nil {
		return nil, err
	}

	_, existing := d.GetOk("target_server_id")
	bodyParams := map[string]interface{}{
		"name":                      "MigrationTask",
		"priority":                  1,
		"type":                      d.Get("type"),
		"os_type":                   d.Get("os_type"),
		"region_name":               region,
		"region_id":                 region,
		"project_name":              getProjectName(d, cfg),
		"project_id":                getProjectID(d, cfg, region),
		"source_server":             buildCreateTaskSourceServerBodyParams(d),
		"target_server":             target,
		"vm_template_id":            utils.ValueIgnoreEmpty(d.Get("vm_template_id")),
		"migration_ip":              utils.ValueIgnoreEmpty(d.Get("migration_ip")),
		"use_public_ip":             utils.ValueIgnoreEmpty(d.Get("use_public_ip")),
		"start_target_server":       utils.ValueIgnoreEmpty(d.Get("start_target_server")),
		"syncing":                   utils.ValueIgnoreEmpty(d.Get("syncing")),
		"exist_server":              existing,
		"auto_start":                utils.ValueIgnoreEmpty(d.Get("auto_start")),
		"use_ipv6":                  utils.ValueIgnoreEmpty(d.Get("use_ipv6")),
		"start_network_check":       utils.ValueIgnoreEmpty(d.Get("start_network_check")),
		"speed_limit":               utils.ValueIgnoreEmpty(d.Get("migrate_speed_limit")),
		"over_speed_threshold":      utils.ValueIgnoreEmpty(d.Get("over_speed_threshold")),
		"is_need_consistency_check": utils.ValueIgnoreEmpty(d.Get("is_need_consistency_check")),
		"need_migration_test":       utils.ValueIgnoreEmpty(d.Get("need_migration_test")),
	}

	return bodyParams, nil
}

func buildCreateTaskSourceServerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id": d.Get("source_server_id"),
	}

	return bodyParams
}

func buildCreateTaskTargetServerBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	bodyParams := make(map[string]interface{})

	if v, ok := d.GetOk("target_server_id"); ok {
		serverID := v.(string)
		ecsClient, err := cfg.ComputeV1Client(cfg.GetRegion(d))
		if err != nil {
			return nil, fmt.Errorf("error creating compute client: %s", err)
		}

		server, err := cloudservers.Get(ecsClient, serverID).Extract()
		if err != nil {
			return nil, fmt.Errorf("error retrieving ECS instance %s: %s", serverID, err)
		}

		bodyParams["name"] = server.Name
		bodyParams["vm_id"] = serverID
	}

	v, ok := d.GetOk("target_server_disks")
	if !ok {
		defaultDisks, err := buildCreateTaskTargetServerDisksDefaultBodyParams(d, cfg, d.Get("source_server_id").(string))
		if err != nil {
			return nil, err
		}
		bodyParams["disks"] = defaultDisks
	} else {
		bodyParams["disks"] = buildCreateTaskTargetServerDisksBodyParams(v)
	}

	return bodyParams, nil
}

func buildCreateTaskTargetServerDisksDefaultBodyParams(d *schema.ResourceData, cfg *config.Config,
	sid string) ([]map[string]interface{}, error) {
	smsClient, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating SMS client: %s", err)
	}

	log.Printf("[DEBUG] filtering SMS source servers by id %s", sid)
	sourceServer, err := getSourceServer(smsClient, sid)
	if err != nil {
		return nil, fmt.Errorf("unable to find the source server %s: %s", sid, err)
	}

	sourceServerDiskParams := utils.PathSearch("init_target_server.disks", sourceServer,
		make([]interface{}, 0)).([]interface{})
	params := make([]map[string]interface{}, len(sourceServerDiskParams))
	for i, v := range sourceServerDiskParams {
		raw := v.(map[string]interface{})
		params[i] = map[string]interface{}{
			"name":             raw["name"],
			"size":             raw["size"],
			"device_use":       raw["device_use"],
			"physical_volumes": buildCreateTaskTargetServerDisksDefaultPhysicalVolumesBodyParams(raw["physical_volumes"]),
		}
	}

	return params, nil
}

func getSourceServer(client *golangsdk.ServiceClient, sourceServerId string) (interface{}, error) {
	getHttpUrl := "v3/sources/{source_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{source_id}", sourceServerId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening source server: %s", err)
	}

	return getRespBody, nil
}

func buildCreateTaskTargetServerDisksDefaultPhysicalVolumesBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"name":        raw["name"],
				"size":        raw["size"],
				"device_use":  raw["device_use"],
				"file_system": raw["file_system"],
				"mount_point": raw["mount_point"],
				"index":       raw["index"],
				"uuid":        raw["uuid"],
				"used_size":   raw["used_size"],
			}
		}
		return params
	}

	return nil
}

func buildCreateTaskTargetServerDisksBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"name":             raw["name"],
				"device_use":       raw["device_type"],
				"size":             convertMBtoBytes(int64(raw["size"].(int))),
				"used_size":        convertMBtoBytes(int64(raw["used_size"].(int))),
				"disk_id":          raw["disk_id"],
				"physical_volumes": buildCreateTaskTargetServerDisksPhysicalVolumesBodyParams(raw["physical_volumes"]),
			}
		}
		return params
	}

	return nil
}

func buildCreateTaskTargetServerDisksPhysicalVolumesBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"name":        raw["name"],
				"size":        convertMBtoBytes(int64(raw["size"].(int))),
				"used_size":   convertMBtoBytes(int64(raw["used_size"].(int))),
				"device_use":  raw["device_type"],
				"file_system": raw["file_system"],
				"mount_point": raw["mount_point"],
				"index":       raw["index"],
				"uuid":        raw["uuid"],
			}
		}
		return params
	}

	return nil
}

func operationMigrateTask(client *golangsdk.ServiceClient, id, operation string) error {
	opts := tasks.ActionOpts{
		Operation: operation,
	}
	return tasks.Action(client, id, opts)
}

func resourceMigrateTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	createBodyParams, err := buildCreateTaskBodyParams(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	createHttpUrl := "v3/tasks"
	createPath := smsClient.Endpoint + createHttpUrl
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(createBodyParams),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	createResp, err := smsClient.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating SMS migrate task: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening creating migrate task response: %s", err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SMS migrate task: can not found migrate task id in return")
	}

	d.SetId(id)

	if d.Get("action").(string) == "start" {
		if err := operationMigrateTask(smsClient, id, "start"); err != nil {
			return diag.Errorf("failed to start migrate task: %s", err)
		}

		if err := waitForTaskStateRunning(ctx, smsClient, d.Timeout(schema.TimeoutCreate), id); err != nil {
			return diag.Errorf("failed to run migrate task: %s", err)
		}
	}

	if _, ok := d.GetOk("speed_limit"); ok {
		err = updateSmsTaskSpeedLimit(smsClient, id, d.Get("speed_limit").(*schema.Set).List())
		if err != nil {
			return diag.Errorf("error set task speed limit: %s", err)
		}
	}

	if _, ok := d.GetOk("configurations"); ok {
		err = updateSmsTaskConfigurationSetting(smsClient, id, d.Get("configurations").(*schema.Set).List())
		if err != nil {
			return diag.Errorf("error set task configurations: %s", err)
		}
	}

	return resourceMigrateTaskRead(ctx, d, meta)
}

func resourceMigrateTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	smsTask, err := GetSmsTask(smsClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching SMS migrate task")
	}

	log.Printf("[DEBUG] Retrieved SMS migrate task %s: %+v", d.Id(), smsTask)
	mErr := multierror.Append(nil,
		d.Set("region", utils.PathSearch("region_name", smsTask, nil)),
		d.Set("project_id", utils.PathSearch("project_id", smsTask, nil)),
		d.Set("type", utils.PathSearch("type", smsTask, nil)),
		d.Set("os_type", utils.PathSearch("os_type", smsTask, nil)),
		d.Set("vm_template_id", utils.PathSearch("vm_template_id", smsTask, nil)),
		d.Set("source_server_id", utils.PathSearch("source_server.id", smsTask, nil)),
		d.Set("target_server_id", utils.PathSearch("target_server.vm_id", smsTask, nil)),
		d.Set("target_server_name", utils.PathSearch("target_server.name", smsTask, nil)),
		d.Set("target_server_disks", flattenSmsTaskTargetServerDisks(
			utils.PathSearch("target_server.disks", smsTask, make([]interface{}, 0)).([]interface{}))),
		d.Set("start_target_server", utils.PathSearch("start_target_server", smsTask, nil)),
		d.Set("migration_ip", utils.PathSearch("migration_ip", smsTask, nil)),
		d.Set("state", utils.PathSearch("state", smsTask, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", smsTask, nil)),
		d.Set("migrate_speed", utils.PathSearch("migrate_speed", smsTask, nil)),
		d.Set("migrate_speed_limit", utils.PathSearch("speed_limit", smsTask, nil)),
		d.Set("use_public_ip", utils.PathSearch("use_public_ip", smsTask, nil)),
		d.Set("syncing", utils.PathSearch("syncing", smsTask, nil)),
		d.Set("use_ipv6", utils.PathSearch("use_ipv6", smsTask, nil)),
		d.Set("need_migration_test", utils.PathSearch("need_migration_test", smsTask, nil)),
	)

	getHttpUrl := "v3/tasks/{task_id}/speed-limit"
	if speedLimit, err := getTaskRelatedPropsByOnlyUrl(smsClient, d.Id(), getHttpUrl); err == nil {
		mErr = multierror.Append(mErr,
			d.Set("speed_limit", flattenSmsTaskSpeedLimit(
				utils.PathSearch("speed_limit", speedLimit, make([]interface{}, 0)).([]interface{}))),
		)
	} else {
		log.Printf("[WARN] error fetching task speed limit (%s): %s", d.Id(), err)
	}

	getHttpUrl = "v3/tasks/{task_id}/configuration-setting"
	if configurations, err := getTaskRelatedPropsByOnlyUrl(smsClient, d.Id(), getHttpUrl); err == nil {
		mErr = multierror.Append(mErr,
			d.Set("migrate_type", utils.PathSearch("migrate_type", configurations, nil)),
			d.Set("configurations", flattenSmsTaskConfigurations(
				utils.PathSearch("configurations", configurations, make([]interface{}, 0)).([]interface{}))),
		)
	} else {
		log.Printf("[WARN] error fetching task configurations (%s): %s", d.Id(), err)
	}

	getHttpUrl = "v3/tasks/{task_id}/passphrase"
	if passphrase, err := getTaskRelatedPropsByOnlyUrl(smsClient, d.Id(), getHttpUrl); err == nil {
		mErr = multierror.Append(mErr,
			d.Set("passphrase", utils.PathSearch("passphrase", passphrase, nil)),
		)
	} else {
		log.Printf("[WARN] error fetching task passphrase (%s): %s", d.Id(), err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SMS migrate task fields: %s", err)
	}

	return nil
}

func flattenSmsTaskTargetServerDisks(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		raw := params.(map[string]interface{})
		m := map[string]interface{}{
			"device_type": utils.PathSearch("device_use", raw, nil),
			"name":        utils.PathSearch("name", raw, nil),
			"size":        convertBytestoMB(int64(utils.PathSearch("size", raw, float64(0)).(float64))),
			"used_size":   convertBytestoMB(int64(utils.PathSearch("used_size", raw, float64(0)).(float64))),
			"disk_id":     utils.PathSearch("disk_id", raw, nil),
			"physical_volumes": flattenSmsTaskTargetServerDisksPhysicalVolumes(
				utils.PathSearch("physical_volumes", raw, nil)),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsTaskTargetServerDisksPhysicalVolumes(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"device_type": utils.PathSearch("device_use", raw, nil),
				"name":        utils.PathSearch("name", raw, nil),
				"size":        convertBytestoMB(int64(utils.PathSearch("size", raw, float64(0)).(float64))),
				"used_size":   convertBytestoMB(int64(utils.PathSearch("used_size", raw, float64(0)).(float64))),
				"file_system": utils.PathSearch("file_system", raw, nil),
				"mount_point": utils.PathSearch("mount_point", raw, nil),
				"index":       utils.PathSearch("index", raw, nil),
				"uuid":        utils.PathSearch("uuid", raw, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}

func GetSmsTask(client *golangsdk.ServiceClient, taskID string) (interface{}, error) {
	getHttpUrl := "v3/tasks/{task_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening task: %s", err)
	}

	return getRespBody, nil
}

func flattenSmsTaskConfigurations(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"config_key":    utils.PathSearch("config_key", params, nil),
			"config_value":  utils.PathSearch("config_value", params, nil),
			"config_status": utils.PathSearch("config_status", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsTaskSpeedLimit(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"start":                utils.PathSearch("start", params, nil),
			"end":                  utils.PathSearch("end", params, nil),
			"speed":                utils.PathSearch("speed", params, nil),
			"over_speed_threshold": utils.PathSearch("over_speed_threshold", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func getTaskRelatedPropsByOnlyUrl(client *golangsdk.ServiceClient, taskId string, getHttpUrl string) (interface{}, error) {
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening task related properties: %s", err)
	}

	return getRespBody, nil
}

func resourceMigrateTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	if d.HasChange("action") {
		action := d.Get("action").(string)
		if action != "" {
			if err := operationMigrateTask(smsClient, d.Id(), action); err != nil {
				return diag.Errorf("failed to %s migrate task: %s", action, err)
			}
		}
	}

	if d.HasChange("speed_limit") {
		err = updateSmsTaskSpeedLimit(smsClient, d.Id(), d.Get("speed_limit").(*schema.Set).List())
		if err != nil {
			return diag.Errorf("error updating task speed limit: %s", err)
		}
	}

	if d.HasChange("configurations") {
		err = updateSmsTaskConfigurationSetting(smsClient, d.Id(), d.Get("configurations").(*schema.Set).List())
		if err != nil {
			return diag.Errorf("error updating task configurations: %s", err)
		}
	}

	return resourceMigrateTaskRead(ctx, d, meta)
}

func updateSmsTaskConfigurationSetting(smsClient *golangsdk.ServiceClient, taskID string, configurationSettingParams []interface{}) error {
	updateHttpUrl := "v3/tasks/{task_id}/configuration-setting"
	updatePath := smsClient.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{task_id}", taskID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildUpdateTaskConfigurationSettingBodyParams(configurationSettingParams)),
	}

	_, err := smsClient.Request("POST", updatePath, &updateOpt)
	return err
}

func buildUpdateTaskConfigurationSettingBodyParams(rawParams []interface{}) map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	configurationSettingParams := make([]interface{}, len(rawParams))
	for i, v := range rawParams {
		configurationSettingParams[i] = map[string]interface{}{
			"config_key":    utils.PathSearch("config_key", v, nil),
			"config_value":  utils.PathSearch("config_value", v, nil),
			"config_status": utils.ValueIgnoreEmpty(utils.PathSearch("config_status", v, nil)),
		}
	}

	bodyParams := map[string]interface{}{
		"configurations": configurationSettingParams,
	}

	return bodyParams
}

func updateSmsTaskSpeedLimit(smsClient *golangsdk.ServiceClient, taskID string, speedLimitParams []interface{}) error {
	updateHttpUrl := "v3/tasks/{task_id}/speed-limit"
	updatePath := smsClient.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{task_id}", taskID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildUpdateTaskSpeedLimitBodyParams(speedLimitParams)),
	}

	_, err := smsClient.Request("POST", updatePath, &updateOpt)
	return err
}

func buildUpdateTaskSpeedLimitBodyParams(rawParams []interface{}) map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	speedLimitParams := make([]interface{}, len(rawParams))
	for i, v := range rawParams {
		speedLimitParams[i] = map[string]interface{}{
			"start": utils.PathSearch("start", v, nil),
			"end":   utils.PathSearch("end", v, nil),
			"speed": utils.PathSearch("speed", v, nil),
			"over_speed_threshold": utils.ValueIgnoreEmpty(
				utils.PathSearch("over_speed_threshold", v, nil)),
		}
	}

	bodyParams := map[string]interface{}{
		"speed_limit": speedLimitParams,
	}

	return bodyParams
}

func resourceMigrateTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	deleteHttpUrl := "v3/tasks/{task_id}"
	deletePath := smsClient.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{task_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = smsClient.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMS migrate task")
	}

	return nil
}

func waitForTaskStateRunning(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"READY"},
		Target:       []string{"RUNNING"},
		Refresh:      taskStateRefreshFunc(client, id),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func taskStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		migTask, err := tasks.Get(client, id)
		if err != nil {
			return nil, "ERROR", err
		}

		return migTask, migTask.State, nil
	}
}

func convertBytestoMB(bytes int64) int64 {
	return bytes / 1024 / 1024
}

func convertMBtoBytes(mb int64) int64 {
	return mb * 1024 * 1024
}
