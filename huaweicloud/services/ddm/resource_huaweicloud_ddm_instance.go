// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ctxType string

// @API DDM POST /v1/{project_id}/instances
// @API DDM GET /v1/{project_id}/instances/{instance_id}
// @API DDM PUT /v1/{project_id}/instances/{instance_id}/modify-name
// @API DDM PUT /v1/{project_id}/instances/{instance_id}/modify-security-group
// @API DDM GET /v2/{project_id}/flavors
// @API DDM PUT /v3/{project_id}/instances/{instance_id}/flavor
// @API DDM POST /v2/{project_id}/instances/{instance_id}/action/enlarge
// @API DDM POST /v2/{project_id}/instances/{instance_id}/action/reduce
// @API DDM PUT /v3/{project_id}/instances/{instance_id}/admin-user
// @API DDM DELETE /v1/{project_id}/instances/{instance_id}
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceDdmInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdmInstanceCreate,
		UpdateContext: resourceDdmInstanceUpdate,
		ReadContext:   resourceDdmInstanceRead,
		DeleteContext: resourceDdmInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				Description: `Specifies the name of the DDM instance.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a product.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the number of nodes.`,
			},
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of an Engine.`,
			},
			"availability_zones": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the list of availability zones.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a VPC.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a security group.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a subnet.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project id.`,
			},
			"param_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of parameter group.`,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the time zone.`,
			},
			"admin_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the username of the administrator.`,
			},
			"admin_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				Description: `Specifies the password of the administrator.`,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			// charge info: charging_mode, period_unit, period, auto_renew, auto_pay
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"delete_rds_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether data stored on the associated DB instances is deleted`,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DDM instance.`,
			},
			"access_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the address for accessing the DDM instance.`,
			},
			"access_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the port for accessing the DDM instance.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the engine version.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Elem:        InstanceNodeInfoRefSchema(),
				Computed:    true,
				Description: `Indicates the node information.`,
			},
		},
	}
}

func InstanceNodeInfoRefSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DDM instance node.`,
			},
			"port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the port of the DDM instance node.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the IP address of the DDM instance node.`,
			},
		},
	}
	return &sc
}

func resourceDdmInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createInstance: create DDM instance
	var (
		createInstanceHttpUrl = "v1/{project_id}/instances"
		createInstanceProduct = "ddm"
	)
	createInstanceClient, err := cfg.NewServiceClient(createInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	createInstancePath := createInstanceClient.Endpoint + createInstanceHttpUrl
	createInstancePath = strings.ReplaceAll(createInstancePath, "{project_id}", createInstanceClient.ProjectID)

	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createInstanceOpt.JSONBody = utils.RemoveNil(buildCreateInstanceBodyParams(d, cfg))
	createInstanceResp, err := createInstanceClient.Request("POST", createInstancePath, &createInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating DDM instance: %s", err)
	}

	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var instanceId string
	var delayTime time.Duration = 200
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		orderId := utils.PathSearch("order_id", createInstanceRespBody, "").(string)
		if orderId == "" {
			return diag.Errorf("unable to find order_id of the DDM instance from the API response")
		}
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		// wait for order success
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for replica order resource %s complete: %s", orderId, err)
		}
		instanceId = resourceId
		delayTime = 20
	} else {
		instanceId = utils.PathSearch("id", createInstanceRespBody, "").(string)
	}
	if instanceId == "" {
		return diag.Errorf("unable to find the DDM instance ID from the API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      ddmInstanceStatusRefreshFunc(instanceId, createInstanceClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        delayTime * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to running: %s", instanceId, err)
	}

	d.SetId(instanceId)

	if _, ok := d.GetOk("parameters"); ok {
		err = initializeParameters(ctx, d, createInstanceClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDdmInstanceRead(ctx, d, meta)
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance":     buildCreateInstanceInstanceChildBody(d, cfg),
		"extend_param": buildCreateInstanceExtendParamChildBody(d),
	}
	return bodyParams
}

