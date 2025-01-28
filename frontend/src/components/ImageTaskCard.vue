<template>
  <div class="task-card card">
    <div class="task-header">
      <span :class="['status-indicator', task.status]">
        {{ getStatusIcon(task.status) }}
        {{ formatStatus(task.status) }}
      </span>
      <span class="task-date">
        {{
          new Date(task.scheduled_for).toLocaleString("en-MY", {
            timeZone: "Asia/Kuala_Lumpur",
            dateStyle: "short",
            timeStyle: "short",
          })
        }}
      </span>
    </div>
    <div class="task-content">
      <div class="task-image">
        <img
          :src="task.imagePreviewUrl || placeholderImage"
          alt="Task preview"
          @error="handleImageError"
        />
      </div>
      <div class="task-info">
        <div class="task-dimensions">
          <span class="dimension-label">Target Size:</span>
          <span class="dimension-value"
            >{{ task.target_width }}×{{ task.target_height }}</span
          >
        </div>
        <button
          v-if="task.status === 'completed'"
          class="btn btn-secondary view-resized-btn"
          @click="$emit('view-resized', task)"
        >
          View Resized Image
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from "vue";
import placeholderImage from "../assets/placeholder-image.svg";

const props = defineProps({
  task: {
    type: Object,
    required: true,
  },
});

defineEmits(["view-resized"]);

const getStatusIcon = (status) => {
  const icons = {
    pending: "⏳",
    processing: "⚙️",
    completed: "✅",
    failed: "❌",
  };
  return icons[status] || "❓";
};

const formatStatus = (status) => {
  return status.charAt(0).toUpperCase() + status.slice(1);
};

const handleImageError = (event) => {
  if (event.target.src !== placeholderImage) {
    event.target.src = placeholderImage;
  }
};
</script>

<style scoped>
.task-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
}

.task-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid #eee;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.9em;
}

.status-indicator.pending {
  color: #f59e0b;
}
.status-indicator.processing {
  color: #3b82f6;
}
.status-indicator.completed {
  color: #10b981;
}
.status-indicator.failed {
  color: #ef4444;
}

.task-date {
  font-size: 0.85em;
  color: #666;
}

.task-content {
  padding: 12px;
}

.task-image {
  width: 100%;
  height: 160px;
  overflow: hidden;
  border-radius: 4px;
  margin-bottom: 12px;
}

.task-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.task-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-dimensions {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9em;
}

.dimension-label {
  color: #666;
}

.dimension-value {
  font-weight: 500;
}

.view-resized-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px;
  border-radius: 6px;
  background: #e2e8f0;
  border: none;
  cursor: pointer;
  transition: background 0.2s;
  color: #000000;
  font-weight: 500;
}

.view-resized-btn:hover {
  background: #cbd5e1;
  color: #000000;
}
</style>
