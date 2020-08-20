package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/stacks"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/stacktemplates"

	"github.com/hashicorp/errwrap"
)

func resourceRTSStackV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceRTSStackV1Create,
		Read:   resourceRTSStackV1Read,
		Update: resourceRTSStackV1Update,
		Delete: resourceRTSStackV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"template_body": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStackTemplate,
				StateFunc: func(v interface{}) string {
					template, _ := normalizeStackTemplate(v)
					return template
				},
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"files": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"environment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateJsonString,
				StateFunc: func(v interface{}) string {
					json, _ := normalizeJsonString(v)
					return json
				},
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"timeout_mins": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"disable_rollback": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"outputs": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"capabilities": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"notification_topics": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceTemplateOptsV1(d *schema.ResourceData) *stacks.Template {
	var template = new(stacks.Template)
	if _, ok := d.GetOk("template_body"); ok {
		rawTemplate := d.Get("template_body").(string)
		template.Bin = []byte(rawTemplate)
	}
	if _, ok := d.GetOk("template_url"); ok {
		rawTemplateUrl := d.Get("template_url").(string)
		template.URL = rawTemplateUrl

	}
	if _, ok := d.GetOk("files"); ok {
		rawFiles := make(map[string]string)
		for key, val := range d.Get("files").(map[string]interface{}) {
			rawFiles[key] = val.(string)
		}
		template.Files = rawFiles

	}
	return template
}

func resourceEnvironmentV1(d *schema.ResourceData) *stacks.Environment {
	rawTemplate := d.Get("environment").(string)
	environment := new(stacks.Environment)
	environment.Bin = []byte(rawTemplate)
	return environment
}
func resourceParametersV1(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("parameters").(map[string]interface{}) {
		m[key] = val.(string)
	}

	return m
}
func resourceRTSStackV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	stackName := d.Get("name").(string)

	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating RTS client: %s", err)
	}

	rollback := d.Get("disable_rollback").(bool)
	createOpts := stacks.CreateOpts{
		Name:            stackName,
		TemplateOpts:    resourceTemplateOptsV1(d),
		DisableRollback: &rollback,
		EnvironmentOpts: resourceEnvironmentV1(d),
		Parameters:      resourceParametersV1(d),
		Timeout:         d.Get("timeout_mins").(int),
	}

	n, err := stacks.Create(orchestrationClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating stack: %s", err)
	}
	d.SetId(n.ID)

	log.Printf("[INFO] Stack %s created successfully", stackName)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATE_IN_PROGRESS"},
		Target:     []string{"CREATE_COMPLETE"},
		Refresh:    waitForRTSStackActive(orchestrationClient, stackName),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()

	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for Stack (%s) to become ACTIVE: %s",
			stackName, stateErr)
	}

	return resourceRTSStackV1Read(d, meta)

}

func resourceRTSStackV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating RTS Client: %s", err)
	}

	stack, err := stacks.Get(orchestrationClient, d.Id()).Extract()
	if err != nil {

		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Removing stack %s as it's already gone", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Stack: %s", err)

	}

	//Checking for stack status explicitly as Get API reports 404 or
	// gets stack with DELETE_COMPLETE status if the stack is not available depending on what is passed in Get (stackname or id)
	if stack.Status == "DELETE_COMPLETE" {
		log.Printf("[WARN] Removing stack %s as it's already gone", d.Id())
		d.SetId("")
		return nil
	}

	//setting id again to manage import as import is better done using stackname for user's ease
	d.SetId(stack.ID)

	d.Set("disable_rollback", stack.DisableRollback)

	originalParams := d.Get("parameters").(map[string]interface{})
	err = d.Set("parameters", flattenStackParameters(stack.Parameters, originalParams))
	if err != nil {
		return err
	}

	d.Set("status_reason", stack.StatusReason)
	d.Set("name", stack.Name)
	d.Set("outputs", flattenStackOutputs(stack.Outputs))
	d.Set("capabilities", stack.Capabilities)
	d.Set("notification_topics", stack.NotificationTopics)
	d.Set("timeout_mins", stack.Timeout)
	d.Set("status", stack.Status)
	d.Set("region", GetRegion(d, config))

	out, err := stacktemplates.Get(orchestrationClient, stack.Name, stack.ID).Extract()
	if err != nil {
		return err
	}

	sTemplate := BytesToString(out)
	template, error := normalizeStackTemplate(sTemplate)
	if error != nil {
		return errwrap.Wrapf("template body contains an invalid JSON or YAML: {{err}}", err)
	}
	d.Set("template_body", template)

	return nil
}

