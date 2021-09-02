package waf

import (
	"regexp"
	"time"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/valuelists"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// ResourceWafReferenceTableV1 the resource of managing a reference table within HuaweiCloud.
func ResourceWafReferenceTableV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafReferenceTableCreate,
		Read:   resourceWafReferenceTableRead,
		Update: resourceWafReferenceTableUpdate,
		Delete: resourceWafReferenceTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^([\\w]{1,64})$"),
					"The name can contains of 1 to 64 characters."+
						"Only letters, digits and underscores (_) are allowed."),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"url", "user-agent", "ip", "params", "cookie", "referer", "header",
				}, false),
			},
			"conditions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 30,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(1, 2048),
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},

			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceWafReferenceTableCreate create a record of reference table.
func resourceWafReferenceTableCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	opt := valuelists.CreateOpts{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Values:      utils.ExpandToStringList(d.Get("conditions").([]interface{})),
		Description: d.Get("description").(string),
	}
	logp.Printf("[DEBUG] Create WAF reference table options: %#v", opt)

	r, err := valuelists.Create(client, opt)
	if err != nil {
		return fmtp.Errorf("error creating WAF reference table: %s", err)
	}
	logp.Printf("[DEBUG] Waf reference table created: %#v", r)
	d.SetId(r.Id)

	return resourceWafReferenceTableRead(d, meta)
}

// resourceWafReferenceTableRead read a record of reference table by id.
func resourceWafReferenceTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	r, err := valuelists.Get(client, d.Id())
	if err != nil {
		return common.CheckDeleted(d, err, "Error obtain WAF reference table information")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", r.Name),
		d.Set("type", r.Type),
		d.Set("conditions", r.Values),
		d.Set("description", r.Description),
		d.Set("creation_time", time.Unix(r.CreationTime/1000, 0).Format("2006-01-02 15:04:05")),
	)

	if mErr.ErrorOrNil() != nil {
		return fmtp.Errorf("error setting WAF reference table fields: %s", err)
	}

	return nil
}

// resourceWafReferenceTableUpdate update record of reference table by id.
func resourceWafReferenceTableUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	desc := d.Get("description").(string)
	opt := valuelists.UpdateValueListOpts{
		Name: d.Get("name").(string),
		// Type is required, but it cannot be changed.
		Type:        d.Get("type").(string),
		Values:      utils.ExpandToStringList(d.Get("conditions").([]interface{})),
		Description: &desc,
	}
	logp.Printf("[DEBUG] Update WAF reference table options: %#v", opt)

	_, err = valuelists.Update(client, d.Id(), opt)
	if err != nil {
		return fmtp.Errorf("error updating WAF reference table: %s", err)
	}

	return resourceWafReferenceTableRead(d, meta)
}

// resourceWafReferenceTableDelete delete the reference table record by id.
func resourceWafReferenceTableDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	_, err = valuelists.Delete(client, d.Id())
	if err != nil {
		return fmtp.Errorf("error deleting WAF reference table: %s", err)
	}

	d.SetId("")
	return nil
}
