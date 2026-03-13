package cse

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSE GET /v4/{project_id}/registry/microservices
func DataSourceMicroservices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMicroservicesRead,

		Schema: map[string]*schema.Schema{
			// Special parameters.
			// These parameters are used to specify the address that used to request the access token and access the
			// microservice engine.
			"auth_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The address that used to request the access token.`,
			},
			"connect_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The address that used to send requests and manage configuration.`,
			},
			"admin_user": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"admin_pass"},
				Description:  `The user name that used to pass the RBAC control.`,
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"admin_user"},
				Description:  `The user password that used to pass the RBAC control.`,
			},

			// Attributes.
			"microservices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the microservice.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the microservice.`,
						},
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application ID of the microservice.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the microservice.`,
						},
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The environment type of the microservice.`,
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The level of the microservice.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the microservice.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the microservice.`,
						},
						"framework": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the framework.`,
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The version of the framework.`,
									},
								},
							},
							Description: `The framework information of the microservice.`,
						},
						"schemas": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of schemas supported by the microservice.`,
						},
						"paths": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The path of the microservice route.`,
									},
									"properties": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The properties of the microservice route.`,
									},
								},
							},
							Description: `The list of paths exposed by the microservice.`,
						},
						"properties": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The extended attributes of the microservice, in key/value format.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the microservice, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the microservice, in RFC3339 format.`,
						},
					},
				},
				Description: `The list of all microservices.`,
			},
		},
	}
}

func listMicroservices(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v4/{project_id}/registry/microservices"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", microserviceDefaultProjectId)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(client.ProjectID),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `POST /v4/token` interface.
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string), d.Get("admin_pass").(string))
	if err != nil {
		return nil, err
	}
	// If the microservice instance has RBAC authentication enabled, the Authorization header will use a special token
	// provided by the CSE service to replace the original IAM authentication information (AKSK authentication) in the
	// request header.
	if token != "" {
		listOpts.MoreHeaders["Authorization"] = token
	}

	requestResp, err := client.Request("GET", listPath, &listOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("services", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenMicroservicesFramework(framework map[string]interface{}) []map[string]interface{} {
	if len(framework) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":    utils.PathSearch("name", framework, nil),
			"version": utils.PathSearch("version", framework, nil),
		},
	}
}

func flattenMicroservicesPaths(paths []interface{}) []map[string]interface{} {
	if len(paths) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(paths))
	for _, path := range paths {
		result = append(result, map[string]interface{}{
			"path":       utils.PathSearch("Path", path, nil),
			"properties": utils.JsonToString(utils.PathSearch("properties", path, make(map[string]interface{})).(map[string]interface{})),
		})
	}
	return result
}

func flattenMicroservices(microservices []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(microservices))

	for _, microservice := range microservices {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("serviceId", microservice, nil),
			"name":        utils.PathSearch("serviceName", microservice, nil),
			"app_id":      utils.PathSearch("appId", microservice, nil),
			"version":     utils.PathSearch("version", microservice, nil),
			"environment": utils.PathSearch("environment", microservice, nil),
			"level":       utils.PathSearch("level", microservice, nil),
			"description": utils.PathSearch("description", microservice, nil),
			"status":      utils.PathSearch("status", microservice, nil),
			"framework": flattenMicroservicesFramework(utils.PathSearch("framework",
				microservice, make(map[string]interface{})).(map[string]interface{})),
			"schemas":    utils.PathSearch("schemas", microservice, make([]interface{}, 0)).([]interface{}),
			"paths":      flattenMicroservicesPaths(utils.PathSearch("paths", microservice, make([]interface{}, 0)).([]interface{})),
			"properties": utils.PathSearch("properties", microservice, make(map[string]interface{})).(map[string]interface{}),
			"created_at": parseMicroservicesTime(utils.PathSearch("timestamp", microservice, "").(string)),
			"updated_at": parseMicroservicesTime(utils.PathSearch("modTimestamp", microservice, "").(string)),
		})
	}

	return result
}

func parseMicroservicesTime(timeStampStr string) string {
	if timeStampStr == "" {
		return ""
	}

	r, err := strconv.Atoi(timeStampStr)
	if err != nil {
		log.Printf("[ERROR] unable to convert the string (%s) to int", timeStampStr)
		return ""
	}

	return utils.FormatTimeStampRFC3339(int64(r), false)
}

func dataSourceMicroservicesRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := common.NewCustomClient(true, d.Get("connect_address").(string), "v2", "default")

	microservices, err := listMicroservices(client, d)
	if err != nil {
		return diag.Errorf("error querying CSE microservices: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("microservices", flattenMicroservices(microservices)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
