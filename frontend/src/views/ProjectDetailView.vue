<template>
  <div class="project-detail-container">
    <header class="header">
      <div class="header-left">
        <button @click="router.push('/projects')" class="back-button">
          ← Back
        </button>
        <h1>{{ project?.name || "Loading..." }}</h1>
      </div>
      <button
        @click="router.push(`/projects/${projectId}/upload`)"
        class="btn btn-primary upload-btn"
      >
        <span class="icon">+</span>
        Upload Images
      </button>
    </header>

    <div v-if="loading" class="loading">
      <span class="loading-spinner"></span>
      <p>Loading project details...</p>
    </div>
    <div v-else-if="error" class="error">
      <span class="icon">⚠️</span>
      {{ error }}
    </div>
    <template v-else>
      <div class="project-info card">
        <div class="info-item">
          <span class="info-label">
            <span class="icon">📁</span>
            Location
          </span>
          <span class="info-value">{{ project.location }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">
            <span class="icon">📅</span>
            Created
          </span>
          <span class="info-value">
            {{
              new Date(project.creation_time).toLocaleDateString("en-GB", {
                day: "2-digit",
                month: "long",
                year: "numeric",
              })
            }}
          </span>
        </div>
      </div>

      <div class="tasks-section">
        <div class="section-header">
          <h2>Image Tasks</h2>
          <span class="badge task-count">{{ tasks?.length || 0 }} tasks</span>
        </div>

        <div v-if="!tasks || tasks.length === 0" class="empty-state">
          <div class="empty-state-icon">🖼️</div>
          <h3>No Image Tasks Yet</h3>
          <p>Start by uploading your first image to resize</p>
          <button
            @click="router.push(`/projects/${projectId}/upload`)"
            class="btn btn-primary"
          >
            <span class="icon">+</span>
            Upload Your First Image
          </button>
        </div>

        <div v-else class="tasks-grid">
          <ImageTaskCard
            v-for="task in tasks"
            :key="task.id"
            :task="task"
            @view-resized="viewResizedImage"
          />
        </div>
      </div>
    </template>

    <div
      v-if="showResizedImage"
      class="modal-overlay"
      @click="closeResizedImage"
    >
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Resized Image</h3>
          <button class="close-button" @click="closeResizedImage">×</button>
        </div>
        <div class="modal-body">
          <div class="image-container">
            <img :src="resizedImageUrl" alt="Resized Image" />
          </div>
          <div class="image-info">
            <p>
              Resized to {{ selectedTask?.target_width }}×{{
                selectedTask?.target_height
              }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  GetProject,
  GetProjectTasks,
  GetImageData,
  GetResizedImageData,
} from "../../wailsjs/go/main/App";
import ImageTaskCard from "../components/ImageTaskCard.vue";
import placeholderImage from "../assets/placeholder-image.svg";

const router = useRouter();
const route = useRoute();
const projectId = parseInt(route.params.id);

const project = ref(null);
const tasks = ref([]);
const loading = ref(true);
const error = ref("");

const showResizedImage = ref(false);
const resizedImageUrl = ref("");
const selectedTask = ref(null);

const loadProjectDetails = async () => {
  try {
    loading.value = true;
    error.value = "";

    const [projectData, tasksData] = await Promise.all([
      GetProject(projectId),
      GetProjectTasks(projectId),
    ]);

    project.value = projectData;

    tasks.value = tasksData
      ? tasksData.map((task) => ({
          ...task,
          imagePreviewUrl: placeholderImage,
        }))
      : [];

    if (tasksData) {
      for (const task of tasksData) {
        if (task.image_path) {
          try {
            console.log(
              "Loading image for task:",
              task.id,
              "Path:",
              task.image_path
            );
            const imageData = await GetImageData(task.image_path);

            if (!imageData) {
              console.error("No image data received for task:", task.id);
              continue;
            }

            const ext = task.image_path.split(".").pop().toLowerCase();
            const mimeType = ext === "png" ? "image/png" : "image/jpeg";

            const imageUrl = `data:${mimeType};base64,${imageData}`;
            console.log("Created data URL for task:", task.id);

            const taskToUpdate = tasks.value.find((t) => t.id === task.id);
            if (taskToUpdate) {
              taskToUpdate.imagePreviewUrl = imageUrl;
              console.log("Updated image URL for task:", task.id);
            }
          } catch (err) {
            console.error(`Failed to load image for task ${task.id}:`, err);
          }
        }
      }
    }
  } catch (err) {
    error.value = "Failed to load project details";
    console.error("Error loading project details:", err);
  } finally {
    loading.value = false;
  }
};

const handleImageError = (event) => {
  console.error("Image failed to load:", event.target.src);
  console.error(
    "Image natural size:",
    event.target.naturalWidth,
    "x",
    event.target.naturalHeight
  );
  if (event.target.src !== placeholderImage) {
    event.target.src = placeholderImage;
  }
};

const getStatusIcon = (status) => {
  const icons = {
    pending: "⏳",
    processing: "⚙️",
    completed: "✅",
    failed: "❌",
  };
  return icons[status] || "⚪️";
};

const formatStatus = (status) => {
  return status.charAt(0).toUpperCase() + status.slice(1);
};

const viewResizedImage = async (task) => {
  try {
    const imageData = await GetResizedImageData(task.image_path);
    if (!imageData) {
      console.error("No resized image data received");
      return;
    }

    // Get file extension dari path
    const ext = task.image_path.split(".").pop().toLowerCase();
    const mimeType = ext === "png" ? "image/png" : "image/jpeg";

    resizedImageUrl.value = `data:${mimeType};base64,${imageData}`;
    selectedTask.value = task;
    showResizedImage.value = true;
  } catch (err) {
    console.error("Failed to load resized image:", err);
  }
};

const closeResizedImage = () => {
  showResizedImage.value = false;
  resizedImageUrl.value = "";
  selectedTask.value = null;
};

onMounted(() => {
  loadProjectDetails();
});
</script>

<style scoped>
.project-detail-container {
  padding: var(--spacing-xl);
  max-width: 1200px;
  margin: 0 auto;
  animation: fadeIn 0.3s ease-out;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-xl);
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.header h1 {
  margin: 0;
  font-size: 1.75rem;
}

.back-button {
  background: none;
  border: none;
  color: #2196f3;
  cursor: pointer;
  font-size: 1rem;
}

.back-button:hover {
  color: #1976d2;
}

.upload-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.upload-btn .icon {
  font-size: 1.1em;
  line-height: 1;
}

.project-info {
  display: flex;
  gap: var(--spacing-xl);
  margin-bottom: var(--spacing-xl);
  padding: var(--spacing-lg);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.info-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
}

.info-value {
  color: var(--text-primary);
  font-weight: 500;
}

.section-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-lg);
}

