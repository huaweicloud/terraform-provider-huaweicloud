package cse

import (
	"context"
	"fmt"
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
	microserviceInstanceNonUpdatableParams = []string{
		"auth_address",
		"connect_address",
		"admin_user",
		"admin_pass",
		"microservice_id",
		"host_name",
		"endpoints",
		"version",
		"properties",
		"health_check",
		"health_check.*.mode",
		"health_check.*.interval",
		"health_check.*.max_retries",
		"health_check.*.port",
		"data_center",
		"data_center.*.name",
		"data_center.*.region",
		"data_center.*.availability_zone",
	}
	// The project ID of the microservice instance is the fixed value "default".
	// No region parameter needs to be defined because this resource does not use IAM authentication.
	microserviceInstanceProjectId     = "default"
	internalPropertyKeys              = []string{"engineID", "engineName"}
	microserviceInstanceNotFoundCodes = []string{
		"400017",
	}
)

// @API CSE POST /v4/token
// @API CSE POST /v4/{project_id}/registry/microservices/{service_id}/instances
// @API CSE GET /v4/{project_id}/registry/microservices/{service_id}/instances/{instance_id}
// @API CSE DELETE /v4/{project_id}/registry/microservices/{service_id}/instances/{instance_id}
func ResourceMicroserviceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicroserviceInstanceCreate,
		ReadContext:   resourceMicroserviceInstanceRead,
		UpdateContext: resourceMicroserviceInstanceUpdate,
		DeleteContext: resourceMicroserviceInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMicroserviceInstanceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(microserviceInstanceNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Special parameters.
			// These parameters are used to specify the address that used to request the access token and access the
			// microservice engine.
			"auth_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The address that used to request the access token.`,
					utils.SchemaDescInput{
						Required: true,
					}),
			},
			"connect_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The address that used to access engine and manages microservice instance.`,
			},
			"admin_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user name that used to pass the RBAC control.",
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"admin_user"},
				Description:  `The user password that used to pass the RBAC control.`,
			},

			// Required parameters.
			"microservice_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the microservice to which the microservice instance belongs.`,
			},
			"host_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The host name of the microservice instance.`,
			},
			"endpoints": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of access addresses of the microservice instance.`,
			},
			// Optional parameters.
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The version of the microservice instance.`,
			},
			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The extended attributes of the microservice instance, in key/value format.`,
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The heartbeat mode of the health check.`,
						},
						"interval": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The heartbeat interval of the health check, in seconds.`,
						},
						"max_retries": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The maximum retry number of the health check.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The port of the health check.`,
						},
					},
				},
				Description: `The health check configuration of the microservice instance.`,
			},
			"data_center": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the data center.`,
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The custom region name of the data center.`,
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The custom availability zone of the data center.`,
						},
					},
				},
				Description: `The data center configuration of the microservice instance.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the microservice instance.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func buildInstanceHealthCheck(healthChecks []interface{}) map[string]interface{} {
	if len(healthChecks) < 1 {
		return nil
	}

	healthCheck := healthChecks[0].(map[string]interface{})

	return map[string]interface{}{
		"mode":     utils.PathSearch("mode", healthCheck, nil),
		"interval": utils.PathSearch("interval", healthCheck, nil),
		"times":    utils.PathSearch("max_retries", healthCheck, nil),
		"port":     utils.ValueIgnoreEmpty(utils.PathSearch("port", healthCheck, nil)),
	}
}

func buildInstanceDataCenter(dataCenters []interface{}) map[string]interface{} {
	if len(dataCenters) < 1 {
		return nil
	}

	dataCenter := dataCenters[0].(map[string]interface{})

	return map[string]interface{}{
		"name":          utils.PathSearch("name", dataCenter, nil),
		"region":        utils.PathSearch("region", dataCenter, nil),
		"availableZone": utils.PathSearch("availability_zone", dataCenter, nil),
	}
}

// buildCustomProperties filters out internal property keys (engineID, engineName) from the properties map.
func buildInstanceProperties(properties map[string]interface{}) map[string]interface{} {
	if len(properties) < 1 {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range properties {
		if !utils.StrSliceContains(internalPropertyKeys, k) {
			result[k] = v
		}
	}

	return result
}

func buildInstanceCreateOpts(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"instance": map[string]interface{}{
			// Required parameters.
			"hostName":  d.Get("host_name").(string),
			"endpoints": d.Get("endpoints").([]interface{}),
			// Optional parameters.
			"version":        utils.ValueIgnoreEmpty(d.Get("version").(string)),
			"properties":     buildInstanceProperties(d.Get("properties").(map[string]interface{})),
			"healthCheck":    buildInstanceHealthCheck(d.Get("health_check").([]interface{})),
			"dataCenterInfo": buildInstanceDataCenter(d.Get("data_center").([]interface{})),
		},
	})
}

