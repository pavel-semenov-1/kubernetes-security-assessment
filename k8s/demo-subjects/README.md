# demo-subjects deployment

We start by creating a namespace

    kubectl create ns demo-subjects

Switch to the newly created namespace

    kubectl config set-context --current --namespace=demo-subjects

Build the backend

     docker build docker/demo-subjects/backend -t demo-subjects/backend

Build the frontend

    docker build docker/demo-subjects/frontend -t demo-subjects/frontend

Create a secret

    kubectl create secret generic datasource-secret --from-literal=password=welcome1

Deploy the application

    helm install demo-subjects k8s/demo-subjects