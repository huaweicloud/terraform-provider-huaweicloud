---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_resources"
description: |-
  Use this data source to get the list of resources within HuaweiCloud Resource Center.
---

# huaweicloud_rc_resources

Use this data source to get the list of resources within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_resources" "test" {}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `region_id` - (Optional, String) Specifies the region ID.

* `ep_id` - (Optional, String) Specifies the enterprise project ID.

* `type` - (Optional, String) Specifies the resource type (provider.type).

* `resource_id` - (Optional, String) Specifies the resource ID.

* `name` - (Optional, String) Specifies the resourceQL name.

* `tags` - (Optional, List) Specifies the tag list.

* `group_id` - (Optional, String) Specifies the resource group ID.

* `sort_key` - (Optional, List) Specifies the sort field list, only **createdAt** and **updatedAt** are supported.

* `sort_value` - (Optional, List) Specifies the sort order list, only **asc** and **desc** are supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `resources` - The list of resources.
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
