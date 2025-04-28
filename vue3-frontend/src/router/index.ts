import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import login from '../login/sign_in.vue'
import register from '../login/sign_up.vue'
import user_main from '../user_ui/user.vue'
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
    {
      path:'/',
      name:'login',
      component: login
    },
    {
      path:'/register',
      name:'register',
      component: register
    },
    {
      path:'/user/main',
      name:'user_main',
      component: user_main
    }
  ],
})

export default router
