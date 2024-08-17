# kubernetes-security-assessment
GitHub repo for my Diploma Thesis

## demo-users deployment
We start by creating a namespace

    kubectl create ns demo-users

Switch to the newly created namespace

    kubectl config set-context --current --namespace=demo-users

Build the backend

     docker build docker/demo-users/backend -t demo-users/backend

Build the frontend

    docker build docker/demo-users/frontend -t demo-users/frontend

Create a secret

    kubectl create secret generic datasource-secret --from-literal=password=welcome1

Deploy the application

    helm install demo-users k8s/demo-users