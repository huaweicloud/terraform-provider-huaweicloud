package css

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParams = []string{"cluster_id", "target_image_id", "upgrade_type",
	"agency", "indices_backup_check", "cluster_load_check"}

// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/inst-type/{inst_type}/image/upgrade
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/upgrade/detail
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/upgrade/{action_id}/retry
func ResourceCssEsCoreUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCssEsCoreUpgradeCreate,
		ReadContext:   resourceCssEsCoreUpgradeRead,
		UpdateContext: resourceCssEsCoreUpgradeUpdate,
		DeleteContext: resourceCssEsCoreUpgradeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"upgrade_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"indices_backup_check": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"cluster_load_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"upgrade_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_nodes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retry_times": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceCssEsCoreUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	err = createEsCoreUpgrade(client, d, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := "detailList | [?status=='RUNNING'] | [0]"
	upgradeDetail, err := getEsUpgradeDetail(client, clusterID, expression)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", upgradeDetail, "").(string)
	if id == "" {
		return diag.Errorf("not found task ID of the core upgrade")
	}
	d.SetId(id)

	// Check whether core upgrade has completed.
	err = checkCoreUpgradeCompleted(ctx, client, clusterID, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		rollBackErr := removeEsCoreUpgradeImpact(client, d)
		if rollBackErr != nil {
			return diag.Errorf("error core upgrade: %v, rolling back upgrade fail: %s", err, rollBackErr)
		}
		return diag.FromErr(err)
	}

	return resourceCssEsCoreUpgradeRead(ctx, d, meta)
}

func resourceCssEsCoreUpgradeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	expression := fmt.Sprintf("detailList | [?id=='%s'] | [0]", d.Id())
	upgradeDetail, err := getEsUpgradeDetail(client, clusterID, expression)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "the CSS es core upgrade detail")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("upgrade_detail", flattenUpgradeDetailResponse(upgradeDetail)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCssEsCoreUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCssEsCoreUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting upgrade task resource is not supported. The resource is only removed from the state," +
		" because it is an irreversible operation."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func createEsCoreUpgrade(client *golangsdk.ServiceClient, d *schema.ResourceData, clusterID string) error {
	createEsCoreUpgradeHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/inst-type/{inst_type}/image/upgrade"
	createEsCoreUpgradePath := client.Endpoint + createEsCoreUpgradeHttpUrl
	createEsCoreUpgradePath = strings.ReplaceAll(createEsCoreUpgradePath, "{project_id}", client.ProjectID)
	createEsCoreUpgradePath = strings.ReplaceAll(createEsCoreUpgradePath, "{cluster_id}", clusterID)
	// currently, `inst_type` only supports "all".
	createEsCoreUpgradePath = strings.ReplaceAll(createEsCoreUpgradePath, "{inst_type}", "all")

	createEsCoreUpgradeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createEsCoreUpgradeOpt.JSONBody = map[string]interface{}{
		"target_image_id":      d.Get("target_image_id"),
		"upgrade_type":         d.Get("upgrade_type"),
		"agency":               d.Get("agency"),
		"indices_backup_check": d.Get("indices_backup_check"),
		"cluster_load_check":   d.Get("cluster_load_check"),
	}

	_, err := client.Request("POST", createEsCoreUpgradePath, &createEsCoreUpgradeOpt)
	if err != nil {
		return fmt.Errorf("error creating the CSS cluster core upgrade task, err: %s", err)
	}

	return nil
}

func getEsUpgradeDetail(client *golangsdk.ServiceClient, clusterID, expression string) (interface{}, error) {
	getEsUpgradeDetailHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/upgrade/detail"
	getEsUpgradeDetailPath := client.Endpoint + getEsUpgradeDetailHttpUrl
	getEsUpgradeDetailPath = strings.ReplaceAll(getEsUpgradeDetailPath, "{project_id}", client.ProjectID)
	getEsUpgradeDetailPath = strings.ReplaceAll(getEsUpgradeDetailPath, "{cluster_id}", clusterID)

	getEsUpgradeDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getEsUpgradeDetailResp, err := client.Request("GET", getEsUpgradeDetailPath, &getEsUpgradeDetailOpt)
	if err != nil {
		return getEsUpgradeDetailResp, err
	}
	getEsUpgradeDetailRespBody, err := utils.FlattenResponse(getEsUpgradeDetailResp)
	if err != nil {
		return getEsUpgradeDetailRespBody, err
	}

	upgradeDetail := utils.PathSearch(expression, getEsUpgradeDetailRespBody, nil)
	if upgradeDetail == nil {
		return upgradeDetail, golangsdk.ErrDefault404{}
	}

	return upgradeDetail, nil
}

func flattenUpgradeDetailResponse(resp interface{}) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"id":          utils.PathSearch("id", resp, nil),
			"start_time":  utils.PathSearch("startTime", resp, nil),
			"end_time":    utils.PathSearch("endTime", resp, nil),
			"status":      utils.PathSearch("status", resp, nil),
			"agency":      utils.PathSearch("agencyName", resp, nil),
			"total_nodes": utils.PathSearch("totalNodes", resp, nil),
			"retry_times": utils.PathSearch("executeTimes", resp, nil),
			"datastore":   flattenUpgradeImageInfoResponse(utils.PathSearch("imageInfo", resp, nil)),
		},
	}

	return rst
}

func flattenUpgradeImageInfoResponse(resp interface{}) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"type":    utils.PathSearch("datastoreType", resp, nil),
			"version": utils.PathSearch("datastoreVersion", resp, nil),
		},
	}

	return rst
}

func checkCoreUpgradeCompleted(ctx context.Context, client *golangsdk.ServiceClient, clusterId, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCESS"},
		Refresh:      coreUpgradeStateRefreshFunc(client, clusterId, id),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) es cluster upgrade task: %s", clusterId, err)
	}
	return nil
}

func coreUpgradeStateRefreshFunc(client *golangsdk.ServiceClient, clusterID, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		expression := fmt.Sprintf("detailList | [?id=='%s'] | [0]", id)
		resp, err := getEsUpgradeDetail(client, clusterID, expression)
		if err != nil {
			return resp, "ERROR", err
		}

		status := utils.PathSearch("status", resp, "").(string)
		return resp, status, nil
	}
}

func removeEsCoreUpgradeImpact(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	retryEsCoreUpgradeHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/upgrade/{action_id}/retry"
	retryEsCoreUpgradePath := client.Endpoint + retryEsCoreUpgradeHttpUrl
	retryEsCoreUpgradePath = strings.ReplaceAll(retryEsCoreUpgradePath, "{project_id}", client.ProjectID)
	retryEsCoreUpgradePath = strings.ReplaceAll(retryEsCoreUpgradePath, "{cluster_id}", d.Get("cluster_id").(string))
	retryEsCoreUpgradePath = strings.ReplaceAll(retryEsCoreUpgradePath, "{action_id}", d.Id())

	retryEsCoreUpgradeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	retryEsCoreUpgradeOpt.JSONBody = map[string]interface{}{
		"retry_mode": "abort",
	}

	_, err := client.Request("PUT", retryEsCoreUpgradePath, &retryEsCoreUpgradeOpt)
	if err != nil {
		return err
	}

	return nil
}
