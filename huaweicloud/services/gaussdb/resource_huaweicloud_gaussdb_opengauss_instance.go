package gaussdb

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ctxType string

const (
	haModeDistributed = "enterprise"
	haModeCentralized = "centralization_standard"
)

// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB POST /v3/{project_id}/instances
// @API GaussDB GET /v3/{project_id}/jobs
// @API GaussDB PUT /v3/{project_id}/instances/{instance_id}/configurations
// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/tags
// @API GaussDB GET /v3.1/{project_id}/instances/{instance_id}/configurations
// @API GaussDB PUT /v3/{project_id}/instances/{instance_id}/name
// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/password
// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/action
// @API GaussDB PUT /v3/{project_id}/instances/{instance_id}/backups/policy
// @API GaussDB PUT /v3/{project_id}/instance/{instance_id}/flavor
// @API GaussDB PUT /v3/{project_id}/configurations/{config_id}/apply
// @API GaussDB DELETE /v3/{project_id}/instances/{instance_id}/tag
// @API GaussDB DELETE /v3/{project_id}/instances/{instance_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrat
func ResourceOpenGaussInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussInstanceCreate,
		ReadContext:   resourceOpenGaussInstanceRead,
		UpdateContext: resourceOpenGaussInstanceUpdate,
		DeleteContext: resourceOpenGaussInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(150 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
				if d.HasChange("coordinator_num") {
					return d.SetNewComputed("private_ips")
				}
				return nil
			},
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ha": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: utils.SuppressCaseDiffs(),
						},
						"replication_mode": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"consistency": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: utils.SuppressCaseDiffs(),
						},
						"instance_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"sharding_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"coordinator_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"replica_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"disk_encryption_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_switch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_single_float_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
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
			"mysql_compatibility_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"advance_features": {
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
			"tags": common.TagsSchema(),
			"force_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			// Attributes
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_window": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"balance_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"error_log_switch_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceOpenGaussInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	// If force_import set, try to import it instead of creating
	if common.HasFilledOpt(d, "force_import") {
		log.Printf("[DEBUG] the Gaussdb OpenGauss instance force_import is set, try to import it instead of creating")
		instanceName := d.Get("name").(string)
		res, err := getGaussDBOpenGaussInstancesByName(client, instanceName)
		if err != nil {
			return diag.FromErr(err)
		}
		instanceId := utils.PathSearch("instances[0].id", res, nil)
		if instanceId == nil {
			return diag.Errorf("unable to retrieve instances by name: %s", instanceName)
		}
		log.Printf("[DEBUG] found existing OpenGauss instance %v with name %s", instanceId, instanceName)
		d.SetId(instanceId.(string))
		return resourceOpenGaussInstanceRead(ctx, d, meta)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBOpenGaussBodyParams(d, region))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB OpenGauss instance: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := utils.PathSearch("instance.id", createRespBody, nil)
	if instanceId == nil {
		return diag.Errorf("error creating GaussDB OpenGauss instance: id is not found in API response")
	}
	d.SetId(instanceId.(string))

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		orderId := utils.PathSearch("order_id", createRespBody, nil)
		if orderId == nil {
			return diag.Errorf("error creating GaussDB OpenGauss instance: order_id is not found in API response")
		}
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		// wait for order success
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for GaussDB OpenGauss order %s complete: %s", orderId, err)
		}
	} else {
		jobId := utils.PathSearch("job_id", createRespBody, nil)
		if jobId == nil {
			return diag.Errorf("error creating GaussDB OpenGauss instance: job_id is not found in API response")
		}
		err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 300, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"BUILD", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      instanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", d.Id(), err)
	}

	if parametersRaw := d.Get("parameters").(*schema.Set); parametersRaw.Len() > 0 {
		if err = initializeParameters(ctx, d, client, parametersRaw.List()); err != nil {
			return diag.FromErr(err)
		}
	}

	if port, ok := d.GetOk("mysql_compatibility_port"); ok && port.(string) != "0" {
		if err = openMysqlCompatibilityPort(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("advance_features"); ok {
		if err = updateAdvanceFeatures(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		err = addInstanceTags(d, client, d.Get("tags").(map[string]interface{}))
		if err != nil {
			return diag.Errorf("error setting tags for GaussDB OpenGauss instance %s: %s", d.Id(), err)
		}
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(360 * time.Second) // lintignore:R018

	return resourceOpenGaussInstanceRead(ctx, d, meta)
}

func getGaussDBOpenGaussInstancesById(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	queryParam := fmt.Sprintf("?id=%s", id)
	return getGaussDBOpenGaussInstances(client, queryParam)
}

func getGaussDBOpenGaussInstancesByName(client *golangsdk.ServiceClient, name string) (interface{}, error) {
	queryParam := fmt.Sprintf("?name=%s", name)
	return getGaussDBOpenGaussInstances(client, queryParam)
}

func getGaussDBOpenGaussInstances(client *golangsdk.ServiceClient, queryParam string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += queryParam

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	return getRespBody, nil
}

func buildCreateGaussDBOpenGaussBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"datastore":             buildCreateGaussDBOpenGaussDatastoreBodyParams(d),
		"ha":                    buildCreateGaussDBOpenGaussHaBodyParams(d),
		"configuration_id":      utils.ValueIgnoreEmpty(d.Get("configuration_id")),
		"port":                  utils.ValueIgnoreEmpty(d.Get("port")),
		"password":              d.Get("password"),
		"backup_strategy":       buildCreateGaussDBOpenGaussBackupStrategyBodyParams(d),
		"enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"disk_encryption_id":    utils.ValueIgnoreEmpty(d.Get("disk_encryption_id")),
		"flavor_ref":            d.Get("flavor"),
		"volume":                buildCreateGaussDBOpenGaussVolumeBodyParams(d),
		"region":                region,
		"availability_zone":     utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		"vpc_id":                d.Get("vpc_id"),
		"subnet_id":             d.Get("subnet_id"),
		"security_group_id":     utils.ValueIgnoreEmpty(d.Get("security_group_id")),
		"charge_info":           buildCreateGaussDBOpenGaussChargeInfoBodyParams(d),
		"time_zone":             utils.ValueIgnoreEmpty(d.Get("time_zone")),
		"sharding_num":          utils.ValueIgnoreEmpty(d.Get("sharding_num")),
		"coordinator_num":       utils.ValueIgnoreEmpty(d.Get("coordinator_num")),
		"replica_num":           utils.ValueIgnoreEmpty(d.Get("replica_num")),
	}
	if v := d.Get("enable_force_switch").(bool); v {
		bodyParams["enable_force_switch"] = v
	}
	if v := d.Get("enable_single_float_ip").(bool); v {
		bodyParams["enable_single_float_ip"] = v
	}
	return bodyParams
}

func buildCreateGaussDBOpenGaussDatastoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("datastore").([]interface{})
	if len(rawParams) > 0 {
		rawParam := rawParams[0].(map[string]interface{})
		return map[string]interface{}{
			"type":    rawParam["engine"],
			"version": rawParam["version"],
		}
	}
	return map[string]interface{}{
		"type": "GaussDB(for openGauss)",
	}
}

func buildCreateGaussDBOpenGaussHaBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("ha").([]interface{})
	rawParam := rawParams[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"mode":             rawParam["mode"],
		"replication_mode": rawParam["replication_mode"],
		"consistency":      rawParam["consistency"],
		"instance_mode":    utils.ValueIgnoreEmpty(rawParam["instance_mode"]),
	}
	return bodyParams
}

func buildCreateGaussDBOpenGaussBackupStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("backup_strategy").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	rawParam := rawParams[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"start_time": rawParam["start_time"],
		"keep_days":  utils.ValueIgnoreEmpty(rawParam["keep_days"]),
	}
	return bodyParams
}

func buildCreateGaussDBOpenGaussVolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("volume").([]interface{})
	mode := d.Get("ha").([]interface{})[0].(map[string]interface{})["mode"].(string)
	dnNum := 1
	if mode == haModeDistributed {
		dnNum = d.Get("sharding_num").(int)
	}
	if mode == haModeCentralized {
		dnNum = d.Get("replica_num").(int) + 1
	}
	raw := rawParams[0].(map[string]interface{})
	dnSize := raw["size"].(int)
	bodyParams := map[string]interface{}{
		"type": raw["type"],
		"size": dnSize * dnNum,
	}
	return bodyParams
}

func buildCreateGaussDBOpenGaussChargeInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("charging_mode"); ok && v == "prePaid" {
		bodyParams := map[string]interface{}{
			"charge_mode":   utils.ValueIgnoreEmpty(d.Get("charging_mode")),
			"period_type":   utils.ValueIgnoreEmpty(d.Get("period_unit")),
			"period_num":    utils.ValueIgnoreEmpty(d.Get("period")),
			"is_auto_renew": d.Get("auto_renew"),
			"is_auto_pay":   "true",
		}
		return bodyParams
	}
	return nil
}

