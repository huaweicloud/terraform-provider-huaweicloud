package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dimensionErrCodes = []string{
	"DLG.0818", // Workspace not found.
	"DLG.6026", // Resource not found.
	"DLG.3902", // Resource ID value is incorrect.
}

var architectureDimensionNonUpdatableParams = []string{
	"workspace_id",
}

// @API DataArtsStudio POST /v2/{project_id}/design/dimensions
// @API DataArtsStudio DELETE /v2/{project_id}/design/dimensions
// @API DataArtsStudio GET /v2/{project_id}/design/dimensions/{id}
// @API DataArtsStudio PUT /v2/{project_id}/design/dimensions
func ResourceArchitectureDimension() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureDimensionCreate,
		ReadContext:   resourceArchitectureDimensionRead,
		UpdateContext: resourceArchitectureDimensionUpdate,
		DeleteContext: resourceArchitectureDimensionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDataArtsStudioImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(architectureDimensionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the dimension is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workspace ID to which the dimension belongs.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The Chinese name of the dimension.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The English name of the dimension.`,
			},
			"dimension_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the dimension.`,
			},
			"l3_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The business object ID to which the dimension belongs.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The asset owner of the dimension.`,
			},
			"attributes": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        dimensionAttributeSchema(),
				Description: `The dimension attribute information.`,
			},
			"datasource": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        dimensionDatasourceSchema(),
				Description: `The data source information of the dimension.`,
			},

			// Optional parameters.
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The alias of the dimension.`,
			},
			"code_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The referenced code table ID.`,
			},
			"configs": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The other configuration information.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the dimension.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the dimension.`,
			},
			"distribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The distribution mode.`,
			},
			"distribute_column": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The DISTRIBUTE BY HASH column.`,
			},
			"hierarchies": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        dimensionHierarchySchema(),
				Description: `The hierarchy attribute definitions of the dimension.`,
			},
			"id_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID field of the dimension.`,
			},
			"l1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Chinese name of the subject domain group.`,
			},
			"l2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Chinese name of the subject domain.`,
			},
			"l2_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The subject domain ID.`,
			},
			"l3": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Chinese name of the business object.`,
			},
			"mappings": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        dimensionMappingSchema(),
				Description: `The table mapping information of the dimension.`,
			},
			"model": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        dimensionModelSchema(),
				Description: `The model information of the dimension.`,
			},
			"model_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The model ID to which the dimension belongs.`,
			},
			"obs_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The OBS external table path.`,
			},
			"self_defined_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        dimensionSelfDefinedFieldsSchema(),
				Description: `The custom extended fields of the dimension.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The publish status of the dimension.`,
			},
			"table_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The table type of the dimension.`,
			},
			"update_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The updater of the dimension.`,
			},

			// Attributes.
			"approval_info": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        dimensionApprovalInfoSchema(),
				Description: `The approval information.`,
			},
			"code_table": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        dimensionCodeTableSchema(),
				Description: `The referenced code table.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"dev_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The development environment version.`,
			},
			"dev_version_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The development environment version name.`,
			},
			"env_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The development and production environment type.`,
			},
			"l1_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subject domain group ID.`,
			},
			"new_biz": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        dimensionNewBizSchema(),
				Description: `The business version management information.`,
			},
			"prod_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The production environment version.`,
			},
			"prod_version_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The production environment version name.`,
			},
		},
	}
}

func dimensionSelfDefinedFieldsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fd_name_ch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Chinese display name of the custom extended field.`,
			},
			"fd_name_en": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The English name of the custom extended field.`,
			},
			"fd_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The value of the custom extended field.`,
			},
			"not_null": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the custom extended field requires a value.`,
			},
		},
	}
}

func dimensionSelfDefinedFieldsComputedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fd_name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese display name of the custom extended field.`,
			},
			"fd_name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the custom extended field.`,
			},
			"fd_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the custom extended field.`,
			},
			"not_null": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the custom extended field requires a value.`,
			},
		},
	}
}

func dimensionQualityInfosSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the quality info.`,
			},
			"alert_conf": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alert configuration.`,
			},
			"attr_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The attribute ID.`,
			},
			"biz_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business entity type.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the quality info.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"data_quality_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data quality ID.`,
			},
			"data_quality_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data quality name.`,
			},
			"expression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The regular expression configuration.`,
			},
			"extend_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended information.`,
			},
			"from_standard": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether it is from data standard quality configuration.`,
			},
			"result_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The result description.`,
			},
			"show_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether to display the regular expression.`,
			},
			"table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The table ID.`,
			},
			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the quality info.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
		},
	}
}

func dimensionNewBizSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business version ID.`,
			},
			"biz_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business entity type.`,
			},
		},
	}
}

func dimensionSecrecyLevelsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The secrecy level ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The secrecy level name.`,
			},
			"slevel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The secrecy level rank.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the secrecy level.`,
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data security primary key.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the secrecy level.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the secrecy level.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"new_biz": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        dimensionNewBizSchema(),
				Description: `The business version management information.`,
			},
		},
	}
}

func dimensionAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The data type of the attribute.`,
			},
			"is_primary_key": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether the attribute is a primary key.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The Chinese name of the attribute.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The English name of the attribute.`,
			},
			"ordinal": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The sequence number of the attribute.`,
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The alias of the attribute.`,
			},
			"code_table_field_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The code table field ID of the attribute.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the attribute.`,
			},
			"data_type_extend": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The data type extend field of the attribute.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the attribute.`,
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the attribute.`,
			},
			"is_biz_primary": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the attribute is a business primary key.`,
			},
			"is_partition_key": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether the attribute is a partition key.`,
			},
			"not_null": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether the attribute is not null.`,
			},
			"self_defined_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        dimensionSelfDefinedFieldsSchema(),
				Description: `The custom extended fields of the attribute.`,
			},
			"stand_row_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the associated data standard.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publish status of the attribute.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the attribute.`,
			},
			"dimension_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The dimension ID of the attribute.`,
			},
			"domain_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain type of the attribute.`,
			},
			"stand_row_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the associated data standard.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the attribute.`,
			},
			"quality_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionQualityInfosSchema(),
				Description: `The quality information of the attribute.`,
			},
			"secrecy_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionSecrecyLevelsSchema(),
				Description: `The secrecy levels of the attribute.`,
			},
		},
	}
}

func dimensionDatasourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dw_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The data connection ID.`,
			},
			"dw_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The data connection type.`,
			},
			"biz_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The business object ID.`,
			},
			"biz_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The business object type.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The database name.`,
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The data source ID.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The DLI queue name.`,
			},
			"schema": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The DWS schema name.`,
			},
		},
	}
}

func dimensionAttrDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alias of the attribute.`,
			},
			"code_table_field_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code table field ID.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the attribute.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"data_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data type of the attribute.`,
			},
			"data_type_extend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data type extend field of the attribute.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the attribute.`,
			},
			"dimension_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The dimension ID.`,
			},
			"domain_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain type of the attribute.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The attribute ID.`,
			},
			"is_biz_primary": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is a business primary key.`,
			},
			"is_partition_key": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is a partition key.`,
			},
			"is_primary_key": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is a primary key.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the attribute.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the attribute.`,
			},
			"not_null": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is not null.`,
			},
			"ordinal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sequence number of the attribute.`,
			},
			"stand_row_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The associated data standard ID.`,
			},
			"stand_row_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The associated data standard name.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publish status of the attribute.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"quality_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionQualityInfosSchema(),
				Description: `The quality information of the attribute.`,
			},
			"secrecy_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionSecrecyLevelsSchema(),
				Description: `The secrecy levels of the attribute.`,
			},
			"self_defined_fields": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionSelfDefinedFieldsComputedSchema(),
				Description: `The custom extended fields of the attribute.`,
			},
		},
	}
}

func dimensionHierarchyAttrSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attr": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionAttrDetailSchema(),
				Description: `The referenced attribute details.`,
			},
			"detail_attrs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionAttrDetailSchema(),
				Description: `The detail attributes.`,
			},
			"attr_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The attribute ID.`,
			},
			"attr_name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The referenced attribute code.`,
			},
			"detail_attr_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The detail attribute IDs.`,
			},
			"detail_attr_name_ens": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The detail attribute English names.`,
			},
			"hierarchies_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The hierarchy ID.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The hierarchy attribute ID.`,
			},
			"level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The hierarchy level.`,
			},
		},
	}
}

func dimensionHierarchySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The hierarchy ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The hierarchy name.`,
			},
			"attrs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionHierarchyAttrSchema(),
				Description: `The attributes contained in the hierarchy.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the hierarchy.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the hierarchy.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
		},
	}
}

func dimensionMappingJoinFieldsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"field1_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The field 1 ID.`,
			},
			"field2_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The field 2 ID.`,
			},
			"field1_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The field 1 name.`,
			},
			"field2_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The field 2 name.`,
			},
		},
	}
}

func dimensionMappingSourceTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table1_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The table 1 ID.`,
			},
			"table2_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The table 2 ID.`,
			},
			"join_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The join type.`,
			},
			"table1_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The table 1 name.`,
			},
			"table2_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The table 2 name.`,
			},
			"join_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        dimensionMappingJoinFieldsSchema(),
				Description: `The ON condition fields.`,
			},
		},
	}
}

func dimensionMappingDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_attr_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The target attribute name.`,
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The detail ID.`,
			},
			"mapping_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The mapping ID.`,
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The remark of the mapping detail.`,
			},
			"src_attr_ids": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The source attribute IDs.`,
			},
			"src_table_ids": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The source table IDs.`,
			},
			"target_attr_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The target attribute ID.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the mapping detail.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the mapping detail.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
		},
	}
}

func dimensionMappingSourceFieldsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_field_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The target field code.`,
			},
			"field_ids": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The source field IDs, multiple IDs separated by commas.`,
			},
			"field_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The source field name list.`,
			},
			"target_field_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The target field ID.`,
			},
			"transform_expression": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The transform expression.`,
			},
			"changed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the field has changed.`,
			},
		},
	}
}

func dimensionMappingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The mapping name.`,
			},
			"source_tables": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        dimensionMappingSourceTablesSchema(),
				Description: `The source table information of the mapping.`,
			},
			"details": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        dimensionMappingDetailsSchema(),
				Description: `The mapping details.`,
			},
			"source_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        dimensionMappingSourceFieldsSchema(),
				Description: `The source field information of the mapping.`,
			},
			"src_model_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The source model ID in relational modeling.`,
			},
			"src_model_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The source model name in relational modeling.`,
			},
			"view_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The collected view source.`,
			},
			"target_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The target table ID.`,
			},
			"target_table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The target table name.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The mapping ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The mapping description.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the mapping.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the mapping.`,
			},
		},
	}
}

func dimensionModelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workspace name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workspace type.`,
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The workspace ID.`,
			},
		},
	}
}

func dimensionBizInfoObjSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"biz_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publish status of the business.`,
			},
			"biz_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business entity type.`,
			},
			"biz_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business version.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the business.`,
			},
			"directory_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory tree.`,
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The email address of the creator.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business information ID.`,
			},
			"l1": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subject domain group Chinese name.`,
			},
			"l2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subject domain Chinese name.`,
			},
			"l3": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business object Chinese name.`,
			},
			"msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval or submission message.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the business.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the business.`,
			},
			"submit_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The submit time.`,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The tenant ID.`,
			},
		},
	}
}

func dimensionApprovalInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"approval_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval status.`,
			},
			"approval_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval time.`,
			},
			"approval_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval type.`,
			},
			"approver": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approver.`,
			},
			"biz_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business ID.`,
			},
			"biz_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The serialized business details.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval ID.`,
			},
			"msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval information.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business Chinese name.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business English name.`,
			},
			"submit_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The submit time.`,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID.`,
			},
			"directory_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory tree.`,
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approver email.`,
			},
			"biz_info_obj": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        dimensionBizInfoObjSchema(),
				Description: `The business details object.`,
			},
		},
	}
}

func dimensionCodeTableFieldValuesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code table field value ID.`,
			},
			"fd_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code table attribute ID.`,
			},
			"fd_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code table attribute value.`,
			},
			"ordinal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sequence number of the field value.`,
			},
		},
	}
}

func dimensionCodeTableFieldsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code table field ID.`,
			},
			"code_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code table ID to which the field belongs.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the field.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the field.`,
			},
			"ordinal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sequence number of the field.`,
			},
			"data_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data type of the field.`,
			},
			"data_type_extend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended data type information of the field.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the field.`,
			},
			"domain_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain type of the field.`,
			},
			"is_unique_key": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the field is unique.`,
			},
			"count_field_values": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The total count of field values.`,
			},
			"code_table_field_values": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionCodeTableFieldValuesSchema(),
				Description: `The code table field values.`,
			},
		},
	}
}

func dimensionCodeTableSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code table ID.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the code table.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the code table.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publish status of the code table.`,
			},
			"tb_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the code table.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the code table.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"directory_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory tree.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the code table.`,
			},
			"directory_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory ID.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"new_biz": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        dimensionNewBizSchema(),
				Description: `The business version management information.`,
			},
			"approval_info": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        dimensionApprovalInfoSchema(),
				Description: `The approval information.`,
			},
			"code_table_fields": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dimensionCodeTableFieldsSchema(),
				Description: `The code table field information.`,
			},
		},
	}
}

func buildDimensionMappingsSourceTablesJoinFields(raw []interface{}) []map[string]interface{} {
	if len(raw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(raw))
	for _, val := range raw {
		result = append(result, map[string]interface{}{
			"field1_id":   utils.PathSearch("field1_id", val, nil),
			"field2_id":   utils.PathSearch("field2_id", val, nil),
			"field1_name": utils.ValueIgnoreEmpty(utils.PathSearch("field1_name", val, nil)),
			"field2_name": utils.ValueIgnoreEmpty(utils.PathSearch("field2_name", val, nil)),
		})
	}
	return result
}

func buildDimensionAttributes(raw []interface{}) []map[string]interface{} {
	if len(raw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(raw))
	for _, val := range raw {
		result = append(result, map[string]interface{}{
			"name_ch":             utils.PathSearch("name_ch", val, nil),
			"name_en":             utils.PathSearch("name_en", val, nil),
			"data_type":           utils.PathSearch("data_type", val, nil),
			"is_primary_key":      utils.PathSearch("is_primary_key", val, nil),
			"ordinal":             utils.PathSearch("ordinal", val, nil),
			"alias":               utils.ValueIgnoreEmpty(utils.PathSearch("alias", val, nil)),
			"code_table_field_id": utils.ValueIgnoreEmpty(utils.PathSearch("code_table_field_id", val, nil)),
			"create_by":           utils.ValueIgnoreEmpty(utils.PathSearch("create_by", val, nil)),
			"data_type_extend":    utils.ValueIgnoreEmpty(utils.PathSearch("data_type_extend", val, nil)),
			"description":         utils.ValueIgnoreEmpty(utils.PathSearch("description", val, nil)),
			"id":                  utils.ValueIgnoreEmpty(utils.PathSearch("id", val, nil)),
			"is_biz_primary":      utils.ValueIgnoreEmpty(utils.PathSearch("is_biz_primary", val, nil)),
			"is_partition_key":    utils.ValueIgnoreEmpty(utils.PathSearch("is_partition_key", val, nil)),
			"not_null":            utils.ValueIgnoreEmpty(utils.PathSearch("not_null", val, nil)),
			"stand_row_id":        utils.ValueIgnoreEmpty(utils.PathSearch("stand_row_id", val, nil)),
			"self_defined_fields": buildDimensionSelfDefinedFields(
				utils.PathSearch("self_defined_fields", val, make([]interface{}, 0)).([]interface{})),
			"secrecy_levels":      nil, // computed-only
			"quality_infos":       nil, // computed-only
		})
	}
	return result
}

func buildDimensionDatasource(raw []interface{}) interface{} {
	if len(raw) == 0 {
		return nil
	}

	return map[string]interface{}{
		"dw_id":      utils.PathSearch("dw_id", raw[0], nil),
		"dw_type":    utils.PathSearch("dw_type", raw[0], nil),
		"biz_id":     utils.ValueIgnoreEmpty(utils.PathSearch("biz_id", raw[0], nil)),
		"biz_type":   utils.ValueIgnoreEmpty(utils.PathSearch("biz_type", raw[0], nil)),
		"db_name":    utils.ValueIgnoreEmpty(utils.PathSearch("db_name", raw[0], nil)),
		"id":         utils.ValueIgnoreEmpty(utils.PathSearch("id", raw[0], nil)),
		"queue_name": utils.ValueIgnoreEmpty(utils.PathSearch("queue_name", raw[0], nil)),
		"schema":     utils.ValueIgnoreEmpty(utils.PathSearch("schema", raw[0], nil)),
	}
}

func buildDimensionHierarchies(raw []interface{}) []map[string]interface{} {
	if len(raw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(raw))
	for _, val := range raw {
		result = append(result, map[string]interface{}{
			"id":   utils.ValueIgnoreEmpty(utils.PathSearch("id", val, nil)),
			"name": utils.ValueIgnoreEmpty(utils.PathSearch("name", val, nil)),
		})
	}
	return result
}