func buildCreateInstanceInstanceChildBody(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"flavor_id":             utils.ValueIgnoreEmpty(d.Get("flavor_id")),
		"node_num":              utils.ValueIgnoreEmpty(d.Get("node_num")),
		"engine_id":             utils.ValueIgnoreEmpty(d.Get("engine_id")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"available_zones":       d.Get("availability_zones").(*schema.Set).List(), // The ordering of the AZ list returned by the API is unknown.
		"vpc_id":                utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"security_group_id":     utils.ValueIgnoreEmpty(d.Get("security_group_id")),
		"subnet_id":             utils.ValueIgnoreEmpty(d.Get("subnet_id")),
		"param_group_id":        utils.ValueIgnoreEmpty(d.Get("param_group_id")),
		"time_zone":             utils.ValueIgnoreEmpty(d.Get("time_zone")),
		"admin_user_name":       utils.ValueIgnoreEmpty(d.Get("admin_user")),
		"admin_user_password":   utils.ValueIgnoreEmpty(d.Get("admin_password")),
	}

	return params
}

func buildCreateInstanceExtendParamChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"charge_mode":   utils.ValueIgnoreEmpty(d.Get("charging_mode")),
		"period_type":   utils.ValueIgnoreEmpty(d.Get("period_unit")),
		"period_num":    utils.ValueIgnoreEmpty(d.Get("period")),
		"is_auto_renew": d.Get("auto_renew"),
		"is_auto_pay":   "true",
	}

	return params
}

func resourceDdmInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Id()

	var (
		updateInstanceProduct = "ddm"
	)
	updateClient, err := cfg.NewServiceClient(updateInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	if d.HasChange("name") {
		err = updateInstanceName(ctx, d, updateClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("security_group_id") {
		err = updateInstanceSecurityGroup(ctx, d, updateClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("flavor_id") {
		err = updateInstanceFlavor(ctx, d, updateClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("node_num") {
		err = updateInstanceNodeNum(ctx, d, updateClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("admin_password") {
		err = updateInstanceAdminPassword(ctx, d, updateClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("parameters") {
		ctx, err = updateInstanceParameters(ctx, d, updateClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), instanceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the DDM instance (%s): %s", instanceId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "ddm",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDdmInstanceRead(ctx, d, meta)
}

func updateInstanceName(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// updateInstanceName: update DDM instance name
	var (
		httpUrl = "v1/{project_id}/instances/{instance_id}/modify-name"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateInstanceNameBodyParams(d)
	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating DDM instance name: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      ddmInstanceStatusRefreshFunc(d.Id(), client),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to running: %s", d.Id(), err)
	}
	return nil
}

func updateInstanceSecurityGroup(_ context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// updateInstanceSecurityGroup: update DDM instance security group
	var (
		httpUrl = "v1/{project_id}/instances/{instance_id}/modify-security-group"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateInstanceSecurityGroupBodyParams(d)
	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating DDM instance security group: %s", err)
	}
	return nil
}

func updateInstanceNodeNum(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// updateInstanceNodeNum: update DDM instance node num
	var (
		updateInstanceNodeEnlargeNumHttpUrl = "v2/{project_id}/instances/{instance_id}/action/enlarge"
		updateInstanceNodeReduceNumHttpUrl  = "v2/{project_id}/instances/{instance_id}/action/reduce"
	)

	var httpUrl string
	var nodeNumber int
	oldNodeNumRaw, newNodeNumRaw := d.GetChange("node_num")
	oldNodeNum := oldNodeNumRaw.(int)
	newNodeNum := newNodeNumRaw.(int)

	if oldNodeNum < newNodeNum {
		httpUrl = updateInstanceNodeEnlargeNumHttpUrl
		nodeNumber = newNodeNum - oldNodeNum
	} else {
		httpUrl = updateInstanceNodeReduceNumHttpUrl
		nodeNumber = oldNodeNum - newNodeNum
	}
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateInstanceNodeNumBodyParams(d, nodeNumber)
	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating DDM instance node number: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      ddmInstanceStatusRefreshFunc(d.Id(), client),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        100 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to running: %s", d.Id(), err)
	}
	return nil
}

func updateInstanceFlavor(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// updateInstanceFlavor: update DDM instance flavor
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/flavor"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	flavorId := utils.ValueIgnoreEmpty(d.Get("flavor_id"))
	engineId := utils.ValueIgnoreEmpty(d.Get("engine_id"))
	specCode, getSpecCodeErr := getSpecCodeByFlavorId(client, flavorId.(string), engineId.(string))
	if getSpecCodeErr != nil {
		return getSpecCodeErr
	}
	updateOpt.JSONBody = buildUpdateInstanceFlavorBodyParams(d, specCode)
	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating DDM instance flavor: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      ddmInstanceStatusRefreshFunc(d.Id(), client),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        100 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to running: %s", d.Id(), err)
	}
	return nil
}

func getSpecCodeByFlavorId(client *golangsdk.ServiceClient, flavorId, engineId string) (string, error) {
	// getDdmFlavors: Query the List of DDM flavors
	var (
		getDdmFlavorsHttpUrl = "v2/{project_id}/flavors"
	)

	getDdmFlavorsPath := client.Endpoint + getDdmFlavorsHttpUrl
	getDdmFlavorsPath = strings.ReplaceAll(getDdmFlavorsPath, "{project_id}", client.ProjectID)

	getDdmFlavorsQueryParams := buildGetFlavorsQueryParams(engineId, 0)
	getDdmFlavorsPath += getDdmFlavorsQueryParams
	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		getDdmFlavorsResp, err := client.Request("GET", getDdmFlavorsPath, &getInstanceOpt)
		if err != nil {
			return "", err
		}
		getDdmFlavorsRespBody, err := utils.FlattenResponse(getDdmFlavorsResp)
		if err != nil {
			return "", err
		}
		specCode, pageRes := flattenGetFlavorsResponseBody(getDdmFlavorsRespBody, flavorId)
		if specCode != "" {
			return specCode, nil
		}
		if pageRes.offset+pageRes.limit >= pageRes.x86Total && pageRes.offset+pageRes.limit >= pageRes.armTotal {
			break
		}
		getDdmFlavorsPath = updatePathOffset(getDdmFlavorsPath, pageRes.offset+pageRes.limit)
	}
	return "", fmt.Errorf("can not found flavor by flavorId: %s", flavorId)
}

func flattenGetFlavorsResponseBody(resp interface{}, flavorId string) (string, *queryRes) {
	if resp == nil {
		return "", &queryRes{}
	}
	curJson := utils.PathSearch("computeFlavorGroups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	var offset, limit, x86Total, armTotal float64
	var specCode string
	for _, v := range curArray {
		specCode = flattenFlavors(v, flavorId)
		if specCode != "" {
			return specCode, &queryRes{}
		}
		offset = utils.PathSearch("offset", v, float64(0)).(float64)
		limit = utils.PathSearch("limit", v, float64(0)).(float64)
		flavorCPUArch := utils.PathSearch("groupType", v, nil)
		if flavorCPUArch == "X86" {
			x86Total = utils.PathSearch("total", v, float64(0)).(float64)
		} else {
			armTotal = utils.PathSearch("total", v, float64(0)).(float64)
		}
	}
	return "", &queryRes{
		offset:   int(offset),
		limit:    int(limit),
		x86Total: int(x86Total),
		armTotal: int(armTotal),
	}
}

func flattenFlavors(resp interface{}, flavorId string) string {
	if resp == nil {
		return ""
	}
	curJson := utils.PathSearch("computeFlavors", resp, make([]interface{}, 0))
	for _, v := range curJson.([]interface{}) {
		id := utils.PathSearch("id", v, nil)
		if id == flavorId {
			return utils.PathSearch("code", v, nil).(string)
		}
	}
	return ""
}

func buildGetFlavorsQueryParams(engineId string, offset int) string {
	res := ""
	res = fmt.Sprintf("%s?engine_id=%v", res, engineId)
	res = fmt.Sprintf("%s&offset=%v", res, offset)
	return res
}

func updateInstanceAdminPassword(_ context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// updateInstanceAdminPassword: update DDM instance admin password
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/admin-user"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}",
		client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateInstanceAdminPasswordBodyParams(d)
	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating DDM instance admin password: %s", err)
	}
	return nil
}

func initializeParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	needRestart, err := modifyParameters(ctx, d, client, d.Get("parameters").(*schema.Set).List())
	if err != nil {
		return err
	}

	if needRestart {
		return restartDdmInstance(ctx, client, d.Id(), "soft", d.Timeout(schema.TimeoutUpdate))
	}
	return nil
}

func restartDdmInstance(ctx context.Context, client *golangsdk.ServiceClient, instanceId, restartType string, timeout time.Duration) error {
	var (
		httpUrl = "v1/{project_id}/instances/{instance_id}/action"
	)

	restartPath := client.Endpoint + httpUrl
	restartPath = strings.ReplaceAll(restartPath, "{project_id}", client.ProjectID)
	restartPath = strings.ReplaceAll(restartPath, "{instance_id}", instanceId)

	restartOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	restartOpt.JSONBody = utils.RemoveNil(buildCreateRestartBodyParams(restartType))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", restartPath, &restartOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddmInstanceStatusRefreshFunc(instanceId, client),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error restarting instance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      ddmInstanceStatusRefreshFunc(instanceId, client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to running: %s", instanceId, err)
	}
	return nil
}

func updateInstanceParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) (context.Context, error) {
	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	changes := ns.Difference(os).List()
	if len(changes) > 0 {
		needRestart, err := modifyParameters(ctx, d, client, changes)
		if err != nil {
			return ctx, nil
		}
		if needRestart {
			// Sending parametersChanged to Read to warn users the instance needs a reboot.
			ctx = context.WithValue(ctx, ctxType("parametersChanged"), "true")
		}
	}
	return ctx, nil
}

func modifyParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	parameters []interface{}) (bool, error) {
	// updateInstanceParameters: update DDM instance parameters
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/configurations"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildInstanceParametersBodyParams(parameters)

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddmInstanceStatusRefreshFunc(d.Id(), client),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return false, fmt.Errorf("error updating DDM instance parameters: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(resp.(*http.Response))
	if err != nil {
		return false, err
	}

	return utils.PathSearch("needRestart", updateRespBody, false).(bool), nil
}

func buildUpdateInstanceNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": utils.ValueIgnoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func buildUpdateInstanceSecurityGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"security_group_id": utils.ValueIgnoreEmpty(d.Get("security_group_id")),
	}
	return bodyParams
}

func buildUpdateInstanceNodeNumBodyParams(d *schema.ResourceData, nodeNumber int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"flavor_id":   utils.ValueIgnoreEmpty(d.Get("flavor_id")),
		"group_id":    utils.ValueIgnoreEmpty(d.Get("param_group_id")),
		"node_number": nodeNumber,
		"is_auto_pay": true,
	}
	return bodyParams
}

func buildUpdateInstanceFlavorBodyParams(d *schema.ResourceData, specCode string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"spec_code":   specCode,
		"group_id":    utils.ValueIgnoreEmpty(d.Get("param_group_id")),
		"is_auto_pay": true,
	}
	return bodyParams
}

func buildUpdateInstanceAdminPasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     utils.ValueIgnoreEmpty(d.Get("admin_user")),
		"password": utils.ValueIgnoreEmpty(d.Get("admin_password")),
	}
	return bodyParams
}

func buildInstanceParametersBodyParams(parameters []interface{}) map[string]interface{} {
	values := make(map[string]string)
	for _, v := range parameters {
		key := v.(map[string]interface{})["name"].(string)
		value := v.(map[string]interface{})["value"].(string)
		values[key] = value
	}
	bodyParams := map[string]interface{}{
		"values": values,
	}
	return bodyParams
}

func buildCreateRestartBodyParams(restartType string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"restart": map[string]interface{}{
			"type": restartType,
		},
	}
	return bodyParams
}

func resourceDdmInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getInstance: Query DDM instance
	var (
		getInstanceHttpUrl = "v1/{project_id}/instances/{instance_id}"
		getInstanceProduct = "ddm"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)
	getInstancePath = strings.ReplaceAll(getInstancePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DdmInstance")
	}

	getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	status := utils.PathSearch("status", getInstanceRespBody, nil)
	if status == "DELETED" {
		return diag.FromErr(mErr.ErrorOrNil())
	}

	azCodes := utils.PathSearch("available_zone", getInstanceRespBody, "")
	availabilityZones := strings.Split(azCodes.(string), ",")
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("status", status),
		d.Set("name", utils.PathSearch("name", getInstanceRespBody, nil)),
		d.Set("availability_zones", availabilityZones),
		d.Set("vpc_id", utils.PathSearch("vpc_id", getInstanceRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", getInstanceRespBody, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", getInstanceRespBody, nil)),
		d.Set("node_num", utils.PathSearch("node_count", getInstanceRespBody, nil)),
		d.Set("access_ip", utils.PathSearch("access_ip", getInstanceRespBody, nil)),
		d.Set("access_port", utils.PathSearch("access_port", getInstanceRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getInstanceRespBody, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", getInstanceRespBody, nil)),
		d.Set("nodes", flattenGetInstanceResponseBodyNodeInfoRef(getInstanceRespBody)),
		d.Set("admin_user", utils.PathSearch("admin_user_name", getInstanceRespBody, nil)),
	)
	warn := setRdsInstanceParameters(ctx, d, getInstanceClient)
	var diagnostics diag.Diagnostics
	diagnostics = append(diagnostics, diag.FromErr(mErr.ErrorOrNil())...)
	diagnostics = append(diagnostics, warn...)

	return diagnostics
}

