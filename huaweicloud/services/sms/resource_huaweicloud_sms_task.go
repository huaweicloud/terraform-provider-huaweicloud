package sms

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/sms/v3/sources"
	"github.com/chnsz/golangsdk/openstack/sms/v3/tasks"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceMigrateTask is the impl of huaweicloud_sms_task
func ResourceMigrateTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrateTaskCreate,
		ReadContext:   resourceMigrateTaskRead,
		UpdateContext: resourceMigrateTaskUpdate,
		DeleteContext: resourceMigrateTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
		},
	}
}

func buildDefaultTargetServerPVRequest(rawPVs []sources.PhysicalVolumes) []tasks.PVRequest {
	if len(rawPVs) == 0 {
		return nil
	}

	pvs := make([]tasks.PVRequest, len(rawPVs))
	for i, pv := range rawPVs {
		pvs[i] = tasks.PVRequest{
			Name:       pv.Name,
			Size:       pv.Size,
			DeviceType: pv.DeviceType,
			FileSystem: pv.FileSystem,
			MountPoint: pv.MountPoint,
			Index:      &pv.Index,
			UUID:       pv.UUID,
			UsedSize:   pv.UsedSize,
		}
	}

	return pvs
}

func buildDefaultTargetServerDiskRequest(d *schema.ResourceData, cfg *config.Config, sid string) ([]tasks.DiskRequest, error) {
	smsClient, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating SMS client: %s", err)
	}

	log.Printf("[DEBUG] filtering SMS source servers by id %s", sid)
	server, err := sources.Get(smsClient, sid)
	if err != nil {
		return nil, fmt.Errorf("unable to find the source server %s: %s", sid, err)
	}

	sourceDisks := server.InitTargetServer.Disks
	disks := make([]tasks.DiskRequest, len(sourceDisks))
	for i, d := range sourceDisks {
		disks[i] = tasks.DiskRequest{
			Name:            d.Name,
			Size:            d.Size,
			DeviceType:      d.DeviceUse,
			PhysicalVolumes: buildDefaultTargetServerPVRequest(d.PhysicalVolumes),
		}
	}

	return disks, nil
}

func buildTargetServerPVRequest(raw []interface{}) []tasks.PVRequest {
	if len(raw) == 0 {
		return nil
	}

	pvs := make([]tasks.PVRequest, len(raw))
	for i, pv := range raw {
		item := pv.(map[string]interface{})
		idx := item["index"].(int)
		pvs[i] = tasks.PVRequest{
			Name:       item["name"].(string),
			Size:       convertMBtoBytes(item["size"].(int64)),
			DeviceType: item["device_type"].(string),
			FileSystem: item["file_system"].(string),
			MountPoint: item["mount_point"].(string),
			Index:      &idx,
			UUID:       item["uuid"].(string),
			UsedSize:   convertMBtoBytes(item["used_size"].(int64)),
		}
	}

	return pvs
}

func buildTargetServerDiskRequest(d *schema.ResourceData, cfg *config.Config, sid string) ([]tasks.DiskRequest, error) {
	v, ok := d.GetOk("target_server_disks")
	if !ok {
		return buildDefaultTargetServerDiskRequest(d, cfg, sid)
	}

	disksRaw := v.([]interface{})
	disks := make([]tasks.DiskRequest, len(disksRaw))
	for i, d := range disksRaw {
		item := d.(map[string]interface{})
		disks[i] = tasks.DiskRequest{
			Name:            item["name"].(string),
			DeviceType:      item["device_type"].(string),
			Size:            convertMBtoBytes(item["size"].(int64)),
			UsedSize:        convertMBtoBytes(item["used_size"].(int64)),
			DiskId:          item["disk_id"].(string),
			PhysicalVolumes: buildTargetServerPVRequest(item["physical_volumes"].([]interface{})),
		}
	}

	return disks, nil
}