.section-header h2 {
  margin: 0;
  font-size: 1.5rem;
}

.task-count {
  background-color: var(--primary-light);
  color: var(--primary-color);
}

.tasks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.task-card {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  border: 1px solid var(--border-color);
}

.task-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  background-color: var(--background-alt);
  border-bottom: 1px solid var(--border-color);
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-weight: 500;
  font-size: 0.875rem;
  padding: 4px 8px;
  border-radius: 4px;
}

.status-indicator.pending {
  background-color: #fff3dc;
  color: #b25e00;
}

.status-indicator.processing {
  background-color: #e3f2fd;
  color: #0062cc;
}

.status-indicator.completed {
  background-color: #e8f5e9;
  color: #2e7d32;
}

.status-indicator.failed {
  background-color: #ffebee;
  color: #c62828;
}

.task-date {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.task-content {
  padding: var(--spacing-md);
}

.task-image {
  width: 100%;
  aspect-ratio: 16/9;
  background: white;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: var(--spacing-md);
  border: 1px solid var(--border-color);
}

.task-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

.task-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.task-dimensions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-md);
  background-color: var(--background-alt);
  border-radius: 4px;
}

.dimension-label {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.dimension-value {
  font-weight: 500;
  color: var(--text-primary);
}

.empty-state {
  text-align: center;
  padding: var(--spacing-xl);
  background: var(--surface-color);
  border-radius: var(--border-radius-lg);
  border: 1px solid var(--border-color);
}

.empty-state-icon {
  font-size: 3rem;
  margin-bottom: var(--spacing-md);
  opacity: 0.5;
}

.empty-state h3 {
  margin-bottom: var(--spacing-xs);
  color: var(--text-primary);
}

.empty-state p {
  color: var(--text-secondary);
  margin-bottom: var(--spacing-lg);
}

.view-resized-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-xs);
  margin-top: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-md);
  background-color: var(--background-alt);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--text-primary);
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.view-resized-btn:hover {
  background-color: var(--primary-light);
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.view-resized-btn .icon {
  font-size: 1.1em;
  line-height: 1;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* Add modal styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  animation: modalFadeIn 0.2s ease-out;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: var(--text-primary);
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  line-height: 1;
}

.close-button:hover {
  background-color: var(--background-alt);
  color: var(--text-primary);
}

.modal-body {
  padding: var(--spacing-lg);
  overflow-y: auto;
}

.image-container {
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--background-alt);
  border-radius: 4px;
  padding: var(--spacing-md);
}

.image-container img {
  max-width: 100%;
  max-height: 70vh;
  object-fit: contain;
}

.image-info {
  margin-top: var(--spacing-md);
  text-align: center;
  color: var(--text-secondary);
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
