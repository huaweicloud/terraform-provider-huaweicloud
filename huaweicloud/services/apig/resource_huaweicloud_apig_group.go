package apig

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apigroups"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/env-variables
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/env-variables
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/env-variables/{env_variable_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/api-groups
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains/{domain_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/sl-domain-access-settings
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apis/publish
func ResourceApigGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the dedicated instance is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the group belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The group name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The group description.",
			},
			"environment": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variable": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The variable name.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The variable value.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the variable that the group has.",
									},
									"variable_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "schema: Deprecated; The ID of the variable that the group has.",
										Deprecated:  "Use 'id' instead",
									},
								},
							},
							Description: "The array of one or more environment variables.",
						},
						"environment_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the environment to which the variables belongs.",
						},
					},
				},
				Description: "The array of one or more environments of the associated group.",
			},
			"url_domains": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							Description: utils.SchemaDesc(
								`The associated domain name.`,
								utils.SchemaDescInput{
									Computed: true,
								},
							),
						},
						"min_ssl_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`The minimum SSL protocol version.`,
								utils.SchemaDescInput{
									Computed: true,
								},
							),
						},
						"is_http_redirect_to_https": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`Whether to enable redirection from HTTP to HTTPS.`,
								utils.SchemaDescInput{
									Computed: true,
								},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The associated domain information of the group.`,
					utils.SchemaDescInput{
						Computed: true,
					},
				),
			},
			"domain_access_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether to use the debugging domain name to access the APIs within the group.",
			},
			"force_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to delete all sub-resources (for API) from this group.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the group, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the group, in RFC3339 format.`,
			},
			// Deprecated
			"registration_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The registration time.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The latest update time of the group.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},
		},
	}
}

func createEnvironmentVariables(client *golangsdk.ServiceClient, instanceId, groupId string,
	environmentSet *schema.Set) error {
	for _, env := range environmentSet.List() {
		envMap := env.(map[string]interface{})
		envId := envMap["environment_id"].(string)
		for _, v := range envMap["variable"].(*schema.Set).List() {
			variable := v.(map[string]interface{})
			opt := environments.CreateVariableOpts{
				Name:    variable["name"].(string),
				Value:   variable["value"].(string),
				GroupId: groupId,
				EnvId:   envId,
			}
			if _, err := environments.CreateVariable(client, instanceId, opt).Extract(); err != nil {
				return err
			}
		}
	}
	return nil
}

func removeEnvironmentVariables(client *golangsdk.ServiceClient, instanceId string,
	environmentSet *schema.Set) error {
	for _, env := range environmentSet.List() {
		envMap := env.(map[string]interface{})
		for _, v := range envMap["variable"].(*schema.Set).List() {
			variable := v.(map[string]interface{})
			err := environments.DeleteVariable(client, instanceId, variable["id"].(string)).ExtractErr()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)

		opt = apigroups.GroupOpts{
			Name:        d.Get("name").(string),
			Description: utils.String(d.Get("description").(string)),
		}
	)
	resp, err := apigroups.Create(client, instanceId, opt).Extract()
	if err != nil {
		return diag.Errorf("error creating dedicated group: %s", err)
	}
	d.SetId(resp.Id)

	groupId := d.Id()
	if environmentSet, ok := d.GetOk("environment"); ok {
		err = createEnvironmentVariables(client, instanceId, groupId, environmentSet.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if domains, ok := d.GetOk("url_domains"); ok {
		err = associateDomain(client, instanceId, groupId, domains.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// The API will not report an error if it is called repeatedly.
	err = updateDomianAccessEnabled(client, instanceId, groupId, d.Get("domain_access_enabled").(bool))
	if err != nil {
		// This feature is not available in some region, so use log.Printf to record the error.
		log.Printf("[ERROR] Update debugging domain name access status failed: %s", err)
	}

	return resourceGroupRead(ctx, d, meta)
}

func queryEnvironmentVariables(client *golangsdk.ServiceClient, instanceId, groupId string) ([]environments.Variable, error) {
	opt := environments.ListVariablesOpts{
		GroupId: groupId,
	}
	pages, err := environments.ListVariables(client, instanceId, opt).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error getting environment variable list from server: %s", err)
	}
	result, err := environments.ExtractVariables(pages)
	if err != nil {
		return nil, fmt.Errorf("error extract environment variables: %s", err)
	}
	return result, nil
}

// Classify all environment variables belonging to the API group according to the APIG environment.
func flattenEnvironmentVariables(variables []environments.Variable) []map[string]interface{} {
	if len(variables) < 1 {
		return nil
	}
	// Store all variables of the same environment in the corresponding list,
	// and generate a map with the environment ID as the key name.
	environmentMap := make(map[string]interface{})
	for _, variable := range variables {
		varMap := map[string]interface{}{
			"name":  variable.Name,
			"value": variable.Value,
			"id":    variable.Id,
		}
		if val, ok := environmentMap[variable.EnvId]; !ok {
			environmentMap[variable.EnvId] = []map[string]interface{}{
				varMap,
			}
		} else {
			environmentMap[variable.EnvId] = append(val.([]map[string]interface{}), varMap)
		}
	}
	// Generate a schema set according to the key value of the map.
	result := make([]map[string]interface{}, 0, len(environmentMap))
	for k, v := range environmentMap {
		envMap := map[string]interface{}{
			"variable":       v,
			"environment_id": k,
		}
		result = append(result, envMap)
	}

	return result
}

func flattenUrlDomain(urlDomains []apigroups.UrlDomian) []map[string]interface{} {
	if len(urlDomains) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(urlDomains))
	for i, v := range urlDomains {
		result[i] = map[string]interface{}{
			"name":                      v.DomainName,
			"min_ssl_version":           v.MinSSLVersion,
			"is_http_redirect_to_https": v.IsHttpRedirectToHttps,
		}
	}

	return result
}

func resourceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Id()
	)

	resp, err := apigroups.Get(client, instanceId, groupId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("url_domains", flattenUrlDomain(resp.UrlDomians)),
		d.Set("domain_access_enabled", resp.SlDomainAccessEnabled),
		d.Set("force_destroy", d.Get("force_destroy")),
		d.Set("created_at", resp.RegistraionTime),
		d.Set("updated_at", resp.UpdateTime),
		// Deprecated attributes
		d.Set("registration_time", resp.RegistraionTime),
		d.Set("update_time", resp.UpdateTime),
	)
	var variables []environments.Variable
	if variables, err = queryEnvironmentVariables(client, instanceId, groupId); err != nil {
		return diag.FromErr(err)
	}
	mErr = multierror.Append(mErr, d.Set("environment", flattenEnvironmentVariables(variables)))

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving dedicated group fields： %s", mErr)
	}
	return nil
}

