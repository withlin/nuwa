/*
Copyright 2019 yametech.

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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var sidecarlog = logf.Log.WithName("sidecar-resource")

func (r *Sidecar) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-nuwa-nip-io-v1-sidecar,mutating=true,failurePolicy=fail,groups=nuwa.nip.io,resources=sidecars,verbs=create;update,versions=v1,name=msidecar.kb.io

var _ webhook.Defaulter = &Sidecar{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Sidecar) Default() {
	sidecarlog.Info("sidecar default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.

}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-nuwa-nip-io-v1-sidecar,mutating=false,failurePolicy=fail,groups=nuwa.nip.io,resources=sidecars,versions=v1,name=vsidecar.kb.io

var _ webhook.Validator = &Sidecar{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Sidecar) ValidateCreate() error {
	sidecarlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Sidecar) ValidateUpdate(old runtime.Object) error {
	sidecarlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Sidecar) ValidateDelete() error {
	sidecarlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
