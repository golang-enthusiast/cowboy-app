## Description
The cowboy app consists of:
- Main application, also known as "cowboy" workers (cmd/app). The workers just sit and wait for messages from the cron job to start the encounter. As soon as the encounter begins, the workers select a random target and send a "shoot" message. The encounter continues until only one winner remains. The winner sends a "winner" message to all consumers, which they print.
- migrations (/migrations). Migrate initial cowboy list to the DynamoDB.
- cron job (cmd/cron). Sends cowboys encounter message at the same time in parallel.

## Technologies
- Golang 1.16
- AWS DynamoDB (cowboys persistent storage)
- AWS SQS (communication between cowboys)
- Kubernetes (container orchestration solution) 
- Kubernetes CRD (cowboy custom resource definition)
- Helm/Skaffold (deployment)

## Pre requirements
- golang 1.16+
- minikube: https://minikube.sigs.k8s.io/docs/start
- helm: https://helm.sh
- skaffold: https://skaffold.dev
- kubebuilder https://github.com/kubernetes-sigs/kubebuilder

## How to start the application:
- Install CRD to k8s cluster:
    - git clone https://github.com/golang-enthusiast/distributed-cowboys
    - make install
    - make run
- Run the main application:
    - make run