---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_pt_applicable_instances"
description: |-
  Use this data source to get the list of instances that a GeminiDB parameter template can be applied to.
---

# huaweicloud_geminidb_pt_applicable_instances

Use this data source to get the list of instances that a GeminiDB parameter template can be applied to.

## Example Usage

### Basic Usage

```hcl
variable "config_id" {}

data "huaweicloud_geminidb_pt_applicable_instances" "test" {
  config_id = var.config_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the applicable instances.
  If omitted, the provider-level region will be used.

* `config_id` - (Required, String) Specifies the ID of the parameter template.

* `instance_name` - (Optional, String) Specifies the instance name to filter.

* `instance_id` - (Optional, String) Specifies the instance ID to filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of instances that the parameter template can be applied to.
  The [instances](#geminidb_applicable_instances_instances) structure is documented below.

<a name="geminidb_applicable_instances_instances"></a>
The `instances` block supports:

* `id` - The instance ID.

* `name` - The instance name.
