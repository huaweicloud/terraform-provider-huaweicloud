package eps

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API EPS POST /v2/enterprises/enterprise-projects/authority
func ResourceAuthority() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthorityCreate,
		ReadContext:   resourceAuthorityRead,
		DeleteContext: resourceAuthorityDelete,
	}
}

func resourceAuthorityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("bss", cfg.Region)
	if err != nil {
		return diag.Errorf("error creating BSS client: %s", err)
	}

	httpUrl := "v2/enterprises/enterprise-projects/authority"
	createPath := client.Endpoint + httpUrl

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("POST", createPath, &opts)
	if err != nil {
		return diag.Errorf("error granting the enterprise project: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceAuthorityRead(ctx, d, meta)
}

func resourceAuthorityRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAuthorityDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Unable to reclaim the enterprise project authorization of the current account during this resource delete.
The authority info can only be removed from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
