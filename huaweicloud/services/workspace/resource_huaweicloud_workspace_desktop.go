package workspace

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/desktops"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/jobs"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func desktopVolumeSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// @API Workspace POST /v2/{project_id}/desktops/rebuild
// @API Workspace POST /v2/{project_id}/desktops/resize
// @API Workspace GET /v2/{project_id}/desktops/{desktop_id}
// @API Workspace DELETE /v2/{project_id}/desktops/{desktop_id}
// @API Workspace POST /v2/{project_id}/desktops/{id}/tags/action
// @API Workspace POST /v2/{project_id}/volumes
// @API Workspace GET /v2/{project_id}/desktops/{desktop_id}/networks
// @API Workspace PUT /v2/{project_id}/desktops/{desktop_id}/networks
// @API Workspace POST /v2/{project_id}/desktops
// @API Workspace POST /v2/{project_id}/volumes/expand
// @API Workspace GET /v2/{project_id}/workspace-sub-jobs
func ResourceDesktop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopCreate,
		ReadContext:   resourceDesktopRead,
		UpdateContext: resourceDesktopUpdate,
		DeleteContext: resourceDesktopDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

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
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"market", "gold", "private",
				}, false),
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sudo", "default", "administrators", "users",
				}, false),
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     desktopVolumeSchemaResource(),
			},
			"data_volume": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     desktopVolumeSchemaResource(),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nic": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"email_notification": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_user": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"power_action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"power_action_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDesktopRootVolume(volumes []interface{}) *desktops.Volume {
	if len(volumes) < 1 {
		return nil
	}

	volume := volumes[0].(map[string]interface{})
	result := desktops.Volume{
		Type: volume["type"].(string),
		Size: volume["size"].(int),
	}
	return &result
}

func buildDesktopDataVolumes(volumes []interface{}) []desktops.Volume {
	if len(volumes) < 1 {
		return nil
	}

	result := make([]desktops.Volume, len(volumes))
	for i, val := range volumes {
		volume := val.(map[string]interface{})
		result[i] = desktops.Volume{
			Type: volume["type"].(string),
			Size: volume["size"].(int),
		}
	}
	return result
}

func buildDesktopNics(nics []interface{}) []desktops.Nic {
	if len(nics) < 1 {
		return nil
	}

	result := make([]desktops.Nic, len(nics))
	for i, val := range nics {
		volume := val.(map[string]interface{})
		result[i] = desktops.Nic{
			NetworkId: volume["network_id"].(string),
		}
	}
	return result
}

func buildDesktopSecurityGroups(securityGroups *schema.Set) []desktops.SecurityGroup {
	if securityGroups.Len() < 1 {
		return nil
	}

	result := make([]desktops.SecurityGroup, securityGroups.Len())
	for i, val := range securityGroups.List() {
		result[i] = desktops.SecurityGroup{
			ID: val.(string),
		}
	}
	return result
}

func buildDesktopCreateOpts(d *schema.ResourceData, conf *config.Config) desktops.CreateOpts {
	result := desktops.CreateOpts{
		Desktops: []desktops.DesktopConfig{
			{
				UserName:    d.Get("user_name").(string),
				UserEmail:   d.Get("user_email").(string),
				UserGroup:   d.Get("user_group").(string),
				DesktopName: d.Get("name").(string),
			},
		},
		DesktopType:         "DEDICATED",
		ProductId:           d.Get("flavor_id").(string),
		RootVolume:          buildDesktopRootVolume(d.Get("root_volume").([]interface{})),
		AvailabilityZone:    d.Get("availability_zone").(string),
		ImageType:           d.Get("image_type").(string),
		ImageId:             d.Get("image_id").(string),
		VpcId:               d.Get("vpc_id").(string),
		EmailNotification:   utils.Bool(d.Get("email_notification").(bool)),
		DataVolumes:         buildDesktopDataVolumes(d.Get("data_volume").([]interface{})),
		Nics:                buildDesktopNics(d.Get("nic").([]interface{})),
		SecurityGroups:      buildDesktopSecurityGroups(d.Get("security_groups").(*schema.Set)),
		Tags:                utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		EnterpriseProjectId: conf.GetEnterpriseProjectID(d),
	}
	return result
}

func waitForWorkspaceJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string, timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"WAITING", "RUNNING"},
		Target:       []string{"SUCCESS"},
		Refresh:      refreshWorkspaceJobFunc(client, jobId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", err
	}

	return resp.(jobs.Job).Entities.DesktopId, nil
}

func refreshWorkspaceJobFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := jobs.ListOpts{
			JobId: jobId,
		}
		resp, err := jobs.List(client, opts)
		if err != nil {
			return resp, "", err
		}
		if resp.TotalCount < 1 {
			return resp, "", fmt.Errorf("unable to find any job details")
		}

		for _, job := range resp.Jobs {
			if job.Status == "SUCCESS" {
				continue
			}
			return job, job.Status, nil
		}

		return resp.Jobs[0], "SUCCESS", nil
	}
}

func resourceDesktopCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	createOpts := buildDesktopCreateOpts(d, conf)
	resp, err := desktops.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating Workspace desktop: %s", err)
	}
	desktopId, err := waitForWorkspaceJobCompleted(ctx, client, resp.JobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
	}
	log.Printf("[DEBUG] The job (%s) has been completed", resp.JobId)

	d.SetId(desktopId)

	if action, ok := d.GetOk("power_action"); ok {
		if action == "os-start" {
			log.Printf("[WARN] the power action (os-start) is invalid after desktop created")
		} else if err = updateDesktopPowerAction(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceDesktopRead(ctx, d, meta)
}

func flattenDesktopRootVolume(volume desktops.VolumeResp) []map[string]interface{} {
	if volume == (desktops.VolumeResp{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":       volume.Type,
			"size":       volume.Size,
			"id":         volume.VolumeId,
			"name":       volume.Name,
			"device":     volume.Device,
			"created_at": volume.CreatedAt,
		},
	}
}

func flattenDesktopDataVolumes(volumes []desktops.VolumeResp) []map[string]interface{} {
	if len(volumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(volumes))
	for i, volume := range volumes {
		result[i] = map[string]interface{}{
			"type":       volume.Type,
			"size":       volume.Size,
			"id":         volume.VolumeId,
			"name":       volume.Name,
			"device":     volume.Device,
			"created_at": volume.CreatedAt,
		}
	}

	// Since the volumes in the response body are unordered, they are sorted by device.
	sort.Slice(result, func(i, j int) bool {
		a := result[i]
		b := result[j]

		return a["device"].(string) <= b["device"].(string)
	})

	return result
}

func flattenDesktopSecurityGroups(securityGroups []desktops.SecurityGroup) []interface{} {
	if len(securityGroups) < 1 {
		return nil
	}

	result := make([]interface{}, len(securityGroups))
	for i, securityGroup := range securityGroups {
		result[i] = securityGroup.ID
	}
	return result
}

func getDesktopNetwork(client *golangsdk.ServiceClient, desktopId string) ([]map[string]interface{}, error) {
	network, err := desktops.GetNetwork(client, desktopId)
	if err != nil {
		return nil, fmt.Errorf("error getting desktop network info: %s", err)
	}

	if len(network) < 1 {
		return nil, fmt.Errorf("unable to find any network information under Workspace desktop (%s)", desktopId)
	}

	nic := []map[string]interface{}{
		{
			"network_id": network[0].Subnet.ID,
		},
	}

	return nic, err
}

func resourceDesktopRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	desktopId := d.Id()
	resp, err := desktops.Get(client, desktopId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace desktop")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("flavor_id", resp.Product.ID),
		d.Set("user_name", resp.UserName),
		d.Set("root_volume", flattenDesktopRootVolume(resp.RootVolume)),
		d.Set("data_volume", flattenDesktopDataVolumes(resp.DataVolumes)),
		d.Set("availability_zone", resp.AvailabilityZone),
		d.Set("user_group", resp.UserGroup),
		d.Set("name", resp.Name),
		d.Set("tags", utils.TagsToMap(resp.Tags)),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("status", resp.Status),
	)

	if imageId, ok := resp.Metadata["metering.image_id"]; ok {
		mErr = multierror.Append(mErr, d.Set("image_id", imageId))
	} else {
		mErr = multierror.Append(mErr, fmt.Errorf("the image_id field does not found in metadata structure"))
	}

	securityGroups := resp.SecurityGroups
	if len(securityGroups) < 1 {
		mErr = multierror.Append(mErr, fmt.Errorf("the security_groups field does not found in API response"))
	} else {
		mErr = multierror.Append(mErr, d.Set("security_groups", flattenDesktopSecurityGroups(securityGroups)))
	}

	nicVal, err := getDesktopNetwork(client, desktopId)
	if err != nil {
		// This feature is not available in some region, so use log.Printf to record the error.
		log.Printf("[ERROR] %s", err)
	} else {
		mErr = multierror.Append(mErr, d.Set("nic", nicVal))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting desktop fields: %s", err)
	}
	return nil
}

