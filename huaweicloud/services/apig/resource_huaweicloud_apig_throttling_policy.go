package apig

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/throttles"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

const (
	timeSecond = "SECOND"
	timeMinute = "MINUTE"
	timeHour   = "HOUR"
	timeDay    = "DAY"

	typeExclusive   = "API-based"
	typeShared      = "API-shared"
	typeUser        = "USER"
	typeApplication = "APP"
)

var (
	policyType = map[string]int{
		typeExclusive: 1,
		typeShared:    2,
	}
)

// ResourceApigThrottlingPolicyV2 is a provider resource of the APIG throttling policy.
func ResourceApigThrottlingPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigThrottlingPolicyV2Create,
		Read:   resourceApigThrottlingPolicyV2Read,
		Update: resourceApigThrottlingPolicyV2Update,
		Delete: resourceApigThrottlingPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceApigThrottlingPolicyResourceImportState,
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
					regexp.MustCompile("^([\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z_0-9]{2,63})$"),
					"The name consists of 3 to 64 characters and only letters, digits, underscore (_) and chinese "+
						"characters are allowed. The name must start with a letter or chinese character."),
			},
			"period": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_api_requests": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
					"The description contain a maximum of 255 characters, "+
						"and the angle brackets (< and >) are not allowed."),
			},
			"max_app_requests": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_ip_requests": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_user_requests": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  timeMinute,
				ValidateFunc: validation.StringInSlice([]string{
					timeSecond, timeMinute, timeHour, timeDay,
				}, false),
			},
			"user_throttles": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 30,
				Elem:     specialThrottleSchemaResource(),
			},
			"app_throttles": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 30,
				Elem:     specialThrottleSchemaResource(),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  typeExclusive,
				ValidateFunc: validation.StringInSlice([]string{
					typeExclusive, typeShared,
				}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func specialThrottleSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_api_requests": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"throttling_object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"throttling_object_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildApigThrottlingPolicyParameters(d *schema.ResourceData,
	config *config.Config) (throttles.ThrottlingPolicyOpts, error) {
	opt := throttles.ThrottlingPolicyOpts{
		Name:           d.Get("name").(string),
		TimeInterval:   d.Get("period").(int),
		TimeUnit:       d.Get("period_unit").(string),
		ApiCallLimits:  d.Get("max_api_requests").(int),
		UserCallLimits: d.Get("max_user_requests").(int),
		AppCallLimits:  d.Get("max_app_requests").(int),
		IpCallLimits:   d.Get("max_ip_requests").(int),
		Description:    d.Get("description").(string),
	}
	pType := d.Get("type").(string)
	if val, ok := policyType[pType]; ok {
		opt.Type = val
	} else {
		return opt, fmtp.Errorf("Wrong throttling policy type: %s", pType)
	}
	return opt, nil
}

// addSpecThrottlingPolicies is a method which to add one or more special throttling policies to throttling policy
// resource by the list of special throttling policies.
func addSpecThrottlingPolicies(client *golangsdk.ServiceClient, policies *schema.Set,
	instanceId, policyId, specType string) error {
	for _, policy := range policies.List() {
		raw := policy.(map[string]interface{})
		specOpts := throttles.SpecThrottleCreateOpts{
			ObjectType: specType,
			ObjectId:   raw["throttling_object_id"].(string),
			CallLimits: raw["max_api_requests"].(int),
		}
		_, err := throttles.CreateSpecThrottle(client, instanceId, policyId, specOpts).Extract()
		if err != nil {
			return err
		}
	}
	return nil
}

// removeSpecThrottlingPolicies is a method which to remove the special throttling policy form throttling policy
// resource by specifies special throttling policy ID.
func removeSpecThrottlingPolicies(client *golangsdk.ServiceClient, policies *schema.Set,
	instanceId, policyId string) error {
	for _, policy := range policies.List() {
		raw := policy.(map[string]interface{})
		err := throttles.DeleteSpecThrottle(client, instanceId, policyId, raw["id"].(string)).ExtractErr()
		if err != nil {
			return err
		}
	}
	return nil
}

