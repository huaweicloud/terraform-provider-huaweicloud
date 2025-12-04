package lb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v2/{project_id}/elb/whitelists
func DataSourceLbWhitelists() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLbWhitelistsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"whitelist_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_whitelist": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"whitelist": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"whitelists": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     whitelistsWhitelistsSchema(),
			},
		},
	}
}

func whitelistsWhitelistsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_whitelist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"whitelist": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceLbWhitelistsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v2/{project_id}/elb/whitelists"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getQueryParams := buildGetLbWhitelistsQueryParams(d)
	getPath += getQueryParams

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving ELB whitelists: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("whitelists", flattenWhitelistsBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetLbWhitelistsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("whitelist_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("enable_whitelist"); ok {
		enableWhitelist, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&enable_whitelist=%v", res, enableWhitelist)
	}
	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&listener_id=%v", res, v)
	}
	if v, ok := d.GetOk("whitelist"); ok {
		res = fmt.Sprintf("%s&whitelist=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenWhitelistsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("whitelists", resp, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"listener_id":      utils.PathSearch("listener_id", v, nil),
			"enable_whitelist": utils.PathSearch("enable_whitelist", v, nil),
			"whitelist":        utils.PathSearch("whitelist", v, nil),
		})
	}
	return rst
}
