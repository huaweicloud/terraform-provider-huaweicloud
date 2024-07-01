package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/throttles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type PolicyType string

const (
	PolicyTypeExclusive   PolicyType = "API-based"
	PolicyTypeShared      PolicyType = "API-shared"
	PolicyTypeUser        PolicyType = "USER"
	PolicyTypeApplication PolicyType = "APP"

	includeSpecialThrottle int = 1
)

var (
	policyType = map[string]int{
		string(PolicyTypeExclusive): 1,
		string(PolicyTypeShared):    2,
	}
)

// ResourceApigThrottlingPolicyV2 is a provider resource of the APIG throttling policy.
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}/throttle-specials
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}/throttle-specials
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttles
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/throttles
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}/throttle-specials/{strategy_Id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}/throttle-specials/{strategy_Id}
func ResourceApigThrottlingPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceThrottlingPolicyCreate,
		ReadContext:   resourceThrottlingPolicyRead,
		UpdateContext: resourceThrottlingPolicyUpdate,
		DeleteContext: resourceThrottlingPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceThrottlingPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the throttling policy is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the throttling policy belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the throttling policy.",
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The period of time for limiting the number of API calls.",
			},
			"max_api_requests": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The maximum number of times an API can be accessed within a specified period..",
			},
			"max_app_requests": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of times the API can be accessed by an app within the same period.",
			},
			"max_ip_requests": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The maximum number of times the API can be accessed by an IP address within the same " +
					"period.",
			},
			"max_user_requests": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of times the API can be accessed by a user within the same period.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     string(PolicyTypeExclusive),
				Description: "The type of the request throttling policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description about the API throttling policy.",
			},
			"period_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "MINUTE",
				Description: "The time unit for limiting the number of API calls.",
			},
			"user_throttles": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    30,
				Elem:        specialThrottleSchemaResource(),
				Description: "The array of one or more special throttling policies for IAM user limit.",
			},
			"app_throttles": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    30,
				Elem:        specialThrottleSchemaResource(),
				Description: "The array of one or more special throttling policies for APP limit.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the throttling policy.",
			},
		},
	}
}

func specialThrottleSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_api_requests": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The maximum number of times an API can be accessed within a specified period.",
			},
			"throttling_object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The object ID which the special throttling policy belongs.",
			},
			"throttling_object_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The object name which the special user/application throttling policy belongs.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the special user/application throttling policy.",
			},
		},
	}
}

func buildThrottlingPolicyOpts(d *schema.ResourceData) (throttles.ThrottlingPolicyOpts, error) {
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
	policyType, ok := policyType[pType]
	if !ok {
		return opt, fmt.Errorf("invalid throttling policy type: %s", pType)
	}
	opt.Type = policyType
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

func resourceThrottlingPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	// build throttling policy create options according to terraform configuration.
	opts, err := buildThrottlingPolicyOpts(d)
	if err != nil {
		return diag.Errorf("unable to get the create option of the throttling policy: %s", err)
	}
	resp, err := throttles.Create(client, instanceId, opts).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "throttling policy")
	}
	d.SetId(resp.Id)

	// After throttling policy resoruce created, bind user throttling policies and appliation throttling policies to
	// resource according to configuration.
	if policies, ok := d.GetOk("user_throttles"); ok {
		err := addSpecThrottlingPolicies(client, policies.(*schema.Set), instanceId, d.Id(), string(PolicyTypeUser))
		if err != nil {
			return diag.Errorf("error creating special user throttling policy: %s", err)
		}
	}
	if policies, ok := d.GetOk("app_throttles"); ok {
		err := addSpecThrottlingPolicies(client, policies.(*schema.Set), instanceId, d.Id(), string(PolicyTypeApplication))
		if err != nil {
			return diag.Errorf("error creating special application throttling policy: %s", err)
		}
	}
	return resourceThrottlingPolicyRead(ctx, d, meta)
}

func analyseThrottlingPolicyType(pType int) *string {
	for k, v := range policyType {
		if v == pType {
			return &k
		}
	}
	return nil
}

