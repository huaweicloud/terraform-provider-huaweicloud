---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_subscription_orders"
description: |-
  Use this data source to get the list of subscription orders.
---

# huaweicloud_secmaster_subscription_orders

Use this data source to get the list of subscription orders.

## Example Usage

```hcl
data "huaweicloud_secmaster_subscription_orders" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `page` - (Optional, String) Specifies the order resource detail information enumeration.  
  The valid values are as follows:
  + **DEFAULT**: Get the list of opened resources, excluding package packages
  + **PURCHASE**: Return the number of ECS under the tenant's name on the basis of DEFAULT.
  + **RESOURCE_LIST**: Return the package list based on DEFAULT.
  + **SMN**: Return the list of subscribed SMN topics.

  Defaults to **DEFAULT**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `csb_version` - The current version information of the tenant.

* `ecs_count` - The number of ECS instances under the current project of the tenant.

* `resources` - The resource list.

  The [resources](#resources_struct) structure is documented below.

* `subscription_count` - The number of topic subscriptions.

* `subscriptions` - The tenant subscription information.

  The [subscriptions](#subscriptions_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_type_name` - The resource name.

* `resource_size` - The resource specification.

* `cloud_service` - The order source.

* `resource_type` - The resource type.

* `resource_spec_code` - The resource specification code.

* `to_period` - Whether the current resource can be converted from pay-per-use to yearly/monthly.

* `create_time` - The creation timestamp of the resource.

* `update_time` - The update timestamp of the resource.

* `expire_time` - The expiration timestamp of the resource.

* `resource_status` - The resource status.

* `order_id` - The order ID.

* `charging_mode` - The charging mode.

* `quota_reset_mode` - The reset mode.

* `quota_reset_cycle_type` - The reset cycle type.

* `quota_reset_cycle` - The reset cycle.

* `amount` - The remaining amount of the resource package.

* `original_amount` - The original usage of the resource package.

* `measure_name` - The measurement unit name of the resource package.

* `tag_list` - The tag list.

  The [tag_list](#tag_list_struct) structure is documented below.

* `usages` - The resource usage list.

  The [usages](#usages_struct) structure is documented below.

<a name="tag_list_struct"></a>
The `tag_list` block supports:

* `key` - The tag key.

* `value` - The tag value.

* `create_time` - The tag creation timestamp.

* `update_time` - The tag update timestamp.

<a name="usages_struct"></a>
The `usages` block supports:

* `unit` - The usage unit.

* `resource_type_name` - The resource type name.

* `source_resource_spec_code` - The source resource specification code.

* `resource_spec_code` - The resource specification code.

* `source_type` - The source resource type code.

* `used_percent` - The usage percentage.

* `quota` - The total quota.

* `used` - The used amount.

* `free` - The remaining amount.

<a name="subscriptions_struct"></a>
The `subscriptions` block supports:

* `owner` - The project ID of the tenant.

* `endpoint` - The subscription endpoint.

* `protocol` - The endpoint protocol.

* `subscription_urn` - The SMN subscription URN.

* `topic_urn` - The URN of the subscribed topic.

* `status` - The subscription status.
