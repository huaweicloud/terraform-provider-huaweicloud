package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk/openstack/smn/v2/subscriptions"
)

func resourceSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriptionCreate,
		Read:   resourceSubscriptionRead,
		Delete: resourceSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateStringList(v, k, []string{"email", "sms", "http", "https"})
				},
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}
	topicUrn := d.Get("topic_urn").(string)
	createOpts := subscriptions.CreateOps{
		Endpoint: d.Get("endpoint").(string),
		Protocol: d.Get("protocol").(string),
		Remark:   d.Get("remark").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	subscription, err := subscriptions.Create(client, createOpts, topicUrn).Extract()
	if err != nil {
		return fmt.Errorf("Error getting subscription from result: %s", err)
	}
	log.Printf("[DEBUG] Create : subscription.SubscriptionUrn %s", subscription.SubscriptionUrn)
	if subscription.SubscriptionUrn != "" {
		d.SetId(subscription.SubscriptionUrn)
		d.Set("subscription_urn", subscription.SubscriptionUrn)
		return resourceSubscriptionRead(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceSubscriptionCreate.")
}

func resourceSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}

	log.Printf("[DEBUG] Deleting subscription %s", d.Id())

	id := d.Id()
	result := subscriptions.Delete(client, id)
	if result.Err != nil {
		return err
	}

	log.Printf("[DEBUG] Successfully deleted subscription %s", id)
	return nil
}

func resourceSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}

	log.Printf("[DEBUG] Getting subscription %s", d.Id())

	id := d.Id()
	subscriptionslist, err := subscriptions.List(client).Extract()
	if err != nil {
		return fmt.Errorf("Error Get subscriptionslist: %s", err)
	}
	log.Printf("[DEBUG] list : subscriptionslist %#v", subscriptionslist)
	for _, subscription := range subscriptionslist {
		if subscription.SubscriptionUrn == id {
			log.Printf("[DEBUG] subscription: %#v", subscription)
			d.Set("topic_urn", subscription.TopicUrn)
			d.Set("endpoint", subscription.Endpoint)
			d.Set("protocol", subscription.Protocol)
			d.Set("subscription_urn", subscription.SubscriptionUrn)
			d.Set("owner", subscription.Owner)
			d.Set("remark", subscription.Remark)
			d.Set("status", subscription.Status)
		}
	}

	log.Printf("[DEBUG] Successfully get subscription %s", id)
	return nil
}