func updateEnvironmentVariables(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		oldRaws, newRaws = d.GetChange("environment")
		addRaws          = newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
		removeRaws       = oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
		instanceId       = d.Get("instance_id").(string)
		groupId          = d.Id()
	)
	if err := removeEnvironmentVariables(client, instanceId, removeRaws); err != nil {
		return err
	}
	return createEnvironmentVariables(client, instanceId, groupId, addRaws)
}

func associateDomain(client *golangsdk.ServiceClient, instanceId, groupId string, domains []interface{}) error {
	for _, v := range domains {
		domain := v.(map[string]interface{})
		opts := apigroups.AssociateDomainOpts{
			InstanceId:            instanceId,
			GroupId:               groupId,
			UrlDomain:             domain["name"].(string),
			MinSSLVersion:         domain["min_ssl_version"].(string),
			IsHttpRedirectToHttps: domain["is_http_redirect_to_https"].(bool),
		}
		_, err := apigroups.AssociateDomain(client, opts)
		if err != nil {
			return fmt.Errorf("error binding domain name to the API group (%s): %s", groupId, err)
		}
	}
	return nil
}

func getDomainIdByName(client *golangsdk.ServiceClient, instanceId, groupId, domainName string) (string, error) {
	resp, err := apigroups.Get(client, instanceId, groupId).Extract()
	if err != nil {
		return "", fmt.Errorf("error retrieving dedicated group(%s): %s", groupId, err)
	}

	if len(resp.UrlDomians) == 0 {
		return "", fmt.Errorf("unable to find any domain name information under dedicated group: %s", groupId)
	}

	for _, v := range resp.UrlDomians {
		if v.DomainName == domainName {
			return v.Id, nil
		}
	}

	return "", golangsdk.ErrDefault404{}
}

func disAssociateDomain(client *golangsdk.ServiceClient, instanceId, groupId string, domains []interface{}) error {
	for _, v := range domains {
		domain := v.(map[string]interface{})
		domainName := domain["name"].(string)
		domainId, err := getDomainIdByName(client, instanceId, groupId, domainName)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] The domain name (%s) has been disassociated.", domainName)
				continue
			}
			return err
		}

		err = apigroups.DisAssociateDomain(client, instanceId, groupId, domainId)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateAssociateDomian(client *golangsdk.ServiceClient, d *schema.ResourceData, instanceId, groupId string) error {
	var (
		oldRaws, newRaws = d.GetChange("url_domains")
		addRaws          = newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
		removeRaws       = oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	)
	if removeRaws.Len() > 0 {
		if err := disAssociateDomain(client, instanceId, groupId, removeRaws.List()); err != nil {
			return err
		}
	}

	if addRaws.Len() > 0 {
		return associateDomain(client, instanceId, groupId, addRaws.List())
	}

	return nil
}

