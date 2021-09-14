package huaweicloud

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCCEAddonV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCEAddonV3Create,
		Read:   resourceCCEAddonV3Read,
		Delete: resourceCCEAddonV3Delete,

		Importer: &schema.ResourceImporter{
			State: resourceCCEAddonV3Import,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic": {
							Type:         schema.TypeMap,
							Optional:     true,
							ForceNew:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							ExactlyOneOf: []string{"values.0.basic", "values.0.basic_json"},
						},
						"basic_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
							ExactlyOneOf: []string{"values.0.basic", "values.0.basic_json"},
						},
						"custom": {
							Type:          schema.TypeMap,
							Optional:      true,
							ForceNew:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							ConflictsWith: []string{"values.0.custom_json"},
						},
						"custom_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
							ConflictsWith: []string{"values.0.custom"},
						},
						"flavor": {
							Type:          schema.TypeMap,
							Optional:      true,
							ForceNew:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							ConflictsWith: []string{"values.0.flavor_json"},
						},
						"flavor_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
							ConflictsWith: []string{"values.0.flavor"},
						},
					},
				},
			},
		},
	}
}

func getValuesValues(d *schema.ResourceData) (basic, custom, flavor map[string]interface{}, err error) {
	values := d.Get("values").([]interface{})
	if len(values) == 0 {
		basic = map[string]interface{}{}
		return
	}
	valuesMap := values[0].(map[string]interface{})

	if basicRaw := valuesMap["basic"].(map[string]interface{}); len(basicRaw) != 0 {
		basic = basicRaw
	}
	if basicJsonRaw := valuesMap["basic_json"].(string); basicJsonRaw != "" {
		err = json.Unmarshal([]byte(basicJsonRaw), &basic)
		if err != nil {
			err = fmtp.Errorf("Error unmarshalling basic json: %s", err)
			return
		}
	}

	if customRaw := valuesMap["custom"].(map[string]interface{}); len(customRaw) != 0 {
		custom = customRaw
	}
	if customJsonRaw := valuesMap["custom_json"].(string); customJsonRaw != "" {
		err = json.Unmarshal([]byte(customJsonRaw), &custom)
		if err != nil {
			err = fmtp.Errorf("Error unmarshalling custom json: %s", err)
			return
		}
	}

	if flavorRaw := valuesMap["flavor"].(map[string]interface{}); len(flavorRaw) != 0 {
		flavor = flavorRaw
	}
	if flavorJsonRaw := valuesMap["flavor_json"].(string); flavorJsonRaw != "" {
		err = json.Unmarshal([]byte(flavorJsonRaw), &flavor)
		if err != nil {
			err = fmtp.Errorf("Error unmarshalling flavor json %s", err)
			return
		}
	}

	return
}

func resourceCCEAddonV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	cceClient, err := config.CceAddonV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Unable to create HuaweiCloud CCE client : %s", err)
	}

	var cluster_id = d.Get("cluster_id").(string)

	basic, custom, flavor, err := getValuesValues(d)
	if err != nil {
		return fmtp.Errorf("error getting values for CCE addon: %s", err)
	}

	createOpts := addons.CreateOpts{
		Kind:       "Addon",
		ApiVersion: "v3",
		Metadata: addons.CreateMetadata{
			Anno: addons.Annotations{
				AddonInstallType: "install",
			},
		},
		Spec: addons.RequestSpec{
			Version:           d.Get("version").(string),
			ClusterID:         cluster_id,
			AddonTemplateName: d.Get("template_name").(string),
			Values: addons.Values{
				Basic:  basic,
				Custom: custom,
				Flavor: flavor,
			},
		},
	}

	create, err := addons.Create(cceClient, createOpts, cluster_id).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCEAddon: %s", err)
	}

	d.SetId(create.Metadata.Id)

	logp.Printf("[DEBUG] Waiting for HuaweiCloud CCEAddon (%s) to become available", create.Metadata.Id)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"installing", "abnormal"},
		Target:       []string{"running", "available"},
		Refresh:      waitForCCEAddonActive(cceClient, create.Metadata.Id, cluster_id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCEAddon: %s", err)
	}

	return resourceCCEAddonV3Read(d, meta)
}

func resourceCCEAddonV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	cceClient, err := config.CceAddonV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	var cluster_id = d.Get("cluster_id").(string)

	n, err := addons.Get(cceClient, d.Id(), cluster_id).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving HuaweiCloud CCEAddon: %s", err)
	}

	d.Set("cluster_id", n.Spec.ClusterID)
	d.Set("version", n.Spec.Version)
	d.Set("template_name", n.Spec.AddonTemplateName)
	d.Set("status", n.Status.Status)
	d.Set("description", n.Spec.Description)

	return nil
}

func resourceCCEAddonV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	cceClient, err := config.CceAddonV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCEAddon Client: %s", err)
	}

	var cluster_id = d.Get("cluster_id").(string)

	err = addons.Delete(cceClient, d.Id(), cluster_id).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud CCE Addon: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting", "Available", "Unavailable"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCCEAddonDelete(cceClient, d.Id(), cluster_id),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud CCE Addon: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCCEAddonActive(cceAddonClient *golangsdk.ServiceClient, id, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := addons.Get(cceAddonClient, id, clusterID).Extract()
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Status, nil
	}
}

func waitForCCEAddonDelete(cceClient *golangsdk.ServiceClient, id, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud CCE Addon %s.\n", id)

		r, err := addons.Get(cceClient, id, clusterID).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud CCE Addon %s", id)
				return r, "Deleted", nil
			}
		}
		if r.Status.Status == "Deleting" {
			return r, "Deleting", nil
		}
		logp.Printf("[DEBUG] HuaweiCloud CCE Addon %s still available.\n", id)
		return r, "Available", nil
	}
}

func resourceCCEAddonV3Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmtp.Errorf("Invalid format specified for CCE Addon. Format must be <cluster id>/<addon id>")
		return nil, err
	}

	clusterID := parts[0]
	addonID := parts[1]

	d.SetId(addonID)
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}
