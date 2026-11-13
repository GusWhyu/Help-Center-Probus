<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import {
  X,
  ExternalLink,
  Headphones,
  ShoppingCart,
  GraduationCap,
  MessageCircle
} from 'lucide-vue-next'

const isOpen = ref(false)
const widgetRef = ref(null)

const contacts = [
  {
    id: 1,
    name: 'CS Probussystem',
    role: 'Customer Service & Bantuan Teknis',
    numberDisplay: '+62 819-1087-5514',
    waNumber: '6281910875514',
    badge: 'CS Support',
    badgeClass: 'bg-emerald-50 text-emerald-700 border-emerald-200',
    icon: Headphones,
    defaultMessage: 'Halo CS Probussystem, saya butuh bantuan terkait kendala sistem Probus.',
  },
  {
    id: 2,
    name: 'Msm - E Commerce',
    role: 'Online Booking, Channel & Integrasi',
    numberDisplay: '+62 878-6144-6124',
    waNumber: '6287861446124',
    badge: 'E-Commerce',
    badgeClass: 'bg-sky-50 text-sky-700 border-sky-200',
    icon: ShoppingCart,
    defaultMessage: 'Halo Msm E-Commerce, saya ingin konsultasi terkait layanan integrasi & E-Commerce Probus.',
  },
  {
    id: 3,
    name: 'Trainer Probussystem',
    role: 'Panduan Operasional & Training Modul',
    numberDisplay: '+62 819-1087-5514',
    waNumber: '6281910875514',
    badge: 'Trainer',
    badgeClass: 'bg-indigo-50 text-indigo-700 border-indigo-200',
    icon: GraduationCap,
    defaultMessage: 'Halo Trainer Probussystem, saya ingin panduan dan jadwal training modul sistem Probus.',
  },
]

function toggleWidget() {
  isOpen.value = !isOpen.value
}

function closeWidget() {
  isOpen.value = false
}

function openWhatsApp(contact) {
  const url = `https://wa.me/${contact.waNumber}?text=${encodeURIComponent(contact.defaultMessage)}`
  window.open(url, '_blank')
}

