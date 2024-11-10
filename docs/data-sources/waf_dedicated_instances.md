---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_dedicated_instances"
description: |-
  Use this data source to get a list of WAF dedicated instances.
---

# huaweicloud_waf_dedicated_instances

Use this data source to get a list of WAF dedicated instances.

## Example Usage

```hcl
variable instance_name {}

data "huaweicloud_waf_dedicated_instances" "test" {
  name = var.instance_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the WAF dedicated instance.
  If omitted, the provider-level region will be used.

* `id` - (Optional, String) Specifies the ID of WAF dedicated instance.

* `name` - (Optional, String) Specifies the name of WAF dedicated instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of WAF dedicated instance.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The WAF dedicated instances list.

The `instances` block supports:

* `id` - The ID of WAF dedicated instance.

* `name` - The name of WAF dedicated instance.

* `available_zone` - The available zone name of the WAF dedicated instance.

* `cpu_architecture` - The ECS CPU architecture of WAF dedicated instance.

* `ecs_flavor` - The flavor of the ECS used by the WAF instance.

* `vpc_id` - The VPC ID of WAF dedicated instance.

* `subnet_id` - The subnet ID of WAF dedicated instance VPC.

* `security_group` - The security group of the instance. This is an array of security group IDs.

* `server_id` - The ID of the ECS hosting the dedicated engine.

* `service_ip` - The service plane IP address of the dedicated WAF instance.

* `run_status` - The running status of the instance. Values are:
  + `0` - Creating.
  + `1` - Running.
  + `2` - Deleting.
  + `3` - Deleted.
  + `4` - Creation failed.
  + `5` - Frozen.
  + `6` - Abnormal.
  + `7` - Updating.
  + `8` - Update failed.

* `access_status` - The access status of the instance. `0`: inaccessible, `1`: accessible.

* `upgradable` - Whether the dedicated WAF instance can be upgraded. `0`: Cannot be upgraded; `1`: Can be upgraded.
