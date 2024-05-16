---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_service_groups"
description: |-
  Use this data source to get the list of CFW service groups.
---

# huaweicloud_cfw_service_groups

Use this data source to get the list of CFW service groups.

## Example Usage

```hcl
variable "object_id" {}

data "huaweicloud_cfw_service_groups" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID.

* `key_word` - (Optional, String) Specifies the keyword of the service group description.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.

* `name` - (Optional, String) Specifies the name of the service group.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `service_groups` - Service group list

  The [service_groups](#data_service_groups_struct) structure is documented below.

<a name="data_service_groups_struct"></a>
The `service_groups` block supports:

* `id` - The service group ID.

* `name` - The name of the service group.

* `type` - The type of the Service group.

* `ref_count` - The number of times this service group has been referenced.

* `description` - The description of the service group.

* `protocols` - The protocols of the service group.
