package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	ID          int
	Title       string
	Description string
}

var tasks = []Task{
	{ID: 1, Title: "Task 1", Description: "Description 1"},
	{ID: 2, Title: "Task 2", Description: "Description 2"},
	{ID: 3, Title: "Task 3", Description: "Description 3"},
}

func main() {
	http.HandleFunc("/", tasksHandler)
	http.HandleFunc("/create", createTaskHandler)
	http.HandleFunc("/delete/", deleteTaskHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	err := tmpl.Execute(w, map[string]interface{}{"Tasks": tasks})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		title := r.FormValue("title")
		description := r.FormValue("description")
		newTask := Task{
			ID:          len(tasks) + 1,
			Title:       title,
			Description: description,
		}
		tasks = append(tasks, newTask)

		// Render only the new task HTML
		taskHTML := fmt.Sprintf(
			`<li id="task-%d" class="bg-gray-200 p-2 mb-2 rounded shadow">
                %s - %s
                <div>
                    <button
                        hx-delete="/delete/%d"
                        hx-target="#task-%d"
                        hx-swap="outerHTML"
                        class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                    >
                        Delete
                    </button>
                </div>
            </li>`,
			newTask.ID, newTask.Title, newTask.Description, newTask.ID, newTask.ID,
		)
		w.Write([]byte(taskHTML))
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/delete/"):]
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	// Return an empty response to remove the task element from the DOM
	w.Write([]byte(""))
}
