key_alias             = "tf-test-obs-key"
bucket_name           = "tf-test-obs-bucket"
object_name           = "tf-test-obs-object"
object_upload_content = <<EOT
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
EOT
target_bucket_configuration = {
  bucket = "tf-test-obs-bucket-target"
}
bandwidth_policy_configurations = [
  {
    max_bandwidth = 1
    start         = "03:00"
    end           = "04:00"
  }
]
