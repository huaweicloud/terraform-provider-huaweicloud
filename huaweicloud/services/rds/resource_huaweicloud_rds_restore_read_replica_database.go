package rds

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var restoreReadReplicaDatabaseNonUpdatableParams = []string{"instance_id", "databases"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/log-replay/database
func ResourceRdsRestoreReadReplicaDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRestoreReadReplicaDatabaseCreate,
		UpdateContext: resourceRestoreReadReplicaDatabaseUpdate,
		ReadContext:   resourceRestoreReadReplicaDatabaseRead,
		DeleteContext: resourceRestoreReadReplicaDatabaseDelete,

		CustomizeDiff: config.FlexibleForceNew(restoreReadReplicaDatabaseNonUpdatableParams),

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
			"databases": {
				Type:     schema.TypeSet,
				Elem:     restoreReadReplicaDatabaseSchema(),
				Required: true,
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

func restoreReadReplicaDatabaseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"old_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"new_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRestoreReadReplicaDatabaseCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	httpUrl := "v3/{project_id}/instances/{instance_id}/log-replay/database"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
			"x-language":   "en-us",
		},
		JSONBody: utils.RemoveNil(buildRestoreReadReplicaDatabaseBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error restoring read replica database for instance (%s): %s", instanceID, err)
	}

	d.SetId(instanceID)

	return nil
}

func buildRestoreReadReplicaDatabaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawDatabases := d.Get("databases").(*schema.Set).List()

	if len(rawDatabases) == 0 {
		return nil
	}

	databases := make([]map[string]interface{}, len(rawDatabases))

	for i, d := range rawDatabases {
		dMap := d.(map[string]interface{})
		databases[i] = map[string]interface{}{
			"old_name": dMap["old_name"],
			"new_name": dMap["new_name"],
		}
	}
	return map[string]interface{}{
		"databases": databases,
	}
}

func resourceRestoreReadReplicaDatabaseUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRestoreReadReplicaDatabaseRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRestoreReadReplicaDatabaseDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS PostgreSQL replica-to-primary restore resource is not supported. " +
		"The resource is not only removed from the state."
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
