package apig

import (
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/responses"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceApigResponseV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigResponseV2Create,
		Read:   resourceApigResponseV2Read,
		Update: resourceApigResponseV2Update,
		Delete: resourceApigResponseV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceApigResponseResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
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
					regexp.MustCompile("^[A-Za-z0-9_-]{1,64}$"), "The name consists of 1 to 64 characters. "+
						"Only letters, digits, hyphens(-), and underscores (_) are allowed."),
			},
			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AUTH_FAILURE", "AUTH_HEADER_MISSING", "AUTHORIZER_FAILURE", "AUTHORIZER_CONF_FAILURE",
								"AUTHORIZER_IDENTITIES_FAILURE", "BACKEND_UNAVAILABLE", "BACKEND_TIMEOUT", "THROTTLED",
								"UNAUTHORIZED", "ACCESS_DENIED", "NOT_FOUND", "REQUEST_PARAMETERS_FAILURE",
								"DEFAULT_4XX", "DEFAULT_5XX",
							}, false),
						},
						"body": {
							Type: schema.TypeString,
							// If parameter body omitted, The API will return 'The parameters must be specified'.
							Required: true,
						},
						"status_code": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(200, 599),
						},
					},
				},
			},
			"create_time": {
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

// 'error_type' is the key of the response mapping, and 'body' and 'status_code' are the structural elements of the
// mapping value.
func buildApigGroupCustomResponses(d *schema.ResourceData) map[string]responses.ResponseInfo {
	result := make(map[string]responses.ResponseInfo)
	respSet := d.Get("rule").(*schema.Set)
	for _, response := range respSet.List() {
		rule := response.(map[string]interface{})
		errorType := rule["error_type"].(string)
		result[errorType] = responses.ResponseInfo{
			Body:   rule["body"].(string),
			Status: rule["status_code"].(int),
		}
	}
	return result
}

func resourceApigResponseV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt := responses.ResponseOpts{
		Name:       d.Get("name").(string),
		Responses:  buildApigGroupCustomResponses(d),
		InstanceId: d.Get("instance_id").(string),
		GroupId:    d.Get("group_id").(string),
	}
	resp, err := responses.Create(client, opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG group: %s", err)
	}
	d.SetId(resp.Id)

	return resourceApigResponseV2Read(d, meta)
}

func setApigGroupCustomResponses(d *schema.ResourceData, respMap map[string]responses.ResponseInfo) error {
	result := make([]map[string]interface{}, 0)
	for errorType, rule := range respMap {
		if rule.IsDefault {
			// The IsDefault of the modified response will be marked as false,
			// record these responses and skip other unmodified responses (IsDefault is true).
			continue
		}
		result = append(result, map[string]interface{}{
			"error_type":  errorType,
			"body":        rule.Body,
			"status_code": rule.Status,
		})
	}
	return d.Set("rule", result)
}

func setApigResponseParamters(d *schema.ResourceData, config *config.Config, resp *responses.Response) error {
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		setApigGroupCustomResponses(d, resp.Responses),
		d.Set("create_time", resp.CreateTime),
		d.Set("update_time", resp.UpdateTime),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceApigResponseV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	groupId := d.Get("group_id").(string)
	resp, err := responses.Get(client, instanceId, groupId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "error getting custom response from server")
	}
	if err = setApigResponseParamters(d, config, resp); err != nil {
		return fmtp.Errorf("Error saving custom response to state: %s", err)
	}
	return nil
}

func resourceApigResponseV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	// Only updating the name will cause all the response rules that have been set to be reset, so no matter whether
	// the response rules are updated or not, the response rules must be carried in the update opt.
	opt := responses.ResponseOpts{
		Name:       d.Get("name").(string),
		Responses:  buildApigGroupCustomResponses(d),
		InstanceId: d.Get("instance_id").(string),
		GroupId:    d.Get("group_id").(string),
	}
	_, err = responses.Update(client, d.Id(), opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud APIG custom response: %s", err)
	}
	return resourceApigResponseV2Read(d, meta)
}

func resourceApigResponseV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	groupId := d.Get("group_id").(string)
	err = responses.Delete(client, instanceId, groupId, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud custom response (%s) from the APIG group (%s): %s",
			d.Id(), groupId, err)
	}
	d.SetId("")
	return nil
}

// Some resources of the APIG service are associated with dedicated instances and groups,
// but their IDs cannot be found on the console.
// This method is used to solve the above problem by importing resources by associating ID and their name.
func resourceApigResponseResourceImportState(d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmtp.Errorf("Invalid format specified for import ids and name, " +
			"must be <instance_id>/<group_id>/<name>")
	}

	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := parts[0]
	groupId := parts[1]
	opt := responses.ListOpts{
		InstanceId: instanceId,
		GroupId:    groupId,
	}
	pages, err := responses.List(client, opt).AllPages()
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error getting custom response list from server: %s", err)
	}
	resp, err := responses.ExtractResponses(pages)
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error extract custom responses: %s", err)
	}
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmtp.Errorf("Unable to find any custom response from server")
	}
	// Since there are no parameters about custom responses in the query options, we need to get the response list and
	// filter by the response name.
	name := parts[2]
	for _, val := range resp {
		if val.Name == name {
			d.SetId(val.Id)
			d.Set("instance_id", instanceId)
			d.Set("group_id", groupId)
			return []*schema.ResourceData{d}, nil
		}
	}
	return []*schema.ResourceData{d}, fmtp.Errorf("Unable to find the custom response (%s) from server")
}
