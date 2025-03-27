package cci

import (
	"context"
	"fmt"
	"log"
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

var networkNonUpdatableParams = []string{"namespace", "name", "ip_families", "subnets", "subnets.*.subnet_id"}

// @API CCI POST /apis/yangtse/v2/namespaces/{namespace}/networks
// @API CCI GET /apis/yangtse/v2/namespaces/{namespace}/networks/{name}
// @API CCI PUT /apis/yangtse/v2/namespaces/{namespace}/networks/{name}
// @API CCI DELETE /apis/yangtse/v2/namespaces/{namespace}/networks/{name}
func ResourceV2Network() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NetworkCreate,
		ReadContext:   resourceV2NetworkRead,
		UpdateContext: resourceV2NetworkUpdate,
		DeleteContext: resourceV2NetworkDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2NetworkImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(networkNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the CCI network.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the annotations of the CCI network.`,
			},
			"ip_families": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the IP families of the CCI network.`,
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the security group IDs of the CCI network.`,
			},
			"subnets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the subnet ID of the CCI network.`,
						},
					},
				},
				Description: `Specifies the subnets of the CCI network.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the CCI network.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the CCI network.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the namespace.`,
			},
			"finalizers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The finalizers of the namespace.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the namespace.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the namespace.`,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the CCI network.`,
						},
						"conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the CCI network conditions.`,
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Tthe status of the CCI network conditions.`,
									},
									"last_transition_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The last transition time of the CCI network conditions.`,
									},
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The reason of the CCI network conditions.`,
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The message of the CCI network conditions.`,
									},
								},
							},
							Description: `Tthe conditions of the CCI network.`,
						},
						"subnet_attrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the CCI network.`,
									},
									"subnet_v4_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The subnet IPv4 ID of the CCI network.`,
									},
									"subnet_v6_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The subnet IPv6 ID of the CCI network.`,
									},
								},
							},
							Description: `The subnet attributes of the CCI network.`,
						},
					},
				},
				Description: `The status of the namespace.`,
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

func resourceV2NetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks"
	createNetworkPath := client.Endpoint + createNetworkHttpUrl
	createNetworkPath = strings.ReplaceAll(createNetworkPath, "{namespace}", d.Get("namespace").(string))
	createNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createNetworkOpt.JSONBody = utils.RemoveNil(buildCreateV2NetworkParams(d))

	resp, err := client.Request("POST", createNetworkPath, &createNetworkOpt)
	if err != nil {
		return diag.Errorf("error creating CCI Network: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" || name == "" {
		return diag.Errorf("unable to find CCI Network name or namespace from API response")
	}
	d.SetId(ns + "/" + name)

	err = waitForCreateV2NetworkStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceV2NetworkRead(ctx, d, meta)
}

func buildCreateV2NetworkParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":        d.Get("name"),
			"namespace":   d.Get("namespace"),
			"annotations": d.Get("annotations"),
		},
		"spec": map[string]interface{}{
			"ipFamilies":     d.Get("ip_families"),
			"networkType":    "underlay_neutron",
			"securityGroups": d.Get("security_group_ids"),
			"subnets":        buildV2NetworkSubnetsParams(d.Get("subnets").([]interface{})),
		},
	}

	return bodyParams
}

