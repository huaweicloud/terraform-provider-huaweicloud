---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_vulnerabilities"
description: |-
  Use this data source to query the SecMaster vulnerabilities within HuaweiCloud.
---

# huaweicloud_secmaster_vulnerabilities

Use this data source to query the SecMaster vulnerabilities within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "from_date" {}
variable "to_date" {}

data "huaweicloud_secmaster_vulnerabilities" "test" {
  workspace_id = var.workspace_id
  from_date    = var.from_date
  to_date      = var.to_date
  sort_by      = "create_time"
  order        = "DESC"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the vulnerabilities.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to query vulnerabilities.

* `sort_by` - (Optional, String) Specifies the field to sort the results by.
  Valid values are **create_time** and **update_time**.

* `order` - (Optional, String) Specifies the sorting direction. The value can be **DESC** or **ASC**.

* `from_date` - (Optional, String) Specifies the start time of the time range to query.
  For example **2023-02-20T00:00:00.000Z**.

* `to_date` - (Optional, String) Specifies the end time of the time range to query.
  For example **2023-02-27T23:59:59.999Z**.

* `condition` - (Optional, List) Specifies the query conditions.
  The [condition](#condition_struct) structure is documented below.

<a name="condition_struct"></a>
The `condition` block supports:

* `conditions` - (Optional, List) Specifies the query conditions.
  The [conditions](#conditions_struct) structure is documented below.

* `logics` - (Optional, List) Specifies the logic conditions.

<a name="conditions_struct"></a>
The `conditions` block supports:

* `name` - (Optional, String) Specifies the field name for the condition.

* `data` - (Optional, List) Specifies the field value list for the condition.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of vulnerabilities.
  The [data](#secmaster_vulnerability_data) structure is documented below.

<a name="secmaster_vulnerability_data"></a>
The `data` block supports:

* `id` - The vulnerability ID.

* `format_version` - The format version of the vulnerability.

* `version` - The version of the vulnerability.

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `create_time` - The creation time of the vulnerability.

* `update_time` - The last update time of the vulnerability.

* `dataclass_ref` - The data class reference.
  The [dataclass_ref](#secmaster_vulnerability_dataclass_ref) structure is documented below.

* `data_object` - The vulnerability details.
  The [data_object](#secmaster_vulnerability_data_object) structure is documented below.

<a name="secmaster_vulnerability_dataclass_ref"></a>
The `dataclass_ref` block supports:

* `id` - The data class ID.

* `name` - The data class name.

<a name="secmaster_vulnerability_data_object"></a>
The `data_object` block supports:

* `vul_name` - The vulnerability name.

* `first_observed_time` - The first observed time of the vulnerability.

* `batch_number` - The batch number of the vulnerability.

* `description` - The description of the vulnerability.

* `resource_num` - The number of affected resources.

* `domain_id` - The domain ID.

* `workspace_id` - The workspace ID.

* `remediation` - The remediation of the vulnerability.
  The [remediation](#secmaster_vulnerability_remediation) structure is documented below.

* `domain_name` - The domain name.

* `update_time` - The last update time of the vulnerability.

* `is_deleted` - Whether the vulnerability is deleted.

* `project_id` - The project ID.

* `extend_properties` - The extended properties.
  The [extend_properties](#secmaster_vulnerability_extend_properties) structure is documented below.

* `region_name` - The region name.

* `id` - The vulnerability notice ID.

* `vulnerability_type` - The vulnerability type.
  The [vulnerability_type](#secmaster_vulnerability_type) structure is documented below.

* `create_time` - The creation time of the vulnerability.

* `last_observed_time` - The last observed time of the vulnerability.

* `resource` - The affected resource.
  The [resource](#secmaster_vulnerability_resource) structure is documented below.

* `count` - The count of vulnerabilities.

* `region_id` - The region ID.

* `vulnerability` - The vulnerability details.
  The [vulnerability](#secmaster_vulnerability_vulnerability) structure is documented below.

* `dataclass_id` - The data class ID.

* `version` - The version of the vulnerability.

* `data_source` - The data source.
  The [data_source](#secmaster_vulnerability_data_source) structure is documented below.

* `arrive_time` - The arrive time of the vulnerability.

* `environment` - The environment information.
  The [environment](#secmaster_vulnerability_environment) structure is documented below.

* `trigger_flag` - The trigger flag of the vulnerability.

* `handled` - The vulnerability handled status.

<a name="secmaster_vulnerability_remediation"></a>
The `remediation` block supports:

* `recommendation` - The recommendation for the vulnerability type.

<a name="secmaster_vulnerability_extend_properties"></a>
The `extend_properties` block supports:

* `operations` - The operations.
  The [operations](#secmaster_vulnerability_operations) structure is documented below.

<a name="secmaster_vulnerability_operations"></a>
The `operations` block supports:

* `is_build_in` - Whether the operation is built-in. Valid values are **true** and **false**.

<a name="secmaster_vulnerability_type"></a>
The `vulnerability_type` block supports:

* `id` - The vulnerability type ID.

* `category` - The category of the vulnerability type.

* `category_en` - The English category of the vulnerability type.

* `category_zh` - The Chinese category of the vulnerability type.

* `vulnerability_type` - The vulnerability type.

* `vulnerability_type_en` - The English vulnerability type.

* `vulnerability_type_zh` - The Chinese vulnerability type.

<a name="secmaster_vulnerability_resource"></a>
The `resource` block supports:

* `id` - The resource ID.

* `name` - The resource name.

* `type` - The resource type.

* `provider` - The resource provider.

* `region_id` - The resource region ID.

* `domain_id` - The resource domain ID.

* `project_id` - The resource project ID.

* `ep_id` - The resource enterprise project ID.

* `tags` - The resource tags.

<a name="secmaster_vulnerability_vulnerability"></a>
The `vulnerability` block supports:

* `id` - The vulnerability ID.

* `type` - The vulnerability type.

* `url` - The vulnerability URL.

* `status` - The vulnerability status.

* `level` - The vulnerability level.

* `reason` - The vulnerability reason.

* `solution` - The vulnerability solution.

* `repair_severity` - The vulnerability repair severity.

* `related` - The related vulnerabilities.

* `tags` - The tags of the vulnerability.

<a name="secmaster_vulnerability_data_source"></a>
The `data_source` block supports:

* `domain_id` - The domain ID.

* `project_id` - The project ID.

* `region_id` - The region ID.

* `company_name` - The company name.

* `source_type` - The source type.

* `product_name` - The product name.

* `product_feature` - The product feature.

<a name="secmaster_vulnerability_environment"></a>
The `environment` block supports:

* `domain_id` - The domain ID.

* `domain_name` - The domain name.

* `project_id` - The project ID.

* `region_id` - The region ID.

* `region_name` - The region name.

* `vendor_type` - The vendor type.
