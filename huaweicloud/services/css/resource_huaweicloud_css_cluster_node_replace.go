package css

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nodeReplaceNonUpdatableParams = []string{"cluster_id", "node_id", "agency", "migrate_data"}

// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/instance/{instance_id}/replace
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}
func ResourceCssClusterNodeReplace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeReplaceCreate,
		ReadContext:   resourceNodeReplaceRead,
		UpdateContext: resourceNodeReplaceUpdate,
		DeleteContext: resourceNodeReplaceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(nodeReplaceNonUpdatableParams),

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agency": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"migrate_data": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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

func resourceNodeReplaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	err = createNodeReplace(client, d, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check whether the cluster node replacement is complete
	err = checkClusterOperationResult(ctx, client, clusterID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterID)

	return nil
}

func resourceNodeReplaceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNodeReplaceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNodeReplaceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting node replace resource is not supported. The resource is only removed from the state," +
		" the cluster instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func createNodeReplace(client *golangsdk.ServiceClient, d *schema.ResourceData, clusterID string) error {
	createNodeReplaceHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/instance/{instance_id}/replace"
	createNodeReplacePath := client.Endpoint + createNodeReplaceHttpUrl
	createNodeReplacePath = strings.ReplaceAll(createNodeReplacePath, "{project_id}", client.ProjectID)
	createNodeReplacePath = strings.ReplaceAll(createNodeReplacePath, "{cluster_id}", clusterID)
	createNodeReplacePath = strings.ReplaceAll(createNodeReplacePath, "{instance_id}", d.Get("node_id").(string))

	createNodeReplaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createNodeReplacePath += fmt.Sprintf("?migrateData=%t", d.Get("migrate_data").(bool))
	if agency, ok := d.GetOk("agency"); ok {
		createNodeReplacePath += fmt.Sprintf("&agency=%s", agency.(string))
	}

	_, err := client.Request("PUT", createNodeReplacePath, &createNodeReplaceOpt)
	if err != nil {
		return fmt.Errorf("error creating the CSS cluster node replace, err: %s", err)
	}

	return nil
}