func buildGaussDBOpenGaussParameters(params []interface{}) map[string]interface{} {
	values := make(map[string]interface{})
	for _, v := range params {
		key := v.(map[string]interface{})["name"].(string)
		value := v.(map[string]interface{})["value"]
		values[key] = value
	}
	bodyParams := map[string]interface{}{
		"values": values,
	}
	return bodyParams
}

func initializeParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, parametersRaw []interface{}) error {
	updateOpts := buildGaussDBOpenGaussParameters(parametersRaw)
	restartRequired, err := modifyParameters(ctx, client, d, schema.TimeoutCreate, &updateOpts)
	if err != nil {
		return err
	}

	if restartRequired {
		return restartGaussDBOpenGaussInstance(ctx, client, d)
	}
	return nil
}

func openMysqlCompatibilityPort(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/mysql-compatibility"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateMysqlCompatibilityPortBodyParams(d)
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error opening GaussDB OpenGauss instance (%s) MySQL compatibility port: %s", d.Id(), err)
	}

	updateRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error opening GaussDB OpenGauss instance MySQL compatibility port: job_id is not " +
			"found in API response")
	}
	return checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 10, d.Timeout(timeout))
}

func restartGaussDBOpenGaussInstance(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/restart"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error restarting GaussDB OpenGauss instance (%s): %s", d.Id(), err)
	}

	updateRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error restarting GaussDB OpenGauss instance name: job_id is not found in API response")
	}

	return checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 30, d.Timeout(schema.TimeoutUpdate))
}

func resourceOpenGaussInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	res, err := getGaussDBOpenGaussInstancesById(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	instance := utils.PathSearch("instances[0]", res, nil)
	if instance == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	var dnNum = 1
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("status", utils.PathSearch("status", instance, nil)),
		d.Set("type", utils.PathSearch("type", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", instance, nil)),
		d.Set("db_user_name", utils.PathSearch("db_user_name", instance, nil)),
		d.Set("time_zone", utils.PathSearch("time_zone", instance, nil)),
		d.Set("flavor", utils.PathSearch("flavor_ref", instance, nil)),
		d.Set("port", strconv.Itoa(int(utils.PathSearch("port", instance, float64(0)).(float64)))),
		d.Set("switch_strategy", utils.PathSearch("switch_strategy", instance, nil)),
		d.Set("maintenance_window", utils.PathSearch("maintenance_window", instance, nil)),
		d.Set("public_ips", utils.PathSearch("public_ips", instance, nil)),
		d.Set("charging_mode", utils.PathSearch("charge_info.charge_mode", instance, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", instance, nil)),
		d.Set("ha", flattenGaussDBOpenGaussResponseBodyHa(instance)),
		d.Set("datastore", flattenGaussDBOpenGaussResponseBodyDatastore(instance)),
		d.Set("backup_strategy", flattenGaussDBOpenGaussResponseBodyBackupStrategy(instance)),
		setOpenGaussNodesAndRelatedNumbers(d, instance, &dnNum),
		d.Set("volume", flattenGaussDBOpenGaussResponseBodyVolume(instance, dnNum)),
		d.Set("mysql_compatibility_port", utils.PathSearch("mysql_compatibility.port", instance, nil)),
		setOpenGaussPrivateIpsAndEndpoints(d, instance),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", instance, make([]interface{}, 0)))),
	)

	mErr = multierror.Append(mErr, setBalanceStatus(d, client))
	mErr = multierror.Append(mErr, setErrorLogSwitchStatus(d, client))
	mErr = multierror.Append(mErr, setAdvanceFeatures(d, client))

	diagErr := setGaussDBMySQLParameters(ctx, d, client)
	resErr := append(diag.FromErr(mErr.ErrorOrNil()), diagErr...)

	return resErr
}

func flattenGaussDBOpenGaussResponseBodyHa(instance interface{}) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"mode":             strings.ToLower(utils.PathSearch("type", instance, "").(string)),
			"replication_mode": utils.PathSearch("ha.replication_mode", instance, nil),
			"consistency":      utils.PathSearch("ha.consistency", instance, nil),
			"instance_mode":    utils.PathSearch("instance_mode", instance, nil),
		},
	}
	return rst
}

func flattenGaussDBOpenGaussResponseBodyDatastore(instance interface{}) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"engine":  utils.PathSearch("datastore.type", instance, nil),
			"version": utils.PathSearch("datastore.version", instance, nil),
		},
	}
	return rst
}

