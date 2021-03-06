/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pod

import (
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"

	"github.com/kubernetes-sigs/multi-tenancy/incubator/virtualcluster/pkg/syncer/constants"
	"github.com/kubernetes-sigs/multi-tenancy/incubator/virtualcluster/pkg/syncer/conversion"
	"github.com/kubernetes-sigs/multi-tenancy/incubator/virtualcluster/pkg/syncer/metrics"
	"github.com/kubernetes-sigs/multi-tenancy/incubator/virtualcluster/pkg/syncer/reconciler"
	"github.com/kubernetes-sigs/multi-tenancy/incubator/virtualcluster/pkg/syncer/utils"
)

func (c *controller) StartDWS(stopCh <-chan struct{}) error {
	if !cache.WaitForCacheSync(stopCh, c.podSynced, c.serviceSynced, c.nsSynced) {
		return fmt.Errorf("failed to wait for caches to sync before starting Pod dws")
	}
	return c.multiClusterPodController.Start(stopCh)
}

func (c *controller) Reconcile(request reconciler.Request) (reconciler.Result, error) {
	klog.Infof("reconcile pod %s/%s %s event for cluster %s", request.Namespace, request.Name, request.Event, request.Cluster.Name)
	vPod := request.Obj.(*v1.Pod)
	c.updateClusterVNodePodMap(request.Cluster.Name, vPod, request.Event)

	var operation string
	switch request.Event {
	case reconciler.AddEvent:
		operation = "pod_add"
		defer recordOperation(operation, time.Now())
		err := c.reconcilePodCreate(request.Cluster.Name, request.Namespace, request.Name, vPod)
		recordError(operation, err)
		if err != nil {
			klog.Errorf("failed reconcile pod %s/%s CREATE of cluster %s %v", request.Namespace, request.Name, request.Cluster.Name, err)
			return reconciler.Result{Requeue: true}, err
		}
	case reconciler.UpdateEvent:
		operation = "pod_update"
		defer recordOperation(operation, time.Now())
		err := c.reconcilePodUpdate(request.Cluster.Name, request.Namespace, request.Name, vPod)
		recordError(operation, err)
		if err != nil {
			klog.Errorf("failed reconcile pod %s/%s UPDATE of cluster %s %v", request.Namespace, request.Name, request.Cluster.Name, err)
			return reconciler.Result{Requeue: true}, err
		}
	case reconciler.DeleteEvent:
		operation = "pod_delete"
		defer recordOperation(operation, time.Now())
		err := c.reconcilePodRemove(request.Cluster.Name, request.Namespace, request.Name, vPod)
		recordError(operation, err)
		if err != nil {
			klog.Errorf("failed reconcile pod %s/%s DELETE of cluster %s %v", request.Namespace, request.Name, request.Cluster.Name, err)
			return reconciler.Result{Requeue: true}, err
		}
	}
	return reconciler.Result{}, nil
}

func isPodScheduled(pod *v1.Pod) bool {
	_, cond := podutil.GetPodCondition(&pod.Status, v1.PodScheduled)
	return cond != nil && cond.Status == v1.ConditionTrue
}

func createNotSupportEvent(pod *v1.Pod) *v1.Event {
	eventTime := metav1.Now()
	return &v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "syncer",
		},
		InvolvedObject: v1.ObjectReference{
			APIVersion: "v1",
			Kind:       "Pod",
			Name:       pod.Name,
		},
		Type:                "Warning",
		Reason:              "NotSupported",
		Message:             "The Pod has nodeName set in the spec which is not supported for now",
		FirstTimestamp:      eventTime,
		LastTimestamp:       eventTime,
		ReportingController: "syncer",
	}
}

