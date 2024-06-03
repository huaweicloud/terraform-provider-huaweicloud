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

// @API RAM POST /v1/shared-principals/search
func DataSourceRAMSharedPrincipals() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceRAMSharedPrincipalsRead,
		Schema: map[string]*schema.Schema{
			"resource_owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principal": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_urn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_share_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared_principals": {
				Type:     schema.TypeList,
				Elem:     sharedPrincipalsSchema(),
				Computed: true,
			},
		},
	}
}

func sharedPrincipalsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_share_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
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

func datasourceRAMSharedPrincipalsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getRAMSharedPrincipalsHttpUrl = "v1/shared-principals/search"
		getRAMSharedPrincipalsProduct = "ram"
	)
	getRAMSharedPrincipalsClient, err := cfg.NewServiceClient(getRAMSharedPrincipalsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	getRAMSharedPrincipalsPath := getRAMSharedPrincipalsClient.Endpoint + getRAMSharedPrincipalsHttpUrl

	getRAMSharedPrincipalsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	allSharedPrincipals := make([]interface{}, 0)
	var nextMarker string
	getRAMSharedPrincipalsOpt.JSONBody = utils.RemoveNil(buildgetRAMSharedPrincipalsBodyParams(d))
	getRAMSharedPrincipalsJSONBody := getRAMSharedPrincipalsOpt.JSONBody.(map[string]interface{})
	for {
		getRAMSharedPrincipalsResp, err := getRAMSharedPrincipalsClient.Request("POST", getRAMSharedPrincipalsPath, &getRAMSharedPrincipalsOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving RAM shared principals")
		}

		getRAMSharedPrincipalsRespBody, err := utils.FlattenResponse(getRAMSharedPrincipalsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		sharedPrincipals := utils.PathSearch("shared_principals", getRAMSharedPrincipalsRespBody, make([]interface{}, 0)).([]interface{})
		if len(sharedPrincipals) > 0 {
			allSharedPrincipals = append(allSharedPrincipals, sharedPrincipals...)
		}

		nextMarker = utils.PathSearch("page_info.next_marker", getRAMSharedPrincipalsRespBody, "").(string)
		if nextMarker == "" {
			break
		}
		getRAMSharedPrincipalsJSONBody["marker"] = nextMarker
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("shared_principals", flattenGetSharedPrincipalsResponseBody(allSharedPrincipals)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildgetRAMSharedPrincipalsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resource_urn":   utils.ValueIgnoreEmpty(d.Get("resource_urn")),
		"resource_owner": d.Get("resource_owner"),
		"limit":          pageLimit,
	}

	if v, ok := d.GetOk("principal"); ok {
		bodyParams["principals"] = []interface{}{v}
	}

	if v, ok := d.GetOk("resource_share_id"); ok {
		bodyParams["resource_share_ids"] = []interface{}{v}
	}
	return bodyParams
}

func flattenGetSharedPrincipalsResponseBody(curArray []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"resource_share_id": utils.PathSearch("resource_share_id", v, nil),
			"id":                utils.PathSearch("id", v, nil),
			"created_at":        utils.PathSearch("created_at", v, nil),
			"updated_at":        utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
