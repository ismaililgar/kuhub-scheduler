package job

import (
	"fmt"
)

type Registry struct {
	jobs map[string]Job
}

func NewRegistry() *Registry {
	return &Registry{
		jobs: make(map[string]Job),
	}
}

func (r *Registry) Register(j Job) {
	r.jobs[j.Name()] = j
}

func (r *Registry) Get(name string) (Job, error) {
	j, ok := r.jobs[name]
	if !ok {
		return nil, fmt.Errorf("job bulunamadı: %s", name)
	}
	return j, nil
}

func (r *Registry) All() []Job {
	all := make([]Job, 0, len(r.jobs))
	for _, j := range r.jobs {
		all = append(all, j)
	}
	return all
}
