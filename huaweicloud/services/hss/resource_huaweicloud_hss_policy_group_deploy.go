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

var policyGroupDeployNonUpdatableParams = []string{"enterprise_project_id", "target_policy_group_id", "operate_all", "host_id_list"}

// @API HSS POST /v5/{project_id}/policy/deploy
func ResourcePolicyGroupDeploy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyGroupDeployCreate,
		ReadContext:   resourcePolicyGroupDeployRead,
		UpdateContext: resourcePolicyGroupDeployUpdate,
		DeleteContext: resourcePolicyGroupDeployDelete,

		CustomizeDiff: config.FlexibleForceNew(policyGroupDeployNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region to which the HSS policy group resource belongs.",
			},
			"target_policy_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the policy group to be deployed.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project.",
			},
			"operate_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to deploy the policy on all hosts.`,
			},
			"host_id_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the ID list of servers where the policy group needs to be deployed.`,
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

func createPolicyGroupDeploy(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	requestPath := client.Endpoint + "v5/{project_id}/policy/deploy"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	enterpriseProjectID := cfg.GetEnterpriseProjectID(d)
	if enterpriseProjectID != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", enterpriseProjectID)
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildPolicyGroupDeployBodyParams(d)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func buildPolicyGroupDeployBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"target_policy_group_id": d.Get("target_policy_group_id"),
		"operate_all":            d.Get("operate_all"),
		"host_id_list":           utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("host_id_list").([]interface{}))),
	}
}

func resourcePolicyGroupDeployCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		product             = "hss"
		targetPolicyGroupId = d.Get("target_policy_group_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err := createPolicyGroupDeploy(client, d, cfg); err != nil {
		return diag.Errorf("error deploying HSS policy group: %s", err)
	}

	d.SetId(targetPolicyGroupId)

	return resourcePolicyGroupDeployRead(ctx, d, meta)
}

func resourcePolicyGroupDeployRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourcePolicyGroupDeployUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourcePolicyGroupDeployDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to deploy HSS policy group. Deleting this resource
will not change the current HSS policy group, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
