// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM GET /v1/{project_id}/instances/{instance_id}/users
func DataSourceDdmAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmAccountsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of DDM instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the DDM account.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the DDM account.`,
			},
			"accounts": {
				Type:        schema.TypeList,
				Elem:        AccountsAccountSchema(),
				Computed:    true,
				Description: `Indicates the list of DDM account.`,
			},
		},
	}
}

func AccountsAccountSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the DDM account.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DDM account.`,
			},
			"permissions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the basic permissions of the DDM account.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the DDM account.`,
			},
			"schemas": {
				Type:        schema.TypeList,
				Elem:        AccountsAccountSchemaSchema(),
				Computed:    true,
				Description: `Indicates the schemas that associated with the account.`,
			},
		},
	}
	return &sc
}

func AccountsAccountSchemaSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the associated schema.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the schema description.`,
			},
		},
	}
	return &sc
}

func resourceDdmAccountsRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDdmAccounts: Query the List of DDM account
	var (
		getDdmAccountsHttpUrl = "v1/{project_id}/instances/{instance_id}/users"
		getDdmAccountsProduct = "ddm"
	)
	getDdmAccountsClient, err := cfg.NewServiceClient(getDdmAccountsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM Client: %s", err)
	}

	getDdmAccountsPath := getDdmAccountsClient.Endpoint + getDdmAccountsHttpUrl
	getDdmAccountsPath = strings.ReplaceAll(getDdmAccountsPath, "{project_id}", getDdmAccountsClient.ProjectID)
	getDdmAccountsPath = strings.ReplaceAll(getDdmAccountsPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))

	getDdmAccountsResp, err := pagination.ListAllItems(
		getDdmAccountsClient,
		"offset",
		getDdmAccountsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DDM account")
	}

	getDdmAccountsRespJson, err := json.Marshal(getDdmAccountsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDdmAccountsRespBody any
	err = json.Unmarshal(getDdmAccountsRespJson, &getDdmAccountsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("accounts", flattenGetAccountsResponseBodyAccount(d, getDdmAccountsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetAccountsResponseBodyAccount(d *schema.ResourceData, resp any) []any {
	if resp == nil {
		return nil
	}

	name := d.Get("name").(string)
	status := d.Get("status").(string)
	curJson := utils.PathSearch("users", resp, make([]any, 0))
	curArray := curJson.([]any)
	rst := make([]any, 0, len(curArray))
	for _, v := range curArray {
		accountName := utils.PathSearch("name", v, nil)
		accountStatus := utils.PathSearch("status", v, nil)
		if name != "" && name != accountName {
			continue
		}
		if status != "" && status != accountStatus {
			continue
		}
		rst = append(rst, map[string]any{
			"name":        accountName,
			"status":      accountStatus,
			"permissions": utils.PathSearch("base_authority", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"schemas":     flattenAccountSchemas(v),
		})
	}
	return rst
}

func flattenAccountSchemas(resp any) []any {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("databases", resp, make([]any, 0))
	curArray := curJson.([]any)
	rst := make([]any, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]any{
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}
	return rst
}
