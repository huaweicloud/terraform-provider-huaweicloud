package dsc

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

var batchDeleteDatabaseInstanceNonUpdatableParams = []string{"db_ids"}

// @API DSC POST /v1/{project_id}/asset-center/database/instances/batch-delete
func ResourceBatchDeleteDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchDeleteDatabaseInstanceCreate,
		ReadContext:   resourceBatchDeleteDatabaseInstanceRead,
		UpdateContext: resourceBatchDeleteDatabaseInstanceUpdate,
		DeleteContext: resourceBatchDeleteDatabaseInstanceDelete,

		CustomizeDiff: config.FlexibleForceNew(batchDeleteDatabaseInstanceNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"db_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildBatchDeleteDatabaseInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"db_ids": utils.ExpandToStringList(d.Get("db_ids").([]interface{})),
	}
}

func resourceBatchDeleteDatabaseInstanceCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/asset-center/database/instances/batch-delete"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildBatchDeleteDatabaseInstanceBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch deleting DSC database instances: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("msg", utils.PathSearch("msg", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBatchDeleteDatabaseInstanceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteDatabaseInstanceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteDatabaseInstanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch delete DSC database instances.
Deleting this resource will not restore the deleted database instances or undo the delete action, but will only
remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
