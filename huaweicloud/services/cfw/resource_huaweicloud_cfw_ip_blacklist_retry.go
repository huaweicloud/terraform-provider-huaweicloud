package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var ipBlacklistRetryNonUpdatableParams = []string{"fw_instance_id", "name"}

// @API CFW POST /v1/{project_id}/ptf/ip-blacklist/retry
func ResourceIpBlacklistRetry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpBlacklistRetryCreate,
		ReadContext:   resourceIpBlacklistRetryRead,
		UpdateContext: resourceIpBlacklistRetryUpdate,
		DeleteContext: resourceIpBlacklistRetryDelete,

		CustomizeDiff: config.FlexibleForceNew(ipBlacklistRetryNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The parameter `name` is optional in the API documentation, but it is actually required.
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildIpBlacklistRetryQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?fw_instance_id=%s&name=%s", d.Get("fw_instance_id"), d.Get("name"))
}

func resourceIpBlacklistRetryCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/ptf/ip-blacklist/retry"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildIpBlacklistRetryQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrying CFW IP blacklist import: %s", err)
	}

	d.SetId(d.Get("fw_instance_id").(string))

	return nil
}

func resourceIpBlacklistRetryRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceIpBlacklistRetryUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceIpBlacklistRetryDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to retry the failed IP blacklist import. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
