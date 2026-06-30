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
	nodeSessionsNonUpdatableParams = []string{"instance_id", "node_id", "processes"}
)

// @API TaurusDB DELETE /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/processes
func ResourceTaurusDBNodeSessionsKill() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBNodeSessionsKillCreate,
		ReadContext:   resourceTaurusDBNodeSessionsKillRead,
		UpdateContext: resourceTaurusDBNodeSessionsKillUpdate,
		DeleteContext: resourceTaurusDBNodeSessionsKillDelete,

		CustomizeDiff: config.FlexibleForceNew(nodeSessionsNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the TaurusDB instance.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the node in the TaurusDB instance.`,
			},
			"processes": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `Specifies the IDs of user session threads to be terminated.`,
			},
			"processes_killed": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `The IDs of terminated user session threads in requested processes.`,
			},
			"processes_not_found": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `The IDs of user session threads that were not found in requested processes.`,
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

func resourceTaurusDBNodeSessionsKillCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/processes"
		product    = "gaussdb"
		instanceId = d.Get("instance_id").(string)
		nodeId     = d.Get("node_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)
	requestPath = strings.ReplaceAll(requestPath, "{node_id}", nodeId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildCreateTaurusDBNodeSessionsKillBodyParams(d),
	}

	resp, err := client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error killing TaurusDB node sessions: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing response: %s", err)
	}

	processesKilled := utils.PathSearch("processes_killed", respBody, make([]interface{}, 0))
	killedIntIds := make([]int, 0)
	if killedList, ok := processesKilled.([]interface{}); ok {
		killedIntIds = expandToIntList(killedList)
	}

	processesNotFound := utils.PathSearch("processes_not_found", respBody, make([]interface{}, 0))
	notFoundIntIds := make([]int, 0)
	if notFoundList, ok := processesNotFound.([]interface{}); ok {
		notFoundIntIds = expandToIntList(notFoundList)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("processes_killed", killedIntIds),
		d.Set("processes_not_found", notFoundIntIds),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting TaurusDB node sessions delete state: %s", mErr)
	}
	return nil
}

func buildCreateTaurusDBNodeSessionsKillBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"processes": utils.ExpandToIntList(d.Get("processes").([]interface{})),
	}
	return bodyParams
}

func expandToIntList(v []interface{}) []int {
	s := make([]int, 0, len(v))
	for _, val := range v {
		if intVal, ok := val.(int); ok {
			s = append(s, intVal)
		}
		if fv, ok := val.(float64); ok {
			s = append(s, int(fv))
		}
	}
	return s
}

func resourceTaurusDBNodeSessionsKillRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}
func resourceTaurusDBNodeSessionsKillUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}
func resourceTaurusDBNodeSessionsKillDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting TaurusDB node processes delete resource is not supported. The TaurusDB node processes " +
		"delete resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
