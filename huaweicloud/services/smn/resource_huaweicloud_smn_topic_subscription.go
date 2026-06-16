package smn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var topicSubscriptionNonUpdatableParams = []string{
	"token",
	"topic_urn",
	"endpoint",
}

// @API SMN GET /v2/notifications/subscriptions/subscribe
func ResourceTopicSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicSubscriptionCreate,
		ReadContext:   resourceTopicSubscriptionRead,
		UpdateContext: resourceTopicSubscriptionUpdate,
		DeleteContext: resourceTopicSubscriptionDelete,

		CustomizeDiff: config.FlexibleForceNew(topicSubscriptionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"topic_urn": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"endpoint"},
			},
			"endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildTopicSubscriptionQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?token=%v", d.Get("token"))

	if v, ok := d.GetOk("topic_urn"); ok {
		queryParams = fmt.Sprintf("%s&topic_urn=%v", queryParams, v)
	}

	if v, ok := d.GetOk("endpoint"); ok {
		queryParams = fmt.Sprintf("%s&endpoint=%v", queryParams, v)
	}

	return queryParams
}

func resourceTopicSubscriptionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/notifications/subscriptions/subscribe"
	)

	client, err := cfg.NewServiceClient("smn", region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildTopicSubscriptionQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error subscribing SMN topic: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	return diag.FromErr(d.Set("subscription_urn", utils.PathSearch("subscription_urn", respBody, nil)))
}

func resourceTopicSubscriptionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceTopicSubscriptionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceTopicSubscriptionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting topic subscription resource is not supported. The resource is only removed from the " +
		"state, the resource remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
