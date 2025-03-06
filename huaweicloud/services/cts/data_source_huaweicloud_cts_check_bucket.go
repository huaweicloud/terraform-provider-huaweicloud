package cts

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCheckBucket() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCtsCheckBucketRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the OBS bucket name.`,
			},
			"bucket_location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the OBS bucket location.`,
			},
			"kms_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the Key ID used for encrypting transferred trace files.`,
			},
			"is_support_trace_files_encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether trace files are encrypted during transfer to an OBS bucket.`,
			},
			"buckets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The OBS bucket information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS bucket name.`,
						},
						"bucket_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS bucket location.`,
						},
						"kms_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Key ID used for encrypting transferred trace files.`,
						},
						"is_support_trace_files_encryption": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether trace files are encrypted during transfer to an OBS bucket.`,
						},
						"check_bucket_response": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The check result of the OBS bucket.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The error code.`,
									},
									"error_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The error message.`,
									},
									"response_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The returned HTTP status code.`,
									},
									"success": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the transfer is successful.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type CheckBucketDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCheckBucketDSWrapper(d *schema.ResourceData, meta interface{}) *CheckBucketDSWrapper {
	return &CheckBucketDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCtsCheckBucketRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCheckBucketDSWrapper(d, meta)
	checkObsBucketsRst, err := wrapper.CheckObsBuckets()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.checkObsBucketsToSchema(checkObsBucketsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CTS POST /v3/{domain_id}/checkbucket
func (w *CheckBucketDSWrapper) CheckObsBuckets() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cts")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{domain_id}/checkbucket"
	uri = strings.ReplaceAll(uri, "{domain_id}", w.Config.DomainID)
	params := map[string]any{
		"bucket_location":                   w.Get("bucket_location"),
		"bucket_name":                       w.Get("bucket_name"),
		"is_support_trace_files_encryption": w.Get("is_support_trace_files_encryption"),
		"kms_id":                            w.Get("kms_id"),
	}
	params = utils.RemoveNil(params)
	reqBody := map[string]any{
		"buckets": []map[string]any{params},
	}
	return httphelper.New(client).
		Method("POST").
		URI(uri).
		Body(reqBody).
		Request().
		Result()
}

func (w *CheckBucketDSWrapper) checkObsBucketsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("buckets", schemas.SliceToList(body.Get("buckets"),
			func(buckets gjson.Result) any {
				return map[string]any{
					"bucket_name":                       buckets.Get("bucket_name").Value(),
					"bucket_location":                   buckets.Get("bucket_location").Value(),
					"kms_id":                            buckets.Get("kms_id").Value(),
					"is_support_trace_files_encryption": buckets.Get("is_support_trace_files_encryption").Value(),
					"check_bucket_response": schemas.SliceToList(buckets.Get("check_bucket_response"),
						func(checkBucketResponse gjson.Result) any {
							return map[string]any{
								"error_code":    checkBucketResponse.Get("error_code").Value(),
								"error_message": checkBucketResponse.Get("error_message").Value(),
								"response_code": checkBucketResponse.Get("response_code").Value(),
								"success":       checkBucketResponse.Get("success").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
