import './style.css';
import './app.css';

import {
    CreateTask,
    GetTasks,
    UpdateTask,
    DeleteTask,
    ToggleTaskStatus,
    GetTasksByStatus,
    GetTasksByPriority,
    GetOverdueTasks
} from '../wailsjs/go/main/App';

// Global state
let tasks = [];
let currentFilter = {
    status: null,
    priority: null,
    sortBy: 'created_at',
    sortOrder: 'desc'
};
let editingTask = null;
let isDarkTheme = false;

// Priority and Status constants
const PRIORITY = {
    LOW: 0,
    MEDIUM: 1,
    HIGH: 2
};

const STATUS = {
    ACTIVE: 0,
    COMPLETED: 1
};

// Initialize app
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
});

function initializeApp() {
    renderApp();
    loadTasks();
    setupEventListeners();
    loadTheme();
}

function renderApp() {
    document.querySelector('#app').innerHTML = `
        <div class="header">
            <h1>📝 Todo App</h1>
            <button class="theme-toggle" onclick="toggleTheme()">
                <span id="theme-icon">🌙</span>
            </button>
        </div>
        
        <div class="container">
            <!-- Task Form -->
            <div class="task-form">
                <h2 style="margin-bottom: 1rem; color: var(--text-primary);">Add New Task</h2>
                <form id="task-form">
                    <div class="form-group">
                        <label for="task-title">Title *</label>
                        <input type="text" id="task-title" placeholder="Enter task title..." required>
                    </div>
                    
                    <div class="form-group">
                        <label for="task-description">Description</label>
                        <textarea id="task-description" placeholder="Enter task description..."></textarea>
                    </div>
                    
                    <div class="form-row-2">
                        <div class="form-group">
                            <label for="task-priority">Priority</label>
                            <select id="task-priority">
                                <option value="${PRIORITY.LOW}">Low</option>
                                <option value="${PRIORITY.MEDIUM}" selected>Medium</option>
                                <option value="${PRIORITY.HIGH}">High</option>
                            </select>
                        </div>
                        
                        <div class="form-group">
                            <label for="task-due-date">Due Date</label>
                            <input type="datetime-local" id="task-due-date">
                        </div>
                    </div>
                    
                    <div class="form-row">
                        <div></div>
                        <button type="submit" class="btn" id="submit-btn">Add Task</button>
                    </div>
                </form>
            </div>
            
            <!-- Filters -->
            <div class="filters">
                <div class="filter-group">
                    <label for="status-filter">Status:</label>
                    <select id="status-filter">
                        <option value="">All Tasks</option>
                        <option value="${STATUS.ACTIVE}">Active</option>
                        <option value="${STATUS.COMPLETED}">Completed</option>
                    </select>
                </div>
                
                <div class="filter-group">
                    <label for="priority-filter">Priority:</label>
                    <select id="priority-filter">
                        <option value="">All Priorities</option>
                        <option value="${PRIORITY.LOW}">Low</option>
                        <option value="${PRIORITY.MEDIUM}">Medium</option>
                        <option value="${PRIORITY.HIGH}">High</option>
                    </select>
                </div>
                
                <div class="filter-group">
                    <label for="sort-filter">Sort by:</label>
                    <select id="sort-filter">
                        <option value="created_at">Date Created</option>
                        <option value="due_date">Due Date</option>
                        <option value="priority">Priority</option>
                        <option value="title">Title</option>
                    </select>
                </div>
                
                <div class="filter-group">
                    <label for="order-filter">Order:</label>
                    <select id="order-filter">
                        <option value="desc">Descending</option>
                        <option value="asc">Ascending</option>
                    </select>
                </div>
                
                <button class="btn btn-secondary" onclick="clearFilters()">Clear Filters</button>
            </div>
            
            <!-- Task List -->
            <div id="task-list" class="task-list">
                <!-- Tasks will be rendered here -->
            </div>
        </div>
        
        <!-- Modal for delete confirmation -->
        <div id="delete-modal" class="modal-overlay" style="display: none;">
            <div class="modal">
                <div class="modal-header">
                    <h3 class="modal-title">Delete Task</h3>
                    <button class="modal-close" onclick="closeModal()">&times;</button>
                </div>
                <p>Are you sure you want to delete this task? This action cannot be undone.</p>
                <div class="modal-actions">
                    <button class="btn btn-secondary" onclick="closeModal()">Cancel</button>
                    <button class="btn btn-danger" onclick="confirmDelete()">Delete</button>
                </div>
            </div>
        </div>
    `;
}

function setupEventListeners() {
    // Task form submission
    document.getElementById('task-form').addEventListener('submit', handleTaskSubmit);
    
    // Filter changes
    document.getElementById('status-filter').addEventListener('change', applyFilters);
    document.getElementById('priority-filter').addEventListener('change', applyFilters);
    document.getElementById('sort-filter').addEventListener('change', applyFilters);
    document.getElementById('order-filter').addEventListener('change', applyFilters);
}

