package apig

import (
	"regexp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apigroups"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceApigGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigGroupV2Create,
		Read:   resourceApigGroupV2Read,
		Update: resourceApigGroupV2Update,
		Delete: resourceApigGroupV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceApigInstanceSubResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z_0-9]{2,63}$"),
					"The name consists of 3 to 64 characters, starting with a letter. "+
						"Only letters, digits and underscores (_) are allowed. "+
						"Chinese characters must be in UTF-8 or Unicode format."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
					"The description contain a maximum of 255 characters, "+
						"and the angle brackets (< and >) are not allowed."),
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
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile("^[A-Za-z][\\w_-]{2,31}$"),
											"The name consists of 3 to 32 characters, starting with a letter. "+
												"Only letters, digits, hyphens (-) and underscores (_) are allowed."),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile("^[\\w:/.-]{1,255}$"),
											"The value consists of 1 to 255 characters, only letters, digit and "+
												"following special characters are allowed: _-/.:"),
									},
									"variable_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"environment_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"registraion_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createApigGroupEnvironmentVariables(client *golangsdk.ServiceClient, instanceId, groupId string,
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

func removeApigGroupEnvironmentVariables(client *golangsdk.ServiceClient, instanceId string,
	environmentSet *schema.Set) error {
	for _, env := range environmentSet.List() {
		envMap := env.(map[string]interface{})
		for _, v := range envMap["variable"].(*schema.Set).List() {
			variable := v.(map[string]interface{})
			err := environments.DeleteVariable(client, instanceId, variable["variable_id"].(string)).ExtractErr()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceApigGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	desc := d.Get("description").(string)
	opt := apigroups.GroupOpts{
		Name:        d.Get("name").(string),
		Description: &desc,
	}
	logp.Printf("[DEBUG] Create Option: %#v", opt)
	instanceId := d.Get("instance_id").(string)
	resp, err := apigroups.Create(client, instanceId, opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG group: %s", err)
	}
	d.SetId(resp.Id)
	if environments, ok := d.GetOk("environment"); ok {
		err = createApigGroupEnvironmentVariables(client, instanceId, d.Id(), environments.(*schema.Set))
		if err != nil {
			return fmtp.Errorf("Binding environment variables failed: %s", err)
		}
	}
	return resourceApigGroupV2Read(d, meta)
}

// Classify all environment variables belonging to the API group according to the APIG environment.
func setApigGroupEnvironmentVariables(d *schema.ResourceData, variables []environments.Variable) error {
	// Store all variables of the same environment in the corresponding list,
	// and generate a map with the environment ID as the key name.
	environmentMap := make(map[string]interface{})
	for _, variable := range variables {
		varMap := map[string]interface{}{
			"name":        variable.Name,
			"value":       variable.Value,
			"variable_id": variable.Id,
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

	if len(result) == 0 {
		return d.Set("environment", nil)
	}
	return d.Set("environment", result)
}

func setApigGroupParamters(d *schema.ResourceData, config *config.Config, resp *apigroups.Group) error {
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("registraion_time", resp.RegistraionTime),
		d.Set("update_time", resp.UpdateTime),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func getApigGroupEnvironmentVariables(d *schema.ResourceData,
	client *golangsdk.ServiceClient) ([]environments.Variable, error) {
	instanceId := d.Get("instance_id").(string)
	listOpt := environments.ListVariablesOpts{
		GroupId: d.Id(),
	}
	pages, err := environments.ListVariables(client, instanceId, listOpt).AllPages()
	if err != nil {
		return []environments.Variable{}, fmtp.Errorf("Error getting environment variable list from server: %s", err)
	}
	result, err := environments.ExtractVariables(pages)
	if err != nil {
		return []environments.Variable{}, fmtp.Errorf("Error extract environment variables: %s", err)
	}
	return result, nil
}

func resourceApigGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := apigroups.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "error getting APIG group")
	}
	if err = setApigGroupParamters(d, config, resp); err != nil {
		return fmtp.Errorf("Error saving group to state: %s", err)
	}
	// Saving environment variables to state file.
	variables, err := getApigGroupEnvironmentVariables(d, client)
	if err != nil {
		return err
	}
	if err = setApigGroupEnvironmentVariables(d, variables); err != nil {
		return fmtp.Errorf("Error saving variables to state: %s", err)
	}
	return nil
}

func updateApigGroupEnvironmentVariables(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldRaws, newRaws := d.GetChange("environment")
	addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	removeRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	instanceId := d.Get("instance_id").(string)
	if err := removeApigGroupEnvironmentVariables(client, instanceId, removeRaws); err != nil {
		return err
	}
	if err := createApigGroupEnvironmentVariables(client, instanceId, d.Id(), addRaws); err != nil {
		return err
	}
	return nil
}

func resourceApigGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt := apigroups.GroupOpts{}
	if d.HasChange("name") {
		opt.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		desc := d.Get("description").(string)
		opt.Description = &desc
	}
	if opt != (apigroups.GroupOpts{}) {
		logp.Printf("[DEBUG] Update Option: %#v", opt)
		instanceId := d.Get("instance_id").(string)
		_, err = apigroups.Update(client, instanceId, d.Id(), opt).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud APIG group (%s): %s", d.Id(), err)
		}
	}
	if d.HasChange("environment") {
		if err := updateApigGroupEnvironmentVariables(d, client); err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud APIG environment variables for the group (%s): %s",
				d.Id(), err)
		}
	}
	return resourceApigGroupV2Read(d, meta)
}

func resourceApigGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	err = apigroups.Delete(client, instanceId, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud APIG group from the instance (%s): %s", instanceId, err)
	}
	d.SetId("")
	return nil
}
