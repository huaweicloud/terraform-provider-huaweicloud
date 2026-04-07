---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_policy_dry_run_configuration"
description: |-
  Manages an Organizations policy dry-run configuration resource within HuaweiCloud.
---

# huaweicloud_organizations_policy_dry_run_configuration

Manages an Organizations policy dry-run configuration resource within HuaweiCloud.

~> Destroying this resource will disable the policy dry-run configuration.

-> Before using this resource, ensure that the service control policy has been enabled.

## Example Usage

```hcl
variable "root_id" {}
variable "bucket_name" {}
variable "region_id" {}
variable "agency_name" {}

resource "huaweicloud_organizations_policy_dry_run_configuration" "test" {
  root_id     = var.root_id
  policy_type = "service_control_policy"
  status      = "enabled"
  bucket_name = var.bucket_name
  region_id   = var.region_id
  agency_name = var.agency_name
}
```

## Argument Reference

The following arguments are supported:

* `root_id` - (Required, String, NonUpdatable) Specifies the ID of the organization's root.

* `policy_type` - (Required, String, NonUpdatable) Specifies the type of the policy.  
  The valid values are as follows:
  + **service_control_policy**: Service control policy.

* `status` - (Optional, String) Specifies the status of the policy dry-run.  
  The valid values are as follows:
  + **enabled**
  + **disabled**

* `bucket_name` - (Optional, String) Specifies the name of the OBS bucket.

* `region_id` - (Optional, String) Specifies the region where the OBS bucket is located.

* `bucket_prefix` - (Optional, String) Specifies the prefix of the OBS bucket.

* `agency_name` - (Optional, String) Specifies the name of the IAM agency.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is formatted as `<root_id>/<policy_type>`.

* `created_at` - The creation time of the policy dry-run configuration.

* `updated_at` - The latest update time of the policy dry-run configuration.

## Import

The policy dry-run configuration can be imported using the `root_id` and `policy_type`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_organizations_policy_dry_run_configuration.test <root_id>/<policy_type>
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.
