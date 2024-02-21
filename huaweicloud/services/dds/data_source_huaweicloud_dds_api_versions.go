package dds

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

// @API DDS GET
func DataSourceDdsAPIVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdsAPIVersionsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Elem:     ddsVersionsSchema(),
				Computed: true,
			},
		},
	}
}

func ddsVersionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     ddsLinkSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func ddsLinkSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"href": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rel": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDdsAPIVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// listDdsAPIVersions: Query the List of DDS API versions.
	var (
		listDdsAPIVersionsHttpUrl = ""
		listDdsAPIVersionsProduct = "dds"
	)
	listDdsAPIVersionsClient, err := conf.NewServiceClient(listDdsAPIVersionsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS Client: %s", err)
	}

	listDdsAPIVersionsPath := listDdsAPIVersionsClient.Endpoint + listDdsAPIVersionsHttpUrl

	listDdsAPIVersionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listDdsAPIVersionsOptResp, err := listDdsAPIVersionsClient.Request("GET", listDdsAPIVersionsPath,
		&listDdsAPIVersionsOpt)

	if err != nil {
		return diag.Errorf("error retrieving DDS API versions: %s", err)
	}

	listDdsAPIVersionsRespBody, err := utils.FlattenResponse(listDdsAPIVersionsOptResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("versions", flattenListDdsAPIVersionsResponseBodyInstance(listDdsAPIVersionsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDdsAPIVersionsResponseBodyInstance(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("versions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":      utils.PathSearch("id", v, nil),
			"status":  utils.PathSearch("status", v, nil),
			"updated": utils.PathSearch("updated", v, nil),
			"links":   flattenLinksDatastore(v),
		})
	}
	return rst
}

func flattenLinksDatastore(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("links", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"href": utils.PathSearch("href", v, nil),
			"rel":  utils.PathSearch("rel", v, nil),
		})
	}
	return rst
}
