#!/bin/bash -e

TOKEN=$(kubectl get secret kube-job-assistant-token -n kja-test-deploy -o jsonpath="{.data.token}" | base64 -d)
CLUSTER_NAME=$(kubectl config view -o jsonpath='{.clusters[0].name}')
CLUSTER_SERVER=$(kubectl config view -o jsonpath='{.clusters[0].cluster.server}')
CLUSTER_CA=$(kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}')

cat <<EOF
apiVersion: v1
kind: Config
users:
- name: kube-job-assistant
  user:
    token: ${TOKEN}
clusters:
- name: ${CLUSTER_NAME}
  cluster:
    certificate-authority-data: ${CLUSTER_CA}
    server: ${CLUSTER_SERVER}
contexts:
- name: kube-job-assistant-context
  context:
    cluster: ${CLUSTER_NAME}
    user: kube-job-assistant
    namespace: kja-test-deploy
current-context: kube-job-assistant-context
EOF