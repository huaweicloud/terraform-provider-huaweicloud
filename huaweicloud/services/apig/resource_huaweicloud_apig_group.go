package apig

import (
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/apigroups"
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

	return resourceApigGroupV2Read(d, meta)
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

func resourceApigGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := apigroups.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting APIG v2 group: %s", err)
	}
	if err = setApigGroupParamters(d, config, resp); err != nil {
		return fmtp.Errorf("Error saving group to state: %s", err)
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
	logp.Printf("[DEBUG] Update Option: %#v", opt)
	instanceId := d.Get("instance_id").(string)
	_, err = apigroups.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud APIG group (%s): %s", d.Id(), err)
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
