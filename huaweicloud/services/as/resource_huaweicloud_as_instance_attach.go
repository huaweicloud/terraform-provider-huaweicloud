package as

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AS POST /autoscaling-api/v1/{project_id}/scaling_group_instance/{scaling_group_id}/action
// @API AS GET /autoscaling-api/v1/{project_id}/scaling_group_instance/{scaling_group_id}/list
func ResourceASInstanceAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceAttachCreate,
		ReadContext:   resourceInstanceAttachRead,
		UpdateContext: resourceInstanceAttachUpdate,
		DeleteContext: resourceInstanceAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"standby": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"append_instance": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"standby"},
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceInstanceAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	asClient, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	groupID := d.Get("scaling_group_id").(string)
	instanceID := d.Get("instance_id").(string)

	createActions := []instances.BatchOpts{
		{
			Action:    "ADD",
			Instances: []string{instanceID},
		},
	}
	if err := doBatchAction(ctx, asClient, d.Timeout(schema.TimeoutCreate), groupID, createActions); err != nil {
		return diag.Errorf("error attaching instance %s to AS group %s: %s", instanceID, groupID, err)
	}

	resourceID := fmt.Sprintf("%s/%s", groupID, instanceID)
	d.SetId(resourceID)

	// waiting for the ECS instance is INSERVICE after attaching
	if err := waitASGroupInstancesInService(ctx, asClient, groupID, instanceID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("the instance %s is not INSERVICE: %s", instanceID, err)
	}

	actionList := make([]instances.BatchOpts, 0)
	if d.Get("protected").(bool) {
		actionList = append(actionList, instances.BatchOpts{
			Action:    "PROTECT",
			Instances: []string{instanceID},
		})
	}

	if d.Get("standby").(bool) {
		isAppend := "no"
		if d.Get("append_instance").(bool) {
			isAppend = "yes"
		}

		actionList = append(actionList, instances.BatchOpts{
			Action:    "ENTER_STANDBY",
			Instances: []string{instanceID},
			AppendEcs: isAppend,
		})
	}

	if err := doBatchAction(ctx, asClient, d.Timeout(schema.TimeoutCreate), groupID, actionList); err != nil {
		return diag.Errorf("error updating instance %s in AS group %s: %s", instanceID, groupID, err)
	}

	if len(actionList) > 0 {
		// waiting for the ECS instance is INSERVICE after updating
		if err := waitASGroupInstancesInService(ctx, asClient, groupID, instanceID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("the instance %s is not INSERVICE: %s", instanceID, err)
		}
	}

	return resourceInstanceAttachRead(ctx, d, meta)
}

func resourceInstanceAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	asClient, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <scaling_group_id>/<instance_id>")
	}

	groupID := parts[0]
	instanceID := parts[1]
	ins, err := getGroupInstanceByID(asClient, groupID, instanceID)
	if err != nil {
		// When the group does not exist or the instance is not in the group, the method `getGroupInstanceByID` will
		// specially handle these two scenarios into a 404 error code.
		return common.CheckDeletedDiag(d, err, "error retrieving AS instance attach")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("scaling_group_id", ins.GroupID),
		d.Set("instance_id", ins.ID),
		d.Set("instance_name", ins.Name),
		d.Set("health_status", ins.HealthStatus),
		d.Set("status", ins.LifeCycleStatus),
		d.Set("protected", ins.Protected),
		d.Set("standby", ins.LifeCycleStatus == "STANDBY"),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceAttachUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	asClient, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	groupID := d.Get("scaling_group_id").(string)
	instanceID := d.Get("instance_id").(string)
	actionList := make([]instances.BatchOpts, 0)

	if d.HasChange("protected") {
		action := "UNPROTECT"
		if d.Get("protected").(bool) {
			action = "PROTECT"
		}
		actionList = append(actionList, instances.BatchOpts{
			Action:    action,
			Instances: []string{instanceID},
		})
	}

	if d.HasChange("standby") {
		var isAppend string
		action := "EXIT_STANDBY"
		if d.Get("standby").(bool) {
			action = "ENTER_STANDBY"
			isAppend = "no"
			if d.Get("append_instance").(bool) {
				isAppend = "yes"
			}
		}
		actionList = append(actionList, instances.BatchOpts{
			Action:    action,
			Instances: []string{instanceID},
			AppendEcs: isAppend,
		})
	}

	if err := doBatchAction(ctx, asClient, d.Timeout(schema.TimeoutUpdate), groupID, actionList); err != nil {
		return diag.Errorf("error updating instance %s in AS group %s: %s", instanceID, groupID, err)
	}

	// waiting for the ECS instance is INSERVICE
	if err := waitASGroupInstancesInService(ctx, asClient, groupID, instanceID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("the instance %s is not INSERVICE: %s", instanceID, err)
	}

	return resourceInstanceAttachRead(ctx, d, meta)
}

func resourceInstanceAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	asClient, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	groupID := d.Get("scaling_group_id").(string)
	instanceID := d.Get("instance_id").(string)
	createActions := []instances.BatchOpts{
		{
			Action:    "REMOVE",
			Instances: []string{instanceID},
		},
	}
	if err := doBatchAction(ctx, asClient, d.Timeout(schema.TimeoutDelete), groupID, createActions); err != nil {
		// When removing a non-existing instance from the scaling group, the error_code is "AS.4030", which means that
		// the batch deletion of cloud servers failed.
		// This error message is ambiguous and cannot be regarded as successfully deleted, so the checkDeleted
		// verification is not performed.
		return diag.Errorf("error disattaching instance %s from AS group %s: %s", instanceID, groupID, err)
	}

	// waiting for the ECS instance deleted
	if err := waitASGroupInstanceDeleted(ctx, asClient, groupID, instanceID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error disattaching instance %s from AS group %s: %s", instanceID, groupID, err)
	}
	return nil
}

func doBatchAction(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	groupID string, actions []instances.BatchOpts) error {
	for _, opt := range actions {
		instanceID := opt.Instances[0]
		log.Printf("[DEBUG] try to %s instance %s from AS group %s", strings.ToLower(opt.Action), instanceID, groupID)
		log.Printf("[DEBUG] the action options: %#v", opt)

		err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			if err := instances.BatchAction(client, groupID, opt).ExtractErr(); err != nil {
				if isRetryableError(err) {
					// waiting for the AS group is INSERVICE and try again
					if waitErr := waitASGroupInstancesInService(ctx, client, groupID, "", timeout); waitErr != nil {
						return resource.NonRetryableError(fmt.Errorf("the AS group %s is not INSERVICE: %s", groupID, waitErr))
					}
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func getGroupInstanceByID(client *golangsdk.ServiceClient, groupID, instanceID string) (*instances.Instance, error) {
	page, err := instances.List(client, groupID, nil).AllPages()
	if err != nil {
		// When the group does not exist, the query instance list API response reports the following error:
		// {"error":{"code":"AS.2007","message":"The AS group does not exist."}.
		// It needs to be specially processed into 404
		return nil, common.ConvertExpected400ErrInto404Err(err, "error.code", "AS.2007")
	}

	allInstances, err := page.(instances.InstancePage).Extract()
	if err != nil {
		return nil, err
	}

	for _, ins := range allInstances {
		if ins.ID == instanceID {
			return &ins, nil
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/autoscaling-api/v1/{project_id}/scaling_group_instance/{scaling_group_id}/list",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the instance (%s) does not exist in AS group (%s)", instanceID, groupID)),
		},
	}
}

func waitASGroupInstancesInService(ctx context.Context, client *golangsdk.ServiceClient, groupID, instanceID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"INSERVICE"},
		Refresh:      refreshInstancesStatus(client, groupID, instanceID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitASGroupInstanceDeleted(ctx context.Context, client *golangsdk.ServiceClient, groupID, instanceID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"REMOVING"},
		Target:       []string{"DELETED"},
		Refresh:      checkInstanceDeleted(client, groupID, instanceID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshInstancesStatus(asClient *golangsdk.ServiceClient, groupID, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		allIns, err := getAllInstancesInGroup(asClient, groupID)
		if err != nil {
			return nil, "ERROR", err
		}

		for _, ins := range allIns {
			if instanceID != "" && utils.PathSearch("instance_id", ins, "").(string) != instanceID {
				continue
			}

			// the status may be PENDING, PENDING_WAIT, REMOVING, REMOVING_WAIT, ENTERING_STANDBY
			if strings.Contains(utils.PathSearch("life_cycle_state", ins, "").(string), "ING") {
				return allIns, "PENDING", nil
			}
		}
		// the status may be INSERVICE, STANDBY
		return allIns, "INSERVICE", nil
	}
}

func checkInstanceDeleted(asClient *golangsdk.ServiceClient, groupID, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		allIns, err := getAllInstancesInGroup(asClient, groupID)
		if err != nil {
			return nil, "ERROR", err
		}

		for _, ins := range allIns {
			if utils.PathSearch("instance_id", ins, "").(string) == instanceID {
				return allIns, "REMOVING", nil
			}
		}

		return allIns, "DELETED", nil
	}
}

func isRetryableError(err error) bool {
	if apiErr, ok := err.(golangsdk.ErrDefault400); ok {
		var respBody interface{}
		if jsonErr := json.Unmarshal(apiErr.Body, &respBody); jsonErr != nil {
			log.Printf("[WARN] failed to unmarshal the response body: %s", jsonErr)
			return false
		}
		// AS.2033: "You are not allowed to perform the operation when the AS group is in current [xxx] status."
		// AS.0003: "AS group lock conflict."
		errCode := utils.PathSearch("error.code", respBody, "")
		if errCode == "AS.2033" || errCode == "AS.0003" {
			return true
		}
	}
	return false
}
