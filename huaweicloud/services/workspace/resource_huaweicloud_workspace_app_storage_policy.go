package workspace

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace PUT /v1/{project_id}/storages-policy/actions/create-statements
// @API Workspace GET /v1/{project_id}/storages-policy/actions/list-statements
func ResourceAppStoragePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppStoragePolicyCreate,
		ReadContext:   resourceAppStoragePolicyRead,
		UpdateContext: resourceAppStoragePolicyUpdate,
		DeleteContext: resourceAppStoragePolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppStoragePolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the custom storage permission policy is located.",
			},
			"server_actions": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The collection of permissions that server can use to access storage.",
			},
			"client_actions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The collection of permissions that client can use to access storage.",
			},
		},
	}
}

func buildAppStoragePolicyCreateOpts(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"roam_actions": d.Get("server_actions").(*schema.Set).List(),
		"actions":      d.Get("client_actions").(*schema.Set).List(),
	}
}

func updateAppStoragePolicy(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	httpUrl := "v1/{project_id}/storages-policy/actions/create-statements"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAppStoragePolicyCreateOpts(d)),
	}

	requestResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return "", err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return "", err
	}

	policyId := utils.PathSearch("policy_statement_id", respBody, "").(string)
	if policyId == "" {
		return "", fmt.Errorf("unable to find the custom storage permission policy ID from the API response")
	}
	return policyId, nil
}

func resourceAppStoragePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	policyId, err := updateAppStoragePolicy(client, d)
	if err != nil {
		return diag.Errorf("error creating custom storage permission policy of Workspace APP: %s", err)
	}
	d.SetId(policyId)

	return resourceAppStoragePolicyRead(ctx, d, meta)
}

func listAppStoragePermissionPolicies(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/storages-policy/actions/list-statements?limit=100"
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		offset = 0
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, fmt.Errorf("error getting list of storage permission policies: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(items) < 1 {
			break
		}
		result = append(result, items...)
		offset += len(items)
	}
	return result, nil
}

func GetAppCustomStoragePolicy(client *golangsdk.ServiceClient) (interface{}, error) {
	policies, err := listAppStoragePermissionPolicies(client)
	if err != nil {
		return nil, err
	}
	policy := utils.PathSearch("[?!(contains(policy_statement_id, 'DEFAULT'))]|[0]", policies, nil)
	if policy == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("the custom storage permission policy has been removed from the Workspace APP service"),
			},
		}
	}
	return policy, nil
}

func resourceAppStoragePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	policy, err := GetAppCustomStoragePolicy(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Custom storage permission policy of Workspace APP")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("server_actions", utils.PathSearch("roam_actions", policy, nil)),
		d.Set("client_actions", utils.PathSearch("actions", policy, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("unable to setting resource fields of the custom storage permission policy: %s", err)
	}
	return nil
}

func resourceAppStoragePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	_, err = updateAppStoragePolicy(client, d)
	if err != nil {
		return diag.Errorf("error updating custom storage permission policy of Workspace APP: %s", err)
	}

	return resourceAppStoragePolicyRead(ctx, d, meta)
}

func resourceAppStoragePolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Deleting this resource will not initialize (restore) the policy configuration and just only remove the
tfstate record for this resource.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceAppStoragePolicyImportState(_ context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	var (
		re       = regexp.MustCompile(`^\d{18}$`)
		importId = d.Id()
		policyId string
	)
	if re.MatchString(importId) {
		// The imported ID is custom storage permission ID.
		policyId = importId
	} else {
		var (
			cfg    = meta.(*config.Config)
			region = cfg.GetRegion(d)
		)
		client, err := cfg.NewServiceClient("appstream", region)
		if err != nil {
			return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
		}

		// Refresh the resource ID using field value of attribute 'policy_statement_id'.
		policy, err := GetAppCustomStoragePolicy(client)
		if err != nil {
			return nil, err
		}
		policyId = utils.PathSearch("policy_statement_id", policy, "").(string)
		if policyId == "" {
			return nil, fmt.Errorf("unable to find the custom storage permission policy ID from the API response")
		}
	}

	d.SetId(policyId)
	return []*schema.ResourceData{d}, nil
}
