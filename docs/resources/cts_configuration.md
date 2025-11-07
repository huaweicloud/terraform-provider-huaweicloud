---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_configuration"
description: |-
  Using this resource to manage CTS configuration within HuaweiCloud.
---

# huaweicloud_cts_configuration

Using this resource to manage CTS configuration within HuaweiCloud.

-> This resource is only a configuration resource for managing CTS configuration. Deleting this resource will not
restore the CTS configuration to the default value, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_cts_configuration" "test" {
  is_sync_global_trace       = true
  is_support_read_only       = true
  support_read_only_services = ["ECS", "EVS"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the CTS configuration.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `is_sync_global_trace` - (Required, Bool) Specifies whether to synchronize global service logs from the central
  region.

* `is_support_read_only` - (Required, Bool) Specifies whether to enable the reporting of read-only audit logs for all
  cloud services.

* `support_read_only_services` - (Optional, List) Specifies the cloud services that enable read-only audit logs.  
  The value is a list of service names, such as **ECS**, **EVS**, **VPC**, etc.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
