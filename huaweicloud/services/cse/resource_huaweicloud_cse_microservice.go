package cse

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	microserviceNonUpdatableParams = []string{
		"auth_address",
		"connect_address",
		"admin_user",
		"admin_pass",
		"name",
		"app_name",
		"version",
		"environment",
		"level",
		"description",
	}
	// The project ID of the microservice instance is the fixed value "default".
	// No region parameter needs to be defined because this resource does not use IAM authentication.
	microserviceDefaultProjectId = "default"
	microserviceNotFoundCodes    = []string{
		"400012",
	}
)

// @API CSE GET /v4/token
// @API CSE POST /v4/{project_id}/registry/microservices
// @API CSE GET /v4/{project_id}/registry/microservices/{service_id}
// @API CSE DELETE /v4/{project_id}/registry/microservices/{service_id}
func ResourceMicroservice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicroserviceCreate,
		ReadContext:   resourceMicroserviceRead,
		UpdateContext: resourceMicroserviceUpdate,
		DeleteContext: resourceMicroserviceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMicroserviceImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(microserviceNonUpdatableParams),

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
				Description: `The address that used to access engine and manages microservice.`,
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
				Description:  "The user password that used to pass the RBAC control.",
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the microservice.",
			},
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the microservice application.",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version of the microservice.",
			},

			// Optional parameters.
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The environment (stage) type of the microservice.",
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The level of the microservice.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the microservice.",
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the microservice.",
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

func getAuthAddress(d *schema.ResourceData) string {
	if v, ok := d.GetOk("auth_address"); ok {
		return v.(string)
	}
	// Using the connect address as the auth address if its empty.
	// The behavior of the connect address is required.
	return d.Get("connect_address").(string)
}

func buildMicroserviceCreateOpts(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"service": map[string]interface{}{
			// Required parameters.
			"serviceName": d.Get("name").(string),
			"appId":       d.Get("app_name").(string),
			"version":     d.Get("version").(string),
			// Optional parameters.
			"environment": utils.ValueIgnoreEmpty(d.Get("environment").(string)),
			"level":       utils.ValueIgnoreEmpty(d.Get("level").(string)),
			"description": utils.ValueIgnoreEmpty(d.Get("description").(string)),
		},
	})
}

func createMicroservice(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v4/{project_id}/registry/microservices"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", microserviceDefaultProjectId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(client.ProjectID),
		JSONBody:         buildMicroserviceCreateOpts(d),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `GET /v4/token` interface.
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string), d.Get("admin_pass").(string))
	if err != nil {
		return nil, err
	}
	// If the microservice has RBAC authentication enabled, the Authorization header will use a special token provided
	// by the CSE service to replace the original IAM authentication information (AKSK authentication) in the request
	// header.
	if token != "" {
		createOpts.MoreHeaders["Authorization"] = token
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceMicroserviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := common.NewCustomClient(true, d.Get("connect_address").(string))

	resp, err := createMicroservice(client, d)
	if err != nil {
		return diag.Errorf("error creating microservice: %s", err)
	}

	microserviceId := utils.PathSearch("serviceId", resp, "").(string)
	if microserviceId == "" {
		return diag.Errorf("unable to find the microservice ID from the API response")
	}
	d.SetId(microserviceId)

	return resourceMicroserviceRead(ctx, d, meta)
}

