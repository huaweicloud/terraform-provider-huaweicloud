package hss

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

var appWhitelistPolicyProcessNonUpdatableParams = []string{
	"policy_id", "process_status", "process_hash_list", "enterprise_project_id",
}

// @API HSS PUT /v5/{project_id}/app/{policy_id}/process
func ResourceAppWhitelistPolicyProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppWhitelistPolicyProcessCreate,
		ReadContext:   resourceAppWhitelistPolicyProcessRead,
		UpdateContext: resourceAppWhitelistPolicyProcessUpdate,
		DeleteContext: resourceAppWhitelistPolicyProcessDelete,

		CustomizeDiff: config.FlexibleForceNew(appWhitelistPolicyProcessNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located. If omitted, the provider-level region will be used.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the whitelist policy.",
			},
			"process_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the trust status of the process.",
			},
			"process_hash_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specifies the list of process hash values to mark.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
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

func buildAppWhitelistPolicyProcessQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	rst := ""
	if epsID := cfg.GetEnterpriseProjectID(d); epsID != "" {
		rst = fmt.Sprintf("?enterprise_project_id=%s", epsID)
	}
	return rst
}

func buildAppWhitelistPolicyProcessBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"process_status":    d.Get("process_status"),
		"process_hash_list": utils.ExpandToStringList(d.Get("process_hash_list").([]interface{})),
	}
}

func resourceAppWhitelistPolicyProcessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		policyID = d.Get("policy_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/app/{policy_id}/process"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", policyID)
	requestPath += buildAppWhitelistPolicyProcessQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAppWhitelistPolicyProcessBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error marking app whitelist policy process status: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(id)

	return resourceAppWhitelistPolicyProcessRead(ctx, d, meta)
}

func resourceAppWhitelistPolicyProcessRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceAppWhitelistPolicyProcessUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceAppWhitelistPolicyProcessDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to mark app whitelist policy process status. Deleting
	this resource will not affect the marking status, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
