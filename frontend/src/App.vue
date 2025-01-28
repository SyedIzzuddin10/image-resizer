<script setup>
import { onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

const handleLogout = () => {
  // Buang semua auth-related data
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  localStorage.removeItem("isAuthenticated");

  router.push("/login");
};

onMounted(() => {
  const token = localStorage.getItem("token");
  if (!token && router.currentRoute.value.meta.requiresAuth) {
    router.push("/login");
  }
});
</script>

<template>
  <div class="app-container">
    <nav v-if="$route.meta.requiresAuth" class="navbar">
      <div class="container nav-content">
        <router-link to="/projects" class="nav-brand">
          <div class="brand-content">
            <div class="logo">R</div>
            <h1>Resizer</h1>
          </div>
        </router-link>
        <button @click="handleLogout" class="btn btn-secondary logout-btn">
          <span class="icon">↪</span>
          Logout
        </button>
      </div>
    </nav>
    <main :class="['main-content', { 'with-nav': $route.meta.requiresAuth }]">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<style>
/* Global styles are now in style.css */
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--background-color);
}

.navbar {
  background-color: var(--surface-color);
  border-bottom: 1px solid var(--border-color);
  padding: var(--spacing-md) 0;
  position: sticky;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(8px);
  background-color: rgba(255, 255, 255, 0.9);
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-brand {
  text-decoration: none;
}

.brand-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.logo {
  width: 2rem;
  height: 2rem;
  background: var(--accent-gradient);
  border-radius: var(--border-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 1.25rem;
}

.nav-brand h1 {
  font-size: 1.25rem;
  background: var(--accent-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.logout-btn {
  font-size: 0.875rem;
  padding: 0.5rem 1rem;
}

.icon {
  font-size: 1rem;
}

.main-content {
  flex: 1;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--spacing-xl);
  transition: padding 0.3s ease;
}

.main-content.with-nav {
  padding-top: var(--spacing-lg);
}

/* Dialog Styles */
.dialog {
  background: var(--surface-color);
  border-radius: var(--border-radius-lg);
  padding: var(--spacing-lg);
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--border-color);
}

.dialog-header {
  margin-bottom: var(--spacing-lg);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-lg);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--border-color);
}

/* Form Styles */
.form-container {
  max-width: 480px;
  margin: var(--spacing-xl) auto;
  padding: var(--spacing-xl);
  background: var(--surface-color);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--border-color);
}

.form-title {
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
  text-align: center;
  font-size: 1.75rem;
}

.form-subtitle {
  text-align: center;
  color: var(--text-secondary);
  margin-bottom: var(--spacing-lg);
  font-size: 0.875rem;
}

.form-group {
  margin-bottom: var(--spacing-md);
}

.form-label {
  display: block;
  margin-bottom: var(--spacing-xs);
  color: var(--text-secondary);
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: var(--spacing-sm);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  font-family: inherit;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.1);
}

.form-error {
  color: var(--danger-color);
  font-size: 0.875rem;
  margin-top: var(--spacing-xs);
}

/* Animation Classes */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-up-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.slide-up-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>
