package waf

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

var nonUpdatableDedicatedInstanceActionParams = []string{
	"instance_id",
	"action",
	"params",
}

// @API WAF POST /v1/{project_id}/premium-waf/instance/{instance_id}/action
func ResourceDedicatedInstanceAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedInstanceActionCreate,
		ReadContext:   resourceDedicatedInstanceActionRead,
		UpdateContext: resourceDedicatedInstanceActionUpdate,
		DeleteContext: resourceDedicatedInstanceActionDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableDedicatedInstanceActionParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"instancename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"float_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"run_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"upgradable": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cloud_service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCreateDedicatedInstanceActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"action": d.Get("action"),
		"params": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("params").([]interface{}))),
	}
}

func resourceDedicatedInstanceActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + "v1/{project_id}/premium-waf/instance/{instance_id}/action"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDedicatedInstanceActionBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating WAF dedicated instance action: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating WAF dedicated instance action: ID is not found in API response")
	}

	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", id),
		d.Set("instancename", utils.PathSearch("instancename", respBody, nil)),
		d.Set("server_id", utils.PathSearch("serverId", respBody, nil)),
		d.Set("zone", utils.PathSearch("zone", respBody, nil)),
		d.Set("arch", utils.PathSearch("arch", respBody, nil)),
		d.Set("cpu_flavor", utils.PathSearch("cpu_flavor", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", respBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", respBody, nil)),
		d.Set("service_ip", utils.PathSearch("service_ip", respBody, nil)),
		d.Set("service_ipv6", utils.PathSearch("service_ipv6", respBody, nil)),
		d.Set("float_ip", utils.PathSearch("floatIp", respBody, nil)),
		d.Set("security_group_ids", utils.ExpandToStringList(
			utils.PathSearch("security_group_ids", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("run_status", utils.PathSearch("run_status", respBody, nil)),
		d.Set("access_status", utils.PathSearch("access_status", respBody, nil)),
		d.Set("upgradable", utils.PathSearch("upgradable", respBody, nil)),
		d.Set("cloud_service_type", utils.PathSearch("cloudServiceType", respBody, nil)),
		d.Set("resource_type", utils.PathSearch("resourceType", respBody, nil)),
		d.Set("resource_spec_code", utils.PathSearch("resourceSpecCode", respBody, nil)),
		d.Set("specification", utils.PathSearch("specification", respBody, nil)),
		d.Set("hosts", flattenDedicatedInstanceActionHosts(
			utils.PathSearch("hosts", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("volume_type", utils.PathSearch("volume_type", respBody, nil)),
		d.Set("cluster_id", utils.PathSearch("cluster_id", respBody, nil)),
		d.Set("pool_id", utils.PathSearch("pool_id", respBody, nil)),
		d.Set("charge_mode", utils.PathSearch("charge_mode", respBody, nil)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return resourceDedicatedInstanceActionRead(ctx, d, meta)
}

func flattenDedicatedInstanceActionHosts(hostsResp []interface{}) []map[string]interface{} {
	if len(hostsResp) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(hostsResp))
	for _, v := range hostsResp {
		result = append(result, map[string]interface{}{
			"id":       utils.PathSearch("id", v, nil),
			"hostname": utils.PathSearch("hostname", v, nil),
		})
	}

	return result
}

func resourceDedicatedInstanceActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDedicatedInstanceActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDedicatedInstanceActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to operate WAF dedicated instance. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from the
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
