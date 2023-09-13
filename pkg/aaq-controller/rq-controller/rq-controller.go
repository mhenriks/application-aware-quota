package rq_controller

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	_ "kubevirt.io/api/core/v1"
	"kubevirt.io/applications-aware-quota/pkg/client"
	"strings"

	v1alpha12 "kubevirt.io/applications-aware-quota/staging/src/kubevirt.io/applications-aware-quota-api/pkg/apis/core/v1alpha1"
	"kubevirt.io/client-go/log"
)

type enqueueState string

const (
	Immediate enqueueState = "Immediate"
	Forget    enqueueState = "Forget"
	BackOff   enqueueState = "BackOff"
	RQSuffix  string       = "-non-schedulable-resources-managed-rq-x"
)

type RQController struct {
	arqInformer cache.SharedIndexInformer
	rqInformer  cache.SharedIndexInformer
	arqQueue    workqueue.RateLimitingInterface
	aaqCli      client.AAQClient
	stop        <-chan struct{}
}

func NewRQController(aaqCli client.AAQClient,
	rqInformer cache.SharedIndexInformer,
	arqInformer cache.SharedIndexInformer,
	stop <-chan struct{},
) *RQController {
	ctrl := RQController{
		rqInformer:  rqInformer,
		aaqCli:      aaqCli,
		arqInformer: arqInformer,
		arqQueue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "arq-queue"),
		stop:        stop,
	}

	return &ctrl
}

func (ctrl *RQController) addAllArqsInNamespace(ns string) {
	objs, err := ctrl.arqInformer.GetIndexer().ByIndex(cache.NamespaceIndex, ns)
	if err != nil {
		log.Log.Infof("AaqGateController: Error failed to list pod from podInformer")
	}
	for _, obj := range objs {
		arq := obj.(*v1alpha12.ApplicationsResourceQuota)
		key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(arq)
		if err != nil {
			return
		}
		ctrl.arqQueue.Add(key)
	}
}

// When a ApplicationsResourceQuotas is deleted, enqueue all gated pods for revaluation
func (ctrl *RQController) deleteArq(obj interface{}) {
	arq := obj.(*v1alpha12.ApplicationsResourceQuota)
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(arq)
	if err != nil {
		return
	}
	ctrl.arqQueue.Add(key)
	return
}

// When a ApplicationsResourceQuotas is updated, enqueue all gated pods for revaluation
func (ctrl *RQController) addArq(obj interface{}) {
	arq := obj.(*v1alpha12.ApplicationsResourceQuota)
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(arq)
	if err != nil {
		return
	}
	ctrl.arqQueue.Add(key)
	return
}

// When a ApplicationsResourceQuotas is updated, enqueue all gated pods for revaluation
func (ctrl *RQController) updateArq(old, cur interface{}) {
	curArq := cur.(*v1alpha12.ApplicationsResourceQuota)
	oldArq := old.(*v1alpha12.ApplicationsResourceQuota)

	if !ResourceListEqual(curArq.Spec.Hard, oldArq.Spec.Hard) {
		key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(curArq)
		if err != nil {
			return
		}
		ctrl.arqQueue.Add(key)
	}

	return
}

func (ctrl *RQController) deleteRQ(obj interface{}) {
	rq := obj.(*v1.ResourceQuota)
	arq := v1alpha12.ApplicationsResourceQuota{
		ObjectMeta: metav1.ObjectMeta{Name: strings.TrimSuffix(rq.Name, RQSuffix)},
	}
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(arq)
	if err != nil {
		return
	}
	ctrl.arqQueue.Add(key)
	return
}
func (ctrl *RQController) updateRQ(old, curr interface{}) {
	curRq := curr.(*v1.ResourceQuota)
	oldRq := old.(*v1.ResourceQuota)
	if !ResourceListEqual(curRq.Spec.Hard, oldRq.Spec.Hard) {
		arq := v1alpha12.ApplicationsResourceQuota{
			ObjectMeta: metav1.ObjectMeta{Name: strings.TrimSuffix(curRq.Name, RQSuffix)},
		}
		key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(arq)
		if err != nil {
			return
		}
		ctrl.arqQueue.Add(key)
	}
	return
}

func (ctrl *RQController) runWorker() {
	for ctrl.Execute() {
	}
}

func (ctrl *RQController) Execute() bool {
	key, quit := ctrl.arqQueue.Get()
	if quit {
		return false
	}
	defer ctrl.arqQueue.Done(key)

	err, enqueueState := ctrl.execute(key.(string))
	if err != nil {
		klog.Errorf(fmt.Sprintf("AaqGateController: Error with key: %v err: %v", key, err))
	}
	switch enqueueState {
	case BackOff:
		ctrl.arqQueue.AddRateLimited(key)
	case Forget:
		ctrl.arqQueue.Forget(key)
	case Immediate:
		ctrl.arqQueue.Add(key)
	}

	return true
}

