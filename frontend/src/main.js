import { createApp } from "vue";
import App from "./App.vue";
import "./style.css";
import router from "./router";

// Buang semua authentication state on app start
localStorage.removeItem("isAuthenticated");
localStorage.removeItem("user");

const app = createApp(App);
app.use(router);
app.mount("#app");
