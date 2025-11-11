---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_instances_by_tags"
description: |-
  Use this data source to query CBH instances by tags within HuaweiCloud.
---

# huaweicloud_cbh_instances_by_tags

Use this data source to query CBH instances by tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbh_instances_by_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the CBH instances.
  If omitted, the provider-level region will be used.

* `without_any_tag` - (Optional, Bool) Specifies whether to query all resources without tags.
  If this parameter is set to **true**, all resources without tags are queried.
  In this case, the `tags`, `tags_any`, `not_tags`, and `not_tags_any` fields are ignored.

* `tags` - (Optional, List) Specifies the tags value. The resources to be queried contain tags listed in `tags`.
  Each resource to be queried contains a maximum of `50` keys. Each tag key can have a maximum of `10` tag values.
  The tag value corresponding to each tag key can be an empty array but the structure cannot be missing.
  The key must be unique, and the values of the same key must be unique. Resources containing all tags are returned.
  Keys in this list are in AND relationship. While values in the key-value pairs are in OR relationship.
  If no tag filtering condition is specified, full data is returned.

  The [tags](#cbh_instances_tags) structure is documented below.

* `tags_any` - (Optional, List) Specifies any tags value. The resources to be queried contain any tags listed in `tags_any`.
  Each resource to be queried contains a maximum of `50` keys. Each tag key can have a maximum of `10` tag values.
  The tag value corresponding to each tag key can be an empty array but the structure cannot be missing.
  Each tag key must be unique, and each value of the same key must be unique.
  The response returns resources containing the tags in this list. Keys in this list are in an OR relationship and values
  in each key-value structure are also in an OR relationship.
  If no tag filtering condition is specified, full data is returned.

  The [tags_any](#cbh_instances_tags) structure is documented below.

* `not_tags` - (Optional, List) Specifies the not tags value. The resources to be queried do not contain tags listed in
  `not_tags`. Each resource to be queried contains a maximum of `50` keys. Each tag key can have a maximum of `10` tag
  values. The tag value corresponding to each tag key can be an empty array but the structure cannot be missing.
  Each tag key must be unique, and each value of the same key must be unique.
  The response returns resources containing no tags in this list. Keys in this list are in an AND relationship while
  values in each key-value structure are in an OR relationship.
  If no filtering condition is specified, full data is returned.

  The [not_tags](#cbh_instances_tags) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies the not contained any tags value. The resources to be queried do not
  contain any tags listed in `not_tags_any`. Each resource to be queried contains a maximum of `50` keys. Each tag key
  can have a maximum of `10` tag values. The tag value corresponding to each tag key can be an empty array but the
  structure cannot be missing. Each tag key must be unique, and each value of the same key must be unique.
  The response returns resources containing no tag in this list. Keys in this list are in the OR relationship and values
  in each key-value structure are also in the OR relationship.
  If no tag filtering criteria is specified, full data is returned.

  The [not_tags_any](#cbh_instances_tags) structure is documented below.

* `sys_tags` - (Optional, List) Specifies the system tags value.
  Only users with the op_service permission can use this field to filter resources.
  Only one tag structure is contained when this API is called by Tag Management Service (TMS).
  The key: **_sys_enterprise_project_id** and the value: Enterprise project ID list.
  Currently, key contains only one value. `0` indicates the default enterprise project.
  Field `sys_tags` and tenant tag filtering conditions (`without_any_tag`, `tags`, `tags_any`, `not_tags`, and
  `not_tags_any`) cannot be used at the same time.
  If `sys_tags` is not specified, filter the resources with other tag filtering criteria. If no tag filtering criteria
  is specified, full data is returned.

  The [sys_tags](#cbh_instances_tags) structure is documented below.

* `matches` - (Optional, List) Specifies the match values.
  The tag key is the field to be matched, for example, **resource_name**.
  The value is a matched value. The key is a fixed dictionary value and cannot be a duplicate key or unsupported key.
  Check whether fuzzy match is required based on the key value. For example, if key is set to resource_name,
  fuzzy search (case-insensitive) is performed by default. If value is empty, exact match is performed. Most services
  do not have resources without names. In this case, an empty list is returned.
  Field `resource_id` indicates exact match. Only **resource_name** is used now.

  The [matches](#cbh_instances_matches) structure is documented below.

<a name="cbh_instances_tags"></a>
The `tags`, `tags_any`, `not_tags`, `not_tags_any`, and `sys_tags` blocks support:

* `key` - (Required, String) The key of the tag.

* `values` - (Required, List) The list of values for the tag.

<a name="cbh_instances_matches"></a>
The `matches` block supports:

* `key` - (Required, String) The field to match against.

* `value` - (Required, String) The value to match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_count` - The total number of records that match the filter criteria.

* `resources` - The list of CBH instances that match the filter parameters.

  The [resources](#cbh_instances_resources) structure is documented below.

<a name="cbh_instances_resources"></a>
The `resources` block supports:

* `resource_id` - The ID of the resource.

* `resource_name` - The name of the resource.

* `resource_detail` - The detailed information of the instance.
  The [resource_detail](#cbh_instances_resource_detail) structure is documented below.

* `tags` - The list of tags associated with the instance.
  The [tags](#cbh_instances_resource_tags) structure is documented below.

* `sys_tags` - The list of system tags associated with the instance.
  The [sys_tags](#cbh_instances_resource_tags) structure is documented below.

<a name="cbh_instances_resource_detail"></a>
The `resource_detail` block supports:

* `name` - The CBH instance name.

* `server_id` - The ID of the server where the CBH instance is deployed.

* `instance_id` - The CBH instance ID.

* `alter_permit` - Whether the CBH instance can be expanded.

* `enterprise_project_id` - The ID of the enterprise project.

* `period_num` - The number of subscription periods of a CBH instance.

* `start_time` - The start time of the CBH instance in timestamp format.

* `end_time` - The end time of the CBH instance in timestamp format.

* `created_time` - The creation time of the CBH instance in UTC format.

* `upgrade_time` - The upgrade schedule of the CBH instance, in the format of timestamp.

* `update` - Whether the CBH instance can be upgraded. Valid values are:
  + **OLD**: The current version is the latest one.
  + **NEW**: Can be upgraded.
  + **CROSS_OS**: Can be upgraded to any later versions.
  + **ROLLBACK**: Can be rolled back.

* `bastion_version` - The current version of the CBH instance.

* `az_info` - The availability zone information of the instance.
  The [az_info](#cbh_instances_az_info) structure is documented below.

* `status_info` - The status information of the instance.
  The [status_info](#cbh_instances_status_info) structure is documented below.

* `resource_info` - The resource information of the instance.
  The [resource_info](#cbh_instances_resource_info) structure is documented below.

* `network` - The network information of the instance.
  The [network](#cbh_instances_network) structure is documented below.

* `ha_info` - The high availability information of the instance.
  The [ha_info](#cbh_instances_ha_info) structure is documented below.

<a name="cbh_instances_az_info"></a>
The `az_info` block supports:

* `region` - The ID of the AZ where the CBH instance locates.

* `zone` - The ID of the AZ where the CBH instance locates. (In primary/standby mode, the ID of the AZ where the primary
  instance locates is required.)

* `availability_zone_display` - The AZ where the CBH instance locates. (In primary/standby mode, the ID of the AZ where
  the primary instance locates is required.)

* `slave_zone` - The AZ where the standby CBH instance locates.

* `slave_zone_display` - The AZ where the standby CBH instance locates.

<a name="cbh_instances_status_info"></a>
The `status_info` block supports:

* `status` - The status of the instance. Valid values are:
  + **SHOUTOFF**: Closed
  + **ACTIVE**: Running
  + **DELETING**: Deleting
  + **BUILD**: Creating
  + **DELETED**: Deleting
  + **ERROR**: Faulty
  + **HAWAIT**: Waiting for the standby node to be created
  + **FROZEN**: Frozen
  + **UPGRADING**: Upgrading
  + **UNPAID**: Pending Payment
  + **RESIZE**: Changing specifications
  + **DILATATION**: Expanding capacity
  + **HA**: Configuring HA

* `task_status` - The task status of the instance. Valid valus are:
  + **powering-on**: Started
  + **powering-off**: Stopped
  + **rebooting**: Rebooting
  + **delete_wait**: Deleting
  + **frozen**: Frozen
  + **NO_TASK**: Running
  + **unfrozen**: Unfrozen
  + **alter**: Changing
  + **updating**: Upgrading
  + **configuring-ha**: Configuring HA
  + **data-migrating**: Migrating data
  + **rollback**: Rolling back to the previous version.
  + **traffic-switchover**: Traffic switching

* `create_instance_status` - The creation status of the instance. Valid valus are:
  + **waiting-for-payment**: Waiting for payment
  + **creating-network**: Creating a network.
  + **creating-server**: Creating the service.
  + **tranfering-horizontal-network**: Establishing network connections.
  + **adding-policy-route**: Adding a routing policy.
  + **configing-dns**: Configuring DNS.
  + **starting-cbs-service**: The service is running.
  + **setting-init-conf**: Initializing
  + **buying-EIP**: Buying an EIP.

* `instance_status` - The status of the instance. Valid values are:
  + **building**: Creating
  + **deleting**: Deleting
  + **deleted**: Deleted
  + **unpaid**: Unpaid
  + **upgrading**: Upgrading
  + **resizing**: Resizing
  + **abnormal**: Abnormal
  + **error**: Faulty
  + **ok**: Normal

* `instance_description` - The description of the instance status.

* `fail_reason` - The failure reason if the instance creation fails.

<a name="cbh_instances_resource_info"></a>
The `resource_info` block supports:

* `specification` - The specification of the instance.

* `order_id` - The order ID.

* `resource_id` - The resource ID.

* `data_disk_size` - The size of the data disk.

* `disk_resource_id` - The list of disk resource IDs.

<a name="cbh_instances_network"></a>
The `network` block supports:

* `vip` - The floating IP address of the CBH instance. (This field is returned when the instance is deployed in
  primary/standby mode.)

* `web_port` - The port used for accessing the CBH instance with a web browser.

* `public_ip` - The EIP bound to the CBH instance.

* `public_id` - The ID of the EIP bound to the CBH instance, in the UUID format.

* `private_ip` - The private IP address of the CBH instance.

* `vpc_id` - The ID of the VPC where the CBH instance locates.

* `subnet_id` - The ID of the subnet where the CBH instance locates.

* `security_group_id` - The ID of the Security group where the CBH instance locates.

<a name="cbh_instances_ha_info"></a>
The `ha_info` block supports:

* `ha_id` - The IDs of the primary and standby instances.

* `instance_type` - The type of the instance (**master** or **slave**).

<a name="cbh_instances_resource_tags"></a>
The `tags` and `sys_tags` blocks in resources support:

* `key` - The key of the tag.

* `value` - The value of the tag.
