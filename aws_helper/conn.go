package aws_helper

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
)

type AWSClient struct {
	stsconn *sts.STS
	S3conn  *s3.S3
	Dynconn *dynamodb.DynamoDB
	region  string
}

type Config struct {
	Region        string
	Use_mfa       bool
	Mfa_device_id string
	Mfa_token     string
	AWSAccount    Account
}

type Account struct {
	Profile   string            `yaml:"profile"`
	RoamRole  string            `yaml:"roam-role"`
	AccountID string            `yaml:"account_id"`
	Secondary Secondary_account `yaml:"secondary"`
}

type Secondary_account struct {
	Account_id string `yaml:"id"`
	Role       string `yaml:"role"`
	Region     string `yaml:"region"`
}

func (c *Config) Connect() interface{} {

	var client AWSClient

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SECURITY_TOKEN")
	os.Unsetenv("AWS_SESSION_TOKEN")
	os.Unsetenv("AWS_DEFAULT_REGION")

	os.Unsetenv("TF_VAR_secondary_access_key_id")
	os.Unsetenv("TF_VAR_secondary_secret_access_key")
	os.Unsetenv("TF_VAR_secondary_security_token")
	os.Unsetenv("TF_VAR_secondary_session_token")
	os.Unsetenv("TF_VAR_secondary_region")

	screds := &credentials.SharedCredentialsProvider{Profile: c.AWSAccount.Profile}

	log.Printf("[INFO] Using aws shared credentials profile: %s\n", c.AWSAccount.Profile)

	awsConfig := &aws.Config{
		Credentials: credentials.NewCredentials(screds),
		Region:      aws.String(c.Region),
		MaxRetries:  aws.Int(3),
	}

	sess := session.New(awsConfig)

	if len(c.AWSAccount.AccountID) > 0 && len(c.AWSAccount.RoamRole) > 0 {

		log.Println("[INFO] Initializing STS Connection")
		client.stsconn = sts.New(sess)

		params := &sts.AssumeRoleInput{}

		if c.Use_mfa {

			params = &sts.AssumeRoleInput{
				RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", c.AWSAccount.AccountID, c.AWSAccount.RoamRole)),
				RoleSessionName: aws.String(fmt.Sprintf("%s-%s", c.AWSAccount.AccountID, c.AWSAccount.RoamRole)),
				DurationSeconds: aws.Int64(3600),
				SerialNumber:    aws.String(c.Mfa_device_id),
				TokenCode:       aws.String(c.Mfa_token),
			}

		} else {

			params = &sts.AssumeRoleInput{
				RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", c.AWSAccount.AccountID, c.AWSAccount.RoamRole)),
				RoleSessionName: aws.String(fmt.Sprintf("%s-%s", c.AWSAccount.AccountID, c.AWSAccount.RoamRole)),
				DurationSeconds: aws.Int64(3600),
			}

		}

		sts_resp, sts_err := client.stsconn.AssumeRole(params)

		if sts_err != nil {
			log.Fatalf("Unable to assume role: %v", sts_err.Error())
		}

		os.Setenv("AWS_ACCESS_KEY_ID", *sts_resp.Credentials.AccessKeyId)
		os.Setenv("AWS_SECRET_ACCESS_KEY", *sts_resp.Credentials.SecretAccessKey)
		os.Setenv("AWS_SECURITY_TOKEN", *sts_resp.Credentials.SessionToken)
		os.Setenv("AWS_SESSION_TOKEN", *sts_resp.Credentials.SessionToken)
		os.Setenv("AWS_DEFAULT_REGION", c.Region)

		if c.AWSAccount.Secondary != (Secondary_account{}) {

			params2 := &sts.AssumeRoleInput{
				RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", c.AWSAccount.Secondary.Account_id, c.AWSAccount.Secondary.Role)),
				RoleSessionName: aws.String(fmt.Sprintf("%s-%s", c.AWSAccount.Secondary.Account_id, c.AWSAccount.Secondary.Role)),
				DurationSeconds: aws.Int64(3600),
			}

			stssec_resp, stssec_err := client.stsconn.AssumeRole(params2)

			if stssec_err != nil {
				log.Fatalf("Unable to assume role: %v", stssec_err.Error())
			}

			os.Setenv("TF_VAR_secondary_access_key_id", *stssec_resp.Credentials.AccessKeyId)
			os.Setenv("TF_VAR_secondary_secret_access_key", *stssec_resp.Credentials.SecretAccessKey)
			os.Setenv("TF_VAR_secondary_security_token", *stssec_resp.Credentials.SessionToken)
			os.Setenv("TF_VAR_secondary_session_token", *stssec_resp.Credentials.SessionToken)
			os.Setenv("TF_VAR_secondary_region", c.AWSAccount.Secondary.Region)

		}

		return c.assumeConnect(sts_resp)

	} else {

		profile_creds := credentials.Value{}
		var profile_err error

		if profile_creds, profile_err = screds.Retrieve(); profile_err != nil {
			log.Fatalf("[ERROR] Failed to get aws credentials for profile: %s with error: %s\n", c.AWSAccount.Profile, profile_err.Error())
		}

		os.Setenv("AWS_ACCESS_KEY_ID", profile_creds.AccessKeyID)
		os.Setenv("AWS_SECRET_ACCESS_KEY", profile_creds.SecretAccessKey)
		if len(profile_creds.SessionToken) > 0 {
			os.Setenv("AWS_SECURITY_TOKEN", profile_creds.SessionToken)
			os.Setenv("AWS_SESSION_TOKEN", profile_creds.SessionToken)
		}
		os.Setenv("AWS_DEFAULT_REGION", c.Region)
	}

	log.Println("[INFO] Initializing S3 Connection")
	client.S3conn = s3.New(sess)
	client.Dynconn = dynamodb.New(sess)

	return &client

}

func (c *Config) assumeConnect(sts *sts.AssumeRoleOutput) interface{} {

	var client AWSClient

	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(*sts.Credentials.AccessKeyId, *sts.Credentials.SecretAccessKey, *sts.Credentials.SessionToken),
		Region:      aws.String(c.Region),
		MaxRetries:  aws.Int(3),
	}

	sess := session.New(awsConfig)

	log.Println("[INFO] Initializing S3 Connection")
	client.S3conn = s3.New(sess)
	client.Dynconn = dynamodb.New(sess)

	return &client

}
