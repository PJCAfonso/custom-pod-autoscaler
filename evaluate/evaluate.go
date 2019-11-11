/*
Copyright 2019 The Custom Pod Autoscaler Authors.

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

package evaluate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jthomperoo/custom-pod-autoscaler/config"
	"github.com/jthomperoo/custom-pod-autoscaler/models"
)

const invalidEvaluationMessage = "Invalid evaluation returned by evaluator: %s"

type executeWithPiper interface {
	ExecuteWithPipe(command string, value string, timeout int) (*bytes.Buffer, error)
}

// Evaluator handles triggering the evaluation logic to decide how to scale a resource
type Evaluator struct {
	Config   *config.Config
	Executer executeWithPiper
}

// ErrInvalidEvaluation occurs when the the evaluator reports success but returns invalid JSON
type ErrInvalidEvaluation struct {
	Message string
}

func (e *ErrInvalidEvaluation) Error() string {
	return e.Message
}

// GetEvaluation uses the metrics provided to determine a set of evaluations
func (e *Evaluator) GetEvaluation(resourceMetrics *models.ResourceMetrics) (*models.Evaluation, error) {
	// Convert metrics into JSON
	metricJSON, err := json.Marshal(resourceMetrics.Metrics)
	if err != nil {
		// Should not occur, panic
		log.Panic(err)
	}

	// Execute the Evaluate command with the metric JSON
	outb, err := e.Executer.ExecuteWithPipe(e.Config.Evaluate, string(metricJSON), e.Config.EvaluateTimeout)
	if err != nil {
		log.Println(outb.String())
		return nil, err
	}
	evaluation := &models.Evaluation{}
	err = json.Unmarshal(outb.Bytes(), evaluation)
	if err != nil {
		return nil, err
	}
	if evaluation.TargetReplicas == nil {
		return nil, &ErrInvalidEvaluation{
			Message: fmt.Sprintf(invalidEvaluationMessage, outb.String()),
		}
	}
	return evaluation, nil
}