func (ctrl *RQController) execute(key string) (error, enqueueState) {
	arqNS, arqName, err := cache.SplitMetaNamespaceKey(key)
	arqObj, exists, err := ctrl.arqInformer.GetIndexer().GetByKey(arqNS + "/" + arqName)
	if err != nil {
		return err, Immediate
	} else if !exists {
		err = ctrl.aaqCli.CoreV1().ResourceQuotas(arqNS).Delete(context.Background(), arqName+RQSuffix, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return err, Immediate
		} else {
			return nil, Forget
		}
	}

	arq := arqObj.(*v1alpha12.ApplicationsResourceQuota)
	nonSchedulableResourcesLimitations := filterNonScheduableResources(arq.Spec.Hard)
	if len(nonSchedulableResourcesLimitations) == 0 {
		err = ctrl.aaqCli.CoreV1().ResourceQuotas(arqNS).Delete(context.Background(), arqName+RQSuffix, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			log.Log.Infof("wallak herer")
			return err, Immediate
		} else {
			return nil, Forget
		}
	}

	rqObj, exists, err := ctrl.rqInformer.GetIndexer().GetByKey(arq.Namespace + "/" + arq.Name + RQSuffix)
	if err != nil {
		return err, Immediate
	} else if !exists {
		rq := &v1.ResourceQuota{
			ObjectMeta: metav1.ObjectMeta{
				Name: arq.Name + RQSuffix,
				Labels: map[string]string{
					"aaq.managed.rq": "true",
				},
			},
			Spec: v1.ResourceQuotaSpec{
				Hard: nonSchedulableResourcesLimitations,
			},
		}
		rq, err = ctrl.aaqCli.CoreV1().ResourceQuotas(arqNS).Create(context.Background(), rq, metav1.CreateOptions{})
		if err != nil {
			return err, Immediate
		} else {
			return err, Forget
		}
	}
	rq := rqObj.(*v1.ResourceQuota)
	if ResourceListEqual(rq.Spec.Hard, nonSchedulableResourcesLimitations) {
		return nil, Forget
	}
	rq.Spec.Hard = nonSchedulableResourcesLimitations
	_, err = ctrl.aaqCli.CoreV1().ResourceQuotas(arqNS).Update(context.Background(), rq, metav1.UpdateOptions{})
	if err != nil {
		return err, Immediate
	}

	return nil, Forget
}

func filterNonScheduableResources(resourceList v1.ResourceList) v1.ResourceList {
	scheduableResources := getSchedulableResources()
	for _, resourceName := range scheduableResources {
		delete(resourceList, resourceName)
	}
	return resourceList
}

// ResourceListEqual checks if two ResourceList maps are equal.
func ResourceListEqual(list1, list2 v1.ResourceList) bool {
	// Check if the lengths of the two maps are different
	if len(list1) != len(list2) {
		return false
	}

	// Iterate over the keys and values in the first map
	for key, value1 := range list1 {
		// Check if the key exists in the second map
		value2, exists := list2[key]
		if !exists {
			return false
		}

		// Check if the values are equal
		if !value1.Equal(value2) {
			return false
		}
	}

	return true
}

// getSchedulableResources returns a list of resource names that are not counted in the resource quota.
func getSchedulableResources() []v1.ResourceName {
	// Add the resource names that should not be counted in the quota here
	return []v1.ResourceName{
		v1.ResourcePods,
		v1.ResourceCPU,
		v1.ResourceMemory,
		v1.ResourceEphemeralStorage,
		v1.ResourceRequestsCPU,
		v1.ResourceRequestsMemory,
		v1.ResourceRequestsEphemeralStorage,
		v1.ResourceLimitsCPU,
		v1.ResourceLimitsMemory,
		v1.ResourceLimitsEphemeralStorage,
	}
}

func (ctrl *RQController) Run(threadiness int) error {
	/*
		defer utilruntime.HandleCrash()

		_, err := ctrl.rqInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			UpdateFunc: ctrl.updateRQ,
			DeleteFunc: ctrl.deleteRQ,
		})
		if err != nil {
			return err
		}
		_, err = ctrl.arqInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			DeleteFunc: ctrl.deleteArq,
			UpdateFunc: ctrl.updateArq,
			AddFunc:    ctrl.addArq,
		})
		if err != nil {
			return err
		}
		klog.Info("Starting Arq controller")
		defer klog.Info("Shutting down Arq controller")

		for i := 0; i < threadiness; i++ {
			go wait.Until(ctrl.runWorker, time.Second, ctrl.stop)
		}

		<-ctrl.stop
	
	*/
	return nil

}