// Task management functions
async function loadTasks() {
    try {
        const filterJSON = JSON.stringify(currentFilter);
        const result = await GetTasks(filterJSON);
        tasks = JSON.parse(result);
        renderTasks();
    } catch (error) {
        console.error('Error loading tasks:', error);
        showNotification('Error loading tasks', 'error');
    }
}

async function handleTaskSubmit(e) {
    e.preventDefault();
    
    const title = document.getElementById('task-title').value.trim();
    const description = document.getElementById('task-description').value.trim();
    const priority = parseInt(document.getElementById('task-priority').value);
    const dueDate = document.getElementById('task-due-date').value;
    
    if (!title) {
        showNotification('Please enter a task title', 'error');
        return;
    }
    
    const taskData = {
        title,
        description,
        priority,
        dueDate: dueDate ? new Date(dueDate).toISOString() : null
    };
    
    try {
        if (editingTask) {
            // Update existing task
            const updateData = {
                id: editingTask.id,
                title,
                description,
                priority,
                status: editingTask.status,
                dueDate: dueDate ? new Date(dueDate).toISOString() : null
            };
            
            const result = await UpdateTask(JSON.stringify(updateData));
            const updatedTask = JSON.parse(result);
            
            const index = tasks.findIndex(t => t.id === updatedTask.id);
            if (index !== -1) {
                tasks[index] = updatedTask;
            }
            
            editingTask = null;
            document.getElementById('submit-btn').textContent = 'Add Task';
            showNotification('Task updated successfully', 'success');
        } else {
            // Create new task
            const result = await CreateTask(JSON.stringify(taskData));
            const newTask = JSON.parse(result);
            tasks.unshift(newTask);
            showNotification('Task created successfully', 'success');
        }
        
        // Reset form
        document.getElementById('task-form').reset();
        document.getElementById('task-due-date').value = '';
        renderTasks();
        
    } catch (error) {
        console.error('Error saving task:', error);
        showNotification('Error saving task', 'error');
    }
}

async function toggleTaskStatus(taskId) {
    try {
        const result = await ToggleTaskStatus(taskId);
        const updatedTask = JSON.parse(result);
        
        const index = tasks.findIndex(t => t.id === taskId);
        if (index !== -1) {
            tasks[index] = updatedTask;
        }
        
        renderTasks();
        showNotification('Task status updated', 'success');
    } catch (error) {
        console.error('Error updating task status:', error);
        showNotification('Error updating task status', 'error');
    }
}

async function deleteTask(taskId) {
    try {
        await DeleteTask(taskId);
        tasks = tasks.filter(t => t.id !== taskId);
        renderTasks();
        showNotification('Task deleted successfully', 'success');
    } catch (error) {
        console.error('Error deleting task:', error);
        showNotification('Error deleting task', 'error');
    }
}

function editTask(task) {
    editingTask = task;
    
    document.getElementById('task-title').value = task.title;
    document.getElementById('task-description').value = task.description || '';
    document.getElementById('task-priority').value = task.priority;
    document.getElementById('task-due-date').value = task.dueDate ? 
        new Date(task.dueDate).toISOString().slice(0, 16) : '';
    
    document.getElementById('submit-btn').textContent = 'Update Task';
    
    // Scroll to form
    document.querySelector('.task-form').scrollIntoView({ behavior: 'smooth' });
}

function renderTasks() {
    const taskList = document.getElementById('task-list');
    
    if (tasks.length === 0) {
        taskList.innerHTML = `
            <div class="empty-state">
                <h3>No tasks found</h3>
                <p>Create your first task to get started!</p>
            </div>
        `;
        return;
    }
    
    const filteredTasks = getFilteredTasks();
    
    taskList.innerHTML = filteredTasks.map(task => `
        <div class="task-item ${task.status === STATUS.COMPLETED ? 'completed' : ''}" data-task-id="${task.id}">
            <div class="task-header">
                <input type="checkbox" 
                       class="task-checkbox" 
                       ${task.status === STATUS.COMPLETED ? 'checked' : ''}
                       onchange="toggleTaskStatus('${task.id}')">
                
                <div class="task-content">
                    <div class="task-title">${escapeHtml(task.title)}</div>
                    ${task.description ? `<div class="task-description">${escapeHtml(task.description)}</div>` : ''}
                    
                    <div class="task-meta">
                        <span class="task-priority ${getPriorityClass(task.priority)}">
                            ${getPriorityIcon(task.priority)} ${getPriorityText(task.priority)}
                        </span>
                        
                        ${task.dueDate ? `
                            <span class="task-due-date ${isOverdue(task.dueDate) ? 'overdue' : ''}">
                                📅 ${formatDate(task.dueDate)}
                            </span>
                        ` : ''}
                        
                        <span>Created: ${formatDate(task.createdAt)}</span>
                    </div>
                </div>
                
                <div class="task-actions">
                    <button class="btn btn-sm btn-secondary" onclick="editTask(${JSON.stringify(task).replace(/"/g, '&quot;')})">
                        Edit
                    </button>
                    <button class="btn btn-sm btn-danger" onclick="showDeleteModal('${task.id}')">
                        Delete
                    </button>
                </div>
            </div>
        </div>
    `).join('');
}

