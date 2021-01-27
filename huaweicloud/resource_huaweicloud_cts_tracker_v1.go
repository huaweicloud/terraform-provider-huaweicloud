package huaweicloud

import (
	"time"

	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cts/v1/tracker"
)

func resourceCTSTrackerV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCTSTrackerCreate,
		Read:   resourceCTSTrackerRead,
		Update: resourceCTSTrackerUpdate,
		Delete: resourceCTSTrackerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tracker_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_prefix_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateName,
			},
			"is_support_smn": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"topic_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operations": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"is_send_all_key_operation": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"need_notify_user_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}

}

func resourceCTSTrackerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	if d.Get("is_support_smn").(bool) == true && d.Get("topic_id").(string) == "" {
		return fmt.Errorf("Error 'topic_id' is required if 'is_support_smn' is set true")
	}

	createOpts := tracker.CreateOptsWithSMN{
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		SimpleMessageNotification: tracker.SimpleMessageNotification{
			IsSupportSMN:          d.Get("is_support_smn").(bool),
			TopicID:               d.Get("topic_id").(string),
			Operations:            resourceCTSOperations(d),
			IsSendAllKeyOperation: d.Get("is_send_all_key_operation").(bool),
			NeedNotifyUserList:    resourceCTSNeedNotifyUserList(d),
		},
	}

	trackers, err := tracker.Create(ctsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating CTS tracker : %s", err)
	}

	d.SetId(trackers.TrackerName)
	//lintignore:R018
	time.Sleep(20 * time.Second)
	return resourceCTSTrackerRead(d, meta)
}

func resourceCTSTrackerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}
	listOpts := tracker.ListOpts{
		TrackerName:    d.Get("tracker_name").(string),
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Status:         d.Get("status").(string),
	}
	trackers, err := tracker.List(ctsClient, listOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Removing cts tracker %s as it's already gone", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving cts tracker: %s", err)
	}

	ctsTracker := trackers[0]

	d.Set("tracker_name", ctsTracker.TrackerName)
	d.Set("bucket_name", ctsTracker.BucketName)
	d.Set("status", ctsTracker.Status)
	d.Set("file_prefix_name", ctsTracker.FilePrefixName)
	d.Set("is_support_smn", ctsTracker.SimpleMessageNotification.IsSupportSMN)
	d.Set("topic_id", ctsTracker.SimpleMessageNotification.TopicID)
	d.Set("is_send_all_key_operation", ctsTracker.SimpleMessageNotification.IsSendAllKeyOperation)
	d.Set("operations", ctsTracker.SimpleMessageNotification.Operations)
	d.Set("need_notify_user_list", ctsTracker.SimpleMessageNotification.NeedNotifyUserList)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCTSTrackerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}
	var updateOpts tracker.UpdateOptsWithSMN

	//as bucket_name is mandatory while updating tracker
	updateOpts.BucketName = d.Get("bucket_name").(string)

	updateOpts.SimpleMessageNotification.TopicID = d.Get("topic_id").(string)

	updateOpts.SimpleMessageNotification.Operations = resourceCTSOperations(d)

	updateOpts.SimpleMessageNotification.NeedNotifyUserList = resourceCTSNeedNotifyUserList(d)

	updateOpts.SimpleMessageNotification.IsSupportSMN = d.Get("is_support_smn").(bool)

	if d.HasChange("file_prefix_name") {
		updateOpts.FilePrefixName = d.Get("file_prefix_name").(string)
	}
	if d.HasChange("status") {
		updateOpts.Status = d.Get("status").(string)
	}
	if d.HasChange("is_send_all_key_operation") {
		updateOpts.SimpleMessageNotification.IsSendAllKeyOperation = d.Get("is_send_all_key_operation").(bool)
	}

	_, err = tracker.Update(ctsClient, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating cts tracker: %s", err)
	}
	//lintignore:R018
	time.Sleep(20 * time.Second)
	return resourceCTSTrackerRead(d, meta)
}

func resourceCTSTrackerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	result := tracker.Delete(ctsClient)
	if result.Err != nil {
		return err
	}
	//lintignore:R018
	time.Sleep(20 * time.Second)
	log.Printf("[DEBUG] Successfully deleted cts tracker %s", d.Id())

	return nil
}

func resourceCTSOperations(d *schema.ResourceData) []string {
	rawOperations := d.Get("operations").(*schema.Set)
	operation := make([]string, (rawOperations).Len())
	for i, raw := range rawOperations.List() {
		operation[i] = raw.(string)
	}
	return operation
}

func resourceCTSNeedNotifyUserList(d *schema.ResourceData) []string {
	rawNotify := d.Get("need_notify_user_list").(*schema.Set)
	notify := make([]string, (rawNotify).Len())
	for i, raw := range rawNotify.List() {
		notify[i] = raw.(string)
	}
	return notify
}
