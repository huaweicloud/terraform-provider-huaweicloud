---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_source_server_overview"
description: |-
  Use this data source to obtain a summary of source servers.
---

# huaweicloud_sms_source_server_overview

Use this data source to obtain a summary of source servers.

## Example Usage

```hcl
data "huaweicloud_sms_source_server_overview" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `waiting` - Indicates the number of servers that are in a waiting migration status.

* `replicate` - Indicates the number of servers that are in a replicating migration status.

* `syncing` - Indicates the number of servers that are in a synchronizing migration status.

* `stopped` - Indicates the number of servers that are in a paused migration status.

* `deleting` - Indicates the number of servers that are in a deleting migration status.

* `cutovering` - Indicates the number of servers whose paired target servers are being launched.

* `unavailable` - Indicates the number of servers that fail the environment check.

* `stopping` - Indicates the number of servers that are in a stopping migration status.

* `skipping` - Indicates the number of servers that are in a skipping migration status.

* `finished` - Indicates the number of servers whose paired target servers have been launched.

* `initialize` - Indicates the number of servers that are in an initializing migration status.

* `error` - Indicates the number of servers that are in an error migration status.

* `cloning` - Indicates the number of servers whose paired target servers are being cloned.

* `unconfigured` - Indicates the number of servers that do not have target server configurations.
