package ram

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const pageLimit = 20

// @API RAM POST /v1/shared-resources/search
func DataSourceRAMSharedResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceRAMSharedResourcesRead,
		Schema: map[string]*schema.Schema{
			"resource_owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principal": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_urns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_share_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared_resources": {
				Type:     schema.TypeList,
				Elem:     sharedResourcesSchema(),
				Computed: true,
			},
		},
	}
}

func sharedResourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_share_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceRAMSharedResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getRAMSharedResourcesHttpUrl = "v1/shared-resources/search"
		getRAMSharedResourcesProduct = "ram"
	)
	getRAMSharedResourcesClient, err := cfg.NewServiceClient(getRAMSharedResourcesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	getRAMSharedResourcesPath := getRAMSharedResourcesClient.Endpoint + getRAMSharedResourcesHttpUrl

	getRAMSharedResourcesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	allSharedResources := make([]interface{}, 0)
	var nextMarker string
	getRAMSharedResourcesOpt.JSONBody = utils.RemoveNil(buildgetRAMSharedResourcesBodyParams(d))
	getRAMSharedResourcesJSONBody := getRAMSharedResourcesOpt.JSONBody.(map[string]interface{})
	for {
		getRAMSharedResourcesResp, err := getRAMSharedResourcesClient.Request("POST", getRAMSharedResourcesPath, &getRAMSharedResourcesOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving RAM shared resources")
		}

		getRAMSharedResourcesRespBody, err := utils.FlattenResponse(getRAMSharedResourcesResp)
		if err != nil {
			return diag.FromErr(err)
		}
		sharedResources := utils.PathSearch("shared_resources", getRAMSharedResourcesRespBody, make([]interface{}, 0)).([]interface{})
		if len(sharedResources) > 0 {
			allSharedResources = append(allSharedResources, sharedResources...)
		}

		nextMarker = utils.PathSearch("page_info.next_marker", getRAMSharedResourcesRespBody, "").(string)
		if nextMarker == "" {
			break
		}
		getRAMSharedResourcesJSONBody["marker"] = nextMarker
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("shared_resources", flattenGetSharedResourcesResponseBody(allSharedResources)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildgetRAMSharedResourcesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"principal":          utils.ValueIgnoreEmpty(d.Get("principal")),
		"resource_urns":      utils.ValueIgnoreEmpty(d.Get("resource_urns")),
		"resource_ids":       utils.ValueIgnoreEmpty(d.Get("resource_ids")),
		"resource_owner":     d.Get("resource_owner"),
		"resource_share_ids": utils.ValueIgnoreEmpty(d.Get("resource_share_ids")),
		"resource_type":      utils.ValueIgnoreEmpty(d.Get("resource_type")),
		"resource_region":    utils.ValueIgnoreEmpty(d.Get("resource_region")),
		"limit":              pageLimit,
	}
	return bodyParams
}

func flattenGetSharedResourcesResponseBody(curArray []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"resource_urn":      utils.PathSearch("resource_urn", v, nil),
			"resource_type":     utils.PathSearch("resource_type", v, nil),
			"resource_share_id": utils.PathSearch("resource_share_id", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"created_at":        utils.PathSearch("created_at", v, nil),
			"updated_at":        utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
