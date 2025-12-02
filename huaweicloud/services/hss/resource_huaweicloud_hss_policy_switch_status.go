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

// @API HSS POST /v5/{project_id}/policy/switch
func ResourcePolicySwitchStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicySwitchStatusCreate,
		ReadContext:   resourcePolicySwitchStatusRead,
		UpdateContext: resourcePolicySwitchStatusUpdate,
		DeleteContext: resourcePolicySwitchStatusDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"enterprise_project_id",
			"policy_name",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable": {
				Type:     schema.TypeBool,
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

func updatePolicySwitchStatus(client *golangsdk.ServiceClient, d *schema.ResourceData, epsId string) error {
	requestPath := client.Endpoint + "v5/{project_id}/policy/switch"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"policy_name": d.Get("policy_name"),
			"enable":      d.Get("enable"),
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourcePolicySwitchStatusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err := updatePolicySwitchStatus(client, d, epsId); err != nil {
		return diag.Errorf("error changing HSS policy status in create operation: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourcePolicySwitchStatusRead(ctx, d, meta)
}

func resourcePolicySwitchStatusRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourcePolicySwitchStatusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err := updatePolicySwitchStatus(client, d, epsId); err != nil {
		return diag.Errorf("error changing HSS policy status in update operation: %s", err)
	}

	return resourcePolicySwitchStatusRead(ctx, d, meta)
}

func resourcePolicySwitchStatusDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to switch HSS policy status. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from the
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
