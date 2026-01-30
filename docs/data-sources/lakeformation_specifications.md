---
subcategory: "LakeFormation"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lakeformation_specifications"
description: |-
  Use this data source to get the list of LakeFormation instance specifications within HuaweiCloud.
---

# huaweicloud_lakeformation_specifications

Use this data source to get the list of LakeFormation instance specifications within HuaweiCloud.

## Example Usage

### Query all instance specifications

```hcl
data "huaweicloud_lakeformation_specifications" "test" {}
```

### Query the instance specifications using spec_code filter parameter

```hcl
data "huaweicloud_lakeformation_specifications" "test" {
  spec_code = "hws.resource.type.lakeformation.qps"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `spec_code` - (Optional, String) Specifies the code of the specification to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

`id` - The data source ID.

* `specifications` - The list of specifications that matched filter parameters.
  The [specifications](#lakeformation_specifications_attr) structure is documented below.

<a name="lakeformation_specifications_attr"></a>
The `specifications` block supports:

* `spec_code` - The code of the specification.

* `resource_type` - The resource type of the specification.

* `stride` - The stride of the specification.

* `unit` - The unit of the specification.

* `min_stride_num` - The minimum stride number of the specification.

* `max_stride_num` - The maximum stride number of the specification.

* `usage_measure_id` - The usage measurement unit ID of the specification.

* `usage_factor` - The usage factor of the specification.

* `usage_value` - The usage value of the specification.

* `free_usage_value` - The free usage value of the specification.

* `stride_num_whitelist` - The stride number whitelist of the specification.
