module vhive-bench/client

go 1.15

replace vhive-bench/client/util => ./util

replace vhive-bench/client/setup => ./setup

replace vhive-bench/client/experiments => ./experiments

replace vhive-bench/client/experiments/benchmarking => ./experiments/benchmarking

replace vhive-bench/client/experiments/visualization => ./experiments/visualization

replace vhive-bench/client/experiments/networking => ./experiments/networking

require (
	github.com/aws/aws-lambda-go v1.20.0
	github.com/aws/aws-sdk-go v1.35.20
	github.com/go-gota/gota v0.10.1
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/sirupsen/logrus v1.7.0
	gonum.org/v1/gonum v0.8.1
	gonum.org/v1/plot v0.8.0
	google.golang.org/grpc v1.33.2
)