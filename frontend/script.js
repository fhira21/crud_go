/**
 * ==============================================================================
 * Employee Management - Frontend JavaScript (Vanilla JS & Fetch API)
 * ==============================================================================
 *
 * HOW THE JAVASCRIPT FETCH API WORKS:
 * 1. The browser's native `fetch()` function sends HTTP requests across the network
 *    to our Golang REST API backend (e.g. GET, POST, PUT, DELETE).
 * 2. `fetch()` returns a JavaScript "Promise" that resolves into a Response object.
 * 3. We use `async/await` syntax to write asynchronous network calls cleanly.
 * 4. For POST and PUT, we include headers: { "Content-Type": "application/json" }
 *    and convert JS objects to JSON strings using `JSON.stringify(payload)`.
 * 5. We parse the Golang JSON response with `await response.json()`.
 */

// If frontend is served by Gin (same origin), API_BASE_URL can simply be "/api/employees"
// If frontend is opened via file:// or Live Server on another port (e.g. 5500), it falls back to http://localhost:8080/api/employees
const API_BASE_URL = window.location.origin.includes("localhost:8080") || window.location.origin.includes("127.0.0.1:8080")
    ? "/api/employees"
    : "http://localhost:8080/api/employees";

// Local cache of employees for instant search filtering
let employeesCache = [];
let employeeToDelete = null;

// DOM Elements
const employeeTableBody = document.getElementById("employeeTableBody");
const emptyState = document.getElementById("emptyState");
const loadingState = document.getElementById("loadingState");
const employeeCountBadge = document.getElementById("employeeCount");
const searchInput = document.getElementById("searchInput");

// Form Modal Elements
const employeeModal = document.getElementById("employeeModal");
const employeeForm = document.getElementById("employeeForm");
const modalTitle = document.getElementById("modalTitle");
const employeeIdInput = document.getElementById("employeeId");
const employeeNameInput = document.getElementById("employeeName");
const employeeNPWPInput = document.getElementById("employeeNPWP");
const employeeAddressInput = document.getElementById("employeeAddress");
const openAddModalBtn = document.getElementById("openAddModalBtn");
const closeModalBtn = document.getElementById("closeModalBtn");
const cancelModalBtn = document.getElementById("cancelModalBtn");
const saveEmployeeBtn = document.getElementById("saveEmployeeBtn");

// Delete Modal Elements
const deleteModal = document.getElementById("deleteModal");
const deleteEmployeeName = document.getElementById("deleteEmployeeName");
const deleteEmployeeId = document.getElementById("deleteEmployeeId");
const confirmDeleteBtn = document.getElementById("confirmDeleteBtn");
const cancelDeleteBtn = document.getElementById("cancelDeleteBtn");
const closeDeleteModalBtn = document.getElementById("closeDeleteModalBtn");

// Toast Notification Element
const toast = document.getElementById("toast");
const toastIcon = document.getElementById("toastIcon");
const toastMessage = document.getElementById("toastMessage");
let toastTimeout = null;

// ==============================================================================
// 1. Initial Page Load & Event Listeners
// ==============================================================================
document.addEventListener("DOMContentLoaded", () => {
    fetchEmployees();

    // Modal open / close listeners
    openAddModalBtn.addEventListener("click", openAddModal);
    closeModalBtn.addEventListener("click", closeModal);
    cancelModalBtn.addEventListener("click", closeModal);

    // Close modal if user clicks outside the modal card
    employeeModal.addEventListener("click", (e) => {
        if (e.target === employeeModal) closeModal();
    });

    // Delete confirmation modal listeners
    cancelDeleteBtn.addEventListener("click", closeDeleteModal);
    closeDeleteModalBtn.addEventListener("click", closeDeleteModal);
    deleteModal.addEventListener("click", (e) => {
        if (e.target === deleteModal) closeDeleteModal();
    });
    confirmDeleteBtn.addEventListener("click", handleConfirmDelete);

    // Form submit listener (handles both Create and Update)
    employeeForm.addEventListener("submit", handleFormSubmit);

    // Search filter listener
    searchInput.addEventListener("input", handleSearch);
});

