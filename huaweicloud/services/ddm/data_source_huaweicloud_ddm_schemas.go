// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM GET /v1/{project_id}/instances/{instance_id}/databases
func DataSourceDdmSchemas() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmSchemasRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of DDM instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the DDM schema.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the DDM schema.`,
			},
			"shard_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sharding mode of the schema.`,
				ValidateFunc: validation.StringInSlice([]string{
					"cluster", "single",
				}, false),
			},
			"schemas": {
				Type:        schema.TypeList,
				Elem:        SchemasSchemaSchema(),
				Computed:    true,
				Description: `Indicates the list of DDM schema.`,
			},
		},
	}
}

func SchemasSchemaSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the DDM schema.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DDM schema.`,
			},
			"shard_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the sharding mode of the schema.`,
			},
			"shard_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of shards in the same working mode.`,
			},
			"data_nodes": {
				Type:        schema.TypeList,
				Elem:        SchemasSchemaDataNodeSchema(),
				Computed:    true,
				Description: `Indicates the RDS instances associated with the schema.`,
			},
		},
	}
	return &sc
}

func SchemasSchemaDataNodeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the RDS instance associated with the schema.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the username for logging in to the associated RDS instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the associated RDS instance.`,
			},
		},
	}
	return &sc
}

func resourceDdmSchemasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var mErr *multierror.Error

	// getDdmSchemas: Query the List of DDM schema
	var (
		getDdmSchemasHttpUrl = "v1/{project_id}/instances/{instance_id}/databases"
		getDdmSchemasProduct = "ddm"
	)
	getDdmSchemasClient, err := cfg.NewServiceClient(getDdmSchemasProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	getDdmSchemasPath := getDdmSchemasClient.Endpoint + getDdmSchemasHttpUrl
	getDdmSchemasPath = strings.ReplaceAll(getDdmSchemasPath, "{project_id}", getDdmSchemasClient.ProjectID)
	getDdmSchemasPath = strings.ReplaceAll(getDdmSchemasPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getDdmSchemasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	var currentTotal int
	getDdmSchemasPath += buildGetDdmSchemasQueryParams(currentTotal)

	schemas := make([]interface{}, 0)
	for {
		getDdmSchemasResp, err := getDdmSchemasClient.Request("GET", getDdmSchemasPath, &getDdmSchemasOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving DdmSchemas")
		}
		getDdmSchemasRespBody, err := utils.FlattenResponse(getDdmSchemasResp)
		if err != nil {
			return diag.FromErr(err)
		}
		total, curPageSize, res := flattenGetSchemasResponseBodySchema(d, getDdmSchemasRespBody)
		schemas = append(schemas, res...)
		currentTotal += curPageSize
		if currentTotal == total {
			break
		}
		getDdmSchemasPath = updatePathOffset(getDdmSchemasPath, currentTotal)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("schemas", schemas),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetSchemasResponseBodySchema(d *schema.ResourceData, resp interface{}) (total, curPageSize int, rst []interface{}) {
	if resp == nil {
		return 0, 0, nil
	}
	total = int(utils.PathSearch("total", resp, float64(0)).(float64))
	curJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst = make([]interface{}, 0, len(curArray))
	name := d.Get("name").(string)
	status := d.Get("status").(string)
	shardMode := d.Get("shard_mode").(string)
	for _, v := range curArray {
		schemaName := utils.PathSearch("name", v, nil)
		schemaStatus := utils.PathSearch("status", v, nil)
		schemaShardMode := utils.PathSearch("shard_mode", v, nil)
		if name != "" && name != schemaName {
			continue
		}
		if status != "" && status != schemaStatus {
			continue
		}
		if shardMode != "" && shardMode != schemaShardMode {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"name":         schemaName,
			"status":       schemaStatus,
			"shard_mode":   schemaShardMode,
			"shard_number": utils.PathSearch("shard_number", v, nil),
			"data_nodes":   flattenSchemaDataNode(v),
		})
	}
	return total, len(curArray), rst
}

func flattenSchemaDataNode(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("used_rds", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":     utils.PathSearch("id", v, nil),
			"name":   utils.PathSearch("name", v, nil),
			"status": utils.PathSearch("status", v, nil),
		})
	}
	return rst
}

func buildGetDdmSchemasQueryParams(offset int) string {
	return fmt.Sprintf("?offset=%v", offset)
}
