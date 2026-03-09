package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var deleteIpBlacklistNonUpdatableParams = []string{"fw_instance_id", "effect_scope"}

// @API CFW DELETE /v1/{project_id}/ptf/ip-blacklist
func ResourceDeleteIpBlacklist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeleteIpBlacklistCreate,
		ReadContext:   resourceDeleteIpBlacklistRead,
		UpdateContext: resourceDeleteIpBlacklistUpdate,
		DeleteContext: resourceDeleteIpBlacklistDelete,

		CustomizeDiff: config.FlexibleForceNew(deleteIpBlacklistNonUpdatableParams),

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
			"effect_scope": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
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

func buildDeleteIpBlacklistQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?fw_instance_id=%s", d.Get("fw_instance_id"))
}

func buildDeleteIpBlacklistBodyParams(d *schema.ResourceData) map[string]interface{} {
	effectScopeInput := d.Get("effect_scope").([]interface{})
	if len(effectScopeInput) == 0 {
		return nil
	}

	return map[string]interface{}{
		"effect_scope": utils.ExpandToIntList(effectScopeInput),
	}
}

func resourceDeleteIpBlacklistCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/ptf/ip-blacklist"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDeleteIpBlacklistQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildDeleteIpBlacklistBodyParams(d),
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting CFW imported IP blacklist: %s", err)
	}

	d.SetId(d.Get("fw_instance_id").(string))

	return nil
}

func resourceDeleteIpBlacklistRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDeleteIpBlacklistUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDeleteIpBlacklistDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to delete the imported IP blacklist. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
