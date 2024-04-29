---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificate_push"
description: ""
---

# huaweicloud_ccm_certificate_push

Push a **SSL** ceritificate to **ELB** or **WAF** within HuaweiCloud.

## Example Usage

```hcl
variable "certificate_id" {}
variable "project_name1" {}
variable "project_name2" {}

resource "huaweicloud_ccm_certificate_push" "test"{
  region         = "cn-north-4"
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

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the certificate region. Changing this creates a new
  private certificate resource. Now only support cn-north-4 (china) and ap-southeast-1 (international).

* `certificate_id` - (Required, String, ForceNew) Specifies the certificate which will be push.
  Changing this parameter will create a new resource.

* `service` - (Required, String, ForceNew) Specifies the service which certificate will be push to.
  Changing this creates a new resource. Valid values are **ELB** and **WAF**.

* `targets` - (Required, List) Specifies the projects which certificate will be push to.
  The [targets](#block-targets) structure is documented below.

<a name="block-targets"></a>
The `targets` block supports:

* `project_name` - (Required, String) Specifies the project which certificate will be push to.
  Only the **default** projects with the same region name are supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the SSL certificate ID.

* `targets` - Indicates the projects which certificate pushed to.
  The [targets](#block-targets-attr) structure is documented below.

<a name="block-targets-attr"></a>
The `targets` block supports:

* `cert_id` - Indicates the pushed certificate new ID.

* `cert_name` - Indicates the pushed certificate name.
