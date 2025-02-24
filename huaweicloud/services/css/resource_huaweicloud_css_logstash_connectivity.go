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

var ConnectivityNonUpdatableParams = []string{"cluster_id", "address_and_ports",
	"address_and_ports.*.address", "address_and_ports.*.port"}

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/checkconnection
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}
func ResourceLogstashConnectivity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogstashConnectivityCreate,
		ReadContext:   resourceLogstashConnectivityRead,
		UpdateContext: resourceLogstashConnectivityUpdate,
		DeleteContext: resourceLogstashConnectivityDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(ConnectivityNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address_and_ports": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"connectivity_results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceLogstashConnectivityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	resp, err := createConnectivity(client, d, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check whether the cluster restart is complete.
	err = checkClusterOperationResult(ctx, client, clusterID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterID)

	d.Set("connectivity_results", flattenConnectivityResults(resp))

	return nil
}

func resourceLogstashConnectivityRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLogstashConnectivityUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLogstashConnectivityDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting connectivity resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func createConnectivity(client *golangsdk.ServiceClient, d *schema.ResourceData, clusterID string) (interface{}, error) {
	createConnectivityHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/checkconnection"
	createConnectivityPath := client.Endpoint + createConnectivityHttpUrl
	createConnectivityPath = strings.ReplaceAll(createConnectivityPath, "{project_id}", client.ProjectID)
	createConnectivityPath = strings.ReplaceAll(createConnectivityPath, "{cluster_id}", clusterID)

	createConnectivityOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createConnectivityOpt.JSONBody = map[string]interface{}{
		"addressAndPorts": d.Get("address_and_ports"),
	}

	resp, err := client.Request("POST", createConnectivityPath, &createConnectivityOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating the CSS logstash connectivity test, err: %s", err)
	}

	createConnectivityRespBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return createConnectivityRespBody, nil
}

func flattenConnectivityResults(resp interface{}) []interface{} {
	curJson := utils.PathSearch("result", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		status := utils.PathSearch("status", v, float64(0)).(float64)
		rst = append(rst, map[string]interface{}{
			"address": utils.PathSearch("address", v, nil),
			"port":    utils.PathSearch("port", v, nil),
			"status":  convertConnectivityStatus(int(status)),
		})
	}
	return rst
}

func convertConnectivityStatus(status int) string {
	switch status {
	case 1:
		return "connection succeeded"
	case 0:
		return "address unreachable"
	case 2:
		return "port unreachable"
	case 3:
		return "domain name cannot be resolved"
	case -2:
		return "wrong location"
	}

	return ""
}
