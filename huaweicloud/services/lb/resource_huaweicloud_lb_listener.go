package lb

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
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

var listenerNonUpdatableParams = []string{"protocol", "protocol_port", "loadbalancer_id", "tenant_id"}

// @API ELB POST /v2/{project_id}/elb/listeners
// @API ELB GET /v2/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB POST /v2.0/{project_id}/listeners/{listener_id}/tags/action
// @API ELB GET /v2/{project_id}/elb/listeners/{listener_id}
// @API ELB GET /v2.0/{project_id}/listeners/{listener_id}/tags
// @API ELB PUT /v2/{project_id}/elb/listeners/{listener_id}
// @API ELB DELETE /v2/{project_id}/elb/listeners/{listener_id}
func ResourceListener() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListenerV2Create,
		ReadContext:   resourceListenerV2Read,
		UpdateContext: resourceListenerV2Update,
		DeleteContext: resourceListenerV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(listenerNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
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
			},
			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"client_ca_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sni_container_refs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"insert_headers": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem:     listenerInsertHeadersSchema(),
			},
			"tls_ciphers_policy": {
				Type:     schema.TypeString,
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
			"tags": common.TagsSchema(),
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
			},
			"admin_state_up": {
				Type:       schema.TypeBool,
				Default:    true,
				Optional:   true,
				Deprecated: "admin_state_up is deprecated",
			},
			"connection_limit": {
				Type:       schema.TypeInt,
				Optional:   true,
				Computed:   true,
				Deprecated: "connection_limit is deprecated",
			},
		},
	}
}

func listenerInsertHeadersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"x_forwarded_elb_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"x_forwarded_host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
		},
	}
	return &sc
}

func resourceListenerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/elb/listeners"
		product = "elbv2"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateListenerBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB listener: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error retrieving ELB listener: %s", err)
	}
	listenerId := utils.PathSearch("listener.id", createRespBody, "").(string)
	if listenerId == "" {
		return diag.Errorf("error creating ELB listener: ID is not found in API response")
	}

	d.SetId(listenerId)

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "listeners", listenerId, tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of ELB listener %s: %s", listenerId, tagErr)
		}
	}

	return resourceListenerV2Read(ctx, d, meta)
}

func buildCreateListenerBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"protocol":                    d.Get("protocol"),
		"protocol_port":               d.Get("protocol_port"),
		"loadbalancer_id":             d.Get("loadbalancer_id"),
		"tenant_id":                   utils.ValueIgnoreEmpty(d.Get("tenant_id")),
		"admin_state_up":              utils.ValueIgnoreEmpty(d.Get("admin_state_up")),
		"name":                        utils.ValueIgnoreEmpty(d.Get("name")),
		"description":                 utils.ValueIgnoreEmpty(d.Get("description")),
		"http2_enable":                utils.ValueIgnoreEmpty(d.Get("http2_enable")),
		"default_pool_id":             utils.ValueIgnoreEmpty(d.Get("default_pool_id")),
		"default_tls_container_ref":   utils.ValueIgnoreEmpty(d.Get("default_tls_container_ref")),
		"client_ca_tls_container_ref": utils.ValueIgnoreEmpty(d.Get("client_ca_tls_container_ref")),
		"sni_container_refs":          d.Get("sni_container_refs").(*schema.Set).List(),
		"insert_headers":              buildListenerInsertHeaders(d),
		"tls_ciphers_policy":          utils.ValueIgnoreEmpty(d.Get("tls_ciphers_policy")),
		"protection_status":           utils.ValueIgnoreEmpty(d.Get("protection_status")),
		"protection_reason":           utils.ValueIgnoreEmpty(d.Get("protection_reason")),
	}
	bodyParams := map[string]interface{}{
		"listener": params,
	}
	return bodyParams
}

func buildListenerInsertHeaders(d *schema.ResourceData) map[string]interface{} {
	if rawInsertHeaders, ok := d.GetOk("insert_headers"); ok {
		if v, ok := rawInsertHeaders.([]interface{})[0].(map[string]interface{}); ok {
			xForwardedElbIp, _ := strconv.ParseBool(v["x_forwarded_elb_ip"].(string))
			xForwardedHost, _ := strconv.ParseBool(v["x_forwarded_host"].(string))
			params := map[string]interface{}{
				"X-Forwarded-ELB-IP": xForwardedElbIp,
				"X-Forwarded-Host":   xForwardedHost,
			}
			return params
		}
	}
	return nil
}