func updateDesktopFlavor(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	opts := desktops.ProductUpdateOpts{
		Desktops: []desktops.DesktopUpdateConfig{
			{
				DesktopId: d.Id(),
			},
		},
		ProductId: d.Get("flavor_id").(string),
		Mode:      "STOP_DESKTOP",
	}
	resp, err := desktops.UpdateProduct(client, opts)
	if err != nil {
		return fmt.Errorf("error updating desktop product: %s", err)
	}

	for _, job := range resp {
		_, err = waitForWorkspaceJobCompleted(ctx, client, job.ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for the job (%s) completed: %s", job.ID, err)
		}
		log.Printf("[DEBUG] The job (%s) has been completed", job.ID)
	}
	log.Printf("[DEBUG] All jobs has been completed")
	return nil
}

func updateDesktopVolumes(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	desktopId := d.Id()
	expandSlice := make([]desktops.ExpandVolumeConfig, 0)

	if d.HasChange("root_volume") {
		expandSlice = append(expandSlice, desktops.ExpandVolumeConfig{
			DesktopId: desktopId,
			VolumeId:  d.Get("root_volume.0.id").(string),
			NewSize:   d.Get("root_volume.0.size").(int),
		})
	}

	lengthDiff := 0
	if d.HasChange("data_volume") {
		oldVal, newVal := d.GetChange("data_volume")
		oldRaw := oldVal.([]interface{})
		newRaw := newVal.([]interface{})
		newLen := len(newRaw)
		oldLen := len(oldRaw)
		if newLen < oldLen {
			return fmt.Errorf("the number of volumes cannot be reduced")
		}
		lengthDiff = newLen - oldLen

		for i, val := range oldRaw {
			oldVolume := val.(map[string]interface{})
			newVolume := newRaw[i].(map[string]interface{})
			if newVolume["type"].(string) != oldVolume["type"].(string) {
				return fmt.Errorf("volume type does not support updates")
			}
			if newVolume["size"].(int) < oldVolume["size"].(int) {
				return fmt.Errorf("volume (%v) size (old:%v, new:%v) cannot be smaller than the size before the change",
					oldVolume["name"], oldVolume["size"], newVolume["size"])
			} else if newVolume["size"].(int) > oldVolume["size"].(int) {
				expandSlice = append(expandSlice, desktops.ExpandVolumeConfig{
					DesktopId: desktopId,
					VolumeId:  oldVolume["id"].(string),
					NewSize:   newVolume["size"].(int),
				})
			}
		}

		if lengthDiff > 0 {
			newVolumeSlice := make([]desktops.Volume, 0, lengthDiff)
			for i := newLen - lengthDiff; i < newLen; i++ {
				newVolume := newRaw[i].(map[string]interface{})
				newVolumeSlice = append(newVolumeSlice, desktops.Volume{
					Type: newVolume["type"].(string),
					Size: newVolume["size"].(int),
				})
			}
			newVolumeOpts := desktops.NewVolumeOpts{
				VolumeConfigs: []desktops.NewVolumeConfig{
					{
						DesktopId: desktopId,
						Volumes:   newVolumeSlice,
					},
				},
			}
			log.Printf("[DEBUG] The new volumeOpts is: %#v", newVolumeOpts)
			resp, err := desktops.NewVolumes(client, newVolumeOpts)
			if err != nil {
				return fmt.Errorf("failed to add volume: %s", err)
			}
			_, err = waitForWorkspaceJobCompleted(ctx, client, resp.JobId, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return fmt.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
			}
			log.Printf("[DEBUG] The job (%s) has been completed", resp.JobId)
		}
	}

	if len(expandSlice) > 0 {
		expandOpts := desktops.VolumeExpandOpts{
			VolumeConfigs: expandSlice,
		}
		log.Printf("[DEBUG] The new expandOpts is: %#v", expandOpts)
		resp, err := desktops.ExpandVolumes(client, expandOpts)
		if err != nil {
			return fmt.Errorf("failed to expand volume size: %s", err)
		}
		_, err = waitForWorkspaceJobCompleted(ctx, client, resp.JobId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
		}
		log.Printf("[DEBUG] The job (%s) has been completed", resp.JobId)
	}
	return nil
}

