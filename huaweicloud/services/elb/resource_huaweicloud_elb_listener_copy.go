package elb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var listenerCopyNonUpdatableParams = []string{
	"listener_id",
	"loadbalancer_id",
	"name",
	"protocol_port",
	"port_ranges",
	"port_ranges.*.start_port",
	"port_ranges.*.end_port",
	"reuse_pool",
}

// @API ELB POST /v3/{project_id}/elb/listeners/{listener_id}/clone
// @API ELB GET /v3/{project_id}/elb/listeners/{listener_id}
// @API ELB DELETE /v3/{project_id}/elb/listeners/{listener_id}
// @API ELB DELETE /v3/{project_id}/elb/listeners/{listener_id}/force
// @API ELB GET /v3/{project_id}/elb/jobs/{job_id}
func ResourceListenerCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListenerCopyCreate,
		ReadContext:   resourceListenerCopyRead,
		UpdateContext: resourceListenerCopyUpdate,
		DeleteContext: resourceListenerCopyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(listenerCopyNonUpdatableParams),

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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"end_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"reuse_pool": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"client_ca_tls_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http2_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"insert_headers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"forwarded_elb_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_port": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_for_port": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_host": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_proto": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"real_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_elb_id": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_tls_certificate_id": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_tls_protocol": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_tls_cipher": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_tls_protocol_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_tls_cipher_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_for_processing_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_clientcert_subjectdn_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_clientcert_subjectdn_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_clientcert_issuerdn_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_clientcert_issuerdn_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_clientcert_fingerprint_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_clientcert_fingerprint_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_clientcert_clientverify_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_clientcert_clientverify_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_clientcert_serialnumber_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_clientcert_serialnumber_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_clientcert_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_clientcert_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_clientcert_ciphers_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_clientcert_ciphers_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_clientcert_end_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_clientcert_end_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_tls_alpn_protocol_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_tls_alpn_protocol_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_tls_sni_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_tls_sni_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_tls_ja3_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_tls_ja3_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarded_tls_ja4_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_tls_ja4_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sni_container_refs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sni_match_algo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_member_retry": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"keepalive_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"client_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"member_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipgroup": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipgroup_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ipgroup": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"transparent_client_ip_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"proxy_protocol_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enhance_l7policy_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"quic_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quic_listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_quic_upgrade": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"protection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gzip_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ssl_early_data_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"connection_attr": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nat64_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildListenerCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"target_listener_params": []map[string]interface{}{
			{
				"loadbalancer_id": d.Get("loadbalancer_id"),
				"name":            utils.ValueIgnoreEmpty(d.Get("name")),
				"protocol_port":   utils.ValueIgnoreEmpty(d.Get("protocol_port")),
				"port_ranges":     buildPortRangestBodyParams(d.Get("port_ranges").([]interface{})),
				"reuse_pool":      d.Get("reuse_pool"),
			},
		},
	}

	return bodyParams
}

func buildPortRangestBodyParams(portList []interface{}) []map[string]interface{} {
	if len(portList) == 0 {
		return nil
	}

	ports := make([]map[string]interface{}, 0, len(portList))
	for _, v := range portList {
		raw, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"start_port": raw["start_port"],
			"end_port":   raw["end_port"],
		}

		ports = append(ports, params)
	}

	return ports
}

func resourceListenerCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		listenerId = d.Get("listener_id").(string)
		httpUrl    = "v3/{project_id}/elb/listeners/{listener_id}/clone"
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{listener_id}", listenerId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildListenerCopyBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error copying the listener: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	targetListenerId := utils.PathSearch("listener_list[0].id", respBody, "").(string)
	if targetListenerId == "" {
		return diag.Errorf("error copying the listener : unable to find the target listener ID from the API response")
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	d.SetId(targetListenerId)

	err = waitForListenerCopyCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the clone of listener (%s) to complete: %s", d.Id(), err)
	}

	return resourceListenerCopyRead(ctx, d, meta)
}

func waitForListenerCopyCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshListenerCopyStatusFunc(client, jobId),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshListenerCopyStatusFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v3/{project_id}/elb/jobs/{job_id}"
		getJobPath := client.Endpoint + httpUrl
		getJobPath = strings.ReplaceAll(getJobPath, "{project_id}", client.ProjectID)
		getJobPath = strings.ReplaceAll(getJobPath, "{job_id}", jobId)
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", getJobPath, &getOpt)
		if err != nil {
			return resp, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return resp, "ERROR", err
		}

		jobStatus := utils.PathSearch("job.status", respBody, "").(string)
		if utils.StrSliceContains([]string{"COMPLETE"}, jobStatus) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func GetListenerCopy(client *golangsdk.ServiceClient, listenerId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/elb/listeners/{listener_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{listener_id}", listenerId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceListenerCopyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	listener, err := GetListenerCopy(client, d.Id())
	if err != nil {
		// When the listener does not exist, the query details API return `404`
		return common.CheckDeletedDiag(d, err, "error retrieving listener")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("loadbalancer_id", utils.PathSearch("listener.loadbalancers[0].id", listener, nil)),
		d.Set("name", utils.PathSearch("listener.name", listener, nil)),
		d.Set("protocol_port", utils.PathSearch("listener.protocol_port", listener, nil)),
		d.Set("port_ranges", flattenListenerCopyPortRanges(
			utils.PathSearch("listener.port_ranges", listener, make([]interface{}, 0)).([]interface{}))),
		d.Set("client_ca_tls_container_ref", utils.PathSearch("listener.client_ca_tls_container_ref", listener, nil)),
		d.Set("connection_limit", utils.PathSearch("listener.connection_limit", listener, nil)),
		d.Set("created_at", utils.PathSearch("listener.created_at", listener, nil)),
		d.Set("default_pool_id", utils.PathSearch("listener.default_pool_id", listener, nil)),
		d.Set("default_tls_container_ref", utils.PathSearch("listener.default_tls_container_ref", listener, nil)),
		d.Set("description", utils.PathSearch("listener.description", listener, nil)),
		d.Set("http2_enable", utils.PathSearch("listener.http2_enable", listener, nil)),
		d.Set("insert_headers", flattenListenerInsertHeaders(utils.PathSearch("listener.insert_headers", listener, nil))),
		d.Set("enterprise_project_id", utils.PathSearch("listener.enterprise_project_id", listener, nil)),
		d.Set("protocol", utils.PathSearch("listener.protocol", listener, nil)),
		d.Set("sni_container_refs", utils.PathSearch("listener.sni_container_refs", listener, nil)),
		d.Set("sni_match_algo", utils.PathSearch("listener.sni_match_algo", listener, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("listener.tags", listener, nil))),
		d.Set("updated_at", utils.PathSearch("listener.updated_at", listener, nil)),
		d.Set("tls_ciphers_policy", utils.PathSearch("listener.tls_ciphers_policy", listener, nil)),
		d.Set("security_policy_id", utils.PathSearch("listener.security_policy_id", listener, nil)),
		d.Set("enable_member_retry", utils.PathSearch("listener.enable_member_retry", listener, nil)),
		d.Set("keepalive_timeout", utils.PathSearch("listener.keepalive_timeout", listener, nil)),
		d.Set("client_timeout", utils.PathSearch("listener.client_timeout", listener, nil)),
		d.Set("member_timeout", utils.PathSearch("listener.member_timeout", listener, nil)),
		d.Set("ipgroup", flattenListenerIpGroup(utils.PathSearch("listener.ipgroup", listener, nil))),
		d.Set("transparent_client_ip_enable", utils.PathSearch("listener.transparent_client_ip_enable", listener, nil)),
		d.Set("proxy_protocol_enable", utils.PathSearch("listener.proxy_protocol_enable", listener, nil)),
		d.Set("enhance_l7policy_enable", utils.PathSearch("listener.enhance_l7policy_enable", listener, nil)),
		d.Set("quic_config", flattenListenerQuicConfig(utils.PathSearch("listener.quic_config", listener, nil))),
		d.Set("protection_status", utils.PathSearch("listener.protection_status", listener, nil)),
		d.Set("protection_reason", utils.PathSearch("listener.protection_reason", listener, nil)),
		d.Set("gzip_enable", utils.PathSearch("listener.gzip_enable", listener, nil)),
		d.Set("ssl_early_data_enable", utils.PathSearch("listener.ssl_early_data_enable", listener, nil)),
		d.Set("cps", utils.PathSearch("listener.cps", listener, nil)),
		d.Set("connection_attr", utils.PathSearch("listener.connection", listener, nil)),
		d.Set("nat64_enable", utils.PathSearch("listener.nat64_enable", listener, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListenerCopyPortRanges(portList []interface{}) []interface{} {
	if len(portList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(portList))
	for _, v := range portList {
		result = append(result, map[string]interface{}{
			"start_port": utils.PathSearch("start_port", v, nil),
			"end_port":   utils.PathSearch("end_port", v, nil),
		})
	}

	return result
}

func flattenListenerInsertHeaders(insertHeaders interface{}) []map[string]interface{} {
	if insertHeaders == nil {
		return nil
	}

	result := map[string]interface{}{
		"forwarded_elb_ip":                         utils.PathSearch("X-Forwarded-ELB-IP", insertHeaders, nil),
		"forwarded_port":                           utils.PathSearch("X-Forwarded-Port", insertHeaders, nil),
		"forwarded_for_port":                       utils.PathSearch("X-Forwarded-For-Port", insertHeaders, nil),
		"forwarded_host":                           utils.PathSearch("X-Forwarded-Host", insertHeaders, nil),
		"forwarded_proto":                          utils.PathSearch("X-Forwarded-Proto", insertHeaders, nil),
		"real_ip":                                  utils.PathSearch("X-Real-IP", insertHeaders, nil),
		"forwarded_elb_id":                         utils.PathSearch("X-Forwarded-ELB-ID", insertHeaders, nil),
		"forwarded_tls_certificate_id":             utils.PathSearch("X-Forwarded-TLS-Certificate-ID", insertHeaders, nil),
		"forwarded_tls_protocol":                   utils.PathSearch("X-Forwarded-TLS-Protocol", insertHeaders, nil),
		"forwarded_tls_cipher":                     utils.PathSearch("X-Forwarded-TLS-Cipher", insertHeaders, nil),
		"forwarded_tls_protocol_alias":             utils.PathSearch("X-Forwarded-TLS-Protocol-alias", insertHeaders, nil),
		"forwarded_tls_cipher_alias":               utils.PathSearch("X-Forwarded-TLS-Cipher-alias", insertHeaders, nil),
		"forwarded_for_processing_mode":            utils.PathSearch("X-Forwarded-For-Processing-Mode", insertHeaders, nil),
		"forwarded_clientcert_subjectdn_enable":    utils.PathSearch("X-Forwarded-Clientcert-subjectdn-enable", insertHeaders, nil),
		"forwarded_clientcert_subjectdn_alias":     utils.PathSearch("X-Forwarded-Clientcert-subjectdn-alias", insertHeaders, nil),
		"forwarded_clientcert_issuerdn_enable":     utils.PathSearch("X-Forwarded-Clientcert-issuerdn-enable", insertHeaders, nil),
		"forwarded_clientcert_issuerdn_alias":      utils.PathSearch("X-Forwarded-Clientcert-issuerdn-alias", insertHeaders, nil),
		"forwarded_clientcert_fingerprint_enable":  utils.PathSearch("X-Forwarded-Clientcert-fingerprint-enable", insertHeaders, nil),
		"forwarded_clientcert_fingerprint_alias":   utils.PathSearch("X-Forwarded-Clientcert-fingerprint-alias", insertHeaders, nil),
		"forwarded_clientcert_clientverify_enable": utils.PathSearch("X-Forwarded-Clientcert-clientverify-enable", insertHeaders, nil),
		"forwarded_clientcert_clientverify_alias":  utils.PathSearch("X-Forwarded-Clientcert-clientverify-alias", insertHeaders, nil),
		"forwarded_clientcert_serialnumber_enable": utils.PathSearch("X-Forwarded-Clientcert-serialnumber-enable", insertHeaders, nil),
		"forwarded_clientcert_serialnumber_alias":  utils.PathSearch("X-Forwarded-Clientcert-serialnumber-alias", insertHeaders, nil),
		"forwarded_clientcert_enable":              utils.PathSearch("X-Forwarded-Clientcert-enable", insertHeaders, nil),
		"forwarded_clientcert_alias":               utils.PathSearch("X-Forwarded-Clientcert-alias", insertHeaders, nil),
		"forwarded_clientcert_ciphers_enable":      utils.PathSearch("X-Forwarded-Clientcert-ciphers-enable", insertHeaders, nil),
		"forwarded_clientcert_ciphers_alias":       utils.PathSearch("X-Forwarded-Clientcert-ciphers-alias", insertHeaders, nil),
		"forwarded_clientcert_end_enable":          utils.PathSearch("X-Forwarded-Clientcert-end-enable", insertHeaders, nil),
		"forwarded_clientcert_end_alias":           utils.PathSearch("X-Forwarded-Clientcert-end-alias", insertHeaders, nil),
		"forwarded_tls_alpn_protocol_enable":       utils.PathSearch("X-Forwarded-Tls-Alpn-Protocol-enable", insertHeaders, nil),
		"forwarded_tls_alpn_protocol_alias":        utils.PathSearch("X-Forwarded-Tls-Alpn-Protocol-alias", insertHeaders, nil),
		"forwarded_tls_sni_enable":                 utils.PathSearch("X-Forwarded-Tls-Sni-enable", insertHeaders, nil),
		"forwarded_tls_sni_alias":                  utils.PathSearch("X-Forwarded-Tls-Sni-alias", insertHeaders, nil),
		"forwarded_tls_ja3_enable":                 utils.PathSearch("X-Forwarded-Tls-Ja3-enable", insertHeaders, nil),
		"forwarded_tls_ja3_alias":                  utils.PathSearch("X-Forwarded-Tls-Ja3-alias", insertHeaders, nil),
		"forwarded_tls_ja4_enable":                 utils.PathSearch("X-Forwarded-Tls-Ja4-enable", insertHeaders, nil),
		"forwarded_tls_ja4_alias":                  utils.PathSearch("X-Forwarded-Tls-Ja4-alias", insertHeaders, nil),
	}

	return []map[string]interface{}{result}
}

func flattenListenerIpGroup(ipGroup interface{}) []map[string]interface{} {
	if ipGroup == nil {
		return nil
	}

	result := map[string]interface{}{
		"ipgroup_id":     utils.PathSearch("ipgroup_id", ipGroup, nil),
		"enable_ipgroup": utils.PathSearch("enable_ipgroup", ipGroup, nil),
		"type":           utils.PathSearch("type", ipGroup, nil),
	}

	return []map[string]interface{}{result}
}

func flattenListenerQuicConfig(quicConfig interface{}) []map[string]interface{} {
	if quicConfig == nil {
		return nil
	}

	result := map[string]interface{}{
		"quic_listener_id":    utils.PathSearch("quic_listener_id", quicConfig, nil),
		"enable_quic_upgrade": utils.PathSearch("enable_quic_upgrade", quicConfig, nil),
	}

	return []map[string]interface{}{result}
}

func resourceListenerCopyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceListenerCopyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	var httpUrl string
	if d.Get("force_delete").(bool) {
		httpUrl = "v3/{project_id}/elb/listeners/{listener_id}/force"
	} else {
		httpUrl = "v3/{project_id}/elb/listeners/{listener_id}"
	}

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{listener_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the listener does not exist, the response HTTP status code of the deletion API is `404`.
		return common.CheckDeletedDiag(d, err, "error deleting listener")
	}

	return nil
}