// ==============================================================================
// 2. READ: Fetch All Employees (GET /api/employees)
// ==============================================================================
/**
 * Sends a GET request to the Golang backend to fetch all employee records from MySQL.
 */
async function fetchEmployees() {
    showLoading(true);
    try {
        const response = await fetch(API_BASE_URL, {
            method: "GET",
            headers: { "Accept": "application/json" }
        });

        const result = await response.json();

        if (!response.ok || !result.success) {
            throw new Error(result.message || "Failed to load employees.");
        }

        employeesCache = result.data || [];
        renderTable(employeesCache);
    } catch (error) {
        console.error("Error fetching employees:", error);
        showToast(`Error: ${error.message}`, "error");
        renderTable([]);
    } finally {
        showLoading(false);
    }
}

// ==============================================================================
// 3. Render Table & Empty State
// ==============================================================================
/**
 * Renders employee rows into the <tbody> element.
 */
function renderTable(employees) {
    employeeTableBody.innerHTML = "";
    employeeCountBadge.textContent = `${employees.length} Total`;

    if (employees.length === 0) {
        emptyState.classList.remove("hidden");
        return;
    }

    emptyState.classList.add("hidden");

    employees.forEach(emp => {
        const tr = document.createElement("tr");

        tr.innerHTML = `
            <td><strong>#${escapeHtml(emp.id)}</strong></td>
            <td><strong>${escapeHtml(emp.name)}</strong></td>
            <td><code>${escapeHtml(emp.npwp)}</code></td>
            <td>${escapeHtml(emp.address)}</td>
            <td>
                <div class="table-actions">
                    <button class="btn btn-sm btn-outline-primary" onclick="openEditModal(${emp.id})">
                        ✏️ Edit
                    </button>
                    <button class="btn btn-sm btn-outline-danger" onclick="openDeleteModal(${emp.id}, '${escapeForJsString(emp.name)}')">
                        🗑️ Delete
                    </button>
                </div>
            </td>
        `;

        employeeTableBody.appendChild(tr);
    });
}

// ==============================================================================
// 4. CREATE & UPDATE: Form Submission (POST & PUT)
// ==============================================================================
/**
 * Handles form submission:
 * - If employeeId is empty => CREATE (POST /api/employees)
 * - If employeeId exists => UPDATE (PUT /api/employees/:id)
 */
async function handleFormSubmit(e) {
    e.preventDefault();

    const id = employeeIdInput.value.trim();
    const name = employeeNameInput.value.trim();
    const npwp = employeeNPWPInput.value.trim();
    const address = employeeAddressInput.value.trim();

    // Client-side validation
    if (!name || !npwp || !address) {
        showToast("Please fill in all required fields.", "error");
        return;
    }

    const payload = { name, npwp, address };
    const isEditMode = Boolean(id);

    // Determine HTTP method and endpoint
    const url = isEditMode ? `${API_BASE_URL}/${id}` : API_BASE_URL;
    const method = isEditMode ? "PUT" : "POST";

    setSaveButtonLoading(true);

    try {
        // Send JSON data to Golang backend using fetch API
        const response = await fetch(url, {
            method: method,
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json"
            },
            body: JSON.stringify(payload)
        });

        const result = await response.json();

        if (!response.ok || !result.success) {
            throw new Error(result.message || "Failed to save employee.");
        }

        // Show success notification
        showToast(
            isEditMode ? "Employee updated successfully!" : "Employee added successfully!",
            "success"
        );

        // Close modal and refresh the table data without full page reload
        closeModal();
        fetchEmployees();
    } catch (error) {
        console.error("Save error:", error);
        showToast(`Error: ${error.message}`, "error");
    } finally {
        setSaveButtonLoading(false);
    }
}

// ==============================================================================
// 5. UPDATE: Open Edit Modal & Populate Data
// ==============================================================================
/**
 * Opens modal populated with selected employee's data.
 */
function openEditModal(id) {
    const employee = employeesCache.find(emp => emp.id === id);
    if (!employee) {
        showToast("Employee record not found in cache.", "error");
        return;
    }

    modalTitle.textContent = "Edit Employee";
    employeeIdInput.value = employee.id;
    employeeNameInput.value = employee.name;
    employeeNPWPInput.value = employee.npwp;
    employeeAddressInput.value = employee.address;

    employeeModal.classList.remove("hidden");
    employeeNameInput.focus();
}

