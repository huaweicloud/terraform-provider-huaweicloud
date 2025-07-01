package rds

import (
	"context"
	"fmt"
	"reflect"
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

var pgPluginParameterNonUpdatableParams = []string{"instance_id", "name"}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/parameter/{name}
// @API RDS GET /v3/{project_id}/instances/{instance_id}/parameter/{name}
// @API RDS GET /v3/{project_id}/instances
func ResourcePgPluginParameter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgPluginParameterCreateOrUpdate,
		ReadContext:   resourcePgPluginParameterRead,
		UpdateContext: resourcePgPluginParameterCreateOrUpdate,
		DeleteContext: resourcePgPluginParameterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRdsPgPluginParameterImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(pgPluginParameterNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of RDS instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the plugin parameter.`,
			},
			"values": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of plugin parameter values.`,
			},
			"restart_required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether a reboot is required.`,
			},
			"default_values": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the default values of the plugin parameter.`,
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

func resourcePgPluginParameterCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	valuesRaw := d.Get("values").(*schema.Set).List()
	values := make([]string, 0)
	values = append(values, defaultValues[d.Get("name").(string)]...)
	values = append(values, utils.ExpandToStringList(valuesRaw)...)
	err = updatePgPluginParameter(ctx, d, client, values, schema.TimeoutCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	d.SetId(fmt.Sprintf("%s/%s", instanceId, name))

	return resourcePgPluginParameterRead(ctx, d, meta)
}

func resourcePgPluginParameterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	respBody, err := getPgPluginParameter(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	respValueString := utils.PathSearch("value", respBody, nil)
	if respValueString == nil {
		return diag.Errorf("error found plugin parameter values")
	}
	respValues := strings.Split(respValueString.(string), ",")

	defaults := defaultValues[d.Get("name").(string)]
	defaultsMap := make(map[string]bool)
	for _, value := range defaults {
		defaultsMap[value] = true
	}
	values := make([]string, 0)
	for _, value := range respValues {
		if !defaultsMap[value] {
			values = append(values, value)
		}
	}
	if len(values) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("values", values),
		d.Set("restart_required", utils.PathSearch("restart_required", respBody, nil)),
		d.Set("default_values", defaults),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePgPluginParameterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	err = updatePgPluginParameter(ctx, d, client, defaultValues[d.Get("name").(string)], schema.TimeoutDelete)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updatePgPluginParameter(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, values []string,
	timeout string) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/parameter/{name}"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{name}", d.Get("name").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdatePgPluginParameterBodyParams(values))
	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error creating RDS PostgreSQL plugin parameter values: %s", err)
	}
	return waitForUpdateParameterValueCompleted(ctx, d, client, values, timeout)
}

func buildUpdatePgPluginParameterBodyParams(values []string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"value": strings.Join(values, ","),
	}
	return bodyParams
}

func waitForUpdateParameterValueCompleted(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	values []string, timeout string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      pluginParameterValueRefreshFunc(d, client, values),
		Timeout:      d.Timeout(timeout),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for RDS PostgreSQL plugin parameter values update completed: %s ", err)
	}
	return nil
}

func pluginParameterValueRefreshFunc(d *schema.ResourceData, client *golangsdk.ServiceClient, values []string) resource.StateRefreshFunc {
	targetValues := make([]string, 0)
	targetValues = append(targetValues, defaultValues[d.Get("name").(string)]...)
	targetValues = append(targetValues, values...)
	targetValuesMap := make(map[string]bool)
	for _, value := range targetValues {
		targetValuesMap[value] = true
	}
	return func() (interface{}, string, error) {
		respBody, err := getPgPluginParameter(d, client)
		if err != nil {
			return nil, "ERROR", err
		}

		respValues := utils.PathSearch("value", respBody, nil)
		if respValues == nil {
			return nil, "ERROR", fmt.Errorf("error found plugin parameter values")
		}
		valuesMap := make(map[string]bool)
		for _, value := range values {
			valuesMap[value] = true
		}
		if !reflect.DeepEqual(targetValuesMap, valuesMap) {
			return respBody, "PENDING", nil
		}
		return respBody, "COMPLETED", nil
	}
}

func getPgPluginParameter(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/parameter/{name}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{name}", d.Get("name").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL plugin parameter values: %s", err)
	}
	return utils.FlattenResponse(getResp)
}

func resourceRdsPgPluginParameterImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}

	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("name", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
