# VPC Setup

## Terraform VPC Module

Using terraform. Copy the following configuration into a new folder called `vpc`.

This will create a VPC as well as an S3 bucket, using `bucket_name`, that is accessible through the VPC using a VPC Endpoint to AWS S3. All clusters deployed into this VPC will use this bucket for state and configuration data. Cluster configuration is stored with the cluster name as a prefix.

```terraform
provider "aws" {
  region = "us-west-2"
}

module "vpc" {
  source = "github.com/coldog/k8pack/vpc"

  asset_bucket = "<bucket_name>"
  vpc_name     = "k8s"
}
```

Now run `terraform apply`.

Results:

```
vpc_id = "..."
```

If you are using a custom VPC you must replicate this behaviour of the VPC endpoint and create an S3 bucket to store configuration.
