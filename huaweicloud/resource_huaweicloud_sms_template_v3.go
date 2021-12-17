package huaweicloud

import (
	"context"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/sms/v3/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/sms/v3/region"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sms "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/sms/v3"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
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
			Create: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_template": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"projectid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_server_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volumetype": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_volume_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nics": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"publicip_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"publicip_bandwidth_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disktype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSMSTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey
	securityToken := config.SecurityToken
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		WithSecurityToken(securityToken).
		Build()
	client := sms.NewSmsClient(
		sms.SmsClientBuilder().
			WithRegion((region.ValueOf(GetRegion(d, config)))).
			WithCredential(auth).
			Build())
	request := &model.CreateTemplateRequest{}
	templatebody := getTemplatebody(ctx, d, meta)
	request.Body = &model.CreateTemplateReq{
		Template: templatebody,
	}

	var response *model.CreateTemplateResponse
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		response, err = client.CreateTemplate(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.DiagErrorf("Error create Huaweicloud SMS template client: %s", err)
	} else {
		d.SetId(*response.Id)
		d.Set("name", d.Get("name"))
		d.Set("is_template", d.Get("is_template"))
		d.Set("region", d.Get("region"))
		d.Set("projectid", d.Get("projectid"))
	}
	return resourceSMSTemplateRead(ctx, d, meta)
}

func resourceSMSTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey
	securityToken := config.SecurityToken
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		WithSecurityToken(securityToken). // 在临时aksk场景下使用
		Build()

	client := sms.NewSmsClient(
		sms.SmsClientBuilder().
			WithRegion(region.ValueOf("ap-southeast-1")).
			WithCredential(auth).
			Build())

	request := &model.ShowTemplateRequest{}
	request.Id = d.Get("id").(string)

	var response *model.ShowTemplateResponse
	err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		var err error
		response, err = client.ShowTemplate(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmtp.DiagErrorf("Error read Huaweicloud SMS template client: %s", err)
	} else {
		d.Set("name", d.Get("name"))
		d.Set("is_template", d.Get("is_template"))
		d.Set("region", d.Get("region"))
		d.Set("projectid", d.Get("projectid"))
	}
	return nil
}

func resourceSMSTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey
	securityToken := config.SecurityToken
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		WithSecurityToken(securityToken). // 在临时aksk场景下使用
		Build()
	client := sms.NewSmsClient(
		sms.SmsClientBuilder().
			WithRegion((region.ValueOf(GetRegion(d, config)))).
			WithCredential(auth).
			Build())
	request := &model.DeleteTemplateRequest{}
	request.Id = d.Get("id").(string)

	var response *model.DeleteTemplateResponse
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error
		response, err = client.DeleteTemplate(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.DiagErrorf("Error delete Huaweicloud SMS template client: %s", err)
	}
	return nil
}

func resourceSMSTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey
	securityToken := config.SecurityToken
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		WithSecurityToken(securityToken).
		Build()

	client := sms.NewSmsClient(
		sms.SmsClientBuilder().
			WithRegion((region.ValueOf(GetRegion(d, config)))).
			WithCredential(auth).
			Build())
	request := &model.UpdateTemplateRequest{}
	request.Id = d.Get("id").(string)
	templatebody := getTemplatebody(ctx, d, meta)
	request.Body = &model.UpdateTemplateReq{
		Template: templatebody,
	}

	var response *model.UpdateTemplateResponse
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		var err error
		response, err = client.UpdateTemplate(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmtp.DiagErrorf("Error update Huaweicloud SMS template client: %s", err)
	} else {
		d.Set("name", d.Get("name"))
		d.Set("is_template", d.Get("is_template"))
		d.Set("region", d.Get("region"))
		d.Set("projectid", d.Get("projectid"))
	}
	return resourceSMSTemplateRead(ctx, d, meta)
}

