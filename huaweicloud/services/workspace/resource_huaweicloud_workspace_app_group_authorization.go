package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/app-groups/authorizations
// @API Workspace GET /v1/{project_id}/app-groups/actions/list-authorizations
// @API Workspace POST /v1/{project_id}/app-groups/actions/batch-delete-authorization
func ResourceAppGroupAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppGroupAuthorizationCreate,
		ReadContext:   resourceAppGroupAuthorizationRead,
		DeleteContext: resourceAppGroupAuthorizationDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"app_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the application group.`,
			},
			"accounts": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: `The list of the accounts to be authorized.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The ID of the user (group).",
						},
						"account": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of the user (group).",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The type of the object to be authorized.",
						},
					},
				},
			},
		},
	}
}

func resourceAppGroupAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-groups/authorizations"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAppAuthorizationBodyParams(d)),
	}
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error authorizing accounts to apllication group: %s", err)
	}

	// When `accounts` contains non-existent objects，the API respond status code is `200`,
	// so call the query interface to determine whether the authorization is successful.
	authorizedAccounts, err := getAppAuthorizations(client, d.Get("app_group_id").(string))
	if err != nil {
		return diag.Errorf("error retrieving authorized accounts: %s", err)
	}

	accounts := d.Get("accounts").([]interface{})
	failedAccounts := getFailedAuthAccounts(accounts, authorizedAccounts)
	if len(failedAccounts) != 0 {
		return diag.Errorf("unable to authorize for some accounts: %s", strings.Join(failedAccounts, "\n"))
	}

	d.SetId(d.Get("app_group_id").(string))

	return resourceAppGroupAuthorizationRead(ctx, d, meta)
}

func getFailedAuthAccounts(inputObjects, resAuthorizedAccounts []interface{}) []string {
	failedObjects := make([]string, 0)
	for _, v := range inputObjects {
		account := utils.PathSearch("account", v, "").(string)
		accountType := utils.PathSearch("type", v, "").(string)
		isExisit := false
		for _, item := range resAuthorizedAccounts {
			resObject := utils.PathSearch("account", item, "").(string)
			resObjectType := utils.PathSearch("account_type", item, "").(string)
			if account == resObject && accountType == resObjectType {
				isExisit = true
				break
			}
		}

		if !isExisit {
			failedObjects = append(failedObjects, fmt.Sprintf("%s | %s;", account, accountType))
		}
	}

	return failedObjects
}

func buildAppAuthorizationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"app_group_ids": []string{d.Get("app_group_id").(string)},
		"accounts":      buildAccounts(d.Get("accounts").([]interface{})),
	}
}

func buildAccounts(accounts []interface{}) []map[string]interface{} {
	rest := make([]map[string]interface{}, len(accounts))
	for i, v := range accounts {
		rest[i] = map[string]interface{}{
			"id":           utils.ValueIgnoreEmpty(utils.PathSearch("id", v, "")),
			"account":      utils.PathSearch("account", v, ""),
			"account_type": utils.PathSearch("type", v, ""),
		}
	}
	return rest
}

func resourceAppGroupAuthorizationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func getAppAuthorizations(client *golangsdk.ServiceClient, appGroupId string) ([]interface{}, error) {
	var (
		// The dafault value of limit is `10`.
		httpUrl            = "v1/{project_id}/app-groups/actions/list-authorizations?limit=100"
		offset             = 0
		authorizedAccounts = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += fmt.Sprintf("&app_group_id=%v", appGroupId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &getOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		authorizations := utils.PathSearch("authorizations", respBody, make([]interface{}, 0)).([]interface{})
		authorizedAccounts = append(authorizedAccounts, authorizations...)
		if len(authorizedAccounts) == int(utils.PathSearch("count", respBody, float64(0)).(float64)) {
			break
		}
		offset += len(authorizedAccounts)
	}
	return authorizedAccounts, nil
}

func resourceAppGroupAuthorizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-groups/actions/batch-delete-authorization"
	)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAppAuthorizationBodyParams(d)),
	}

	// In any case, the delete interface response status is `200`.
	_, err = client.Request("POST", deletePath, &deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting APP group authorization: %s", err)
	}
	return nil
}
