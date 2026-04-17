package vpcep

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

var vpcepServiceAddServersNonUpdatableParams = []string{"vpc_endpoint_service_id", "server_resources",
	"server_resources.*.resource_id", "server_resources.*.availability_zone_id",
}

// @API VPCEP POST /v2/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/add-server-resources
func ResourceVPCEndpointServiceAddServers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCEndpointServiceAddServersCreate,
		ReadContext:   resourceVPCEndpointServiceAddServersRead,
		UpdateContext: resourceVPCEndpointServiceAddServersUpdate,
		DeleteContext: resourceVPCEndpointServiceAddServersDelete,

		CustomizeDiff: config.FlexibleForceNew(vpcepServiceAddServersNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_endpoint_service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"availability_zone_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
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

func resourceVPCEndpointServiceAddServersCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                  = meta.(*config.Config)
		region               = cfg.GetRegion(d)
		httpUrl              = "v2/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/add-server-resources"
		vpcEndpointServiceId = d.Get("vpc_endpoint_service_id").(string)
	)

	client, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vpc_endpoint_service_id}", vpcEndpointServiceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildAddServersRequestBody(d),
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error adding servers to VPC endpoint service: %s", err)
	}

	d.SetId(vpcEndpointServiceId)

	return nil
}

func buildAddServersRequestBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"server_resources": buildAddServerResources(d.Get("server_resources")),
	}
	return bodyParams
}

func buildAddServerResources(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"resource_id":          utils.ValueIgnoreEmpty(raw["resource_id"]),
				"availability_zone_id": utils.ValueIgnoreEmpty(raw["availability_zone_id"]),
			}
		}
		return rst
	}
	return nil
}

func resourceVPCEndpointServiceAddServersRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVPCEndpointServiceAddServersUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVPCEndpointServiceAddServersDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting VPCEP service add servers resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
