package huaweicloud

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3.0/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityRoleCreate,
		Read:   resourceIdentityRoleRead,
		Update: resourceIdentityRoleUpdate,
		Delete: resourceIdentityRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"references": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceIdentityRoleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	policy := policies.Policy{}
	policyDoc := d.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return fmtp.Errorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Policy:      policy,
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	role, err := policies.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud Role: %s", err)
	}

	d.SetId(role.ID)

	return resourceIdentityRoleRead(d, meta)
}

func resourceIdentityRoleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	role, err := policies.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "role")
	}

	logp.Printf("[DEBUG] Retrieved HuaweiCloud Role: %#v", role)

	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("type", role.Type)
	d.Set("references", role.References)

	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return fmtp.Errorf("Error marshalling policy: %s", err)
	}
	d.Set("policy", string(policy))

	return nil
}

func resourceIdentityRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	policy := policies.Policy{}
	policyDoc := d.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return fmtp.Errorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Policy:      policy,
	}

	logp.Printf("[DEBUG] Update Options: %#v", createOpts)

	_, err = policies.Update(identityClient, d.Id(), createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud Role: %s", err)
	}

	return resourceIdentityRoleRead(d, meta)
}

func resourceIdentityRoleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	err = policies.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Role: %s", err)
	}

	return nil
}
