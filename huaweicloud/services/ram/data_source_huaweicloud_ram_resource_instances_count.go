package ram

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM POST /v1/resource-shares/resource-instances/count
func DataSourceResourceInstancesCount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceInstancesCountRead,
		Schema: map[string]*schema.Schema{
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceResourceInstancesCountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/resource-shares/resource-instances/count"
		product = "ram"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResourceInstancesCountBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving RAM resource instances count: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("total_count", utils.PathSearch("total_count", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildResourceInstancesCountBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"without_any_tag": d.Get("without_any_tag"),
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsInput := v.([]interface{})
		tags := make([]map[string]interface{}, 0, len(tagsInput))
		for _, item := range tagsInput {
			tag, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			m := map[string]interface{}{
				"key": tag["key"],
			}
			if v, ok := tag["values"]; ok && v != nil {
				m["values"] = utils.ExpandToStringList(v.([]interface{}))
			}

			tags = append(tags, m)
		}

		params["tags"] = tags
	}

	if v, ok := d.GetOk("matches"); ok {
		matchesInput := v.([]interface{})
		matches := make([]map[string]interface{}, 0, len(matchesInput))
		for _, item := range matchesInput {
			match, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			matches = append(matches, map[string]interface{}{
				"key":   match["key"],
				"value": match["value"],
			})
		}

		params["matches"] = matches
	}

	return params
}
