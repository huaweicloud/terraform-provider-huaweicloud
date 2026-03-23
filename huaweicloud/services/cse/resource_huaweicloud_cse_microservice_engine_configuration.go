package cse

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const microserviceEngineConfigurationDefaultProjectId = "default"

var microserviceEngineConfigurationNonUpdatableParams = []string{
	"auth_address",
	"connect_address",
	"admin_user",
	"admin_pass",
	"key",
	"value_type",
	"enterprise_project_id",
	"tags",
}

// @API CSE POST /v4/token
// @API CSE POST /v1/{project_id}/kie/kv
// @API CSE GET /v1/{project_id}/kie/kv/{kv_id}
// @API CSE PUT /v1/{project_id}/kie/kv/{kv_id}
// @API CSE DELETE /v1/{project_id}/kie/kv/{kv_id}
// @API CSE GET /v1/{project_id}/kie/kv
func ResourceMicroserviceEngineConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicroserviceEngineConfigurationCreate,
		ReadContext:   resourceMicroserviceEngineConfigurationRead,
		UpdateContext: resourceMicroserviceEngineConfigurationUpdate,
		DeleteContext: resourceMicroserviceEngineConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMicroserviceEngineConfigurationImportState,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(microserviceEngineConfigurationNonUpdatableParams),
			config.MergeDefaultTags(),
		),

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
				Description: `The address that used to access the engine and manage configuration.`,
			},
			"admin_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user name that used to pass the RBAC control.`,
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"admin_user"},
				Description:  `The user password that used to pass the RBAC control.`,
			},

			// Required parameters.
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The configuration key (item name).`,
			},
			"value_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the configuration value.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The configuration value.`,
			},

			// Optional parameters.
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The configuration status.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The enterprise project ID to which the microservice engine configuration belongs.`,
			},
			"tags": common.TagsSchema(
				`The key/value pairs to associate with the configuration.`,
			),

			// Attributes.
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the configuration, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the configuration, in RFC3339 format.`,
			},
			"create_revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The create version of the configuration.`,
			},
			"update_revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The update version of the configuration.`,
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

func buildMicroserviceEngineConfigurationCreateOpts(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"key":        d.Get("key").(string),
		"value_type": d.Get("value_type").(string),
		"value":      d.Get("value").(string),
		"labels":     d.Get("tags").(map[string]interface{}),
		"status":     utils.ValueIgnoreEmpty(d.Get("status").(string)),
	})
}

func createMicroserviceEngineConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData,
	authInfo MicroserviceEngineAuthInfo) (interface{}, error) {
	var (
		httpUrl     = "v1/{project_id}/kie/kv"
		createPath  = client.Endpoint + httpUrl
		authAddress = getAuthAddress(d)
		adminUser   = d.Get("admin_user").(string)
		adminPass   = d.Get("admin_pass").(string)
	)

	// The project ID of the microservices is the fixed value "default".
	// No region parameter needs to be defined because this data source does not use IAM authentication.
	createPath = strings.ReplaceAll(createPath, "{project_id}", microserviceEngineConfigurationDefaultProjectId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(authInfo.EnterpriseProjectId),
		JSONBody:         buildMicroserviceEngineConfigurationCreateOpts(d),
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
		createOpts.MoreHeaders["Authorization"] = token
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return "", err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceMicroserviceEngineConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
		// Creating a microservice engine configuration requires building a client based on the
		// microservice engine's connection address, which does not use IAM authentication.
		client                     = common.NewCustomClient(true, d.Get("connect_address").(string))
		microserviceEngineAuthInfo = MicroserviceEngineAuthInfo{
			AuthAddress:         getAuthAddress(d),
			AdminUser:           d.Get("admin_user").(string),
			AdminPass:           d.Get("admin_pass").(string),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}
	)

	respBody, err := createMicroserviceEngineConfiguration(client, d, microserviceEngineAuthInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	configId := utils.PathSearch("id", respBody, "").(string)
	if configId == "" {
		return diag.Errorf("unable to find the CSE microservice configuration ID from the API response")
	}
	d.SetId(configId)

	return resourceMicroserviceEngineConfigurationRead(ctx, d, meta)
}

// formatKieKvTime converts KIE KV create_time/update_time (Unix milliseconds as float64) to RFC3339.
func formatKieKvTime(val interface{}) string {
	if val == nil {
		return ""
	}

	f, ok := val.(float64)
	if !ok || f == 0 {
		return ""
	}

	return utils.FormatTimeStampRFC3339(int64(f)/1000, false)
}

// GetMicroserviceEngineConfiguration queries a single KIE KV configuration by ID.
func GetMicroserviceEngineConfiguration(client *golangsdk.ServiceClient, authInfo MicroserviceEngineAuthInfo,
	configId string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/kie/kv/{kv_id}"
		getPath = client.Endpoint + httpUrl
	)

	// The project ID of the microservice engine configuration is the fixed value "default".
	// No region parameter needs to be defined because this data source does not use IAM authentication.
	getPath = strings.ReplaceAll(getPath, "{project_id}", microserviceEngineConfigurationDefaultProjectId)
	getPath = strings.ReplaceAll(getPath, "{kv_id}", configId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(authInfo.EnterpriseProjectId),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `POST /v4/token` interface.
	token, err := GetAuthorizationToken(authInfo.AuthAddress, authInfo.AdminUser, authInfo.AdminPass)
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

func resourceMicroserviceEngineConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
		// Creating a microservice engine configuration requires building a client based on the
		// microservice engine's connection address, which does not use IAM authentication.
		client                     = common.NewCustomClient(true, d.Get("connect_address").(string))
		microserviceEngineAuthInfo = MicroserviceEngineAuthInfo{
			AuthAddress:         getAuthAddress(d),
			AdminUser:           d.Get("admin_user").(string),
			AdminPass:           d.Get("admin_pass").(string),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}
		configId = d.Id()
	)

	respBody, err := GetMicroserviceEngineConfiguration(client, microserviceEngineAuthInfo, configId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying configuration (%s)", configId))
	}

	mErr := multierror.Append(nil,
		d.Set("key", utils.PathSearch("key", respBody, nil)),
		d.Set("value_type", utils.PathSearch("value_type", respBody, nil)),
		d.Set("value", utils.PathSearch("value", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("enterprise_project_id", microserviceEngineAuthInfo.EnterpriseProjectId),
		d.Set("tags", utils.PathSearch("labels", respBody, nil)),
		d.Set("created_at", formatKieKvTime(utils.PathSearch("create_time", respBody, nil))),
		d.Set("updated_at", formatKieKvTime(utils.PathSearch("update_time", respBody, nil))),
		d.Set("create_revision", utils.PathSearch("create_revision", respBody, nil)),
		d.Set("update_revision", utils.PathSearch("update_revision", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildMicroserviceEngineConfigurationUpdateOpts(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"value":  d.Get("value").(string),
		"status": utils.ValueIgnoreEmpty(d.Get("status").(string)),
	})
}

func updateMicroserviceEngineConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData, authInfo MicroserviceEngineAuthInfo,
	configId string) error {
	var (
		httpUrl    = "v1/{project_id}/kie/kv/{kv_id}"
		updatePath = client.Endpoint + httpUrl
	)

	// The project ID of the microservice engine configuration is the fixed value "default".
	// No region parameter needs to be defined because this data source does not use IAM authentication.
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", microserviceEngineConfigurationDefaultProjectId)
	updatePath = strings.ReplaceAll(updatePath, "{kv_id}", configId)

	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(authInfo.EnterpriseProjectId),
		JSONBody:         buildMicroserviceEngineConfigurationUpdateOpts(d),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `POST /v4/token` interface.
	token, err := GetAuthorizationToken(authInfo.AuthAddress, authInfo.AdminUser, authInfo.AdminPass)
	if err != nil {
		return err
	}
	// If the microservice instance has RBAC authentication enabled, the Authorization header will use a special token
	// provided by the CSE service to replace the original IAM authentication information (AKSK authentication) in the
	// request header.
	if token != "" {
		updateOpts.MoreHeaders["Authorization"] = token
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	return err
}

func resourceMicroserviceEngineConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
		// Creating a microservice engine configuration requires building a client based on the
		// microservice engine's connection address, which does not use IAM authentication.
		client                     = common.NewCustomClient(true, d.Get("connect_address").(string))
		microserviceEngineAuthInfo = MicroserviceEngineAuthInfo{
			AuthAddress:         getAuthAddress(d),
			AdminUser:           d.Get("admin_user").(string),
			AdminPass:           d.Get("admin_pass").(string),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}
	)

	if err := updateMicroserviceEngineConfiguration(client, d, microserviceEngineAuthInfo, d.Id()); err != nil {
		return diag.Errorf("error updating configuration (%s): %s", d.Id(), err)
	}

	return resourceMicroserviceEngineConfigurationRead(ctx, d, meta)
}

func deleteMicroserviceEngineConfiguration(client *golangsdk.ServiceClient, authInfo MicroserviceEngineAuthInfo, configId string) error {
	var (
		httpUrl    = "v1/{project_id}/kie/kv/{kv_id}"
		deletePath = client.Endpoint + httpUrl
	)

	// The project ID of the microservice engine configuration is the fixed value "default".
	// No region parameter needs to be defined because this data source does not use IAM authentication.
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", microserviceEngineConfigurationDefaultProjectId)
	deletePath = strings.ReplaceAll(deletePath, "{kv_id}", configId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(authInfo.EnterpriseProjectId),
	}
	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `POST /v4/token` interface.
	token, err := GetAuthorizationToken(authInfo.AuthAddress, authInfo.AdminUser, authInfo.AdminPass)
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
	return err
}

func resourceMicroserviceEngineConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
		// Creating a microservice engine configuration requires building a client based on the
		// microservice engine's connection address, which does not use IAM authentication.
		client                     = common.NewCustomClient(true, d.Get("connect_address").(string))
		microserviceEngineAuthInfo = MicroserviceEngineAuthInfo{
			AuthAddress:         getAuthAddress(d),
			AdminUser:           d.Get("admin_user").(string),
			AdminPass:           d.Get("admin_pass").(string),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}
		configId = d.Id()
	)

	if err := deleteMicroserviceEngineConfiguration(client, microserviceEngineAuthInfo, configId); err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting configuration (%s)", configId))
	}

	return nil
}

