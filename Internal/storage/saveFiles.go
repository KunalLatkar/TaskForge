package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"project/Internal/task"
)

type FileStorage struct {
	baseDir string // Root directory
}

var statuses = []task.Status{
	task.CompletedStatus,
	task.PendingStatus,
	task.FailedStatus,
	task.ProcessingStatus,
}

func NewFileStorage(baseDir string) (*FileStorage, error) {
	fs := &FileStorage{baseDir: baseDir}

	if err := fs.initDirectories(); err != nil {
		return nil, err
	}
	return fs, nil
}

func (fs *FileStorage) initDirectories() error {

	for _, status := range statuses {
		dir := filepath.Join(fs.baseDir, string(status))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("Failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (fs *FileStorage) Save(t *task.Task) error {
	dir := filepath.Join(fs.baseDir, string(t.Status))
	filename := filepath.Join(dir, fmt.Sprintf("%s.json", t.ID))

	data, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("Failed to write the file: %w", err)
	}

	fmt.Printf("Saved the task %s to %s\n", t.ID, filename)
	return nil
}

func (fs *FileStorage) Load(id string) (*task.Task, error) {

	for _, status := range statuses {
		filename := filepath.Join(fs.baseDir, string(status), fmt.Sprintf("%s.json", id))

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("Failed to read the file %w", err)
		}

		var t task.Task
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, fmt.Errorf("Failed to unmarshal the task %w", err)
		}

		return &t, nil
	}

	return nil, fmt.Errorf("Task not found: %s", id)
}

func (fs *FileStorage) LoadAll() ([]*task.Task, error) {
	var tasks []*task.Task

	for _, status := range statuses {
		dir := filepath.Join(fs.baseDir, string(status))

		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to read the directory %s: %w", dir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filename := filepath.Join(dir, entry.Name())
			data, err := os.ReadFile(filename)
			if err != nil {
				fmt.Printf("Failed to read %s: %v\n", filename, err)
			}

			var t task.Task

			if err := json.Unmarshal(data, &t); err != nil {
				fmt.Printf("Failed to unmarshal the %s: %v\n", filename, err)
				continue
			}

			tasks = append(tasks, &t)
		}
	}

	fmt.Printf("Loaded %d tasks from disk\n", len(tasks))
	return tasks, nil
}

func (fs *FileStorage) Delete(id string) error {
	for _, status := range statuses {
		filename := filepath.Join(fs.baseDir, string(status), fmt.Sprintf("%s.json", id))

		if err := os.Remove(filename); err == nil {
			fmt.Printf("Deleted the task %s\n", filename)
			return nil
		}
	}

	return fmt.Errorf("task file not found %s\n", id)
}

func (fs *FileStorage) Move(t *task.Task, newStatus task.Status) error {

	for _, oldStatus := range statuses {

		if oldStatus == newStatus {
			continue
		}

		oldPath := filepath.Join(fs.baseDir, string(oldStatus), fmt.Sprintf("%s.json", t.ID))

		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to move the old file %w", err)
		}

	}

	t.Status = newStatus

	if err := fs.Save(t); err != nil {
		return fmt.Errorf("Failed to save the task at new location %w", err)
	}

	return nil
}