func buildDimensionMappingsSourceTables(raw []interface{}) []map[string]interface{} {
	if len(raw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(raw))
	for _, val := range raw {
		result = append(result, map[string]interface{}{
			"join_type":   utils.ValueIgnoreEmpty(utils.PathSearch("join_type", val, nil)),
			"table1_id":   utils.ValueIgnoreEmpty(utils.PathSearch("table1_id", val, nil)),
			"table2_id":   utils.ValueIgnoreEmpty(utils.PathSearch("table2_id", val, nil)),
			"table1_name": utils.ValueIgnoreEmpty(utils.PathSearch("table1_name", val, nil)),
			"table2_name": utils.ValueIgnoreEmpty(utils.PathSearch("table2_name", val, nil)),
			"join_fields": buildDimensionMappingsSourceTablesJoinFields(
				utils.PathSearch("join_fields", val, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func buildDimensionMappingsSourceFields(raw []interface{}) []map[string]interface{} {
	if len(raw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(raw))
	for _, val := range raw {
		result = append(result, map[string]interface{}{
			"field_ids":            utils.ValueIgnoreEmpty(utils.PathSearch("field_ids", val, nil)),
			"field_names":          utils.ValueIgnoreEmpty(utils.PathSearch("field_names", val, nil)),
			"target_field_name":    utils.ValueIgnoreEmpty(utils.PathSearch("target_field_name", val, nil)),
			"transform_expression": utils.ValueIgnoreEmpty(utils.PathSearch("transform_expression", val, nil)),
		})
	}
	return result
}

func buildDimensionMappingsDetails(raw []interface{}) []map[string]interface{} {
	if len(raw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(raw))
	for _, val := range raw {
		result = append(result, map[string]interface{}{
			"target_attr_name": utils.PathSearch("target_attr_name", val, nil),
			"remark":           utils.ValueIgnoreEmpty(utils.PathSearch("remark", val, nil)),
			"src_attr_ids":     utils.ValueIgnoreEmpty(utils.PathSearch("src_attr_ids", val, nil)),
			"src_table_ids":    utils.ValueIgnoreEmpty(utils.PathSearch("src_table_ids", val, nil)),
		})
	}
	return result
}

func buildDimensionMappings(raw []interface{}) []map[string]interface{} {
	if len(raw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(raw))
	for _, val := range raw {
		result = append(result, map[string]interface{}{
			"name":           utils.PathSearch("name", val, nil),
			"src_model_id":   utils.ValueIgnoreEmpty(utils.PathSearch("src_model_id", val, nil)),
			"src_model_name": utils.ValueIgnoreEmpty(utils.PathSearch("src_model_name", val, nil)),
			"view_text":      utils.ValueIgnoreEmpty(utils.PathSearch("view_text", val, nil)),
			"source_tables":  buildDimensionMappingsSourceTables(utils.PathSearch("source_tables", val, make([]interface{}, 0)).([]interface{})),
			"source_fields":  buildDimensionMappingsSourceFields(utils.PathSearch("source_fields", val, make([]interface{}, 0)).([]interface{})),
			"details":        buildDimensionMappingsDetails(utils.PathSearch("details", val, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func buildDimensionModel(raw []interface{}) interface{} {
	if len(raw) == 0 {
		return nil
	}

	return map[string]interface{}{
		"name": utils.PathSearch("name", raw[0], nil),
		"type": utils.PathSearch("type", raw[0], nil),
		"id":   utils.ValueIgnoreEmpty(utils.PathSearch("id", raw[0], nil)),
	}
}

func buildDimensionSelfDefinedFields(raw []interface{}) []map[string]interface{} {
	if len(raw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(raw))
	for _, val := range raw {
		result = append(result, map[string]interface{}{
			"fd_name_ch": utils.PathSearch("fd_name_ch", val, nil),
			"fd_name_en": utils.PathSearch("fd_name_en", val, nil),
			"fd_value":   utils.PathSearch("fd_value", val, nil),
			"not_null":   utils.PathSearch("not_null", val, nil),
		})
	}
	return result
}

func buildCreateOrUpdateDimensionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name_ch":             d.Get("name_ch"),
		"name_en":             d.Get("name_en"),
		"dimension_type":      d.Get("dimension_type"),
		"l3_id":               d.Get("l3_id"),
		"owner":               d.Get("owner"),
		"attributes":          buildDimensionAttributes(d.Get("attributes").([]interface{})),
		"datasource":          buildDimensionDatasource(d.Get("datasource").([]interface{})),
		"alias":               utils.ValueIgnoreEmpty(d.Get("alias")),
		"code_table_id":       utils.ValueIgnoreEmpty(d.Get("code_table_id")),
		"configs":             utils.ValueIgnoreEmpty(d.Get("configs")),
		"create_by":           utils.ValueIgnoreEmpty(d.Get("create_by")),
		"description":         utils.ValueIgnoreEmpty(d.Get("description")),
		"distribute":          utils.ValueIgnoreEmpty(d.Get("distribute")),
		"distribute_column":   utils.ValueIgnoreEmpty(d.Get("distribute_column")),
		"hierarchies":         buildDimensionHierarchies(d.Get("hierarchies").([]interface{})),
		"l1":                  utils.ValueIgnoreEmpty(d.Get("l1")),
		"l2":                  utils.ValueIgnoreEmpty(d.Get("l2")),
		"l2_id":               utils.ValueIgnoreEmpty(d.Get("l2_id")),
		"l3":                  utils.ValueIgnoreEmpty(d.Get("l3")),
		"mappings":            buildDimensionMappings(d.Get("mappings").([]interface{})),
		"model":               buildDimensionModel(d.Get("model").([]interface{})),
		"model_id":            utils.ValueIgnoreEmpty(d.Get("model_id")),
		"obs_location":        utils.ValueIgnoreEmpty(d.Get("obs_location")),
		"self_defined_fields": buildDimensionSelfDefinedFields(d.Get("self_defined_fields").([]interface{})),
		"status":              utils.ValueIgnoreEmpty(d.Get("status")),
		"table_type":          utils.ValueIgnoreEmpty(d.Get("table_type")),
	}
}

func resourceArchitectureDimensionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/design/dimensions"
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateDimensionBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.value.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the DataArts Architecture dimension ID from the API response")
	}
	d.SetId(id)

	return resourceArchitectureDimensionRead(ctx, d, meta)
}

func flattenDimensionQualityInfos(raw interface{}) []map[string]interface{} {
	rawArray, _ := raw.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("id", val, nil),
			"alert_conf":         utils.PathSearch("alert_conf", val, nil),
			"attr_id":            utils.PathSearch("attr_id", val, nil),
			"biz_type":           utils.PathSearch("biz_type", val, nil),
			"create_by":          utils.PathSearch("create_by", val, nil),
			"create_time":        utils.PathSearch("create_time", val, nil),
			"data_quality_id":    utils.PathSearch("data_quality_id", val, nil),
			"data_quality_name":  utils.PathSearch("data_quality_name", val, nil),
			"expression":         utils.PathSearch("expression", val, nil),
			"extend_info":        utils.PathSearch("extend_info", val, nil),
			"from_standard":      utils.PathSearch("from_standard", val, nil),
			"result_description":  utils.PathSearch("result_description", val, nil),
			"show_control":       utils.PathSearch("show_control", val, nil),
			"table_id":           utils.PathSearch("table_id", val, nil),
			"update_by":          utils.PathSearch("update_by", val, nil),
			"update_time":        utils.PathSearch("update_time", val, nil),
		})
	}
	return rst
}

func flattenDimensionSecrecyLevels(raw interface{}) []map[string]interface{} {
	rawArray, _ := raw.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", val, nil),
			"name":        utils.PathSearch("name", val, nil),
			"slevel":      utils.PathSearch("slevel", val, nil),
			"description": utils.PathSearch("description", val, nil),
			"uuid":        utils.PathSearch("uuid", val, nil),
			"create_by":   utils.PathSearch("create_by", val, nil),
			"create_time": utils.PathSearch("create_time", val, nil),
			"update_by":   utils.PathSearch("update_by", val, nil),
			"update_time": utils.PathSearch("update_time", val, nil),
			"new_biz":     flattenDimensionNewBiz(utils.PathSearch("new_biz", val, nil)),
		})
	}
	return rst
}

func flattenDimensionAttributes(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"name_ch":             utils.PathSearch("name_ch", val, nil),
			"name_en":             utils.PathSearch("name_en", val, nil),
			"data_type":           utils.PathSearch("data_type", val, nil),
			"is_primary_key":      utils.PathSearch("is_primary_key", val, false),
			"ordinal":             strconv.FormatInt(int64(utils.PathSearch("ordinal", val, float64(0)).(float64)), 10),
			"alias":               utils.PathSearch("alias", val, nil),
			"code_table_field_id": utils.PathSearch("code_table_field_id", val, nil),
			"create_by":           utils.PathSearch("create_by", val, nil),
			"data_type_extend":    utils.PathSearch("data_type_extend", val, nil),
			"description":         utils.PathSearch("description", val, nil),
			"id":                  utils.PathSearch("id", val, nil),
			"is_biz_primary":      utils.PathSearch("is_biz_primary", val, false),
			"is_partition_key":    utils.PathSearch("is_partition_key", val, false),
			"not_null":            utils.PathSearch("not_null", val, false),
			"stand_row_id":        utils.PathSearch("stand_row_id", val, nil),
			"status":              utils.PathSearch("status", val, nil),
			"create_time":         utils.PathSearch("create_time", val, nil),
			"dimension_id":        utils.PathSearch("dimension_id", val, nil),
			"domain_type":         utils.PathSearch("domain_type", val, nil),
			"stand_row_name":      utils.PathSearch("stand_row_name", val, nil),
			"update_time":         utils.PathSearch("update_time", val, nil),
			"self_defined_fields": flattenDimensionSelfDefinedFields(utils.PathSearch("self_defined_fields", val, make([]interface{}, 0)).([]interface{})),
			"quality_infos":       flattenDimensionQualityInfos(utils.PathSearch("quality_infos", val, make([]interface{}, 0))),
			"secrecy_levels":      flattenDimensionSecrecyLevels(utils.PathSearch("secrecy_levels", val, make([]interface{}, 0))),
		})
	}
	return rst
}

