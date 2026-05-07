package drs

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

// @API DRS GET /v5/{project_id}/connections
func DataSourceDrsConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inst_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fetch_all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"driver_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"endpoint": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"endpoint_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_user": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_password": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_sharding": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"endpoint_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"db_port": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"db_user": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"db_password": {
													Type:      schema.TypeString,
													Computed:  true,
													Sensitive: true,
												},
												"instance_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"instance_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"db_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"vpc": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ssl": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ssl_link": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"ssl_cert_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssl_cert_key": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
									"ssl_cert_check_sum": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssl_cert_password": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildDrsConnectionsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("connection_id"); ok {
		queryParams += fmt.Sprintf("&connection_id=%s", v.(string))
	}
	if v, ok := d.GetOk("db_type"); ok {
		queryParams += fmt.Sprintf("&db_type=%s", v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams += fmt.Sprintf("&name=%s", v.(string))
	}
	if v, ok := d.GetOk("inst_id"); ok {
		queryParams += fmt.Sprintf("&inst_id=%s", v.(string))
	}
	if v, ok := d.GetOk("ip"); ok {
		queryParams += fmt.Sprintf("&ip=%s", v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		queryParams += fmt.Sprintf("&description=%s", v.(string))
	}
	if v, ok := d.GetOk("create_time"); ok {
		queryParams += fmt.Sprintf("&create_time=%s", v.(string))
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		queryParams += fmt.Sprintf("&enterprise_project_id=%s", v.(string))
	}
	if v, ok := d.GetOk("fetch_all"); ok {
		queryParams += fmt.Sprintf("&fetch_all=%t", v.(bool))
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams += fmt.Sprintf("&sort_key=%s", v.(string))
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams += fmt.Sprintf("&sort_dir=%s", v.(string))
	}

	return queryParams
}

func dataSourceDrsConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/connections"
		result  = make([]interface{}, 0)
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	rawQuery := buildDrsConnectionsQueryParams(d)

	rawQuery = strings.ReplaceAll(rawQuery, "&fetch_all=true", "")
	rawQuery = strings.ReplaceAll(rawQuery, "&fetch_all=false", "")
	rawQuery = strings.TrimPrefix(rawQuery, "&")

	fetchAll := d.Get("fetch_all").(bool)

	if fetchAll {
		rawQuery = removeQueryParam(rawQuery, "offset")
		rawQuery = removeQueryParam(rawQuery, "limit")
	}

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	limit := 10
	if fetchAll {
		limit = 100
	}
	offset := 0

	for {
		var currentPath string

		if rawQuery == "" {
			currentPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
		} else {
			currentPath = fmt.Sprintf("%s?%s&limit=%d", listPath, rawQuery, limit)
		}

		if fetchAll {
			currentPath = fmt.Sprintf("%s&offset=%d", currentPath, offset)
		}

		resp, err := client.Request("GET", currentPath, &reqOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		connections := utils.PathSearch("connections", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, connections...)

		if len(connections) < limit {
			break
		}

		if !fetchAll {
			break
		}
		offset += len(connections)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("connections", flattenDrsConnections(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func removeQueryParam(rawQuery, param string) string {
	pairs := strings.Split(rawQuery, "&")
	var newPairs []string
	for _, p := range pairs {
		if p == "" {
			continue
		}

		key := p
		if idx := strings.IndexByte(p, '='); idx >= 0 {
			key = p[:idx]
		}
		if key == param {
			continue
		}
		newPairs = append(newPairs, p)
	}
	return strings.Join(newPairs, "&")
}

func flattenDrsConnections(respBody interface{}) []interface{} {
	connections, ok := respBody.([]interface{})
	if !ok || len(connections) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(connections))
	for _, conn := range connections {
		result = append(result, map[string]interface{}{
			"connection_id":         utils.PathSearch("connection_id", conn, nil),
			"name":                  utils.PathSearch("name", conn, nil),
			"create_time":           utils.PathSearch("create_time", conn, 0),
			"db_type":               utils.PathSearch("db_type", conn, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", conn, nil),
			"description":           utils.PathSearch("description", conn, nil),
			"config":                flattenDrsConnectionConfig(conn),
			"endpoint":              flattenDrsEndpoint(conn),
			"vpc":                   flattenDrsVpcInfo(conn),
			"ssl":                   flattenDrsSslConfig(conn),
		})
	}
	return result
}

func flattenDrsConnectionConfig(resp interface{}) []interface{} {
	cfg := utils.PathSearch("config", resp, nil)
	if cfg == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"driver_name": utils.PathSearch("driver_name", cfg, nil),
		},
	}
}

func flattenDrsEndpoint(resp interface{}) []interface{} {
	endpoint := utils.PathSearch("endpoint", resp, nil)
	if endpoint == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":              utils.PathSearch("id", endpoint, nil),
			"endpoint_name":   utils.PathSearch("endpoint_name", endpoint, nil),
			"ip":              utils.PathSearch("ip", endpoint, nil),
			"db_port":         utils.PathSearch("db_port", endpoint, nil),
			"db_user":         utils.PathSearch("db_user", endpoint, nil),
			"db_password":     utils.PathSearch("db_password", endpoint, nil),
			"instance_id":     utils.PathSearch("instance_id", endpoint, nil),
			"instance_name":   utils.PathSearch("instance_name", endpoint, nil),
			"db_name":         utils.PathSearch("db_name", endpoint, nil),
			"source_sharding": flattenDrsSourceSharding(endpoint),
		},
	}
}

func flattenDrsSourceSharding(resp interface{}) []interface{} {
	shardingList := utils.PathSearch("source_sharding", resp, make([]interface{}, 0)).([]interface{})
	if len(shardingList) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(shardingList))
	for _, shard := range shardingList {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", shard, nil),
			"endpoint_name": utils.PathSearch("endpoint_name", shard, nil),
			"ip":            utils.PathSearch("ip", shard, nil),
			"db_port":       utils.PathSearch("db_port", shard, nil),
			"db_user":       utils.PathSearch("db_user", shard, nil),
			"db_password":   utils.PathSearch("db_password", shard, nil),
			"instance_id":   utils.PathSearch("instance_id", shard, nil),
			"instance_name": utils.PathSearch("instance_name", shard, nil),
			"db_name":       utils.PathSearch("db_name", shard, nil),
		})
	}
	return result
}

func flattenDrsVpcInfo(resp interface{}) []interface{} {
	vpc := utils.PathSearch("vpc", resp, nil)
	if vpc == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"vpc_id":            utils.PathSearch("vpc_id", vpc, nil),
			"subnet_id":         utils.PathSearch("subnet_id", vpc, nil),
			"security_group_id": utils.PathSearch("security_group_id", vpc, nil),
		},
	}
}

func flattenDrsSslConfig(resp interface{}) []interface{} {
	ssl := utils.PathSearch("ssl", resp, nil)
	if ssl == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"ssl_link":           utils.PathSearch("ssl_link", ssl, false),
			"ssl_cert_name":      utils.PathSearch("ssl_cert_name", ssl, nil),
			"ssl_cert_key":       utils.PathSearch("ssl_cert_key", ssl, nil),
			"ssl_cert_check_sum": utils.PathSearch("ssl_cert_check_sum", ssl, nil),
			"ssl_cert_password":  utils.PathSearch("ssl_cert_password", ssl, nil),
		},
	}
}
