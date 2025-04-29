import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import login from '../login/sign_in.vue'
import register from '../login/sign_up.vue'
import user_main from '../user_ui/user.vue'
import user_profile from '../user_ui/user_profile.vue'
import user_model from '../user_ui/user_model.vue'
import user_task from '../user_ui/user_task.vue'
import admin_main from '../admin_ui/admin.vue'
import admin_profile from '../admin_ui/admin_profile.vue'
import admin_model from '../admin_ui/admin_model.vue'
import admin_task from '../admin_ui/admin_task.vue'
import admin_user from '../admin_ui/admin_user.vue'

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
    },
    {
      path:'/user/profile',
      name:'user_profile',
      component: user_profile
    },
    {
      path:'/user/model',
      name:'user_model',
      component: user_model
    },
    {
      path:'/user/task',
      name:'user_task',
      component: user_task
    },
    {
      path: '/admin/main',
      name: 'admin_main',
      component: admin_main
    },
    {
      path: '/admin/profile',
      name: 'admin_profile',
      component: admin_profile
    },
    {
      path: '/admin/model',
      name: 'admin_model',
      component: admin_model
    },
    {
      path: '/admin/task',
      name: 'admin_task',
      component: admin_task
    },
    {
      path: '/admin/user',
      name: 'admin_user',
      component: admin_user
    }
  ],
})

export default router