func getMicroserviceEngineConfigurationByKey(client *golangsdk.ServiceClient, authInfo MicroserviceEngineAuthInfo,
	keyName string) (interface{}, error) {
	configurations, err := listMicroserviceEngineConfigurations(client, authInfo)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch(fmt.Sprintf("[?key=='%s']|[0]", keyName), configurations, nil), nil
}

func resourceMicroserviceEngineConfigurationImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	var (
		authAddr, connectAddr, keyName, adminUser, adminPwd, enterpriseProjectId string
		mErr                                                                     *multierror.Error

		importedId   = d.Id()
		addressRegex = `https://\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`
		re           = regexp.MustCompile(fmt.Sprintf(`^(%[1]s)?/?(%[1]s)/(.*)$`, addressRegex))
		formatErr    = fmt.Errorf("the imported microservice ID specifies an invalid format, want "+
			"'<auth_address>/<connect_address>/<key>' or '<auth_address>/<connect_address>/<key>/<admin_user>/<admin_pass>', but got '%s'",
			importedId)
	)

	if !re.MatchString(d.Id()) {
		return nil, formatErr
	}

	resp := re.FindAllStringSubmatch(d.Id(), -1)
	// If the imported ID matches the address regular expression, the length of the response result must be greater than 1.
	if len(resp[0]) == 4 {
		authAddr = resp[0][1]
		connectAddr = resp[0][2]
		mErr = multierror.Append(mErr,
			d.Set("auth_address", resp[0][1]),
			d.Set("connect_address", resp[0][2]),
		)

		parts := strings.Split(resp[0][3], "/")
		switch len(parts) {
		case 1:
			keyName = parts[0]
		case 2:
			keyName = parts[0]
			enterpriseProjectId = parts[1]
		case 3:
			keyName = parts[0]
			adminUser = parts[1]
			adminPwd = parts[2]
		case 4:
			keyName = parts[0]
			adminUser = parts[1]
			adminPwd = parts[2]
			enterpriseProjectId = parts[3]
		default:
			return nil, formatErr
		}

		client := common.NewCustomClient(true, connectAddr)
		authInfo := MicroserviceEngineAuthInfo{
			AuthAddress:         authAddr,
			AdminUser:           adminUser,
			AdminPass:           adminPwd,
			EnterpriseProjectId: enterpriseProjectId,
		}
		engineDetail, lookupErr := getMicroserviceEngineConfigurationByKey(client, authInfo, keyName)
		if lookupErr != nil {
			return nil, lookupErr
		}

		d.SetId(utils.PathSearch("id", engineDetail, "").(string))
		mErr = multierror.Append(mErr,
			d.Set("auth_address", authAddr),
			d.Set("connect_address", connectAddr),
			d.Set("key", keyName),
			d.Set("admin_user", adminUser),
			d.Set("admin_pass", adminPwd),
			d.Set("enterprise_project_id", enterpriseProjectId),
		)
		return []*schema.ResourceData{d}, mErr.ErrorOrNil()
	}

	return nil, formatErr
}
