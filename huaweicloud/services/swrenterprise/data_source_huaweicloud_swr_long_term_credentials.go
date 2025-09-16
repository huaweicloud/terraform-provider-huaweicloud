package swrenterprise

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

// @API SWR GET /v2/{project_id}/instances/{instance_id}/long-term-credentials
func DataSourceSwrEnterpriseLongTermCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseLongTermCredentialsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"auth_tokens": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the namespaces.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the namespace ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the credential name.`,
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to enable the credential.`,
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user ID.`,
						},
						"user_profile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user profile.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time.`,
						},
						"expire_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the expired time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSwrEnterpriseLongTermCredentialsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	listLongTermCredentialsHttpUrl := "v2/{project_id}/instances/{instance_id}/long-term-credentials"
	listLongTermCredentialsPath := client.Endpoint + listLongTermCredentialsHttpUrl
	listLongTermCredentialsPath = strings.ReplaceAll(listLongTermCredentialsPath, "{project_id}", client.ProjectID)
	listLongTermCredentialsPath = strings.ReplaceAll(listLongTermCredentialsPath, "{instance_id}", d.Get("instance_id").(string))
	listLongTermCredentialsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listLongTermCredentialsPath + fmt.Sprintf("?limit=100&offset=%v", offset)
		listLongTermCredentialsResp, err := client.Request("GET", currentPath, &listLongTermCredentialsOpt)
		if err != nil {
			return diag.Errorf("error querying SWR long term credentials: %s", err)
		}
		listLongTermCredentialsRespBody, err := utils.FlattenResponse(listLongTermCredentialsResp)
		if err != nil {
			return diag.Errorf("error flattening SWR long term credentials response: %s", err)
		}

		tokens := utils.PathSearch("auth_tokens", listLongTermCredentialsRespBody, make([]interface{}, 0)).([]interface{})
		if len(tokens) == 0 {
			break
		}
		for _, token := range tokens {
			results = append(results, map[string]interface{}{
				"id":           utils.PathSearch("token_id", token, nil),
				"name":         utils.PathSearch("name", token, nil),
				"enable":       utils.PathSearch("enable", token, nil),
				"user_id":      utils.PathSearch("user_id", token, nil),
				"user_profile": utils.PathSearch("user_profile", token, nil),
				"created_at":   utils.PathSearch("created_at", token, nil),
				"expire_date":  utils.PathSearch("expire_date", token, nil),
			})
		}

		// offset must be the multiple of limit
		offset += 100
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("auth_tokens", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
