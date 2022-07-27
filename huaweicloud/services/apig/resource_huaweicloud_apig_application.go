package apig

import (
	"log"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/applications"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceApigApplicationV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigApplicationV2Create,
		Read:   resourceApigApplicationV2Read,
		Update: resourceApigApplicationV2Update,
		Delete: resourceApigApplicationV2Delete,
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
			"app_codes": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile("^[A-Za-z0-9+=][\\w!@#$%+-/=]{63,179}$"),
						"The code consists of 64 to 180 characters, starting with a letter, digit, "+
							"plus sign (+) or slash (/). Only letters, digits and following special special "+
							"characters are allowed: !@#$%+-_/="),
				},
			},
			"secret_action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RESET",
				}, false),
			},
			"registraion_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_secret": {
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

func buildApigApplicaitonParameters(d *schema.ResourceData) applications.AppOpts {
	return applications.AppOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}

func createApigV2ApplicationCode(client *golangsdk.ServiceClient, instanceId, appId, code string) error {
	opt := applications.AppCodeOpts{
		AppCode: code,
	}
	if _, err := applications.CreateAppCode(client, instanceId, appId, opt).Extract(); err != nil {
		return err
	}
	return nil
}

func createApigV2ApplicationCodes(client *golangsdk.ServiceClient, instanceId, appId string,
	codes []interface{}) error {
	for _, v := range codes {
		if err := createApigV2ApplicationCode(client, instanceId, appId, v.(string)); err != nil {
			return fmtp.Errorf("Error creating APIG v2 application code: %s", err)
		}
	}
	return nil
}

func resourceApigApplicationV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opts := buildApigApplicaitonParameters(d)
	log.Printf("[DEBUG] Create Options: %#v", opts)
	instanceId := d.Get("instance_id").(string)
	resp, err := applications.Create(client, instanceId, opts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 application at the instance (%s): %s", instanceId, err)
	}
	d.SetId(resp.Id)
	if v, ok := d.GetOk("app_codes"); ok {
		if err := createApigV2ApplicationCodes(client, instanceId, d.Id(), v.(*schema.Set).List()); err != nil {
			return fmtp.Errorf("Error creating APIG v2 application codes: %s", err)
		}
	}
	return resourceApigApplicationV2Read(d, meta)
}

func getApigApplicationCodesFromServer(d *schema.ResourceData,
	client *golangsdk.ServiceClient) ([]applications.AppCode, error) {
	allPages, err := applications.ListAppCode(client, d.Get("instance_id").(string), d.Id(),
		applications.ListCodeOpts{}).AllPages()
	if err != nil {
		return []applications.AppCode{}, err
	}
	results, err := applications.ExtractAppCodes(allPages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func setApigApplicationCodes(d *schema.ResourceData, config *config.Config, resp *applications.Application) error {
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	results, err := getApigApplicationCodesFromServer(d, client)
	if err != nil {
		return fmtp.Errorf("Error getting APIG v2 application codes: %s", err)
	}
	codes := make([]interface{}, len(results))
	for i, v := range results {
		codes[i] = v.Code
	}
	// The application code is sort by create time on server, not code.
	return d.Set("app_codes", schema.NewSet(schema.HashString, codes))
}

func setApigApplicationParamters(d *schema.ResourceData, config *config.Config, resp *applications.Application) error {
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("registraion_time", resp.RegistraionTime),
		d.Set("update_time", resp.UpdateTime),
		d.Set("app_key", resp.AppKey),
		d.Set("app_secret", resp.AppSecret),
		setApigApplicationCodes(d, config, resp),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceApigApplicationV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := applications.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "error retrieving APIG application")
	}

	return setApigApplicationParamters(d, config, resp)
}

func isCodeInApplication(codes []applications.AppCode, code string) (string, bool) {
	for _, s := range codes {
		if s.Code == code {
			return s.Id, true
		}
	}
	return "", false
}

func removeApigV2ApplicationCode(d *schema.ResourceData, client *golangsdk.ServiceClient,
	results []applications.AppCode, code string) error {
	instanceId := d.Get("instance_id").(string)
	id, ok := isCodeInApplication(results, code)
	if !ok {
		return fmtp.Errorf("Code is not in the application (%s)", d.Id())
	}
	if err := applications.RemoveAppCode(client, instanceId, d.Id(), id).ExtractErr(); err != nil {
		return fmtp.Errorf("Error removing code (%s) form the application (%s) : %s", code, d.Id(), err)
	}
	return nil
}

func updateApigApplicationCodes(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldRaws, newRaws := d.GetChange("app_codes")
	addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	removeRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	if len(removeRaws.List()) != 0 {
		results, err := getApigApplicationCodesFromServer(d, client)
		if err != nil {
			return fmtp.Errorf("Error getting APIG v2 application codes: %s", err)
		}
		for _, v := range removeRaws.List() {
			if err := removeApigV2ApplicationCode(d, client, results, v.(string)); err != nil {
				return err
			}
		}
	}
	instanceId := d.Get("instance_id").(string)
	if len(addRaws.List()) != 0 {
		for _, v := range addRaws.List() {
			if err := createApigV2ApplicationCode(client, instanceId, d.Id(), v.(string)); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceApigApplicationV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt := applications.AppOpts{}
	if d.HasChange("name") {
		opt.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		opt.Description = d.Get("description").(string)
	}
	if opt != (applications.AppOpts{}) {
		log.Printf("[DEBUG] Update Options: %#v", opt)
		instanceId := d.Get("instance_id").(string)
		_, err = applications.Update(client, instanceId, d.Id(), opt).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud APIG v2 application (%s): %s", d.Id(), err)
		}
	}
	if d.HasChange("app_codes") {
		err = updateApigApplicationCodes(d, client)
		if err != nil {
			return fmtp.Errorf("Updating HuaweiCloud APIG v2 application code failed: %s", err)
		}
	}
	if d.HasChange("secret_action") {
		if v, ok := d.GetOk("secret_action"); ok && v.(string) == "RESET" {
			if _, err := applications.ResetAppSecret(client, d.Get("instance_id").(string), d.Id(),
				applications.SecretResetOpts{}).Extract(); err != nil {
				return fmtp.Errorf("Reseting HuaweiCloud APIG v2 application secret failed: %s", err)
			}
		}
	}
	return resourceApigApplicationV2Read(d, meta)
}

func resourceApigApplicationV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	err = applications.Delete(client, instanceId, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud APIG v2 application from the instance (%s): %s",
			instanceId, err)
	}
	d.SetId("")
	return nil
}

func resourceApigInstanceSubResourceImportState(d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <instance_id>/<id>")
	}
	d.SetId(parts[1])
	d.Set("instance_id", parts[0])
	return []*schema.ResourceData{d}, nil
}