function getFilteredTasks() {
    let filtered = [...tasks];
    
    if (currentFilter.status !== null) {
        filtered = filtered.filter(task => task.status === currentFilter.status);
    }
    
    if (currentFilter.priority !== null) {
        filtered = filtered.filter(task => task.priority === currentFilter.priority);
    }
    
    // Sort tasks
    filtered.sort((a, b) => {
        let aValue, bValue;
        
        switch (currentFilter.sortBy) {
            case 'title':
                aValue = a.title.toLowerCase();
                bValue = b.title.toLowerCase();
                break;
            case 'priority':
                aValue = a.priority;
                bValue = b.priority;
                break;
            case 'due_date':
                aValue = a.dueDate ? new Date(a.dueDate) : new Date(0);
                bValue = b.dueDate ? new Date(b.dueDate) : new Date(0);
                break;
            case 'created_at':
            default:
                aValue = new Date(a.createdAt);
                bValue = new Date(b.createdAt);
                break;
        }
        
        if (currentFilter.sortOrder === 'asc') {
            return aValue > bValue ? 1 : -1;
        } else {
            return aValue < bValue ? 1 : -1;
        }
    });
    
    return filtered;
}

function applyFilters() {
    const statusFilter = document.getElementById('status-filter').value;
    const priorityFilter = document.getElementById('priority-filter').value;
    const sortFilter = document.getElementById('sort-filter').value;
    const orderFilter = document.getElementById('order-filter').value;
    
    currentFilter = {
        status: statusFilter === '' ? null : parseInt(statusFilter),
        priority: priorityFilter === '' ? null : parseInt(priorityFilter),
        sortBy: sortFilter,
        sortOrder: orderFilter
    };
    
    renderTasks();
}

function clearFilters() {
    document.getElementById('status-filter').value = '';
    document.getElementById('priority-filter').value = '';
    document.getElementById('sort-filter').value = 'created_at';
    document.getElementById('order-filter').value = 'desc';
    
    currentFilter = {
        status: null,
        priority: null,
        sortBy: 'created_at',
        sortOrder: 'desc'
    };
    
    renderTasks();
}

// Modal functions
let taskToDelete = null;

function showDeleteModal(taskId) {
    taskToDelete = taskId;
    document.getElementById('delete-modal').style.display = 'flex';
}

function closeModal() {
    document.getElementById('delete-modal').style.display = 'none';
    taskToDelete = null;
}

function confirmDelete() {
    if (taskToDelete) {
        deleteTask(taskToDelete);
        closeModal();
    }
}

// Theme functions
function toggleTheme() {
    isDarkTheme = !isDarkTheme;
    document.documentElement.setAttribute('data-theme', isDarkTheme ? 'dark' : 'light');
    document.getElementById('theme-icon').textContent = isDarkTheme ? '☀️' : '🌙';
    localStorage.setItem('theme', isDarkTheme ? 'dark' : 'light');
}

function loadTheme() {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        isDarkTheme = true;
        document.documentElement.setAttribute('data-theme', 'dark');
        document.getElementById('theme-icon').textContent = '☀️';
    }
}

// Utility functions
function getPriorityText(priority) {
    switch (priority) {
        case PRIORITY.LOW: return 'Low';
        case PRIORITY.MEDIUM: return 'Medium';
        case PRIORITY.HIGH: return 'High';
        default: return 'Unknown';
    }
}

function getPriorityClass(priority) {
    switch (priority) {
        case PRIORITY.LOW: return 'low';
        case PRIORITY.MEDIUM: return 'medium';
        case PRIORITY.HIGH: return 'high';
        default: return '';
    }
}

function getPriorityIcon(priority) {
    switch (priority) {
        case PRIORITY.LOW: return '🟢';
        case PRIORITY.MEDIUM: return '🟡';
        case PRIORITY.HIGH: return '🔴';
        default: return '⚪';
    }
}

function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 1) {
        return 'Today';
    } else if (diffDays === 2) {
        return 'Yesterday';
    } else if (diffDays <= 7) {
        return `${diffDays - 1} days ago`;
    } else {
        return date.toLocaleDateString();
    }
}

function isOverdue(dateString) {
    const dueDate = new Date(dateString);
    const now = new Date();
    return dueDate < now;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showNotification(message, type = 'info') {
    // Simple notification - you could enhance this with a proper notification system
    console.log(`${type.toUpperCase()}: ${message}`);
}

// Global functions for onclick handlers
window.toggleTaskStatus = toggleTaskStatus;
window.editTask = editTask;
window.showDeleteModal = showDeleteModal;
window.closeModal = closeModal;
window.confirmDelete = confirmDelete;
window.toggleTheme = toggleTheme;
window.clearFilters = clearFilters;
