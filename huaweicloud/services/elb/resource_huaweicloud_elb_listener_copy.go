package elb

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var listenerCopyNonUpdatableParams = []string{
	"listener_id", "protocol", "protocol_port", "port_ranges", "loadbalancer_id", "reuse_pool",
}

// @API ELB POST /v3/{project_id}/elb/listeners/{listener_id}/clone
// @API ELB POST /v2.0/{project_id}/listeners/{listener_id}/tags/action
// @API ELB GET /v3/{project_id}/elb/listeners/{listener_id}
// @API ELB GET /v2.0/{project_id}/listeners/{listener_id}/tags
// @API ELB PUT /v3/{project_id}/elb/listeners/{listener_id}
// @API ELB DELETE /v3/{project_id}/elb/listeners/{listener_id}/force
// @API ELB DELETE /v3/{project_id}/elb/listeners/{listener_id}
// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}
func ResourceListenerCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListenerCopyCreate,
		ReadContext:   resourceListenerV3Read,
		UpdateContext: resourceListenerV3Update,
		DeleteContext: resourceListenerV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(listenerCopyNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"port_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"end_port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"reuse_pool": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http2_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"forward_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"forward_port": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"forward_request_port": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"forward_host": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"forward_elb": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"forward_proto": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"real_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"forward_tls_certificate": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"forward_tls_cipher": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"forward_tls_protocol": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"access_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"ip_group"},
			},
			"ip_group": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"access_policy"},
			},
			"ip_group_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"access_policy"},
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"server_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sni_certificate": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idle_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"request_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"response_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"advanced_forwarding_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gzip_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_member_retry": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"proxy_protocol_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"sni_match_algo": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_early_data_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"quic_listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_quic_upgrade": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"quic_listener_id"},
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"max_connection": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"cps": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"transparent_client_ip_enable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"nat64_enable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"tags": common.TagsSchema(),
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceListenerCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/listeners/{listener_id}/clone"
		product = "elbv3"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{listener_id}", d.Get("listener_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateListenerCopyBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB listener copy: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	listenerId := utils.PathSearch("listener_list[0].id", createRespBody, nil)
	if listenerId == nil {
		return diag.Errorf("error creating ELB listener copy: ID is not found in API response")
	}
	d.SetId(listenerId.(string))

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating ELB listener copy: job_id is not found in API response")
	}
	err = checkLoadBalancerJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if diagErr := initializeListenerParams(ctx, client, d); diagErr != nil {
		return diagErr
	}

	if tagRaw := d.Get("tags").(map[string]interface{}); len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "listeners", d.Id(), tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags for ELB listener copy %s: %s", d.Id(), tagErr)
		}
	}
	return resourceListenerV3Read(ctx, d, meta)
}

func buildCreateListenerCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":            utils.ValueIgnoreEmpty(d.Get("name")),
		"loadbalancer_id": d.Get("loadbalancer_id"),
		"protocol_port":   utils.ValueIgnoreEmpty(d.Get("protocol_port").(int)),
		"port_ranges":     buildCreateListenerCopyPortRangesBodyParams(d),
	}
	if v, ok := d.GetOk("reuse_pool"); ok {
		reusePool, _ := strconv.ParseBool(v.(string))
		bodyParams["reuse_pool"] = reusePool
	}
	return map[string]interface{}{
		"target_listener_params": []interface{}{bodyParams},
	}
}

func buildCreateListenerCopyPortRangesBodyParams(d *schema.ResourceData) []interface{} {
	portRanges := d.Get("port_ranges").(*schema.Set)
	if portRanges.Len() == 0 {
		return nil
	}
	res := make([]interface{}, 0, portRanges.Len())
	for _, portRange := range portRanges.List() {
		if v, ok := portRange.(map[string]interface{}); ok {
			res = append(res, map[string]interface{}{
				"start_port": utils.ValueIgnoreEmpty(v["start_port"]),
				"end_port":   utils.ValueIgnoreEmpty(v["end_port"]),
			})
		}
	}
	return res
}

func initializeListenerParams(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	updateListenerChanges := []string{"description", "ca_certificate", "default_pool_id", "idle_timeout",
		"request_timeout", "response_timeout", "server_certificate", "access_policy", "ip_group", "ip_group_enable",
		"forward_eip", "forward_port", "forward_request_port", "forward_host", "tls_ciphers_policy", "sni_certificate",
		"http2_enable", "gzip_enable", "advanced_forwarding_enabled", "protection_status", "protection_reason",
		"forward_elb", "forward_proto", "real_ip", "forward_tls_certificate", "forward_tls_cipher", "forward_tls_protocol",
		"enable_member_retry", "proxy_protocol_enable", "sni_match_algo", "security_policy_id", "ssl_early_data_enable",
		"quic_listener_id", "enable_quic_upgrade", "max_connection", "cps", "nat64_enable", "transparent_client_ip_enable",
	}
	if d.HasChanges(updateListenerChanges...) {
		return updateListener(ctx, d, client)
	}

	return nil
}
