package rms

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var resourceAggregationPendingRequestDeleteNonUpdatableParams = []string{"requester_account_id"}

// @API Config DELETE /v1/resource-manager/domains/{domain_id}/aggregators/pending-aggregation-request/{requester_account_id}
func ResourceRmsResourceAggregationPendingRequestDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAggregationPendingRequestDeleteCreate,
		ReadContext:   resourceAggregationPendingRequestDeleteRead,
		UpdateContext: resourceAggregationPendingRequestDeleteUpdate,
		DeleteContext: resourceAggregationPendingRequestDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(resourceAggregationPendingRequestDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"requester_account_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceAggregationPendingRequestDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/pending-aggregation-request/{requester_account_id}"
		product = "rms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{domain_id}", cfg.DomainID)
	createPath = strings.ReplaceAll(createPath, "{requester_account_id}", d.Get("requester_account_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Config resource aggregation pending request delete: %s", err)
	}

	d.SetId(d.Get("requester_account_id").(string))

	return nil
}

func resourceAggregationPendingRequestDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAggregationPendingRequestDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAggregationPendingRequestDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting Config resource aggregation pending request delete is not supported. The resource is only " +
		"removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
