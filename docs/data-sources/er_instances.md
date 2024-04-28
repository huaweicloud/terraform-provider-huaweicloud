---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_instances"
description: ""
---

# huaweicloud_er_instances

Use this data source to filter ER instances within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_er_instances" "test" {
  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the ER instances are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID used to query specified ER instance.

* `name` - (Optional, String) Specifies the name used to filter the ER instances.
  The valid length is limited from `1` to `64`, only Chinese and English letters, digits, underscores (_) and
  hyphens (-) are allowed.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the ER instances to be queried.

* `owned_by_self` - (Optional, Bool) Specifies whether resources belong to the current renant.

* `status` - (Optional, String) Specifies the status used to filter the ER instances.

* `tags` - (Optional, Map) Specifies the key/value pairs used to filter the ER instances.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - All instances that match the filter parameters.  
  The [object](#er_data_instances) structure is documented below.

<a name="er_data_instances"></a>
The `instances` block supports:

* `id` - The ER instance ID.

* `asn` - The BGP AS number of the ER instance.

* `name` - The name of the ER instance.

* `description` - The description of the ER instance.

* `status` - The current status of the ER instance.

* `enterprise_project_id` - The ID of enterprise project to which the ER instance belongs.

* `tags` - The key/value pairs to associate with the ER instance.

* `created_at` - The creation time of the ER instance.

* `updated_at` - The last update time of the ER instance.

* `enable_default_propagation` - Whether to enable the propagation of the default route table.

* `enable_default_association` - Whether to enable the association of the default route table.

* `auto_accept_shared_attachments` - Whether to automatically accept the creation of shared attachment.

* `default_propagation_route_table_id` - The ID of the default propagation route table.

* `default_association_route_table_id` - The ID of the default association route table.

* `availability_zones` - The availability zone list where the ER instance is located.