func flattenDimensionDatasource(raw interface{}) []map[string]interface{} {
	if raw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"dw_id":      utils.PathSearch("dw_id", raw, nil),
			"dw_type":    utils.PathSearch("dw_type", raw, nil),
			"biz_id":     utils.PathSearch("biz_id", raw, nil),
			"biz_type":   utils.PathSearch("biz_type", raw, nil),
			"db_name":    utils.PathSearch("db_name", raw, nil),
			"id":         utils.PathSearch("id", raw, nil),
			"queue_name": utils.PathSearch("queue_name", raw, nil),
			"schema":     utils.PathSearch("schema", raw, nil),
		},
	}
}

func flattenDimensionAttrDetail(raw interface{}) []map[string]interface{} {
	if raw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"alias":               utils.PathSearch("alias", raw, nil),
			"code_table_field_id": utils.PathSearch("code_table_field_id", raw, nil),
			"create_by":           utils.PathSearch("create_by", raw, nil),
			"create_time":         utils.PathSearch("create_time", raw, nil),
			"data_type":           utils.PathSearch("data_type", raw, nil),
			"data_type_extend":    utils.PathSearch("data_type_extend", raw, nil),
			"description":         utils.PathSearch("description", raw, nil),
			"dimension_id":        utils.PathSearch("dimension_id", raw, nil),
			"domain_type":         utils.PathSearch("domain_type", raw, nil),
			"id":                  utils.PathSearch("id", raw, nil),
			"is_biz_primary":      utils.PathSearch("is_biz_primary", raw, false),
			"is_partition_key":    utils.PathSearch("is_partition_key", raw, false),
			"is_primary_key":      utils.PathSearch("is_primary_key", raw, false),
			"name_ch":             utils.PathSearch("name_ch", raw, nil),
			"name_en":             utils.PathSearch("name_en", raw, nil),
			"not_null":            utils.PathSearch("not_null", raw, false),
			"ordinal":             utils.PathSearch("ordinal", raw, nil),
			"stand_row_id":        utils.PathSearch("stand_row_id", raw, nil),
			"stand_row_name":      utils.PathSearch("stand_row_name", raw, nil),
			"status":              utils.PathSearch("status", raw, nil),
			"update_time":         utils.PathSearch("update_time", raw, nil),
			"quality_infos":       flattenDimensionQualityInfos(utils.PathSearch("quality_infos", raw, make([]interface{}, 0))),
			"secrecy_levels":      flattenDimensionSecrecyLevels(utils.PathSearch("secrecy_levels", raw, make([]interface{}, 0))),
			"self_defined_fields": flattenDimensionSelfDefinedFields(
				utils.PathSearch("self_defined_fields", raw, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenDimensionDetailAttrs(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"alias":               utils.PathSearch("alias", val, nil),
			"code_table_field_id": utils.PathSearch("code_table_field_id", val, nil),
			"create_by":           utils.PathSearch("create_by", val, nil),
			"create_time":         utils.PathSearch("create_time", val, nil),
			"data_type":           utils.PathSearch("data_type", val, nil),
			"data_type_extend":    utils.PathSearch("data_type_extend", val, nil),
			"description":         utils.PathSearch("description", val, nil),
			"dimension_id":        utils.PathSearch("dimension_id", val, nil),
			"domain_type":         utils.PathSearch("domain_type", val, nil),
			"id":                  utils.PathSearch("id", val, nil),
			"is_biz_primary":      utils.PathSearch("is_biz_primary", val, false),
			"is_partition_key":    utils.PathSearch("is_partition_key", val, false),
			"is_primary_key":      utils.PathSearch("is_primary_key", val, false),
			"name_ch":             utils.PathSearch("name_ch", val, nil),
			"name_en":             utils.PathSearch("name_en", val, nil),
			"not_null":            utils.PathSearch("not_null", val, false),
			"ordinal":             utils.PathSearch("ordinal", val, nil),
			"stand_row_id":        utils.PathSearch("stand_row_id", val, nil),
			"stand_row_name":      utils.PathSearch("stand_row_name", val, nil),
			"status":              utils.PathSearch("status", val, nil),
			"update_time":         utils.PathSearch("update_time", val, nil),
			"quality_infos":       flattenDimensionQualityInfos(utils.PathSearch("quality_infos", val, make([]interface{}, 0))),
			"secrecy_levels":      flattenDimensionSecrecyLevels(utils.PathSearch("secrecy_levels", val, make([]interface{}, 0))),
			"self_defined_fields": flattenDimensionSelfDefinedFields(
				utils.PathSearch("self_defined_fields", val, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenDimensionHierarchyAttrs(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"attr":            flattenDimensionAttrDetail(utils.PathSearch("attr", val, nil)),
			"detail_attrs":    flattenDimensionDetailAttrs(utils.PathSearch("detail_attrs", val, make([]interface{}, 0)).([]interface{})),
			"attr_id":         utils.PathSearch("attr_id", val, nil),
			"attr_name_en":    utils.PathSearch("attr_name_en", val, nil),
			"detail_attr_ids": utils.PathSearch("detail_attr_ids", val, make([]interface{}, 0)),
			"hierarchies_id":  utils.PathSearch("hierarchies_id", val, nil),
			"id":              utils.PathSearch("id", val, nil),
			"level":           utils.PathSearch("level", val, nil),
		})
	}
	return rst
}

func flattenDimensionHierarchies(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", val, nil),
			"name":       utils.PathSearch("name", val, nil),
			"attrs":      flattenDimensionHierarchyAttrs(utils.PathSearch("attrs", val, make([]interface{}, 0)).([]interface{})),
			"create_by":  utils.PathSearch("create_by", val, nil),
			"created_at": utils.PathSearch("create_time", val, nil),
			"updated_by": utils.PathSearch("update_by", val, nil),
			"updated_at": utils.PathSearch("update_time", val, nil),
		})
	}
	return rst
}

