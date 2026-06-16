package geminidb

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

// @API GeminiDB GET /v3/{project_id}/instances/{instance_id}/available-flavors
func DataSourceAvailableFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailableFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"current_flavor": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     availableFlavorsSchema(),
			},
			"optional_flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     availableFlavorsSchema(),
			},
		},
	}
}

func availableFlavorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"vcpus": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func dataSourceAvailableFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v3/{project_id}/instances/{instance_id}/available-flavors?limit=100"
		offset        = 0
		result        = make([]interface{}, 0)
		currentFlavor interface{}
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the GeminiDB instance available flavors: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		currentFlavor = utils.PathSearch("current_flavor", getRespBody, nil)
		flavors := utils.PathSearch("optional_flavors.list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(flavors) == 0 {
			break
		}

		result = append(result, flavors...)
		offset += len(flavors)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("current_flavor", flattenCurrentFlavor(currentFlavor)),
		d.Set("optional_flavors", flattenAvailableFlavors(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAvailableFlavors(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"vcpus":     utils.PathSearch("vcpus", v, nil),
			"ram":       utils.PathSearch("ram", v, nil),
			"spec_code": utils.PathSearch("spec_code", v, nil),
			"az_status": utils.PathSearch("az_status", v, nil),
		})
	}

	return result
}

func flattenCurrentFlavor(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"vcpus":     utils.PathSearch("vcpus", resp, nil),
		"ram":       utils.PathSearch("ram", resp, nil),
		"spec_code": utils.PathSearch("spec_code", resp, nil),
		"az_status": utils.PathSearch("az_status", resp, nil),
	}

	return []map[string]interface{}{result}
}
