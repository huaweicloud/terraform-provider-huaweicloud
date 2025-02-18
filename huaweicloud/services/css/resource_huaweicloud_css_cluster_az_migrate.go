package css

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var azMigrateNonUpdatableParams = []string{"cluster_id", "instance_type", "source_az", "target_az",
	"migrate_type", "agency", "indices_backup_check"}

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/inst-type/{inst_type}/azmigrate
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/upgrade/detail
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/upgrade/{action_id}/retry
// @API CSS POST /v2.0/{project_id}/clusters/{cluster_id}/restart
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}
func ResourceCssClusterAzMigrate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCssClusterAzMigrateCreate,
		ReadContext:   resourceCssClusterAzMigrateRead,
		UpdateContext: resourceCssClusterAzMigrateUpdate,
		DeleteContext: resourceCssClusterAzMigrateDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(azMigrateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_az": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_az": {
				Type:     schema.TypeString,
				Required: true,
			},
			"migrate_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"indices_backup_check": {
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceCssClusterAzMigrateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	// Check whether the cluster status is available.
	err = checkClusterOperationResult(ctx, client, clusterID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	err = createAzMigrate(client, d, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := "detailList | [?status=='RUNNING'] | [0]"
	taskDetail, err := getAzMigrateDetail(client, clusterID, expression)
	if err != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("id", taskDetail, "").(string)
	if id == "" {
		return diag.Errorf("not found task ID of the az migrate task")
	}
	d.SetId(id)

	// Check whether az migrate has completed.
	err = checkAzMigrateCompleted(ctx, client, clusterID, id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		rollBackErr := removeEsCoreUpgradeImpact(client, d)
		if rollBackErr != nil {
			return diag.Errorf("error az migtare: %v, rolling az migrate fail: %s", err, rollBackErr)
		}
		return diag.FromErr(err)
	}

	// After the availability zone migration is successful, the cluster needs to be restarted before it can be used.
	err = restartCluster(client, clusterID, "role", "all")
	if err != nil {
		return diag.Errorf("error restart CSS cluster, err: %s", err)
	}

	// Check whether the cluster status is available.
	err = checkClusterOperationResult(ctx, client, clusterID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCssClusterAzMigrateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCssClusterAzMigrateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCssClusterAzMigrateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting cluster az migrate resource is not supported. The  resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func createAzMigrate(client *golangsdk.ServiceClient, d *schema.ResourceData, clusterID string) error {
	instType := d.Get("instance_type").(string)
	createAzMigrateHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/inst-type/{inst_type}/azmigrate"
	createAzMigratePath := client.Endpoint + createAzMigrateHttpUrl
	createAzMigratePath = strings.ReplaceAll(createAzMigratePath, "{project_id}", client.ProjectID)
	createAzMigratePath = strings.ReplaceAll(createAzMigratePath, "{cluster_id}", clusterID)
	createAzMigratePath = strings.ReplaceAll(createAzMigratePath, "{inst_type}", instType)

	updateInstanceAzOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateInstanceAzOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
		"source_az":            d.Get("source_az"),
		"target_az":            d.Get("target_az"),
		"migrate_type":         d.Get("migrate_type"),
		"agency":               d.Get("agency"),
		"indices_backup_check": utils.ValueIgnoreEmpty(d.Get("indices_backup_check")),
	})
	_, err := client.Request("POST", createAzMigratePath, &updateInstanceAzOpt)
	if err != nil {
		return fmt.Errorf("error creating CSS cluster az migrate task: %s", err)
	}

	return nil
}

func checkAzMigrateCompleted(ctx context.Context, client *golangsdk.ServiceClient, clusterId, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCESS"},
		Refresh:      azMigrateStateRefreshFunc(client, clusterId, id),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) es cluster az migrate task: %s", clusterId, err)
	}
	return nil
}

func azMigrateStateRefreshFunc(client *golangsdk.ServiceClient, clusterID, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		expression := fmt.Sprintf("detailList | [?id=='%s'] | [0]", id)
		resp, err := getAzMigrateDetail(client, clusterID, expression)
		if err != nil {
			return resp, "ERROR", err
		}

		status := utils.PathSearch("status", resp, "").(string)
		return resp, status, nil
	}
}

func getAzMigrateDetail(client *golangsdk.ServiceClient, clusterID, expression string) (interface{}, error) {
	getAzMigrateDetailHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/upgrade/detail?action_mode=AZ_MIGRATION"
	getAzMigrateDetailPath := client.Endpoint + getAzMigrateDetailHttpUrl
	getAzMigrateDetailPath = strings.ReplaceAll(getAzMigrateDetailPath, "{project_id}", client.ProjectID)
	getAzMigrateDetailPath = strings.ReplaceAll(getAzMigrateDetailPath, "{cluster_id}", clusterID)

	getAzMigrateDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getAzMigrateDetailResp, err := client.Request("GET", getAzMigrateDetailPath, &getAzMigrateDetailOpt)
	if err != nil {
		return getAzMigrateDetailResp, err
	}
	getAzMigrateDetailRespBody, err := utils.FlattenResponse(getAzMigrateDetailResp)
	if err != nil {
		return getAzMigrateDetailRespBody, err
	}

	taskDetail := utils.PathSearch(expression, getAzMigrateDetailRespBody, nil)
	if taskDetail == nil {
		return taskDetail, golangsdk.ErrDefault404{}
	}

	return taskDetail, nil
}
