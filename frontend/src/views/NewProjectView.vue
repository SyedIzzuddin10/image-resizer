<template>
  <div class="new-project-container">
    <header class="page-header">
      <button class="back-button" @click="router.push('/projects')">
        ← Back
      </button>
      <h1>Create New Project</h1>
    </header>

    <div class="form-container card">
      <form @submit.prevent="handleSubmit" class="project-form">
        <div class="form-group">
          <label for="projectName">Project Name</label>
          <input
            type="text"
            id="projectName"
            v-model="project.name"
            required
            placeholder="Enter project name"
          />
          <small class="help-text">
            This name will be used to create the project folder
          </small>
        </div>

        <div class="form-group">
          <label for="description">Project Description</label>
          <textarea
            id="description"
            v-model="project.description"
            rows="3"
            placeholder="Enter project description (optional)"
          ></textarea>
        </div>

        <div class="action-buttons">
          <button
            type="button"
            @click="router.push('/projects')"
            class="btn secondary"
          >
            Cancel
          </button>
          <button type="submit" :disabled="isCreating" class="btn btn-primary">
            {{ isCreating ? "Creating..." : "Create Project" }}
          </button>
        </div>

        <p v-if="error" class="error-message">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { CreateProject } from "../../wailsjs/go/main/App";

const router = useRouter();
const isCreating = ref(false);
const error = ref("");

const project = ref({
  name: "",
  description: "",
});

const handleSubmit = async () => {
  try {
    isCreating.value = true;
    error.value = "";

    const newProject = await CreateProject(
      project.value.name,
      project.value.description
    );

    router.push(`/projects/${newProject.id}`);
  } catch (err) {
    error.value = "Failed to create project. Please try again.";
    console.error("Error creating project:", err);
  } finally {
    isCreating.value = false;
  }
};
</script>

<style scoped>
.new-project-container {
  padding: var(--spacing-xl);
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-xl);
}

.page-header h1 {
  margin: 0;
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-primary);
}

.back-button {
  background: none;
  border: none;
  color: #2196f3;
  cursor: pointer;
  font-size: 1rem;
  transition: transform 0.2s ease;
}

.back-button:hover {
  color: #1976d2;
  transform: translateX(-2px);
}

.form-container {
  padding: var(--spacing-xl);
  background: var(--surface-color);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.project-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  max-width: 600px;
  margin: 0 auto;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.form-group label {
  font-weight: 500;
  color: var(--text-primary);
  font-size: 0.95rem;
  margin-bottom: 4px;
}

.form-group input,
.form-group textarea {
  padding: 12px 16px;
  border: 1.5px solid var(--border-color);
  border-radius: 12px;
  font-size: 1rem;
  background-color: var(--surface-color);
  transition: all 0.2s ease;
}

.form-group input:hover,
.form-group textarea:hover {
  border-color: var(--primary-color);
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--primary-light);
}

.help-text {
  color: var(--text-secondary);
  font-size: 0.875rem;
  margin-top: 4px;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
  margin-top: var(--spacing-xl);
}

.action-buttons button {
  padding: 12px 24px;
  border-radius: 12px;
  font-weight: 500;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-buttons button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background-color: #e0e0e0 !important;
  color: #757575 !important;
  border: 1.5px solid #e0e0e0 !important;
  box-shadow: none !important;
}

.action-buttons button:disabled span {
  color: #757575 !important;
}

.action-buttons .btn.secondary {
  background-color: var(--surface-color);
  border: 1.5px solid var(--border-color);
  color: var(--text-secondary);
}

.action-buttons .btn.secondary:hover:not(:disabled) {
  background-color: var(--background-alt);
  border-color: var(--text-secondary);
  transform: translateY(-1px);
  color: var(--text-primary);
}

.action-buttons .btn.secondary:active:not(:disabled) {
  transform: translateY(0);
}

.error-message {
  color: var(--danger-color);
  background-color: var(--danger-light);
  padding: var(--spacing-md);
  border-radius: 12px;
  margin-top: var(--spacing-md);
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.error-message::before {
  content: "⚠️";
}
</style>
