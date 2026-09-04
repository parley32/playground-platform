package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type mockTaskService struct {
	task  Task
	tasks []Task
	err   error
	store Store
}

func (m mockTaskService) ListSingleTask(ID string) (Task, error) { return m.task, m.err }

func (m mockTaskService) List() ([]Task, error) { return m.tasks, m.err }

func (m mockTaskService) Add(string) (Task, error) { return m.task, m.err }

func (m mockTaskService) Delete(string) error { return m.err }

func (m mockTaskService) Check(string) error { return m.err }

func TestListSingleItem(t *testing.T) {
	tests := []struct {
		name       string
		wantTask   Task
		err        error
		returnCode int
	}{
		{
			name:       "single item success",
			wantTask:   Task{ID: 1, Text: "Learn Go", Done: false},
			err:        nil,
			returnCode: http.StatusOK,
		},
		{
			name:       "single item invalid ID",
			wantTask:   Task{},
			err:        ErrInvalidID,
			returnCode: http.StatusBadRequest,
		}, {
			name:       "single item not found ",
			wantTask:   Task{},
			err:        ErrNotFound,
			returnCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := mockTaskService{
				task: test.wantTask,
				err:  test.err,
			}

			handler := NewTaskHandler(service)

			rec := httptest.NewRecorder()

			handler.ListSingleItem(rec, "1")

			if rec.Code != test.returnCode {
				t.Fatalf("Return code invalid")
			}

			if test.err != nil {

				if !strings.Contains(rec.Body.String(), test.err.Error()) {
					t.Fatalf("body = %q, want it to contain %q", rec.Body.String(), test.err.Error())
				}
			}

			if test.err == nil {
				contentType := rec.Header().Get("Content-Type")
				if !strings.HasPrefix(contentType, "application/json") {
					t.Fatalf("Content-Type = %q, want application/json", contentType)
				}
				var got Task
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("Stored task not a valid JSON")
				}

				if !reflect.DeepEqual(got, test.wantTask) {
					t.Fatalf("got %+v, expected %+v", got, test.wantTask)
				}
			}
		})
	}
}

func TestListItems(t *testing.T) {
	tests := []struct {
		name       string
		wantTask   []Task
		err        error
		returnCode int
	}{
		{
			name:       "list items success",
			wantTask:   []Task{{ID: 1, Text: "Learn Go", Done: false}},
			err:        nil,
			returnCode: http.StatusOK,
		},
		{
			name:       "list items failure",
			wantTask:   []Task{},
			err:        ErrInternal,
			returnCode: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := mockTaskService{
				tasks: test.wantTask,
				err:   test.err,
			}

			handler := NewTaskHandler(service)

			rec := httptest.NewRecorder()

			handler.List(rec)

			if rec.Code != test.returnCode {
				t.Fatalf("got %q, wanted %q", rec.Code, test.returnCode)
			}

			if test.err != nil {
				if !strings.Contains(rec.Body.String(), test.err.Error()) {
					t.Fatalf("body = %q, want it to contain %q", rec.Body.String(), test.err.Error())
				}
			}

			if test.err == nil {
				var got []Task
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("response is not valid JSON: %v", err)
				}

				if !reflect.DeepEqual(got, test.wantTask) {
					t.Fatalf("got %+v, expected %+v", got, test.wantTask)
				}
			}
		})
	}
}

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name       string
		wantTask   Task
		err        error
		returnCode int
	}{
		{
			name:       "create item success",
			wantTask:   Task{ID: 1, Text: "Learn Go", Done: false},
			err:        nil,
			returnCode: http.StatusCreated,
		},
		{
			name:       "invalid item name, failure",
			wantTask:   Task{},
			err:        ErrInvalidTaskName,
			returnCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := mockTaskService{
				task: test.wantTask,
				err:  test.err,
			}

			handler := NewTaskHandler(service)
			req := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				strings.NewReader(`{"description":"Learn Go"}`),
			)
			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			if rec.Code != test.returnCode {
				t.Fatalf("got %d, wanted %d", rec.Code, test.returnCode)
			}

			if test.err != nil {
				if !strings.Contains(rec.Body.String(), test.err.Error()) {
					t.Fatalf("body = %q, want it to contain %q", rec.Body.String(), test.err.Error())
				}
			}

			if test.err == nil {
				var got Task
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("response is not valid JSON: %v", err)
				}

				if !reflect.DeepEqual(got, test.wantTask) {
					t.Fatalf("got %+v, expected %+v", got, test.wantTask)
				}
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name       string
		want       string
		err        error
		returnCode int
	}{
		{
			name:       "delete item success",
			want:       "task deleted successfully",
			err:        nil,
			returnCode: http.StatusOK,
		},
		{
			name:       "invalid item id, failure",
			err:        ErrInvalidID,
			returnCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := mockTaskService{
				err: test.err,
			}

			handler := NewTaskHandler(service)
			rec := httptest.NewRecorder()

			handler.Delete(rec, "1")

			if rec.Code != test.returnCode {
				t.Fatalf("got %d, wanted %d", rec.Code, test.returnCode)
			}

			if test.err != nil {
				if !strings.Contains(rec.Body.String(), test.err.Error()) {
					t.Fatalf("body = %q, want it to contain %q", rec.Body.String(), test.err.Error())
				}
			}

			if test.err == nil {
				var got string
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("response is not valid JSON: %v", err)
				}

				if !reflect.DeepEqual(got, test.want) {
					t.Fatalf("got %+v, expected %+v", got, test.want)
				}
			}
		})
	}
}

func TestCheckItem(t *testing.T) {
	tests := []struct {
		name       string
		want       string
		err        error
		returnCode int
	}{
		{
			name:       "check item success",
			want:       "task completed successfully",
			err:        nil,
			returnCode: http.StatusOK,
		},
		{
			name:       "invalid item id, failure",
			err:        ErrInvalidID,
			returnCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := mockTaskService{
				err: test.err,
			}

			handler := NewTaskHandler(service)
			rec := httptest.NewRecorder()

			handler.Check(rec, "1")

			if rec.Code != test.returnCode {
				t.Fatalf("got %d, wanted %d", rec.Code, test.returnCode)
			}

			if test.err != nil {
				if !strings.Contains(rec.Body.String(), test.err.Error()) {
					t.Fatalf("body = %q, want it to contain %q", rec.Body.String(), test.err.Error())
				}
			}

			if test.err == nil {
				var got string
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("response is not valid JSON: %v", err)
				}

				if !reflect.DeepEqual(got, test.want) {
					t.Fatalf("got %+v, expected %+v", got, test.want)
				}
			}
		})
	}
}