func flattenGaussDBOpenGaussResponseBodyBackupStrategy(instance interface{}) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"start_time": utils.PathSearch("backup_strategy.start_time", instance, nil),
			"keep_days":  utils.PathSearch("backup_strategy.keep_days", instance, nil),
		},
	}
	return rst
}

func flattenGaussDBOpenGaussResponseBodyVolume(instance interface{}, dnNum int) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"type": utils.PathSearch("volume.type", instance, nil),
			"size": int(utils.PathSearch("volume.size", instance, float64(0)).(float64)) / dnNum,
		},
	}
	return rst
}

func setOpenGaussNodesAndRelatedNumbers(d *schema.ResourceData, instance interface{}, dnNum *int) *multierror.Error {
	shardingNum := 0
	coordinatorNum := 0

	nodesArray := utils.PathSearch("nodes", instance, make([]interface{}, 0)).([]interface{})
	nodesList := make([]map[string]interface{}, 0, len(nodesArray))
	for _, v := range nodesArray {
		name := utils.PathSearch("name", v, "").(string)
		node := map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              name,
			"status":            utils.PathSearch("status", v, nil),
			"role":              utils.PathSearch("role", v, nil),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
			"private_ip":        utils.PathSearch("private_ip", v, nil),
			"public_ip":         utils.PathSearch("public_ip", v, nil),
		}
		nodesList = append(nodesList, node)

		if strings.Contains(name, "_gaussdbv5cn") {
			coordinatorNum++
		} else if strings.Contains(name, "_gaussdbv5dn") {
			shardingNum++
		}
	}

	if shardingNum > 0 && coordinatorNum > 0 {
		*dnNum = shardingNum / d.Get("replica_num").(int)
		mErr := multierror.Append(
			d.Set("nodes", nodesList),
			d.Set("sharding_num", dnNum),
			d.Set("coordinator_num", coordinatorNum),
		)
		return mErr
	}

	// If the HA mode is centralized, the HA structure of API response is nil.
	replicaNum := utils.PathSearch("replica_num", instance, float64(0)).(float64)
	*dnNum = int(replicaNum) + 1
	mErr := multierror.Append(
		d.Set("nodes", nodesList),
		d.Set("replica_num", replicaNum),
	)
	return mErr
}

func setOpenGaussPrivateIpsAndEndpoints(d *schema.ResourceData, instance interface{}) *multierror.Error {
	privateIps := utils.PathSearch("private_ips", instance, make([]interface{}, 0)).([]interface{})
	if len(privateIps) == 0 {
		return nil
	}

	port := utils.PathSearch("port", instance, float64(0)).(float64)
	ipList := strings.Split(privateIps[0].(string), "/")
	endpoints := make([]string, 0, len(ipList))
	for i := 0; i < len(ipList); i++ {
		ipList[i] = strings.Trim(ipList[i], " ")
		endpoint := fmt.Sprintf("%s:%v", ipList[i], port)
		endpoints = append(endpoints, endpoint)
	}

	mErr := multierror.Append(
		d.Set("private_ips", ipList),
		d.Set("endpoints", endpoints),
	)

	return mErr
}

func setGaussDBMySQLParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) diag.Diagnostics {
	var (
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/configurations"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		log.Printf("[WARN] error retrieving GaussDB OpenGauss(%s) parameters: %s", d.Id(), err)
		return nil
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	rawParameters := d.Get("parameters").(*schema.Set)
	rawParameterMap := make(map[string]bool)
	for _, rawParameter := range rawParameters.List() {
		rawParameterMap[rawParameter.(map[string]interface{})["name"].(string)] = true
	}

	var configurationRestart bool
	var params []map[string]interface{}
	parametersList := utils.PathSearch("configuration_parameters", getRespBody, make([]interface{}, 0)).([]interface{})
	for _, v := range parametersList {
		name := utils.PathSearch("name", v, "").(string)
		value := utils.PathSearch("value", v, nil)
		if utils.PathSearch("restart_required", v, false).(bool) {
			configurationRestart = true
		}
		if rawParameterMap[name] {
			p := map[string]interface{}{
				"name":  name,
				"value": value,
			}
			params = append(params, p)
		}
	}

	var diagnostics diag.Diagnostics
	if len(params) > 0 {
		if err = d.Set("parameters", params); err != nil {
			log.Printf("error saving parameters to GaussDB OpenGauss instance (%s): %s", d.Id(), err)
		}
		if ctx.Value(ctxType("parametersChanged")) == "true" {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Parameters Changed",
				Detail:   "Parameters changed which needs reboot.",
			})
		}
	}
	if configurationRestart && ctx.Value(ctxType("configurationChanged")) == "true" {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Configuration Changed",
			Detail:   "Configuration changed, the instance may needs reboot.",
		})
	}
	if len(diagnostics) > 0 {
		return diagnostics
	}
	return nil
}

