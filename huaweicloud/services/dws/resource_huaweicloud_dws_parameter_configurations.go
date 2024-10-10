package dws

import (
	"context"
	"fmt"
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

// @API DWS PUT /v2/{project_id}/clusters/{cluster_id}/configurations/{configuration_id}
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}/configurations
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}/configurations/{configuration_id}
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}
func ResourceParameterConfigurations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParameterConfigurationsCreate,
		ReadContext:   resourceParameterConfigurationsRead,
		UpdateContext: resourceParameterConfigurationsUpdate,
		DeleteContext: resourceParameterConfigurationsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The DWS cluster ID.`,
			},
			"configurations": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The list of the DWS cluster parameter configurations.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the parameter.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the parameter.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The value of the parameter.`,
						},
					},
				},
			},
		},
	}
}

func resourceParameterConfigurationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		clusterId = d.Get("cluster_id").(string)
	)
	// For the same DWS cluster, it is not supported to run multiple tasks at the same time.
	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = modifyConfigurations(ctx, client, clusterId, d.Get("configurations").([]interface{}), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterId)
	return resourceParameterConfigurationsRead(ctx, d, meta)
}

func modifyConfigurations(ctx context.Context, client *golangsdk.ServiceClient, clusterId string, configurations []interface{},
	timeout time.Duration) error {
	if err := waitClusterTaskStateCompleted(ctx, client, timeout, clusterId); err != nil {
		return err
	}

	group, err := getParameterGroup(client, clusterId)
	if err != nil {
		return err
	}

	httpUrl := "v2/{project_id}/clusters/{cluster_id}/configurations/{configuration_id}"
	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)
	modifyPath = strings.ReplaceAll(modifyPath, "{cluster_id}", clusterId)
	modifyPath = strings.ReplaceAll(modifyPath, "{configuration_id}", utils.PathSearch("id", group, "").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"configurations": configurations,
		},
	}

	_, err = client.Request("PUT", modifyPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error setting the DWS cluster (%s) paramsters: %s", clusterId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshParamesterStateFun(client, clusterId),
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the parameters modification to completed: %s", err)
	}
	return nil
}

func refreshParamesterStateFun(client *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getParameterGroup(client, clusterId)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		if utils.StrSliceContains([]string{"In-Sync", "Pending-Reboot"}, status) {
			return respBody, "COMPLETED", nil
		}

		if utils.StrSliceContains([]string{"Sync-Failure"}, status) {
			return respBody, status, nil
		}
		// The status value is `Applying`.
		return respBody, "PENDING", nil
	}
}

func resourceParameterConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	resp, err := GetParameterConfigurations(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving parameter configurations")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("configurations",
			flattenConfigurations(d.Get("configurations").([]interface{}), resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConfigurations(configurations []interface{}, resp interface{}) interface{} {
	if resp == nil {
		return nil
	}

	parameters := make([]map[string]interface{}, len(configurations))
	for i, v := range configurations {
		name := utils.PathSearch("name", v, "").(string)
		parameterType := utils.PathSearch("type", v, "").(string)
		express := fmt.Sprintf("configurations[?name=='%s']|[0].{values:values}", name)
		valueRaw := utils.PathSearch(express, resp, nil)
		typeExpress := fmt.Sprintf("values[?type=='%s']|[0].value", parameterType)
		// For multi-option parameters that require restarting the cluster, the input parameters are connected by commas (,),
		// the obtained values ​​are connected by semicolons (;). The backend is used to distinguish before and after parameter adjustment.
		value := strings.ReplaceAll(utils.PathSearch(typeExpress, valueRaw, "").(string), ";", ",")
		parameters[i] = map[string]interface{}{
			"name":  name,
			"value": value,
			"type":  parameterType,
		}
	}

	return parameters
}

func GetParameterConfigurations(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	group, err := getParameterGroup(client, clusterId)
	if err != nil {
		return nil, err
	}

	getHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/configurations/{configuration_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	getPath = strings.ReplaceAll(getPath, "{configuration_id}", utils.PathSearch("id", group, "").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(getResp)
}

func getParameterGroup(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v1.0/{project_id}/clusters/{cluster_id}/configurations"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	resp, err := client.Request("GET", getPath, &createOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("configurations[?type=='hiddenParameterGroup']|[0]", respBody, nil), nil
}

func resourceParameterConfigurationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		clusterId = d.Get("cluster_id").(string)
	)

	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = modifyConfigurations(ctx, client, clusterId, d.Get("configurations").([]interface{}), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceParameterConfigurationsRead(ctx, d, meta)
}

func resourceParameterConfigurationsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only used to modify the DWS cluster parameter configurations. Deleting this resource will	not clear
	the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
