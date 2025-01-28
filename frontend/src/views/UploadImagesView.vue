<template>
  <div class="upload-page">
    <header class="page-header">
      <button class="back-button" @click="router.back()">← Back</button>
      <h1>Upload Images</h1>
    </header>

    <div class="upload-container card">
      <div
        class="upload-area"
        @click="fileInput?.click()"
        @dragover.prevent
        @drop.prevent="handleDrop"
      >
        <input
          type="file"
          ref="fileInput"
          @change="handleFileSelect"
          multiple
          accept="image/*"
          class="hidden"
        />
        <div class="upload-icon">
          <span class="material-icons">cloud_upload</span>
        </div>
        <h3>Drop your images here</h3>
        <p>or click to browse files</p>
      </div>

      <div v-if="selectedFiles.length > 0" class="selected-files">
        <h3 class="section-title">Selected Images</h3>
        <div class="files-grid">
          <div
            v-for="file in selectedFiles"
            :key="file.name"
            class="file-preview"
          >
            <div class="preview-image">
              <img
                :src="getFilePreview(file)"
                :alt="file.name"
                @error="handleImageError"
              />
              <button
                class="remove-button"
                @click="removeFile(file)"
                title="Remove image"
              >
                <span class="material-icons">close</span>
              </button>
            </div>
            <div class="file-info">
              <span class="file-name">{{ file.name }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="settings-panel">
        <h3 class="section-title">Processing Settings</h3>
        <div class="settings-grid">
          <div class="form-group">
            <label for="width">Target Width (px)</label>
            <input
              type="number"
              id="width"
              v-model="uploadSettings.width"
              min="1"
              max="10000"
              placeholder="Enter width"
            />
          </div>
          <div class="form-group">
            <label for="height">Target Height (px)</label>
            <input
              type="number"
              id="height"
              v-model="uploadSettings.height"
              min="1"
              max="10000"
              placeholder="Enter height"
            />
          </div>
          <div class="form-group full-width">
            <label for="scheduleTime">Schedule Time</label>
            <input
              type="datetime-local"
              id="scheduleTime"
              v-model="uploadSettings.scheduleTime"
            />
          </div>
        </div>
      </div>

      <div class="error-message" v-if="error">
        <span class="icon">⚠️</span>
        {{ error }}
      </div>

      <div class="action-buttons">
        <button class="btn secondary" @click="router.back()">Cancel</button>
        <button
          class="btn btn-primary"
          @click="handleUpload"
          :disabled="selectedFiles.length === 0 || isUploading"
        >
          <span v-if="isUploading">
            Processing... ({{ processedCount }}/{{ selectedFiles.length }})
          </span>
          <span v-else>Upload and Process</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from "vue";
import { useRouter, useRoute } from "vue-router";
import { CreateImageTask, SaveUploadedFile } from "../../wailsjs/go/main/App";
import placeholderImage from "../assets/placeholder-image.svg";

const router = useRouter();
const route = useRoute();
const projectId = parseInt(route.params.id);

const fileInput = ref(null);
const selectedFiles = ref([]);
const isUploading = ref(false);
const error = ref("");
const processedCount = ref(0);
const filePreviewUrls = ref(new Map());

const uploadSettings = ref({
  width: 800,
  height: 600,
  scheduleTime: (() => {
    const now = new Date();
    // Date untuk Malaysia timezone (UTC+8)
    const malaysiaTime = new Date(now.getTime() + 8 * 60 * 60 * 1000);
    return malaysiaTime.toISOString().slice(0, 16);
  })(),
});

const handleDrop = (event) => {
  event.preventDefault();
  const files = Array.from(event.dataTransfer.files).filter((file) =>
    file.type.startsWith("image/")
  );
  if (files.length === 0) return;

  selectedFiles.value = files;

  files.forEach((file) => {
    try {
      filePreviewUrls.value.set(file, URL.createObjectURL(file));
    } catch (err) {
      console.error(`Error creating preview for ${file.name}:`, err);
      filePreviewUrls.value.set(file, placeholderImage);
    }
  });
};

const handleFileSelect = (event) => {
  const files = Array.from(event.target.files);
  selectedFiles.value = files;

  // Create object url untuk previews
  files.forEach((file) => {
    try {
      filePreviewUrls.value.set(file, URL.createObjectURL(file));
    } catch (err) {
      console.error(`Error creating preview for ${file.name}:`, err);
      filePreviewUrls.value.set(file, placeholderImage);
    }
  });
};

const getFilePreview = (file) => {
  const url = filePreviewUrls.value.get(file);
  return url || placeholderImage;
};

const handleImageError = (event) => {
  if (event.target.src !== placeholderImage) {
    event.target.src = placeholderImage;
  }
};

const removeFile = (fileToRemove) => {
  const previewUrl = filePreviewUrls.value.get(fileToRemove);
  if (previewUrl && previewUrl !== placeholderImage) {
    URL.revokeObjectURL(previewUrl);
  }
  filePreviewUrls.value.delete(fileToRemove);
  selectedFiles.value = selectedFiles.value.filter(
    (file) => file !== fileToRemove
  );
};

const handleUpload = async () => {
  if (selectedFiles.value.length === 0) return;

  try {
    isUploading.value = true;
    error.value = "";
    processedCount.value = 0;

    for (const file of selectedFiles.value) {
      try {
        const buffer = await file.arrayBuffer();
        const fileData = Array.from(new Uint8Array(buffer));

        const savedPath = await SaveUploadedFile(
          projectId,
          fileData,
          file.name
        );

        // Create a Date object from the local time input and convert to UTC ISO string
        const scheduledDate = new Date(uploadSettings.value.scheduleTime);
        const scheduledISOString = scheduledDate.toISOString();

        const task = await CreateImageTask(
          projectId,
          savedPath,
          uploadSettings.value.width,
          uploadSettings.value.height,
          scheduledISOString
        );
        processedCount.value++;
      } catch (err) {
        console.error(`Error processing ${file.name}:`, err);
        error.value = `Error processing ${file.name}: ${err.message}`;
      }
    }

    router.push(`/projects/${projectId}`);
  } catch (err) {
    error.value = "Failed to process images. Please try again.";
    console.error("Error processing images:", err);
  } finally {
    isUploading.value = false;
  }
};

// Clearkan object URLs when component destroyed
onBeforeUnmount(() => {
  for (const [_, url] of filePreviewUrls.value.entries()) {
    if (url && url !== placeholderImage) {
      URL.revokeObjectURL(url);
    }
  }
});
</script>

<style scoped>
.upload-page {
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

.upload-container {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xl);
  padding: var(--spacing-xl);
  background: var(--surface-color);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.upload-area {
  border: 2px dashed var(--border-color);
  border-radius: 16px;
  padding: var(--spacing-xl);
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: var(--background-alt);
}

.upload-area:hover {
  border-color: var(--primary-color);
  background: var(--primary-light);
}

.upload-icon {
  margin-bottom: var(--spacing-md);
}

.upload-icon .material-icons {
  font-size: 48px;
  color: var(--primary-color);
}

.upload-area h3 {
  margin: 0;
  margin-bottom: var(--spacing-xs);
  color: var(--text-primary);
  font-weight: 600;
}

.upload-area p {
  margin: 0;
  color: var(--text-secondary);
}

.section-title {
  margin: 0;
  margin-bottom: var(--spacing-md);
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
}

.files-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--spacing-md);
}