func setBalanceStatus(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/balance"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		log.Printf("[WARN] error retrieving GaussDB OpenGauss(%s) balance status: %s", d.Id(), err)
		return nil
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return err
	}
	return d.Set("balance_status", utils.PathSearch("balanced", getRespBody, nil))
}

func setErrorLogSwitchStatus(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/error-log/switch/status"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		log.Printf("[WARN] error retrieving GaussDB OpenGauss(%s) error log switch status: %s", d.Id(), err)
		return nil
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return err
	}
	return d.Set("error_log_switch_status", utils.PathSearch("status", getRespBody, nil))
}

func setAdvanceFeatures(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	features, err := getAdvanceFeatures(client, d.Id())
	if err != nil {
		log.Printf("[WARN] error retrieving GaussDB OpenGauss(%s) error log switch status: %s", d.Id(), err)
		return nil
	}

	rawAdvanceFeatures := d.Get("advance_features").(*schema.Set)
	rawAdvanceFeatureMap := make(map[string]bool)
	for _, rawAdvanceFeature := range rawAdvanceFeatures.List() {
		rawAdvanceFeatureMap[rawAdvanceFeature.(map[string]interface{})["name"].(string)] = true
	}
	rst := make([]interface{}, 0, len(features))
	for _, v := range features {
		name := utils.PathSearch("name", v, "").(string)
		if rawAdvanceFeatureMap[name] {
			rst = append(rst, map[string]string{
				"name":  utils.PathSearch("name", v, "").(string),
				"value": utils.PathSearch("value", v, "").(string),
			})
		}
	}
	return d.Set("advance_features", rst)
}

func getAdvanceFeatures(client *golangsdk.ServiceClient, instanceId string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/advance-features"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	features := utils.PathSearch("features", getRespBody, make([]interface{}, 0)).([]interface{})

	return features, nil
}

func resourceOpenGaussInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}
	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if d.HasChange("name") {
		if err = updateInstanceName(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password") {
		if err = updateInstancePassword(d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("sharding_num") {
		if err = expandInstanceShardingNumber(ctx, d, client, bssClient); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("coordinator_num") {
		if err = expandInstanceCoordinatorNumber(ctx, d, client, bssClient); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("volume") {
		if err = updateInstanceVolumeSize(ctx, d, client, bssClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("backup_strategy") {
		if err = updateInstanceBackupStrategy(d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("flavor") {
		if err = updateInstanceFlavor(ctx, d, client, bssClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("configuration_id") {
		ctx, err = updateConfiguration(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}

		// if parameters is set, it should be modified
		if params, ok := d.GetOk("parameters"); ok {
			updateOpts := buildGaussDBOpenGaussParameters(params.(*schema.Set).List())
			_, err = modifyParameters(ctx, client, d, schema.TimeoutUpdate, &updateOpts)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("parameters") && !d.HasChanges("configuration_id") {
		ctx, err = updateParameters(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("mysql_compatibility_port") {
		res, err := getGaussDBOpenGaussInstancesById(client, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		port := utils.PathSearch("instances[0].mysql_compatibility.port", res, nil)
		if port == nil || port.(string) == "0" {
			err = openMysqlCompatibilityPort(ctx, client, d, schema.TimeoutUpdate)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err = updateMysqlCompatibilityPort(ctx, client, d, schema.TimeoutUpdate)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("advance_features") {
		if err = updateAdvanceFeatures(ctx, client, d, schema.TimeoutUpdate); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "gaussdb",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err = cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		err = updateInstanceTags(d, client)
		if err != nil {
			return diag.Errorf("error updating tags of GaussDB OpenGauss instance %q: %s", d.Id(), err)
		}
	}

	return resourceOpenGaussInstanceRead(ctx, d, meta)
}

func updateInstanceName(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/name"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceNameBodyParams(d))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB OpenGauss instance (%s) name: %s", d.Id(), err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error updating GaussDB OpenGauss instance name: job_id is not found in API response")
	}

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 2, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateInstanceNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name"),
	}
	return bodyParams
}

func updateInstancePassword(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/password"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateInstancePasswordBodyParams(d))

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB OpenGauss instance (%s) password: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateInstancePasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"password": d.Get("password"),
	}
	return bodyParams
}

func expandInstanceShardingNumber(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	oRaw, nRaw := d.GetChange("sharding_num")
	if nRaw.(int) < oRaw.(int) {
		return fmt.Errorf("error expanding shard for instance: new num must be larger than the old one")
	}
	expandSize := nRaw.(int) - oRaw.(int)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildExpandInstanceShardingNumberBodyParams(expandSize))
	return updateInstanceVolumeAndRelatedHaNumbers(ctx, client, bssClient, d, updateOpt)
}

func buildExpandInstanceShardingNumberBodyParams(expandSize int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"expand_cluster": map[string]interface{}{
			"shard": map[string]interface{}{
				"count": expandSize,
			},
		},
		"is_auto_pay": "true",
	}
	return bodyParams
}

func expandInstanceCoordinatorNumber(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	oRaw, nRaw := d.GetChange("coordinator_num")
	if nRaw.(int) < oRaw.(int) {
		return fmt.Errorf("error expanding coordinator for instance: new number must be larger than the old one")
	}
	expandSize := nRaw.(int) - oRaw.(int)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildExpandInstanceCoordinatorNumberBodyParams(d, expandSize))
	return updateInstanceVolumeAndRelatedHaNumbers(ctx, client, bssClient, d, updateOpt)
}

func buildExpandInstanceCoordinatorNumberBodyParams(d *schema.ResourceData, expandSize int) map[string]interface{} {
	coordinators := make([]interface{}, 0)
	azList := strings.Split(d.Get("availability_zone").(string), ",")
	for i := 0; i < expandSize; i++ {
		coordinators = append(coordinators, map[string]interface{}{
			"az_code": azList[0],
		})
	}
	bodyParams := map[string]interface{}{
		"expand_cluster": map[string]interface{}{
			"coordinators": coordinators,
		},
		"is_auto_pay": "true",
	}
	return bodyParams
}

func updateInstanceVolumeSize(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	volumeRaw := d.Get("volume").([]interface{})
	dnSize := volumeRaw[0].(map[string]interface{})["size"].(int)
	dnNum := 1
	if d.Get("ha.0.mode").(string) == haModeDistributed {
		dnNum = d.Get("sharding_num").(int)
	}
	if d.Get("ha.0.mode").(string) == haModeCentralized {
		dnNum = d.Get("replica_num").(int) + 1
	}

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceVolumeSizeBodyParams(dnSize, dnNum))
	return updateInstanceVolumeAndRelatedHaNumbers(ctx, client, bssClient, d, updateOpt)
}

func buildUpdateInstanceVolumeSizeBodyParams(dnSize, dnNum int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enlarge_volume": map[string]interface{}{
			"size": dnSize * dnNum,
		},
		"is_auto_pay": "true",
	}
	return bodyParams
}

func updateInstanceVolumeAndRelatedHaNumbers(ctx context.Context, client, bssClient *golangsdk.ServiceClient,
	d *schema.ResourceData, opts golangsdk.RequestOpts) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/action"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", updatePath, &opts)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating GaussDB OpenGauss instance (%s): %s", d.Id(), err)
	}

	updateRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		orderId := utils.PathSearch("orderId", updateRespBody, nil)
		if orderId == nil {
			return fmt.Errorf("error updating GaussDB OpenGauss instance: order_id is not found in API response")
		}
		// wait for order success
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	} else {
		jobId := utils.PathSearch("job_id", updateRespBody, nil)
		if jobId == nil {
			return fmt.Errorf("error updating GaussDB OpenGauss instance: job_id is not found in API response")
		}
		err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 180, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"MODIFYING", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      instanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to become ready: %s", d.Id(), err)
	}

	return nil
}

func updateInstanceBackupStrategy(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/backups/policy"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceBackupStrategyBodyParams(d))

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB OpenGauss instance (%s) backup strategy: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateInstanceBackupStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("backup_strategy").([]interface{})
	rawParam := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"start_time":          rawParam["start_time"],
		"keep_days":           rawParam["keep_days"],
		"period":              "1,2,3,4,5,6,7",
		"differential_period": "30",
	}
	bodyParams := map[string]interface{}{
		"backup_policy": params,
	}
	return bodyParams
}

func updateInstanceFlavor(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instance/{instance_id}/flavor"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceFlavorBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating GaussDB OpenGauss instance (%s) flavor: %s", d.Id(), err)
	}

	updateRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		orderId := utils.PathSearch("order_id", updateRespBody, nil)
		if orderId == nil {
			return fmt.Errorf("error updating GaussDB OpenGauss instance flavor: order_id is not found in API response")
		}
		// wait for order success
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	} else {
		jobId := utils.PathSearch("job_id", updateRespBody, nil)
		if jobId == nil {
			return fmt.Errorf("error updating GaussDB OpenGauss instance flavor: job_id is not found in API response")
		}
		err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 180, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"MODIFYING", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      instanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to become ready: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateInstanceFlavorBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"flavor_ref": d.Get("flavor"),
	}
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		bodyParams["is_auto_pay"] = true
	}
	return bodyParams
}

