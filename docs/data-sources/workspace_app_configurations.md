---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_configurations"
description: |-
  Use this data source to get enterprise system configuration list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_configurations

Use this data source to get enterprise system configuration list of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
variable "items" {
  type = list(string)
}

data "huaweicloud_workspace_app_configurations" "test" {
  items = var.items
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the Workspace APP are located.  
  If omitted, the provider-level region will be used.

* `items` - (Required, List) Specifies the list of configuration keys to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - The list of configuration information.  
  The [configurations](#data_attr_app_configurations) structure is documented below.

<a name="data_attr_app_configurations"></a>
The `configurations` block supports:

* `config_key` - The key of the configuration.

* `config_value` - The value corresponding to the configuration key.
