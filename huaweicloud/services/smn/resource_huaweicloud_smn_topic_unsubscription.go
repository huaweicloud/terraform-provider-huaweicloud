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

var topicUnsubscriptionNonUpdatableParams = []string{
	"subscription_urn",
}

// @API SMN GET /v2/notifications/subscriptions/unsubscribe
func ResourceTopicUnsubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicUnsubscriptionCreate,
		ReadContext:   resourceTopicUnsubscriptionRead,
		UpdateContext: resourceTopicUnsubscriptionUpdate,
		DeleteContext: resourceTopicUnsubscriptionDelete,

		CustomizeDiff: config.FlexibleForceNew(topicUnsubscriptionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTopicUnsubscriptionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/notifications/subscriptions/unsubscribe"
	)

	client, err := cfg.NewServiceClient("smn", region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = fmt.Sprintf("%s?subscription_urn=%s", requestPath, d.Get("subscription_urn").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error unsubscribing SMN topic: %s", err)
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

	return diag.FromErr(d.Set("message", utils.PathSearch("message", respBody, nil)))
}

func resourceTopicUnsubscriptionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceTopicUnsubscriptionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceTopicUnsubscriptionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting topic unsubscription resource is not supported. The resource is only removed from the " +
		"state, the resource remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
