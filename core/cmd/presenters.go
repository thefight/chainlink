package cmd

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// JAID represents a JSON API ID.
// It implements the api2go MarshalIdentifier and UnmarshalIdentitier interface.
type JAID struct {
	ID string `json:"-"`
}

// GetID implements the api2go MarshalIdentifier interface.
func (jaid JAID) GetID() string {
	return jaid.ID
}

// SetID implements the api2go UnmarshalIdentitier interface.
func (jaid *JAID) SetID(value string) error {
	jaid.ID = value

	return nil
}

// Job represents a V2 Job
type Job struct {
	JAID
	Name         string `json:"name"`
	Type         string `json:"type"`
	PipelineSpec struct {
		ID           int32  `json:"ID"`
		DotDAGSource string `json:"dotDagSource"`
	} `json:"pipelineSpec"`
}

// GetName implements the api2go EntityNamer interface
func (j Job) GetName() string {
	return "specDBs"
}

// GetTaskTypes extracts the tasks types from the dependency graph
func (j Job) GetTaskTypes() ([]string, error) {
	types := []string{}
	dag := pipeline.NewTaskDAG()
	dag.UnmarshalText([]byte(j.PipelineSpec.DotDAGSource))

	tasks, err := dag.TasksInDependencyOrder()
	if err != nil {
		return nil, err
	}

	// Gather in reverse order as dependency tasks start from output
	for i := len(tasks) - 1; i >= 0; i-- {
		t := tasks[i]
		types = append(types, fmt.Sprintf("%s %s", t.DotID(), t.Type()))
	}

	return types, nil
}

// FriendlyTaskTypes returns the tasks tasks as a string separated by newlines.
func (j Job) FriendlyTaskTypes() []string {
	taskTypes, err := j.GetTaskTypes()
	if err != nil {
		return []string{"unknown"}
	}

	return taskTypes
}

// ToRow returns the job as a table row
func (j Job) ToRow() [][]string {
	row := [][]string{}

	for _, t := range j.FriendlyTaskTypes() {
		row = append(row, []string{
			j.ID,
			j.Name,
			j.Type,
			t,
		})
	}
	return row
}
