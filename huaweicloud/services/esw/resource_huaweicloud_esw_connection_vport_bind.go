package esw

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

var connectionVportBindNonUpdatableParams = []string{"connection_id", "port_id"}

// @API ESW POST /v3/{project_id}/l2cg/connections/{connection_id}/vports/bind
// @API ESW POST /v3/{project_id}/l2cg/connections/{connection_id}/vports/unbind
func ResourceConnectionVportBind() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionVportBindCreate,
		ReadContext:   resourceConnectionVportBindRead,
		UpdateContext: resourceConnectionVportBindUpdate,
		DeleteContext: resourceConnectionVportBindDelete,

		CustomizeDiff: config.FlexibleForceNew(connectionVportBindNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_id": {
				Type:     schema.TypeString,
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

func resourceConnectionVportBindCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/l2cg/connections/{connection_id}/vports/bind"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{connection_id}", d.Get("connection_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateConnectionVportBindBodyParams(d))
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ESW connection vport bind: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("connection_id").(string), d.Get("port_id").(string)))

	return resourceConnectionVportBindRead(ctx, d, meta)
}

func buildCreateConnectionVportBindBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vport": map[string]interface{}{
			"id": d.Get("port_id").(string),
		},
	}
	return bodyParams
}

func resourceConnectionVportBindRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceConnectionVportBindUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceConnectionVportBindDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/l2cg/connections/{connection_id}/vports/unbind"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{connection_id}", d.Get("connection_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteConnectionVportBindBodyParams(d))
	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting ESW connection vport bind: %s", err)
	}

	return nil
}

func buildDeleteConnectionVportBindBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vport": map[string]interface{}{
			"id": d.Get("port_id").(string),
		},
	}
	return bodyParams
}
