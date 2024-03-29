package iotda

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/products
func DataSourceDeviceCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceCertificatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
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
						"created_at": {
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
				},
			},
		},
	}
}

func dataSourceDeviceCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)

	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allCertificates []model.CertificatesRspDto
		limit           = int32(50)
		offset          int32
	)

	for {
		listOpts := model.ListCertificatesRequest{
			AppId:  utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			Limit:  utils.Int32(limit),
			Offset: &offset,
		}

		listResp, listErr := client.ListCertificates(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA device CA certificates: %s", listErr)
		}

		if len(*listResp.Certificates) == 0 {
			break
		}

		allCertificates = append(allCertificates, *listResp.Certificates...)
		offset += limit
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	targetCertificates := filterListCertificates(allCertificates, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("certificates", flattenCertificates(targetCertificates)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func convertStatus(status *bool) string {
	if status != nil && *status {
		return "Verified"
	}

	return "Unverified"
}

func filterListCertificates(certificates []model.CertificatesRspDto, d *schema.ResourceData) []model.CertificatesRspDto {
	if len(certificates) == 0 {
		return nil
	}

	rst := make([]model.CertificatesRspDto, 0, len(certificates))
	for _, v := range certificates {
		if certificateID, ok := d.GetOk("certificate_id"); ok &&
			fmt.Sprint(certificateID) != utils.StringValue(v.CertificateId) {
			continue
		}

		if cn, ok := d.GetOk("cn"); ok &&
			fmt.Sprint(cn) != utils.StringValue(v.CnName) {
			continue
		}

		statusResp := convertStatus(v.Status)
		if status, ok := d.GetOk("status"); ok &&
			fmt.Sprint(status) != statusResp {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenCertificates(certificates []model.CertificatesRspDto) []interface{} {
	if len(certificates) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(certificates))
	for _, v := range certificates {
		rst = append(rst, map[string]interface{}{
			"id":             v.CertificateId,
			"cn":             v.CnName,
			"owner":          v.Owner,
			"status":         convertStatus(v.Status),
			"verify_code":    v.VerifyCode,
			"created_at":     v.CreateDate,
			"effective_date": v.EffectiveDate,
			"expiry_date":    v.ExpiryDate,
		})
	}

	return rst
}
