// Generated by PMS #558
package css

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceCssEsEnabledElbCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCssEsEnabledElbCertificatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the Elasticsearth cluster.`,
			},
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The certificates list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of SL certificate. Divided into server certificate and CA certificate.`,
						},
					},
				},
			},
		},
	}
}

type EsEnabledElbCertificatesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newEsEnabledElbCertificatesDSWrapper(d *schema.ResourceData, meta interface{}) *EsEnabledElbCertificatesDSWrapper {
	return &EsEnabledElbCertificatesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCssEsEnabledElbCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newEsEnabledElbCertificatesDSWrapper(d, meta)
	listElbCertsRst, err := wrapper.ListElbCerts()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listElbCertsToSchema(listElbCertsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/elb/certificates
func (w *EsEnabledElbCertificatesDSWrapper) ListElbCerts() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "css")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/{project_id}/clusters/{cluster_id}/elb/certificates"
	uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *EsEnabledElbCertificatesDSWrapper) listElbCertsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("certificates", schemas.SliceToList(body.Get("certificates"),
			func(certificates gjson.Result) any {
				return map[string]any{
					"id":   certificates.Get("id").Value(),
					"name": certificates.Get("name").Value(),
					"type": certificates.Get("type").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
