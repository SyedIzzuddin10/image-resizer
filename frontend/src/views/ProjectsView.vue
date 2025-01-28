<template>
  <div class="projects-container">
    <header class="header">
      <div class="header-content">
        <h1>My Projects</h1>
        <p class="header-subtitle">
          Manage and organize your image resizing projects
        </p>
      </div>
      <button
        @click="router.push('/projects/new')"
        class="btn btn-primary new-project-btn"
      >
        <span class="icon">+</span>
        New Project
      </button>
    </header>

    <div v-if="loading" class="loading">
      <span class="loading-spinner"></span>
      <p>Loading your projects...</p>
    </div>

    <div v-else-if="error" class="error">
      <span class="icon">⚠️</span>
      <p>{{ error }}</p>
    </div>

    <div v-else-if="projects.length === 0" class="empty-state">
      <div class="empty-state-icon">📁</div>
      <h2>No Projects Yet</h2>
      <p>Create your first project to start resizing images</p>
      <button
        @click="router.push('/projects/new')"
        class="btn btn-primary create-first-btn"
      >
        <span class="icon">+</span>
        Create Your First Project
      </button>
    </div>

    <div v-else class="projects-grid">
      <div
        v-for="project in projects"
        :key="project.id"
        class="project-card card"
      >
        <div class="project-header">
          <h3>{{ project.name }}</h3>
          <span class="badge">
            {{
              new Date(project.creation_time).toLocaleDateString("en-GB", {
                month: "short",
                day: "numeric",
                year: "numeric",
              })
            }}
          </span>
        </div>

        <p class="description">
          {{ project.description || "No description provided" }}
        </p>

        <div class="project-actions">
          <button
            @click="viewProject(project)"
            class="btn btn-primary view-btn"
          >
            View Project
          </button>
          <button
            @click="deleteProject(project.id)"
            class="btn btn-danger delete-btn"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import {
  ListProjects,
  CreateProject,
  DeleteProject,
  MessageDialog,
} from "../../wailsjs/go/main/App";

const router = useRouter();
const projects = ref([]);
const loading = ref(true);
const error = ref("");

const loadProjects = async () => {
  try {
    loading.value = true;
    error.value = "";
    const result = await ListProjects();
    projects.value = result || [];
  } catch (err) {
    error.value = "Failed to load projects";
    console.error("Error loading projects:", err);
    projects.value = [];
  } finally {
    loading.value = false;
  }
};

const deleteProject = async (id) => {
  // Use the native dialog
  const confirmed = await MessageDialog(
    "Delete Project",
    "Are you sure you want to delete this project?",
    "question"
  );

  if (confirmed) {
    try {
      await DeleteProject(id);
      projects.value = projects.value.filter((p) => p.id !== id);
    } catch (err) {
      error.value = "Failed to delete project";
      console.error("Error deleting project:", err);
    }
  }
};

const viewProject = (project) => {
  router.push(`/projects/${project.id}`);
};

onMounted(() => {
  loadProjects();
});
</script>

<style scoped>
.projects-container {
  animation: fadeIn 0.3s ease-out;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-xl);
  padding-bottom: var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.header-content h1 {
  margin-bottom: var(--spacing-xs);
}

.header-subtitle {
  color: var(--text-secondary);
  font-size: 0.875rem;
  margin: 0;
}

.new-project-btn {
  font-weight: 500;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--spacing-lg);
  animation: slideUp 0.3s ease-out;
}

.project-card {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  background: var(--surface-color);
}

.project-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-sm);
}

.project-header h3 {
  margin: 0;
  font-size: 1.125rem;
}

.description {
  color: var(--text-secondary);
  font-size: 0.875rem;
  line-height: 1.5;
  margin: 0;
  flex: 1;
}

.project-actions {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: var(--spacing-sm);
  margin-top: auto;
}

.view-btn,
.delete-btn {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
}

.icon {
  font-size: 1rem;
}

.empty-state {
  text-align: center;
  padding: var(--spacing-xl) var(--spacing-lg);
  background: var(--surface-color);
  border-radius: var(--border-radius-lg);
  border: 1px solid var(--border-color);
}

.empty-state h2 {
  margin-bottom: var(--spacing-sm);
  font-size: 1.25rem;
}

.empty-state p {
  color: var(--text-secondary);
  margin-bottom: var(--spacing-lg);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-xl);
}

.loading p {
  color: var(--text-secondary);
  margin: 0;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
