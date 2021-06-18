package huaweicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/scm/v3/certificates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const MAX_ERROR_MESSAGE_LEN = 200
const ELLIPSIS = "..."

func resourceScmCertificatePushV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceScmCertificatePushV3Create,
		Read:   resourceScmCertificatePushV3Read,
		Delete: resourceScmCertificatePushV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CDN", "WAF", "Enhance_ELB",
				}, false),
			},
			"target_project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceScmCertificatePushV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ScmV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud Certificate client: %s", err)
	}

	certificateId := d.Get("certificate_id").(string)

	pushOpts := certificates.PushOpts{
		TargetProject: d.Get("target_project").(string),
		TargetService: d.Get("target_service").(string),
	}

	err = certificates.Push(elbClient, certificateId, pushOpts).ExtractErr()
	if err != nil {
		d.SetId("")
		// Parse 'err' to print more error messages.
		errMsg := processErr(err)
		return fmt.Errorf(errMsg)
	}

	d.SetId(generateId(certificateId, pushOpts))

	return resourceScmCertificatePushV3Read(d, meta)
}

func generateId(certificateId string, opt certificates.PushOpts) (id string) {
	id = certificateId + "_" + opt.TargetService + "_" + opt.TargetProject
	return
}

// The ErrDefault500 to print only "Internal Server Error" are not clear enough.
// Parse the 'err' object to print more error messages.
func processErr(err error) string {
	// errMsg: The error message to be printed.
	errMsg := fmt.Sprintf("Push certificate service error: %s", err)
	if err500, ok := err.(golangsdk.ErrDefault500); ok {
		errBody := string(err500.Body)
		// Maybe the text in the body is very long, only 200 characters printed。
		if len(errBody) >= MAX_ERROR_MESSAGE_LEN {
			errBody = errBody[0:MAX_ERROR_MESSAGE_LEN] + ELLIPSIS
		}
		// If 'err' is an ErrDefault500 object, the following information will be printed.
		log.Printf("[ERROR] Push certificate service error. URL: %s, Body: %s",
			err500.URL, errBody)
		errMsg = fmt.Sprintf("Push certificate service error: "+
			"Bad request with: [%s %s], error message: %s", err500.Method, err500.URL, errBody)
	} else {
		// If 'err' is other error object, the default information will be printed.
		log.Printf("[ERROR] Push certificate service error: %s, \n%#v", err.Error(), err)
		errMsg = fmt.Sprintf("Push certificate service error: %s", err)
	}
	return errMsg
}

func resourceScmCertificatePushV3Read(d *schema.ResourceData, meta interface{}) error {
	// no API to read pushed info. -- 2021-6-16
	return nil
}

func resourceScmCertificatePushV3Delete(d *schema.ResourceData, meta interface{}) error {
	// no API to remove pushed services. -- 2021-6-16
	return nil
}
