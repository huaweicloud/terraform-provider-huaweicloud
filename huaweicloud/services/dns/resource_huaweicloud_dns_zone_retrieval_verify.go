package dns

import (
	"context"
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

var zoneRetrievalVerifyNonUpdatableParams = []string{
	"retrieval_id",
}

// @API DNS POST /v2/retrieval/verification/{id}
func ResourceDNSZoneRetrievalVerify() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSZoneRetrievalVerifyCreate,
		UpdateContext: resourceDNSZoneRetrievalVerifyUpdate,
		ReadContext:   resourceDNSZoneRetrievalVerifyRead,
		DeleteContext: resourceDNSZoneRetrievalVerifyDelete,

		CustomizeDiff: config.FlexibleForceNew(zoneRetrievalVerifyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"retrieval_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the retrieval ID.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the message.`,
			},
		},
	}
}

func resourceDNSZoneRetrievalVerifyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	createHttpUrl := "v2/retrieval/verification/{id}"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{id}", d.Get("retrieval_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{202},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DNS retrieval verify: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("message", utils.PathSearch("message", createRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDNSZoneRetrievalVerifyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDNSZoneRetrievalVerifyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDNSZoneRetrievalVerifyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DNS zone retrieval verify resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
