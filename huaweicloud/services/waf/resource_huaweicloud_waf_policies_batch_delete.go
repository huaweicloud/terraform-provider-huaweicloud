package waf

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatablePoliciesBatchDeleteParams = []string{
	"policy_ids",
	"enterprise_project_id",
}

// @API WAF POST /v1/{project_id}/waf/policies/batch-delete
func ResourcePoliciesBatchDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoliciesBatchDeleteCreate,
		ReadContext:   resourcePoliciesBatchDeleteRead,
		UpdateContext: resourcePoliciesBatchDeleteUpdate,
		DeleteContext: resourcePoliciesBatchDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatablePoliciesBatchDeleteParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildPoliciesBatchDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_ids": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("policy_ids").([]interface{}))),
	}
}

func resourcePoliciesBatchDeleteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "waf"
		httpUrl = "v1/{project_id}/waf/policies/batch-delete"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId != "" {
		requestPath += "?enterprise_project_id=" + epsId
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildPoliciesBatchDeleteBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch deleting WAF policies: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourcePoliciesBatchDeleteRead(ctx, d, meta)
}

func resourcePoliciesBatchDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourcePoliciesBatchDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourcePoliciesBatchDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to batch delete WAF policies. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from the
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
