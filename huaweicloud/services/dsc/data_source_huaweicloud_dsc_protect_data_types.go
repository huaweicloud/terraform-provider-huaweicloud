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

// @API DSC GET /v1/{project_id}/data-protect-types
func DataSourceDscProtectDataTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscProtectDataTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"life_cycle": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataTypeDetailSchema(),
			},
		},
	}
}

func dataTypeDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"life_cycle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDscProtectDataTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "dsc"
		httpUrl    = "v1/{project_id}/data-protect-types"
		lifeCycle  = d.Get("life_cycle").(string)
		limit      = 1000
		offset     = 0
		allDetails = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	for {
		currentPath := fmt.Sprintf("%s?life_cycle=%s&offset=%d&limit=%d", requestPath, lifeCycle, offset, limit)
		requestOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving DSC protect data types: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		detailsRaw, ok := respBody.([]interface{})

		if !ok {
			detailsRaw = make([]interface{}, 0)
		}

		allDetails = append(allDetails, detailsRaw...)
		if len(detailsRaw) < limit {
			break
		}

		offset += len(detailsRaw)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenDataTypeDetails(allDetails)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataTypeDetails(details []interface{}) []interface{} {
	if len(details) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(details))
	for _, item := range details {
		result = append(result, map[string]interface{}{
			"category":     utils.PathSearch("category", item, nil),
			"create_time":  utils.PathSearch("create_time", item, nil),
			"data_type":    utils.PathSearch("data_type", item, nil),
			"data_type_id": utils.PathSearch("data_type_id", item, nil),
			"protect_id":   utils.PathSearch("id", item, nil),
			"is_deleted":   utils.PathSearch("is_deleted", item, nil),
			"life_cycle":   utils.PathSearch("life_cycle", item, nil),
			"update_time":  utils.PathSearch("update_time", item, nil),
		})
	}

	return result
}
