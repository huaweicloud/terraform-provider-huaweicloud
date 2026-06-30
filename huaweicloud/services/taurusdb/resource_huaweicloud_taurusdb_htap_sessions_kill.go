package taurusdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	htapSessionsKillNonUpdatableParams = []string{"instance_id", "process_list"}
)

// @API TaurusDB DELETE /v3/{project_id}/instances/{instance_id}/htap/process
func ResourceTaurusDBHtapSessionsKill() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapSessionsKillCreate,
		ReadContext:   resourceTaurusDBHtapSessionsKillRead,
		UpdateContext: resourceTaurusDBHtapSessionsKillUpdate,
		DeleteContext: resourceTaurusDBHtapSessionsKillDelete,

		CustomizeDiff: config.FlexibleForceNew(htapSessionsKillNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"process_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceTaurusDBHtapSessionsKillCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/htap/process"
		product    = "gaussdb"
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProviderClient.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildCreateTaurusDBHtapSessionsKillBodyParams(d),
	}

	resp, err := client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error killing TaurusDB HTAP instance sessions: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing response: %s", err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("msg", utils.PathSearch("msg", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCreateTaurusDBHtapSessionsKillBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"process_list": utils.ExpandToStringList(d.Get("process_list").([]interface{})),
	}
	return bodyParams
}

func resourceTaurusDBHtapSessionsKillRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceTaurusDBHtapSessionsKillUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceTaurusDBHtapSessionsKillDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting TaurusDB HTAP sessions kill resource is not supported. The TaurusDB HTAP sessions " +
		"kill resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
