import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/chat'
  },
  {
    path: '/chat',
    name: 'Chat',
    component: () => import('../views/ChatView.vue')
  },
  {
    path: '/packages',
    name: 'Packages',
    component: () => import('../views/PackagesView.vue')
  },
  {
    path: '/groups',
    name: 'Groups',
    component: () => import('../views/GroupsView.vue')
  },
  {
    path: '/installations',
    name: 'Installations',
    component: () => import('../views/InstallationsView.vue')
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('../views/SettingsView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
