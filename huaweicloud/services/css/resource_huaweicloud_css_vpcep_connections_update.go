package css

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var vpcepserviceConnectionsUpdateNonUpdatableParams = []string{
	"cluster_id", "action", "endpoint_id_list",
}

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/connections
func ResourceCssVpcepserviceConnectionsUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCssVpcepserviceConnectionsUpdateCreate,
		ReadContext:   resourceCssVpcepserviceConnectionsUpdateRead,
		UpdateContext: resourceCssVpcepserviceConnectionsUpdateUpdate,
		DeleteContext: resourceCssVpcepserviceConnectionsUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(vpcepserviceConnectionsUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint_id_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceCssVpcepserviceConnectionsUpdateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/connections"
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildUpdateRequestBody(d),
	}

	_, err = client.Request("POST", requestPath, &createOpt)
	if err != nil {
		return diag.Errorf("error updating CSS cluster (%s) vpcep service connections: %s", clusterId, err)
	}

	d.SetId(clusterId)

	return nil
}

func buildUpdateRequestBody(d *schema.ResourceData) map[string]interface{} {
	requestBody := map[string]interface{}{
		"action":           d.Get("action").(string),
		"endpoint_id_list": utils.ExpandToStringList(d.Get("endpoint_id_list").([]interface{})),
	}
	return requestBody
}

func resourceCssVpcepserviceConnectionsUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCssVpcepserviceConnectionsUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCssVpcepserviceConnectionsUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "This resource only removes the state. The VPC endpoint connections on the CSS cluster remain unchanged."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