func GetMicroservice(client *golangsdk.ServiceClient, authAddress, adminUser, adminPass, microserviceId string) (interface{}, error) {
	httpUrl := "v4/{project_id}/registry/microservices/{service_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", microserviceDefaultProjectId)
	getPath = strings.ReplaceAll(getPath, "{service_id}", microserviceId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(client.ProjectID),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `GET /v4/token` interface.
	token, err := GetAuthorizationToken(authAddress, adminUser, adminPass)
	if err != nil {
		return nil, err
	}
	// If the microservice has RBAC authentication enabled, the Authorization header will use a special token provided
	// by the CSE service to replace the original IAM authentication information (AKSK authentication) in the request
	// header.
	if token != "" {
		getOpts.MoreHeaders["Authorization"] = token
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceMicroserviceRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var (
		client         = common.NewCustomClient(true, d.Get("connect_address").(string))
		authAddress    = getAuthAddress(d)
		adminUser      = d.Get("admin_user").(string)
		adminPass      = d.Get("admin_pass").(string)
		microserviceId = d.Id()
	)

	respBody, err := GetMicroservice(client, authAddress, adminUser, adminPass, microserviceId)
	if err != nil {
		// When the microservice does not exist, the error code returned is 400, and the error information is:
		// {"errorCode": "400012", "errorMessage": "Micro-service does not exist", "detail": "Service does not exist."}
		// So, we need to convert the error code to 404 error and specify the key of the error code to be "errorCode".
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", microserviceNotFoundCodes...),
			fmt.Sprintf("error retrieving microservice (%s)", microserviceId))
	}

	mErr := multierror.Append(nil,
		// Required parameters.
		d.Set("name", utils.PathSearch("service.serviceName", respBody, "").(string)),
		d.Set("app_name", utils.PathSearch("service.appId", respBody, "").(string)),
		d.Set("version", utils.PathSearch("service.version", respBody, "").(string)),
		// Optional parameters.
		d.Set("environment", utils.PathSearch("service.environment", respBody, "").(string)),
		d.Set("level", utils.PathSearch("service.level", respBody, "").(string)),
		d.Set("description", utils.PathSearch("service.description", respBody, "").(string)),
		// Attributes.
		d.Set("status", utils.PathSearch("service.status", respBody, "").(string)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMicroserviceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func deleteMicroservice(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl        = "v4/{project_id}/registry/microservices/{service_id}"
		microserviceId = d.Id()
	)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", microserviceDefaultProjectId)
	deletePath = strings.ReplaceAll(deletePath, "{service_id}", microserviceId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(client.ProjectID),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `GET /v4/token` interface.
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string), d.Get("admin_pass").(string))
	if err != nil {
		return err
	}
	// If the microservice has RBAC authentication enabled, the Authorization header will use a special token provided
	// by the CSE service to replace the original IAM authentication information (AKSK authentication) in the request
	// header.
	if token != "" {
		deleteOpts.MoreHeaders["Authorization"] = token
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return err
	}
	return nil
}

func resourceMicroserviceDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := common.NewCustomClient(true, d.Get("connect_address").(string))

	// The current configuration is force deletion that delete microservices and related configuration and binding
	// instances
	err := deleteMicroservice(client, d)
	if err != nil {
		// When the microservice does not exist, the error code returned is 400, and the error information is:
		// {"errorCode": "400012", "errorMessage": "Micro-service does not exist", "detail": "Service does not exist."}
		// So, we need to convert the error code to 404 error and specify the key of the error code to be "errorCode".
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", microserviceNotFoundCodes...),
			fmt.Sprintf("error deleting microservice (%s)", d.Id()))
	}

	return nil
}

func resourceMicroserviceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	var (
		authAddr, connectAddr, importedIdWithoutAddrs, microserviceId, adminUser, adminPwd string
		mErr                                                                               *multierror.Error

		importedId = d.Id()
		formatErr  = fmt.Errorf("the imported microservice ID specifies an invalid format, want "+
			"'<auth_address>/<connect_address>/<id>' or '<auth_address>/<connect_address>/<id>/<admin_user>/<admin_pass>', but got '%s'",
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
	// If the imported ID matches the address regular expression, the length of the response result must be greater
	// than 1.
	switch len(resp[0]) {
	case 4:
		authAddr = resp[0][1]
		connectAddr = resp[0][2]
		importedIdWithoutAddrs = resp[0][3]
		if authAddr == "" {
			// Using the connect address as the auth address if the auth address input is omitted.
			authAddr = connectAddr
		}
	default:
		return nil, formatErr
	}

	mErr = multierror.Append(mErr,
		d.Set("auth_address", authAddr),
		d.Set("connect_address", connectAddr),
	)

	parts := strings.Split(importedIdWithoutAddrs, "/")
	switch len(parts) {
	case 1:
		microserviceId = parts[0]
	case 3:
		microserviceId = parts[0]
		adminUser = parts[1]
		adminPwd = parts[2]

		mErr = multierror.Append(mErr,
			d.Set("admin_user", adminUser),
			d.Set("admin_pass", adminPwd),
		)
	default:
		return nil, formatErr
	}

	d.SetId(microserviceId)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
