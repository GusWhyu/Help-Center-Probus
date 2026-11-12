import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import SystemDetailView from '../views/SystemDetailView.vue'
import SolutionsView from '../views/SolutionsView.vue'
import GuideView from '../views/GuideView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
  },
  {
    path: '/system/:id',
    name: 'system-detail',
    component: SystemDetailView,
  },
  {
    path: '/system/:systemId/guide/:moduleId?',
    name: 'system-guide',
    component: GuideView,
  },
  {
    path: '/guide/:systemId?/:moduleId?',
    name: 'guide',
    component: GuideView,
  },
  {
    path: '/solutions',
    name: 'solutions',
    component: SolutionsView,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
