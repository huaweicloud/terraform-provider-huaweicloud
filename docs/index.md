# HuaweiCloud Provider

The HuaweiCloud provider is used to interact with the many resources
supported by HuaweiCloud. The provider needs to be configured with
the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the HuaweiCloud Provider
provider "huaweicloud" {
  region      = "cn-north-1"
  domain_name = "my-account-name"
  access_key  = "my-access-key"
  secret_key  = "my-secret-key"
}

# Create a VPC
resource "huaweicloud_vpc" "example" {
  name = "my_vpc"
  cidr = "192.168.0.0/16"
}
```

## Authentication

The Huawei Cloud provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static credentials ###

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `access_key` and `secret_key`
in-line in the provider block:

Usage:

```hcl
provider "huaweicloud" {
  region      = "cn-north-1"
  domain_name = "my-account-name"
  access_key  = "my-access-key"
  secret_key  = "my-secret-key"
}
```
-> **NOTE:** `domain_name`, [Account name](https://support.huaweicloud.com/en-us/usermanual-iam/iam_01_0552.html) need to be set if using IAM or prePaid resources.

### Environment variables

You can provide your credentials via the `OS_ACCESS_KEY` and
`OS_SECRET_KEY`, environment variables, representing your Huawei
Cloud Access Key and Secret Key, respectively.

```hcl
provider "huaweicloud" {}
```

Usage:

```sh
$ export OS_ACCESS_KEY="anaccesskey"
$ export OS_SECRET_KEY="asecretkey"
$ export OS_REGION_NAME="cn-north-1"
$ export OS_DOMAIN_NAME="account-name"
$ terraform plan
```

## Configuration Reference

The following arguments are supported:

* `region` - (Required) This is the Huawei Cloud region. It must be provided,
  but it can also be sourced from the `OS_REGION_NAME` environment variables.

* `domain_name` - (Optional, Required for IAM and prePaid resources) The
  [Account name](https://support.huaweicloud.com/en-us/usermanual-iam/iam_01_0552.html)
  of IAM to scope to. If omitted, the `OS_DOMAIN_NAME` environment variable is used.

* `access_key` - (Optional) The access key of the HuaweiCloud to use.
  If omitted, the `OS_ACCESS_KEY` environment variable is used.

* `secret_key` - (Optional) The secret key of the HuaweiCloud to use.
  If omitted, the `OS_SECRET_KEY` environment variable is used.

* `tenant_name` - (Optional) The Name of the Tenant/Project to login with.
  If omitted, the `OS_TENANT_NAME` or `OS_PROJECT_NAME` environment variable are used.

* `auth_url` - (Optional, Required before 1.14.0) The Identity authentication URL. If omitted, the
  `OS_AUTH_URL` environment variable is used. This is not required if you use Huawei Cloud.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `OS_INSECURE` environment variable is used.

* `max_retries` - (Optional) This is the maximum number of times an API
  call is retried, in the case where requests are being throttled or
  experiencing transient failures. The delay between the subsequent API
  calls increases exponentially. If omitted, default value is `5`.


## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `OS_REGION_NAME` - The region in which to create the resources.

* `OS_ACCESS_KEY` - The access key of the HuaweiCloud to use.

* `OS_SECRET_KEY` - The secret key of the HuaweiCloud to use.


You should be able to use any HuaweiCloud environment to develop on as long as the
above environment variables are set.
