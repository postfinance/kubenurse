# kubenurse

kubenurse is a little service that monitors all network connections in a kubernetes cluster and exports the taken metrics as prometheus endpoint.

## What is this chart about?

This version of the Helm chart uses Traefik version <b>2.x</b> that includes
Middlewares and IngressRoutes, instead of Ingress and backend/frontend logic
that existed in version <b>1.x</b>.

What has changed in Traefik version 2.x:
https://doc.traefik.io/traefik/migration/v1-to-v2/

Besides that, this chart is also designed to enable access to kubenurse via a
subpath. For example, `dummy.com/kubenurse/*`.