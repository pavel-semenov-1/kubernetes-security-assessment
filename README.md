# kubernetes-security-assessment
GitHub repository for my Diploma Thesis

Clone the repository and follow the instructions below to deploy KSA dashboard.

## Rancher Desktop 

Rancher Desktop is simple to install. Just download the installation package for your OS from the official page: https://rancherdesktop.io/. Depending on the platform 

## Demo applications (optional)

[Demo Users app deployment](k8s/demo-users/README.md)

[Demo Subjects app deployment](k8s/demo-subjects/README.md)

## KSA Dashboard Build

In order to build all the images that the dashboard requires run the following script:

```
docker/rebuildKSA.sh
```

## KSA Dashboard Deploy

First, we need to create the namespace:
```
kubectl create ns ksa
```

Then, navigate to the `k8s/ksa` folder and execute the installation script:
```
cd k8s/ksa
./reinstall.sh
```

After a minute or so you should be able to access the dashboard at http://localhost:3123

## Troubleshooting

If the dashboard is not accessible, try to list the pods inside **ksa** namespace:
```
kubectl -n ksa get po
```

If the pod is not running, try to view its event by executing
```
kubectl -n ksa describe po <pod-name>
```

Otherwise you can view the pod logs by
```
kubectl -n ksa logs <pod-name>
```