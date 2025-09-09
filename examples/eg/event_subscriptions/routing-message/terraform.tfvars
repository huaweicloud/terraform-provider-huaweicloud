vpc_name              = "tf_test_vpc"
subnet_name           = "tf_test_subnet"
security_group_name   = "tf_test_security_group"
bucket_name           = "tf-test-bucket"
instance_name         = "tf_test_kafka"
topic_name            = "tf-test-topic"
connection_name       = "tf-test-connect"
object_name           = "tf-test-obs-object"
object_upload_content = <<EOT
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
EOT
