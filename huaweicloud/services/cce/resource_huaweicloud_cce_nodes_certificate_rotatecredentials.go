package cce

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

var rotateNodesCredentialsNonUpdatableParams = []string{
	"cluster_id", "api_version", "kind", "node_list", "node_list.*.node_id",
}

// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/rotate-cert
func ResourceRotateNodesCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRotateNodesCredentialsCreate,
		ReadContext:   resourceRotateNodesCredentialsRead,
		UpdateContext: resourceRotateNodesCredentialsUpdate,
		DeleteContext: resourceRotateNodesCredentialsDelete,

		CustomizeDiff: config.FlexibleForceNew(rotateNodesCredentialsNonUpdatableParams),

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
			"api_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
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
		},
	}
}

func resourceRotateNodesCredentialsCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterID = d.Get("cluster_id").(string)
		httpUrl   = "api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/rotate-cert"
	)
	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = buildRotateNodesCredentialsBodyParams(d)

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error rotating CCE cluster nodes credentials: %s", err)
	}

	d.SetId(clusterID)

	return nil
}

func buildRotateNodesCredentialsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"apiVersion": d.Get("api_version"),
		"kind":       d.Get("kind"),
	}

	if v, ok := d.GetOk("node_list"); ok {
		bodyParams["nodeList"] = buildNodeListBodyParams(v.([]interface{}))
	}

	return bodyParams
}

func buildNodeListBodyParams(nodeList []interface{}) []map[string]interface{} {
	if len(nodeList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(nodeList))
	for i, v := range nodeList {
		node := v.(map[string]interface{})
		result[i] = map[string]interface{}{
			"nodeID": node["node_id"],
		}
	}
	return result
}

func resourceRotateNodesCredentialsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRotateNodesCredentialsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRotateNodesCredentialsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for rotating cluster node credentials. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
