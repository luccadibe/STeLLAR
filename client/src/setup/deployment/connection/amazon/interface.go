// MIT License
//
// Copyright (c) 2020 Theodor Amariucai
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package amazon

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"vhive-bench/client/util"
)

const (
	awsRegion = "us-west-1"
	s3Bucket  = "benchmarking-aws"
)

//AWSSingleton is an object used to interact with AWS through the methods it exports.
var AWSSingleton *instance

type instance struct {
	localZip []byte
	// S3Key is the bucket location in which this specific deployment will be uploaded
	S3Key         string
	NamePrefix    string
	region        string
	stage         string
	session       *session.Session
	s3Svc         *s3.S3
	lambdaSvc     *lambda.Lambda
	apiGatewaySvc *apigateway.APIGateway
	apiTemplate   []byte
}

//InitializeSingleton will create a new Amazon instance to interact with different AWS services.
func InitializeSingleton() {
	sessionInstance := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory location: %s", err.Error())
	}

	var apiTemplatePath string
	if strings.Contains(path, "connection") {
		apiTemplatePath = "../../../vHive-API-template-prod-oas30-apigateway.json"
	} else {
		apiTemplatePath = "./vHive-API-template-prod-oas30-apigateway.json"
	}

	apiTemplateFile := util.ReadFile(apiTemplatePath)
	apiTemplateByteValue, err := ioutil.ReadAll(apiTemplateFile)
	if err != nil {
		log.Fatalf("Could not read API template JSON when initializing AWS connection: %s", err.Error())
	}

	AWSSingleton = &instance{
		NamePrefix:    "vHive_",
		region:        awsRegion,
		stage:         "prod",
		session:       sessionInstance,
		lambdaSvc:     lambda.New(sessionInstance),
		apiGatewaySvc: apigateway.New(sessionInstance),
		s3Svc:         s3.New(sessionInstance),
		apiTemplate:   apiTemplateByteValue,
	}
}

//UploadZIPToS3 helps get around the 50MB image size limit for AWS functions.
func UploadZIPToS3(localZipPath string, sizeMB float64) {
	log.Infof(`Deploying to AWS and package size (~%vMB) > 50 MB, will now attempt to upload to Amazon S3.`, sizeMB)
	AWSSingleton.S3Key = fmt.Sprintf("benchmarking%vMB.zip", sizeMB)

	if _, err := AWSSingleton.s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(AWSSingleton.S3Key),
	}); err == nil {
		log.Infof("Object %q was already found in S3 bucket %q, skipping upload.", AWSSingleton.S3Key, s3Bucket)
		return
	}

	zipFile, err := os.Open(localZipPath)
	if err != nil {
		log.Fatalf("Failed to open zip file %q: %v", localZipPath, err)
	}

	uploadOutput, err := s3manager.NewUploader(AWSSingleton.session).Upload(&s3manager.UploadInput{
		Bucket: aws.String("benchmarking-aws"),
		Key:    aws.String(AWSSingleton.S3Key),
		Body:   zipFile,
	})
	if err != nil {
		log.Fatalf("Unable to upload %q to %q, %v", AWSSingleton.S3Key, s3Bucket, err.Error())
	}

	log.Infof("Successfully uploaded %q to bucket %q (%s)", AWSSingleton.S3Key, s3Bucket, uploadOutput.Location)
}

//SetLocalZip sets the location of the zipped binary file for the function to be deployed.
func SetLocalZip(path string) {
	zipBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read local zipped binary: %s", err.Error())
	}
	AWSSingleton.localZip = zipBytes
}
