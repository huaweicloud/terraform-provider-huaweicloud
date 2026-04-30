package dds

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS POST /v3/{project_id}/backups/download-policy
// @API DDS GET /v3/{project_id}/backups/download-policy
// @API DDS PUT /v3/{project_id}/backups/download-policy
func ResourceBackupDownloadPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupDownloadPolicyCreate,
		ReadContext:   resourceBackupDownloadPolicyRead,
		UpdateContext: resourceBackupDownloadPolicyUpdate,
		DeleteContext: resourceBackupDownloadPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildCreateBackupDownloadPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action": d.Get("action"),
	}

	return bodyParams
}

func resourceBackupDownloadPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	policyInfo, err := GetBackupDownloadPolicyInfo(client)
	if err != nil {
		return diag.Errorf("error querying backup download policy information: %s", err)
	}

	if policyInfo == nil || len(policyInfo.(map[string]interface{})) == 0 {
		err = settingBackupDownloadPolicy(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		respAction := utils.PathSearch("action", policyInfo, "").(string)
		policyId := utils.PathSearch("id", policyInfo, "").(string)
		targetAction := d.Get("action").(string)
		if respAction != targetAction {
			if policyId == "" {
				return diag.Errorf("error updating backup download policy: unable to find policy ID from the API response")
			}

			err = updateBackupDownloadPolicy(client, d, policyId)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(client.ProjectID)

	return resourceBackupDownloadPolicyRead(ctx, d, meta)
}

func settingBackupDownloadPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/backups/download-policy"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         buildCreateBackupDownloadPolicyBodyParams(d),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error creating backup download policy: %s", err)
	}

	return nil
}

func GetBackupDownloadPolicyInfo(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v3/{project_id}/backups/download-policy"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceBackupDownloadPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	policy, err := GetBackupDownloadPolicyInfo(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the backup download policy")
	}

	if policy == nil || len(policy.(map[string]interface{})) == 0 {
		// When the backup download policy does not exist, the response HTTP status code of the query API is 200
		// and return an empty object.
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving the backup download policy")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("action", utils.PathSearch("action", policy, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBackupDownloadPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	if d.HasChanges("action") {
		policyInfo, err := GetBackupDownloadPolicyInfo(client)
		if err != nil {
			return diag.Errorf("error querying backup download policy information: %s", err)
		}

		policyId := utils.PathSearch("id", policyInfo, "").(string)
		if policyId == "" {
			return diag.Errorf("error updating backup download policy: unable to find policy ID from the API response")
		}

		err = updateBackupDownloadPolicy(client, d, policyId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceBackupDownloadPolicyRead(ctx, d, meta)
}

func buildUpdateBackupDownloadPolicyBodyParams(d *schema.ResourceData, policyId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":     policyId,
		"action": d.Get("action"),
	}

	return bodyParams
}

func updateBackupDownloadPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, policyId string) error {
	httpUrl := "v3/{project_id}/backups/download-policy"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         buildUpdateBackupDownloadPolicyBodyParams(d, policyId),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating the backup download policy: %s", err)
	}

	return nil
}

func resourceBackupDownloadPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DDS backup download policy resource is not supported. The resource is only removed from the " +
		"state, the resource remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