.file-preview {
  border-radius: 12px;
  overflow: hidden;
  background: var(--background-alt);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.file-preview:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.preview-image {
  position: relative;
  aspect-ratio: 16/9;
}

.preview-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.remove-button {
  position: absolute;
  top: 8px;
  right: 8px;
  background: rgba(0, 0, 0, 0.5);
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  color: white;
  backdrop-filter: blur(4px);
}

.remove-button:hover {
  background: rgba(255, 59, 48, 0.8);
  transform: scale(1.1);
}

.file-info {
  padding: var(--spacing-sm) var(--spacing-md);
}

.file-name {
  font-size: 0.9rem;
  color: var(--text-primary);
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.settings-panel {
  background: var(--background-alt);
  padding: var(--spacing-lg);
  border-radius: 12px;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.form-group.full-width {
  grid-column: 1 / -1;
}

.form-group label {
  font-weight: 500;
  color: var(--text-primary);
  font-size: 0.95rem;
}

.form-group input {
  padding: 12px 16px;
  border: 1.5px solid var(--border-color);
  border-radius: 12px;
  font-size: 1rem;
  background-color: var(--surface-color);
  transition: all 0.2s ease;
}

.form-group input:hover {
  border-color: var(--primary-color);
}

.form-group input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--primary-light);
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
  margin-top: var(--spacing-lg);
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
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.hidden {
  display: none;
}
</style>
