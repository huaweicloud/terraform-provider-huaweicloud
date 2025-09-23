package dws

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/db-manager/om-user/action
func ResourceOmAccountAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOmAccountActionCreate,
		ReadContext:   resourceOmAccountActionRead,
		DeleteContext: resourceOmAccountActionDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the DWS cluster ID.`,
			},
			"operation": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the operation type of the OM account.`,
			},
		},
	}
}

func resourceOmAccountActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		httpUrl   = "v1/{project_id}/clusters/{cluster_id}/db-manager/om-user/action"
		clusterId = d.Get("cluster_id").(string)
		action    = d.Get("operation").(string)
	)

	// If the operation is "increaseOmUserPeriod", only the first interface sent will take effect during concurrency.
	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"operation": action,
		},
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error operating (%s) OM account for DWS cluster (%s): %s", action, clusterId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// If the OM account switch is disabled, when the expiration time is extended, the status code is 200, the error code is 1,
	// and the error message is not empty.
	errMsg := utils.PathSearch("error_msg", respBody, "").(string)
	if errMsg != "" {
		return diag.Errorf("error extending the validity period for the OM account: %s", errMsg)
	}

	d.SetId(clusterId)
	return nil
}

func resourceOmAccountActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOmAccountActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the OM account. Deleting this resource will
	not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
