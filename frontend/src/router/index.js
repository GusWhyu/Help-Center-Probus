import { createRouter, createWebHistory } from 'vue-router'

const HomeView = () => import('../views/HomeView.vue')
const SystemDetailView = () => import('../views/SystemDetailView.vue')
const SolutionsView = () => import('../views/SolutionsView.vue')
const SolutionDetailView = () => import('../views/SolutionDetailView.vue')
const GuideView = () => import('../views/GuideView.vue')

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
  {
    path: '/solutions/:id',
    name: 'solution-detail',
    component: SolutionDetailView,
  },
  {
    path: '/solution/:id',
    redirect: to => `/solutions/${to.params.id}`,
  },
  {
    path: '/contact',
    name: 'contact',
    beforeEnter() {
      window.location.href = 'https://probussystem.com/contact'
    },
  },
  {
    path: '/hubungi-kami',
    beforeEnter() {
      window.location.href = 'https://probussystem.com/contact'
    },
  },
  {
    path: '/contact-us',
    beforeEnter() {
      window.location.href = 'https://probussystem.com/contact'
    },
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
