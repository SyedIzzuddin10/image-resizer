<template>
  <div class="login-container">
    <div class="login-content">
      <div class="brand-section">
        <div class="logo">R</div>
        <h1 class="brand-name">Resizer</h1>
        <p class="brand-tagline">
          Effortlessly resize and optimize your images
        </p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="username" class="form-label">Username</label>
          <div class="input-container">
            <span class="input-icon">👤</span>
            <input
              type="text"
              id="username"
              v-model="username"
              required
              placeholder="Enter your username"
              class="form-input"
              :class="{ 'error-input': error }"
            />
          </div>
        </div>

        <div class="form-group">
          <label for="password" class="form-label">Password</label>
          <div class="input-container">
            <span class="input-icon">🔒</span>
            <input
              :type="showPassword ? 'text' : 'password'"
              id="password"
              v-model="password"
              required
              placeholder="Enter your password"
              class="form-input"
              :class="{ 'error-input': error }"
            />
            <button
              type="button"
              class="toggle-password"
              @click="showPassword = !showPassword"
            >
              {{ showPassword ? "👁️" : "👁️‍🗨️" }}
            </button>
          </div>
        </div>

        <transition name="fade">
          <p v-if="error" class="error-message">
            <span class="icon">⚠️</span>
            {{ error }}
          </p>
        </transition>

        <button
          type="submit"
          class="btn btn-primary btn-block"
          :disabled="isLoading"
        >
          <span v-if="isLoading" class="loading-spinner"></span>
          {{ isLoading ? "Signing in..." : "Sign in" }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { Login } from "../../wailsjs/go/main/App";

const router = useRouter();
const username = ref("");
const password = ref("");
const error = ref("");
const isLoading = ref(false);
const showPassword = ref(false);

const handleLogin = async () => {
  try {
    isLoading.value = true;
    error.value = "";

    const response = await Login(username.value, password.value);
    if (response && response.user) {
      localStorage.setItem("user", JSON.stringify(response.user));
      localStorage.setItem("isAuthenticated", "true");
      router.push("/projects");
    }
  } catch (err) {
    error.value = "Invalid credentials. Please try again.";
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-md);
  background-color: var(--background-color);
}

.login-content {
  width: 100%;
  max-width: 400px;
  padding: var(--spacing-xl);
  background: var(--surface-color);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-lg);
  animation: fadeIn 0.5s ease-out;
}

.brand-section {
  text-align: center;
  margin-bottom: var(--spacing-xl);
}

.logo {
  width: 4rem;
  height: 4rem;
  background: var(--accent-gradient);
  border-radius: var(--border-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 2rem;
  margin: 0 auto var(--spacing-md);
  box-shadow: var(--shadow-md);
}

.brand-name {
  font-size: 2rem;
  background: var(--accent-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: var(--spacing-xs);
}

.brand-tagline {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.form-group {
  margin: 0;
}

.form-label {
  display: block;
  margin-bottom: var(--spacing-xs);
  color: var(--text-primary);
  font-weight: 500;
  font-size: 0.875rem;
}

.input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: var(--spacing-sm);
  color: var(--text-secondary);
  font-size: 1rem;
}

.form-input {
  width: 100%;
  padding: 0.75rem var(--spacing-md) 0.75rem 2.5rem;
  background-color: var(--background-color);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-full);
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--primary-light);
}

.toggle-password {
  position: absolute;
  right: var(--spacing-sm);
  background: none;
  border: none;
  padding: var(--spacing-xs);
  cursor: pointer;
  color: var(--text-secondary);
  transition: color 0.2s ease;
}

.toggle-password:hover {
  color: var(--primary-color);
}

.error-input {
  border-color: var(--danger-color);
}

.error-message {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  color: var(--danger-color);
  font-size: 0.875rem;
  background-color: rgba(239, 68, 68, 0.1);
  padding: var(--spacing-sm);
  border-radius: var(--border-radius-md);
  margin: 0;
}

.btn-block {
  width: 100%;
  padding: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  font-weight: 500;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-5px);
}
</style>
