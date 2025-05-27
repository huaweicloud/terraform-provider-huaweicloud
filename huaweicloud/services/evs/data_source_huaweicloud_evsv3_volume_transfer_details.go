package evs

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

// @API EVS GET /v3/{project_id}/os-volume-transfer/detail
func DataSourceEvsV3VolumeTransferDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsV3VolumeTransferDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"transfers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of volume transfers.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The volume transfer ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The volume transfer name.`,
						},
						"volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The volume ID.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the transfer was created.`,
						},
						"links": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        datasourceV3TransferDetailLinksSchema(),
							Description: `The links to the cloud disk transfer record.`,
						},
					},
				},
			},
		},
	}
}

func datasourceV3TransferDetailLinksSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The corresponding shortcut link.`,
			},
			"rel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The shortcut link marker name.`,
			},
		},
	}
}

func buildV3TransferDetailsQueryPath(requestPath string, offset int) string {
	if offset == 0 {
		return requestPath
	}

	return fmt.Sprintf("%s?offset=%d", requestPath, offset)
}

func dataSourceEvsV3VolumeTransferDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v3/{project_id}/os-volume-transfer/detail"
		product      = "evs"
		offset       = 0
		allTransfers []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		resp, err := client.Request("GET", buildV3TransferDetailsQueryPath(requestPath, offset), &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving EVS v3 volume transfers: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		transfers := utils.PathSearch("transfers", respBody, make([]interface{}, 0)).([]interface{})
		if len(transfers) == 0 {
			break
		}

		allTransfers = append(allTransfers, transfers...)
		offset += len(transfers)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("transfers", flattenDatasourceV3VolumeTransfersDetails(allTransfers)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDatasourceV3VolumeTransfersDetails(allTransfers []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(allTransfers))
	for _, v := range allTransfers {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"volume_id":  utils.PathSearch("volume_id", v, nil),
			"created_at": utils.PathSearch("created_at", v, nil),
			"links":      flattenDatasourceV3TransferDetailsLinks(v),
		})
	}

	return rst
}

func flattenDatasourceV3TransferDetailsLinks(respBody interface{}) []interface{} {
	links := utils.PathSearch("links", respBody, make([]interface{}, 0)).([]interface{})
	if len(links) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(links))
	for _, v := range links {
		rst = append(rst, map[string]interface{}{
			"href": utils.PathSearch("href", v, nil),
			"rel":  utils.PathSearch("rel", v, nil),
		})
	}

	return rst
}
