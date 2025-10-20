package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var syncNetworkPolicyNonUpdatableParams = []string{
	"cluster_id", "enterprise_project_id",
}

// @API HSS GET /v5/{project_id}/container-network/{cluster_id}/policy-sync
func ResourceContainerNetworkPolicySync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerNetworkPolicySyncCreate,
		ReadContext:   resourceContainerNetworkPolicySyncRead,
		UpdateContext: resourceContainerNetworkPolicySyncUpdate,
		DeleteContext: resourceContainerNetworkPolicySyncDelete,

		CustomizeDiff: config.FlexibleForceNew(syncNetworkPolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located. If omitted, the provider-level region will be used.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the cluster to synchronize network policies for.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
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

func buildContainerNetworkPolicySyncQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	rst := ""
	if epsID := cfg.GetEnterpriseProjectID(d); epsID != "" {
		rst = fmt.Sprintf("?enterprise_project_id=%s", epsID)
	}
	return rst
}

func resourceContainerNetworkPolicySyncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "hss"
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := "v5/{project_id}/container-network/{cluster_id}/policy-sync"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterId)
	requestPath += buildContainerNetworkPolicySyncQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("GET", client.Endpoint+requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error synchronizing container network policies for cluster %s: %s", clusterId, err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(id)

	return resourceContainerNetworkPolicySyncRead(ctx, d, meta)
}

func resourceContainerNetworkPolicySyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceContainerNetworkPolicySyncUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceContainerNetworkPolicySyncDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to synchronize container network policies. Deleting
	this resource will not affect the synchronization status, but will only remove the resource information from
	the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
