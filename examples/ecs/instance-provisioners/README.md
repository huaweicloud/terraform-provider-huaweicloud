# Using Provisioner over SSH to an ECS instance

This example provides best practice code for using Terraform to create an ECS instance with SSH provisioners in
HuaweiCloud. The example demonstrates how to provision an ECS instance with a Public IP address and run a `remote-exec`
provisioner over SSH.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ECS instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `keypair_name` - The name of the SSH keypair
* `private_key_path` - The path to the private key file
* `instance_name` - The name of the ECS instance
* `bandwidth_name` - The name of the EIP bandwidth (required if EIP address is not provided)
* `instance_remote_exec_inline` - The inline script for remote execution

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `availability_zone` - The availability zone for the ECS instance (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `security_group_name` - The name of the security group (default: "")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_cpu_core_count` - The number of CPU cores for the ECS instance (default: 2)
* `instance_memory_size` - The memory size in GB for the ECS instance (default: 4)
* `instance_image_id` - The image ID of the ECS instance (default: "")
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The operating system of the ECS instance image (default: "Ubuntu")
* `associate_eip_address` - The EIP address to associate with the ECS instance (default: "")
* `eip_type` - The type of the EIP (default: "5_bgp")
* `bandwidth_size` - The size of the bandwidth (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")
* `instance_user_data` - The user data script for the ECS instance (default: "")
* `instance_system_disk_type` - The system disk type for the ECS instance (default: "SSD")
* `instance_system_disk_size` - The system disk size in GB for the ECS instance (default: 40)
* `security_group_ids` - The list of security group IDs for the ECS instance (default: [])

-> Either `security_group_name` or `security_group_ids` must be specified. If both are specified, `security_group_ids`
  takes precedence. If `security_group_ids` is specified, ensure that port `22` is enabled in the security group
  rules under it.

## SSH Connection Methods

### Using Private Key to Connect to ECS Instance

For `huaweicloud_compute_instance` resource, if `user_data` is specified, `admin_pass` setting will not take effect.
Therefore, we need to use a key-pair to log in to the ECS instance.

```hcl
resource "null_resource" "test" {
  provisioner "remote-exec" {
    connection {
      user        = "root"
      private_key = file(var.private_key_path)
      host        = huaweicloud_compute_eip_associate.test.public_ip
    }

    inline = var.instance_remote_exec_inline
  }
}
```

### Using Password to Connect to ECS Instance

If you just want to log in to an ECS instance and the instance is constructed with a password, please use the following
configuration:

```hcl
resource "null_resource" "test" {
  provisioner "remote-exec" {
    connection {
      user     = "root"
      password = "your password"
      host     = huaweicloud_compute_eip_associate.test.public_ip
    }

    inline = var.instance_remote_exec_inline
  }
}
```

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                    = "your_vpc_name"
  subnet_name                 = "your_subnet_name"
  security_group_name         = "your_security_group_name"
  keypair_name                = "your_keypair_name"
  private_key_path            = "your_private_key_path"
  instance_name               = "your_instance_name"
  bandwidth_name              = "your_bandwidth_name"
  instance_remote_exec_inline = "your_remote_exec_inline""
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the ECS instance takes about few minutes
* After the creation is successful, the provisioner starts to run
* The private key file must be accessible and have proper permissions
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.57.0 |
