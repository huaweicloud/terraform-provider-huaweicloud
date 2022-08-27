package dew

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kps/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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

			"is_managed": {
				Type:     schema.TypeBool,
				Optional: true,
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

func dataSourceKeypairsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcKmsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS v3 client: %s", err)
	}

	var marker string
	var allKeypairs []model.Keypair
	for {
		response, err := client.ListKeypairs(&model.ListKeypairsRequest{Marker: utils.StringIgnoreEmpty(marker)})
		if err != nil {
			return diag.Errorf("error fetching keypair: %s", err)
		}

		if response.Keypairs != nil {
			for _, k := range *response.Keypairs {
				allKeypairs = append(allKeypairs, *k.Keypair)
			}
		}

		if response.PageInfo.NextMarker != nil {
			marker = *response.PageInfo.NextMarker
		} else {
			break
		}
	}

	filter := map[string]interface{}{
		"Name":            d.Get("name"),
		"PublicKey":       d.Get("public_key"),
		"Fingerprint":     d.Get("fingerprint"),
		"IsKeyProtection": d.Get("is_managed"),
	}

	filterKeypairs, err := utils.FilterSliceWithField(allKeypairs, filter)
	if err != nil {
		return diag.Errorf("erroring filting keypair list: %s", err)
	}

	keypairsToSet, ids, err := flattenKaypairs(filterKeypairs)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hashcode.Strings(ids))

	mErr := d.Set("keypairs", keypairsToSet)
	if mErr != nil {
		return diag.Errorf("set keypairs err:%s", mErr)
	}

	return nil
}

func flattenKaypairs(keypairs []interface{}) ([]map[string]interface{}, []string, error) {
	result := make([]map[string]interface{}, len(keypairs))
	ids := make([]string, len(keypairs))

	for i, item := range keypairs {
		val := item.(model.Keypair)
		keypair := map[string]interface{}{
			"name":        val.Name,
			"public_key":  val.PublicKey,
			"fingerprint": val.Fingerprint,
			"is_managed":  val.IsKeyProtection,
		}

		scope, err := parseEncodeValue(val.Scope.MarshalJSON())
		if err != nil {
			return nil, nil, fmt.Errorf("can not parse the value of %q from response: %s", "scope", err)
		}
		if scope == scopeDomainValue {
			scope = scopeDomainLabel
		}

		keypair["scope"] = scope

		result[i] = keypair
		ids[i] = *val.Name
	}

	return result, ids, nil
}
