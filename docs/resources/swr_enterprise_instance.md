---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance"
description: |-
  Manages a SWR enterprise instance resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_instance

Manages a SWR enterprise instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "spec" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_swr_enterprise_instance" "test" {
  name                  = var.name
  spec                  = var.spec
  vpc_id                = var.vpc_id
  subnet_id             = var.subnet_id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the instance.

* `spec` - (Required, String, NonUpdatable) Specifies the specification of the instance. Value can be **swr.ee.professional**.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID .

* `subnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.

* `anonymous_access` - (Optional, Bool) Specifies whether to enable anonymous access. Default to **false**.

* `public_network_access_control_status` - (Optional, String) Specifies the public network access control status.
  Options are **Disable** and **Enable**. Default to **Disable**.

* `public_network_access_white_ip_list` - (Optional, List) Specifies the public network access white IP list.
  The [public_network_access_white_ip_list](#block--public_network_access_white_ip_list) structure is documented below.

  -> It can be updated only when `public_network_access_control_status` is **Enable**.

* `delete_dns` - (Optional, Bool) Specifies whether to delete DNS resources when deleting instance. Default to **false**.

* `delete_obs` - (Optional, Bool) Specifies whether to delete OBS bucket when deleting instance. Default to **false**.

* `obs_encrypt` - (Optional, Bool, NonUpdatable) Specifies whether the OBS bucket is encrypted. Default to **false**.

* `encrypt_type` - (Optional, String, NonUpdatable) Specifies the encrypt type.

* `obs_bucket_name` - (Optional, String, NonUpdatable) Specifies the OBS bucket name.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the instance.

* `description` - (Optional, String, NonUpdatable) Specifies the description.

<a name="block--public_network_access_white_ip_list"></a>
The `public_network_access_white_ip_list` block supports:

* `ip` - (Required, String) Specifies the IP address or CIDR block.

* `description` - (Optional, String) Specifies the description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `charge_mode` - Indicates the charge mode of instance.

* `status` - Indicates the instance status.

* `statistics` - Indicates the statistic infos.
  The [statistics](#attrblock--statistics) structure is documented below.

* `access_address` - Indicates the access address of instance.

* `user_def_obs` - Indicates whether the user specifies the OBS bucket.

* `version` - Indicates the instance version.

* `vpc_cidr` - Indicates the range of available subnets for the VPC.

* `vpc_name` - Indicates the VPC name.

* `subnet_cidr` - Indicates the range of available subnets for the subnet.

* `subnet_name` - Indicates the subnet name.

* `created_at` - Indicates the creation time.

* `expires_at` - Indicates the expired time.

* `updated_at` - Indicates the last update time.

<a name="attrblock--statistics"></a>
The `statistics` block supports:

* `image_repo_quota` - Indicates the image repo quota.

* `intranet_endpoint_count` - Indicates the intranet endpoint count.

* `intranet_endpoint_quota` - Indicates the intranet endpoint quota.

* `long_term_quota` - Indicates the long term quota.

* `namespace_quota` - Indicates the namespace quota.

* `notify_policy_count` - Indicates the notify policy count.

* `notify_policy_quota` - Indicates the notify policy quota.

* `replica_policy_count` - Indicates the replica policy count.

* `replica_policy_quota` - Indicates the replica policy quota.

* `replica_registry_count` - Indicates the replica registry count.

* `replica_registry_quota` - Indicates the replica registry quota.

* `retention_policy_count` - Indicates the retention policy count.

* `retention_policy_quota` - Indicates the retention policy quota.

* `sign_policy_count` - Indicates the sign policy count.

* `sign_policy_quota` - Indicates the sign policy quota.

* `storage_used` - Indicates the storage used.

* `total_image_count` - Indicates the total image count.

* `total_namespace_count` - Indicates the total namespace count.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 40 minutes.
* `update` - Default is 20 minutes.

## Import

The instance can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `delete_obs`, `delete_dns`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_swr_enterprise_instance" "test" {
    ...

  lifecycle {
    ignore_changes = [
      delete_obs, delete_dns,
    ]
  }
}
```
