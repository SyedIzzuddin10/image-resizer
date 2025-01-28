import { createRouter, createWebHashHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import ProjectsView from "../views/ProjectsView.vue";
import ProjectDetailView from "../views/ProjectDetailView.vue";
import NewProjectView from "../views/NewProjectView.vue";
import UploadImagesView from "../views/UploadImagesView.vue";

const routes = [
  {
    path: "/",
    redirect: "/login",
  },
  {
    path: "/login",
    name: "login",
    component: LoginView,
  },
  {
    path: "/projects",
    name: "projects",
    component: ProjectsView,
    meta: { requiresAuth: true },
  },
  {
    path: "/projects/new",
    name: "new-project",
    component: NewProjectView,
    meta: { requiresAuth: true },
  },
  {
    path: "/projects/:id",
    name: "project-detail",
    component: ProjectDetailView,
    meta: { requiresAuth: true },
  },
  {
    path: "/projects/:id/upload",
    name: "upload-images",
    component: UploadImagesView,
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

// Navigation guard
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem("isAuthenticated") === "true";

  // If trying to access a protected route without authentication
  if (to.meta.requiresAuth && !isAuthenticated) {
    next("/login");
    return;
  }

  // If trying to access login while authenticated
  if (to.path === "/login" && isAuthenticated) {
    next("/projects");
    return;
  }

  next();
});

export default router;
