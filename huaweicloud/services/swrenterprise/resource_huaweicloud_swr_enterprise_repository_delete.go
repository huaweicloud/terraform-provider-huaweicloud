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

var enterpriseRepositoryDeleteNonUpdatableParams = []string{
	"instance_id", "namespace_name", "repository_name",
}

// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}
func ResourceSwrEnterpriseRepositoryDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseRepositoryDeleteCreate,
		UpdateContext: resourceSwrEnterpriseRepositoryDeleteUpdate,
		ReadContext:   resourceSwrEnterpriseRepositoryDeleteRead,
		DeleteContext: resourceSwrEnterpriseRepositoryDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(enterpriseRepositoryDeleteNonUpdatableParams),

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceSwrEnterpriseRepositoryDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	repositoryName := d.Get("repository_name").(string)

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{namespace_name}", namespaceName)
	deletePath = strings.ReplaceAll(deletePath, "{repository_name}", url.PathEscape(strings.ReplaceAll(repositoryName, "/", "%2F")))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting SWR enterprise repository: %s", err)
	}

	d.SetId(instanceId + "/" + namespaceName + "/" + repositoryName)

	return nil
}

func resourceSwrEnterpriseRepositoryDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseRepositoryDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseRepositoryDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR enterprise repository delete resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