func updateDomianAccessEnabled(client *golangsdk.ServiceClient, instanceId, groupId string, domainAccessEnabled bool) error {
	opt := apigroups.UpdateDomainAccessEnabledOpts{
		InstanceId:            instanceId,
		GroupId:               groupId,
		SlDomainAccessEnabled: utils.Bool(domainAccessEnabled),
	}
	err := apigroups.UpdateDomainAccessEnabled(client, opt)
	if err != nil {
		return fmt.Errorf("error updating debugging domain name access (%s): %s", groupId, err)
	}
	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Id()
	)

	if d.HasChanges("name", "description") {
		opt := apigroups.GroupOpts{
			Name:        d.Get("name").(string),
			Description: utils.String(d.Get("description").(string)),
		}
		_, err = apigroups.Update(client, instanceId, groupId, opt).Extract()
		if err != nil {
			return diag.Errorf("error updating dedicated group (%s): %s", groupId, err)
		}
	}

	if d.HasChange("environment") {
		if err := updateEnvironmentVariables(client, d); err != nil {
			return diag.Errorf("error updating environment variables: %s", err)
		}
	}

	if d.HasChanges("url_domains") {
		if err := updateAssociateDomian(client, d, instanceId, groupId); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("domain_access_enabled") {
		err := updateDomianAccessEnabled(client, instanceId, groupId, d.Get("domain_access_enabled").(bool))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGroupRead(ctx, d, meta)
}

func queryApisInGroup(client *golangsdk.ServiceClient, instanceId, groupId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/apis?group_id={group_id}&limit=500"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{group_id}", groupId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffsset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffsset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error querying API list: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		apiRecords := utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{})
		if len(apiRecords) < 1 {
			break
		}
		result = append(result, apiRecords...)
		offset += len(apiRecords)
	}
	return result, nil
}

func parsePublishInfo(apiRecords []interface{}) map[string][]interface{} {
	result := make(map[string][]interface{})
	for _, apiRecord := range apiRecords {
		apiId := utils.PathSearch("id", apiRecord, "").(string)
		if apiId == "" {
			log.Printf("[ERROR] The API ID field is missing.")
			continue
		}
		publishedEnvs := utils.PathSearch("run_env_id", apiRecord, "").(string)
		for _, env := range strings.Split(strings.TrimSpace(publishedEnvs), "|") {
			if _, ok := result[env]; ok {
				result[env] = append(result[env], apiId)
				continue
			}
			result[env] = []interface{}{apiId}
		}
	}
	return result
}

func batchUnpublishApis(client *golangsdk.ServiceClient, instanceId, envId string, apiIds []interface{}) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apis/publish?action=offline"
	unpublishPath := client.Endpoint + httpUrl
	unpublishPath = strings.ReplaceAll(unpublishPath, "{project_id}", client.ProjectID)
	unpublishPath = strings.ReplaceAll(unpublishPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"apis":   apiIds,
			"env_id": envId,
		},
	}

	_, err := client.Request("POST", unpublishPath, &opt)
	return err
}

func forceUnpublishApis(client *golangsdk.ServiceClient, instanceId, groupId string) error {
	apiRecords, err := queryApisInGroup(client, instanceId, groupId)
	if err != nil {
		return err
	}
	// Unpublish first if they have published APIs.
	for envId, unpublishApiIds := range parsePublishInfo(apiRecords) {
		if err = batchUnpublishApis(client, instanceId, envId, unpublishApiIds); err != nil {
			log.Printf("[ERROR] Error unpublishing API from environment (%s): %s", envId, err)
		}
	}
	// Remove all APIs.
	for _, apiId := range utils.PathSearch("[*].id", apiRecords, make([]interface{}, 0)).([]interface{}) {
		if err = deleteApi(client, instanceId, apiId.(string)); err != nil {
			log.Printf("[ERROR] Error deleting API (%s): %s", apiId, err)
		}
	}
	return nil
}

func resourceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		groupId    = d.Id()
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	if d.Get("force_destroy").(bool) {
		if err = forceUnpublishApis(client, instanceId, groupId); err != nil {
			log.Printf("[ERROR] Unable to clean the sub-resources (for API) under group (%s): %s", groupId, err)
		}
	}

	err = apigroups.Delete(client, instanceId, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting group from the instance (%s): %s", instanceId, err)
	}

	return nil
}

func resourceGroupResourceImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
