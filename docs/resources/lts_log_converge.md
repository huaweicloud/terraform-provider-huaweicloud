---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_log_converge"
description: |-
  Using this resource to converge logs cross accounts to centrally store and analyze logs within HuaweiCloud.
---

# huaweicloud_lts_log_converge

Using this resource to converge logs cross accounts to centrally store and analyze logs within HuaweiCloud.

~> 1. Deleting the log groups or streams of the source account does not affect their target equivalents.
   <br>2. Deleting this resource will not delete the target log group and log stream synchronously, so you need to clean
   them up manually.
   <br>3. Aggregation cannot be configured for source log groups without log streams.

## Example Usage

### Use existing log group and log stream as target log group and target log stream

```hcl
variable "organization_id" {}
variable "administrator_account_id" {}
variable "administrator_project_id" {}
variable "member_account_id" {}
variable "source_log_group_id" {}
variable "target_log_group_name" {}
variable "target_log_group_id" {}
variable "source_log_stream_id" {}
variable "target_log_stream_name" {}
variable "target_log_stream_id" {}
variable "target_log_stream_ttl" {}

resource "huaweicloud_lts_log_converge" "test" {
  organization_id       = var.organization_id
  management_account_id = var.administrator_account_id
  management_project_id = var.administrator_project_id
  member_account_id     = var.member_account_id

  log_mapping_config {
    source_log_group_id   = var.source_log_group_id
    target_log_group_name = var.target_log_group_name
    target_log_group_id   = var.target_log_group_id

    log_stream_config {
      source_log_stream_id   = var.source_log_stream_id
      target_log_stream_name = var.target_log_stream_name
      target_log_stream_id   = var.target_log_stream_id
      target_log_stream_ttl  = var.target_log_stream_ttl
    }
  }
}
```

### Use the log group and log stream automatically created by the service as the target log group and target log stream

```hcl
variable "organization_id" {}
variable "administrator_account_id" {}
variable "administrator_project_id" {}
variable "member_account_id" {}
variable "source_log_group_id" {}
variable "target_log_group_name" {}
variable "source_log_stream_id" {}
variable "target_log_stream_name" {}
variable "target_log_stream_ttl" {}

resource "huaweicloud_lts_log_converge" "test" {
  organization_id       = var.organization_id
  management_account_id = var.administrator_account_id
  management_project_id = var.administrator_project_id
  member_account_id     = var.member_account_id

  log_mapping_config {
    source_log_group_id   = var.source_log_group_id
    target_log_group_name = var.target_log_group_name

    log_stream_config {
      source_log_stream_id   = var.source_log_stream_id
      target_log_stream_name = var.target_log_stream_name
      target_log_stream_ttl  = var.target_log_stream_ttl
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `organization_id` - (Required, String, ForceNew) Specifies the organization ID to which the converged logs belong.

  Changing this parameter will create a new resource.

* `management_account_id` - (Required, String, ForceNew) Specifies the administrator account ID used to manage log converge.

  Changing this parameter will create a new resource.

* `member_account_id` - (Required, String, ForceNew) Specifies the member account ID to which the converged logs belong.

  Changing this parameter will create a new resource.

* `log_mapping_config` - (Required, List) Specifies the log converge configurations.  
  The [log_mapping_config](#converge_log_mapping_config) structure is documented below.

* `management_project_id` - (Optional, String) Specifies the administrator project ID that required for
  first-time use.

<a name="converge_log_mapping_config"></a>
The `log_mapping_config` block supports:

* `source_log_group_id` - (Required, String) Specifies the ID of the log group for source side.

* `target_log_group_name` - (Required, String) Specifies the name of the log group for target side.

* `target_log_group_id` - (Optional, String) Specifies the ID of the log group for target side.

  -> If you want to use an existing log group, this parameter (`target_log_group_id`) is required.

* `log_stream_config` - (Optional, List) Specifies the log streams converged under the current log group.  
  The [log_stream_config](#converge_log_streams_config) structure is documented below.

<a name="converge_log_streams_config"></a>
The `log_stream_config` block supports:

* `source_log_stream_id` - (Required, String) Specifies the ID of the log stream for source side.

* `target_log_stream_ttl` - (Required, Int) Specifies the ID of the log stream for source side.

* `target_log_stream_name` - (Required, String) Specifies the ID of the log stream for source side.

* `target_log_stream_id` - (Optional, String) Specifies the ID of the log stream for source side.

  -> If you want to use an existing log stream, this parameter (`target_log_stream_id`) is required.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consisting of `organization_id` and `member_account_id`, the format is `<organization_id>/<member_account_id>`.

* `log_mapping_config` - The log converge configurations  
  The [log_mapping_config](#converge_log_mapping_config_attr) structure is documented below.

* `created_at` - The creation time of the log converge configuration.

* `updated_at` - The latest update time of the log converge configuration.

<a name="converge_log_mapping_config_attr"></a>
The `log_mapping_config` block supports:

* `log_stream_config` - The log streams converged under the current log group.  
  The [log_stream_config](#converge_log_streams_config_attr) structure is documented below.

<a name="converge_log_streams_config_attr"></a>
The `log_stream_config` block supports:

* `target_log_stream_eps_id` - The enterprise project ID of the log stream for target side.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 15 minutes.
* `update` - Default is 15 minutes.
* `delete` - Default is 10 minutes.

## Import

The log converge can be imported using the `organization_id` and `member_account_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_lts_log_converge.test <organization_id>/<member_account_id>
```
