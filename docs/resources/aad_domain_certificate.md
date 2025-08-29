---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_domain_certificate"
description: |-
  Manages an AAD domain certificate resource within HuaweiCloud.
---

# huaweicloud_aad_domain_certificate

Manages an AAD domain certificate resource within HuaweiCloud.

-> This resource is a one-time action resource using to upload the certificate corresponding to the domain. Deleting
   this resource will not change the current request record, but will only remove the resource information from the tf
   state file.

## Example Usage

```hcl
variable "domain_id" {}
variable "cert_name" {}

resource "huaweicloud_aad_domain_certificate" "test" {
  domain_id = var.domain_id
  op_type   = 1
  cert_name = var.cert_name
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, String, NonUpdatable) Specifies the domain ID.

* `op_type` - (Required, Int, NonUpdatable) Specifies the operation type.  
  The valid values are as follows:
  + **0**: Upload new certificate.
  + **1**: Replace with existing certificate.

* `cert_name` - (Required, String, NonUpdatable) Specifies the certificate name.

* `cert_file` - (Optional, String, NonUpdatable) Specifies the certificate file content.  
  + Required when `op_type` is `0` (Required when uploading a new certificate).
  + Set empty when `op_type` is `1` (Empty when replacing with an existing certificate).

* `cert_key_file` - (Optional, String, NonUpdatable) Specifies the private key file content.  
  + Required when `op_type` is `0` (Required when uploading a new certificate).
  + Set empty when `op_type` is `1` (Empty when replacing with an existing certificate).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is also the domain ID.

* `domain_name` - The domain name.

* `cert_info` - The certificate information.  
  The [cert_info](#cert_info_struct) structure is documented below.

<a name="cert_info_struct"></a>
The `cert_info` block supports:

* `cert_name` - The certificate name.

* `id` - The certificate ID.

* `apply_domain` - The applicable domain.

* `expire_time` - The expiration time.

* `expire_status` - The expiration status.

## Import

The AAD domain certificate can be imported using the `domain_id`, e.g.

```bash
terraform import huaweicloud_aad_domain_certificate.test <domain_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `op_type`, `cert_file`, and `cert_key_file`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_aad_domain_certificate" "test" {
  ...

  lifecycle {
    ignore_changes = [
      op_type, cert_file, cert_key_file,
    ]
  }
}
```
