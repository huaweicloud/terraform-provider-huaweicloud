package smn

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/smn/v2/subscriptions"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubscriptionCreate,
		ReadContext:   resourceSubscriptionRead,
		DeleteContext: resourceSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSubscriptionImport,
		},

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
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"email", "sms", "http", "https", "functionstage", "functiongraph",
				}, false),
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.SmnV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating smn client: %s", err)
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
		return diag.Errorf("error getting subscription from result: %s", err)
	}
	log.Printf("[DEBUG] Create: subscription.SubscriptionUrn %s", subscription.SubscriptionUrn)

	d.SetId(subscription.SubscriptionUrn)
	return resourceSubscriptionRead(ctx, d, meta)
}

func resourceSubscriptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.SmnV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating smn client: %s", err)
	}
	id := d.Id()
	topicUrn := d.Get("topic_urn").(string)

	log.Printf("[DEBUG] Getting subscription %s", id)

	subscriptionslist, err := subscriptions.ListFromTopic(client, topicUrn).Extract()
	if err != nil {
		return diag.Errorf("error Get subscriptionslist: %s", err)
	}
	log.Printf("[DEBUG] list: subscriptionslist %#v", subscriptionslist)

	var subscriptionToSet subscriptions.SubscriptionGet
	for _, subscription := range subscriptionslist {
		if subscription.SubscriptionUrn == id {
			subscriptionToSet = subscription
			break
		}
	}

	if (subscriptionToSet == subscriptions.SubscriptionGet{}) {
		return diag.Errorf("can't find subscription: %s", id)
	}

	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("topic_urn", subscriptionToSet.TopicUrn),
		d.Set("endpoint", subscriptionToSet.Endpoint),
		d.Set("protocol", subscriptionToSet.Protocol),
		d.Set("subscription_urn", subscriptionToSet.SubscriptionUrn),
		d.Set("owner", subscriptionToSet.Owner),
		d.Set("remark", subscriptionToSet.Remark),
		d.Set("status", subscriptionToSet.Status),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SMN topic fields: %s", err)
	}

	log.Printf("[DEBUG] Successfully get subscription %s", id)
	return nil
}

func resourceSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.SmnV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating smn client: %s", err)
	}

	log.Printf("[DEBUG] Deleting subscription %s", d.Id())

	id := d.Id()
	result := subscriptions.Delete(client, id)
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}

	log.Printf("[DEBUG] Successfully deleted subscription %s", id)
	return nil
}

func resourceSubscriptionImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")

	topicUrn := strings.Join(parts[:len(parts)-1], ":")

	d.SetId(d.Id())
	d.Set("topic_urn", topicUrn)

	return []*schema.ResourceData{d}, nil
}

// urn:smn:cn-north-4:0970dd7a1300f5672ff2c003c60ae115:AUTO_ALARM_NOTIFY_TOPIC_MYSQL_mysql_0970dd7a1300f5672ff2c003c60ae115
// urn:smn:cn-north-4:0970dd7a1300f5672ff2c003c60ae115:AUTO_ALARM_NOTIFY_TOPIC_MYSQL_mysql_0970dd7a1300f5672ff2c003c60ae115:a3ee64b608334ef2a47deef3b8454b14
// urn:smn:cn-north-4:0970dd7a1300f5672ff2c003c60ae115:topic_2
// urn:smn:cn-north-4:0970dd7a1300f5672ff2c003c60ae115:topic_2:a2aa5a1f66df494184f4e108398de1a6
// urn:smn:cn-north-4:0970dd7a1300f5672ff2c003c60ae115:topic_2:a2aa5a1f66df494184f4e108398de1a6
