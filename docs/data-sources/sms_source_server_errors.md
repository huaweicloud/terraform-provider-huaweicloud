---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_source_server_errors"
description: |-
  Use this data source to query the list of source servers that failed to be migrated and the reported error messages.
---

# huaweicloud_sms_source_server_errors

Use this data source to query the list of source servers that failed to be migrated and the reported error messages.

## Example Usage

```hcl
data "huaweicloud_sms_source_server_errors" "test" {}
```

## Argument Reference

The following arguments are supported:

* `migproject` - (Optional, String) Specifies the migration project ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `migration_errors` - Indicates the details of failed source servers.

  The [migration_errors](#migration_errors_struct) structure is documented below.

<a name="migration_errors_struct"></a>
The `migration_errors` block supports:

* `error_json` - Indicates the error message in JSON format.

* `host_name` - Indicates the host name of the source server.

* `name` - Indicates the source server name in SMS.

* `source_id` - Indicates the source server ID.

* `source_ip` - Indicates the IP address of the source server.

* `target_ip` - Indicates the IP address of the target server.
