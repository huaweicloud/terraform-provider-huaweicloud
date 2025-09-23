package iotda

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/cetificates
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

func buildDeviceCertificatesQueryParams(d *schema.ResourceData) string {
	if spaceId, ok := d.GetOk("space_id"); ok {
		return fmt.Sprintf("&app_id=%v", spaceId)
	}

	return ""
}

func dataSourceDeviceCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		isDerived       = WithDerivedAuth(cfg, region)
		httpUrl         = "v5/iot/{project_id}/certificates?limit=50"
		allCertificates = make([]interface{}, 0)
		offset          = 0
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildDeviceCertificatesQueryParams(d)
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpts)
		if err != nil {
			return diag.Errorf("error querying IoTDA device CA certificates: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		certificates := utils.PathSearch("certificates", respBody, make([]interface{}, 0)).([]interface{})
		if len(certificates) == 0 {
			break
		}

		allCertificates = append(allCertificates, certificates...)
		offset += len(certificates)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("certificates", flattenCertificates(filterListCertificates(allCertificates, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListCertificates(certificates []interface{}, d *schema.ResourceData) []interface{} {
	if len(certificates) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(certificates))
	for _, v := range certificates {
		if certificateID, ok := d.GetOk("certificate_id"); ok &&
			fmt.Sprint(certificateID) != utils.PathSearch("certificate_id", v, "").(string) {
			continue
		}

		if cn, ok := d.GetOk("cn"); ok &&
			fmt.Sprint(cn) != utils.PathSearch("cn_name", v, "").(string) {
			continue
		}

		statusResp := flattenStatus(utils.PathSearch("status", v, false).(bool))
		if status, ok := d.GetOk("status"); ok &&
			fmt.Sprint(status) != statusResp {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenCertificates(certificates []interface{}) []interface{} {
	if len(certificates) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(certificates))
	for _, v := range certificates {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("certificate_id", v, nil),
			"cn":             utils.PathSearch("cn_name", v, nil),
			"owner":          utils.PathSearch("owner", v, nil),
			"status":         flattenStatus(utils.PathSearch("status", v, false).(bool)),
			"verify_code":    utils.PathSearch("verify_code", v, nil),
			"created_at":     utils.PathSearch("create_date", v, nil),
			"effective_date": utils.PathSearch("effective_date", v, nil),
			"expiry_date":    utils.PathSearch("expiry_date", v, nil),
		})
	}

	return rst
}
