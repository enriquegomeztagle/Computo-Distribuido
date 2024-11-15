#!/bin/bash

IMAGE_NAME="proyecto-kubernetes"
JOB_NAME="proyecto-kubernetes-job"
DEPLOYMENT_NAME="heml-proyecto-kubernetes"
NAMESPACE="default"
CHART_NAME="heml-proyecto-kubernetes"
CHART_PATH="./heml-proyecto-kubernetes"

echo "### Eliminando recursos de Kubernetes ###"
kubectl delete job $JOB_NAME --ignore-not-found
kubectl delete deployment $DEPLOYMENT_NAME --ignore-not-found
kubectl delete pod --all --namespace=$NAMESPACE --ignore-not-found
kubectl delete svc $DEPLOYMENT_NAME --ignore-not-found

echo "### Eliminando imagen Docker local ###"
docker rmi -f $IMAGE_NAME

echo "### Reconstruyendo imagen Docker ###"
docker build -t $IMAGE_NAME .

echo "### Cargando imagen en Kind ###"
kind load docker-image $IMAGE_NAME

echo "### Verificando estructura de Helm Chart ###"
if [ ! -d "$CHART_PATH" ]; then
  echo "### Creando Helm Chart ###"
  helm create $CHART_PATH
  rm -rf $CHART_PATH/templates/{service.yaml,serviceaccount.yaml,deployment.yaml,hpa.yaml,ingress.yaml,tests/test-connection.yaml}
fi

echo "### Configurando archivos de Helm Chart ###"
cat <<EOF > $CHART_PATH/templates/job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Values.job.name }}
spec:
  template:
    spec:
      containers:
      - name: {{ .Values.job.name }}
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
        imagePullPolicy: {{ .Values.image.pullPolicy }}
        command: {{ .Values.job.command | toJson }}
      restartPolicy: {{ .Values.job.restartPolicy }}
EOF

cat <<EOF > $CHART_PATH/values.yaml
image:
  repository: $IMAGE_NAME
  tag: latest
  pullPolicy: Never

job:
  name: $JOB_NAME
  command: ["bash", "-c", "./run_tests.sh"]
  restartPolicy: Never
EOF

echo "### Desplegando Helm Chart ###"
helm upgrade --install $CHART_NAME $CHART_PATH

echo "### Verificando estado de los recursos ###"
while true; do
  kubectl get jobs
  kubectl get pods
  echo "----------------------------------------"
  sleep 5
done &

while true; do
  JOB_STATUS=$(kubectl get job $JOB_NAME -o jsonpath='{.status.conditions[?(@.type=="Complete")].status}' 2>/dev/null)
  if [ "$JOB_STATUS" == "True" ]; then
    echo "### El Job ha finalizado exitosamente ###"
    break
  fi
  sleep 5
done

echo "### Mostrando logs del Job ###"
kubectl logs -l job-name=$JOB_NAME

echo "### Ejecutando prueba de conexión de Helm ###"
helm test $CHART_NAME
