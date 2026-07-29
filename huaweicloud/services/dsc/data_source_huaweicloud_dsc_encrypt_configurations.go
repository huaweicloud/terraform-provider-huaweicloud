package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/sdg/server/mask/algorithms/encryption-configurations
func DataSourceEncryptConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEncryptConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the encryption configurations are located.`,
			},

			// Required parameters.
			"algorithm_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the encryption algorithm.`,
			},

			// Optional parameters.
			"configuration_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the encryption configuration.`,
			},

			// Attributes.
			"configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        encryptConfigurationsSchema(),
				Description: `The list of the encryption configurations.`,
			},
			"access_permission": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the user has the access permission.`,
			},
		},
	}
}

func encryptConfigurationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the encryption configuration.`,
			},
			"configuration_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the encryption configuration.`,
			},
			"algorithm_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the encryption algorithm.`,
			},
			"algorithm_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the encryption algorithm.`,
			},
			"enable_rotate": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the key rotation is enabled.`,
			},
			"encrypt_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The encryption mode.`,
			},
			"filling_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The filling method used for encryption masking.`,
			},
			"kms_context": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The KMS context information.`,
				Elem:        encryptConfigurationsKmsContextSchema(),
			},
			"mask_task_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of the masking tasks.`,
			},
			"rotate_period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The key rotation period, in days.`,
			},
		},
	}
}

func encryptConfigurationsKmsContextSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kms_key_alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alias of the KMS key.`,
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the KMS key.`,
			},
			"kms_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The region where the KMS key is located.`,
			},
		},
	}
}

func buildEncryptConfigurationsQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("&algorithm_type=%v", d.Get("algorithm_type"))
	if v, ok := d.GetOk("configuration_name"); ok {
		queryParams = fmt.Sprintf("%s&configuration_name=%v", queryParams, v)
	}

	return queryParams
}

func listEncryptConfigurations(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, interface{}, error) {
	var (
		httpUrl          = "v1/{project_id}/sdg/server/mask/algorithms/encryption-configurations"
		limit            = 200
		offset           = 0
		results          = make([]interface{}, 0)
		accessPermission interface{}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d%s", listPath, limit, buildEncryptConfigurationsQueryParams(d))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, nil, err
		}

		if offset == 0 {
			accessPermission = utils.PathSearch("access_permission", respBody, nil)
		}

		configuration := utils.PathSearch("configuration_list", respBody, make([]interface{}, 0)).([]interface{})
		results = append(results, configuration...)
		if len(configuration) < limit {
			break
		}

		offset += len(configuration)
	}

	return results, accessPermission, nil
}

func dataSourceEncryptConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	encryptConfigs, accessPermission, err := listEncryptConfigurations(client, d)
	if err != nil {
		return diag.Errorf("error retrieving encryption configurations: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("access_permission", accessPermission),
		d.Set("configurations", flattenEncryptConfigurations(encryptConfigs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEncryptConfigurations(encryptConfigs []interface{}) []map[string]interface{} {
	if len(encryptConfigs) < 1 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(encryptConfigs))
	for _, v := range encryptConfigs {
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"configuration_name": utils.PathSearch("configuration_name", v, nil),
			"algorithm_name":     utils.PathSearch("algorithm_name", v, nil),
			"algorithm_type":     utils.PathSearch("algorithm_type", v, nil),
			"enable_rotate":      utils.PathSearch("enable_rotate", v, nil),
			"encrypt_mode":       utils.PathSearch("encrypt_mode", v, nil),
			"filling_method":     utils.PathSearch("filling_method", v, nil),
			"kms_context": flattenEncryptConfigurationsKmsContext(utils.PathSearch("kms_context",
				v, make(map[string]interface{}, 0)).(map[string]interface{})),
			"mask_task_num": utils.PathSearch("mask_task_num", v, nil),
			"rotate_period": utils.PathSearch("rotate_period", v, nil),
		})
	}
	return rst
}

func flattenEncryptConfigurationsKmsContext(kmsContext map[string]interface{}) []map[string]interface{} {
	if len(kmsContext) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"kms_key_alias": utils.PathSearch("kms_key_alias", kmsContext, nil),
			"kms_key_id":    utils.PathSearch("kms_key_id", kmsContext, nil),
			"kms_region":    utils.PathSearch("kms_region", kmsContext, nil),
		},
	}
}
