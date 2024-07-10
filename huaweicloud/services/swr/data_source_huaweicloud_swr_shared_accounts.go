package swr

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/manage/namespaces/{namespace}/repositories/{repository}/access-domains
func DataSourceSharedAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharedAccountsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"shared_accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shared_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"permit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deadline": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSharedAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	// Query the list of SWR shared accounts.
	resp, err := getSharedAccounts(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("shared_accounts", flattenSharedAccountsResponse(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getSharedAccounts(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)

	listSharedAccountsHttpUrl := "v2/manage/namespaces/{namespace}/repositories/{repository}/access-domains"
	listSharedAccountsPath := client.Endpoint + listSharedAccountsHttpUrl
	listSharedAccountsPath = strings.ReplaceAll(listSharedAccountsPath, "{namespace}", organization)
	listSharedAccountsPath = strings.ReplaceAll(listSharedAccountsPath, "{repository}", repository)

	listSharedAccountsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listSharedAccountsResp, err := client.Request("GET", listSharedAccountsPath, &listSharedAccountsOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying SWR shared accounts: %s", err)
	}
	listSharedAccountsRespBody, err := utils.FlattenResponse(listSharedAccountsResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening SWR shared accounts: %s", err)
	}
	sharedAccounts := listSharedAccountsRespBody.([]interface{})

	return sharedAccounts, nil
}

func flattenSharedAccountsResponse(rawParams []interface{}) []interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	results := make([]interface{}, len(rawParams))
	for i, v := range rawParams {
		createdAt := utils.PathSearch("created", v, nil).(string)
		updatedAt := utils.PathSearch("updated", v, nil).(string)
		results[i] = map[string]interface{}{
			"organization":   utils.PathSearch("namespace", v, nil),
			"repository":     utils.PathSearch("repository", v, nil),
			"shared_account": utils.PathSearch("access_domain", v, nil),
			"status":         utils.PathSearch("status", v, false),
			"permit":         utils.PathSearch("permit", v, nil),
			"deadline":       utils.PathSearch("deadline", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"creator_id":     utils.PathSearch("creator_id", v, nil),
			"created_by":     utils.PathSearch("creator_name", v, nil),
			"created_at":     utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(createdAt)/1000, false),
			"updated_at":     utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(updatedAt)/1000, false),
		}
	}
	return results
}
