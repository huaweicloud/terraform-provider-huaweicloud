package apig

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceApigEnvironmentV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigEnvironmentV2Create,
		Read:   resourceApigEnvironmentV2Read,
		Update: resourceApigEnvironmentV2Update,
		Delete: resourceApigEnvironmentV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceApigEnvironmentResourceImportState,
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
					regexp.MustCompile("^[A-Za-z][\\w0-9]{2,63}$"),
					"The name consists of 3 to 64 characters, starting with a letter. "+
						"Only letters, digits and underscores (_) are allowed."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
					"The description contain a maximum of 255 characters, "+
						"and the angle brackets (< and >) are not allowed."),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildApigEnvironmentParameters(d *schema.ResourceData) environments.EnvironmentOpts {
	desc := d.Get("description").(string)
	return environments.EnvironmentOpts{
		Name:        d.Get("name").(string),
		Description: &desc,
	}
}

func resourceApigEnvironmentV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opts := buildApigEnvironmentParameters(d)
	instanceId := d.Get("instance_id").(string)
	resp, err := environments.Create(client, instanceId, opts).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Error creating HuaweiCloud APIG environment")
	}
	d.SetId(resp.Id)
	return resourceApigEnvironmentV2Read(d, meta)
}

func setApigEnvironmentParamters(d *schema.ResourceData, config *config.Config, resp *environments.Environment) error {
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("create_time", resp.CreateTime),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

// The GetEnvironmentFormServer is a method to get specifies environment form server by instance id and environment id.
func GetEnvironmentFormServer(client *golangsdk.ServiceClient, instanceId,
	envId string) (*environments.Environment, error) {
	allPages, err := environments.List(client, instanceId, environments.ListOpts{}).AllPages()
	if err != nil {
		return nil, err
	}
	envs, err := environments.ExtractEnvironments(allPages)
	if err != nil {
		return nil, err
	}
	for _, v := range envs {
		if v.Id == envId {
			return &v, nil
		}
	}
	return nil, fmtp.Errorf("The environment does not exist")
}

func resourceApigEnvironmentV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	env, err := GetEnvironmentFormServer(client, instanceId, d.Id())
	if err != nil {
		return common.CheckDeleted(d, err, fmt.Sprintf("Unable to get the environment (%s) form server", d.Id()))
	}
	return setApigEnvironmentParamters(d, config, env)
}

func resourceApigEnvironmentV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt := environments.EnvironmentOpts{
		Name: d.Get("name").(string), // Due to API restrictions, the name must be provided.
	}
	if d.HasChange("description") {
		desc := d.Get("description").(string)
		opt.Description = &desc
	}
	instanceId := d.Get("instance_id").(string)
	_, err = environments.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud APIG environment (%s): %s", d.Id(), err)
	}

	return resourceApigEnvironmentV2Read(d, meta)
}

func resourceApigEnvironmentV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	err = environments.Delete(client, instanceId, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud APIG environment from the instance (%s): %s", instanceId, err)
	}
	d.SetId("")
	return nil
}

// The ID cannot find on console, so we need to import by environment name.
func resourceApigEnvironmentResourceImportState(d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <instance_id>/<env name>")
	}
	instanceId := parts[0]
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	name := parts[1]
	opt := environments.ListOpts{
		Name: name,
	}
	pages, err := environments.List(client, instanceId, opt).AllPages()
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error retrieving environment: %s", err)
	}
	resp, err := environments.ExtractEnvironments(pages)
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmtp.Errorf("Unable to find the environment (%s) form server: %s", name, err)
	}
	d.SetId(resp[0].Id)
	d.Set("instance_id", instanceId)
	return []*schema.ResourceData{d}, nil
}