func (c *controller) reconcilePodCreate(cluster, namespace, name string, vPod *v1.Pod) error {
	// load deleting pod, don't create any pod on super master.
	if vPod.DeletionTimestamp != nil {
		return c.reconcilePodUpdate(cluster, namespace, name, vPod)
	}

	targetNamespace := conversion.ToSuperMasterNamespace(cluster, namespace)
	_, err := c.podLister.Pods(targetNamespace).Get(name)
	if err == nil {
		return c.reconcilePodUpdate(cluster, namespace, name, vPod)
	}

	if vPod.Spec.NodeName != "" && !isPodScheduled(vPod) {
		// For now, we skip vPod that has NodeName set to prevent tenant from deploying DaemonSet or DaemonSet alike CRDs.
		tenantClient, err := c.multiClusterPodController.GetClusterClient(cluster)
		if err != nil {
			return fmt.Errorf("failed to create client from cluster %s config: %v", cluster, err)
		}
		event := createNotSupportEvent(vPod)
		vEvent := conversion.BuildVirtualPodEvent(cluster, event, vPod)
		_, err = tenantClient.CoreV1().Events(vPod.Namespace).Create(vEvent)
		return err
	}

	newObj, err := conversion.BuildMetadata(cluster, targetNamespace, vPod)
	if err != nil {
		return err
	}

	pPod := newObj.(*v1.Pod)

	// check if the secret in super master is ready we must create pod after sync the secret.
	saName := "default"
	if pPod.Spec.ServiceAccountName != "" {
		saName = pPod.Spec.ServiceAccountName
	}

	pSecret, err := utils.GetSecret(c.client, targetNamespace, saName)
	if err != nil {
		return fmt.Errorf("failed to get secret: %v", err)
	}

	if pSecret.Labels[constants.SyncStatusKey] != constants.SyncStatusReady {
		return fmt.Errorf("secret for pod is not ready")
	}

	tenantClient, err := c.multiClusterPodController.GetClusterClient(cluster)
	if err != nil {
		return err
	}
	vSecret, err := utils.GetSecret(tenantClient.CoreV1(), namespace, saName)
	if err != nil {
		return fmt.Errorf("failed to get secret: %v", err)
	}

	services, err := c.getPodRelatedServices(cluster, pPod)
	if err != nil {
		return fmt.Errorf("failed to list services from cluster %s cache: %v", cluster, err)
	}

	if len(services) == 0 {
		return fmt.Errorf("service is not ready")
	}

	nameServer, err := c.getClusterNameServer(c.client, cluster)
	if err != nil {
		return fmt.Errorf("failed to find nameserver: %v", err)
	}

	var ms = []conversion.PodMutator{
		conversion.PodMutateDefault(vPod, vSecret, pSecret, services, nameServer),
		conversion.PodMutateAutoMountServiceAccountToken(c.config.DisableServiceAccountToken),
		conversion.PodAddExtensionMeta(vPod),
	}

	if c.config.EnableTenantKubeConfig {
		clientConfig, err := c.multiClusterPodController.GetClusterClientConfig(cluster)
		if err != nil {
			return fmt.Errorf("failed to get cluster config")
		}

		kubeConfigBytes, err := createKubeConfigByServiceAccount(tenantClient.CoreV1().ServiceAccounts(vPod.Namespace), tenantClient.CoreV1().Secrets(vPod.Namespace), clientConfig, saName)
		if err != nil {
			return fmt.Errorf("failed to create kubeconfig from service account: %v", err)
		}

		secret := conversion.BuildKubeConfigSecret(cluster, vPod, kubeConfigBytes)
		secret.Namespace = targetNamespace

		_, err = c.client.Secrets(targetNamespace).Create(secret)
		if err != nil && !errors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create kubeconfig secret for pod: %v", err)
		}

		ms = append(ms, conversion.PodMutateKubeConfig(vPod, secret, c.config.TenantKubeConfigMountPath))
	}

	err = conversion.VC(c.multiClusterPodController, cluster).Pod(pPod).Mutate(ms...)
	if err != nil {
		return fmt.Errorf("failed to mutate pod: %v", err)
	}

	_, err = c.client.Pods(targetNamespace).Create(pPod)
	if errors.IsAlreadyExists(err) {
		klog.Infof("pod %s/%s of cluster %s already exist in super master", namespace, name, cluster)
		return nil
	}
	return err
}

func (c *controller) getClusterNameServer(client v1core.ServicesGetter, cluster string) (string, error) {
	svc, err := client.Services(conversion.ToSuperMasterNamespace(cluster, constants.TenantDNSServerNS)).Get(constants.TenantDNSServerServiceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return "", nil
		}
		return "", err
	}

	return svc.Spec.ClusterIP, nil
}

