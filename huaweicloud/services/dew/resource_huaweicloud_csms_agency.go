package dew

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

var agencyNonUpdatableParams = []string{"secret_type"}

// @API DEW POST /v1/csms/agencies
func ResourceAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgencyCreate,
		ReadContext:   resourceAgencyRead,
		UpdateContext: resourceAgencyUpdate,
		DeleteContext: resourceAgencyDelete,

		CustomizeDiff: config.FlexibleForceNew(agencyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secret_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"agencies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agency_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agency_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAgencyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"secret_type": d.Get("secret_type"),
	}

	return bodyParams
}

func resourceAgencyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/csms/agencies"
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAgencyBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating agency: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	agencyList := utils.PathSearch("agencies", respBody, make([]interface{}, 0)).([]interface{})

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("agencies", flattenAgencies(agencyList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAgencies(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"agency_id":   utils.PathSearch("agency_id", v, nil),
			"agency_name": utils.PathSearch("agency_name", v, nil),
		})
	}

	return rst
}

func resourceAgencyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceAgencyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceAgencyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
