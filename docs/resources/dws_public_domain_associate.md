---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_public_domain_associate"
description: |- 
  Use this resource to bind public domain name to DWS cluster within HuaweiCloud.
---
# huaweicloud_dws_public_domain_associate

Use this resource to bind public domain name to DWS cluster within HuaweiCloud.

-> Before using this resource, make sure that the EIP has been bound to the DWS cluster.
   And the `public_bind_type` parameter of the `huaweicloud_dws_cluster` resource cannot be set to **auto_assign**.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "public_domain_name" {}

resource "huaweicloud_dws_public_domain_associate" "test" {
  cluster_id  = var.dws_cluster_id
  domain_name = var.public_domain_name
  ttl         = 1000
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the DWS cluster ID to which the public domain name to be
  associated belongs. Changing this creates a new resource.

* `domain_name` - (Required, String) Specifies the public domain name.  
  The valid length is limited from `4` to `64`, only English letters, digits and hyphens (-) are
  allowed.  
  The name must start with an English letter.

* `ttl` - (Optional, Int) Specifies cache period of the SOA record set, in seconds.  
  The valid value ranges from `300` to `2,147,483,647`. The default value is `300`.

  -> If you want to modify the `ttl` value, the `domain_name` value must also be modified at the same time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also `domain_name`.

## Import

The resource can be imported using the `cluster_id` and `domain_name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dws_public_domain_associate.test <cluster_id>/<domain_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `ttl`. It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dws_public_domain_associate" "test" {
  ...

  lifecycle {
    ignore_changes = [
      ttl,
    ]
  }
}
```
