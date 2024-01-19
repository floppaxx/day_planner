// Replace this base URL with the actual URL where your Go backend is hosted
const backendBaseUrl = 'http://localhost:8001';
 
function date() {
    const currentDate = new Date();
    GlobalcurrentDate = currentDate.toISOString().split('T')[0];
    writeDatetoPage(currentDate);
    updateTaskList();
    
}

function writeDatetoPage(date) {
    var Weekdays = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
    var Months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
    var nextDayOfWeek = Weekdays[date.getDay()];
    var nextMonth = Months[date.getMonth()];
    document.getElementById("date").innerHTML = nextDayOfWeek + ", " + date.getDate() + " " + nextMonth + " " + date.getFullYear();
}

var dayOffset = 0;
var GlobalcurrentDate = "";

function nextDay() {
    
    dayOffset++;
    const currentDate = new Date();
    const nextDate = new Date(currentDate);
    nextDate.setDate(currentDate.getDate() + dayOffset);
    GlobalcurrentDate = nextDate.toISOString().split('T')[0];
    writeDatetoPage(nextDate);
    updateTaskList();
}

function previousDay() {
    dayOffset--;
    const currentDate = new Date();
    const prevDate = new Date(currentDate);
    prevDate.setDate(currentDate.getDate() + dayOffset);
    GlobalcurrentDate = prevDate.toISOString().split('T')[0];
    writeDatetoPage(prevDate);
    updateTaskList();
}

var taskList = [];

function addTask() {
    const taskNameInput = document.getElementById('taskName');
    const taskName = taskNameInput.value.trim();


    if (taskName !== '') {
        const currentDate = new Date(GlobalcurrentDate);

        // Create a new task object
        const task = {
            name: taskName,
            date: currentDate.toISOString().split('T')[0],
            done: false
        };

        // Add the task to the backend
        fetch(`${backendBaseUrl}/tasks`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(task),
        })
        .then(response => response.json())
        .then(data => {
            // Update the task list with the returned data
            taskList.push(data);
            updateTaskList(true);
        })
        .catch(error => {
            console.error('Error adding task:', error);
        });

        // Clear the input field
        taskNameInput.value = '';
    } else {
        alert('Task name cannot be empty!');
    }
}

function updateTaskList(state) {
    const currentDate = new Date(GlobalcurrentDate);
    
    if (sessionStorage.getItem('tasks') === null || state == true) {
        sessionStorage.clear();
        fetch(`${backendBaseUrl}/tasks`)
        .then(response => response.json())
        .then(data => {
            // Filter tasks for the current date
            sessionStorage.setItem('tasks', JSON.stringify(data));

            taskList = taskList.filter(task => {
                const taskDate = new Date(task.date);
                const taskDateWithoutTime = new Date(taskDate.getFullYear(), taskDate.getMonth(), taskDate.getDate());
                const currentDateWithoutTime = new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate());
                
                return taskDateWithoutTime.getTime() === currentDateWithoutTime.getTime();
            });
            // Display tasks
            const taskListElement = document.getElementById('taskList');
            taskListElement.innerHTML = '';


            taskList.forEach(task => {
                const isCompleted = task.done;

                const taskElement = document.createElement('div');
                taskElement.className = `task ${isCompleted ? 'completed' : ''}`;
                taskElement.innerHTML = `
                    <div class="taskName">${task.name}</div>
                    <div class="taskActions">
                        <img src="img/check-icon.svg" onclick="completeTask(${task.id})">
                        <img src="img/delete-icon.svg" onclick="deleteTask(${task.id})">
                    </div>
                `;
                taskListElement.appendChild(taskElement);
            });
        })
        .catch(error => {
            console.error('Error fetching tasks:', error);
        });
    }
    else
    {
        taskList = JSON.parse(sessionStorage.getItem('tasks'));
        // Display tasks
        
        taskList = taskList.filter(task => {
            const taskDate = new Date(task.date);
            const taskDateWithoutTime = new Date(taskDate.getFullYear(), taskDate.getMonth(), taskDate.getDate());
            const currentDateWithoutTime = new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate());
            
            return taskDateWithoutTime.getTime() === currentDateWithoutTime.getTime();
        });
        
        const taskListElement = document.getElementById('taskList');
        taskListElement.innerHTML = '';

        taskList.forEach(task => {
            const isCompleted = task.done;
            ('HUI');

            const taskElement = document.createElement('div');
            taskElement.className = `task ${isCompleted ? 'completed' : ''}`;
            taskElement.innerHTML = `
                <div class="taskName">${task.name}</div>
                <div class="taskActions">
                    <img src="img/check-icon.svg" onclick="completeTask(${task.id})">
                    <img src="img/delete-icon.svg" onclick="deleteTask(${task.id})">
                </div>
            `;
            (taskElement);
            taskListElement.appendChild(taskElement);
        });
}
}

function completeTask(taskId) {
    const storedTasks = sessionStorage.getItem('tasks');
    const tasks = storedTasks ? JSON.parse(storedTasks) : [];

    const taskIndex = tasks.findIndex(task => task.id === taskId);
    if (taskIndex !== -1) {
        // Toggle the 'done' status
        tasks[taskIndex].done = !tasks[taskIndex].done;

        // Update the task in session storage
        sessionStorage.setItem('tasks', JSON.stringify(tasks));

        // Update the task list and re-render UI
        taskList = tasks;
        updateTaskList();
    }

    // Update the task on the backend
    fetch(`${backendBaseUrl}/tasks/${taskId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ done: tasks[taskIndex].done }),
    })
    .then(response => response.json())
    .then(data => {
        // Update the task list with the returned data
        taskList = taskList.map(t => (t.id === taskId ? data : t));
        updateTaskList();
    })
    .catch(error => {
        console.error('Error updating task:', error);
    });
}

function deleteTask(taskId) {
    // Delete the task on the backend
    fetch(`${backendBaseUrl}/tasks/${taskId}`, {
        method: 'DELETE',
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Error deleting task');
        }
        // Remove the task from session storage
        const storedTasks = sessionStorage.getItem('tasks');
        const tasks = storedTasks ? JSON.parse(storedTasks) : [];
        const updatedTasks = tasks.filter(task => task.id !== taskId);
        sessionStorage.setItem('tasks', JSON.stringify(updatedTasks));

        // Remove the task from the task list
        taskList = updatedTasks;
        updateTaskList();
    })
    .catch(error => {
        console.error('Error deleting task:', error);
    });
}



date();
