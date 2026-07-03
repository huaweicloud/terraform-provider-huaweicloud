package deprecated

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cts/v1/tracker"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceCTSTrackerV1() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceCTSTrackerV1Read,
		DeprecationMessage: "CTS tracker has been deprecated.",

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
	cfg := meta.(*config.Config)
	trackerClient, err := cfg.CtsV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating cts Client: %s", err)
	}

	listOpts := tracker.ListOpts{
		TrackerName:    d.Get("tracker_name").(string),
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Status:         d.Get("status").(string),
	}

	refinedTrackers, err := tracker.List(trackerClient, listOpts)

	if err != nil {
		return fmt.Errorf("unable to retrieve cts tracker: %s", err)
	}

	if len(refinedTrackers) < 1 {
		return errors.New("your query returned no results, please change your search criteria and try again")
	}

	if len(refinedTrackers) > 1 {
		return errors.New("your query returned more than one result, please try a more specific search criteria")
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

	d.Set("region", cfg.GetRegion(d))

	return nil
}
