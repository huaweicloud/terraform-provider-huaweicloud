package hss

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

var closeHoneypotPortPolicyNonUpdatableParams = []string{"policy_id", "host_id", "enterprise_project_id"}

// @API HSS DELETE /v5/{project_id}/honeypot-port/host-policy/{policy_id}
func ResourceCloseHoneypotPortPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloseHoneypotPortPolicyCreate,
		ReadContext:   resourceCloseHoneypotPortPolicyRead,
		UpdateContext: resourceCloseHoneypotPortPolicyUpdate,
		DeleteContext: resourceCloseHoneypotPortPolicyDelete,

		CustomizeDiff: config.FlexibleForceNew(closeHoneypotPortPolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_id": {
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

func buildCloseHoneypotPortPolicyBodyParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?host_id=%v", d.Get("host_id"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%s", queryParams, epsId)
	}

	return queryParams
}

func resourceCloseHoneypotPortPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v5/{project_id}/honeypot-port/host-policy/{policy_id}"
		epsId    = cfg.GetEnterpriseProjectID(d)
		policyId = d.Get("policy_id").(string)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", policyId)
	requestPath += buildCloseHoneypotPortPolicyBodyParams(d, epsId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error closing dynamic port honeypot port policy: %s", err)
	}

	d.SetId(policyId)

	return resourceCloseHoneypotPortPolicyRead(ctx, d, meta)
}

func resourceCloseHoneypotPortPolicyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceCloseHoneypotPortPolicyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceCloseHoneypotPortPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to close honeypot port policy. Deleting
	this resource will not clear the corresponding close records, but will only remove the resource information from
	the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
