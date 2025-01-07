package cce

import (
	"context"
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

// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}
// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/clustercertrevoke
var certificateRevokeNonUpdatableParams = []string{"cluster_id", "user_id", "agency_id"}

func ResourceCertificateRevoke() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateRevokeCreate,
		ReadContext:   resourceCertificateRevokeRead,
		UpdateContext: resourceCertificateRevokeUpdate,
		DeleteContext: resourceCertificateRevokeDelete,

		CustomizeDiff: config.FlexibleForceNew(certificateRevokeNonUpdatableParams),

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
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"user_id", "agency_id"},
			},
			"agency_id": {
				Type:     schema.TypeString,
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

func buildCertificateRevokeCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"userId":   utils.ValueIgnoreEmpty(d.Get("user_id")),
		"agencyId": utils.ValueIgnoreEmpty(d.Get("agency_id")),
	}

	return bodyParams
}

func resourceCertificateRevokeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	// Wait for the cce cluster to become available
	clusterID := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(client, clusterID, []string{"Available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateCluster.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE cluster to become available: %s", err)
	}

	var (
		createCertificateRevokeHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/clustercertrevoke"
	)

	createCertificateRevokePath := client.Endpoint + createCertificateRevokeHttpUrl
	createCertificateRevokePath = strings.ReplaceAll(createCertificateRevokePath, "{project_id}", client.ProjectID)
	createCertificateRevokePath = strings.ReplaceAll(createCertificateRevokePath, "{cluster_id}", d.Get("cluster_id").(string))

	createCertificateRevokeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCertificateRevokeCreateOpts(d)),
	}

	_, err = client.Request("POST", createCertificateRevokePath, &createCertificateRevokeOpt)
	if err != nil {
		return diag.Errorf("error revoking CCE cluster certificate: %s", err)
	}

	d.SetId(d.Get("cluster_id").(string))

	return resourceCertificateRevokeRead(ctx, d, meta)
}

func resourceCertificateRevokeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCertificateRevokeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCertificateRevokeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting certificate revoke resource is not supported. The certificate revoke resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
