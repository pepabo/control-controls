# control-controls

control-controls control controls of AWS Security Hub across all regions.

## Usage

Export current security standards controls as a controls.yml.

``` console
$ control-controls export > controls.yml
2022-04-14T15:08:59+09:00 INF Fetching controls from eu-north-1
2022-04-14T15:09:04+09:00 INF Fetching controls from ap-south-1
2022-04-14T15:09:07+09:00 INF Fetching controls from eu-west-3
2022-04-14T15:09:12+09:00 INF Fetching controls from eu-west-2
2022-04-14T15:09:16+09:00 INF Fetching controls from eu-west-1
2022-04-14T15:09:21+09:00 INF Fetching controls from ap-northeast-3
2022-04-14T15:09:22+09:00 INF Fetching controls from ap-northeast-2
2022-04-14T15:09:24+09:00 INF Fetching controls from ap-northeast-1
2022-04-14T15:09:25+09:00 INF Fetching controls from sa-east-1
2022-04-14T15:09:30+09:00 INF Fetching controls from ca-central-1
2022-04-14T15:09:34+09:00 INF Fetching controls from ap-southeast-1
2022-04-14T15:09:36+09:00 INF Fetching controls from ap-southeast-2
2022-04-14T15:09:39+09:00 INF Fetching controls from eu-central-1
2022-04-14T15:09:43+09:00 INF Fetching controls from us-east-1
2022-04-14T15:09:47+09:00 INF Fetching controls from us-east-2
2022-04-14T15:09:50+09:00 INF Fetching controls from us-west-1
2022-04-14T15:09:53+09:00 INF Fetching controls from us-west-2
$
```

<details>

<summary>exported controls.yml is here</summary>

``` yaml
autoEnable: true
standards:
  aws-foundational-security-best-practices/v/1.0.0:
    enable: true
    controls:
      enable: [APIGateway.5, AutoScaling.1, AutoScaling.2, CloudTrail.1, CloudTrail.2, CloudTrail.4, CloudTrail.5, Config.1, DynamoDB.1, EC2.19, EC2.2, EC2.21, EC2.6, ECR.3, ELB.10, ELB.5, ELB.7, ES.4, ES.5, ES.6, ES.7, ES.8, IAM.1, IAM.2, IAM.3, IAM.5, IAM.6, IAM.7, IAM.8, NetworkFirewall.6, RDS.11, RDS.17, RDS.18, RDS.19, RDS.2, RDS.20, RDS.21, RDS.22, RDS.23, RDS.25, RDS.3, RDS.5, Redshift.4, Redshift.6, Redshift.8, S3.1, S3.10, S3.11, S3.12, S3.2, S3.3, S3.4, S3.5, S3.6, S3.9, SQS.1, SSM.1, SSM.4]
  cis-aws-foundations-benchmark/v/1.2.0:
    enable: true
    controls:
      enable: [CIS.1.1, CIS.1.10, CIS.1.11, CIS.1.13, CIS.1.14, CIS.1.16, CIS.1.2, CIS.1.22, CIS.1.3, CIS.1.4, CIS.1.5, CIS.1.6, CIS.1.7, CIS.1.8, CIS.1.9, CIS.2.1, CIS.2.2, CIS.2.3, CIS.2.4, CIS.2.5, CIS.2.6, CIS.2.7, CIS.2.8, CIS.2.9, CIS.3.1, CIS.3.10, CIS.3.11, CIS.3.12, CIS.3.13, CIS.3.14, CIS.3.2, CIS.3.3, CIS.3.4, CIS.3.5, CIS.3.6, CIS.3.7, CIS.3.8, CIS.3.9, CIS.4.3]
  pci-dss/v/3.2.1:
    enable: false
regions:
  ap-northeast-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  ap-northeast-2:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  ap-northeast-3:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [RDS.16, RDS.24]
  ap-south-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  ap-southeast-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  ap-southeast-2:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  ca-central-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  eu-central-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  eu-north-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  eu-west-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  eu-west-2:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  eu-west-3:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  sa-east-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.13, RDS.4, RDS.6, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  us-east-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CloudFront.1, CloudFront.2, CloudFront.3, CloudFront.4, CloudFront.5, CloudFront.6, CloudFront.7, CloudFront.8, CloudFront.9, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4, WAF.1]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  us-east-2:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  us-west-1:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
  us-west-2:
    standards:
      aws-foundational-security-best-practices/v/1.0.0:
        controls:
          enable: [ACM.1, APIGateway.1, APIGateway.2, APIGateway.3, APIGateway.4, Autoscaling.5, CodeBuild.1, CodeBuild.2, CodeBuild.4, CodeBuild.5, DMS.1, DynamoDB.2, DynamoDB.3, EC2.1, EC2.10, EC2.15, EC2.16, EC2.17, EC2.18, EC2.20, EC2.22, EC2.3, EC2.4, EC2.7, EC2.8, EC2.9, ECS.1, ECS.2, EFS.1, EFS.2, ELB.2, ELB.3, ELB.4, ELB.6, ELB.8, ELB.9, ELBv2.1, EMR.1, ES.1, ES.2, ES.3, ElasticBeanstalk.1, ElasticBeanstalk.2, GuardDuty.1, IAM.21, IAM.4, KMS.1, KMS.2, KMS.3, Lambda.1, Lambda.2, Lambda.5, Opensearch.1, Opensearch.2, Opensearch.3, Opensearch.4, Opensearch.5, Opensearch.6, Opensearch.8, RDS.1, RDS.10, RDS.12, RDS.13, RDS.14, RDS.15, RDS.16, RDS.24, RDS.4, RDS.6, RDS.7, RDS.8, RDS.9, Redshift.1, Redshift.2, Redshift.3, Redshift.7, S3.8, SNS.1, SSM.2, SSM.3, SageMaker.1, SecretsManager.1, SecretsManager.2, SecretsManager.3, SecretsManager.4]
      cis-aws-foundations-benchmark/v/1.2.0:
        controls:
          enable: [CIS.1.12, CIS.1.20, CIS.4.1, CIS.4.2]
```

