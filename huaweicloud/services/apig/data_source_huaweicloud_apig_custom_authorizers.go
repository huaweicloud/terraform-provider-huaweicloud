package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/authorizers
func DataSourceCustomAuthorizers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomAuthorizersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the dedicated instance to which the custom authorizers belong.`,
			},
			"authorizer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the custom authorizer.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the custom authorizer.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the custom authorizer.`,
			},
			"authorizers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All custom authorizers that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the custom authorizer.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the custom authorizer.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the custom authorizer.`,
						},
						"function_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the FGS function.`,
						},
						"function_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URN of the FGS function.`,
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The network architecture types of function.`,
						},
						"function_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the FGS function.`,
						},
						"function_alias_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version alias URI of the FGS function.`,
						},
						"identity": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The parameter identities of the custom authorizer.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the parameter to be verified.`,
									},
									"location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The parameter location of identity.`,
									},
									"validation": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The parameter verification expression of identity.`,
									},
								},
							},
						},
						"cache_age": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum cache age of custom authorizer.`,
						},
						"user_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user data of custom authorizer.`,
						},
						"is_body_send": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to send the body of custom authorizer.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of custom authorizer.`,
						},
					},
				},
			},
		},
	}
}

func buildListCustomAuthorizersParams(d *schema.ResourceData) string {
	res := ""
	if authorizer, ok := d.GetOk("authorizer_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, authorizer)
	}
	if name, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, name)
	}
	if customType, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, customType)
	}
	return res
}

func queryCustomAuthorizers(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/authorizers?limit=500"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	queryParams := buildListCustomAuthorizersParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving custom authorizers under specified "+
				"dedicated instance (%s): %s", instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		authorizers := utils.PathSearch("authorizer_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(authorizers) < 1 {
			break
		}
		result = append(result, authorizers...)
		offset += len(authorizers)
	}
	return result, nil
}

func dataSourceCustomAuthorizersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	authorizers, err := queryCustomAuthorizers(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("authorizers", flattenCustomAuthorizers(authorizers)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCustomAuthorizers(authorizers []interface{}) []interface{} {
	if len(authorizers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(authorizers))
	for _, authorizer := range authorizers {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", authorizer, nil),
			"name":               utils.PathSearch("name", authorizer, nil),
			"type":               utils.PathSearch("type", authorizer, nil),
			"function_type":      utils.PathSearch("authorizer_type", authorizer, nil),
			"function_urn":       utils.PathSearch("authorizer_uri", authorizer, nil),
			"network_type":       utils.PathSearch("network_type", authorizer, nil),
			"function_version":   utils.PathSearch("authorizer_version", authorizer, nil),
			"function_alias_uri": utils.PathSearch("authorizeralias_uri", authorizer, nil),
			"identity":           flattenCustomAuthorizersIdentity(utils.PathSearch("identities", authorizer, make([]interface{}, 0))),
			"cache_age":          int(utils.PathSearch("ttl", authorizer, float64(0)).(float64)),
			"user_data":          utils.PathSearch("user_data", authorizer, nil),
			"is_body_send":       utils.PathSearch("need_body", authorizer, false).(bool),
			"created_at":         utils.PathSearch("create_time", authorizer, nil),
		})
	}
	return result
}

func flattenCustomAuthorizersIdentity(identities interface{}) []map[string]interface{} {
	identitiesArray := identities.([]interface{})
	result := make([]map[string]interface{}, len(identitiesArray))
	for i, identity := range identitiesArray {
		result[i] = map[string]interface{}{
			"name":       utils.PathSearch("name", identity, nil),
			"location":   utils.PathSearch("location", identity, nil),
			"validation": utils.PathSearch("validation", identity, nil),
		}
	}

	return result
}
