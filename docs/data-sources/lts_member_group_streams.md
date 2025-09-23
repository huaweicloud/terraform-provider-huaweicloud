---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_member_group_streams"
description: |-
  Use this data source to get the list of LTS member group streams.
---

# huaweicloud_lts_member_group_streams

Use this data source to get the list of LTS member group streams.

## Example Usage

```hcl
variable "member_account_id" {}

data "huaweicloud_lts_member_group_streams" "test" {
  member_account_id = var.member_account_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the member group streams.  
  If omitted, the provider-level region will be used.

* `member_account_id` - (Required, String) Specifies the ID of the member account.

  -> This is the domain ID of your account, and the account has already joined the organization.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of the log groups.  
  The [groups](#member_group_streams_groups) structure is documented below.

<a name="member_group_streams_groups"></a>
The `groups` block supports:

* `log_group_id` - The ID of the log group.

* `log_group_name` - The name of the log group.

* `log_streams` - The list of log streams.  
  The [log_streams](#member_group_streams_log_streams) structure is documented below.

<a name="member_group_streams_log_streams"></a>
The `log_streams` block supports:

* `log_stream_id` - The ID of the log stream.

* `log_stream_name` - The name of the log stream.
