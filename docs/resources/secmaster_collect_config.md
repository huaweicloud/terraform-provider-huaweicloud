---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collect_config"
description: |-
  Manages a collect config resource within HuaweiCloud.
---

# huaweicloud_secmaster_collect_config

Manages a collect config resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataspace_id" {}
variable "dataspace_name" {}
variable "region_id" {}

resource "huaweicloud_secmaster_collect_config" "test" {
  workspace_id   = var.workspace_id
  dataspace_id   = var.dataspace_id
  dataspace_name = var.dataspace_name
  region_id      = var.region_id

  config {
    source_id      = 1201
    alert          = true
    ttl            = 7
    shards         = 1
    csvc_display   = "数据库审计服务 DBSS"
    source_display = "数据库审计服务告警"
    csvc           = "dbss"
    source_name    = "dbss-alarm"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the collect config
  belongs.

* `dataspace_id` - (Required, String, NonUpdatable) Specifies the ID of the dataspace to which the collect config
  belongs.

* `dataspace_name` - (Required, String, NonUpdatable) Specifies the name of the dataspace to which the collect config
  belongs.

* `region_id` - (Required, String, NonUpdatable) Specifies the region ID of the collect config.

* `domain_id` - (Optional, String, NonUpdatable) Specifies the domain ID of the collect config.

* `config` - (Required, List) Specifies the configuration of the collect config.
  The [config](#collect_config_config) structure is documented below.

* `lts_config` - (Optional, List) Specifies the LTS configuration of the collect config.
  The [lts_config](#collect_config_lts_config) structure is documented below.

<a name="collect_config_config"></a>
The `config` block supports:

* `csvc_display` - (Required, String, NonUpdatable) Specifies the display name of the cloud service.

* `csvc` - (Required, String, NonUpdatable) Specifies the abbreviation of the cloud service.

* `shards` - (Required, Int, NonUpdatable) Specifies the number of shards.

* `source_display` - (Required, String, NonUpdatable) Specifies the display name of the data source.

* `source_id` - (Required, Int, NonUpdatable) Specifies the ID of the data source.

* `ttl` - (Required, Int, NonUpdatable) Specifies the time to live (in days).

* `accounts` - (Optional, List, NonUpdatable) Specifies the list of account IDs.

* `action` - (Optional, String, NonUpdatable) Specifies the action of the collect config.

* `alert` - (Optional, Bool) Specifies whether to enable alert.

* `all_accounts` - (Optional, Bool, NonUpdatable) Specifies whether to apply to all accounts.

* `new_account_auto_access` - (Optional, Bool, NonUpdatable) Specifies whether to automatically access new accounts.

* `source_name` - (Optional, String, NonUpdatable) Specifies the name of the data source.

<a name="collect_config_lts_config"></a>
The `lts_config` block supports:

* `config_name` - (Optional, String, NonUpdatable) Specifies the name of the LTS configuration.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the LTS configuration.

* `enable` - (Optional, String, NonUpdatable) Specifies whether to enable the LTS configuration.

* `log_group_id` - (Optional, String, NonUpdatable) Specifies the ID of the log group.

* `log_stream_id` - (Optional, String, NonUpdatable) Specifies the ID of the log stream.

* `log_type` - (Optional, String, NonUpdatable) Specifies the type of the log.

* `log_type_prefix` - (Optional, String, NonUpdatable) Specifies the prefix of the log type.

* `pipe_alias` - (Optional, String, NonUpdatable) Specifies the alias of the pipe.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the value of `source_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The collect config can be imported using the `region_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_collect_config.test <region_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include `lts_config`. It is generally
recommended running `terraform plan` after importing a resource. You can then decide if changes should be applied to the
resource, or the resource definition should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_collect_config" "test" {
  ...

lifecycle {
  ignore_changes = [
    lts_config
  ]
}
}
```