func (c *controller) getPodRelatedServices(cluster string, pPod *v1.Pod) ([]*v1.Service, error) {
	var services []*v1.Service
	list, err := c.serviceLister.Services(conversion.ToSuperMasterNamespace(cluster, metav1.NamespaceDefault)).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	services = append(services, list...)

	list, err = c.serviceLister.Services(pPod.Namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	services = append(services, list...)

	return services, nil
}

func (c *controller) reconcilePodUpdate(cluster, namespace, name string, vPod *v1.Pod) error {
	targetNamespace := conversion.ToSuperMasterNamespace(cluster, namespace)
	pPod, err := c.podLister.Pods(targetNamespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			// if the pod on super master has been deleted and syncer has not
			// deleted virtual pod with 0 grace period second successfully.
			// we depends on periodic check to do gc.
			return nil
		}
		return err
	}

	if vPod.DeletionTimestamp != nil {
		if pPod.DeletionTimestamp != nil {
			// pPod is under deletion, waiting for UWS bock populate the pod status.
			return nil
		}
		deleteOptions := metav1.NewDeleteOptions(*vPod.DeletionGracePeriodSeconds)
		deleteOptions.Preconditions = metav1.NewUIDPreconditions(string(pPod.UID))
		err = c.client.Pods(targetNamespace).Delete(name, deleteOptions)
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	updatedPod := conversion.CheckPodEquality(pPod, vPod)
	if updatedPod != nil {
		pPod, err = c.client.Pods(targetNamespace).Update(updatedPod)
		if err != nil {
			return err
		}
	}

	// pod has been updated by tenant controller
	if !equality.Semantic.DeepEqual(vPod.Status, pPod.Status) {
		c.enqueuePod(pPod)
	}

	return nil
}

func createKubeConfigByServiceAccount(saClient v1core.ServiceAccountInterface, secretClient v1core.SecretInterface, clientConfig clientcmd.ClientConfig, saName string) ([]byte, error) {
	serviceAccount, err := saClient.Get(saName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	for _, reference := range serviceAccount.Secrets {
		secret, err := secretClient.Get(reference.Name, metav1.GetOptions{})
		if err != nil {
			continue
		}

		if secret.Type == v1.SecretTypeServiceAccountToken {
			token, exists := secret.Data[v1.ServiceAccountTokenKey]
			if !exists {
				return nil, fmt.Errorf("service account token %q for service account %q did not contain token data", secret.Name, saName)
			}

			cfg := &rawConfig

			ctx := cfg.Contexts[cfg.CurrentContext]
			cfg.CurrentContext = saName
			cfg.Contexts = map[string]*clientcmdapi.Context{
				cfg.CurrentContext: ctx,
			}
			ctx.AuthInfo = saName
			cfg.AuthInfos = map[string]*clientcmdapi.AuthInfo{
				ctx.AuthInfo: {
					Token: string(token),
				},
			}
			cluster := cfg.Clusters[ctx.Cluster]
			cluster.Server = "https://kubernetes.default"
			cfg.Clusters = map[string]*clientcmdapi.Cluster{
				ctx.Cluster: cluster,
			}

			out, err := clientcmd.Write(*cfg)
			if err != nil {
				return nil, fmt.Errorf("failed to write serializes the config to yaml")
			}

			return out, nil
		}
	}

	return nil, fmt.Errorf("any available service account token not found")
}

func (c *controller) reconcilePodRemove(cluster, namespace, name string, vPod *v1.Pod) error {
	targetNamespace := conversion.ToSuperMasterNamespace(cluster, namespace)
	opts := &metav1.DeleteOptions{
		PropagationPolicy: &constants.DefaultDeletionPolicy,
	}
	err := c.client.Pods(targetNamespace).Delete(name, opts)
	if errors.IsNotFound(err) {
		klog.Warningf("pod %s/%s of cluster (%s) is not found in super master", namespace, name, cluster)
		return nil
	}
	return err
}

func recordOperation(operation string, start time.Time) {
	metrics.PodOperations.WithLabelValues(operation).Inc()
	metrics.PodOperationsDuration.WithLabelValues(operation).Observe(metrics.SinceInSeconds(start))
}

func recordError(operation string, err error) {
	if err != nil {
		metrics.PodOperationsErrors.WithLabelValues(operation).Inc()
	}
}
