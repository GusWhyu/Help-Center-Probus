<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { MessageSquare, ChevronDown, ExternalLink, Headphones, ShoppingCart, GraduationCap } from 'lucide-vue-next'

const isOpen = ref(false)
const dropdownRef = ref(null)

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

function toggleDropdown() {
  isOpen.value = !isOpen.value
}

function openWhatsApp(contact) {
  const url = `https://wa.me/${contact.waNumber}?text=${encodeURIComponent(contact.defaultMessage)}`
  window.open(url, '_blank')
  isOpen.value = false
}

function handleClickOutside(event) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
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
  <div class="relative inline-block text-left" ref="dropdownRef">
    <!-- Trigger Button in Utility Top Bar -->
    <div class="flex items-center gap-2">
      <span class="text-slate-500 hidden sm:inline">Butuh bantuan cepat?</span>
      <button
        @click.stop="toggleDropdown"
        type="button"
        class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-50 hover:bg-emerald-100/80 text-emerald-800 font-semibold text-xs border border-emerald-200/80 transition-all cursor-pointer shadow-2xs group"
      >
        <MessageSquare class="w-3.5 h-3.5 text-emerald-600 fill-emerald-600 group-hover:scale-110 transition-transform" />
        <span>Hubungi WhatsApp Support</span>
        <ChevronDown
          class="w-3 h-3 text-emerald-700 transition-transform duration-200"
          :class="{ 'rotate-180': isOpen }"
        />
      </button>
    </div>

    <!-- Dropdown Menu Popover -->
    <div
      v-if="isOpen"
      class="absolute right-0 top-full mt-2 w-80 sm:w-96 bg-white rounded-2xl border border-slate-200 shadow-2xl z-50 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-150 text-left"
    >
      <!-- Header of Dropdown -->
      <div class="p-4 bg-gradient-to-r from-emerald-600 to-teal-700 text-white">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <div class="w-7 h-7 rounded-lg bg-white/20 backdrop-blur-md flex items-center justify-center">
              <MessageSquare class="w-4 h-4 text-white fill-white" />
            </div>
            <h4 class="text-sm font-bold tracking-tight">WhatsApp Support Probus</h4>
          </div>
          <span class="text-[10px] bg-white/20 text-white font-bold px-2 py-0.5 rounded-full">
            3 Kontak
          </span>
        </div>
        <p class="text-xs text-emerald-100/90 mt-1">
          Pilih kontak di bawah untuk langsung terhubung via WhatsApp:
        </p>
      </div>

      <!-- 3 WhatsApp Contacts List -->
      <div class="p-2 divide-y divide-slate-100">
        <div
          v-for="contact in contacts"
          :key="contact.id"
          @click="openWhatsApp(contact)"
          class="p-3 rounded-xl hover:bg-slate-50 cursor-pointer transition-all duration-150 flex items-start gap-3 group border border-transparent hover:border-slate-200/70"
        >
          <!-- Icon Box -->
          <div class="w-9 h-9 rounded-xl bg-slate-100 text-slate-700 flex items-center justify-center shrink-0 mt-0.5 group-hover:bg-emerald-50 group-hover:text-emerald-600 transition-colors">
            <component :is="contact.icon" class="w-4 h-4 stroke-[2]" />
          </div>

          <!-- Contact Details -->
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

      <!-- Footer Note -->
      <div class="p-3 bg-slate-50 border-t border-slate-100 text-[11px] text-slate-500 text-center flex items-center justify-center gap-1.5">
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
        <span>Tim support siap melayani Anda pada jam kerja operasional.</span>
      </div>
    </div>
  </div>
</template>
