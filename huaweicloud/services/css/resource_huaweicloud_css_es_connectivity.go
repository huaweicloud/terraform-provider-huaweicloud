package css

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var esConnectivityNonUpdatableParams = []string{"source_cluster_id", "target_cluster_id"}

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/logs/connectivity
func ResourceEsConnectivity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEsConnectivityCreate,
		ReadContext:   resourceEsConnectivityRead,
		UpdateContext: resourceEsConnectivityUpdate,
		DeleteContext: resourceEsConnectivityDelete,

		CustomizeDiff: config.FlexibleForceNew(esConnectivityNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_cluster_id": {
				Type:     schema.TypeString,
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

func resourceEsConnectivityCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterId := d.Get("source_cluster_id").(string)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	createConnectivityHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/connectivity"
	createConnectivityPath := client.Endpoint + createConnectivityHttpUrl
	createConnectivityPath = strings.ReplaceAll(createConnectivityPath, "{project_id}", client.ProjectID)
	createConnectivityPath = strings.ReplaceAll(createConnectivityPath, "{cluster_id}", clusterId)

	createConnectivityOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createConnectivityOpt.JSONBody = map[string]interface{}{
		"target_cluster_id": d.Get("target_cluster_id"),
	}

	_, err = client.Request("POST", createConnectivityPath, &createConnectivityOpt)
	if err != nil {
		return diag.Errorf("error creating the CSS cluster connectivity test: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	return nil
}

func resourceEsConnectivityRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEsConnectivityUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEsConnectivityDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ES connectivity resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
