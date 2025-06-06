//
// Copyright 2021, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"net/http"
	"time"
)

type (
	SidekiqServiceInterface interface {
		GetQueueMetrics(options ...RequestOptionFunc) (*QueueMetrics, *Response, error)
		GetProcessMetrics(options ...RequestOptionFunc) (*ProcessMetrics, *Response, error)
		GetJobStats(options ...RequestOptionFunc) (*JobStats, *Response, error)
		GetCompoundMetrics(options ...RequestOptionFunc) (*CompoundMetrics, *Response, error)
	}

	// SidekiqService handles communication with the sidekiq service
	//
	// GitLab API docs: https://docs.gitlab.com/api/sidekiq_metrics/
	SidekiqService struct {
		client *Client
	}
)

var _ SidekiqServiceInterface = (*SidekiqService)(nil)

// QueueMetrics represents the GitLab sidekiq queue metrics.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-queue-metrics
type QueueMetrics struct {
	Queues map[string]struct {
		Backlog int `json:"backlog"`
		Latency int `json:"latency"`
	} `json:"queues"`
}

// GetQueueMetrics lists information about all the registered queues,
// their backlog and their latency.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-queue-metrics
func (s *SidekiqService) GetQueueMetrics(options ...RequestOptionFunc) (*QueueMetrics, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/sidekiq/queue_metrics", nil, options)
	if err != nil {
		return nil, nil, err
	}

	q := new(QueueMetrics)
	resp, err := s.client.Do(req, q)
	if err != nil {
		return nil, resp, err
	}

	return q, resp, nil
}

// ProcessMetrics represents the GitLab sidekiq process metrics.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-process-metrics
type ProcessMetrics struct {
	Processes []struct {
		Hostname    string     `json:"hostname"`
		Pid         int        `json:"pid"`
		Tag         string     `json:"tag"`
		StartedAt   *time.Time `json:"started_at"`
		Queues      []string   `json:"queues"`
		Labels      []string   `json:"labels"`
		Concurrency int        `json:"concurrency"`
		Busy        int        `json:"busy"`
	} `json:"processes"`
}

// GetProcessMetrics lists information about all the Sidekiq workers registered
// to process your queues.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-process-metrics
func (s *SidekiqService) GetProcessMetrics(options ...RequestOptionFunc) (*ProcessMetrics, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/sidekiq/process_metrics", nil, options)
	if err != nil {
		return nil, nil, err
	}

	p := new(ProcessMetrics)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// JobStats represents the GitLab sidekiq job stats.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-job-statistics
type JobStats struct {
	Jobs struct {
		Processed int `json:"processed"`
		Failed    int `json:"failed"`
		Enqueued  int `json:"enqueued"`
	} `json:"jobs"`
}

// GetJobStats list information about the jobs that Sidekiq has performed.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-job-statistics
func (s *SidekiqService) GetJobStats(options ...RequestOptionFunc) (*JobStats, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/sidekiq/job_stats", nil, options)
	if err != nil {
		return nil, nil, err
	}

	j := new(JobStats)
	resp, err := s.client.Do(req, j)
	if err != nil {
		return nil, resp, err
	}

	return j, resp, nil
}

// CompoundMetrics represents the GitLab sidekiq compounded stats.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-a-compound-response-of-all-the-previously-mentioned-metrics
type CompoundMetrics struct {
	QueueMetrics
	ProcessMetrics
	JobStats
}

// GetCompoundMetrics lists all the currently available information about Sidekiq.
// Get a compound response of all the previously mentioned metrics
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-a-compound-response-of-all-the-previously-mentioned-metrics
func (s *SidekiqService) GetCompoundMetrics(options ...RequestOptionFunc) (*CompoundMetrics, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/sidekiq/compound_metrics", nil, options)
	if err != nil {
		return nil, nil, err
	}

	c := new(CompoundMetrics)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}
