package cdn

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

var domainBatchCopyNonUpdatableParams = []string{"source_domain", "target_domains", "config_list"}

// @API CDN POST /v1.0/cdn/configuration/domains/batch-copy
func ResourceDomainBatchCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainBatchCopyCreate,
		ReadContext:   resourceDomainBatchCopyRead,
		UpdateContext: resourceDomainBatchCopyUpdate,
		DeleteContext: resourceDomainBatchCopyDelete,

		CustomizeDiff: config.FlexibleForceNew(domainBatchCopyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"source_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The source domain whose configuration will be copied.`,
			},
			"target_domains": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The target domain names to which configurations will be copied.`,
			},
			"config_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The configuration items to copy.`,
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

func buildDomainBatchCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"configs": map[string]interface{}{
			"target_domain": d.Get("target_domains"),
			"source_domain": d.Get("source_domain"),
			"config_list":   utils.ExpandToStringList(d.Get("config_list").([]interface{})),
		},
	}
}

func collectBatchCopyErrors(results []interface{}) error {
	var errors []string

	for _, result := range results {
		status := utils.PathSearch("status", result, "").(string)
		domainName := utils.PathSearch("domain_name", result, "").(string)
		reason := utils.PathSearch("reason", result, "").(string)

		if status != "success" && status != "completed" && status != "" {
			errorMsg := fmt.Sprintf("domain '%s' failed with status '%s', reason: '%s'", domainName, status, reason)
			if reason != "" {
				errors = append(errors, errorMsg)
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("batch copy operation failed for %d domain(s): %s", len(errors), strings.Join(errors, "; "))
	}
	return nil
}

func resourceDomainBatchCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1.0/cdn/configuration/domains/batch-copy"
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildDomainBatchCopyBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error performing domain batch copy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	results := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
	if err := collectBatchCopyErrors(results); err != nil {
		return diag.Errorf("error in batch copy operation: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceDomainBatchCopyRead(ctx, d, meta)
}

func resourceDomainBatchCopyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDomainBatchCopyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDomainBatchCopyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for performing domain batch copy operation. Deleting this 
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
