#!/bin/bash
instance_name=$1

instance=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=${instance_name}" --output text --query 'Reservations[*].Instances[*].InstanceId')
aws ec2 stop-instances --instance-id "$instance" --output=text
while [[ $(aws ec2 describe-instances --instance-id "$instance" --query "Reservations[*].Instances[*].State.Name" --output=text) != "stopped" ]]; do echo "Stopping ${instance_name}..." && sleep 5; done
echo "Stopped ${instance_name}!"