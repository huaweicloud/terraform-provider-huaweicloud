// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

// @API ModelArts POST /v2/{project_id}/pools
// @API ModelArts DELETE /v2/{project_id}/pools/{id}
// @API ModelArts GET /v2/{project_id}/pools/{id}
// @API ModelArts PATCH /v2/{project_id}/pools/{id}
func ResourceModelartsResourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelartsResourcePoolCreate,
		UpdateContext: resourceModelartsResourcePoolUpdate,
		ReadContext:   resourceModelartsResourcePoolRead,
		DeleteContext: resourceModelartsResourcePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the resource pool.`,
			},
			"scope": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`List of job types supported by the resource pool.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			"resources": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolResourceFlavorSchema(),
				Required:    true,
				Description: `List of resource specifications in the resource pool.`,
			},
			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"vpc_id"},
				Description:  `The ModelArts network ID of the resource pool.`,
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The prefix of the user-defined node name of the resource pool.`,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"subnet_id", "clusters", "user_login"},
				Description:  `The VPC ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The network ID of a subnet.`,
			},
			"clusters": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolClustersSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The list of the CCE clusters.`,
			},
			"user_login": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolUserLoginSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `The user login info of the resource pool.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Workspace ID, which defaults to 0.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the resource pool.`,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the resource pool.`,
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource ID of the resource pool.`,
			},
		},
	}
}

func modelartsResourcePoolResourceFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource flavor ID.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Number of resources of the corresponding flavors.`,
			},
			"node_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of resource pool nodes.`,
			},
			"max_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The max number of resources of the corresponding flavors.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The VPC ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The network ID of a subnet.`,
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The security group IDs.`,
			},
			"azs": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolResourcesAzsSchema(),
				Optional:    true,
				Description: `AZs for resource pool nodes.`,
			},
			"taints": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolResourcesTaintSchema(),
				Optional:    true,
				Description: `The taints added to nodes.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `The labels of resource pool.`,
			},
			"tags": common.TagsSchema(),
			"post_install": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The script to be executed after security.`,
			},
		},
	}
	return &sc
}

func modelartsResourcePoolResourcesTaintSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The key of the taint.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The value of the taint.`,
			},
			"effect": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The effect of the taint.`,
			},
		},
	}
	return &sc
}

func modelartsResourcePoolResourcesAzsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"az": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The AZ name.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of nodes.`,
			},
		},
	}
	return &sc
}

func modelartsResourcePoolClustersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the CCE cluster.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the CCE cluster.`,
			},
		},
	}
	return &sc
}

func modelartsResourcePoolUserLoginSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `The password of the login user.`,
			},
			"key_pair_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The key pair name of the login user.`,
			},
		},
	}
	return &sc
}

func scopeStatusRefreshFunc(cfg *config.Config, region string, d *schema.ResourceData, scopes []interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getResourcePoolRespBody, err := queryResourcePool(cfg, region, d)
		if err != nil {
			return getResourcePoolRespBody, "ERROR", err
		}

		for _, scope := range scopes {
			scopeStatus := fmt.Sprintf("status.scope[?scopeType=='%s']|[0].state", scope)
			if utils.PathSearch(scopeStatus, getResourcePoolRespBody, "").(string) != "Enabled" {
				return "No matches found", "PENDING", nil
			}
		}
		return "Matched", "COMPLETED", nil
	}
}

func createResourcePoolWaitingForScopesCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, timeout time.Duration) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: scopeStatusRefreshFunc(cfg, region, d, d.Get("scope").(*schema.Set).List()),
		Timeout: timeout,
		// In most cases, the bind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the scope statuses are both completed: %s", err)
	}
	return nil
}

func resourceModelartsResourcePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createResourcePoolHttpUrl = "v2/{project_id}/pools"
		createResourcePoolProduct = "modelarts"
	)
	createResourcePoolClient, err := cfg.NewServiceClient(createResourcePoolProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createResourcePoolPath := createResourcePoolClient.Endpoint + createResourcePoolHttpUrl
	createResourcePoolPath = strings.ReplaceAll(createResourcePoolPath, "{project_id}", createResourcePoolClient.ProjectID)

	createResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createResourcePoolOpt.JSONBody = utils.RemoveNil(buildCreateResourcePoolBodyParams(d))
	createResourcePoolResp, err := createResourcePoolClient.Request("POST", createResourcePoolPath, &createResourcePoolOpt)
	if err != nil {
		return diag.Errorf("error creating Modelarts resource pool: %s", err)
	}

	createResourcePoolRespBody, err := utils.FlattenResponse(createResourcePoolResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("metadata.name", createResourcePoolRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating Modelarts resource pool: ID is not found in API response")
	}
	d.SetId(id.(string))

	if d.Get("charging_mode") == "prePaid" {
		// wait 30 seconds so that the resource pool can be queried
		// lintignore:R018
		time.Sleep(30 * time.Second)
		resourcePool, err := queryResourcePool(cfg, region, d)
		if err != nil {
			return diag.Errorf("error retrieving Modelarts resource pool: %s", err)
		}
		orderId := utils.PathSearch(`metadata.annotations."os.modelarts/order.id"`, resourcePool, nil)
		if orderId == nil {
			return diag.Errorf("error creating Modelarts resource pool: order ID is not found in API response")
		}
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = createResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) creation to complete: %s", d.Id(), err)
	}

	err = createResourcePoolWaitingForScopesCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) creation to complete: %s", d.Id(), err)
	}
	return resourceModelartsResourcePoolRead(ctx, d, meta)
}

func buildCreateResourcePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"apiVersion": "v2",
		"kind":       "Pool",
		"metadata":   buildCreateResourcePoolMetaDataBodyParams(d),
		"spec":       buildCreateResourcePoolSpecBodyParams(d),
	}
	return bodyParams
}

func buildCreateResourcePoolMetaDataBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"labels":      buildCreateResourcePoolMetaDataLabelsBodyParams(d),
		"annotations": buildCreateResourcePoolMetaDataAnnotationsBodyParams(d),
	}
	return params
}

func buildCreateResourcePoolMetaDataLabelsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"os.modelarts/name":         d.Get("name"),
		"os.modelarts/node.prefix":  utils.ValueIgnoreEmpty(d.Get("prefix")),
		"os.modelarts/workspace.id": utils.ValueIgnoreEmpty(d.Get("workspace_id")),
	}
	return params
}

func buildCreateResourcePoolMetaDataAnnotationsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"os.modelarts/description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	if d.Get("charging_mode") == "prePaid" {
		params["os.modelarts/billing.mode"] = "1"
		if d.Get("period_unit") == "month" {
			params["os.modelarts/period.type"] = "2"
		} else {
			params["os.modelarts/period.type"] = "3"
		}
		params["os.modelarts/period.num"] = strconv.Itoa(d.Get("period").(int))
		if d.Get("auto_renew").(string) == "true" {
			params["os.modelarts/auto.renew"] = "1"
		}
		params["os.modelarts/auto.pay"] = "1"
	}
	return params
}

func buildCreateResourcePoolSpecBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"type":      "Dedicate",
		"scope":     utils.ValueIgnoreEmpty(d.Get("scope").(*schema.Set).List()),
		"resources": buildResourcePoolSpecResources(d),
		"userLogin": buildCreateResourcePoolSpecUserLoginBodyParams(d),
		"network":   buildCreateResourcePoolSpecNetworkBodyParams(d),
		"clusters":  buildCreateResourcePoolSpecClustersBodyParams(d),
	}
	return params
}

func buildCreateResourcePoolSpecUserLoginBodyParams(d *schema.ResourceData) map[string]interface{} {
	userLoginRaw := d.Get("user_login").([]interface{})
	if len(userLoginRaw) < 1 {
		return nil
	}
	return map[string]interface{}{
		"password":    utils.ValueIgnoreEmpty(d.Get("user_login.0.password")),
		"keyPairName": utils.ValueIgnoreEmpty(d.Get("user_login.0.key_pair_name")),
	}
}

func buildCreateResourcePoolSpecNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	network := make(map[string]interface{})
	if networkId, ok := d.GetOk("network_id"); ok {
		network["name"] = networkId
	} else {
		network["vpcId"] = d.Get("vpc_id")
		network["subnetId"] = d.Get("subnet_id")
	}
	return network
}

func buildCreateResourcePoolSpecClustersBodyParams(d *schema.ResourceData) []interface{} {
	if clusters, ok := d.GetOk("clusters"); ok {
		providerIDs := make([]interface{}, 0, len(clusters.([]interface{})))
		for _, clusterRaw := range clusters.([]interface{}) {
			cluster := clusterRaw.(map[string]interface{})
			providerIDs = append(providerIDs, map[string]interface{}{
				"providerId": cluster["provider_id"],
			})
		}
		return providerIDs
	}
	return nil
}

func buildResourcePoolSpecResources(d *schema.ResourceData) []map[string]interface{} {
	if rawArray, ok := d.Get("resources").([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"flavor":       utils.ValueIgnoreEmpty(raw["flavor_id"]),
					"count":        utils.ValueIgnoreEmpty(raw["count"]),
					"nodePool":     utils.ValueIgnoreEmpty(raw["node_pool"]),
					"maxCount":     utils.ValueIgnoreEmpty(raw["max_count"]),
					"azs":          buildResourcePoolResourcesAzs(raw["azs"]),
					"network":      buildResourcePoolSpecResourcesNetworkBodyParams(raw),
					"taints":       buildResourcePoolResourcesTaints(raw["taints"]),
					"tags":         utils.ExpandResourceTags(raw["tags"].(map[string]interface{})),
					"labels":       utils.ValueIgnoreEmpty(raw["labels"]),
					"extendParams": buildCreateResourcePoolSpecResourcesPostInstallBodyParams(raw),
				}
			}
		}
		return rst
	}
	return nil
}

func buildResourcePoolSpecResourcesNetworkBodyParams(rawParam map[string]interface{}) map[string]interface{} {
	if vpcId := rawParam["vpc_id"]; len(vpcId.(string)) > 0 {
		return map[string]interface{}{
			"vpc":            utils.ValueIgnoreEmpty(rawParam["vpc_id"]),
			"subnet":         utils.ValueIgnoreEmpty(rawParam["subnet_id"]),
			"securityGroups": utils.ValueIgnoreEmpty(rawParam["security_group_ids"]),
		}
	}
	return nil
}

func buildCreateResourcePoolSpecResourcesPostInstallBodyParams(rawParam map[string]interface{}) map[string]interface{} {
	if postInstall := rawParam["post_install"]; len(postInstall.(string)) > 0 {
		return map[string]interface{}{"post_install": postInstall}
	}
	return nil
}

func createResourcePoolWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)

			createResourcePoolWaitingRespBody, err := queryResourcePool(cfg, region, d)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`status.phase`, createResourcePoolWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status.phase`)
			}

			if utils.PathSearch("status.resources.abnormal", createResourcePoolWaitingRespBody, nil) != nil {
				return nil, "ERROR", fmt.Errorf("error creating resource pool: the resource pool is abnormal")
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"Running",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createResourcePoolWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"Creating",
				"Waiting",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return createResourcePoolWaitingRespBody, "PENDING", nil
			}

			return createResourcePoolWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceModelartsResourcePoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	getModelartsResourcePoolRespBody, err := queryResourcePool(cfg, region, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Modelarts resource pool")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch(`metadata.labels."os.modelarts/name"`,
			getModelartsResourcePoolRespBody, nil)),
		d.Set("prefix", utils.PathSearch(`metadata.labels."os.modelarts/node.prefix"`,
			getModelartsResourcePoolRespBody, nil)),
		d.Set("workspace_id", utils.PathSearch(`metadata.labels."os.modelarts/workspace.id"`,
			getModelartsResourcePoolRespBody, nil)),
		d.Set("description", utils.PathSearch(`metadata.annotations."os.modelarts/description"`,
			getModelartsResourcePoolRespBody, nil)),
		d.Set("scope", utils.PathSearch("spec.scope", getModelartsResourcePoolRespBody, nil)),
		d.Set("resources", flattenGetResourcePoolResponseBodyResources(getModelartsResourcePoolRespBody)),
		d.Set("network_id", utils.PathSearch("spec.network.name", getModelartsResourcePoolRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("spec.network.vpcId", getModelartsResourcePoolRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("spec.network.subnetId", getModelartsResourcePoolRespBody, nil)),
		d.Set("clusters", flattenGetResourcePoolClusterBodyResources(getModelartsResourcePoolRespBody)),
		d.Set("status", utils.PathSearch("status.phase", getModelartsResourcePoolRespBody, nil)),
		d.Set("resource_pool_id", utils.PathSearch(`metadata.labels."os.modelarts/resource.id"`,
			getModelartsResourcePoolRespBody, nil)),
	)
	chargingMode := utils.PathSearch(`metadata.annotations."os.modelarts/billing.mode"`,
		getModelartsResourcePoolRespBody, nil)
	if chargingMode == "1" {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "prePaid"))
	}

	keyPairName := utils.PathSearch("spec.userLogin.keyPairName", getModelartsResourcePoolRespBody, nil)
	if keyPairName != nil {
		rst := make([]interface{}, 0, 1)
		rst = append(rst, map[string]interface{}{
			"key_pair_name": keyPairName,
		})
		mErr = multierror.Append(mErr, d.Set("user_login", rst))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryResourcePool(cfg *config.Config, region string, d *schema.ResourceData) (interface{}, error) {
	var (
		getModelartsResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
		getModelartsResourcePoolProduct = "modelarts"
	)
	getModelartsResourcePoolClient, err := cfg.NewServiceClient(getModelartsResourcePoolProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	getModelartsResourcePoolPath := getModelartsResourcePoolClient.Endpoint + getModelartsResourcePoolHttpUrl
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{project_id}",
		getModelartsResourcePoolClient.ProjectID)
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{id}", d.Id())

	getModelartsResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getModelartsResourcePoolResp, err := getModelartsResourcePoolClient.Request("GET", getModelartsResourcePoolPath,
		&getModelartsResourcePoolOpt)

	if err != nil {
		return nil, err
	}

	getModelartsResourcePoolRespBody, err := utils.FlattenResponse(getModelartsResourcePoolResp)
	if err != nil {
		return nil, err
	}
	return getModelartsResourcePoolRespBody, nil
}

func flattenGetResourcePoolResponseBodyResources(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("spec.resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"flavor_id":          utils.PathSearch("flavor", v, nil),
			"count":              utils.PathSearch("count", v, nil),
			"max_count":          utils.PathSearch("maxCount", v, nil),
			"node_pool":          utils.PathSearch("nodePool", v, nil),
			"vpc_id":             utils.PathSearch("network.vpc", v, nil),
			"subnet_id":          utils.PathSearch("network.subnet", v, nil),
			"security_group_ids": utils.PathSearch("network.securityGroups", v, nil),
			"azs":                flattenResourcePoolResourcesAzs(v),
			"taints":             flattenResourcePoolResourcesTaints(v),
			"labels":             utils.PathSearch("labels", v, nil),
			"post_install":       utils.PathSearch("extendParams.post_install", v, nil),
			"tags":               flattenResourcePoolResourcesTags(v),
		})
	}
	return rst
}

func flattenGetResourcePoolClusterBodyResources(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("spec.clusters", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"provider_id": utils.PathSearch("providerId", v, nil),
			"name":        utils.PathSearch("name", v, nil),
		})
	}
	return rst
}

func flattenResourcePoolResourcesAzs(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("azs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"az":    utils.PathSearch("az", v, nil),
			"count": utils.PathSearch("count", v, nil),
		})
	}
	return rst
}

func flattenResourcePoolResourcesTaints(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("taints", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"value":  utils.PathSearch("value", v, nil),
			"effect": utils.PathSearch("effect", v, nil),
		})
	}
	return rst
}

func flattenResourcePoolResourcesTags(resp interface{}) map[string]interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("tags", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make(map[string]interface{})
	for _, v := range curArray {
		key := utils.PathSearch("key", v, "").(string)
		value := utils.PathSearch("value", v, "").(string)
		rst[key] = value
	}
	return rst
}

func resourceModelartsResourcePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateResourcePoolChanges := []string{
		"description",
		"scope",
		"resources",
	}

	if d.HasChanges(updateResourcePoolChanges...) {
		var (
			updateResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
			updateResourcePoolProduct = "modelarts"
		)
		updateResourcePoolClient, err := cfg.NewServiceClient(updateResourcePoolProduct, region)
		if err != nil {
			return diag.Errorf("error creating ModelArts client: %s", err)
		}

		updateResourcePoolPath := updateResourcePoolClient.Endpoint + updateResourcePoolHttpUrl
		updateResourcePoolPath = strings.ReplaceAll(updateResourcePoolPath, "{project_id}", updateResourcePoolClient.ProjectID)
		updateResourcePoolPath = strings.ReplaceAll(updateResourcePoolPath, "{id}", d.Id())

		updateResourcePoolOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/merge-patch+json"},
		}

		updateResourcePoolOpt.JSONBody = utils.RemoveNil(buildUpdateResourcePoolBodyParams(d))
		updateModelartsResourcePoolResp, err := updateResourcePoolClient.Request("PATCH", updateResourcePoolPath, &updateResourcePoolOpt)
		if err != nil {
			return diag.Errorf("error updating Modelarts resource pool: %s", err)
		}

		if d.Get("charging_mode") == "prePaid" {
			updateModelartsResourcePoolRespBody, err := utils.FlattenResponse(updateModelartsResourcePoolResp)
			if err != nil {
				return diag.FromErr(err)
			}
			orderId := utils.PathSearch(`metadata.annotations."os.modelarts/order.id"`,
				updateModelartsResourcePoolRespBody, nil)
			if orderId == nil {
				return diag.Errorf("error updating Modelarts resource pool: order ID is not found in API response")
			}
			bssClient, err := cfg.BssV2Client(region)
			if err != nil {
				return diag.Errorf("error creating BSS v2 client: %s", err)
			}
			err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = common.WaitOrderAllResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		err = updateResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Modelarts resource pool (%s) update to complete: %s", d.Id(), err)
		}
	}
	return resourceModelartsResourcePoolRead(ctx, d, meta)
}

func buildUpdateResourcePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": buildUpdateResourcePoolMetaDataBodyParams(d),
		"spec":     buildUpdateResourcePoolSpecBodyParams(d),
	}
	return bodyParams
}

func buildUpdateResourcePoolMetaDataBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"annotations": buildUpdateResourcePoolMetaDataAnnotationsBodyParams(d),
	}
	return params
}

func buildUpdateResourcePoolMetaDataAnnotationsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"os.modelarts/description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	if d.Get("charging_mode") == "prePaid" {
		params["os.modelarts/order.id"] = ""
		params["os.modelarts/auto.pay"] = "1"
	}
	return params
}

func buildUpdateResourcePoolSpecBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"scope":     utils.ValueIgnoreEmpty(d.Get("scope").(*schema.Set).List()),
		"resources": buildResourcePoolSpecResources(d),
	}
	return params
}

func buildResourcePoolResourcesAzs(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"az":    utils.ValueIgnoreEmpty(raw["az"]),
					"count": utils.ValueIgnoreEmpty(raw["count"]),
				}
			}
		}
		return rst
	}
	return nil
}

func buildResourcePoolResourcesTaints(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"key":    utils.ValueIgnoreEmpty(raw["key"]),
					"effect": utils.ValueIgnoreEmpty(raw["effect"]),
					"value":  utils.ValueIgnoreEmpty(raw["value"]),
				}
			}
		}
		return rst
	}
	return nil
}

func updateResourcePoolWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)

			getResourcePoolRespBody, err := queryResourcePool(cfg, region, d)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`status.phase`, getResourcePoolRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse 'status.phase' from response body")
			}

			status := fmt.Sprintf("%v", statusRaw)

			unexpectedStatus := []string{
				"Abnormal", "Error", "ScalingFailed", "CreationFailed",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return getResourcePoolRespBody, status, nil
			}

			// check if the resource pool is in the process of expanding capacity
			if rawArray, ok := d.Get("resources").([]interface{}); ok {
				if utils.PathSearch("status.resources.abnormal", getResourcePoolRespBody, nil) != nil {
					return nil, "ERROR", fmt.Errorf("error updating resource pool")
				}

				for _, v := range rawArray {
					raw := v.(map[string]interface{})
					flavor := utils.ValueIgnoreEmpty(raw["flavor_id"])
					count := utils.ValueIgnoreEmpty(raw["count"])

					searchActiveJsonPath := fmt.Sprintf("length(status.resources.available[?flavor=='%s' && count==`%d`])",
						flavor, count)

					log.Println("searchActiveJsonPath: ", searchActiveJsonPath)
					if utils.PathSearch(searchActiveJsonPath, getResourcePoolRespBody, float64(0)).(float64) == 0 {
						return getResourcePoolRespBody, "PENDING", nil
					}
				}
			}

			// check if the resource pool is in the process of changing scope
			if rawArray, ok := d.GetOk("scope"); ok {
				for _, v := range rawArray.(*schema.Set).List() {
					scopeStatus := fmt.Sprintf("status.scope[?scopeType=='%s']|[0].state", v)
					log.Println("scopeStatus: ", scopeStatus)
					if utils.PathSearch(scopeStatus, getResourcePoolRespBody, "").(string) != "Enabled" {
						return getResourcePoolRespBody, "PENDING", nil
					}
				}
			}

			return getResourcePoolRespBody, "COMPLETED", nil
		},
		Timeout:      t,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceModelartsResourcePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if d.Get("charging_mode").(string) == "prePaid" {
		resourcePoolId := d.Get("resource_pool_id")
		if resourcePoolId == nil {
			return diag.Errorf("error getting resource ID from the resource pool(%s)", d.Id())
		}
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{resourcePoolId.(string)}); err != nil {
			return diag.Errorf("error unsubscribing Modelarts resource pool: %s", err)
		}
	} else {
		err := deleteResourcePool(cfg, d, region)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := deleteResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteResourcePool(cfg *config.Config, d *schema.ResourceData, region string) error {
	var (
		deleteResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
		deleteResourcePoolProduct = "modelarts"
	)
	deleteResourcePoolClient, err := cfg.NewServiceClient(deleteResourcePoolProduct, region)
	if err != nil {
		return fmt.Errorf("error creating ModelArts client: %s", err)
	}

	deleteResourcePoolPath := deleteResourcePoolClient.Endpoint + deleteResourcePoolHttpUrl
	deleteResourcePoolPath = strings.ReplaceAll(deleteResourcePoolPath, "{project_id}", deleteResourcePoolClient.ProjectID)
	deleteResourcePoolPath = strings.ReplaceAll(deleteResourcePoolPath, "{id}", d.Id())

	deleteResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteResourcePoolClient.Request("DELETE", deleteResourcePoolPath, &deleteResourcePoolOpt)
	if err != nil {
		return fmt.Errorf("error deleting Modelarts resource pool: %s", err)
	}
	return nil
}

func deleteResourcePoolWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			res, err := queryResourcePool(cfg, region, d)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					res = map[string]string{"code": "COMPLETED"}
					return res, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return res, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
