package waf

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

var nonUpdatablePolicyCopyParams = []string{
	"src_policy_id",
	"dest_policy_name",
	"enterprise_project_id",
}

// @API WAF POST /v1/{project_id}/waf/policies/{src_policy_id}/copy
func ResourcePolicyCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCopyCreate,
		ReadContext:   resourcePolicyCopyRead,
		UpdateContext: resourcePolicyCopyUpdate,
		DeleteContext: resourcePolicyCopyDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatablePolicyCopyParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"src_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_policy_name": {
				Type:     schema.TypeString,
				Required: true,
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

func buildPolicyCopyQueryParams(d *schema.ResourceData, epsId string) string {
	req := fmt.Sprintf("?dest_policy_name=%s", d.Get("dest_policy_name").(string))

	if epsId != "" {
		req = fmt.Sprintf("%s&enterprise_project_id=%s", req, epsId)
	}

	return req
}

func resourcePolicyCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "waf"
		httpUrl = "v1/{project_id}/waf/policies/{src_policy_id}/copy"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{src_policy_id}", d.Get("src_policy_id").(string))
	requestPath += buildPolicyCopyQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error copying WAF policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error copying WAF policy: ID is not found in API response")
	}

	d.SetId(id)

	return resourcePolicyCopyRead(ctx, d, meta)
}

func resourcePolicyCopyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourcePolicyCopyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourcePolicyCopyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to copy WAF policy. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from the
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
