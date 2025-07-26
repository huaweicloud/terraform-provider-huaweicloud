package cse

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			// Authentication and request parameters.
			"auth_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The address that used to request the access token.`,
			},
			"connect_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The address that used to send requests and manage configuration.`,
			},
			"admin_user": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The user name that used to pass the RBAC control.",
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				RequiredWith: []string{"admin_user"},
				Description:  "The user password that used to pass the RBAC control.",
			},
			// Resource parameters.
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The configuration key (item name).",
			},
			"value_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of the configuration value.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The configuration value.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The configuration status.`,
			},
			"tags": common.TagsForceNewSchema(
				`The key/value pairs to associate with the configuration that used to filter resource.`,
			),
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
		},
	}
}

func buildMicroserviceEngineConfigurationCreateOpts(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"key":        d.Get("key").(string),
		"value_type": d.Get("value_type").(string),
		"value":      d.Get("value").(string),
		"status":     d.Get("status").(string),
		"labels":     d.Get("tags").(map[string]interface{}),
	}
}

func resourceMicroserviceEngineConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		client     = common.NewCustomClient(true, d.Get("connect_address").(string), "v1", "default")
		httpUrl    = "kie/kv"
		createPath = client.ResourceBase + httpUrl
	)

	// The connection address is no longer used to obtain an authentication token.
	token, err := GetAuthorizationToken(d.Get("auth_address").(string), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildMicroserviceEngineConfigurationCreateOpts(d),
	}
	if token != "" {
		createOpts.MoreHeaders = map[string]string{
			"Authorization": token,
		}
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating configuration: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	configId := utils.PathSearch("id", respBody, "").(string)
	if configId == "" {
		return diag.Errorf("enable to find the CSE microservice configuration ID from the API response")
	}
	d.SetId(configId)

	return resourceMicroserviceEngineConfigurationRead(ctx, d, meta)
}

func QueryMicroserviceEngineConfiguration(client *golangsdk.ServiceClient, token, configId string) (interface{}, error) {
	httpUrl := "kie/kv/{kv_id}"

	queryPath := client.ResourceBase + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{kv_id}", configId)

	queryOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	if token != "" {
		queryOpts.MoreHeaders = map[string]string{
			"Authorization": token,
		}
	}

	requestResp, err := client.Request("GET", queryPath, &queryOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceMicroserviceEngineConfigurationRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var (
		client   = common.NewCustomClient(true, d.Get("connect_address").(string), "v1", "default")
		configId = d.Id()
	)
	token, err := GetAuthorizationToken(d.Get("auth_address").(string), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		// When the engine does not exist, obtaining a token will cause a request connection exception.
		// To ensure that the resource is available on RFS platform, this situation is specially handled as a 404 error.
		log.Printf("[ERROR] %s", err)
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	respBody, err := QueryMicroserviceEngineConfiguration(client, token, configId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying configuration (%s)", configId))
	}

	mErr := multierror.Append(nil,
		d.Set("key", utils.PathSearch("key", respBody, nil)),
		d.Set("value_type", utils.PathSearch("value_type", respBody, nil)),
		d.Set("value", utils.PathSearch("value", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("tags", utils.PathSearch("labels", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("create_revision", utils.PathSearch("create_revision", respBody, nil)),
		d.Set("update_revision", utils.PathSearch("update_revision", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildMicroserviceEngineConfigurationUpdateOpts(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"value":  d.Get("value").(string),
		"status": d.Get("status").(string),
	}
}

func resourceMicroserviceEngineConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		client   = common.NewCustomClient(true, d.Get("connect_address").(string), "v1", "default")
		httpUrl  = "kie/kv/{kv_id}"
		configId = d.Id()
	)

	updatePath := client.ResourceBase + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{kv_id}", configId)

	token, err := GetAuthorizationToken(d.Get("auth_address").(string), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildMicroserviceEngineConfigurationUpdateOpts(d),
	}
	if token != "" {
		updateOpts.MoreHeaders = map[string]string{
			"Authorization": token,
		}
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating configuration (%s): %s", configId, err)
	}
	return resourceMicroserviceEngineConfigurationRead(ctx, d, meta)
}

func resourceMicroserviceEngineConfigurationDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var (
		client   = common.NewCustomClient(true, d.Get("connect_address").(string), "v1", "default")
		httpUrl  = "kie/kv/{kv_id}"
		configId = d.Id()
	)

	deletePath := client.ResourceBase + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{kv_id}", configId)

	token, err := GetAuthorizationToken(d.Get("auth_address").(string), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if token != "" {
		deleteOpts.MoreHeaders = map[string]string{
			"Authorization": token,
		}
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting configuration (%s)", configId))
	}
	return nil
}

func queryMicroserviceEngineConfigurationByKey(connectAddress, token, keyName string) (interface{}, error) {
	var (
		client  = common.NewCustomClient(true, connectAddress, "v1", "default")
		httpUrl = "kie/kv"
	)
	queryPath := client.ResourceBase + httpUrl

	queryOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if token != "" {
		queryOpts.MoreHeaders = map[string]string{
			"Authorization": token,
		}
	}

	requestResp, err := client.Request("GET", queryPath, &queryOpts)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch(fmt.Sprintf("data[?key=='%s']|[0]", keyName), respBody, nil), nil
}

func resourceMicroserviceEngineConfigurationImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	var (
		token, authAddr, connectAddr, keyName, adminUser, adminPwd string
		err                                                        error
		mErr                                                       *multierror.Error

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
		case 3:
			keyName = parts[0]
			adminUser = parts[1]
			adminPwd = parts[2]
			token, err = GetAuthorizationToken(authAddr, adminUser, adminPwd)
			if err != nil {
				return nil, err
			}
			mErr = multierror.Append(mErr,
				d.Set("admin_user", adminUser),
				d.Set("admin_pass", adminPwd),
			)
		default:
			return nil, formatErr
		}
		engineDetail, err := queryMicroserviceEngineConfigurationByKey(connectAddr, token, keyName)
		if err != nil {
			return nil, err
		}
		d.SetId(utils.PathSearch("id", engineDetail, "").(string))
		return []*schema.ResourceData{d}, mErr.ErrorOrNil()
	}

	return nil, formatErr
}
