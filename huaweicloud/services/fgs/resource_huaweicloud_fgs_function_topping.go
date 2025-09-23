package fgs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

type ToppingStatus string

var (
	ToppingStatusEnable  ToppingStatus = "true"
	ToppingStatusDisable ToppingStatus = "false"
)

// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{func_urn}/collect/{state}
func ResourceFunctionTopping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionToppingCreate,
		ReadContext:   resourceFunctionToppingRead,
		DeleteContext: resourceFunctionToppingDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the function is located.`,
			},
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The URN of the function to be topped.`,
			},
		},
	}
}

func updateFunctionToppingStatus(client *golangsdk.ServiceClient, functionUrn, status string) error {
	httpUrl := "v2/{project_id}/fgs/functions/{func_urn}/collect/{state}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{func_urn}", functionUrn)
	updatePath = strings.ReplaceAll(updatePath, "{state}", status)

	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	_, err := client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return fmt.Errorf("error updating function topping status (target is: %s): %s", status, err)
	}
	return nil
}

func resourceFunctionToppingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		functionUrn = d.Get("function_urn").(string)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	err = updateFunctionToppingStatus(client, functionUrn, string(ToppingStatusEnable))
	if err != nil {
		return diag.FromErr(err)
	}

	// In the same script, a function can only be topped once.
	d.SetId(functionUrn)

	return resourceFunctionToppingRead(ctx, d, meta)
}

func resourceFunctionToppingRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFunctionToppingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		functionUrn = d.Get("function_urn").(string)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	err = updateFunctionToppingStatus(client, functionUrn, string(ToppingStatusDisable))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FunctionGraph function topping")
	}
	return nil
}
