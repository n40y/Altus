// scanner/collection/aws_collect.go
package collection

import (
	"altus/scanner/models"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// CollectAWSResources appelle les API AWS réelles selon le profil du terminal
func CollectAWSResources(cfg aws.Config, data *models.ScanResult) {
	ctx := context.Background()

	// ---- 1. COLLECTE DES UTILISATEURS IAM ----
	iamClient := iam.NewFromConfig(cfg)
	usersOutput, err := iamClient.ListUsers(ctx, &iam.ListUsersInput{})
	if err != nil {
		fmt.Printf("[!] Droits insuffisants ou erreur sur l'API IAM : %v\n", err)
	} else {
		for _, user := range usersOutput.Users {
			data.Nodes = append(data.Nodes, models.Node{
				Name: *user.UserName,
				Type: "IAM_USER",
			})
		}
	}

	// ---- 2. COLLECTE DES BUCKETS S3 ----
	s3Client := s3.NewFromConfig(cfg)
	bucketsOutput, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		fmt.Printf("[!] Droits insuffisants ou erreur sur l'API S3 : %v\n", err)
	} else {
		for _, bucket := range bucketsOutput.Buckets {
			data.Nodes = append(data.Nodes, models.Node{
				Name: *bucket.Name,
				Type: "S3_BUCKET",
			})
		}
	}

	// ---- 3. COLLECTE DES INSTANCES EC2 ----
	ec2Client := ec2.NewFromConfig(cfg)
	instancesOutput, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		fmt.Printf("[!] Droits insuffisants ou erreur sur l'API EC2 : %v\n", err)
	} else {
		for _, reservation := range instancesOutput.Reservations {
			for _, instance := range reservation.Instances {
				nodeName := *instance.InstanceId
				// Récupération optionnelle du tag "Name" de l'instance
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" && *tag.Value != "" {
						nodeName = *tag.Value
						break
					}
				}
				data.Nodes = append(data.Nodes, models.Node{
					Name: nodeName,
					Type: "EC2_INSTANCE",
				})
			}
		}
	}
}
