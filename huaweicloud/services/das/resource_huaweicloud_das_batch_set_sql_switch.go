package das

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var batchSetSqlSwitchNonUpdatableParams = []string{
	"engine_type",
	"switch_on",
	"switch_type",
	"instance_ids",
	"retention_hours",
}

// @API DAS POST /v3/{project_id}/instance/batch-set-sql-switch
func ResourceBatchSetSqlSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchSetSqlSwitchCreate,
		ReadContext:   resourceBatchSetSqlSwitchRead,
		UpdateContext: resourceBatchSetSqlSwitchUpdate,
		DeleteContext: resourceBatchSetSqlSwitchDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(batchSetSqlSwitchNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the SQL switch is located.",
			},

			// Required parameters.
			"engine_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The engine type of the instances.",
			},
			"switch_on": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable the SQL switch.",
			},
			"switch_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of SQL switch to set.",
			},
			"instance_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of instance IDs.",
			},
			"retention_hours": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The retention hours of the SQL data.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					"Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.",
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildBatchSetSqlSwitchBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"switch_on":    d.Get("switch_on").(bool),
		"switch_type":  d.Get("switch_type").(string),
		"instance_ids": d.Get("instance_ids").([]interface{}),
		"engine_type":  d.Get("engine_type").(string),
	}

	if retentionHours, ok := d.GetOk("retention_hours"); ok {
		bodyParams["retention_hours"] = retentionHours.(int)
	}

	return bodyParams
}

func resourceBatchSetSqlSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instance/batch-set-sql-switch"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: buildBatchSetSqlSwitchBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error setting DAS SQL switch batch action: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return resourceBatchSetSqlSwitchRead(ctx, d, meta)
}

func resourceBatchSetSqlSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBatchSetSqlSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBatchSetSqlSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for batch setting SQL switch. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
