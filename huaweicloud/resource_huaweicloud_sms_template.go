package huaweicloud

import(
	"time"
	_"fmt"
	"context"
	_"strconv"
	_"huaweicloud-sdk-go-v3/core/auth/global"
   // _sms "huaweicloud-sdk-go-v3/services/sms/v3"
	_"huaweicloud-sdk-go-v3/services/sms/v3/model"
	_"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/sms/v3/model"
	// region "huaweicloud-sdk-go-v3/services/sms/v3/region"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	sms "github.com/chnsz/golangsdk/openstack/sms"
)


func ResourceSMSTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSMSTemplateCreate,
        ReadContext:   resourceSMSTemplateRead,
		DeleteContext: resourceSMSTemplateDelete,
		UpdateContext: resourceSMSTemplateUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
			},
			"is_template": {
				Type:         schema.TypeBool,
				Required:     true,
			
			},
			"projectid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:	schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}


func resourceSMSTemplateCreate(ctx context.Context,d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	
    config := meta.(*config.Config)
    smsClinet,err := config.ServerMigrationServiceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud sms client: %s", err)
	}
	createOpts := sms.CreateOpts{
		Name:        d.Get("name").(string),
		Is_template: d.Get("is_template").(bool),
		Region:      d.Get("region").(string),
		Projectid:   d.Get("projectid").(string),
	}
	SmsCreate,err:=sms.Create(smsClinet,createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating Template: %s", err)
	}
	// SMSTemplateId := strconv.FormatInt(*SMSTemplate.Id,10)
	d.SetId(SmsCreate.ID)
	logp.Printf("[INFO] SMSTemplate ID: %s", d.Id)

	// return resourceSMSTemplateRead(ctx, d, meta)
	return nil
}


func resourceSMSTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSMSTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSMSTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}