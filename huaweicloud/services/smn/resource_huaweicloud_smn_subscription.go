package smn

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/smn/v2/subscriptions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SMN DELETE /v2/{project_id}/notifications/subscriptions/{subscriptionUrn}
// @API SMN GET /v2/{project_id}/notifications/topics/{topicUrn}/subscriptions
// @API SMN PUT /v2/{project_id}/notifications/topics/{topic_urn}/subscriptions/{subscription_urn}
// @API SMN POST /v2/{project_id}/notifications/topics/{topicUrn}/subscriptions
func ResourceSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubscriptionCreate,
		ReadContext:   resourceSubscriptionRead,
		UpdateContext: resourceSubscriptionUpdate,
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
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extension": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"client_secret": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"keyword": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"sign_secret": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"header": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
			"filter_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The message filter policies of a subscriber.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The filter policy name. The policy name must be unique.`,
						},
						"string_equals": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The string array for exact match.`,
						},
					},
				},
			},
		},
	}
}

func resourceSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Get("topic_urn").(string)
	createOpts := subscriptions.CreateOps{
		Endpoint:  d.Get("endpoint").(string),
		Protocol:  d.Get("protocol").(string),
		Remark:    d.Get("remark").(string),
		Extension: buildExtensionOpts(d.Get("extension").([]interface{})),
	}

	log.Printf("[DEBUG] create Options: %#v", createOpts)
	subscription, err := subscriptions.Create(client, createOpts, topicUrn).Extract()
	if err != nil {
		return diag.Errorf("error getting subscription from result: %s", err)
	}

	log.Printf("[DEBUG] create SMN subscription: %s", subscription.SubscriptionUrn)
	d.SetId(subscription.SubscriptionUrn)
	return resourceSubscriptionRead(ctx, d, meta)
}

func buildExtensionOpts(extensionRaw []interface{}) *subscriptions.ExtensionSpec {
	if len(extensionRaw) == 0 || extensionRaw[0] == nil {
		return nil
	}

	extension := extensionRaw[0].(map[string]interface{})

	res := subscriptions.ExtensionSpec{
		ClientID:     extension["client_id"].(string),
		ClientSecret: extension["client_secret"].(string),
		Keyword:      extension["keyword"].(string),
		SignSecret:   extension["sign_secret"].(string),
		Header:       extension["header"].(map[string]interface{}),
	}

	return &res
}

func resourceSubscriptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	id := d.Id()
	topicUrn := d.Get("topic_urn").(string)

	log.Printf("[DEBUG] fetching subscription: %s", id)
	subscriptionslist, err := subscriptions.ListFromTopic(client, topicUrn).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching the subscriptions")
	}

	var targetSubscription *subscriptions.SubscriptionGet
	for i := range subscriptionslist {
		if subscriptionslist[i].SubscriptionUrn == id {
			targetSubscription = &subscriptionslist[i]
			break
		}
	}

	if targetSubscription == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("topic_urn", targetSubscription.TopicUrn),
		d.Set("endpoint", targetSubscription.Endpoint),
		d.Set("protocol", targetSubscription.Protocol),
		d.Set("subscription_urn", targetSubscription.SubscriptionUrn),
		d.Set("owner", targetSubscription.Owner),
		d.Set("remark", targetSubscription.Remark),
		d.Set("status", targetSubscription.Status),
		d.Set("filter_policies", flattenSubscriptionFilterPolicies(targetSubscription.FilterPolicies)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SMN topic fields: %s", err)
	}
	return nil
}

func resourceSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if d.HasChange("remark") {
		var (
			httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/subscriptions/{subscription_urn}"
			product = "smn"
		)
		client, err := cfg.NewServiceClient(product, region)
		if err != nil {
			return diag.Errorf("error creating SMN client: %s", err)
		}

		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{topic_urn}", d.Get("topic_urn").(string))
		updatePath = strings.ReplaceAll(updatePath, "{subscription_urn}", d.Id())

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateOpt.JSONBody = buildUpdateSubscriptionBodyParams(d)
		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating SMN subscriptions: %s", err)
		}
	}
	return resourceSubscriptionRead(ctx, d, meta)
}

func buildUpdateSubscriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"remark": d.Get("remark"),
	}
	return bodyParams
}

func resourceSubscriptionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	id := d.Id()
	result := subscriptions.Delete(client, id)
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}

	log.Printf("[DEBUG] successfully delete subscription %s", id)
	return nil
}

func flattenSubscriptionFilterPolicies(filterPolicies []subscriptions.FilterPolicy) []map[string]interface{} {
	rst := make([]map[string]interface{}, len(filterPolicies))
	for i, v := range filterPolicies {
		rst[i] = map[string]interface{}{
			"name":          v.Name,
			"string_equals": v.StringEquals,
		}
	}
	return rst
}

func resourceSubscriptionImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	subscriptionUrn := d.Id()
	index := strings.LastIndex(subscriptionUrn, ":")
	if index == -1 {
		return nil, fmt.Errorf("invalid format: the subscription URN is invalid")
	}
	topicUrn := subscriptionUrn[:index]

	d.SetId(subscriptionUrn)
	d.Set("topic_urn", topicUrn)

	return []*schema.ResourceData{d}, nil
}
