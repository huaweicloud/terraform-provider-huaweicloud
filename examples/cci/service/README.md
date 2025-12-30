# Create a CCI Service

This example provides best practice code for using Terraform to create a Cloud Container Instance (CCI) service
in HuaweiCloud with all necessary networking components. The example demonstrates how to create a VPC, subnet, security
group, CCI namespace, ELB load balancer, and CCI service with LoadBalancer type.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* CCI service enabled in target region
* ELB service enabled in target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where CCI service is located
* `access_key` - The access key of IAM user
* `secret_key` - The secret key of IAM user

### Resource Variables

#### Required Variables

* `elb_name` - The name of ELB load balancer
* `service_name` - The name of CCI service
* `namespace_name` - The name of CCI namespace

#### Optional Variables

* `vpc_name` - The name of VPC (default: "tf-test-vpc")
* `vpc_cidr` - The CIDR block of VPC (default: "192.168.0.0/16")
* `subnet_name` - The name of subnet (default: "tf-test-subnet")
* `subnet_cidr` - The CIDR block of subnet (default: "192.168.0.0/24")
* `subnet_gateway_ip` - The gateway IP of subnet (default: "192.168.0.1")
* `security_group_name` - The name of security group (default: "tf-test-secgroup")
* `selector_app` - The app label of selector (default: "test1")
* `service_type` - The type of service (default: "LoadBalancer")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  elb_name       = "tf-test-elb"
  service_name   = "tf-test-service"
  namespace_name = "tf-test-namespace"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Service Configuration

### ELB Load Balancer Configuration

The service uses an ELB load balancer for external access:

```hcl
resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "tf-test-elb"
  cross_vpc_backend = true
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[1]
  ]
}
```

### Service Selector Configuration

The service uses a selector to route traffic to pods:

```hcl
selector = {
  app = "test1"
}
```

### Service Annotations

The service uses annotations to reference the ELB load balancer:

```hcl
annotations = {
  "kubernetes.io/elb.class" = "elb"
  "kubernetes.io/elb.id"    = huaweicloud_elb_loadbalancer.test.id
}
```

### Lifecycle Configuration

The service ignores changes to annotations to prevent unnecessary updates:

```hcl
lifecycle {
  ignore_changes = [
    annotations,
  ]
}
```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The service must be created within an existing namespace
* The service type "LoadBalancer" requires an ELB load balancer
* The selector labels must match pod labels
* The ELB load balancer must be created before the service
* All resources will be created in the specified region
* Service names must be unique within the namespace
* Make sure to have sufficient quota for the resources you plan to create
* The security group is created with default rules deleted for a clean configuration
* The availability zone for the ELB load balancer is automatically selected from the second available zone

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.73.9 |
