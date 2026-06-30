package geminidb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB POST /v3/{project_id}/instances/{instance_id}/big-keys
func DataSourceBigKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigKeysRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"key_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBigKeysParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"key_types": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("key_types").([]interface{}))),
		"limit":     100,
	}
}

func dataSourceBigKeysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/big-keys"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)
	bodyParams := utils.RemoveNil(buildBigKeysParams(d))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		bodyParams["offset"] = offset
		requestOpt.JSONBody = bodyParams

		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving GeminiDB Redis instance big keys: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		bigKeys := utils.PathSearch("keys", respBody, make([]interface{}, 0)).([]interface{})
		if len(bigKeys) == 0 {
			break
		}

		result = append(result, bigKeys...)
		offset += len(bigKeys)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("keys", flattenBigKeys(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBigKeys(mappingInfos []interface{}) []interface{} {
	if len(mappingInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(mappingInfos))
	for _, v := range mappingInfos {
		rst = append(rst, map[string]interface{}{
			"db_id":    utils.PathSearch("db_id", v, nil),
			"key_type": utils.PathSearch("key_type", v, nil),
			"key_name": utils.PathSearch("key_name", v, nil),
			"key_size": utils.PathSearch("key_size", v, nil),
		})
	}

	return rst
}
