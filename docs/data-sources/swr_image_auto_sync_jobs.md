---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_auto_sync_jobs"
description: |-
  Use this data source to get the list of SWR image auto sync jobs.
---

# huaweicloud_swr_image_auto_sync_jobs

Use this data source to get the list of SWR image auto sync jobs.

## Example Usage

```hcl
variable "organization" {}
variable "repository" {}

data "huaweicloud_swr_image_auto_sync_jobs" "test" {
  organization = var.organization
  repository   = var.repository
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `organization` - (Required, String) Specifies the organization name. Enter 1 to 64 characters, starting with a
  lowercase letter and ending with a lowercase letter or digit. Only lowercase letters, digits, periods (.),
  underscores (_), and hyphens (-) are allowed. Periods, underscores, and hyphens cannot be placed next to each other.
  A maximum of two consecutive underscores are allowed.

* `repository` - (Required, String) Specifies the image repository name. Enter 1 to 128 characters. It must start and
  end with a lowercase letter or digit. Only lowercase letters, digits, periods (.), slashes (/), underscores (_), and
  hyphens (-) are allowed. Periods, slashes, underscores, and hyphens cannot be placed next to each other. A maximum of
  two consecutive underscores are allowed. Replace a slash (/) with a dollar sign ($) before you send the request.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - Indicates the jobs.
  The [jobs](#attrblock--jobs) structure is documented below.

<a name="attrblock--jobs"></a>
The `jobs` block supports:

* `id` - Indicates the job ID.

* `status` - Indicates the synchronization status.

* `sync_operator_id` - Indicates the operator account ID.

* `sync_operator_name` - Indicates the operator account name.

* `remote_organization` - Indicates the target organization.

* `remote_region_id` - Indicates the target region.

* `organization` - Indicates the name of the organization.

* `override` - Indicates whether to overwrite.

* `repo_name` - Indicates the repository name.

* `tag` - Indicates the image tag.

* `created_at` - Indicates the time when the job was created. It is the UTC standard time.

* `updated_at` - Indicates the time when the job was updated. It is the UTC standard time.