func flattenSpecThrottlingPolicies(specThrottles []throttles.SpecThrottle) (userThrottles,
	appThrottles []map[string]interface{}, err error) {
	if len(specThrottles) == 0 {
		return nil, nil, nil
	}
	// According to the rules of append method, the maximum memory is expanded to 32,
	// and the average waste of memory is less than the waste caused by directly setting it to 30.
	users := make([]map[string]interface{}, 0)
	apps := make([]map[string]interface{}, 0)
	// The special throttling policies contain IAM user throttles and app throttles.
	for _, throttle := range specThrottles {
		switch throttle.ObjectType {
		case string(PolicyTypeApplication):
			apps = append(apps, map[string]interface{}{
				"max_api_requests":       throttle.CallLimits,
				"throttling_object_id":   throttle.ObjectId,
				"throttling_object_name": throttle.ObjectName,
				"id":                     throttle.ID,
			})
		case string(PolicyTypeUser):
			users = append(users, map[string]interface{}{
				"max_api_requests":       throttle.CallLimits,
				"throttling_object_id":   throttle.ObjectId,
				"throttling_object_name": throttle.ObjectName,
				"id":                     throttle.ID,
			})
		default:
			return users, apps, fmt.Errorf("invalid policy type, want '%v' or '%v', but '%v'", PolicyTypeApplication,
				PolicyTypeUser, throttle.ObjectType)
		}
	}

	return users, apps, nil
}

func resourceThrottlingPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Id()
	)
	resp, err := throttles.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "throttling policy")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("type", analyseThrottlingPolicyType(resp.Type)),
		d.Set("name", resp.Name),
		d.Set("period", resp.TimeInterval),
		d.Set("period_unit", resp.TimeUnit),
		d.Set("max_api_requests", resp.ApiCallLimits),
		d.Set("max_user_requests", resp.UserCallLimits),
		d.Set("max_app_requests", resp.AppCallLimits),
		d.Set("max_ip_requests", resp.IpCallLimits),
		d.Set("description", resp.Description),
		// Attributes
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(resp.CreateTime)/1000, false)),
	)

	if resp.IsIncludeSpecialThrottle == includeSpecialThrottle {
		// Get related special throttling policies.
		pages, err := throttles.ListSpecThrottles(client, instanceId, d.Id(), throttles.SpecThrottlesListOpts{}).AllPages()
		if err != nil {
			return diag.Errorf("error retrieving special throttle: %s", err)
		}
		specResp, err := throttles.ExtractSpecThrottles(pages)
		if err != nil {
			return diag.Errorf("unable to find the special throttles from policy: %s", err)
		}
		userThrottles, appThrottles, err := flattenSpecThrottlingPolicies(specResp)
		if err != nil {
			return diag.Errorf("error retrieving special throttle: %s", err)
		}
		mErr = multierror.Append(mErr,
			d.Set("user_throttles", userThrottles),
			d.Set("app_throttles", appThrottles),
		)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving throttling policy (%s) fields: %s", policyId, err)
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
			addObject := addPolicy["throttling_object_id"].(string)
			// If the two lists contain the objects with the same special throttling policy id, it means that the
			// policy is only updated (the delete and create operations will change policy ID).
			if rmObject == addObject {
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

func resourceThrottlingPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	if d.HasChangesExcept("user_throttles", "app_throttles") {
		opt, err := buildThrottlingPolicyOpts(d)
		if err != nil {
			return diag.Errorf("unable to get the update option of the throttling policy: %s", err)
		}
		instanceId := d.Get("instance_id").(string)
		_, err = throttles.Update(client, instanceId, d.Id(), opt).Extract()
		if err != nil {
			return diag.Errorf("error updating throttling policy: %s", err)
		}
	}
	if d.HasChange("user_throttles") {
		err = updateSpecThrottlingPolicies(d, client, "user_throttles", string(PolicyTypeUser))
		if err != nil {
			return diag.Errorf("error updating special user throttles: %s", err)
		}
	}
	if d.HasChange("app_throttles") {
		err = updateSpecThrottlingPolicies(d, client, "app_throttles", string(PolicyTypeApplication))
		if err != nil {
			return diag.Errorf("error updating special app throttles: %s", err)
		}
	}

	return resourceThrottlingPolicyRead(ctx, d, meta)
}

func resourceThrottlingPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Id()
	)
	if err = throttles.Delete(client, instanceId, policyId).ExtractErr(); err != nil {
		return diag.Errorf("unable to delete the throttling policy (%s): %s", policyId, err)
	}

	return nil
}

// The ID cannot find on the console, so we need to import by throttling policy name.
func resourceThrottlingPolicyImportState(_ context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}

	instanceId := parts[0]
	name := parts[1]
	opt := throttles.ListOpts{
		Name: name,
	}
	pages, err := throttles.List(client, instanceId, opt).AllPages()
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error retrieving throttling policies: %s", err)
	}
	resp, err := throttles.ExtractPolicies(pages)
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmt.Errorf("unable to find the throttling policy (%s) form server: %s", name, err)
	}
	d.SetId(resp[0].Id)

	return []*schema.ResourceData{d}, d.Set("instance_id", instanceId)
}