func getTemplatebody(ctx context.Context, d *schema.ResourceData, meta interface{}) *model.TemplateRequest {
	availabilityZoneTemplateTemplateRequest := d.Get("availability_zone").(string)
	//vpc
	var vpcObject model.VpcObject
	vpc_cidr := d.Get("vpc_cidr").(string)
	vpcObject = model.VpcObject{
		Name: d.Get("vpc_name").(string),
		Id:   d.Get("vpc_id").(string),
		Cidr: &vpc_cidr,
	}
	//网卡信息
	nics0 := d.Get("nics").([]interface{})
	var nicss []model.Nics
	for _, value := range nics0 {
		nicsValue := value.(map[string]interface{})
		tempIp := nicsValue["ip"].(string)
		nics := model.Nics{
			Id:   nicsValue["id"].(string),
			Name: nicsValue["name"].(string),
			Cidr: nicsValue["cidr"].(string),
			Ip:   &tempIp,
		}
		nicss = append(nicss, nics)
	}
	//安全组
	security_groups0 := d.Get("security_groups").([]interface{})
	var security_groups []model.SgObject
	for _, value := range security_groups0 {
		security_groupValue := value.(map[string]interface{})
		security_groups1 := model.SgObject{
			Id:   security_groupValue["id"].(string),
			Name: security_groupValue["name"].(string),
		}
		security_groups = append(security_groups, security_groups1)
	}

	//公网IP
	var publicip model.PublicIp
	publicip = model.PublicIp{
		Type:          d.Get("publicip_type").(string),
		BandwidthSize: (int32)(d.Get("publicip_bandwidth_size").(int)),
	}

	//磁盘信息
	disk0 := d.Get("disk").([]interface{})
	var disk []model.TemplateDisk
	for _, value := range disk0 {
		diskValue := value.(map[string]interface{})

		diskValue1 := model.TemplateDisk{
			Index:    int32(diskValue["index"].(int)),
			Name:     diskValue["name"].(string),
			Disktype: diskValue["disktype"].(string),
			Size:     int64(diskValue["size"].(int)),
		}
		disk = append(disk, diskValue1)
	}

	volumetype := d.Get("volumetype")
	var tempVolumeType model.TemplateRequestVolumetype
	if volumetype == "SAS" {
		tempVolumeType = model.GetTemplateRequestVolumetypeEnum().SAS
	} else if volumetype == "SSD" {
		tempVolumeType = model.GetTemplateRequestVolumetypeEnum().SSD
	} else if volumetype == "SATA" {
		tempVolumeType = model.GetTemplateRequestVolumetypeEnum().SATA
	} else {
		fmtp.Errorf("properties Volumetype type mismatch")
	}
	tempFlavor := d.Get("flavor").(string)
	tempTm := d.Get("target_server_name").(string)
	var tempDVT model.TemplateRequestDataVolumeType
	tempDVT1 := d.Get("data_volume_type").(string)
	if tempDVT1 == "SAS" {
		tempDVT = model.GetTemplateRequestDataVolumeTypeEnum().SAS
	} else if tempDVT1 == "SSD" {
		tempDVT = model.GetTemplateRequestDataVolumeTypeEnum().SSD
	} else if tempDVT1 == "SATA" {
		tempDVT = model.GetTemplateRequestDataVolumeTypeEnum().SATA
	} else {
		fmtp.Errorf("properties Volumetype type mismatch")
	}
	tempTP := d.Get("target_password").(string)
	templatebody := &model.TemplateRequest{
		Name:             d.Get("name").(string),
		IsTemplate:       d.Get("is_template").(bool),
		Region:           d.Get("region").(string),
		Projectid:        d.Get("projectid").(string),
		DataVolumeType:   &tempDVT,
		TargetPassword:   &tempTP,
		Flavor:           &tempFlavor,
		Volumetype:       &tempVolumeType,
		TargetServerName: &tempTm,
		Vpc:              &vpcObject,
		Nics:             &nicss,
		SecurityGroups:   &security_groups,
		Publicip:         &publicip,
		Disk:             &disk,
		AvailabilityZone: &availabilityZoneTemplateTemplateRequest,
	}

	return templatebody
}