func resourceListenerV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/{project_id}/elb/listeners/{listener_id}"
		product = "elbv2"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{listener_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB listener")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("listener.name", getRespBody, nil)),
		d.Set("protocol", utils.PathSearch("listener.protocol", getRespBody, nil)),
		d.Set("protocol_port", utils.PathSearch("listener.protocol_port", getRespBody, nil)),
		d.Set("loadbalancer_id", utils.PathSearch("listener.loadbalancers[0].id", getRespBody, nil)),
		d.Set("description", utils.PathSearch("listener.description", getRespBody, nil)),
		d.Set("default_pool_id", utils.PathSearch("listener.default_pool_id", getRespBody, nil)),
		d.Set("http2_enable", utils.PathSearch("listener.http2_enable", getRespBody, nil)),
		d.Set("default_tls_container_ref", utils.PathSearch("listener.default_tls_container_ref", getRespBody, nil)),
		d.Set("client_ca_tls_container_ref", utils.PathSearch("listener.client_ca_tls_container_ref", getRespBody, nil)),
		d.Set("sni_container_refs", utils.PathSearch("listener.sni_container_refs", getRespBody, nil)),
		d.Set("insert_headers", flattenInsertHeaders(getRespBody)),
		d.Set("tls_ciphers_policy", utils.PathSearch("listener.tls_ciphers_policy", getRespBody, nil)),
		d.Set("protection_status", utils.PathSearch("listener.protection_status", getRespBody, nil)),
		d.Set("protection_reason", utils.PathSearch("listener.protection_reason", getRespBody, nil)),
		d.Set("tenant_id", utils.PathSearch("listener.tenant_id", getRespBody, nil)),
		d.Set("admin_state_up", utils.PathSearch("listener.admin_state_up", getRespBody, nil)),
		d.Set("connection_limit", utils.PathSearch("listener.connection_limit", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("listener.created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("listener.updated_at", getRespBody, nil)),
	)

	// fetch tags
	if resourceTags, err := tags.Get(client, "listeners", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		log.Printf("[WARN] fetching tags of ELB listener failed: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInsertHeaders(listener interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("listener.insert_headers", listener, nil)
	if curJson == nil {
		return nil
	}

	xForwardedElbIp := utils.PathSearch(`"X-Forwarded-ELB-IP"`, curJson, false).(bool)
	xForwardedHost := utils.PathSearch(`"X-Forwarded-Host"`, curJson, false).(bool)
	rst := []map[string]interface{}{
		{
			"x_forwarded_elb_ip": strconv.FormatBool(xForwardedElbIp),
			"x_forwarded_host":   strconv.FormatBool(xForwardedHost),
		},
	}
	return rst
}

func resourceListenerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/elb/listeners/{listener_id}"
		product = "elbv2"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{listener_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateListenerBodyParams(d))
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating ELB listener: %s", err)
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "listeners", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of ELB listener:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceListenerV2Read(ctx, d, meta)
}

func buildUpdateListenerBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"admin_state_up":              d.Get("admin_state_up"),
		"name":                        d.Get("name"),
		"description":                 d.Get("description"),
		"http2_enable":                d.Get("http2_enable"),
		"default_pool_id":             utils.ValueIgnoreEmpty(d.Get("default_pool_id")),
		"default_tls_container_ref":   utils.ValueIgnoreEmpty(d.Get("default_tls_container_ref")),
		"client_ca_tls_container_ref": utils.ValueIgnoreEmpty(d.Get("client_ca_tls_container_ref")),
		"sni_container_refs":          d.Get("sni_container_refs").(*schema.Set).List(),
		"insert_headers":              buildListenerInsertHeaders(d),
		"tls_ciphers_policy":          utils.ValueIgnoreEmpty(d.Get("tls_ciphers_policy")),
		"protection_status":           utils.ValueIgnoreEmpty(d.Get("protection_status")),
		"protection_reason":           d.Get("protection_reason"),
	}
	bodyParams := map[string]interface{}{
		"listener": params,
	}
	return bodyParams
}

func resourceListenerV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/elb/listeners/{listener_id}"
		product = "elb"
	)

	client, err := cfg.NewServiceClient(product, region)
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
		return common.CheckDeletedDiag(d, err, "error deleting ELB listener")
	}

	return nil
}
