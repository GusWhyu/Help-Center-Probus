<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  Phone,
  Mail,
  MapPin,
  ExternalLink,
  Send,
  MessageSquare,
  CheckCircle2,
  Building2,
  Clock,
  Globe,
  ArrowRight,
  Menu,
  X
} from 'lucide-vue-next'

const router = useRouter()
const mobileMenuOpen = ref(false)

// Form state
const form = ref({
  name: '',
  email: '',
  phone: '',
  system: 'PRESTO Resto POS',
  topic: 'Konsultasi Sistem & Implementasi',
  message: '',
})

const isSubmitted = ref(false)
const isSubmitting = ref(false)

const systemOptions = [
  'PRESTO Resto POS',
  'PHIS Hotel PMS',
  'PPOS Point of Sales',
  'Probus Self Order',
  'PR-PO Online Logistik',
  'Order Taker Mobile',
  'Vrestropical POS',
  'PPE 2.0 Management',
  'Barong POS',
  'Konsultasi Umum Sistem',
]

const topicOptions = [
  'Konsultasi Sistem & Implementasi',
  'Permintaan Demo Produk (Custom Demo)',
  'Jadwal Training Modul & Operasional',
  'Bantuan Kendala Teknis Kasir / Hotel',
  'Kerjasama Bisnis & Kemitraan',
]

function handleSubmit() {
  if (!form.value.name || !form.value.phone) {
    alert('Mohon isi nama lengkap dan nomor telepon Anda.')
    return
  }
  isSubmitting.value = true
  setTimeout(() => {
    isSubmitting.value = false
    isSubmitted.value = true
  }, 600)
}

function resetForm() {
  form.value = {
    name: '',
    email: '',
    phone: '',
    system: 'PRESTO Resto POS',
    topic: 'Konsultasi Sistem & Implementasi',
    message: '',
  }
  isSubmitted.value = false
}
</script>