func updateDesktopNetwork(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	nicRaw := d.Get("nic").([]interface{})
	if len(nicRaw) < 1 {
		return nil
	}

	nicVal := nicRaw[0].(map[string]interface{})
	securityGroups := d.Get("security_groups").(*schema.Set)
	listId := make([]string, securityGroups.Len())
	for i, id := range securityGroups.List() {
		listId[i] = id.(string)
	}

	desktopId := d.Id()
	opts := desktops.UpdateNetworkOpts{
		DesktopId:        desktopId,
		VpcId:            d.Get("vpc_id").(string),
		SubnetId:         nicVal["network_id"].(string),
		SecurityGroupIds: listId,
	}

	resp, err := desktops.UpdateNetwork(client, opts)
	if err != nil {
		return fmt.Errorf("error updating the network of the Workspace desktop (%s): %s", desktopId, err)
	}
	_, err = waitForWorkspaceJobCompleted(ctx, client, resp.JobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
	}
	log.Printf("[DEBUG] The job (%s) has been completed", resp.JobId)
	return nil
}

func waitForWorkspaceStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      desktopStatusRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
		// The final status of both startup and restart is "ACTIVE". the parameter only applies to "REBOOT" action.
		ContinuousTargetOccurence: 2,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func desktopStatusRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		desktopId := d.Id()
		resp, err := desktops.Get(client, desktopId)
		if err != nil {
			return resp, "", err
		}

		// In statusMap, key represents the power action of the desktop, and value represents the status after the desktop operation is completed.
		// Use statusMap to make a mapping relationship between the power action of the desktop and the final status of the desktop.
		statusMap := map[string]string{
			"os-start":     "ACTIVE",
			"os-stop":      "SHUTOFF",
			"reboot":       "ACTIVE",
			"os-hibernate": "SHUTOFF",
		}

		powerAction := d.Get("power_action").(string)
		// TaskStatus variable is always an empty string when the desktop power action is completed.
		// If the desktop power action changes from one state to another, taskStatus is an empty string for a long time,
		// whether a desktop action is completed cannot be determined only by taskStatus.
		if resp.TaskStatus == "" && resp.Status == statusMap[powerAction] {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func updateDesktopPowerAction(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	desktopId := d.Id()
	action := d.Get("power_action").(string)
	opts := desktops.ActionOpts{
		DesktopIds: []string{desktopId},
		OpType:     action,
		Type:       d.Get("power_action_type").(string),
	}

	resp, err := desktops.DoAction(client, opts)
	if err != nil {
		return fmt.Errorf("error updating the power action of the Workspace desktop (%s): %s", desktopId, err)
	}

	if resp.JobId == "" {
		err = waitForWorkspaceStatusCompleted(ctx, client, d)
		if err != nil {
			return fmt.Errorf("error waiting for power action (%s) for desktop (%s) failed: %s", action, desktopId, err)
		}
		return nil
	}

	// Job ID is returned when cold migration is started.
	_, err = waitForWorkspaceJobCompleted(ctx, client, resp.JobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
	}
	log.Printf("[DEBUG] The job (%s) has been completed", resp.JobId)
	return nil
}

func resourceDesktopUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	desktopId := d.Id()
	if d.HasChange("flavor_id") {
		if err = updateDesktopFlavor(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("image_type", "image_id") {
		rebuildOpts := desktops.RebuildOpts{
			DesktopIds: []string{desktopId},
			ImageType:  d.Get("image_type").(string),
			ImageId:    d.Get("image_id").(string),
		}
		resp, err := desktops.Rebuild(client, rebuildOpts)
		if err != nil {
			return diag.Errorf("error rebuild Workspace desktop: %s", err)
		}
		_, err = waitForWorkspaceJobCompleted(ctx, client, resp.JobId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
		}
		log.Printf("[DEBUG] The job (%s) has been completed", resp.JobId)
	}

	if d.HasChanges("root_volume", "data_volume") {
		if err = updateDesktopVolumes(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "desktops", desktopId)
		if err != nil {
			return diag.Errorf("error updating tags of Workspace desktop (%s): %s", desktopId, err)
		}
	}

	if d.HasChange("nic") {
		err = updateDesktopNetwork(ctx, client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   desktopId,
			ResourceType: "workspace-desktop",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("power_action") {
		err = updateDesktopPowerAction(ctx, client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDesktopRead(ctx, d, meta)
}

func waitForDesktopDeleted(ctx context.Context, client *golangsdk.ServiceClient, desktopId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE", "DELETING", "SHUTOFF"},
		Target:       []string{"DELETED"},
		Refresh:      refreshDesktopStatusFunc(client, desktopId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshDesktopStatusFunc(client *golangsdk.ServiceClient, desktopId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := desktops.Get(client, desktopId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return resp, "DELETED", nil
			}
			return resp, "ERROR", err
		}
		// During the removal process of desktop, the workspace service cannot perceive the ECS mechine and the API will
		// return an empty status.
		if resp.Status == "" {
			return resp, "DELETING", nil
		}
		// The uppercase characters is the default format for attribute 'status' in the API response.
		return resp, strings.ToUpper(resp.Status), nil
	}
}

func waitForDesktopUserDeleted(ctx context.Context, client *golangsdk.ServiceClient, userName string, timeout time.Duration) error {
	listOpts := users.ListOpts{
		Name: userName,
	}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"ACTIVE"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := users.List(client, listOpts)
			if err != nil {
				return resp, "ERROR", err
			}
			if len(resp.Users) < 1 {
				return resp, "DELETED", nil
			}
			return resp, "ACTIVE", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceDesktopDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	isDeleteUser := d.Get("delete_user").(bool)
	opts := desktops.DeleteOpts{
		DeleteUser:        isDeleteUser,
		EmailNotification: d.Get("email_notification").(bool),
	}
	err = desktops.Delete(client, d.Id(), opts)
	if err != nil {
		return diag.Errorf("error deleting desktop (%s): %s", d.Id(), err)
	}
	// Make sure the desktop has been deleted.
	err = waitForDesktopDeleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("an error occur when delete desktop: %s", err)
	}
	if isDeleteUser {
		// Make sure the related user has been deleted.
		userName := d.Get("user_name").(string)
		err = waitForDesktopUserDeleted(ctx, client, userName, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.Errorf("an error occur when delete user: %s", err)
		}
	}
	return nil
}
