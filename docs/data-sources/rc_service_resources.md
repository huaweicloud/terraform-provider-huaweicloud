---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_service_resources"
description: |-
  Use this data source to get the resources of a specific service and type within HuaweiCloud Resource Center.
---

# huaweicloud_rc_service_resources

Use this data source to get the resources of a specific service and type within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_service_resources" "test" {
  rc_provider = "ecs"
  type        = "cloudservers"
}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `rc_provider` - (Required, String) Specifies the cloud service name.

* `type` - (Required, String) Specifies the esource type name.

* `region_id` - (Optional, String) Specifies the egion ID.

* `ep_id` - (Optional, String) Specifies the enterprise project ID.

* `tag` - (Optional, List) Specifies the tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `resources` - List of resources.
  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `id` - The resource ID.

* `name` - The resource name.

* `provider` - The cloud service name.

* `type` - The resource type.

* `region_id` - The region ID.

* `project_id` - The project ID in OpenStack.

* `project_name` - The project name in OpenStack.

* `ep_id` - The enterprise project ID.

* `ep_name` - The enterprise project name.

* `checksum` - The resource detail checksum.

* `created` - The resource creation time.

* `updated` - The resource update time.

* `provisioning_state` - The resource provisioning state.

* `tags` - The resource tags.

* `properties` - The resource detailed properties.

* `group_id` - The resource group ID.

* `state` - The resource state.
  + **Normal**
  + **Deleted**
