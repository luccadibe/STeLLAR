#!/bin/bash
sudo apt-get update
sudo apt-get install tmux
sudo apt-get install awscli
mkdir -p "latency-samples"
sudo ulimit -n 16384
aws configure