func resourceV2NetworkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	resp, err := GetNetwork(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 network")
	}

	mErr := multierror.Append(
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("finalizers", utils.PathSearch("metadata.finalizers", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("ip_families", utils.PathSearch("spec.ipFamilies", resp, nil)),
		d.Set("security_group_ids", utils.PathSearch("spec.securityGroups", resp, nil)),
		d.Set("subnets", flattenNetworkSubnets(utils.PathSearch("spec.subnets", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", flattenNetworkStatus(utils.PathSearch("status", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNetworkStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"status":       utils.PathSearch("status", status, nil),
			"subnet_attrs": flattenNetworkStatusSubnetAttrs(utils.PathSearch("subnetAttrs", status, make([]interface{}, 0)).([]interface{})),
			"conditions":   flattenNetworkStatusConditions(utils.PathSearch("conditions", status, make([]interface{}, 0)).([]interface{})),
		},
	}

	return rst
}

func flattenNetworkStatusConditions(conditions []interface{}) []interface{} {
	if len(conditions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(conditions))
	for i, v := range conditions {
		rst[i] = map[string]interface{}{
			"type":                 utils.PathSearch("type", v, nil),
			"status":               utils.PathSearch("status", v, nil),
			"last_transition_time": utils.PathSearch("lastTransitionTime", v, nil),
			"reason":               utils.PathSearch("reason", v, nil),
			"message":              utils.PathSearch("message", v, nil),
		}
	}

	return rst
}

func flattenNetworkStatusSubnetAttrs(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, len(resp))
	for i, v := range resp {
		rst[i] = map[string]interface{}{
			"network_id":   utils.PathSearch("networkID", v, nil),
			"subnet_v4_id": utils.PathSearch("subnetV4ID", v, nil),
			"subnet_v6_id": utils.PathSearch("subnetV6ID", v, nil),
		}
	}

	return rst
}

func flattenNetworkSubnets(subnets []interface{}) []interface{} {
	if len(subnets) == 0 {
		return nil
	}

	rst := make([]interface{}, len(subnets))
	for i, v := range subnets {
		rst[i] = map[string]interface{}{
			"subnet_id": utils.PathSearch("subnetID", v, nil),
		}
	}

	return rst
}

func resourceV2NetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	updateNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks/{name}"
	updateNetworkPath := client.Endpoint + updateNetworkHttpUrl
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{namespace}", d.Get("namespace").(string))
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{name}", d.Get("name").(string))
	updateNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateNetworkOpt.JSONBody = utils.RemoveNil(buildUpdateV2NetworkParams(d))

	_, err = client.Request("PUT", updateNetworkPath, &updateNetworkOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 Network: %s", err)
	}
	return resourceV2NetworkRead(ctx, d, meta)
}

func buildUpdateV2NetworkParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       d.Get("kind"),
		"apiVersion": d.Get("api_version"),
		"metadata": map[string]interface{}{
			"name":              d.Get("name"),
			"namespace":         d.Get("namespace"),
			"uid":               d.Get("uid"),
			"resourceVersion":   d.Get("resource_version"),
			"creationTimestamp": d.Get("creation_timestamp"),
			"annotations":       d.Get("annotations"),
			"finalizers":        d.Get("finalizers"),
		},
		"spec": map[string]interface{}{
			"ipFamilies":     utils.ValueIgnoreEmpty(d.Get("ip_families")),
			"networkType":    "underlay_neutron",
			"securityGroups": utils.ValueIgnoreEmpty(d.Get("security_group_ids")),
			"subnets":        utils.ValueIgnoreEmpty(buildV2NetworkSubnetsParams(d.Get("subnets").([]interface{}))),
		},
		"status": map[string]interface{}{
			"status":      d.Get("status.0.status"),
			"conditions":  buildNetworkStatusConditions(d.Get("status.0.conditions").([]interface{})),
			"subnetAttrs": buildNetworkStatusSubnetAttrs(d.Get("status.0.subnet_attrs").([]interface{})),
		},
	}

	return bodyParams
}

func buildNetworkStatusConditions(conditions []interface{}) []interface{} {
	if len(conditions) == 0 {
		return nil
	}

	params := make([]interface{}, len(conditions))
	for i, v := range conditions {
		params[i] = map[string]interface{}{
			"type":               utils.PathSearch("type", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"lastTransitionTime": utils.PathSearch("last_transition_time", v, nil),
			"reason":             utils.PathSearch("reason", v, nil),
			"message":            utils.PathSearch("message", v, nil),
		}
	}

	return params
}

func buildNetworkStatusSubnetAttrs(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	params := make([]interface{}, len(resp))
	for i, v := range resp {
		params[i] = map[string]interface{}{
			"networkID":  utils.PathSearch("network_id", v, nil),
			"subnetV4ID": utils.PathSearch("subnet_v4_id", v, nil),
			"subnetV6ID": utils.ValueIgnoreEmpty(utils.PathSearch("subnet_v6_id", v, nil)),
		}
	}

	return params
}

func buildV2NetworkSubnetsParams(subnets []interface{}) []interface{} {
	if len(subnets) == 0 {
		return nil
	}

	params := make([]interface{}, len(subnets))
	for i, v := range subnets {
		params[i] = map[string]interface{}{
			"subnetID": utils.PathSearch("subnet_id", v, nil),
		}
	}

	return params
}

func resourceV2NetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	deleteNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks/{name}"
	deleteNetworkPath := client.Endpoint + deleteNetworkHttpUrl
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{namespace}", ns)
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{name}", name)
	deleteNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteNetworkPath, &deleteNetworkOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI v2 network: %s", err)
	}

	err = waitForDeleteV2NetworkStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForCreateV2NetworkStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Ready"},
		Refresh:      refreshCreateV2NetworkStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI network to complete: %s", err)
	}
	return nil
}

func refreshCreateV2NetworkStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetNetwork(client, ns, name)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status.status", resp, "").(string)
		if status != "Ready" {
			return resp, "Pending", nil
		}

		return resp, status, nil
	}
}

func waitForDeleteV2NetworkStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Deleted"},
		Refresh:      refreshDeleteV2NetworkStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI v2 network to complete: %s", err)
	}
	return nil
}

func refreshDeleteV2NetworkStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetNetwork(client, ns, name)
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[DEBUG] successfully deleted CCI network: %s", name)
			return "", "Deleted", nil
		}
		return resp, "Pending", nil
	}
}

func resourceV2NetworkImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<namespace>/<name>', but '%s'", importedId)
	}

	d.Set("namespace", parts[0])
	d.Set("name", parts[1])

	return []*schema.ResourceData{d}, nil
}

func GetNetwork(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks/{name}"
	getNetworkPath := client.Endpoint + getNetworkHttpUrl
	getNetworkPath = strings.ReplaceAll(getNetworkPath, "{namespace}", namespace)
	getNetworkPath = strings.ReplaceAll(getNetworkPath, "{name}", name)
	getNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getNetworkResp, err := client.Request("GET", getNetworkPath, &getNetworkOpt)
	if err != nil {
		return getNetworkResp, err
	}

	return utils.FlattenResponse(getNetworkResp)
}