<template>
  <div class="min-h-screen bg-slate-50/40 text-slate-800 font-sans flex flex-col antialiased">
    <!-- Main Navigation Bar -->
    <nav class="border-b border-slate-200/90 bg-white py-3.5 px-4 sm:px-6 lg:px-12 xl:px-16 sticky top-0 z-40 shadow-xs">
      <div class="w-full flex items-center justify-between">
        <!-- Logo & Title -->
        <div class="flex items-center gap-3">
          <div class="flex items-center cursor-pointer" @click="router.push('/')">
            <img
              src="/probus-logo.png"
              alt="Probus System"
              class="h-8 sm:h-9 w-auto object-contain hover:opacity-90 transition-opacity"
            />
          </div>
          <span class="text-slate-300 font-light hidden sm:inline">|</span>
          <span class="text-slate-500 text-sm font-medium hidden sm:inline">Pusat Bantuan Publik</span>
        </div>

        <!-- Desktop Navigation Links (Hubungi Implementator Active) -->
        <div class="hidden md:flex items-center gap-8 text-sm font-medium">
          <router-link to="/" class="text-slate-600 hover:text-slate-950 transition-colors">
            Panduan Sistem
          </router-link>
          <router-link to="/solutions" class="text-slate-600 hover:text-slate-950 transition-colors">
            Bank Solusi
          </router-link>
          <router-link to="/contact" class="text-sky-600 font-semibold hover:text-sky-700 transition-colors relative py-1">
            Hubungi Implementator
            <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-sky-600 rounded-full"></span>
          </router-link>
        </div>

        <!-- Mobile Hamburger Button -->
        <div class="flex md:hidden items-center">
          <button
            @click="mobileMenuOpen = !mobileMenuOpen"
            type="button"
            class="p-2 rounded-xl text-slate-600 hover:text-slate-900 hover:bg-slate-100 focus:outline-none transition-colors"
            aria-label="Toggle Menu"
          >
            <X v-if="mobileMenuOpen" class="w-5 h-5" />
            <Menu v-else class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Mobile Dropdown Drawer -->
      <div
        v-if="mobileMenuOpen"
        class="md:hidden border-t border-slate-100 mt-3 pt-3 pb-2 space-y-1 animate-in fade-in slide-in-from-top-2 duration-150"
      >
        <router-link
          to="/"
          @click="mobileMenuOpen = false"
          class="block px-3 py-2 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-slate-900"
        >
          Panduan Sistem
        </router-link>
        <router-link
          to="/solutions"
          @click="mobileMenuOpen = false"
          class="block px-3 py-2 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-slate-900"
        >
          Bank Solusi
        </router-link>
        <router-link
          to="/contact"
          @click="mobileMenuOpen = false"
          class="block px-3 py-2 rounded-lg text-sm font-semibold bg-sky-50 text-sky-700"
        >
          Hubungi Implementator
        </router-link>
        <div class="pt-2 border-t border-slate-100 px-3">
          <a
            href="tel:+62361256500"
            class="inline-flex items-center gap-2 text-xs font-bold text-slate-700 hover:text-sky-600"
          >
            <Phone class="w-3.5 h-3.5 text-slate-500" />
            <span>Phone: +62 361 256 500</span>
          </a>
        </div>
      </div>
    </nav>

    <!-- 3. Header Hero Section (Matches Screenshot media_1794540993219_f5b98a39.jpg) -->
    <section class="pt-10 pb-8 sm:pt-14 sm:pb-10 px-4 sm:px-6 lg:px-12 xl:px-16 text-center bg-gradient-to-b from-white via-sky-50/20 to-transparent">
      <div class="w-full max-w-4xl mx-auto">
        <h1 class="text-2xl sm:text-3xl lg:text-4xl xl:text-5xl font-extrabold text-slate-900 tracking-tight leading-tight">
          Contact Our Team
        </h1>
        <p class="mt-3 sm:mt-4 text-slate-600 text-sm sm:text-base lg:text-lg leading-relaxed max-w-2xl mx-auto">
          Have questions about our hotel PMS, restaurant POS, or retail systems? Reach out to our consultants in Bali for advice or a custom demo.
        </p>

        <!-- Company Profile Button Requested by User -->
        <div class="mt-6 sm:mt-7 flex flex-wrap items-center justify-center gap-3">
          <a
            href="https://probussystem.com/"
            target="_blank"
            class="inline-flex items-center gap-2 px-5 sm:px-6 py-2.5 sm:py-3 rounded-2xl bg-slate-900 hover:bg-slate-800 text-white font-bold text-xs sm:text-sm shadow-md hover:shadow-lg transition-all group"
          >
            <Globe class="w-4 h-4 text-sky-400" />
            <span>Kunjungi Website Company Profile Probus</span>
            <ExternalLink class="w-4 h-4 text-slate-400 group-hover:translate-x-0.5 transition-transform" />
          </a>
        </div>
      </div>
    </section>

    <!-- 4. Main Contact Grid (Two Columns: Send Us a Message & Office/Contact Cards) -->
    <main class="flex-1 w-full px-4 sm:px-6 lg:px-12 xl:px-16 pb-20">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        
        <!-- LEFT COLUMN: "Send Us a Message" Form (col-span-7) -->
        <div class="lg:col-span-7 bg-white rounded-3xl border border-slate-200/90 p-7 sm:p-9 shadow-xs">
          <div class="flex items-center justify-between border-b border-slate-100 pb-5 mb-6">
            <div>
              <h2 class="text-xl sm:text-2xl font-black text-slate-900 tracking-tight">
                Send Us a Message
              </h2>
              <p class="text-xs sm:text-sm text-slate-500 mt-1">
                Kirim pesan langsung ke tim konsultan implementator Probus System.
              </p>
            </div>
            <div class="w-10 h-10 rounded-2xl bg-sky-50 text-sky-600 flex items-center justify-center shrink-0">
              <MessageSquare class="w-5 h-5" />
            </div>
          </div>

          <!-- Success Alert -->
          <div
            v-if="isSubmitted"
            class="p-6 rounded-2xl bg-emerald-50 border border-emerald-200 text-center space-y-3 animate-in fade-in duration-200"
          >
            <div class="w-12 h-12 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center mx-auto">
              <CheckCircle2 class="w-6 h-6" />
            </div>
            <h3 class="text-base font-bold text-emerald-900">Pesan Berhasil Terkirim!</h3>
            <p class="text-xs sm:text-sm text-emerald-800 leading-relaxed max-w-md mx-auto">
              Terima kasih, <strong>{{ form.name }}</strong>. Tim konsultan Probus System telah menerima informasi Anda dan akan segera menghubungi nomor <strong>{{ form.phone }}</strong>.
            </p>
            <div class="pt-3 flex flex-wrap items-center justify-center gap-3">
              <button
                @click="resetForm"
                class="px-5 py-2.5 rounded-xl bg-emerald-600 text-white text-xs font-semibold hover:bg-emerald-700 transition-colors cursor-pointer shadow-xs"
              >
                Kirim Pesan Lainnya
              </button>
            </div>
          </div>

          <!-- Contact Form -->
          <form v-else @submit.prevent="handleSubmit" class="space-y-5">
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <!-- Name -->
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1.5">
                  Nama Lengkap <span class="text-rose-500">*</span>
                </label>
                <input
                  v-model="form.name"
                  type="text"
                  required
                  placeholder="Contoh: Bpk. Made Wijaya"
                  class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/10 transition-all"
                />
              </div>

              <!-- Phone -->
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1.5">
                  Nomor Telepon <span class="text-rose-500">*</span>
                </label>
                <input
                  v-model="form.phone"
                  type="tel"
                  required
                  placeholder="Contoh: 0812-3456-7890"
                  class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/10 transition-all"
                />
              </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <!-- Email -->
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1.5">
                  Alamat Email
                </label>
                <input
                  v-model="form.email"
                  type="email"
                  placeholder="nama@hotelresto.com"
                  class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/10 transition-all"
                />
              </div>

              <!-- Sistem Yang Ditanyakan -->
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1.5">
                  Sistem / Produk Terkait
                </label>
                <select
                  v-model="form.system"
                  class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm text-slate-800 bg-white focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/10 transition-all"
                >
                  <option v-for="sys in systemOptions" :key="sys" :value="sys">
                    {{ sys }}
                  </option>
                </select>
              </div>
            </div>

            <!-- Topik -->
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1.5">
                Topik Kebutuhan
              </label>
              <select
                v-model="form.topic"
                class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm text-slate-800 bg-white focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/10 transition-all"
              >
                <option v-for="top in topicOptions" :key="top" :value="top">
                  {{ top }}
                </option>
              </select>
            </div>

            <!-- Message -->
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1.5">
                Pesan / Keterangan Kendala
              </label>
              <textarea
                v-model="form.message"
                rows="4"
                placeholder="Tuliskan pertanyaan, modul yang ingin dipelajari, atau deskripsi kendala yang dihadapi..."
                class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/10 transition-all"
              ></textarea>
            </div>

            <!-- Buttons -->
            <div class="pt-2 flex items-center">
              <button
                type="submit"
                :disabled="isSubmitting"
                class="w-full sm:w-auto px-6 h-11 sm:h-12 rounded-xl bg-sky-600 hover:bg-sky-700 text-white font-bold text-sm transition-all duration-200 inline-flex items-center justify-center gap-2 shadow-xs cursor-pointer disabled:opacity-50 whitespace-nowrap"
              >
                <Send class="w-4 h-4 shrink-0" />
                <span class="whitespace-nowrap">{{ isSubmitting ? 'Mengirim...' : 'Kirim Pesan' }}</span>
              </button>
            </div>
          </form>
        </div>

        <!-- RIGHT COLUMN: Location, Direct Contact, & Map Preview (col-span-5) -->
        <div class="lg:col-span-5 space-y-6">
          
          <!-- Top 2 Info Cards: Office Location & Direct Contact -->
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-4">
            
            <!-- Card 1: Office Location -->
            <div class="bg-white rounded-2xl border border-slate-200/90 p-5 shadow-xs flex flex-col justify-between">
              <div>
                <span class="text-[11px] font-bold text-sky-600 uppercase tracking-wider block mb-1">
                  OFFICE LOCATION
                </span>
                <h3 class="text-sm sm:text-base font-extrabold text-slate-900 leading-snug">
                  Pertokoan Sudirman Agung
                </h3>
                <p class="text-xs text-slate-600 mt-2 leading-relaxed">
                  Pertokoan Sudirman Agung Blok C No. 51, Jl. PB Sudirman, Dauh Puri Klod, Denpasar, Bali 80113
                </p>
              </div>
            </div>

            <!-- Card 2: Direct Contact -->
            <div class="bg-white rounded-2xl border border-slate-200/90 p-5 shadow-xs flex flex-col justify-between">
              <div>
                <span class="text-[11px] font-bold text-sky-600 uppercase tracking-wider block mb-1">
                  DIRECT CONTACT
                </span>
                <h3 class="text-sm sm:text-base font-extrabold text-slate-900 leading-snug">
                  Phone & Email
                </h3>
                <div class="text-xs text-slate-600 mt-2 space-y-1 leading-relaxed">
                  <p>
                    <span class="text-slate-400">Phone: </span>
                    <a href="tel:+62361256500" class="font-semibold text-slate-800 hover:text-sky-600 transition-colors">
                      +62 361 256 500
                    </a>
                  </p>
                  <p class="break-all">
                    <span class="text-slate-400">Email: </span>
                    <a href="mailto:info@probussystem.com" class="font-semibold text-slate-800 hover:text-sky-600 transition-colors">
                      info@probussystem.com
                    </a>
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Bottom: Interactive Google Maps Card (Matching Screenshot) -->
          <div class="bg-white rounded-2xl border border-slate-200/90 overflow-hidden shadow-xs">
            <div class="p-3.5 bg-slate-50 border-b border-slate-100 flex items-center justify-between">
              <div class="flex items-center gap-2 text-xs font-bold text-slate-700">
                <MapPin class="w-4 h-4 text-sky-600" />
                <span>Peta Lokasi Kantor Bali</span>
              </div>
              <a
                href="https://maps.google.com/?q=Pertokoan+Sudirman+Agung+Blok+C+No.+51,+Jl.+PB+Sudirman,+Denpasar,+Bali"
                target="_blank"
                class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-sky-50 hover:bg-sky-100 text-sky-700 font-bold text-xs border border-sky-100 transition-colors"
              >
                <span>Open in Maps</span>
                <ExternalLink class="w-3.5 h-3.5" />
              </a>
            </div>

            <!-- Map Iframe Embed -->
            <div class="relative w-full h-64 sm:h-72 bg-slate-100">
              <iframe
                title="Peta Lokasi Probus System Bali"
                width="100%"
                height="100%"
                style="border:0;"
                loading="lazy"
                allowfullscreen
                referrerpolicy="no-referrer-when-downgrade"
                src="https://maps.google.com/maps?q=Pertokoan+Sudirman+Agung+Blok+C+No.+51,+Jl.+PB+Sudirman,+Denpasar,+Bali&t=&z=16&ie=UTF8&iwloc=&output=embed"
              ></iframe>
            </div>

            <div class="p-3 bg-white text-[11px] text-slate-500 border-t border-slate-100 flex items-center justify-between">
              <span>Denpasar, Bali 80113</span>
              <span class="text-emerald-600 font-semibold flex items-center gap-1">
                <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                <span>Kantor Pusat Operasional</span>
              </span>
            </div>
          </div>

          <!-- Extra Banner: Company Profile Promotion -->
          <div class="p-5 rounded-2xl bg-gradient-to-r from-sky-900 to-slate-900 text-white flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 shadow-sm">
            <div>
              <h4 class="text-sm font-bold tracking-tight">Profil Lengkap Probus System</h4>
              <p class="text-xs text-sky-200/80 mt-0.5">
                Kunjungi official website untuk melihat ekosistem produk & klien kami.
              </p>
            </div>
            <a
              href="https://probussystem.com/"
              target="_blank"
              class="shrink-0 px-4 py-2 rounded-xl bg-white text-slate-900 hover:bg-sky-50 font-bold text-xs transition-colors inline-flex items-center gap-1.5"
            >
              <span>probussystem.com</span>
              <ExternalLink class="w-3.5 h-3.5 text-sky-600" />
            </a>
          </div>

        </div>

      </div>
    </main>

    <!-- 5. Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 px-4 sm:px-6 lg:px-12 xl:px-16 text-xs sm:text-sm text-slate-500 mt-auto">
      <div class="w-full flex flex-col sm:flex-row items-center justify-between gap-3 text-center sm:text-left">
        <p class="font-normal">© 2026 CV Mitra Solusi Mandiri</p>
      </div>
    </footer>
  </div>
</template>