func flattenGetInstanceResponseBodyNodeInfoRef(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("nodes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"status": utils.PathSearch("status", v, nil),
			"port":   utils.PathSearch("port", v, nil),
			"ip":     utils.PathSearch("ip", v, nil),
		})
	}
	return rst
}

func setRdsInstanceParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) diag.Diagnostics {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/configurations"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Id())

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		log.Printf("[WARN] error fetching parameters of instance (%s): %s", d.Id(), err)
		return nil
	}
	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		log.Printf("[WARN] error fetching parameters of instance (%s): %s", d.Id(), err)
		return nil
	}

	var getAccountRespBody interface{}
	err = json.Unmarshal(getRespJson, &getAccountRespBody)
	if err != nil {
		log.Printf("[WARN] error fetching parameters of instance (%s): %s", d.Id(), err)
		return nil
	}

	configs := utils.PathSearch("configuration_parameter", getAccountRespBody, make([]interface{}, 0)).([]interface{})

	var paramRestart []string
	var params []map[string]interface{}
	rawParameterList := d.Get("parameters").(*schema.Set).List()
	for _, v := range configs {
		nameRaw := utils.PathSearch("name", v, "").(string)
		valueRaw := utils.PathSearch("value", v, "").(string)
		restartRaw := utils.PathSearch("need_restart", v, "").(string)
		for _, parameter := range rawParameterList {
			name := parameter.(map[string]interface{})["name"]
			if nameRaw == name {
				p := map[string]interface{}{
					"name":  nameRaw,
					"value": valueRaw,
				}
				params = append(params, p)
				if restartRaw == "1" {
					paramRestart = append(paramRestart, nameRaw)
				}
				break
			}
		}
	}

	var diagnostics diag.Diagnostics
	if len(params) > 0 {
		if err = d.Set("parameters", params); err != nil {
			log.Printf("error saving parameters to DDM instance (%s): %s", d.Id(), err)
		}
		if len(paramRestart) > 0 && ctx.Value(ctxType("parametersChanged")) == "true" {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Parameters Changed",
				Detail:   fmt.Sprintf("parameters %s changed which needs restart.", paramRestart),
			})
		}
	}
	if len(diagnostics) > 0 {
		return diagnostics
	}
	return nil
}

func resourceDdmInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// deleteInstance: Delete DDM instance
	var (
		deleteInstanceHttpUrl = "v1/{project_id}/instances/{instance_id}"
		deleteInstanceProduct = "ddm"
	)
	deleteInstanceClient, err := cfg.NewServiceClient(deleteInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribe DDM instance: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"RUNNING", "PENDING"},
			Target:       []string{"DELETED"},
			Refresh:      ddmInstanceStatusRefreshFunc(d.Id(), deleteInstanceClient),
			Timeout:      d.Timeout(schema.TimeoutDelete),
			Delay:        50 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error deleting DDM instance (%s) error: %s", d.Id(), err)
		}
		return nil
	}

	deleteInstancePath := deleteInstanceClient.Endpoint + deleteInstanceHttpUrl
	deleteInstancePath = strings.ReplaceAll(deleteInstancePath, "{project_id}", deleteInstanceClient.ProjectID)
	deleteInstancePath = strings.ReplaceAll(deleteInstancePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	deleteInstanceQueryParams := buildDeleteInstanceQueryParams(d)
	deleteInstancePath += deleteInstanceQueryParams

	deleteInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	_, err = deleteInstanceClient.Request("DELETE", deleteInstancePath, &deleteInstanceOpt)
	if err != nil {
		return diag.Errorf("error deleting DDM instance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING", "PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      ddmInstanceStatusRefreshFunc(d.Id(), deleteInstanceClient),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to deleted: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func buildDeleteInstanceQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("delete_rds_data"); ok {
		res = fmt.Sprintf("%s&delete_rds_data=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func ddmInstanceStatusRefreshFunc(id string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getJobStatusHttpUrl = "v1/{project_id}/instances/{instance_id}"
		)

		getJobStatusPath := client.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", client.ProjectID)
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{instance_id}", fmt.Sprintf("%v", id))

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getJobStatusResp, err := client.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return getJobStatusResp, "DELETED", nil
			}
			return nil, "", err
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("status", getJobStatusRespBody, "").(string)
		if status == "CREATEFAILED" || status == "ERROR" {
			return nil, status, fmt.Errorf("the DDM instance created fail")
		}
		if status == "RUNNING" {
			return getJobStatusRespBody, status, nil
		}
		return getJobStatusRespBody, "PENDING", nil
	}
}
