package css

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/loadbalancers/es-switch
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/es-listeners
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/es-listeners
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/es-listeners/{listener_id}
func ResourceEsLoadbalancerConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEsLoadbalancerConfigCreate,
		ReadContext:   resourceEsLoadbalancerConfigRead,
		UpdateContext: resourceEsLoadbalancerConfigUpdate,
		DeleteContext: resourceEsLoadbalancerConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

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
				ForceNew: true,
			},
			"agency": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"server_cert_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ca_cert_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_cert_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_cert_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"elb_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"authentication_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"loadbalancer": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadbalancerSchema(),
			},
			"listener": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     listenerSchema(),
			},
			"health_monitors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     healthmonitorsSchema(),
			},
		},
	}
}

func loadbalancerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func listenerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ip_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ipGroupSchema(),
			},
		},
	}
}

func ipGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func healthmonitorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEsLoadbalancerConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clusterId := d.Get("cluster_id").(string)
	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	err = openOrCloseLoadbalancer(client, d, true)
	if err != nil {
		return diag.Errorf("error opening CSS loadbalancer: %s", err)
	}

	d.SetId(clusterId)

	createEsListenerHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/es-listeners"
	createEsListenerPath := client.Endpoint + createEsListenerHttpUrl
	createEsListenerPath = strings.ReplaceAll(createEsListenerPath, "{project_id}", client.ProjectID)
	createEsListenerPath = strings.ReplaceAll(createEsListenerPath, "{cluster_id}", clusterId)

	createEsListenerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createEsListenerOpt.JSONBody = utils.RemoveNil(buildCreateEsListenerBodyParams(d))
	_, err = client.Request("POST", createEsListenerPath, &createEsListenerOpt)
	if err != nil {
		return diag.Errorf("error creating CSS es-listener: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"FINISHED"},
		Refresh:      ltCreateRstStateRefreshFunc(d.Id(), client),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the cluster creating es-listener to finish: %s", err)
	}

	return resourceEsLoadbalancerConfigRead(ctx, d, meta)
}

func buildCreateEsListenerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"protocol_port":  d.Get("protocol_port").(int),
		"server_cert_id": utils.ValueIgnoreEmpty(d.Get("server_cert_id")),
		"ca_cert_id":     utils.ValueIgnoreEmpty(d.Get("ca_cert_id")),
		"protocol":       "HTTP",
	}
	if _, ok := d.GetOk("server_cert_id"); ok {
		bodyParams["protocol"] = "HTTPS"
	}

	return bodyParams
}

func resourceEsLoadbalancerConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	getEsListenerRespBody, err := getEsListenerDetail(d.Id(), client)
	if err != nil {
		// The cluster does not exist, http code is 403, key/value of error code is errCode/CSS.0015.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error getting CSS es-listener")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("cluster_id", d.Id()),
		d.Set("loadbalancer_id", utils.PathSearch("loadBalancer.id", getEsListenerRespBody, nil)),
		d.Set("agency", utils.PathSearch("agency", getEsListenerRespBody, nil)),
		d.Set("protocol_port",
			int(utils.PathSearch("listener.protocol_port", getEsListenerRespBody, float64(0)).(float64))),
		d.Set("loadbalancer", flattenLoadbalancerResponse(getEsListenerRespBody)),
		d.Set("listener", flattenListenerResponse(getEsListenerRespBody)),
		d.Set("health_monitors", flattenHealthMonitorsResponse(getEsListenerRespBody)),
		d.Set("server_cert_id", utils.PathSearch("serverCertId", getEsListenerRespBody, nil)),
		d.Set("server_cert_name", utils.PathSearch("serverCertName", getEsListenerRespBody, nil)),
		d.Set("ca_cert_id", utils.PathSearch("cacertId", getEsListenerRespBody, nil)),
		d.Set("ca_cert_name", utils.PathSearch("cacertName", getEsListenerRespBody, nil)),
		d.Set("elb_enabled", utils.PathSearch("elb_enable", getEsListenerRespBody, false)),
		d.Set("authentication_type", utils.PathSearch("authentication_type", getEsListenerRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLoadbalancerResponse(getEsListenerRespBody interface{}) interface{} {
	v := utils.PathSearch("loadBalancer", getEsListenerRespBody, nil)
	result := map[string]interface{}{
		"id":        utils.PathSearch("id", v, nil),
		"name":      utils.PathSearch("name", v, nil),
		"ip":        utils.PathSearch("vip_address", v, nil),
		"public_ip": utils.PathSearch("publicips[0].publicip_address", v, nil),
	}
	return []interface{}{result}
}

func flattenListenerResponse(resp interface{}) interface{} {
	v := utils.PathSearch("listener", resp, nil)
	result := map[string]interface{}{
		"id":            utils.PathSearch("id", v, nil),
		"name":          utils.PathSearch("name", v, nil),
		"protocol":      utils.PathSearch("protocol", v, nil),
		"protocol_port": int(utils.PathSearch("protocol_port", v, float64(0)).(float64)),
		"ip_group":      flattenIpGroupResponse(v),
	}
	return []interface{}{result}
}

func flattenIpGroupResponse(resp interface{}) interface{} {
	v := utils.PathSearch("ipgroup", resp, nil)
	result := map[string]interface{}{
		"id":      utils.PathSearch("ipgroup_id", v, nil),
		"enabled": utils.PathSearch("enable_ipgroup", v, false),
	}
	return []interface{}{result}
}

func flattenHealthMonitorsResponse(resp interface{}) []interface{} {
	curJson := utils.PathSearch("healthmonitors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		result := map[string]interface{}{
			"ip":            utils.PathSearch("address", v, nil),
			"protocol_port": int(utils.PathSearch("protocol_port", v, float64(0)).(float64)),
			"status":        utils.PathSearch("operating_status", v, nil),
		}
		rst = append(rst, result)
	}

	return rst
}

func resourceEsLoadbalancerConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	createEsListenerHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/es-listeners/{listener_id}"
	createEsListenerPath := client.Endpoint + createEsListenerHttpUrl
	createEsListenerPath = strings.ReplaceAll(createEsListenerPath, "{project_id}", client.ProjectID)
	createEsListenerPath = strings.ReplaceAll(createEsListenerPath, "{cluster_id}", d.Id())
	createEsListenerPath = strings.ReplaceAll(createEsListenerPath, "{listener_id}", d.Get("listener.0.id").(string))

	createEsListenerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createEsListenerOpt.JSONBody = utils.RemoveNil(buildUpdateEsListenerBodyParams(d))
	_, err = client.Request("PUT", createEsListenerPath, &createEsListenerOpt)
	if err != nil {
		return diag.Errorf("error updating CSS es-listener: %s", err)
	}

	return resourceEsLoadbalancerConfigRead(ctx, d, meta)
}

func buildUpdateEsListenerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"listener": map[string]interface{}{
			"default_tls_container_ref":   d.Get("server_cert_id"),
			"client_ca_tls_container_ref": d.Get("ca_cert_id"),
		},
	}

	return bodyParams
}

func resourceEsLoadbalancerConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	err = openOrCloseLoadbalancer(client, d, false)
	if err != nil {
		// The elb is already disabled, http code is 400, key/value of error code is errCode/CSS.0001.
		err = common.ConvertExpected400ErrInto404Err(err, "errCode", "CSS.0001")
		// The cluster does not exist, http code is 403, key/value of error code is errCode/CSS.0015.
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error closing CSS loadbalancer")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"FINISHED"},
		Refresh:      ltDeleteRstStateRefreshFunc(d.Id(), client),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the cluster deleting es-listener to finish: %s", err)
	}

	return nil
}

func ltCreateRstStateRefreshFunc(clusterID string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getEsListenerDetail(clusterID, client)
		if err != nil {
			return nil, "ERROR", err
		}

		listener := utils.PathSearch("listener", resp, nil)
		if listener != nil {
			return resp, "FINISHED", nil
		}

		return resp, "PENDING", nil
	}
}

func ltDeleteRstStateRefreshFunc(clusterID string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getEsListenerDetail(clusterID, client)
		if err != nil {
			return nil, "ERROR", err
		}

		elbEnabled := utils.PathSearch("elb_enable", resp, false).(bool)
		listencer := utils.PathSearch("listener", resp, nil)
		if !elbEnabled && listencer == nil {
			return resp, "FINISHED", nil
		}

		return resp, "PENDING", nil
	}
}

func openOrCloseLoadbalancer(client *golangsdk.ServiceClient, d *schema.ResourceData, isOpen bool) error {
	openLoadbalancerHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/loadbalancers/es-switch"
	openLoadbalancerPath := client.Endpoint + openLoadbalancerHttpUrl
	openLoadbalancerPath = strings.ReplaceAll(openLoadbalancerPath, "{project_id}", client.ProjectID)
	openLoadbalancerPath = strings.ReplaceAll(openLoadbalancerPath, "{cluster_id}", d.Get("cluster_id").(string))

	openLoadbalancerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	openLoadbalancerOpt.JSONBody = map[string]interface{}{
		"enable": isOpen,
		"agency": d.Get("agency").(string),
		"elb_id": d.Get("loadbalancer_id").(string),
	}
	_, err := client.Request("POST", openLoadbalancerPath, &openLoadbalancerOpt)
	if err != nil {
		return err
	}

	return nil
}

func getEsListenerDetail(clusterID string, client *golangsdk.ServiceClient) (interface{}, error) {
	getEsListenerHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/es-listeners"

	getEsListenerPath := client.Endpoint + getEsListenerHttpUrl
	getEsListenerPath = strings.ReplaceAll(getEsListenerPath, "{project_id}", client.ProjectID)
	getEsListenerPath = strings.ReplaceAll(getEsListenerPath, "{cluster_id}", clusterID)

	getEsListenerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getEsListenerResp, err := client.Request("GET", getEsListenerPath, &getEsListenerOpt)
	if err != nil {
		return nil, err
	}

	getEsListenerRespBody, err := utils.FlattenResponse(getEsListenerResp)
	if err != nil {
		return nil, err
	}

	return getEsListenerRespBody, nil
}
