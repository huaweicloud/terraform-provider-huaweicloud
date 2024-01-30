package apig

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apigroups"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instanceId}/{path}
// @API APIG POST /v2/{project_id}/apigw/instances/{instanceId}/{path}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instanceId}/api-groups/{groupId}
// @API APIG GET /v2/{project_id}/apigw/instances/{instanceId}/api-groups/{groupId}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instanceId}/api-groups/{groupId}
// @API APIG POST /v2/{project_id}/apigw/instances/{instanceId}/api-groups
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instanceId}/{path}/{id}
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
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5\\w]*$"),
						"Only chinese and english letters, digits and underscores (_) are allowed, and must start "+
							"with a chinese or english letter. Chinese characters must be in UTF-8 or Unicode format."),
					validation.StringLenBetween(3, 64),
				),
				Description: "The group name.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 255),
				),
				Description: "The group description.",
			},
			"environment": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variable": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.All(
											validation.StringMatch(
												regexp.MustCompile(`^[A-Za-z][\w-]*$`),
												"Only letters, digits, hyphens (-) and underscores (_) are allowed, "+
													"and must start with a letter."),
											validation.StringLenBetween(3, 32),
										),
										Description: "The variable name.",
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.All(
											validation.StringMatch(regexp.MustCompile(`^[\w:/.-]*$`),
												"Only letters, digit and following special characters are allowed: _-/.:"),
											validation.StringLenBetween(1, 255),
										),
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
			"registration_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registration time.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use 'updated_at' instead",
				Description: `schema: Deprecated; The latest update time of the group.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the group.`,
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

	if environmentSet, ok := d.GetOk("environment"); ok {
		err = createEnvironmentVariables(client, instanceId, d.Id(), environmentSet.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
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
	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
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