func updateConfiguration(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) (context.Context, error) {
	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}/apply"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{config_id}", d.Get("configuration_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceConfigurationBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return ctx, fmt.Errorf("error updating GaussDB OpenGauss instance (%s) configuration: %s", d.Id(), err)
	}

	updateRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return ctx, err
	}
	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return ctx, fmt.Errorf("error updating GaussDB OpenGauss instance configuration: job_id is not found in API response")
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 2, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return ctx, err
	}

	// Sending configurationChanged to Read to warn users the instance needs a reboot.
	ctx = context.WithValue(ctx, ctxType("configurationChanged"), "true")

	return ctx, nil
}

func buildUpdateInstanceConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_ids": []string{d.Id()},
	}
	return bodyParams
}

func updateParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) (context.Context, error) {
	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	change := ns.Difference(os)
	if change.Len() > 0 {
		updateOpts := buildGaussDBOpenGaussParameters(change.List())
		restartRequired, err := modifyParameters(ctx, client, d, schema.TimeoutUpdate, &updateOpts)
		if err != nil {
			return ctx, err
		}
		if restartRequired {
			// Sending parametersChanged to Read to warn users the instance needs a reboot.
			ctx = context.WithValue(ctx, ctxType("parametersChanged"), "true")
		}
	}

	return ctx, nil
}

func modifyParameters(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string,
	parameterOpts *map[string]interface{}) (bool, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/configurations"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         parameterOpts,
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return false, fmt.Errorf("error updating GaussDB OpenGauss instance (%s) parameters: %s", d.Id(), err)
	}

	updateRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return false, err
	}
	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return false, fmt.Errorf("error updating GaussDB OpenGauss instance parameters: job_id is not found in API response")
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 2, d.Timeout(timeout))
	if err != nil {
		return false, err
	}
	restartRequired := utils.PathSearch("restart_required", updateRespBody, false).(bool)
	return restartRequired, nil
}

func updateMysqlCompatibilityPort(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/mysql-compatibility"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateMysqlCompatibilityPortBodyParams(d)
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating GaussDB OpenGauss instance (%s) MySQL compatibility port: %s", d.Id(), err)
	}

	updateRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error updating GaussDB OpenGauss instance MySQL compatibility port: job_id is not " +
			"found in API response")
	}
	return checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 10, d.Timeout(timeout))
}

func buildUpdateMysqlCompatibilityPortBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"port": d.Get("mysql_compatibility_port"),
	}
	return bodyParams
}

func updateAdvanceFeatures(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/advance-features"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateAdvanceFeaturesBodyParams(d)
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating GaussDB OpenGauss instance (%s) advance features: %s", d.Id(), err)
	}

	return checkAdvanceFeaturesJobFinish(ctx, d, client, timeout)
}

func buildUpdateAdvanceFeaturesBodyParams(d *schema.ResourceData) map[string]interface{} {
	advanceFeatures := d.Get("advance_features").(*schema.Set).List()
	params := make(map[string]string)
	for _, v := range advanceFeatures {
		key := v.(map[string]interface{})["name"].(string)
		value := v.(map[string]interface{})["value"].(string)
		params[key] = value
	}
	bodyParams := map[string]interface{}{
		"params": params,
	}
	return bodyParams
}

func checkAdvanceFeaturesJobFinish(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	timeout string) error {
	rawAdvanceFeatures := d.Get("advance_features").(*schema.Set).List()
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Completed"},
		Refresh:      gaussDBOpenGaussAdvanceFeaturesRefreshFunc(client, d.Id(), rawAdvanceFeatures),
		Timeout:      d.Timeout(timeout),
		Delay:        2,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for GaussDB Opengauss advance features job (%s) to be completed", err)
	}
	return nil
}