/**
 * Opens blank modal for creating a new employee.
 */
function openAddModal() {
    modalTitle.textContent = "Add New Employee";
    employeeForm.reset();
    employeeIdInput.value = "";

    employeeModal.classList.remove("hidden");
    employeeNameInput.focus();
}

function closeModal() {
    employeeModal.classList.add("hidden");
    employeeForm.reset();
    employeeIdInput.value = "";
}

// ==============================================================================
// 6. DELETE: Confirmation & Request (DELETE /api/employees/:id)
// ==============================================================================
/**
 * Opens confirmation modal before deleting.
 */
function openDeleteModal(id, name) {
    employeeToDelete = { id, name };
    deleteEmployeeId.textContent = id;
    deleteEmployeeName.textContent = name;
    deleteModal.classList.remove("hidden");
}

function closeDeleteModal() {
    deleteModal.classList.add("hidden");
    employeeToDelete = null;
}

/**
 * Sends DELETE request to Golang backend.
 */
async function handleConfirmDelete() {
    if (!employeeToDelete) return;

    const id = employeeToDelete.id;
    confirmDeleteBtn.disabled = true;
    confirmDeleteBtn.textContent = "Deleting...";

    try {
        const response = await fetch(`${API_BASE_URL}/${id}`, {
            method: "DELETE",
            headers: { "Accept": "application/json" }
        });

        const result = await response.json();

        if (!response.ok || !result.success) {
            throw new Error(result.message || "Failed to delete employee.");
        }

        showToast("Employee deleted successfully!", "success");
        closeDeleteModal();
        // Refresh table automatically
        fetchEmployees();
    } catch (error) {
        console.error("Delete error:", error);
        showToast(`Error: ${error.message}`, "error");
    } finally {
        confirmDeleteBtn.disabled = false;
        confirmDeleteBtn.textContent = "Delete";
    }
}

// ==============================================================================
// 7. Search & Filter
// ==============================================================================
function handleSearch(e) {
    const term = e.target.value.toLowerCase().trim();
    if (!term) {
        renderTable(employeesCache);
        return;
    }

    const filtered = employeesCache.filter(emp => {
        return (
            emp.name.toLowerCase().includes(term) ||
            emp.npwp.toLowerCase().includes(term) ||
            emp.address.toLowerCase().includes(term) ||
            String(emp.id).includes(term)
        );
    });

    renderTable(filtered);
}

// ==============================================================================
// 8. UI Helpers: Loading, Toast, Escaping
// ==============================================================================
function showLoading(show) {
    if (show) {
        loadingState.classList.remove("hidden");
        employeeTableBody.innerHTML = "";
        emptyState.classList.add("hidden");
    } else {
        loadingState.classList.add("hidden");
    }
}

function setSaveButtonLoading(loading) {
    const spinner = document.getElementById("saveBtnSpinner");
    const text = document.getElementById("saveBtnText");
    if (loading) {
        saveEmployeeBtn.disabled = true;
        spinner.classList.remove("hidden");
        text.textContent = "Saving...";
    } else {
        saveEmployeeBtn.disabled = false;
        spinner.classList.add("hidden");
        text.textContent = "Save Employee";
    }
}

function showToast(message, type = "info") {
    if (toastTimeout) clearTimeout(toastTimeout);

    toastMessage.textContent = message;
    toast.className = `toast toast-${type}`;

    if (type === "success") {
        toastIcon.textContent = "✅";
    } else if (type === "error") {
        toastIcon.textContent = "❌";
    } else {
        toastIcon.textContent = "ℹ️";
    }

    toast.classList.remove("hidden");

    toastTimeout = setTimeout(() => {
        toast.classList.add("hidden");
    }, 3500);
}

function escapeHtml(str) {
    if (str === null || str === undefined) return "";
    return String(str)
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

function escapeForJsString(str) {
    if (!str) return "";
    return String(str).replace(/'/g, "\\'").replace(/"/g, '\\"');
}