func resourceRTSStackV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating RTS Client: %s", err)
	}

	var updateOpts stacks.UpdateOpts

	updateOpts.TemplateOpts = resourceTemplateOptsV1(d)
	updateOpts.EnvironmentOpts = resourceEnvironmentV1(d)
	updateOpts.Parameters = resourceParametersV1(d)

	if d.HasChange("timeout_mins") {

		updateOpts.Timeout = d.Get("timeout_mins").(int)
	}
	if d.HasChange("disable_rollback") {

		rollback := d.Get("disable_rollback").(bool)
		updateOpts.DisableRollback = &rollback
	}

	err = stacks.Update(orchestrationClient, d.Get("name").(string), d.Id(), updateOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error updating Stack: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"UPDATE_IN_PROGRESS",
			"CREATE_COMPLETE",
			"ROLLBACK_IN_PROGRESS"},
		Target:     []string{"UPDATE_COMPLETE"},
		Refresh:    waitForRTSStackUpdate(orchestrationClient, d.Get("name").(string)),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()

	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for updating stack: %s", stateErr)
	}

	log.Printf("[INFO] Successfully updated stack %s", d.Get("name").(string))

	return resourceRTSStackV1Read(d, meta)
}

func resourceRTSStackV1Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating RTS Client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETE_IN_PROGRESS",
			"CREATE_COMPLETE",
			"CREATE_FAILED",
			"UPDATE_COMPLETE",
			"UPDATE_FAILED",
			"CREATE_FAILED",
			"ROLLBACK_COMPLETE",
			"ROLLBACK_IN_PROGRESS"},
		Target:     []string{"DELETE_COMPLETE"},
		Refresh:    waitForRTSStackDelete(orchestrationClient, d.Get("name").(string), d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()

	if stateErr != nil {
		return fmt.Errorf("Error deleting Stack: %s", stateErr)
	}

	d.SetId("")
	return nil
}

func waitForRTSStackActive(orchestrationClient *golangsdk.ServiceClient, stackName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := stacks.Get(orchestrationClient, stackName).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "CREATE_IN_PROGRESS" {
			return n, n.Status, nil
		}

		if n.Status == "CREATE_FAILED" {
			return nil, "", fmt.Errorf("%s: %s", n.Status, n.StatusReason)
		}
		return n, n.Status, nil
	}
}

func waitForRTSStackDelete(orchestrationClient *golangsdk.ServiceClient, stackName string, stackId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := stacks.Get(orchestrationClient, stackName).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted stack %s", stackName)
				return r, "DELETE_COMPLETE", nil
			}
			return r, "DELETE_IN_PROGRESS", err
		}

		//Checking for target status explicitly as Get API reports 404 or
		// gets stack with DELETE_COMPLETE status if the stack is not available depending on what is passed in Get (stackname or id)
		if r.Status == "DELETE_COMPLETE" {
			log.Printf("[INFO] Successfully deleted stack %s", stackName)
			return r, r.Status, nil
		}

		if r.Status != "DELETE_IN_PROGRESS" {
			err = stacks.Delete(orchestrationClient, stackName, stackId).ExtractErr()
			if err != nil {
				if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
					if errCode.Actual == 409 {
						return r, r.Status, nil
					}
				}
				return r, r.Status, err
			}
		}

		if r.Status == "DELETE_FAILED" {
			return r, "", fmt.Errorf("%s: %q", r.Status, r.StatusReason)
		}

		return r, r.Status, nil
	}
}

func waitForRTSStackUpdate(orchestrationClient *golangsdk.ServiceClient, stackName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := stacks.Get(orchestrationClient, stackName).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "UPDATE_IN_PROGRESS" {
			return n, "UPDATE_IN_PROGRESS", nil
		}
		if n.Status == "ROLLBACK_COMPLETE" || n.Status == "ROLLBACK_FAILED" || n.Status == "UPDATE_FAILED" {

			return nil, "", fmt.Errorf("%s: %s", n.Status, n.StatusReason)
		}

		return n, n.Status, nil
	}
}