// updateSpecThrottlingPolicieCallLimit is a method which to udpate the API call limit of the special throttling policy
// by specifies special throttling policy ID.
func updateSpecThrottlingPolicieCallLimit(client *golangsdk.ServiceClient, instanceId, policyId, strategyId string,
	limit int) error {
	opts := &throttles.SpecThrottleUpdateOpts{
		CallLimits: limit,
	}
	_, err := throttles.UpdateSpecThrottle(client, instanceId, policyId, strategyId, opts).Extract()
	if err != nil {
		return err
	}
	return nil
}

func resourceApigThrottlingPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	// build throttling policy create options according to terraform configuration.
	opts, err := buildApigThrottlingPolicyParameters(d, config)
	if err != nil {
		return fmtp.Errorf("Unable to get the create option of the throttling policy: %s", err)
	}
	v, err := throttles.Create(client, instanceId, opts).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Error creating HuaweiCloud throttling policy")
	}
	d.SetId(v.Id)

	// After throttling policy resoruce created, bind user throttling policies and appliation throttling policies to
	// resource according to configuration.
	if policies, ok := d.GetOk("user_throttles"); ok {
		err := addSpecThrottlingPolicies(client, policies.(*schema.Set), instanceId, d.Id(), typeUser)
		if err != nil {
			return fmtp.Errorf("Error creating special user throttling policy: %s", err)
		}
	}
	if policies, ok := d.GetOk("app_throttles"); ok {
		err := addSpecThrottlingPolicies(client, policies.(*schema.Set), instanceId, d.Id(), typeApplication)
		if err != nil {
			return fmtp.Errorf("Error creating special application throttling policy: %s", err)
		}
	}
	return resourceApigThrottlingPolicyV2Read(d, meta)
}

func setApigThrottlingPolicyType(d *schema.ResourceData, pType int) error {
	for k, v := range policyType {
		if v == pType {
			return d.Set("type", k)
		}
	}
	return fmtp.Errorf("The member type (%d) is not supported", pType)
}

func setApigThrottlingPolicyParameters(d *schema.ResourceData, config *config.Config, resp throttles.ThrottlingPolicy) error {
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("period", resp.TimeInterval),
		d.Set("period_unit", resp.TimeUnit),
		d.Set("max_api_requests", resp.ApiCallLimits),
		d.Set("max_user_requests", resp.UserCallLimits),
		d.Set("max_app_requests", resp.AppCallLimits),
		d.Set("max_ip_requests", resp.IpCallLimits),
		d.Set("description", resp.Description),
		d.Set("create_time", resp.CreateTime),
		setApigThrottlingPolicyType(d, resp.Type),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func setSpecThrottlingPolicies(d *schema.ResourceData, specThrottles []throttles.SpecThrottle) error {
	if len(specThrottles) == 0 {
		return nil
	}
	// According to the rules of append method, the maximum memory is expanded to 32,
	// and the average waste of memory is less than the waste caused by directly setting it to 30.
	users := make([]map[string]interface{}, 0)
	apps := make([]map[string]interface{}, 0)
	// The special throttling policies contain IAM user throttles and app throttles.
	for _, throttle := range specThrottles {
		if throttle.ObjectType == typeApplication {
			apps = append(apps, map[string]interface{}{
				"max_api_requests":       throttle.CallLimits,
				"throttling_object_id":   throttle.ObjectId,
				"throttling_object_name": throttle.ObjectName,
				"id":                     throttle.ID,
			})
		} else {
			users = append(users, map[string]interface{}{
				"max_api_requests":       throttle.CallLimits,
				"throttling_object_id":   throttle.ObjectId,
				"throttling_object_name": throttle.ObjectName,
				"id":                     throttle.ID,
			})
		}
	}
	mErr := multierror.Append(nil,
		d.Set("user_throttles", users),
		d.Set("app_throttles", apps),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceApigThrottlingPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	resp, err := throttles.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, fmt.Sprintf("error getting throttle (%s) form server", d.Id()))
	}
	err = setApigThrottlingPolicyParameters(d, config, *resp)
	if err != nil {
		return fmtp.Errorf("Error setting throttles to state: %s", d.Id(), err)
	}

	// Set special throttling policies for IAM user and application to state.
	pages, err := throttles.ListSpecThrottles(client, instanceId, d.Id(), throttles.SpecThrottlesListOpts{}).AllPages()
	if err != nil {
		return fmtp.Errorf("Error retrieving special throttle: %s", err)
	}
	specResp, err := throttles.ExtractSpecThrottles(pages)
	if err != nil {
		return fmtp.Errorf("Unable to find the special throttles from policy: %s", err)
	}
	return setSpecThrottlingPolicies(d, specResp)
}

