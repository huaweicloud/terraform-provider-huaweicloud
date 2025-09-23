---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificate_push"
description: |
  Manages a CCM resource to push SSL certificate to ELB or WAF service within HuaweiCloud.
---

# huaweicloud_ccm_certificate_push

Manages a CCM resource to push SSL certificate to ELB or WAF service within HuaweiCloud.

## Example Usage

### Pushing a certificate to ELB service

```hcl
variable "certificate_id" {}
variable "project_name1" {}
variable "project_name2" {}

resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = var.certificate_id
  service        = "ELB"
  
  targets {
    project_name = var.project_name1
  }
  
  targets {
    project_name = var.project_name2
  }
}
```

### Pushing a certificate to WAF service

```hcl
variable "certificate_id" {}
variable "project_name" {}

resource "huaweicloud_ccm_certificate_push" "test"{
  certificate_id = var.certificate_id
  service        = "WAF"
  
  targets {
    project_name = var.project_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the certificate region. Changing this creates a new
  private certificate resource. Now only support cn-north-4 (china) and ap-southeast-1 (international).

* `certificate_id` - (Required, String, ForceNew) Specifies the ID of the SSL certificate to be pushed.
  Changing this parameter will create a new resource.

* `service` - (Required, String, ForceNew) Specifies the target service for certificate push.
  Changing this creates a new resource. Valid values are **ELB** and **WAF**.

  -> Before pushing the certificate to the WAF service, please confirm that the target region has an available WAF instance.

* `targets` - (Required, List) Specifies the projects which certificate will be push to.
  The [targets](#block-targets) structure is documented below.

<a name="block-targets"></a>
The `targets` block supports:

* `project_name` - (Required, String) Specifies the region where the target service for certificate push is located.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the SSL certificate ID.

* `targets` - Indicates the projects which certificate pushed to.
  The [targets](#block-targets-attr) structure is documented below.

<a name="block-targets-attr"></a>
The `targets` block supports:

* `cert_id` - Indicates the target certificate ID.

* `cert_name` - Indicates the target certificate name.
