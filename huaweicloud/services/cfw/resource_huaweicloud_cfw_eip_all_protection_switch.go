package cfw

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eipAllProtectionSwitchNonUpdatableParams = []string{"fw_instance_id", "bypass_operation"}

// @API CFW POST /v1/{project_id}/eip/protect/all/{fw_instance_id}/operation
func ResourceEipAllProtectionSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipAllProtectionSwitchCreate,
		ReadContext:   resourceEipAllProtectionSwitchRead,
		UpdateContext: resourceEipAllProtectionSwitchUpdate,
		DeleteContext: resourceEipAllProtectionSwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(eipAllProtectionSwitchNonUpdatableParams),

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
			"bypass_operation": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fail_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEipAllProtectionSwitchCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		httpUrl         = "v1/{project_id}/eip/protect/all/{fw_instance_id}/operation"
		fwInstanceID    = d.Get("fw_instance_id").(string)
		bypassOperation = d.Get("bypass_operation").(int)
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{fw_instance_id}", fwInstanceID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"bypass_operation": bypassOperation,
		},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error switching CFW EIP all protection (%s): %v", err, bypassOperation)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error switching CFW EIP all protection, ID is not found in API response")
	}

	d.SetId(id)

	failReason := utils.PathSearch("fail_reason", respBody, "").(string)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("fw_instance_id", id),
		d.Set("object_id", utils.PathSearch("data.object_id", respBody, nil)),
		d.Set("fail_reason", failReason),
	)

	if failReason != "" {
		return diag.Errorf("error switching CFW EIP all protection, fail reason is: %s", failReason)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEipAllProtectionSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceEipAllProtectionSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceEipAllProtectionSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to switch CFW EIP all protection. Deleting this
    resource will not revert the operation or clear the cloud status, but will only remove the resource information
    from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
