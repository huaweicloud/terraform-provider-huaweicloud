package dns

import (
	"context"
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

var zoneRetrievalNonUpdatableParams = []string{
	"zone_name",
}

var dnsZoneRetrievalSchema = map[string]*schema.Schema{
	"region": {
		Type:        schema.TypeString,
		ForceNew:    true,
		Optional:    true,
		Computed:    true,
		Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
	},
	"zone_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: `Specifies the zone name.`,
		DiffSuppressFunc: func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
			return strings.TrimSuffix(oldVal, ".") == strings.TrimSuffix(newVal, ".")
		},
	},
	"enable_force_new": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
		Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
	},
	"retrieval_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `Indicates the retrieval ID`,
	},
	"status": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `Indicates the status.`,
	},
	"created_at": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `Indicates the create time.`,
	},
	"updated_at": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `Indicates the last update time.`,
	},
	"record": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: `Indicates the record detail.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the record host.`,
				},
				"value": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the record value.`,
				},
			},
		},
	},
}

// @API DNS POST /v2/retrieval
// @API DNS GET /v2/retrieval
func ResourceDNSZoneRetrieval() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSZoneRetrievalCreate,
		UpdateContext: resourceDNSZoneRetrievalUpdate,
		ReadContext:   resourceDNSZoneRetrievalRead,
		DeleteContext: resourceDNSZoneRetrievalDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(zoneRetrievalNonUpdatableParams, dnsZoneRetrievalSchema),

		Schema: dnsZoneRetrievalSchema,
	}
}

func resourceDNSZoneRetrievalCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	zoneName := d.Get("zone_name").(string)
	createHttpUrl := "v2/retrieval"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"zone_name": zoneName,
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DNS retrieval: %s", err)
	}

	d.SetId(zoneName)

	return resourceDNSZoneRetrievalRead(ctx, d, meta)
}

func resourceDNSZoneRetrievalRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	getHttpUrl := "v2/retrieval?name={name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{name}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS zone retrieval.")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("zone_name", utils.PathSearch("zone_name", getRespBody, nil)),
		d.Set("retrieval_id", utils.PathSearch("id", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
		d.Set("record", flattenDNSZoneRetrievalRecord(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDNSZoneRetrievalRecord(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("record", resp, nil)
	if param, ok := rawParams.(map[string]interface{}); ok {
		m := map[string]interface{}{
			"host":  utils.PathSearch("host", param, nil),
			"value": utils.PathSearch("value", param, nil),
		}
		return []interface{}{m}
	}

	return nil
}

func resourceDNSZoneRetrievalUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDNSZoneRetrievalDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DNS zone retrieval resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
