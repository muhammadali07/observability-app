# observability-app

git cloning repository
cd observability-app
cd terraform
terraform init
terraform apply
cd ..
cd myapp
go mod tidy
kubectl get pods -n monitoring
kubectl port-forward svc/tempo 4318 -n monitoring
export OTLP_ENDPOINT=localhost:4318
curl http://0.0.0.0:8080/devices
kubectl port-forward svc/grafana 3000:80 -n monitoring
# Datasource URL: http://tempo.monitoring:3100
kubectl apply -f k8s
kubectl get pods -n default
kubectl get svc -n default
kubectl port-forward svc/myapp 8080:8080 -n default
