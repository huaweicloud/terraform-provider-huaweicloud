package coc

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

// @API COC GET /v1/atomics
func DataSourceCocDocumentAtomics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocDocumentAtomicsRead,

		Schema: map[string]*schema.Schema{
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of atomic capabilities.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"atomic_unique_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the unique identifier of an atomic capability.`,
						},
						"atomic_name_zh": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the atomic Chinese name.`,
						},
						"atomic_name_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the atomic English name.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the tag information.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCocDocumentAtomicsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/atomics"
	readPath := client.Endpoint + httpUrl

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	readResp, err := client.Request("GET", readPath, &readOpt)
	if err != nil {
		return diag.Errorf("error retrieving document atomics: %s", err)
	}
	readRespBody, err := utils.FlattenResponse(readResp)
	if err != nil {
		return diag.Errorf("error flattening document atomics: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("data", flattenDocumentAtomics(readRespBody)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting COC document atomics fields: %s", err)
	}

	return nil
}

func flattenDocumentAtomics(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("data", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"atomic_unique_key": utils.PathSearch("atomic_unique_key", v, nil),
			"atomic_name_zh":    utils.PathSearch("atomic_name_zh", v, nil),
			"atomic_name_en":    utils.PathSearch("atomic_name_en", v, nil),
			"tags":              utils.PathSearch("tags", v, nil),
		}
	}
	return rst
}
