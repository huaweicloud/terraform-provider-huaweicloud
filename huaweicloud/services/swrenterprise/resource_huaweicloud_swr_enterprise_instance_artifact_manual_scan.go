package swrenterprise

import (
	"context"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var enterpriseInstanceArtifactManualScanNonUpdatableParams = []string{
	"instance_id", "namespace_name", "repository_name", "reference",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}/artifacts/{reference}/scan
func ResourceSwrEnterpriseInstanceArtifactManualScan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseInstanceArtifactManualScanCreate,
		UpdateContext: resourceSwrEnterpriseInstanceArtifactManualScanUpdate,
		ReadContext:   resourceSwrEnterpriseInstanceArtifactManualScanRead,
		DeleteContext: resourceSwrEnterpriseInstanceArtifactManualScanDelete,

		CustomizeDiff: config.FlexibleForceNew(enterpriseInstanceArtifactManualScanNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace name.`,
			},
			"repository_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the repository name.`,
			},
			"reference": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the artifact digest.`,
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

func resourceSwrEnterpriseInstanceArtifactManualScanCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	repositoryName := d.Get("repository_name").(string)
	reference := d.Get("reference").(string)

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}/artifacts/{reference}/scan"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{namespace_name}", namespaceName)
	deletePath = strings.ReplaceAll(deletePath, "{repository_name}", url.PathEscape(strings.ReplaceAll(repositoryName, "/", "%2F")))
	deletePath = strings.ReplaceAll(deletePath, "{reference}", reference)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error scanning SWR enterprise artifact manually: %s", err)
	}

	d.SetId(instanceId + "/" + namespaceName + "/" + repositoryName + "/" + reference)

	return nil
}

func resourceSwrEnterpriseInstanceArtifactManualScanRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseInstanceArtifactManualScanUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseInstanceArtifactManualScanDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR enterprise artifact manual scan resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