</details>

For example, disable controls (Redshift.4, Redshift.6, Redshift.8).

``` yaml
autoEnable: true
standards:
  aws-foundational-security-best-practices/v/1.0.0:
    enable: true
    controls:
      enable: [APIGateway.5, AutoScaling.1, AutoScaling.2, CloudTrail.1, CloudTrail.2, CloudTrail.4, CloudTrail.5, Config.1, DynamoDB.1, EC2.19, EC2.2, EC2.21, EC2.6, ECR.3, ELB.10, ELB.5, ELB.7, ES.4, ES.5, ES.6, ES.7, ES.8, IAM.1, IAM.2, IAM.3, IAM.5, IAM.6, IAM.7, IAM.8, NetworkFirewall.6, RDS.11, RDS.17, RDS.18, RDS.19, RDS.2, RDS.20, RDS.21, RDS.22, RDS.23, RDS.25, RDS.3, RDS.5, S3.1, S3.10, S3.11, S3.12, S3.2, S3.3, S3.4, S3.5, S3.6, S3.9, SQS.1, SSM.1, SSM.4]
      disable: [Redshift.4, Redshift.6, Redshift.8]
[...]
```

Dry run.

``` console
$ control-controls plan controls.yml
2022-04-14T15:16:54+09:00 INF Checking eu-north-1
2022-04-14T15:17:02+09:00 INF Checking ap-south-1
2022-04-14T15:17:08+09:00 INF Checking eu-west-3
2022-04-14T15:17:15+09:00 INF Checking eu-west-2
2022-04-14T15:17:23+09:00 INF Checking eu-west-1
2022-04-14T15:17:31+09:00 INF Checking ap-northeast-3
2022-04-14T15:17:34+09:00 INF Checking ap-northeast-2
2022-04-14T15:17:37+09:00 INF Checking ap-northeast-1
2022-04-14T15:17:40+09:00 INF Checking sa-east-1
2022-04-14T15:17:49+09:00 INF Checking ca-central-1
2022-04-14T15:17:55+09:00 INF Checking ap-southeast-1
2022-04-14T15:17:59+09:00 INF Checking ap-southeast-2
2022-04-14T15:18:05+09:00 INF Checking eu-central-1
2022-04-14T15:18:13+09:00 INF Checking us-east-1
2022-04-14T15:18:19+09:00 INF Checking us-east-2
2022-04-14T15:18:25+09:00 INF Checking us-west-1
2022-04-14T15:18:31+09:00 INF Checking us-west-2
- eu-north-1::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.4
- eu-north-1::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.6
- eu-north-1::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.8
- ap-south-1::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.4
- ap-south-1::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.6
[...]
- us-west-1::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.6
- us-west-1::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.8
- us-west-2::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.4
- us-west-2::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.6
- us-west-2::standards::aws-foundational-security-best-practices/v/1.0.0::controls::Redshift.8

Plan: 0 to enable, 51 to disable
```

