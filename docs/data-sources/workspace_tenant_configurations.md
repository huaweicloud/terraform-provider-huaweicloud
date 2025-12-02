---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_tenant_configurations"
description: |-
  Use this data source to get the list of Workspace tenant configurations.
---

# huaweicloud_workspace_tenant_configurations

Use this data source to get the list of Workspace tenant configurations.

## Example Usage

```hcl
variable "configuration_name" {}

data "huaweicloud_workspace_tenant_configurations" "test" {
  configuration_name = var.configuration_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the tenant configuration are located.  
  If omitted, the provider-level region will be used.

* `configuration_name` - (Optional, String) Specifies the name of the configuration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - The list of tenant configurations that matched filter parameters.  
  The [configurations](#workspace_configurations_attr) structure is documented below.

<a name="workspace_configurations_attr"></a>
The `configurations` block supports:

* `id` - The ID of the configuration.

* `name` - The name of the configuration.

* `status` - The status of the configuration.

* `values` - The configuration values.  
  The [values](#workspace_configurations_values_attr) structure is documented below.

<a name="workspace_configurations_values_attr"></a>
The `values` block supports:

* `key` - The key of the configuration item.

* `value` - The value of the configuration item.
