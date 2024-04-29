---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_datasource_auths"
description: ""
---

# huaweicloud_dli_datasource_auths

Use this data source to get the list of DLI datasource authentications within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_dli_datasource_auths" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the datasource authentication.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `auths` - The list of the datasource authentications.
  The [auths](#datasource_auths) structure is documented below.

<a name="datasource_auths"></a>
The `auths` block supports:

* `name` - The name of the datasource authentication.

* `type` - The type of the datasource authentication.

* `username` - The login user name of the security cluster.

* `certificate_location` - The OBS path of the security cluster certificate.

* `truststore_location` - The OBS path of the **truststore** configuration file.

* `keystore_location` - The OBS path of the **keystore** configuration file.

* `keytab` - The OBS path of the **keytab** configuration file.

* `krb5_conf` - The OBS path of the **krb5** configuration file.

* `owner` - The user name of owner.

* `created_at` - The creation time of the datasource authentication.

* `updated_at` - The latest update time of the datasource authentication.