function handleClickOutside(event) {
  if (widgetRef.value && !widgetRef.value.contains(event.target)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="fixed bottom-6 right-6 z-50 text-left" ref="widgetRef">
    <!-- 1. POPUP CARD: 3 WhatsApp Contacts -->
    <div
      v-if="isOpen"
      class="absolute bottom-16 right-0 w-80 sm:w-96 bg-white rounded-2xl border border-slate-200 shadow-2xl overflow-hidden animate-in fade-in slide-in-from-bottom-3 duration-200"
    >
      <!-- Header -->
      <div class="p-4 bg-gradient-to-r from-emerald-600 via-emerald-700 to-teal-800 text-white flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center shrink-0">
            <!-- Authentic WhatsApp SVG Icon -->
            <svg class="w-4 h-4 fill-white" viewBox="0 0 24 24">
              <path d="M12.031 6.172c-3.181 0-5.767 2.586-5.768 5.766-.001 1.298.38 2.27 1.019 3.287l-.582 2.128 2.182-.573c.978.58 1.911.928 3.145.929 3.178 0 5.767-2.587 5.768-5.766.001-3.187-2.575-5.77-5.764-5.771zm3.392 8.244c-.144.405-.837.774-1.17.824-.299.045-.677.063-1.092-.069-.252-.08-.575-.187-.988-.365-1.739-.751-2.874-2.502-2.961-2.617-.087-.116-.708-.94-.708-1.793s.448-1.273.607-1.446c.159-.173.346-.217.462-.217l.332.006c.106.005.249-.04.39.298.144.347.491 1.2.534 1.287.043.087.072.188.014.304-.058.116-.087.188-.173.289l-.26.304c-.087.086-.177.18-.076.354.101.174.449.741.964 1.201.662.591 1.221.774 1.394.86s.275.072.376-.043c.101-.116.433-.506.549-.68.116-.173.231-.145.39-.087s1.011.477 1.184.564.289.13.332.202c.045.072.045.419-.1.824zm-3.423-14.416c-6.627 0-12 5.373-12 12 0 2.152.569 4.17 1.564 5.914l-1.664 6.086 6.257-1.64c1.706.924 3.655 1.455 5.726 1.455 6.627 0 12-5.373 12-12 0-6.627-5.373-12-12-12z"/>
            </svg>
          </div>
          <div>
            <h4 class="text-sm font-bold leading-tight">WhatsApp Support Probus</h4>
            <p class="text-[11px] text-emerald-100/90 mt-0.5">3 Jalur Bantuan Langsung</p>
          </div>
        </div>

        <button
          @click="closeWidget"
          class="w-7 h-7 rounded-full bg-white/10 hover:bg-white/20 text-white flex items-center justify-center transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Description Header -->
      <div class="px-4 py-2.5 bg-emerald-50/70 border-b border-emerald-100 text-[11px] text-emerald-900 flex items-center gap-1.5">
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shrink-0"></span>
        <span>Pilih nomor WhatsApp tujuan Anda di bawah ini:</span>
      </div>

      <!-- List of 3 Contacts -->
      <div class="p-2 space-y-1 divide-y divide-slate-100">
        <div
          v-for="contact in contacts"
          :key="contact.id"
          @click="openWhatsApp(contact)"
          class="p-3 rounded-xl hover:bg-slate-50 cursor-pointer transition-all duration-150 flex items-start gap-3 group border border-transparent hover:border-slate-200/70"
        >
          <!-- Icon -->
          <div class="w-9 h-9 rounded-xl bg-slate-100 text-slate-700 flex items-center justify-center shrink-0 mt-0.5 group-hover:bg-emerald-50 group-hover:text-emerald-600 transition-colors">
            <component :is="contact.icon" class="w-4 h-4 stroke-[2]" />
          </div>

          <!-- Info -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between gap-1">
              <h5 class="text-xs sm:text-sm font-bold text-slate-900 group-hover:text-emerald-700 transition-colors">
                {{ contact.name }}
              </h5>
              <span :class="['text-[10px] font-bold px-2 py-0.5 rounded border', contact.badgeClass]">
                {{ contact.badge }}
              </span>
            </div>
            <p class="text-[11px] text-slate-500 mt-0.5 line-clamp-1">
              {{ contact.role }}
            </p>
            <div class="mt-1.5 flex items-center gap-1.5 text-xs font-semibold text-emerald-700">
              <span>{{ contact.numberDisplay }}</span>
              <ExternalLink class="w-3 h-3 text-emerald-500 group-hover:translate-x-0.5 transition-transform" />
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="p-3 bg-slate-50 border-t border-slate-100 text-[11px] text-slate-400 text-center">
        PT Probus System Indonesia • Layanan Bantuan Cepat
      </div>
    </div>

    <!-- 2. FLOATING TRIGGER BUTTON (Pojok Kanan Bawah) -->
    <div class="relative flex items-center">
      <!-- Tooltip Pill (Hover) -->
      <div
        v-if="!isOpen"
        class="hidden sm:flex items-center mr-3 px-3 py-1.5 rounded-full bg-slate-900 text-white text-xs font-semibold shadow-lg whitespace-nowrap opacity-90 hover:opacity-100 transition-opacity"
      >
        <span>Chat WhatsApp Support</span>
      </div>

      <button
        @click.stop="toggleWidget"
        type="button"
        :aria-label="isOpen ? 'Tutup WhatsApp' : 'Buka WhatsApp Support'"
        :class="[
          'w-14 h-14 rounded-full flex items-center justify-center text-white shadow-xl hover:shadow-2xl transition-all duration-300 cursor-pointer',
          isOpen
            ? 'bg-slate-900 hover:bg-slate-800 rotate-90 scale-95'
            : 'bg-emerald-500 hover:bg-emerald-600 hover:scale-110 ring-4 ring-emerald-500/20'
        ]"
      >
        <!-- If Open: Show Close Icon -->
        <X v-if="isOpen" class="w-6 h-6 stroke-[2.5]" />

        <!-- If Closed: Show Official WhatsApp Logo -->
        <div v-else class="relative flex items-center justify-center">
          <svg class="w-7 h-7 fill-white" viewBox="0 0 24 24">
            <path d="M12.031 6.172c-3.181 0-5.767 2.586-5.768 5.766-.001 1.298.38 2.27 1.019 3.287l-.582 2.128 2.182-.573c.978.58 1.911.928 3.145.929 3.178 0 5.767-2.587 5.768-5.766.001-3.187-2.575-5.77-5.764-5.771zm3.392 8.244c-.144.405-.837.774-1.17.824-.299.045-.677.063-1.092-.069-.252-.08-.575-.187-.988-.365-1.739-.751-2.874-2.502-2.961-2.617-.087-.116-.708-.94-.708-1.793s.448-1.273.607-1.446c.159-.173.346-.217.462-.217l.332.006c.106.005.249-.04.39.298.144.347.491 1.2.534 1.287.043.087.072.188.014.304-.058.116-.087.188-.173.289l-.26.304c-.087.086-.177.18-.076.354.101.174.449.741.964 1.201.662.591 1.221.774 1.394.86s.275.072.376-.043c.101-.116.433-.506.549-.68.116-.173.231-.145.39-.087s1.011.477 1.184.564.289.13.332.202c.045.072.045.419-.1.824zm-3.423-14.416c-6.627 0-12 5.373-12 12 0 2.152.569 4.17 1.564 5.914l-1.664 6.086 6.257-1.64c1.706.924 3.655 1.455 5.726 1.455 6.627 0 12-5.373 12-12 0-6.627-5.373-12-12-12z"/>
          </svg>

          <!-- Notification Dot -->
          <span class="absolute -top-1 -right-1 w-3.5 h-3.5 bg-red-500 rounded-full border-2 border-white flex items-center justify-center text-[9px] font-black">
            3
          </span>
        </div>
      </button>
    </div>
  </div>
</template>

