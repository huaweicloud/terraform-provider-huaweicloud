---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_registry"
description: |-
  Manages a SWR enterprise instance registry resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_instance_registry

Manages a SWR enterprise instance registry resource within HuaweiCloud.

## Example Usage

### Create the target repository, with the repository type being huawei-SWR (SWR Basic Edition Repository)

```hcl
variable "instance_id" {}
variable "name" {}
variable "url" {}
variable "access_key" {}
variable "access_secret" {}

resource "huaweicloud_swr_enterprise_instance_registry" "test" {
  instance_id = var.instance_id
  name        = var.name
  type        = "huawei-SWR"
  url         = var.url
  insecure    = true
  description = "desc"
  
  credential {
    access_key    = var.access_key
    access_secret = var.access_secret
    type          = "basic"
  }
}
```

### Create the target repository, with the repository type being swr-pro-internal (SWR Enterprise Edition repository)

```hcl
variable "instance_id" {}
variable "name" {}
variable "url" {}
variable "region_id" {}
variable "project_id" {}
variable "target_instance_id" {}
variable "access_key" {}
variable "access_secret" {}

resource "huaweicloud_swr_enterprise_instance_registry" "test" {
  instance_id = var.instance_id
  name        = var.name
  type        = "swr-pro-internal"
  insecure    = true
  description = "desc"

  url                = var.url
  region_id          = var.region_id
  project_id         = var.project_id
  target_instance_id = var.target_instance_id
  
  credential {
    access_key    = var.access_key
    access_secret = var.access_secret
    type          = "basic"
  }
}
```

### Create the target repository, with the repository type being swr-pro (open-source Harbor repository)

```hcl
variable "instance_id" {}
variable "name" {}
variable "url" {}
variable "access_key" {}
variable "access_secret" {}
variable "domain_name" {}
variable "ip_address" {}

resource "huaweicloud_swr_enterprise_instance_registry" "test" {
  instance_id = var.instance_id
  name        = var.name
  type        = "swr-pro"
  url         = var.url
  insecure    = true
  description = "desc"
  
  credential {
    access_key    = var.access_key
    access_secret = var.access_secret
    type          = "basic"
  }

  dns_conf {
    hosts = {
      "${var.domain_name}": var.ip_address
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `type` - (Required, String, NonUpdatable) Specifies the registry type.
  Valid values are as follows:
  + **huawei-SWR**: SWR Basic Edition Repository
  + **swr-pro-internal**: SWR Enterprise Edition repository
  + **swr-pro**: open-source Harbor repository

* `name` - (Required, String) Specifies the registry name.

* `credential` - (Required, List) Specifies the credential infos.
  The [credential](#block--credential) structure is documented below.

* `insecure` - (Required, Bool) Specifies whether the registry is insecure.

* `url` - (Required, String) Specifies the registry url.

* `description` - (Optional, String) Specifies the registry description.

* `dns_conf` - (Optional, List) Specifies the DNS configuration.
  The [dns_conf](#block--dns_conf) structure is documented below.

* `project_id` - (Optional, String) Specifies the project ID of the target registry.

* `region_id` - (Optional, String) Specifies the region ID of the target registry.

* `target_instance_id` - (Optional, String) Specifies the target enterprise instance ID.

-> `project_id`, `region_id` and `target_instance_id` is required if `type` is **swr-pro-internal**

<a name="block--credential"></a>
The `credential` block supports:

* `type` - (Required, String) Specifies the credential type. Valid value is **basic**.

* `access_key` - (Required, String) Specifies the access key.

* `access_secret` - (Required, String) Specifies the access secret.

<a name="block--dns_conf"></a>
The `dns_conf` block supports:

* `hosts` - (Optional, Map) Specifies the hosts map. Key is domain name, and value is IP address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `registry_id` - Indicates the registry ID.

* `status` - Indicates the registry status.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

## Import

The registry can be imported using `instance_id` and `registry_id`, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_instance_registry.test <instance_id>/<registry_id>
```

Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `credential.0.access_secret`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the registry, or the resource definition should be updated to
align with the registry. Also you can ignore changes as below.

```hcl
resource "huaweicloud_swr_enterprise_instance_registry" "test" {
    ...

  lifecycle {
    ignore_changes = [
      credential.0.access_secret,
    ]
  }
}
```