Apply changes.

``` console
$ control-controls apply controls.yml --disabled-reason 'Redshift is not running.'
2022-04-14T15:43:37+09:00 INF Applying to eu-north-1
2022-04-14T15:43:46+09:00 INF Disable control Control=Redshift.4 Reason="Redshift is not running." Region=eu-north-1 Standard=aws-foundational-security-best-practice
s/v/1.0.0                                                                                                                     
2022-04-14T15:43:47+09:00 INF Disable control Control=Redshift.6 Reason="Redshift is not running." Region=eu-north-1 Standard=aws-foundational-security-best-practice
s/v/1.0.0                                                                                                                     
2022-04-14T15:43:49+09:00 INF Disable control Control=Redshift.8 Reason="Redshift is not running." Region=eu-north-1 Standard=aws-foundational-security-best-practice
s/v/1.0.0                                                                                                                     
2022-04-14T15:43:51+09:00 INF Applying to ap-south-1
2022-04-14T15:43:56+09:00 INF Disable control Control=Redshift.4 Reason="Redshift is not running." Region=ap-south-1 Standard=aws-foundational-security-best-practice
s/v/1.0.0                                                                                                                     
2022-04-14T15:43:57+09:00 INF Disable control Control=Redshift.6 Reason="Redshift is not running." Region=ap-south-1 Standard=aws-foundational-security-best-practice
s/v/1.0.0                                                                                                                     
[...]
2022-04-14T15:46:18+09:00 INF Disable control Control=Redshift.6 Reason="Redshift is not running." Region=us-west-1 Standard=aws-foundational-security-best-practices
/v/1.0.0                                                                                                                     
2022-04-14T15:46:19+09:00 INF Disable control Control=Redshift.8 Reason="Redshift is not running." Region=us-west-1 Standard=aws-foundational-security-best-practices
/v/1.0.0                                                                                                                     
2022-04-14T15:46:20+09:00 INF Applying to us-west-2
2022-04-14T15:46:26+09:00 INF Disable control Control=Redshift.4 Reason="Redshift is not running." Region=us-west-2 Standard=aws-foundational-security-best-practices
/v/1.0.0                                                                                                                     
2022-04-14T15:46:27+09:00 INF Disable control Control=Redshift.6 Reason="Redshift is not running." Region=us-west-2 Standard=aws-foundational-security-best-practices
/v/1.0.0                                                                                                                     
2022-04-14T15:46:29+09:00 INF Disable control Control=Redshift.8 Reason="Redshift is not running." Region=us-west-2 Standard=aws-foundational-security-best-practices
/v/1.0.0                                                                                                                     

Apply complete
```

## Install

**homebrew tap:**

```console
$ brew install pepabo/tap/control-controls
```

**manually:**

Download binany from [releases page](https://github.com/pepabo/control-controls/releases)

**go install:**

```console
$ go install github.com/pepabo/control-controls@latest
```