func flattenDimensionMappingsSourceTablesJoinFields(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"field1_id":   utils.PathSearch("field1_id", val, nil),
			"field2_id":   utils.PathSearch("field2_id", val, nil),
			"field1_name": utils.PathSearch("field1_name", val, nil),
			"field2_name": utils.PathSearch("field2_name", val, nil),
		})
	}
	return rst
}

func flattenDimensionMappingsSourceTables(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"join_type":   utils.PathSearch("join_type", val, nil),
			"table1_id":   utils.PathSearch("table1_id", val, nil),
			"table2_id":   utils.PathSearch("table2_id", val, nil),
			"table1_name": utils.PathSearch("table1_name", val, nil),
			"table2_name": utils.PathSearch("table2_name", val, nil),
			"join_fields": flattenDimensionMappingsSourceTablesJoinFields(
				utils.PathSearch("join_fields", val, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenDimensionMappingsSourceFields(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"field_ids":            utils.PathSearch("field_ids", val, nil),
			"field_names":          utils.PathSearch("field_names", val, make([]interface{}, 0)),
			"target_field_name":    utils.PathSearch("target_field_name", val, nil),
			"target_field_id":      utils.PathSearch("target_field_id", val, nil),
			"transform_expression": utils.PathSearch("transform_expression", val, nil),
			"changed":              utils.PathSearch("changed", val, false),
		})
	}
	return rst
}

func flattenDimensionMappingsDetails(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"target_attr_name": utils.PathSearch("target_attr_name", val, nil),
			"id":               utils.PathSearch("id", val, nil),
			"mapping_id":       utils.PathSearch("mapping_id", val, nil),
			"remark":           utils.PathSearch("remark", val, nil),
			"src_attr_ids":     utils.PathSearch("src_attr_ids", val, nil),
			"src_table_ids":    utils.PathSearch("src_table_ids", val, nil),
			"target_attr_id":   utils.PathSearch("target_attr_id", val, nil),
			"create_by":        utils.PathSearch("create_by", val, nil),
			"create_time":      utils.PathSearch("create_time", val, nil),
			"update_by":        utils.PathSearch("update_by", val, nil),
			"update_time":      utils.PathSearch("update_time", val, nil),
		})
	}
	return rst
}