func buildTargetServerRequest(d *schema.ResourceData, cfg *config.Config, sid string) (tasks.TargetServerRequest, error) {
	var targetServer tasks.TargetServerRequest

	if v, ok := d.GetOk("target_server_id"); ok {
		serverID := v.(string)
		ecsClient, err := cfg.ComputeV1Client(cfg.GetRegion(d))
		if err != nil {
			return targetServer, fmt.Errorf("error creating compute client: %s", err)
		}

		server, err := cloudservers.Get(ecsClient, serverID).Extract()
		if err != nil {
			return targetServer, fmt.Errorf("error retrieving ECS instance %s: %s", serverID, err)
		}

		targetServer.Name = server.Name
		targetServer.VMID = serverID
	}

	targetDisks, err := buildTargetServerDiskRequest(d, cfg, sid)
	if err != nil {
		return tasks.TargetServerRequest{}, err
	}

	targetServer.Disks = targetDisks
	return targetServer, nil
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

func buildMigrateTaskRequest(d *schema.ResourceData, cfg *config.Config) (*tasks.CreateOpts, error) {
	region := cfg.GetRegion(d)

	sourceID := d.Get("source_server_id").(string)
	source := tasks.SourceServerRequest{
		Id: sourceID,
	}

	target, err := buildTargetServerRequest(d, cfg, sourceID)
	if err != nil {
		return nil, err
	}

	_, existing := d.GetOk("target_server_id")
	createOpts := tasks.CreateOpts{
		Name:         "MigrationTask",
		Priority:     1,
		Type:         d.Get("type").(string),
		OsType:       d.Get("os_type").(string),
		Region:       region,
		RegionID:     region,
		Project:      getProjectName(d, cfg),
		ProjectID:    getProjectID(d, cfg, region),
		SourceServer: source,
		TargetServer: target,
		VmTemplateId: d.Get("vm_template_id").(string),
		MigrationIp:  d.Get("migration_ip").(string),
		UsePublicIp:  utils.Bool(d.Get("use_public_ip").(bool)),
		StartServer:  utils.Bool(d.Get("start_target_server").(bool)),
		Syncing:      utils.Bool(d.Get("syncing").(bool)),
		ExistServer:  &existing,
	}

	return &createOpts, nil
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

	createOpts, err := buildMigrateTaskRequest(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	id, err := tasks.Create(smsClient, createOpts)
	if err != nil {
		return diag.Errorf("error creating SMS migrate task: %s", err)
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

	return resourceMigrateTaskRead(ctx, d, meta)
}

func resourceMigrateTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	migTask, err := tasks.Get(smsClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching SMS migrate task")
	}

	log.Printf("[DEBUG] Retrieved SMS migrate task %s: %+v", d.Id(), migTask)
	mErr := multierror.Append(
		d.Set("region", migTask.Region),
		d.Set("project_id", migTask.ProjectID),
		d.Set("type", migTask.Type),
		d.Set("os_type", migTask.OsType),
		d.Set("vm_template_id", migTask.VmTemplateId),
		d.Set("source_server_id", migTask.SourceServer.Id),
		d.Set("target_server_id", migTask.TargetServer.VMID),
		d.Set("target_server_name", migTask.TargetServer.Name),
		d.Set("target_server_disks", flattenTargetServerDisks(migTask.TargetServer.Disks)),
		d.Set("start_target_server", migTask.StartTargetServer),
		d.Set("migration_ip", migTask.MigrationIp),
		d.Set("state", migTask.State),
		d.Set("enterprise_project_id", migTask.EnterpriseProjectId),
		d.Set("migrate_speed", migTask.MigrateSpeed),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SMS migrate task fields: %s", err)
	}

	return nil
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

	return resourceMigrateTaskRead(ctx, d, meta)
}

func resourceMigrateTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	err = tasks.Delete(smsClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting SMS migrate task: %s", err)
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

func flattenTargetServerDisks(disks []tasks.TargetDisk) []map[string]interface{} {
	results := make([]map[string]interface{}, len(disks))
	for i, item := range disks {
		results[i] = map[string]interface{}{
			"device_type":      item.DeviceType,
			"name":             item.Name,
			"size":             convertBytestoMB(item.Size),
			"used_size":        convertBytestoMB(item.UsedSize),
			"disk_id":          item.DiskId,
			"physical_volumes": flattenTargetServerPVs(item.PhysicalVolumes),
		}
	}
	return results
}

func flattenTargetServerPVs(pvs []tasks.TargetPhysicalVolumes) []map[string]interface{} {
	results := make([]map[string]interface{}, len(pvs))
	for i, item := range pvs {
		results[i] = map[string]interface{}{
			"device_type": item.DeviceType,
			"name":        item.Name,
			"size":        convertBytestoMB(item.Size),
			"used_size":   convertBytestoMB(item.UsedSize),
			"file_system": item.FileSystem,
			"mount_point": item.MountPoint,
			"index":       item.Index,
			"uuid":        item.UUID,
		}
	}
	return results
}

func convertBytestoMB(bytes int64) int64 {
	return bytes / 1024 / 1024
}

func convertMBtoBytes(mb int64) int64 {
	return mb * 1024 * 1024
}
