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

// @API RAM POST /v1/resource-share-associations/search
func DataSourceShareAssociations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceShareAssociationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"association_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the association type.`,
			},
			"principal": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the principal associated with the resource share.`,
			},
			"resource_urn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the URN of the resource associated with the resource share.`,
			},
			"resource_share_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of resource share IDs.`,
			},
			// This filter field was not tested due to lack of test environment.
			"resource_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of resource IDs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the association.`,
			},
			"associations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of association details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"associated_entity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated entity.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the association was created.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the association.`,
						},
						"status_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the status to the association.`,
						},
						"association_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The entity type in the association.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the association was last updated.`,
						},
						"external": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the principle is in the same organization as the resource owner.`,
						},
						"resource_share_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `ID of the resource share.`,
						},
						"resource_share_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Name of the resource share.`,
						},
					},
				},
			},
		},
	}
}

// buildShareAssociationsBodyParams The default limit value for paging query is `200`, so the limit value is not
// configured here.
func buildShareAssociationsBodyParams(d *schema.ResourceData, nextMarker string) map[string]interface{} {
	requestParams := map[string]interface{}{
		"association_type":   d.Get("association_type"),
		"principal":          utils.ValueIgnoreEmpty(d.Get("principal")),
		"resource_urn":       utils.ValueIgnoreEmpty(d.Get("resource_urn")),
		"resource_share_ids": utils.ValueIgnoreEmpty(d.Get("resource_share_ids")),
		"resource_ids":       utils.ValueIgnoreEmpty(d.Get("resource_ids")),
		"association_status": utils.ValueIgnoreEmpty(d.Get("status")),
	}

	if nextMarker != "" {
		requestParams["marker"] = nextMarker
	}
	return requestParams
}

func dataSourceShareAssociationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		mErr              *multierror.Error
		nextMarker        string
		httpUrl           = "v1/resource-share-associations/search"
		product           = "ram"
		totalAssociations []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestOpt.JSONBody = utils.RemoveNil(buildShareAssociationsBodyParams(d, nextMarker))
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RAM resource share associations: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		associations := utils.PathSearch("resource_share_associations", respBody, make([]interface{}, 0)).([]interface{})
		if len(associations) > 0 {
			totalAssociations = append(totalAssociations, associations...)
		}

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("associations", flattenResourceShareAssociations(totalAssociations)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResourceShareAssociations(totalAssociations []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(totalAssociations))
	for _, v := range totalAssociations {
		rst = append(rst, map[string]interface{}{
			"associated_entity":   utils.PathSearch("associated_entity", v, nil),
			"created_at":          utils.PathSearch("created_at", v, nil),
			"status":              utils.PathSearch("status", v, nil),
			"status_message":      utils.PathSearch("status_message", v, nil),
			"association_type":    utils.PathSearch("association_type", v, nil),
			"updated_at":          utils.PathSearch("updated_at", v, nil),
			"external":            utils.PathSearch("external", v, nil),
			"resource_share_id":   utils.PathSearch("resource_share_id", v, nil),
			"resource_share_name": utils.PathSearch("resource_share_name", v, nil),
		})
	}
	return rst
}
