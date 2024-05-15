package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/signs
func DataSourceSignatures() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSignaturesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the dedicated instance to which the signatrues belong.`,
			},
			"signature_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of signature to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of signature to be queried.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of signature to be queried.`,
			},
			"algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the algorithm of the signature to be queried.`,
			},
			"signatures": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All signature key that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the signature.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the signature.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the signature.`,
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the signature.`,
						},
						"secret": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The secret of the signature.`,
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The algorithm of the signature.`,
						},
						"bind_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of bound APIs.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the signature.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the signature.`,
						},
					},
				},
			},
		},
	}
}

func buildListSignaturesParams(d *schema.ResourceData) string {
	res := ""
	if signId, ok := d.GetOk("signature_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, signId)
	}
	if name, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, name)
	}
	return res
}

func querySignatures(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/signs?limit=200"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	queryParams := buildListSignaturesParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving signatures under specified "+
				"dedicated instance (%s): %s", instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		signatures := utils.PathSearch("signs", respBody, make([]interface{}, 0)).([]interface{})
		if len(signatures) < 1 {
			break
		}
		result = append(result, signatures...)
		offset += len(signatures)
	}
	return result, nil
}

func dataSourceSignaturesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	signatures, err := querySignatures(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("signatures", filterSignatures(flattenSignatures(signatures), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterSignatures(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if signatureType, ok := d.GetOk("type"); ok &&
			fmt.Sprint(signatureType) != fmt.Sprint(utils.PathSearch("type", v, nil)) {
			continue
		}

		if algorithm, ok := d.GetOk("algorithm"); ok &&
			fmt.Sprint(algorithm) != fmt.Sprint(utils.PathSearch("algorithm", v, nil)) {
			continue
		}
		rst = append(rst, v)
	}
	return rst
}

func flattenSignatures(signatures []interface{}) []interface{} {
	if len(signatures) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(signatures))
	for _, signature := range signatures {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", signature, nil),
			"name":       utils.PathSearch("name", signature, nil),
			"type":       utils.PathSearch("sign_type", signature, nil),
			"key":        utils.PathSearch("sign_key", signature, nil),
			"secret":     utils.PathSearch("sign_secret", signature, nil),
			"algorithm":  utils.PathSearch("sign_algorithm", signature, nil),
			"bind_num":   utils.PathSearch("bind_num", signature, nil),
			"created_at": utils.PathSearch("create_time", signature, nil),
			"updated_at": utils.PathSearch("update_time", signature, nil),
		})
	}
	return result
}
