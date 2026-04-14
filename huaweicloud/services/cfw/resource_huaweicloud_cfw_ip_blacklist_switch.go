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

var ipBlacklistSwitchNonUpdatableParams = []string{"fw_instance_id", "status"}

// @API CFW POST /v1/{project_id}/ptf/ip-blacklist/switch
func ResourceIpBlacklistSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpBlacklistSwitchCreate,
		ReadContext:   resourceIpBlacklistSwitchRead,
		UpdateContext: resourceIpBlacklistSwitchUpdate,
		DeleteContext: resourceIpBlacklistSwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(ipBlacklistSwitchNonUpdatableParams),

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
			"status": {
				Type:     schema.TypeInt,
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

func buildEnableIpBlacklistQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?fw_instance_id=%s", d.Get("fw_instance_id").(string))
}

func resourceIpBlacklistSwitchCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/ptf/ip-blacklist/switch"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildEnableIpBlacklistQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"status": d.Get("status"),
		},
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating CFW IP blacklist switch: %s", err)
	}

	d.SetId(d.Get("fw_instance_id").(string))

	return nil
}

func resourceIpBlacklistSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceIpBlacklistSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceIpBlacklistSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to enable or disable the IP blacklist feature.
    Deleting this resource will not restore the previous switch state on the cloud, but will only remove the resource
    information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