func isBasicParamsChanged(d *schema.ResourceData) bool {
	//lintignore:R019
	if d.HasChanges("name", "period", "max_api_requests", "description", "period_unit") ||
		d.HasChanges("max_app_requests", "max_ip_requests", "max_user_requests", "type") {
		return true
	}
	return false
}

func updateSpecThrottlingPolicies(d *schema.ResourceData, client *golangsdk.ServiceClient,
	paramName, specType string) error {
	oldRaws, newRaws := d.GetChange(paramName)
	addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	removeRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	instanceId := d.Get("instance_id").(string)

	// If only max API requests update, the ID should be the same after the special throttling policy is updated.
	for _, rm := range removeRaws.List() {
		rmPolicy := rm.(map[string]interface{})
		rmObject := rmPolicy["throttling_object_id"].(string)
		for _, add := range addRaws.List() {
			addPolicy := add.(map[string]interface{})
			// If the two lists contain the objects with the same special throttling policy id, it means that the
			// policy is only updated (the delete and create operations will change policy ID).
			if rmObject == addPolicy["throttling_object_id"].(string) {
				strategyId := rmPolicy["id"].(string)
				limit := addPolicy["max_api_requests"].(int)
				// Update specifies special throttling policy by strategy ID.
				err := updateSpecThrottlingPolicieCallLimit(client, instanceId, d.Id(), strategyId, limit)
				if err != nil {
					return err
				}
				removeRaws.Remove(rm)
				addRaws.Remove(add)
			}
		}
	}
	err := removeSpecThrottlingPolicies(client, removeRaws, instanceId, d.Id())
	if err != nil {
		return err
	}
	err = addSpecThrottlingPolicies(client, addRaws, instanceId, d.Id(), specType)
	if err != nil {
		return err
	}
	return nil
}

func resourceApigThrottlingPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	if isBasicParamsChanged(d) {
		opt, err := buildApigThrottlingPolicyParameters(d, config)
		if err != nil {
			return fmtp.Errorf("Unable to get the update option of the throttling policy: %s", err)
		}
		instanceId := d.Get("instance_id").(string)
		_, err = throttles.Update(client, instanceId, d.Id(), opt).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating throttling policy: %s", err)
		}
	}
	if d.HasChange("user_throttles") {
		err = updateSpecThrottlingPolicies(d, client, "user_throttles", typeUser)
		if err != nil {
			return fmtp.Errorf("Error updating special user throttles: %s", err)
		}
	}
	if d.HasChange("app_throttles") {
		err = updateSpecThrottlingPolicies(d, client, "app_throttles", typeApplication)
		if err != nil {
			return fmtp.Errorf("Error updating special app throttles: %s", err)
		}
	}

	return resourceApigThrottlingPolicyV2Read(d, meta)
}

func resourceApigThrottlingPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	if err = throttles.Delete(client, instanceId, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Unable to delete the throttling policy (%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

// The ID cannot find on the console, so we need to import by throttling policy name.
func resourceApigThrottlingPolicyResourceImportState(d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <instance_id>/<name>")
	}

	instanceId := parts[0]
	name := parts[1]
	opt := throttles.ListOpts{
		Name: name,
	}
	pages, err := throttles.List(client, instanceId, opt).AllPages()
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error retrieving throttling policies: %s", err)
	}
	resp, err := throttles.ExtractPolicies(pages)
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmtp.Errorf("Unable to find the throttling policy (%s) form server: %s", name, err)
	}
	d.SetId(resp[0].Id)
	d.Set("instance_id", instanceId)

	return []*schema.ResourceData{d}, nil
}
