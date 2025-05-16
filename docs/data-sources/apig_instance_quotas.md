---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_quotas"
description: |-
  Use this data source to get resource quota list under the specified APIG instance within HuaweiCloud.
---

# huaweicloud_apig_instance_quotas

Use this data source to get resource quota list under the specified APIG instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_apig_instance_quotas" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the APIG instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of the quotas.

  The [quotas](#data_instance_quotas_struct) structure is documented below.

<a name="data_instance_quotas_struct"></a>
The `quotas` block supports:

* `config_id` - The ID of the quota.

* `config_name` - The name of the quota.

* `config_value` - The number of available quotas.

* `used` - The number of quota used.

* `remark` - The description of the quota.

* `config_time` - The creation time of the quota.