func createInstance(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl        = "v4/{project_id}/registry/microservices/{service_id}/instances"
		microserviceId = d.Get("microservice_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", microserviceInstanceProjectId)
	createPath = strings.ReplaceAll(createPath, "{service_id}", microserviceId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(client.ProjectID),
		JSONBody:         buildInstanceCreateOpts(d),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `POST /v4/token` interface.
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string), d.Get("admin_pass").(string))
	if err != nil {
		return nil, err
	}
	// If the microservice instance has RBAC authentication enabled, the Authorization header will use a special token
	// provided by the CSE service to replace the original IAM authentication information (AKSK authentication) in the
	// request header.
	if token != "" {
		createOpts.MoreHeaders["Authorization"] = token
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceMicroserviceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Creating a microservice instance in the microservice engine requires building a client based on the microservice
	// engine's connection address, which does not use IAM authentication.
	client := common.NewCustomClient(true, d.Get("connect_address").(string))

	resp, err := createInstance(client, d)
	if err != nil {
		return diag.Errorf("error creating microservice instance: %s", err)
	}

	instanceId := utils.PathSearch("instanceId", resp, "").(string)
	if instanceId == "" {
		return diag.Errorf("unable to find the instance ID from the API response")
	}
	d.SetId(instanceId)

	return resourceMicroserviceInstanceRead(ctx, d, meta)
}

// GetInstance retrieves the microservice instance details by authorization parameters, microservice ID, and instance ID.
func GetInstance(client *golangsdk.ServiceClient, authAddress, adminUser, adminPass, microserviceId, instanceId string) (interface{}, error) {
	httpUrl := "v4/{project_id}/registry/microservices/{service_id}/instances/{instance_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", microserviceInstanceProjectId)
	getPath = strings.ReplaceAll(getPath, "{service_id}", microserviceId)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(client.ProjectID),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `POST /v4/token` interface.
	token, err := GetAuthorizationToken(authAddress, adminUser, adminPass)
	if err != nil {
		return nil, err
	}
	// If the microservice instance has RBAC authentication enabled, the Authorization header will use a special token
	// provided by the CSE service to replace the original IAM authentication information (AKSK authentication) in the
	// request header.
	if token != "" {
		getOpts.MoreHeaders["Authorization"] = token
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func flattenHealthCheck(healthCheck interface{}) []map[string]interface{} {
	if healthCheck == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"mode":        utils.PathSearch("mode", healthCheck, nil),
			"interval":    utils.PathSearch("interval", healthCheck, nil),
			"max_retries": utils.PathSearch("times", healthCheck, nil),
			"port":        utils.PathSearch("port", healthCheck, nil),
		},
	}
}

func flattenDataCenter(dataCenter interface{}) []map[string]interface{} {
	if dataCenter == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":              utils.PathSearch("name", dataCenter, nil),
			"region":            utils.PathSearch("region", dataCenter, nil),
			"availability_zone": utils.PathSearch("availableZone", dataCenter, nil),
		},
	}
}

func resourceMicroserviceInstanceRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var (
		// Creating a microservice instance in the microservice engine requires building a client based on the
		// microservice engine's connection address, which does not use IAM authentication.
		client         = common.NewCustomClient(true, d.Get("connect_address").(string))
		authAddress    = getAuthAddress(d)
		adminUser      = d.Get("admin_user").(string)
		adminPass      = d.Get("admin_pass").(string)
		microserviceId = d.Get("microservice_id").(string)
		instanceId     = d.Id()
	)

	respBody, err := GetInstance(client, authAddress, adminUser, adminPass, microserviceId, instanceId)
	if err != nil {
		// When the microservice instance does not exist, the error code returned is 400, and the error information is:
		// {"errorCode": "400017", "errorMessage": "Instance does not exist", "detail": "... instance does not exist"}
		// So, we need to convert the error code to 404 error and specify the key of the error code to be "errorCode".
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", microserviceInstanceNotFoundCodes...),
			fmt.Sprintf("error retrieving microservice instance (%s)", instanceId))
	}

	mErr := multierror.Append(nil,
		d.Set("host_name", utils.PathSearch("instance.hostName", respBody, nil)),
		d.Set("endpoints", utils.PathSearch("instance.endpoints", respBody, nil)),
		d.Set("version", utils.PathSearch("instance.version", respBody, nil)),
		d.Set("properties", buildInstanceProperties(utils.PathSearch("instance.properties", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("health_check", flattenHealthCheck(utils.PathSearch("instance.healthCheck", respBody, nil))),
		d.Set("data_center", flattenDataCenter(utils.PathSearch("instance.dataCenterInfo", respBody, nil))),
		d.Set("status", utils.PathSearch("instance.status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMicroserviceInstanceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func deleteInstance(client *golangsdk.ServiceClient, authAddress, adminUser, adminPass, microserviceId, instanceId string) error {
	httpUrl := "v4/{project_id}/registry/microservices/{service_id}/instances/{instance_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", microserviceInstanceProjectId)
	deletePath = strings.ReplaceAll(deletePath, "{service_id}", microserviceId)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(client.ProjectID),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `POST /v4/token` interface.
	token, err := GetAuthorizationToken(authAddress, adminUser, adminPass)
	if err != nil {
		return err
	}
	// If the microservice instance has RBAC authentication enabled, the Authorization header will use a special token
	// provided by the CSE service to replace the original IAM authentication information (AKSK authentication) in the
	// request header.
	if token != "" {
		deleteOpts.MoreHeaders["Authorization"] = token
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return err
	}
	return nil
}

func instanceStatusRefreshFunc(client *golangsdk.ServiceClient, authAddress, adminUser, adminPass,
	microserviceId, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetInstance(client, authAddress, adminUser, adminPass, microserviceId, instanceId)
		if err != nil {
			// Convert 400 error to 404 error if it matches instance not found codes.
			convertedErr := common.ConvertExpected400ErrInto404Err(err, "errorCode", microserviceInstanceNotFoundCodes...)
			if _, ok := convertedErr.(golangsdk.ErrDefault404); ok {
				// When the error code is 404 (or converted to 404), the instance has been deleted.
				return "Resource Not Found", "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		// Instance still exists, continue polling.
		return respBody, "PENDING", nil
	}
}

func resourceMicroserviceInstanceDelete(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var (
		// Creating a microservice instance in the microservice engine requires building a client based on the
		// microservice engine's connection address, which does not use IAM authentication.
		client         = common.NewCustomClient(true, d.Get("connect_address").(string))
		authAddress    = getAuthAddress(d)
		adminUser      = d.Get("admin_user").(string)
		adminPass      = d.Get("admin_pass").(string)
		microserviceId = d.Get("microservice_id").(string)
		instanceId     = d.Id()
	)

	err := deleteInstance(client, authAddress, adminUser, adminPass, microserviceId, instanceId)
	if err != nil {
		// When the microservice instance does not exist, the error code returned is 400, and the error information is:
		// {"errorCode": "400017", "errorMessage": "Instance does not exist", "detail": "... instance does not exist"}
		// So, we need to convert the error code to 404 error and specify the key of the error code to be "errorCode".
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", microserviceInstanceNotFoundCodes...),
			fmt.Sprintf("error deleting microservice instance (%s)", instanceId))
	}

	// Wait for the instance to be deleted.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      instanceStatusRefreshFunc(client, authAddress, adminUser, adminPass, microserviceId, instanceId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the microservice instance (%s) to be deleted: %s", instanceId, err)
	}

	return nil
}

func resourceMicroserviceInstanceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	var (
		authAddr, connectAddr, idWithoutAddrs, microserviceId, instanceId, adminUser, adminPwd string
		mErr                                                                                   *multierror.Error

		importedId = d.Id()
		formatErr  = fmt.Errorf("the imported microservice ID specifies an invalid format, want "+
			"'<auth_address>/<connect_address>/<microservice_id>/<id>' or "+
			"'<auth_address>/<connect_address>/<microservice_id>/<id>/<admin_user>/<admin_pass>', but got '%s'",
			importedId)
	)

	if !connectAddressRegexPattern.MatchString(importedId) {
		return nil, formatErr
	}
	resp := connectAddressRegexPattern.FindAllStringSubmatch(importedId, -1)
	// To prevent panics caused by null resp values ​​due to differences in regular expression implementation, defensive
	// checks are added.
	if len(resp) == 0 {
		return nil, formatErr
	}
	// If the imported ID matches the address regular expression, the length of the response result must be greater than 1.
	switch len(resp[0]) {
	case 4:
		authAddr = resp[0][1]
		connectAddr = resp[0][2]
		idWithoutAddrs = resp[0][3]
		if authAddr == "" {
			authAddr = connectAddr // Using the connect address as the auth address if the auth address input is omitted.
		}
	default:
		return nil, formatErr
	}

	mErr = multierror.Append(mErr,
		d.Set("auth_address", authAddr),
		d.Set("connect_address", connectAddr),
	)

	parts := strings.Split(idWithoutAddrs, "/")
	switch len(parts) {
	case 2:
		microserviceId = parts[0]
		instanceId = parts[1]
	case 4:
		microserviceId = parts[0]
		instanceId = parts[1]
		adminUser = parts[2]
		adminPwd = parts[3]

		mErr = multierror.Append(mErr,
			d.Set("admin_user", adminUser),
			d.Set("admin_pass", adminPwd),
		)
	default:
		return nil, formatErr
	}

	mErr = multierror.Append(mErr, d.Set("microservice_id", microserviceId))
	d.SetId(instanceId)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