func flattenDimensionMappings(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"name":              utils.PathSearch("name", val, nil),
			"src_model_id":      utils.PathSearch("src_model_id", val, nil),
			"src_model_name":    utils.PathSearch("src_model_name", val, nil),
			"id":                utils.PathSearch("id", val, nil),
			"view_text":         utils.PathSearch("view_text", val, nil),
			"target_table_id":   utils.PathSearch("target_table_id", val, nil),
			"target_table_name": utils.PathSearch("target_table_name", val, nil),
			"description":       utils.PathSearch("description", val, nil),
			"created_at":        utils.PathSearch("create_time", val, nil),
			"updated_at":        utils.PathSearch("update_time", val, nil),
			"created_by":        utils.PathSearch("create_by", val, nil),
			"updated_by":        utils.PathSearch("update_by", val, nil),
			"source_tables":     flattenDimensionMappingsSourceTables(utils.PathSearch("source_tables", val, make([]interface{}, 0)).([]interface{})),
			"source_fields":     flattenDimensionMappingsSourceFields(utils.PathSearch("source_fields", val, make([]interface{}, 0)).([]interface{})),
			"details":           flattenDimensionMappingsDetails(utils.PathSearch("details", val, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenDimensionModel(raw interface{}) []map[string]interface{} {
	if raw == nil {
		return nil
	}

	rst := []map[string]interface{}{{
		"name": utils.PathSearch("name", raw, nil),
		"type": utils.PathSearch("type", raw, nil),
		"id":   utils.PathSearch("id", raw, nil),
	}}
	return rst
}

func flattenDimensionSelfDefinedFields(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"fd_name_ch": utils.PathSearch("fd_name_ch", val, nil),
			"fd_name_en": utils.PathSearch("fd_name_en", val, nil),
			"fd_value":   utils.PathSearch("fd_value", val, nil),
			"not_null":   utils.PathSearch("not_null", val, false),
		})
	}
	return rst
}

func flattenDimensionBizInfoObj(raw interface{}) []map[string]interface{} {
	if raw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"biz_status":     utils.PathSearch("biz_status", raw, nil),
			"biz_type":       utils.PathSearch("biz_type", raw, nil),
			"biz_version":    utils.PathSearch("biz_version", raw, nil),
			"create_by":      utils.PathSearch("create_by", raw, nil),
			"directory_path": utils.PathSearch("directory_path", raw, nil),
			"email":          utils.PathSearch("email", raw, nil),
			"id":             utils.PathSearch("id", raw, nil),
			"l1":             utils.PathSearch("l1", raw, nil),
			"l2":             utils.PathSearch("l2", raw, nil),
			"l3":             utils.PathSearch("l3", raw, nil),
			"msg":            utils.PathSearch("msg", raw, nil),
			"name_ch":        utils.PathSearch("name_ch", raw, nil),
			"name_en":        utils.PathSearch("name_en", raw, nil),
			"submit_time":    utils.PathSearch("submit_time", raw, nil),
			"tenant_id":      utils.PathSearch("tenant_id", raw, nil),
		},
	}
}

func flattenDimensionApprovalInfo(raw interface{}) []map[string]interface{} {
	if raw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"approval_status": utils.PathSearch("approval_status", raw, nil),
			"approval_time":   utils.PathSearch("approval_time", raw, nil),
			"approval_type":   utils.PathSearch("approval_type", raw, nil),
			"approver":        utils.PathSearch("approver", raw, nil),
			"biz_id":          utils.PathSearch("biz_id", raw, nil),
			"biz_info":        utils.PathSearch("biz_info", raw, nil),
			"id":              utils.PathSearch("id", raw, nil),
			"msg":             utils.PathSearch("msg", raw, nil),
			"name_ch":         utils.PathSearch("name_ch", raw, nil),
			"name_en":         utils.PathSearch("name_en", raw, nil),
			"submit_time":     utils.PathSearch("submit_time", raw, nil),
			"tenant_id":       utils.PathSearch("tenant_id", raw, nil),
			"directory_path":  utils.PathSearch("directory_path", raw, nil),
			"email":           utils.PathSearch("email", raw, nil),
			"biz_info_obj":    flattenDimensionBizInfoObj(utils.PathSearch("biz_info_obj", raw, nil)),
		},
	}
}

func flattenDimensionCodeTableFieldValues(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"id":      utils.PathSearch("id", val, nil),
			"fd_id":   utils.PathSearch("fd_id", val, nil),
			"fd_value": utils.PathSearch("fd_value", val, nil),
			"ordinal": utils.PathSearch("ordinal", val, nil),
		})
	}
	return rst
}

