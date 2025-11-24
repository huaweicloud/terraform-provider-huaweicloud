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

var containerNetworkClusterSyncNonUpdatableParams = []string{
	"enterprise_project_id",
}

// @API HSS GET /v5/{project_id}/container-network/cluster/sync
func ResourceContainerNetworkClusterSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerNetworkClusterSyncCreate,
		ReadContext:   resourceContainerNetworkClusterSyncRead,
		UpdateContext: resourceContainerNetworkClusterSyncUpdate,
		DeleteContext: resourceContainerNetworkClusterSyncDelete,

		CustomizeDiff: config.FlexibleForceNew(containerNetworkClusterSyncNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
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

func buildContainerNetworkClusterSyncQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func resourceContainerNetworkClusterSyncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/container-network/cluster/sync"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerNetworkClusterSyncQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error syncing HSS container network cluster: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceContainerNetworkClusterSyncRead(ctx, d, meta)
}

func resourceContainerNetworkClusterSyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceContainerNetworkClusterSyncUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceContainerNetworkClusterSyncDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to sync HSS container network cluster. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from the
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
