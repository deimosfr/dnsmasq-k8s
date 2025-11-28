package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// UpdateConfigMapWithRetry updates a ConfigMap with a retry mechanism for handling conflicts.
// It attempts to update the ConfigMap up to 10 times with a random backoff between 500ms and 5s.
// On each retry, it re-fetches the ConfigMap to ensure the latest version is used.
func UpdateConfigMapWithRetry(ctx context.Context, clientset kubernetes.Interface, namespace, name, key, content string) error {
	maxRetries := 10
	minSleep := 500 * time.Millisecond
	maxSleep := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		// Step 1: Get ConfigMap
		configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				// Create ConfigMap if it doesn't exist
				newCM := &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Data: map[string]string{
						key: content,
					},
				}
				_, err = clientset.CoreV1().ConfigMaps(namespace).Create(ctx, newCM, metav1.CreateOptions{})
				if err != nil {
					if errors.IsAlreadyExists(err) {
						// If created concurrently, retry to get it and update
						continue
					}
					// For other errors, we also retry as per requirement
					// But we should probably check if it's a retryable error?
					// For now, we follow the instruction to retry.
					// Sleep before retry
					sleepDuration := minSleep + time.Duration(rand.Int63n(int64(maxSleep-minSleep)))
					select {
					case <-time.After(sleepDuration):
						continue
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				return nil
			}
			// Failed to get ConfigMap, retry
			sleepDuration := minSleep + time.Duration(rand.Int63n(int64(maxSleep-minSleep)))
			select {
			case <-time.After(sleepDuration):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Step 2: Apply Change
		if configMap.Data == nil {
			configMap.Data = make(map[string]string)
		}
		configMap.Data[key] = content

		// Step 3: Update
		_, err = clientset.CoreV1().ConfigMaps(namespace).Update(ctx, configMap, metav1.UpdateOptions{})
		if err == nil {
			return nil
		}

		if errors.IsConflict(err) {
			// Conflict, retry with backoff
			sleepDuration := minSleep + time.Duration(rand.Int63n(int64(maxSleep-minSleep)))
			select {
			case <-time.After(sleepDuration):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		} else {
			// Other error, return immediately? Or retry?
			// The requirement says "10 retries... As the configmap changes in between... error code is something like code 429 conflict issue".
			// It implies we mostly care about conflict. But let's stick to retrying only on conflict for update,
			// or maybe other transient errors.
			// However, the plan said "If other error, return error".
			return err
		}
	}

	return fmt.Errorf("failed to update ConfigMap %s after %d retries", name, maxRetries)
}
