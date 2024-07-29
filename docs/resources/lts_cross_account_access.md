---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_cross_account_access"
description: |-
  Manages an LTS cross account access resource within HuaweiCloud.
---

# huaweicloud_lts_cross_account_access

Manages an LTS cross account access resource within HuaweiCloud.

-> **NOTE:** Somethings you should know before using this resource.
<br/> 1. Before creating this resource, you need to prepare an agency relationship.
<br/> 2. Before data synchronization is complete, data in the target and source log streams may be different.
         Check back later in one hour.
<br/> 3. After you configure cross-account access, if account A deletes the agency from IAM, LTS cannot detect the
         deletion and the cross-account ingestion still takes effect. If the cross-account access configuration is
         no longer used, notify account B to delete it.

## Example Usage

```hcl
variable "name" {}
variable "agency_group_id" {}
variable "agency_stream_id" {}
variable "agency_group_name" {}
variable "agency_stream_name" {}
variable "log_group_id" {}
variable "log_stream_id" {}
variable "log_group_name" {}
variable "log_stream_name" {}
variable "agency_name" {}
variable "agency_domain_name" {}
variable "agency_project_id" {}

resource "huaweicloud_lts_cross_account_access" "test" {
  name               = var.name
  agency_project_id  = var.agency_project_id
  agency_domain_name = var.agency_domain_name
  agency_name        = var.agency_name

  # Relevant information of the delegated party.
  log_agencystream_name = var.agency_stream_name
  log_agencystream_id   = var.agency_stream_id
  log_agencygroup_name  = var.agency_group_name
  log_agencygroup_id    = var.agency_group_id

  # Relevant information of the delegatee party.
  log_stream_name = var.log_stream_name
  log_stream_id   = var.log_stream_id
  log_group_name  = var.log_group_name
  log_group_id    = var.log_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the cross account access.
  Changing this creates a new resource.

* `agency_domain_name` - (Required, String, ForceNew) Specifies the name of the delegator account to verify
  the delegation. Changing this creates a new resource.

* `agency_name` - (Required, String, ForceNew) Specifies the name of the agency created in IAM by the delegator.
  Changing this creates a new resource.

* `agency_project_id` - (Required, String, ForceNew) Specifies the delegator project ID.
  Changing this creates a new resource.

* `log_agencygroup_id` - (Required, String, ForceNew) Specify the log group ID that already exists in the
  delegated account. Changing this creates a new resource.

* `log_agencygroup_name` - (Required, String, ForceNew) Specify the log group name that already exists in the
  delegated account. Changing this creates a new resource.

* `log_agencystream_id` - (Required, String, ForceNew) Specifies the log stream ID that already exists in the
  delegated account. Changing this creates a new resource.

* `log_agencystream_name` - (Required, String, ForceNew) Specifies the log stream name that already exists in the
  delegated account. Changing this creates a new resource.

* `log_group_id` - (Required, String, ForceNew) Specify the log group ID that already exists in the
  delegatee account. Changing this creates a new resource.

* `log_group_name` - (Required, String, ForceNew) Specify the log group name that already exists in the
  delegatee account. Changing this creates a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the log stream ID that already exists in the
  delegatee account. Changing this creates a new resource.

* `log_stream_name` - (Required, String, ForceNew) Specifies the log stream name that already exists in the
  delegatee account. Changing this creates a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the cross account access.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `access_config_type` - The log access configuration type.

* `created_at` - The creation time of the cross account access, in RFC3339 format.
