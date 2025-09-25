package waf

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

var nonUpdatableDomainRouteUpdateParams = []string{
	"instance_id",
	"routes",
	"routes.*.name",
	"routes.*.servers",
	"routes.*.servers.*.back_protocol",
	"routes.*.servers.*.address",
	"routes.*.servers.*.port",
	"routes.*.cname",
}

// @API WAF PUT /v1/{project_id}/waf/instance/{instance_id}/route
func ResourceDomainRouteUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainRouteUpdateCreate,
		ReadContext:   resourceDomainRouteUpdateRead,
		UpdateContext: resourceDomainRouteUpdateUpdate,
		DeleteContext: resourceDomainRouteUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableDomainRouteUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"servers": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"back_protocol": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"address": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"cname": {
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

func buildDomainRouteUpdateBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawRoutes := d.Get("routes").([]interface{})
	routeParams := make([]map[string]interface{}, 0, len(rawRoutes))

	for _, r := range rawRoutes {
		routeMap := r.(map[string]interface{})
		rawServers := routeMap["servers"].([]interface{})

		servers := make([]map[string]interface{}, 0, len(rawServers))
		for _, s := range rawServers {
			serverMap := s.(map[string]interface{})
			server := map[string]interface{}{
				"back_protocol": utils.ValueIgnoreEmpty(serverMap["back_protocol"]),
				"address":       utils.ValueIgnoreEmpty(serverMap["address"]),
				"port":          utils.ValueIgnoreEmpty(serverMap["port"]),
			}
			servers = append(servers, utils.RemoveNil(server))
		}

		route := map[string]interface{}{
			"name":    routeMap["name"],
			"servers": servers,
			"cname":   utils.ValueIgnoreEmpty(routeMap["cname"]),
		}
		routeParams = append(routeParams, utils.RemoveNil(route))
	}

	return routeParams
}

func resourceDomainRouteUpdateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "waf"
		httpUrl    = "v1/{project_id}/waf/instance/{instance_id}/route"
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: buildDomainRouteUpdateBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating WAF domain route: %s", err)
	}

	d.SetId(instanceId)

	return resourceDomainRouteUpdateRead(ctx, d, meta)
}

func resourceDomainRouteUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDomainRouteUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDomainRouteUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to update WAF domain route. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from the
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
