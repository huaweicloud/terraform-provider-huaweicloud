package eg

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eventSubscriptionNonUpdateParams = []string{"subscription_ids", "operation", "enterprise_project_id"}

// @API EG POST /v1/{project_id}/subscriptions/operation
func ResourceEventSubscriptionBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventSubscriptionBatchActionCreate,
		ReadContext:   resourceEventSubscriptionBatchActionRead,
		UpdateContext: resourceEventSubscriptionBatchActionUpdate,
		DeleteContext: resourceEventSubscriptionBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(eventSubscriptionNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the event subscriptions are located.`,
			},
			"subscription_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of subscription IDs to be operated.`,
			},
			"operation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Whether to enable the event subscription.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the subscriptions belong.`,
			},
			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildEventSubscriptionBatchActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"subscription_ids": d.Get("subscription_ids"),
		"operation":        d.Get("operation"),
	}
}

func buildEventSubscriptionBatchActionQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res += fmt.Sprintf("?enterprise_project_id=%s", v.(string))
	}
	return res
}

func checkEventSubscriptionBatchActionResult(respBody interface{}) error {
	failedCount := utils.PathSearch("failed_count", respBody, float64(0)).(float64)
	if failedCount != 0 {
		events := utils.PathSearch("events", respBody, make([]interface{}, 0)).([]interface{})
		var failedSubscriptions []string
		for _, event := range events {
			if errorCode := utils.PathSearch("error_code", event, nil); errorCode != nil {
				subscriptionId := utils.PathSearch("subscription_id", event, "").(string)
				errorMsg := utils.PathSearch("error_msg", event, "").(string)
				failedSubscriptions = append(failedSubscriptions, fmt.Sprintf("subscription %s: %s (%s)", subscriptionId, errorMsg, errorCode))
			}
		}
		return fmt.Errorf("failed to operate %v subscription(s): %s", failedCount, strings.Join(failedSubscriptions, "; "))
	}
	return nil
}

func resourceEventSubscriptionBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/subscriptions/operation"
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildEventSubscriptionBatchActionQueryParams(d)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildEventSubscriptionBatchActionBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("unable to batch operate the event subscriptions: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	err = checkEventSubscriptionBatchActionResult(respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceEventSubscriptionBatchActionRead(ctx, d, meta)
}

func resourceEventSubscriptionBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEventSubscriptionBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEventSubscriptionBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating event subscription status. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
 file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
