package dds

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS PUT /v3/{project_id}/instances/recycle-policy
// @API DDS GET /v3/{project_id}/instances/recycle-policy
func ResourceDDSRecyclePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSRecyclePolicyCreateOrUpdate,
		ReadContext:   resourceDDSRecyclePolicyRead,
		UpdateContext: resourceDDSRecyclePolicyCreateOrUpdate,
		DeleteContext: resourceDDSRecyclePolicyDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"retention_period_in_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceDDSRecyclePolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	err = updateRecyclePolicy(client, d.Get("retention_period_in_days").(int))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.IsNewResource() {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return diag.Errorf("unable to generate ID: %s", err)
		}
		d.SetId(uuid)
	}

	return resourceDDSRecyclePolicyRead(ctx, d, meta)
}

func updateRecyclePolicy(client *golangsdk.ServiceClient, retentionPeriod int) error {
	updateHttpUrl := "v3/{project_id}/instances/recycle-policy"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"recycle_policy": map[string]interface{}{
				"enabled":                  true,
				"retention_period_in_days": retentionPeriod,
			},
		},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating DDS recycle policy: %s", err)
	}

	return nil
}

func resourceDDSRecyclePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getHttpUrl := "v3/{project_id}/instances/recycle-policy"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DDS recycle policy")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("retention_period_in_days", utils.PathSearch("recycle_policy.retention_period_in_days", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDDSRecyclePolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	err = updateRecyclePolicy(client, 7)
	if err != nil {
		return diag.FromErr(err)
	}

	errorMsg := "Deleting recycle policy is unsupported. The resource is removed from the state, and the retention " +
		"period is reset to 7."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
