---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_datasource_auth"
description: ""
---

# huaweicloud_dli_datasource_auth

Manages a DLI datasource authentication resource within HuaweiCloud.  

## Example Usage

### Create a datasource Password authentication

```hcl
  variable "username" {}
  variable "password" {}
  
  resource "huaweicloud_dli_datasource_auth" "test" {
    type     = "passwd"
    name     = "demo"
    username = var.username
    password = var.password
  }
```

### Create a datasource CSS authentication

```hcl
  variable "username" {}
  variable "password" {}
  variable "certificate_location" {}
  
  resource "huaweicloud_dli_datasource_auth" "test" {
    type                 = "CSS"
    name                 = "demo"
    username             = var.username
    password             = var.password
    certificate_location = var.certificate_location
  }
```

### Create a datasource Kafka_SSL authentication

```hcl
  variable "truststore_location" {}
  variable "truststore_password" {}
  variable "keystore_location" {}
  variable "keystore_password" {}
  variable "key_password" {}
  
  resource "huaweicloud_dli_datasource_auth" "test" {
    type                = "Kafka_SSL"
    name                = "demo"
    truststore_location = var.truststore_location
    truststore_password = var.truststore_password
    keystore_location   = var.keystore_location
    keystore_password   = var.keystore_password
    key_password        = var.key_password
  }
```

### Create a datasource Kerberos authentication

```hcl
  variable "username" {}
  variable "krb5_conf" {}
  variable "keytab" {}
  
  resource "huaweicloud_dli_datasource_auth" "test" {
    type      = "KRB"
    name      = "demo"
    username  = var.username
    krb5_conf = var.krb5_conf
    keytab    = var.keytab
  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The name of a datasource authentication.  
  Only letters, digits and underscores (_) are allowed.
  And the name cannot be all digits or start with a underscore.

* `type` - (Required, String, ForceNew) Data source type.  
  The options are as follows:
    + **passwd**: Password.
    + **CSS**: CSS.
    + **KRB**: Kafka_SSL.
    + **Kafka_SSL**: Kafka_SSL.

  Changing this parameter will create a new resource.

* `user_name` - (Optional, String) Specifies the user name for accessing the security cluster or datasource.

* `password` - (Optional, String) The password for accessing the security cluster or datasource.
  This parameter must be used together with `user_name`.

* `certificate_location` - (Optional, String, ForceNew) Path of the security cluster certificate.  
 Currently, only OBS paths and CER files are supported.

  Changing this parameter will create a new resource.

* `truststore_location` - (Optional, String) The OBS path of the **truststore** configuration file.

* `truststore_password` - (Optional, String) The password of the **truststore** configuration file.

* `keystore_location` - (Optional, String) The OBS path of the **keystore** configuration file.

* `keystore_password` - (Optional, String) The password of the **keystore** configuration file.

* `key_password` - (Optional, String, ForceNew) The key password.

  Changing this parameter will create a new resource.

* `krb5_conf` - (Optional, String) The OBS path of the **krb5** configuration file.

* `keytab` - (Optional, String) The OBS path of the **keytab** configuration file.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `name`.

* `owner` - The user name of owner.

## Import

The DLI datasource authentication can be imported using `id` which equals the `name`, e.g.

```bash
$ terraform import huaweicloud_dli_datasource_auth.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`password`, `truststore_password`, `keystore_password`, `key_password`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dli_datasource_auth" "test" {
    ...

  lifecycle {
    ignore_changes = [
      password, truststore_password, keystore_password, key_password
    ]
  }
}
```