func gaussDBOpenGaussAdvanceFeaturesRefreshFunc(client *golangsdk.ServiceClient, instanceId string,
	rawAdvanceFeatures []interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		features, err := getAdvanceFeatures(client, instanceId)
		if err != nil {
			return nil, "Error", err
		}

		featuresMap := make(map[string]string)
		for _, v := range features {
			name := utils.PathSearch("name", v, "").(string)
			value := utils.PathSearch("value", v, "").(string)
			featuresMap[name] = value
		}
		for _, rawAdvanceFeature := range rawAdvanceFeatures {
			name := rawAdvanceFeature.(map[string]interface{})["name"].(string)
			value := rawAdvanceFeature.(map[string]interface{})["value"].(string)
			// if raw advance feature name is not exist, it indicates the name is error, then it is no need to check the value
			if v, ok := featuresMap[name]; ok && v != value {
				return features, "Pending", nil
			}
		}

		return features, "Completed", nil
	}
}

func addInstanceTags(d *schema.ResourceData, client *golangsdk.ServiceClient, addTags map[string]interface{}) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/tags"
	)

	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{project_id}", client.ProjectID)
	addPath = strings.ReplaceAll(addPath, "{instance_id}", d.Id())

	addOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	addOpt.JSONBody = utils.RemoveNil(buildAddInstanceTagsBodyParams(addTags))

	_, err := client.Request("POST", addPath, &addOpt)
	if err != nil {
		return fmt.Errorf("error adding tags to GaussDB OpenGauss instance (%s): %s", d.Id(), err)
	}

	return nil
}

func buildAddInstanceTagsBodyParams(addTags map[string]interface{}) map[string]interface{} {
	tags := make([]interface{}, 0, len(addTags))
	for key, value := range addTags {
		tags = append(tags, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}
	bodyParams := map[string]interface{}{
		"tags": tags,
	}
	return bodyParams
}

func deleteInstanceTags(d *schema.ResourceData, client *golangsdk.ServiceClient, deleteTagKeys []string) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/tag"
	)

	deleteBasePath := client.Endpoint + httpUrl
	deleteBasePath = strings.ReplaceAll(deleteBasePath, "{project_id}", client.ProjectID)
	deleteBasePath = strings.ReplaceAll(deleteBasePath, "{instance_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for _, deleteTagKey := range deleteTagKeys {
		deletePath := deleteBasePath + buildDeleteInstanceTagParamBodyParams(deleteTagKey)
		_, err := client.Request("DELETE", deletePath, &deleteOpt)
		if err != nil {
			return fmt.Errorf("error deleting tag(%s) from GaussDB OpenGauss instance (%s): %s", d.Id(), deleteTagKey, err)
		}
	}

	return nil
}

func buildDeleteInstanceTagParamBodyParams(key string) string {
	return fmt.Sprintf("?key=%s", key)
}

func updateInstanceTags(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	if len(oMap) > 0 {
		keys := make([]string, 0, len(oMap))
		for key := range oMap {
			keys = append(keys, key)
		}
		err := deleteInstanceTags(d, client, keys)
		if err != nil {
			return err
		}
	}

	if len(nMap) > 0 {
		err := addInstanceTags(d, client, nMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceOpenGaussInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribe OpenGauss instance: %s", err)
		}
	} else {
		deletePath := client.Endpoint + httpUrl
		deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
		deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())

		deleteOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error deleting GaussDB OpenGauss instance")
		}

		deleteRespBody, err := utils.FlattenResponse(deleteResp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobId := utils.PathSearch("job_id", deleteRespBody, nil)
		if jobId == nil {
			return diag.Errorf("error deleting GaussDB OpenGauss instance: job_id is not found in API response")
		}

		err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 5, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func checkGaussDBOpenGaussJobFinish(ctx context.Context, client *golangsdk.ServiceClient, jobID string, delay int,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      gaussDBOpenGaussStatusRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        time.Duration(delay) * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for GaussDB Opengauss job (%s) to be completed: %s ", jobID, err)
	}
	return nil
}

func instanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := getGaussDBOpenGaussInstancesById(client, instanceID)
		if err != nil {
			return nil, "", err
		}
		instance := utils.PathSearch("instances[0]", res, nil)
		if instance == nil {
			return struct{}{}, "DELETED", nil
		}
		status := utils.PathSearch("status", instance, "").(string)
		return instance, status, nil
	}
}

func gaussDBOpenGaussStatusRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getJobStatusHttpUrl = "v3/{project_id}/jobs?id={job_id}"
		)

		getJobStatusPath := client.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", client.ProjectID)
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{job_id}", jobId)

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}
		getJobStatusResp, err := client.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			return nil, "Failed", err
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("job.status", getJobStatusRespBody, "")
		return getJobStatusRespBody, status.(string), nil
	}
}
