package fgs

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph POST /v2/{project_id}/fgs/functions/enable-lts-logs
func ResourceLtsLogEnable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsLogEnableCreate,
		ReadContext:   resourceLtsLogEnableRead,
		DeleteContext: resourceLtsLogEnableDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the LTS log function is to be enabled.`,
			},
		},
	}
}

func ignoreErrorWhileLtsLogIsEnabled(err error) error {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return err
		}

		errorMsg := utils.PathSearch("error_msg", apiError, "").(string)
		if regexp.MustCompile(`lts log has been enabled`).MatchString(strings.ToLower(errorMsg)) {
			return nil
		}
	}
	return err
}

func resourceLtsLogEnableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/fgs/functions/enable-lts-logs"
	)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("POST", createPath, &opt)
	if paredErr := ignoreErrorWhileLtsLogIsEnabled(err); paredErr != nil {
		return diag.Errorf("error enabling LTS logs for FunctionGraph: %s", paredErr)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return resourceLtsLogEnableRead(ctx, d, meta)
}

func resourceLtsLogEnableRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLtsLogEnableDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for enabling LTS logs for FunctionGraph. Deleting this resource will
not disable the LTS logs, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