func flattenDimensionCodeTableFields(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, val := range rawArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", val, nil),
			"code_table_id":     utils.PathSearch("code_table_id", val, nil),
			"name_ch":           utils.PathSearch("name_ch", val, nil),
			"name_en":           utils.PathSearch("name_en", val, nil),
			"ordinal":           utils.PathSearch("ordinal", val, nil),
			"data_type":         utils.PathSearch("data_type", val, nil),
			"data_type_extend":  utils.PathSearch("data_type_extend", val, nil),
			"description":       utils.PathSearch("description", val, nil),
			"domain_type":       utils.PathSearch("domain_type", val, nil),
			"is_unique_key":     utils.PathSearch("is_unique_key", val, false),
			"count_field_values": utils.PathSearch("count_field_values", val, nil),
			"code_table_field_values": flattenDimensionCodeTableFieldValues(
				utils.PathSearch("code_table_field_values", val, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenDimensionCodeTable(raw interface{}) []map[string]interface{} {
	if raw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":                utils.PathSearch("id", raw, nil),
			"name_ch":           utils.PathSearch("name_ch", raw, nil),
			"name_en":           utils.PathSearch("name_en", raw, nil),
			"status":            utils.PathSearch("status", raw, nil),
			"tb_version":        utils.PathSearch("tb_version", raw, nil),
			"create_by":         utils.PathSearch("create_by", raw, nil),
			"create_time":       utils.PathSearch("create_time", raw, nil),
			"directory_path":    utils.PathSearch("directory_path", raw, nil),
			"description":       utils.PathSearch("description", raw, nil),
			"directory_id":      utils.PathSearch("directory_id", raw, nil),
			"update_time":       utils.PathSearch("update_time", raw, nil),
			"new_biz":           flattenDimensionNewBiz(utils.PathSearch("new_biz", raw, nil)),
			"approval_info":     flattenDimensionApprovalInfo(utils.PathSearch("approval_info", raw, nil)),
			"code_table_fields": flattenDimensionCodeTableFields(utils.PathSearch("code_table_fields", raw, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenDimensionNewBiz(raw interface{}) []map[string]interface{} {
	if raw == nil {
		return nil
	}

	rst := []map[string]interface{}{{
		"id":       utils.PathSearch("id", raw, nil),
		"biz_type": utils.PathSearch("biz_type", raw, nil),
	}}
	return rst
}

func GetArchitectureDimensionById(client *golangsdk.ServiceClient, workspaceId, dimensionId string) (interface{}, error) {
	httpURL := "v2/{project_id}/design/dimensions/{id}?latest=true"

	getPath := client.Endpoint + httpURL
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", dimensionId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(workspaceId),
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	dimension := utils.PathSearch("data.value", respBody, make(map[string]interface{})).(map[string]interface{})
	if len(dimension) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/design/dimensions/{id}?latest=true",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the architecture dimension (%s) does not exist", dimensionId)),
			},
		}
	}

	return dimension, nil
}

func resourceArchitectureDimensionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		id          = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	dimension, err := GetArchitectureDimensionById(client, workspaceId, id)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", dimensionErrCodes...),
			"error retrieving DataArts Architecture dimension")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", d.Get("workspace_id").(string)),
		d.Set("name_ch", utils.PathSearch("name_ch", dimension, nil)),
		d.Set("name_en", utils.PathSearch("name_en", dimension, nil)),
		d.Set("dimension_type", utils.PathSearch("dimension_type", dimension, nil)),
		d.Set("l3_id", utils.PathSearch("l3_id", dimension, nil)),
		d.Set("owner", utils.PathSearch("owner", dimension, nil)),
		d.Set("attributes", flattenDimensionAttributes(utils.PathSearch("attributes", dimension, make([]interface{}, 0)).([]interface{}))),
		d.Set("datasource", flattenDimensionDatasource(utils.PathSearch("datasource", dimension, nil))),
		d.Set("alias", utils.PathSearch("alias", dimension, nil)),
		d.Set("code_table_id", utils.PathSearch("code_table_id", dimension, nil)),
		d.Set("configs", utils.PathSearch("configs", dimension, nil)),
		d.Set("create_by", utils.PathSearch("create_by", dimension, nil)),
		d.Set("description", utils.PathSearch("description", dimension, nil)),
		d.Set("distribute", utils.PathSearch("distribute", dimension, nil)),
		d.Set("distribute_column", utils.PathSearch("distribute_column", dimension, nil)),
		d.Set("hierarchies", flattenDimensionHierarchies(utils.PathSearch("hierarchies", dimension, make([]interface{}, 0)).([]interface{}))),
		d.Set("id_field", utils.PathSearch("id_field", dimension, nil)),
		d.Set("l1", utils.PathSearch("l1", dimension, nil)),
		d.Set("l2", utils.PathSearch("l2", dimension, nil)),
		d.Set("l2_id", utils.PathSearch("l2_id", dimension, nil)),
		d.Set("l3", utils.PathSearch("l3", dimension, nil)),
		d.Set("mappings", flattenDimensionMappings(utils.PathSearch("mappings", dimension, make([]interface{}, 0)).([]interface{}))),
		d.Set("model", flattenDimensionModel(utils.PathSearch("model", dimension, nil))),
		d.Set("model_id", utils.PathSearch("model_id", dimension, nil)),
		d.Set("obs_location", utils.PathSearch("obs_location", dimension, nil)),
		d.Set("self_defined_fields", flattenDimensionSelfDefinedFields(
			utils.PathSearch("self_defined_fields", dimension, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", utils.PathSearch("status", dimension, nil)),
		d.Set("table_type", utils.PathSearch("table_type", dimension, nil)),
		d.Set("update_by", utils.PathSearch("update_by", dimension, nil)),
		d.Set("approval_info", flattenDimensionApprovalInfo(utils.PathSearch("approval_info", dimension, nil))),
		d.Set("code_table", flattenDimensionCodeTable(utils.PathSearch("code_table", dimension, nil))),
		d.Set("create_time", utils.PathSearch("create_time", dimension, nil)),
		d.Set("dev_version", utils.PathSearch("dev_version", dimension, nil)),
		d.Set("dev_version_name", utils.PathSearch("dev_version_name", dimension, nil)),
		d.Set("env_type", utils.PathSearch("env_type", dimension, nil)),
		d.Set("l1_id", utils.PathSearch("l1_id", dimension, nil)),
		d.Set("new_biz", flattenDimensionNewBiz(utils.PathSearch("new_biz", dimension, nil))),
		d.Set("prod_version", utils.PathSearch("prod_version", dimension, nil)),
		d.Set("prod_version_name", utils.PathSearch("prod_version_name", dimension, nil)),
		d.Set("update_time", utils.PathSearch("update_time", dimension, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceArchitectureDimensionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/design/dimensions"
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateBody := utils.RemoveNil(buildCreateOrUpdateDimensionBodyParams(d))
	updateBody["id"] = d.Id()
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody:         updateBody,
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceArchitectureDimensionRead(ctx, d, meta)
}

func resourceArchitectureDimensionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2/{project_id}/design/dimensions"
	)

	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody: 		  utils.RemoveNil(map[string]interface{}{
			"ids": []string{d.Id()},
		}),
	}

	resp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if utils.PathSearch("data.value", respBody, 0) != 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error deleting DataArts Architecture dimension")
	}

	return nil
}
