package cbr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBR GET /v3/{project_id}/policies
// @API CBR GET /v3/{project_id}/vaults
func DataSourceVaults() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVaultsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region in which to query the vaults.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vault name.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The object type of the vault.",
			},
			"consistent_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The consistent level (specification) of the vault.",
			},
			"protection_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protection type of the vault.",
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The vault sapacity, in GB.",
			},
			"auto_expand_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable automatic expansion of the backup protection type vault.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the vault belongs.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the policy associated with the vault.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vault status.",
			},
			"vaults": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vault ID in UUID format.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vault name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The object type of the vault.",
						},
						"consistent_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The consistent level (specification) of the vault.",
						},
						"protection_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protection type of the vault.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The vault capacity, in GB.",
						},
						"auto_expand_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable automatic expansion of the backup protection type vault.",
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise project ID.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the policy associated with the vault.",
						},
						"allocated": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The allocated capacity of the vault, in GB.",
						},
						"used": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The used capacity, in GB.",
						},
						"spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The specification code.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vault status.",
						},
						"storage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the bucket for the vault.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The key/value pairs to associate with the vault.",
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the ECS instance to be backed up.",
									},
									"excludes": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The array of disk IDs which will be excluded in the backup.",
									},
									"includes": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The array of disk or SFS file system IDs which will be included in the backup.",
									},
								},
							},
							Description: "The array of one or more resources to attach to the vault.",
						},
						"auto_bind": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether automatic association is supported.",
						},
						"bind_rules": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The rules for automatic association.",
						},
					},
				},
			},
		},
	}
}

func buildListVaultsParams(d *schema.ResourceData) string {
	res := "&cloud_type=public"
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&object_type=%v", res, v)
	}
	if v, ok := d.GetOk("protection_type"); ok {
		res = fmt.Sprintf("%s&protect_type=%v", res, v)
	}
	if v, ok := d.GetOk("policy_id"); ok {
		res = fmt.Sprintf("%s&policy_id=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	return res
}

func filterVaults(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("consistent_level"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("billing.consistent_level", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("size"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("billing.size", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("auto_expand_enabled"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("auto_expand", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenAllVaults(client *golangsdk.ServiceClient, vaultList []interface{}) []map[string]interface{} {
	if len(vaultList) < 1 {
		return nil
	}
	results := make([]map[string]interface{}, 0, len(vaultList))
	for _, val := range vaultList {
		vaultId := utils.PathSearch("id", val, "").(string)
		objectType := utils.PathSearch("billing.object_type", val, "").(string)
		vMap := map[string]interface{}{
			"id":                    vaultId,
			"name":                  utils.PathSearch("name", val, ""),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", val, ""),
			"type":                  objectType,
			"protection_type":       utils.PathSearch("billing.protect_type", val, ""),
			"status":                utils.PathSearch("billing.status", val, ""),
			"consistent_level":      utils.PathSearch("billing.consistent_level", val, ""),
			"size":                  utils.PathSearch("billing.size", val, 0),
			"allocated":             utils.PathSearch("billing.allocated", val, 0),
			"used":                  utils.PathSearch("billing.used", val, 0),
			"spec_code":             utils.PathSearch("billing.spec_code", val, ""),
			"storage":               utils.PathSearch("billing.storage_unit", val, ""),
			"auto_expand_enabled":   utils.PathSearch("auto_expand", val, false),
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", val, make(map[string]interface{}))),
			"resources":             flattenVaultResources(objectType, utils.PathSearch("resources", val, make([]interface{}, 0)).([]interface{})),
			"auto_bind":             utils.PathSearch("auto_bind", val, nil),
			"bind_rules":            utils.FlattenTagsToMap(utils.PathSearch("bind_rules.tags", val, nil)),
		}

		// Query the CBR policy which bound to the vault by ID.
		if policies, err := getPoliciesByVaultId(client, vaultId); err != nil {
			log.Printf("[DEBUG] No policy bound to vault (%s): %s", vaultId, err)
		} else {
			vMap["policy_id"] = utils.PathSearch("[0].id", policies, "")
		}
		results = append(results, vMap)
	}

	return results
}

func queryVaults(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/vaults?limit=100"
		offset  = 0
		results = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListVaultsParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error querying vaults: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		vaults := utils.PathSearch("vaults", respBody, make([]interface{}, 0)).([]interface{})
		if len(vaults) < 1 {
			break
		}
		results = append(results, vaults...)
		offset += len(vaults)
	}
	return results, nil
}

func dataSourceVaultsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	vaults, err := queryVaults(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the ID and other parameters.
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)
	mErr := multierror.Append(nil,
		d.Set("vaults", flattenAllVaults(client, filterVaults(vaults, d))),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting data-source fields: %s", err)
	}
	return nil
}
