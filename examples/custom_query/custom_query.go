package main

import (
	"log/slog"

	"github.com/m87/nod"
	sqlite_nod "github.com/m87/nod/sqlite"
)

type Task struct {
	Title   string
	Project string
	Done    bool
}

func (task *Task) ToNode() (*nod.Node, error) {
	status := "open"
	if task.Done {
		status = "done"
	}

	return &nod.Node{
		Core: nod.NodeCore{
			Name:   task.Title,
			Kind:   "task",
			Status: status,
		},
		KV: map[string]*nod.NodeKV{
			"project": {Key: "project", ValueText: nod.Ptr(task.Project)},
		},
	}, nil
}

func (task *Task) FromNode(node *nod.Node) error {
	task.Title = node.Core.Name
	task.Done = node.Core.Status == "done"
	if project := node.KV["project"]; project != nil && project.ValueText != nil {
		task.Project = *project.ValueText
	}
	return nil
}

func (task *Task) IsApplicable(node *nod.Node) bool {
	return node != nil && node.Core.Kind == "task"
}

// TaskConditions maps domain language to nod's storage-level expressions.
// Query callers only use TaskWhere and do not need to know where fields are stored.
type TaskConditions struct{}

var TaskWhere TaskConditions

func (TaskConditions) TitleEquals(title string) nod.Expression {
	return nod.NodeFields.Name.Equals(title)
}

func (TaskConditions) InProject(project string) nod.Expression {
	return nod.KvString("project").Equals(project)
}

func (TaskConditions) IsOpen() nod.Expression {
	return nod.NodeFields.Status.Equals("open")
}

func (conditions TaskConditions) OpenInProject(project string) nod.Expression {
	return nod.And(conditions.IsOpen(), conditions.InProject(project))
}

func main() {
	repo, err := sqlite_nod.NewRepositoryInMemory(slog.Default(), nil)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	tasks := nod.Nodes[Task](repo)
	for _, task := range []*Task{
		{Title: "Add codec example", Project: "nod"},
		{Title: "Release previous version", Project: "nod", Done: true},
		{Title: "Unrelated task", Project: "other"},
	} {
		if _, err := tasks.SaveNode(task); err != nil {
			panic(err)
		}
	}

	// The query uses domain-specific conditions instead of NodeFields or KV keys.
	found, err := nod.NewTypedNodeQuery[Task](repo).
		WithKV().
		Where(TaskWhere.OpenInProject("nod")).
		FindAll()
	if err != nil {
		panic(err)
	}

	for _, task := range found {
		slog.Info("Found task", "title", task.Title, "project", task.Project)
	}
}
