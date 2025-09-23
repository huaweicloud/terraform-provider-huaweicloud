package dew

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

// @API DEW GET /v1.0/{project_id}/kms/aliases
func DataSourceKmsAliases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKmsAliasRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the key ID used to query the alias.`,
			},
			"aliases": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        aliasSchema(),
				Description: `The list of key aliases.`,
			},
		},
	}
}

func aliasSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the account to which the alias belongs.`,
			},
			"key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key ID`,
			},
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alias of the key`,
			},
			"alias_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alias resource locator.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the alias.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the alias.`,
			},
		},
	}
}

func buildListKmsAliasParams(d *schema.ResourceData) string {
	res := "?limit=50"

	if keyId, ok := d.GetOk("key_id"); ok {
		res = fmt.Sprintf("%s&key_id=%v", res, keyId)
	}

	return res
}

func listKmsAlias(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		getHttpUrl = "v1.0/{project_id}/kms/aliases"
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + getHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	queryParams := buildListKmsAliasParams(d)
	listPath += queryParams
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	marker := ""
	for {
		getListPath := listPath
		if marker != "" {
			getListPath = listPath + fmt.Sprintf("&marker=%s", marker)
		}

		requestResp, err := client.Request("GET", getListPath, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving KMS aliases: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		aliases := utils.PathSearch("aliases", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, aliases...)

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceKmsAliasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "kms"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	aliases, err := listKmsAlias(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(id)

	mErr = multierror.Append(mErr,
		d.Set("region", region),
		d.Set("aliases", flattenKmsAlias(aliases)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKmsAlias(aliases []interface{}) []interface{} {
	if len(aliases) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(aliases))
	for _, alias := range aliases {
		result = append(result, map[string]interface{}{
			"domain_id":   utils.PathSearch("domain_id", alias, nil),
			"key_id":      utils.PathSearch("key_id", alias, nil),
			"alias":       utils.PathSearch("alias", alias, nil),
			"alias_urn":   utils.PathSearch("alias_urn", alias, nil),
			"create_time": utils.PathSearch("create_time", alias, nil),
			"update_time": utils.PathSearch("update_time", alias, nil),
		})
	}

	return result
}
