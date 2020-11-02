package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cts/v1/tracker"
)

func dataSourceCTSTrackerV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCTSTrackerV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"file_prefix_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tracker_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_support_smn": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"topic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"is_send_all_key_operation": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"need_notify_user_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceCTSTrackerV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	trackerClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	listOpts := tracker.ListOpts{
		TrackerName:    d.Get("tracker_name").(string),
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Status:         d.Get("status").(string),
	}

	refinedTrackers, err := tracker.List(trackerClient, listOpts)

	if err != nil {
		return fmt.Errorf("Unable to retrieve cts tracker: %s", err)
	}

	if len(refinedTrackers) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedTrackers) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	trackers := refinedTrackers[0]

	log.Printf("[INFO] Retrieved cts tracker %s using given filter", trackers.TrackerName)

	d.SetId(trackers.TrackerName)

	d.Set("tracker_name", trackers.TrackerName)
	d.Set("bucket_name", trackers.BucketName)
	d.Set("file_prefix_name", trackers.FilePrefixName)
	d.Set("status", trackers.Status)
	d.Set("is_support_smn", trackers.SimpleMessageNotification.IsSupportSMN)
	d.Set("topic_id", trackers.SimpleMessageNotification.TopicID)
	d.Set("is_send_all_key_operation", trackers.SimpleMessageNotification.IsSendAllKeyOperation)
	d.Set("operations", trackers.SimpleMessageNotification.Operations)
	d.Set("need_notify_user_list", trackers.SimpleMessageNotification.NeedNotifyUserList)

	d.Set("region", GetRegion(d, config))

	return nil
}
