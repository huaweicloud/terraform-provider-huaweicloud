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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/sign-bindings/binded-signs
func DataSourceApiAssociatedSignatures() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiAssociatedSignaturesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the signatures belong.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API bound to the signature.`,
			},
			"signature_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the signature.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the signature.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the signature.`,
			},
			"env_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the environment where the API is published.`,
			},
			"env_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the environment where the API is published.`,
			},
			"signatures": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All signatures that match the filter parameters.`,
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
							Description: `The signature key.`,
						},
						"secret": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The signature secret.`,
						},
						"env_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the environment where the API is published.`,
						},
						"env_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the environment where the API is published.`,
						},
						"bind_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The bind ID.`,
						},
						"bind_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time that the signature is bound to the API.`,
						},
					},
				},
			},
		},
	}
}

func buildListApiAssociatedSignaturesParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("env_id"); ok {
		res = fmt.Sprintf("%s&env_id=%v", res, v)
	}
	if v, ok := d.GetOk("signature_id"); ok {
		res = fmt.Sprintf("%s&sign_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&sign_name=%v", res, v)
	}
	return res
}

func queryApiAssociatedSignatures(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/sign-bindings/binded-signs?api_id={api_id}"
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{api_id}", apiId)

	queryParams := buildListApiAssociatedSignaturesParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&limit=100&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving associated signatures (bound to the API: %s) under specified "+
				"dedicated instance (%s): %s", apiId, instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		signatures := utils.PathSearch("bindings", respBody, make([]interface{}, 0)).([]interface{})
		if len(signatures) < 1 {
			break
		}
		result = append(result, signatures...)
		offset += len(signatures)
	}
	return result, nil
}

func dataSourceApiAssociatedSignaturesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	signatures, err := queryApiAssociatedSignatures(client, d)
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
		d.Set("signatures", filterAssociatedSignatures(flattenAssociatedSignatures(signatures), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterAssociatedSignatures(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("env_name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("env_name", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("type", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenAssociatedSignatures(signatures []interface{}) []interface{} {
	if len(signatures) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(signatures))
	for _, signature := range signatures {
		result = append(result, map[string]interface{}{
			"id":        utils.PathSearch("sign_id", signature, nil),
			"name":      utils.PathSearch("sign_name", signature, nil),
			"type":      utils.PathSearch("sign_type", signature, nil),
			"key":       utils.PathSearch("sign_key", signature, nil),
			"secret":    utils.PathSearch("sign_secret", signature, nil),
			"env_id":    utils.PathSearch("env_id", signature, nil),
			"env_name":  utils.PathSearch("env_name", signature, nil),
			"bind_id":   utils.PathSearch("id", signature, nil),
			"bind_time": utils.PathSearch("binding_time", signature, nil),
		})
	}
	return result
}
