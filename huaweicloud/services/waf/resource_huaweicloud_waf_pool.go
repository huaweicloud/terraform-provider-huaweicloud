package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsPool = []string{
	"name",
	"type",
	"vpc_id",
	"description",
	"enterprise_project_id",
}

// @API WAF POST /v1/{project_id}/premium-waf/pool
// @API WAF GET /v1/{project_id}/premium-waf/pool/{pool_id}
// @API WAF DELETE /v1/{project_id}/premium-waf/pool/{pool_id}
func ResourceWafPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafPoolCreate,
		ReadContext:   resourceWafPoolRead,
		UpdateContext: resourceWafPoolUpdate,
		DeleteContext: resourceWafPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFPoolImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsPool),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// The query API does not return, so the `Computed` attribute is not added.
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     poolIdNameEntrySchema(),
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     poolIdNameEntrySchema(),
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func poolIdNameEntrySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildWafPoolQueryParams(epsId string) string {
	if epsId == "" {
		return ""
	}

	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func buildCreatePoolBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"region":      region,
		"type":        d.Get("type"),
		"vpc_id":      d.Get("vpc_id"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceWafPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/premium-waf/pool"
		product = "waf"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildWafPoolQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildCreatePoolBodyParams(d, region)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating WAF pool: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating WAF pool: ID is not found in API response")
	}

	d.SetId(id)

	return resourceWafPoolRead(ctx, d, meta)
}

func GetWafPool(client *golangsdk.ServiceClient, poolId, epsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/premium-waf/pool/{pool_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{pool_id}", poolId)
	getPath += buildWafPoolQueryParams(epsId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceWafPoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "waf"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	respBody, err := GetWafPool(client, d.Id(), epsId)
	if err != nil {
		// If the pool does not exist, the response HTTP status code of the details API is `404`.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF pool")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("hosts", flattenPoolIdNameEntries(
			utils.PathSearch("hosts", respBody, make([]interface{}, 0)))),
		d.Set("instances", flattenPoolIdNameEntries(
			utils.PathSearch("instances", respBody, make([]interface{}, 0)))),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPoolIdNameEntries(raw interface{}) []map[string]interface{} {
	entries, ok := raw.([]interface{})
	if !ok || len(entries) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(entries))
	for _, v := range entries {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"service_ip": utils.PathSearch("service_ip", v, nil),
		})
	}

	return rst
}

func resourceWafPoolUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWafPoolDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/premium-waf/pool/{pool_id}"
		product = "waf"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{pool_id}", d.Id())
	deletePath += buildWafPoolQueryParams(epsId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the pool does not exist, the response HTTP status code of the deletion API is `404`.
		return common.CheckDeletedDiag(d, err, "error deleting WAF pool")
	}

	return nil
}

func resourceWAFPoolImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <id>/<enterprise_project_id>")
	}

	d.SetId(parts[0])

	return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[1])
}
