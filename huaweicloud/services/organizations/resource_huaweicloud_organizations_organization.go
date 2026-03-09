package organizations

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations POST /v1/organizations
// @API Organizations GET /v1/organizations/roots
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/create
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/delete
// @API Organizations POST /v1/organizations/policies/enable
// @API Organizations GET /v1/organizations
// @API Organizations GET /v1/organizations/{resource_type}/{resource_id}/tags
// @API Organizations POST /v1/organizations/policies/disable
// @API Organizations DELETE /v1/organizations
func ResourceOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationCreate,
		UpdateContext: resourceOrganizationUpdate,
		ReadContext:   resourceOrganizationRead,
		DeleteContext: resourceOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"enabled_policy_types": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of policy types to be enabled in the organization root.`,
			},
			"root_tags": common.TagsSchema(`The key/value pairs to be attached to the root.`),
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the organization.`,
			},
			"master_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique ID of the organization's management account.`,
			},
			"master_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the organization's management account.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the organization was created.`,
			},
			"root_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the root.`,
			},
			"root_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the root.`,
			},
			"root_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The urn of the root.`,
			},
		},
	}
}

func resourceOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/organizations"
	)

	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating organization: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	organizationId := utils.PathSearch("organization.id", respBody, "").(string)
	if organizationId == "" {
		return diag.Errorf("unable to find the organization ID from the API response")
	}
	d.SetId(organizationId)

	getRootRespBody, err := getRoot(client)
	if err != nil {
		return diag.FromErr(err)
	}

	rootId := utils.PathSearch("roots|[0].id", getRootRespBody, "").(string)
	if v, ok := d.GetOk("root_tags"); ok {
		tagList := utils.ExpandResourceTags(v.(map[string]interface{}))
		err = addTags(client, rootType, rootId, tagList)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("enabled_policy_types"); ok {
		enabledPolicyTypes := v.(*schema.Set).List()
		for _, enabledPolicyType := range enabledPolicyTypes {
			if err = enablePolicy(ctx, d.Timeout(schema.TimeoutCreate), client,
				enabledPolicyType.(string), rootId); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceOrganizationRead(ctx, d, meta)
}

func resourceOrganizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	respBody, err := GetOrganization(client)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			"error retrieving organization",
		)
	}

	getRootRespBody, err := getRoot(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Organizations root")
	}

	rootId := utils.PathSearch("roots|[0].id", getRootRespBody, "").(string)
	policyTypes := utils.PathSearch("roots|[0].policy_types[?status=='enabled'].type", getRootRespBody,
		make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("urn", utils.PathSearch("organization.urn", respBody, nil)),
		d.Set("master_account_id", utils.PathSearch("organization.management_account_id", respBody, nil)),
		d.Set("master_account_name", utils.PathSearch("organization.management_account_name", respBody, nil)),
		d.Set("created_at", utils.PathSearch("organization.created_at", respBody, nil)),
		d.Set("root_id", rootId),
		d.Set("root_name", utils.PathSearch("roots|[0].name", getRootRespBody, nil)),
		d.Set("root_urn", utils.PathSearch("roots|[0].urn", getRootRespBody, nil)),
		d.Set("enabled_policy_types", policyTypes),
	)

	tagMap, err := getTags(client, rootType, rootId)
	if err != nil {
		log.Printf("[WARN] error fetching Organizations tags of root (%s): %s", rootId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("root_tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetOrganization(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v1/organizations"
	getPath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func enablePolicy(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient, policyType,
	rootId string) error {
	err := requestRootPolicy(client, policyType, rootId, "enable")
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      policyStateRefreshFunc(client, policyType, "enabled"),
		Timeout:      timeout,
		Delay:        1 * time.Second,
		PollInterval: 1 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for Organizations policy type (%s) to be enabled: %s", policyType, err)
	}
	return nil
}

func disablePolicy(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient, policyType,
	rootId string) error {
	err := requestRootPolicy(client, policyType, rootId, "disable")
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      policyStateRefreshFunc(client, policyType, "disabled"),
		Timeout:      timeout,
		Delay:        1 * time.Second,
		PollInterval: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for Organizations policy type (%s) to be disabled: %s", policyType, err)
	}
	return nil
}

func requestRootPolicy(client *golangsdk.ServiceClient, policyType, rootId, action string) error {
	httpUrl := fmt.Sprintf("v1/organizations/policies/%s", action)
	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildRequestRootPolicyBodyParams(policyType, rootId)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func buildRequestRootPolicyBodyParams(policyType, rootId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_type": policyType,
		"root_id":     rootId,
	}
	return bodyParams
}

func policyStateRefreshFunc(client *golangsdk.ServiceClient, policyType, target string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getRoot(client)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch(fmt.Sprintf("roots|[0].policy_types[?type=='%s'].status|[0]", policyType), respBody, "").(string)
		if err != nil {
			return nil, "", err
		}

		if status == target {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func resourceOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	getRootRespBody, err := getRoot(client)
	if err != nil {
		return diag.FromErr(err)
	}

	rootId := utils.PathSearch("roots|[0].id", getRootRespBody, "").(string)
	if d.HasChange("root_tags") {
		if err = updateTags(d, client, rootType, rootId, "root_tags"); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enabled_policy_types") {
		oldRaw, newRaw := d.GetChange("enabled_policy_types")
		enabledPolicyTypes := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
		disabledPolicyTypes := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
		timeout := d.Timeout(schema.TimeoutUpdate)
		for _, enabledPolicyType := range enabledPolicyTypes.List() {
			if err = enablePolicy(ctx, timeout, client, enabledPolicyType.(string), rootId); err != nil {
				return diag.FromErr(err)
			}
		}

		for _, disabledPolicyType := range disabledPolicyTypes.List() {
			if err = disablePolicy(ctx, timeout, client, disabledPolicyType.(string), rootId); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceOrganizationRead(ctx, d, meta)
}

func resourceOrganizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/organizations"
	)

	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			"error deleting organization",
		)
	}

	return nil
}
