#! /bin/sh
echo "\n-------------------------Deleting query aggregator deployments-------------------------\n"
echo "\n===================================================================================================\n"
echo "\nRemoving service,deployment"
kubectl delete -f query-aggregator-service.yaml
echo "\n===================================================================================================\n"
echo "\nVerifying the status of pods"
kubectl get pods
echo "\n===================================================================================================\n"