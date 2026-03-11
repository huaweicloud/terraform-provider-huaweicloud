package cse

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	DefaultVersion = "CSE2"

	microserviceEngineNotFoundCodes = []string{
		"SVCSTG.00501116",
		"SVCSTG.00501125",
	}

	microserviceEngineNonUpdatableParams = []string{
		"name",
		"flavor",
		"network_id",
		"auth_type",
		"availability_zones",
		"version",
		"admin_pass",
		"enterprise_project_id",
		"description",
		"eip_id",
		"extend_params",
	}
)

// @API VPC GET /v1/{project_id}/subnets/{subnet_id}
// @API VPC GET /v1/{project_id}/vpcs/{vpc_id}
// @API CSE POST /v2/{project_id}/enginemgr/engines
// @API CSE GET /v2/{project_id}/enginemgr/engines/{engine_id}/jobs/{job_id}
// @API CSE GET /v2/{project_id}/enginemgr/engines/{engine_id}
// @API CSE DELETE /v2/{project_id}/enginemgr/engines/{engine_id}
func ResourceMicroserviceEngine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicroserviceEngineCreate,
		ReadContext:   resourceMicroserviceEngineRead,
		UpdateContext: resourceMicroserviceEngineUpdate,
		DeleteContext: resourceMicroserviceEngineDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEngineImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(microserviceEngineNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the microservice engine is located.`,
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the microservice engine.`,
			},
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The flavor of the microservice engine.`,
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The network ID of the microservice engine.`,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The authentication type of the microservice engine.`,
			},

			// Optional parameters.
			"availability_zones": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The availability zones where the microservice engine is deployed.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     DefaultVersion,
				Description: `The version of the microservice engine.`,
			},
			"admin_pass": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The administrator password of the microservice engine.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The enterprise project ID to which the microservice engine belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the microservice engine.`,
			},
			"eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The EIP ID bound to the microservice engine.`,
			},
			"extend_params": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The extended parameters of the microservice engine.`,
			},

			// Attributes.
			"service_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The service limit of the microservice engine.`,
			},
			"instance_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The instance limit of the microservice engine.`,
			},
			"service_registry_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The service registry addresses of the microservice engine.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The private address of the service registry.`,
						},
						"public": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The public address of the service registry.`,
						},
					},
				},
			},
			"config_center_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The config center addresses of the microservice engine.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The private address of the config center.`,
						},
						"public": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The public address of the config center.`,
						},
					},
				},
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func getSubnetById(client *golangsdk.ServiceClient, networkId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/subnets/{subnet_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{subnet_id}", networkId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func getVpcById(client *golangsdk.ServiceClient, vpcId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/vpcs/{vpc_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{vpc_id}", vpcId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func buildAuthCred(authType, adminPass string) interface{} {
	if authType == "RBAC" && adminPass != "" {
		return map[string]interface{}{
			"pwd": adminPass,
		}
	}
	return nil
}

// buildMicroserviceEngineCreateOpts builds the request body for creating a microservice engine.
// It extracts data from the schema, VPC info, and subnet info, then constructs the API request parameters.
// The function handles field mapping, type conversion, and optional field processing (using utils.ValueIgnoreEmpty).
func buildMicroserviceEngineCreateOpts(d *schema.ResourceData, vpcInfo, subnetInfo interface{}) map[string]interface{} {
	authType := d.Get("auth_type").(string)
	return map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": utils.ValueIgnoreEmpty(d.Get("description").(string)),
		"payment":     "1", // PostPaid
		"flavor":      d.Get("flavor").(string),
		"azList":      utils.ExpandToStringListBySet(d.Get("availability_zones").(*schema.Set)),
		"authType":    authType,
		"vpc":         utils.PathSearch("vpc.name", vpcInfo, ""),
		"vpcId":       utils.PathSearch("subnet.vpc_id", subnetInfo, ""),
		"networkId":   d.Get("network_id").(string),
		"subnetCidr":  utils.PathSearch("subnet.cidr", subnetInfo, ""),
		"publicIpId":  utils.ValueIgnoreEmpty(d.Get("eip_id").(string)),
		"auth_cred":   buildAuthCred(authType, d.Get("admin_pass").(string)),
		"specType":    d.Get("version").(string),
		"inputs":      d.Get("extend_params").(map[string]interface{}),
	}
}

func createMicroserviceEngine(cseClient, vpcClient *golangsdk.ServiceClient, d *schema.ResourceData,
	enterpriseProjectId string) (interface{}, error) {
	networkId := d.Get("network_id").(string)
	subnetInfo, err := getSubnetById(vpcClient, networkId)
	if err != nil {
		return nil, err
	}
	vpcInfo, err := getVpcById(vpcClient, utils.PathSearch("subnet.vpc_id", subnetInfo, "").(string))
	if err != nil {
		return nil, err
	}

	httpUrl := "v2/{project_id}/enginemgr/engines"
	createPath := cseClient.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", cseClient.ProjectID)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(enterpriseProjectId),
		JSONBody:         utils.RemoveNil(buildMicroserviceEngineCreateOpts(d, vpcInfo, subnetInfo)),
	}

	requestResp, err := cseClient.Request("POST", createPath, &createOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func getMicroserviceEngineJob(client *golangsdk.ServiceClient, engineId, jobId, enterpriseProjectId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/enginemgr/engines/{engine_id}/jobs/{job_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{engine_id}", engineId)
	getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(enterpriseProjectId),
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

// refreshMicroserviceEngineJobFunc returns a state refresh function for polling the microservice engine job status.
// It handles error conversion (converting 400/401 errors to 404 when the engine is not found),
// checks for failure statuses, and determines if the job has reached the target status.
// The targets parameter specifies the list of status values that indicate completion.
// If targets is empty and a 404 error occurs, it returns COMPLETED (used for delete operations).
func refreshMicroserviceEngineJobFunc(client *golangsdk.ServiceClient, engineId, jobId, enterpriseProjectId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getMicroserviceEngineJob(client, engineId, jobId, enterpriseProjectId)
		if err != nil {
			parsedErr := common.ConvertExpected400ErrInto404Err(
				common.ConvertExpected401ErrInto404Err(err, "error_code", microserviceEngineNotFoundCodes...),
				"error_code",
				microserviceEngineNotFoundCodes...,
			)
			if _, ok := parsedErr.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "RESOURCE_NOT_FOUND", "COMPLETED", nil
			}
			return nil, "ERROR", parsedErr
		}

		status := utils.PathSearch("status", resp, "").(string)
		if utils.StrSliceContains([]string{"CreateFail", "DeleteFailed", "UpgradeFailed", "ModifyFailed"}, status) {
			return resp, "ERROR", fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func resourceMicroserviceEngineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		enterpriseProjectId = cfg.GetEnterpriseProjectID(d)
	)

	cseClient, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}
	vpcClient, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	createOpts, err := createMicroserviceEngine(cseClient, vpcClient, d, enterpriseProjectId)
	if err != nil {
		return diag.Errorf("error creating microservice engine: %s", err)
	}

	engineId := utils.PathSearch("id", createOpts, "").(string)
	if engineId == "" {
		return diag.Errorf("unable to find the microservice engine ID from the API response")
	}
	d.SetId(engineId)

	log.Printf("[DEBUG] Waiting for the microservice engine to become running, the engine ID is %s", engineId)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: refreshMicroserviceEngineJobFunc(cseClient, engineId,
			strconv.Itoa(int(utils.PathSearch("jobId", createOpts, float64(0)).(float64))), enterpriseProjectId, []string{"Finished"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        180 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the creation of microservice engine (%s) to complete: %s", engineId, err)
	}

	return resourceMicroserviceEngineRead(ctx, d, meta)
}

func GetMicroserviceEngineById(client *golangsdk.ServiceClient, engineId, epsId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/enginemgr/engines/{engine_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{engine_id}", engineId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(epsId),
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func flattenServiceRegistryAddresses(entrypoint map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"private": utils.PathSearch("serviceEndpoint.serviceCenter.masterEntrypoint", entrypoint, ""),
			"public":  utils.PathSearch("publicServiceEndpoint.serviceCenter.masterEntrypoint", entrypoint, ""),
		},
	}
}

func flattenConfigAddresses(entrypoint map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"private": utils.PathSearch("serviceEndpoint.configCenter.masterEntrypoint", entrypoint, ""),
			"public":  utils.PathSearch("publicServiceEndpoint.configCenter.masterEntrypoint", entrypoint, ""),
		},
	}
}

func resourceMicroserviceEngineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		engineId            = d.Id()
		enterpriseProjectId = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	respBody, err := GetMicroserviceEngineById(client, engineId, enterpriseProjectId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(
			common.ConvertExpected401ErrInto404Err(err, "error_code", microserviceEngineNotFoundCodes...),
			"error_code",
			microserviceEngineNotFoundCodes...,
		), fmt.Sprintf("error retrieving microservice engine (%s)", engineId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, "").(string)),
		d.Set("flavor", utils.PathSearch("flavor", respBody, "").(string)),
		d.Set("availability_zones", utils.PathSearch("reference.azList", respBody, make([]interface{}, 0)).([]interface{})),
		d.Set("auth_type", utils.PathSearch("authType", respBody, "").(string)),
		d.Set("version", utils.PathSearch("specType", respBody, "").(string)),
		d.Set("enterprise_project_id", utils.PathSearch("enterpriseProjectId", respBody, "").(string)),
		d.Set("network_id", utils.PathSearch("reference.networkId", respBody, "").(string)),
		d.Set("description", utils.PathSearch("description", respBody, "").(string)),
		d.Set("eip_id", utils.PathSearch("reference.publicIpId", respBody, "").(string)),
		d.Set("extend_params", utils.PathSearch("reference.inputs", respBody, map[string]interface{}{}).(map[string]interface{})),
		d.Set("service_registry_addresses", flattenServiceRegistryAddresses(utils.PathSearch("externalEntrypoint",
			respBody, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("config_center_addresses", flattenConfigAddresses(utils.PathSearch("externalEntrypoint",
			respBody, make(map[string]interface{})).(map[string]interface{}))),
	)

	diagErr := make([]diag.Diagnostic, 0, 3)
	// Attributes
	if serviceLimit := utils.PathSearch("reference.serviceLimit", respBody, "").(string); serviceLimit != "" {
		limit, err := strconv.Atoi(serviceLimit)
		if err != nil {
			// Record and continue.
			diagErr = append(diagErr, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Wrong format",
				Detail:   fmt.Sprintf("Unable to parse the service limit (%#v).", serviceLimit),
			})
		} else {
			mErr = multierror.Append(mErr, d.Set("service_limit", limit))
		}
	}
	if instanceLimit := utils.PathSearch("reference.instanceLimit", respBody, "").(string); instanceLimit != "" {
		limit, err := strconv.Atoi(instanceLimit)
		if err != nil {
			// Record and continue.
			diagErr = append(diagErr, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Wrong format",
				Detail:   fmt.Sprintf("Unable to parse the instance limit (%#v).", instanceLimit),
			})
		} else {
			mErr = multierror.Append(mErr, d.Set("instance_limit", limit))
		}
	}

	diagErr = append(diagErr, diag.FromErr(mErr.ErrorOrNil())...)
	return diagErr
}

func resourceMicroserviceEngineUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func deleteMicroserviceEngine(client *golangsdk.ServiceClient, engineId, epsId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/enginemgr/engines/{engine_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{engine_id}", engineId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(epsId),
	}

	requestResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceMicroserviceEngineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		engineId            = d.Id()
		enterpriseProjectId = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	resp, err := deleteMicroserviceEngine(client, engineId, enterpriseProjectId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(
			common.ConvertExpected401ErrInto404Err(err, "error_code", microserviceEngineNotFoundCodes...),
			"error_code",
			microserviceEngineNotFoundCodes...,
		), fmt.Sprintf("error deleting microservice engine (%s)", engineId))
	}

	log.Printf("[DEBUG] Waiting for the Microservice engine delete complete, the engine ID is %s.", engineId)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: refreshMicroserviceEngineJobFunc(client, engineId,
			strconv.Itoa(int(utils.PathSearch("jobId", resp, float64(0)).(float64))), enterpriseProjectId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        120 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the deletion of microservice engine (%s) to complete: %s", engineId, err)
	}

	return nil
}

func resourceEngineImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	switch len(parts) {
	case 1:
		d.SetId(parts[0])
		return []*schema.ResourceData{d}, nil
	case 2:
		d.SetId(parts[0])
		return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[1])
	}
	return nil, fmt.Errorf("the imported ID specifies an invalid format: want '<id>' or "+
		"'<id>/<enterprise_project_id>', but '%s'", importedId)
}
