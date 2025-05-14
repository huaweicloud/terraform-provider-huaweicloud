package cbc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var resourcesUnsubscribeNonUpdatableParams = []string{"resource_ids"}

const (
	// '1' means unsubscribing resources.
	// '2' means unsubscribing resources only for the renewal period.
	unsubscribeType = 1
	// '2' means the resource is in use.
	resourceInUse = 2
	// '0' means query the main resource and the subsidiary resources.
	// '1' means query only the main resource.
	onlyMainResource = 1
)

// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceResourcesUnsubscribe() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourcesUnsubscribeCreate,
		ReadContext:   resourceResourcesUnsubscribeRead,
		UpdateContext: resourceResourcesUnsubscribeUpdate,
		DeleteContext: resourceResourcesUnsubscribeDelete,

		CustomizeDiff: config.FlexibleForceNew(resourcesUnsubscribeNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The IDs of the resource to be unsubscribed.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceResourcesUnsubscribeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("bssv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS client: %s", err)
	}

	resourceIds := utils.ExpandToStringList(d.Get("resource_ids").([]interface{}))
	err = unsubscribePrePaidResources(client, resourceIds)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForResourcesUnsubscribed(ctx, client, resourceIds, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for all resources to be unsubscribed: %s ", err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomId)

	return nil
}

func unsubscribePrePaidResources(client *golangsdk.ServiceClient, resourceIds []string) error {
	httpUrl := "v2/orders/subscriptions/resources/unsubscribe"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"resource_ids":     resourceIds,
			"unsubscribe_type": unsubscribeType,
		},
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error unsubscribing the resources: %s", err)
	}

	respbody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	failResources := utils.PathSearch("fail_resource_infos", respbody, make([]interface{}, 0)).([]interface{})
	if len(failResources) > 0 {
		failInfos := make([]string, 0)
		for _, v := range failResources {
			failInfos = append(failInfos, fmt.Sprintf("%s | %s;",
				utils.PathSearch("resource_id", v, "").(string),
				utils.PathSearch("error_msg", v, "").(string)))
		}

		log.Printf("[ERROR] error unsubscribing some resources: %s", strings.Join(failInfos, "\n"))
	}

	return nil
}

func waitForResourcesUnsubscribed(ctx context.Context, client *golangsdk.ServiceClient, resourceIds []string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshPrePaidResourcesByIds(client, resourceIds),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshPrePaidResourcesByIds(client *golangsdk.ServiceClient, resourceIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl = "v2/orders/suscriptions/resources/query"
			limit   = 100
			offset  = 0
			result  = make([]interface{}, 0)
			getOpt  = golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			bodyParams = map[string]interface{}{
				"resource_ids":       resourceIds,
				"status_list":        []int{resourceInUse},
				"only_main_resource": onlyMainResource,
				// The 'limit' default value is `10`.
				"limit":  limit,
				"offset": offset,
			}
		)

		listPath := client.Endpoint + httpUrl
		for {
			bodyParams["offset"] = offset
			getOpt.JSONBody = bodyParams

			resp, err := client.Request("POST", listPath, &getOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			resources := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
			result = append(result, resources...)
			if len(resources) < limit {
				break
			}
			offset += len(resources)
		}

		if len(result) > 0 {
			return result, "PENDING", nil
		}
		return result, "COMPLETED", nil
	}
}

func resourceResourcesUnsubscribeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourcesUnsubscribeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourcesUnsubscribeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for unsubscribing the specified resources. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
