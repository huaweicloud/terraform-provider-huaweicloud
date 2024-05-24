---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resources"
description: |-
  Use this data source to get the list of RMS resources.
---

# huaweicloud_rms_resources

Use this data source to get the list of RMS resources.

## Example Usage

```hcl
data "huaweicloud_rms_resources" "test" {
  type = "vpc.vpcs"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the resource name.

* `resource_id` - (Optional, String) Specifies the resource ID.

* `type` - (Optional, String) Specifies the resource type. For example, **vpc.vpcs** and **rds.instances**.

* `region_id` - (Optional, String) Specifies the region to which the resource belongs.

* `enterprise_project_id` - (Optional, String) Specifies the ID of enterprise project to which the resource belongs.

* `tracked` - (Optional, Bool) Specifies whether the resource is tracked.

* `resource_deleted` - (Optional, Bool) Specifies whether the query the deleted resources.

* `tags` - (Optional, Map) Specifies the tags to filter the resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The service details list.

  The [resources](#resources) structure is documented below.

<a name="resources"></a>
The `resources` block supports:

* `id` - The resource ID.

* `name` - The resource name.

* `service` - The service name.

* `type` - The resource type.

* `region_id` - The region to which the resource belongs.

* `project_id` - The ID of project to which the resource belongs.

* `project_name` - The name of project to which the resource belongs.

* `enterprise_project_id` - The ID of enterprise project to which the resource belongs.

* `enterprise_project_name` - The name of enterprise project to which the resource belongs.

* `checksum` - The checksum of the resource.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `provisioning_state` - The provisioning state of the resource.

* `state` - The state of the resource.

* `tags` - The tags of the resource.

* `properties` - The properties of the resource.
