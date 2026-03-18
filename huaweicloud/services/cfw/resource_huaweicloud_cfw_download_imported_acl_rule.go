package cfw

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW GET /v1/{project_id}/acl-rule/import-result
func ResourceDownloadImportedAclRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDownloadImportedAclRuleCreate,
		ReadContext:   resourceDownloadImportedAclRuleRead,
		UpdateContext: resourceDownloadImportedAclRuleUpdate,
		DeleteContext: resourceDownloadImportedAclRuleDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"object_id",
			"export_file_name",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export_file_name": {
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

func resourceDownloadImportedAclRuleCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		httpUrl        = "v1/{project_id}/acl-rule/import-result"
		objectId       = d.Get("object_id").(string)
		exportFileName = d.Get("export_file_name").(string)
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?object_id=%s", objectId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error downloading CFW imported ACL rule: %s", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error reading response body: %s", err)
	}

	if err := os.WriteFile(exportFileName, bodyBytes, 0600); err != nil {
		return diag.Errorf("failed to write file to (%s): %s", exportFileName, err)
	}

	d.SetId(objectId)
	return nil
}

func resourceDownloadImportedAclRuleRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDownloadImportedAclRuleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDownloadImportedAclRuleDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to download imported acl rule. Deleting this resource
    will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
