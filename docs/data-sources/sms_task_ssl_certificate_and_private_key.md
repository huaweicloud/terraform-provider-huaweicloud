---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_task_ssl_certificate_and_private_key"
description: |-
  Use this data source to download the certificate and private key (in PEM format) required for data migration.
---

# huaweicloud_sms_task_ssl_certificate_and_private_key

Use this data source to download the certificate and private key (in PEM format) required for data migration.

## Example Usage

```hcl
variable "task_id" {}

data "huaweicloud_sms_task_ssl_certificate_and_private_key" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Specifies the migration task ID.

* `enable_ca_cert` - (Optional, Bool) Specifies whether to generate a CA certificate. The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `cert` - Indicates the source certificate.

* `private_key` - Indicates the source private key.

* `ca` - Indicates the CA certificate.

* `target_mgmt_cert` - Indicates the certificate of the target server for migration task management.

* `target_mgmt_private_key` - Indicates the private key of the target server for migration task management.

* `target_data_cert` - Indicates the certificate of the target server for data migration.

* `target_data_private_key` - Indicates the private key of the target server for data migration.
