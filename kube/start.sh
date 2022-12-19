#! /bin/sh

#preflight checks
  echo "\n-------------------------Preflight checks start-------------------------\n"

  # 1 - check kubectl is installed and a kubernetes cluster is accessible by the kubectl
  if kubectl version; then
  echo "kubectl - OK "
  else
  echo "\n\n*Please install kubectl/Verify that kubernetes cluster is up and running and accessible with kubectl\n\n"
  exit 1
  fi

  echo "\n-------------------------Preflight checks end -------------------------\n"

  #Installation section
  echo "\n-------------------------Starting query aggregator deployment-------------------------\n"
  echo "\nCreating service,deployment\n"
  kubectl create -f query-aggregator-service.yaml
  echo "\n======================================================================================================\n"
  echo "\nGetting status of pods\n"
  kubectl get pods
  echo "\n======================================================================================================\n"
  echo "\nCluster info\n"
  kubectl cluster-info
  echo "\n======================================================================================================\n"
  echo "\nService info\n"
  kubectl get svc
  echo "\n======================================================================================================\n"