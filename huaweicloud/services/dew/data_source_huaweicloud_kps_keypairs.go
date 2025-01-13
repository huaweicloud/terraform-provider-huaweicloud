package dew

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

// @API DEW GET /v3/{project_id}/keypairs
func DataSourceKeypairs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeypairsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// It is not recommended to use Boolean values as filtering conditions.
			"is_managed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Deprecated: true}),
			},
			"keypairs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_managed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildRequestPathWithMarker(requestPath, marker string) string {
	if marker == "" {
		return requestPath
	}

	return fmt.Sprintf("%s?marker=%s", requestPath, marker)
}

func dataSourceKeypairsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v3/{project_id}/keypairs"
		product       = "kms"
		marker        = ""
		totalKeypairs []interface{}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		resp, err := client.Request("GET", buildRequestPathWithMarker(requestPath, marker), &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving KPS keypairs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		keypairs := utils.PathSearch("keypairs[].keypair", respBody, make([]interface{}, 0)).([]interface{})
		totalKeypairs = append(totalKeypairs, keypairs...)

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("keypairs", flattenKeypairs(filterKeypairs(totalKeypairs, d))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterKeypairs(keypairs []interface{}, d *schema.ResourceData) []interface{} {
	var (
		name        = d.Get("name").(string)
		publicKey   = d.Get("public_key").(string)
		fingerprint = d.Get("fingerprint").(string)
		result      = make([]interface{}, 0, len(keypairs))
	)

	for _, v := range keypairs {
		if name != "" && utils.PathSearch("name", v, "").(string) != name {
			continue
		}

		if publicKey != "" && utils.PathSearch("public_key", v, "").(string) != publicKey {
			continue
		}

		if fingerprint != "" && utils.PathSearch("fingerprint", v, "").(string) != fingerprint {
			continue
		}

		result = append(result, v)
	}

	return result
}

func flattenDatasourceScopeAttribute(scope string) string {
	if scope == scopeDomainValue {
		scope = scopeDomainLabel
	}

	return scope
}

func flattenKeypairs(keypairs []interface{}) []interface{} {
	if len(keypairs) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(keypairs))
	for _, v := range keypairs {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"public_key":  utils.PathSearch("public_key", v, nil),
			"fingerprint": utils.PathSearch("fingerprint", v, nil),
			"is_managed":  utils.PathSearch("is_key_protection", v, nil),
			"scope":       flattenDatasourceScopeAttribute(utils.PathSearch("scope", v, "").(string)),
		})
	}

	return result
}
