package iotda

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	v5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceDeviceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeviceCertificateCreate,
		ReadContext:   resourceDeviceCertificateRead,
		UpdateContext: resourceDeviceCertificateUpdate,
		DeleteContext: resourceDeviceCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"verify_content": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"cn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"verify_code": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"effective_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiry_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDeviceCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := model.AddCertificateRequest{
		Body: &model.CreateCertificateDto{
			Content: d.Get("content").(string),
			AppId:   utils.StringIgnoreEmpty(d.Get("space_id").(string)),
		},
	}

	resp, err := client.AddCertificate(&createOpts)
	if err != nil {
		return diag.Errorf("error creating IoTDA device CA certificate: %s", err)
	}

	if resp.CertificateId == nil {
		return diag.Errorf("error creating IoTDA device CA certificate: id is not found in API response")
	}

	d.SetId(*resp.CertificateId)
	return resourceDeviceCertificateRead(ctx, d, meta)
}

func resourceDeviceCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	detail, err := QueryDeviceCertificate(client, d.Id(), utils.StringIgnoreEmpty(d.Get("space_id").(string)))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device CA certificate")
	}

	status := "Unverified"
	if detail.Status != nil && *detail.Status {
		status = "Verified"
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cn", detail.CnName),
		d.Set("owner", detail.Owner),
		d.Set("status", status),
		d.Set("verify_code", detail.VerifyCode),
		d.Set("effective_date", detail.EffectiveDate),
		d.Set("expiry_date", detail.ExpiryDate),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDeviceCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	opts := model.CheckCertificateRequest{
		CertificateId: d.Id(),
		ActionId:      "verify",
		Body: &model.VerifyCertificateDto{
			VerifyContent: d.Get("verify_content").(string),
		},
	}

	_, err = client.CheckCertificate(&opts)
	if err != nil {
		return diag.Errorf("error verifing IoTDA device CA certificate: %s", err)
	}

	return resourceDeviceCertificateRead(ctx, d, meta)
}

func resourceDeviceCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deleteOpts := &model.DeleteCertificateRequest{
		CertificateId: d.Id(),
	}
	_, err = client.DeleteCertificate(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA device CA certificate: %s", err)
	}

	return nil
}

func QueryDeviceCertificate(client *v5.IoTDAClient, id string, spaceId *string) (*model.CertificatesRspDto, error) {
	var marker *string
	for {
		resp, err := client.ListCertificates(&model.ListCertificatesRequest{
			AppId:  spaceId,
			Limit:  utils.Int32(50),
			Marker: marker,
		})

		if err != nil {
			return nil, err
		}
		if resp.Certificates == nil || len(*resp.Certificates) == 0 {
			break
		}

		for _, v := range *resp.Certificates {
			if utils.StringValue(v.CertificateId) == id {
				return &v, nil
			}
		}
		marker = resp.Page.Marker
	}

	return nil, &sdkerr.ServiceResponseError{StatusCode: http.StatusNotFound}
}
