<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Search,
  ArrowLeft,
  ChevronRight,
  AlertCircle,
  Clock,
  ChevronDown,
  X,
  ThumbsUp,
  CheckCircle2,
  Filter,
  Lightbulb,
  Printer,
  Wifi,
  MonitorCheck,
  Building2,
  Zap,
  HelpCircle,
  Loader2
} from 'lucide-vue-next'
import apiClient from '../api/client'
import WhatsAppSupportDropdown from '../components/WhatsAppSupportDropdown.vue'

const router = useRouter()
const searchQuery = ref('')
const selectedSystemFilter = ref('all')
const isLoading = ref(false)

const systemNameMap = {
  1: 'PHIS',
  2: 'PRESTO',
  3: 'PPOS',
  4: 'SELF ORDER',
  5: 'PR-PO ONLINE',
  6: 'ORDER TAKER',
  7: 'VRESTROPICAL',
  8: 'PPE',
  9: 'BARONG POS',
}

function formatArticleContent(html) {
  if (!html) return ''
  return html.replace(/http:\/\/(127\.0\.0\.1|localhost):8000\/uploads\//g, '/uploads/')
}

function cleanSummary(text) {
  if (!text) return ''
  return text.replace(/&nbsp;/g, ' ').replace(/<[^>]*>?/gm, '').trim()
}

// Articles list from database
const articles = ref([])

// Load articles from database
async function loadArticles() {
  isLoading.value = true
  try {
    const res = await apiClient.get('/articles', { params: { limit: 100 } })
    if (res.data && res.data.data && Array.isArray(res.data.data)) {
      articles.value = res.data.data.map(item => {
        const primaryMenu = item.menu_systems && item.menu_systems.length > 0 ? item.menu_systems[0] : null
        const katId = primaryMenu ? primaryMenu.kategori_id : 2
        const sysCode = systemNameMap[katId] || 'PRESTO'
        const menuName = primaryMenu ? primaryMenu.name : 'Umum'

        // Estimate read time
        const words = (item.content || '').split(/\s+/).length
        const minutes = Math.max(2, Math.min(10, Math.ceil(words / 150)))

        return {
          id: item.id,
          title: item.title,
          system: sysCode,
          kategoriId: katId,
          categoryName: menuName,
          menuName: menuName,
          status: item.status,
          visibility: item.visibility,
          readTime: `${minutes} menit baca`,
          summary: cleanSummary(item.summary) || 'Panduan pemecahan kendala teknis dan operasional sistem Probus.',
          content: formatArticleContent(item.content),
          raw: item,
        }
      })
    }
  } catch (err) {
    console.error('Failed to load articles from database:', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadArticles()
})

// Filter options
const systemOptions = [
  { value: 'all', label: 'Semua Sistem' },
  { value: 'PRESTO', label: 'PRESTO Resto POS' },
  { value: 'PHIS', label: 'Probus Hotel PMS' },
  { value: 'PPOS', label: 'PPOS Point of Sales' },
  { value: 'PPE', label: 'PPE 2.0 Management' },
]

// Filtered articles
const filteredArticles = computed(() => {
  return articles.value.filter(art => {
    const matchSystem = selectedSystemFilter.value === 'all' ||
      art.system === selectedSystemFilter.value ||
      art.kategoriId === Number(selectedSystemFilter.value)
    const q = searchQuery.value.trim().toLowerCase()
    const matchQuery = !q ||
      (art.title && art.title.toLowerCase().includes(q)) ||
      (art.summary && art.summary.toLowerCase().includes(q)) ||
      (art.categoryName && art.categoryName.toLowerCase().includes(q)) ||
      (art.system && art.system.toLowerCase().includes(q))
    return matchSystem && matchQuery
  })
})

// Reader Modal
const readerModalOpen = ref(false)
const activeArticle = ref(null)
const feedbackGiven = ref(false)

function openArticle(art) {
  activeArticle.value = art
  feedbackGiven.value = false
  readerModalOpen.value = true
}

function closeReader() {
  readerModalOpen.value = false
  activeArticle.value = null
}

function giveFeedback() {
  feedbackGiven.value = true
}

function goBackHome() {
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-slate-50/40 text-slate-800 font-sans flex flex-col antialiased">
    <!-- 1. Utility Top Bar -->
    <header class="border-b border-slate-200/70 bg-white text-xs text-slate-500 py-2.5 px-6 sm:px-12 lg:px-20 xl:px-28">
      <div class="max-w-[1600px] mx-auto flex flex-col sm:flex-row items-center justify-between gap-2">
        <div class="flex items-center gap-2 cursor-pointer hover:text-slate-800 transition-colors">
          <span class="text-base">🇮🇩</span>
          <span class="font-medium text-slate-600">Bahasa Indonesia</span>
          <ChevronDown class="w-3.5 h-3.5 text-slate-400" />
        </div>
        <WhatsAppSupportDropdown />
      </div>
    </header>

    <!-- 2. Main Navigation Bar -->
    <nav class="border-b border-slate-200/80 bg-white/95 backdrop-blur-sm py-3.5 px-6 sm:px-12 lg:px-20 xl:px-28 sticky top-0 z-40 shadow-xs">
      <div class="max-w-[1600px] mx-auto flex items-center justify-between">
        <div class="flex items-center gap-3.5">
          <div class="flex items-center gap-2.5 cursor-pointer" @click="goBackHome">
            <div class="w-8 h-8 rounded-lg bg-slate-950 flex items-center justify-center text-white font-extrabold text-base tracking-wider shadow-sm">
              P
            </div>
            <span class="font-black text-lg tracking-wider text-slate-950">PROBUS</span>
          </div>
          <span class="text-slate-300 font-light">|</span>
          <span class="text-slate-500 text-sm font-medium">Bank Solusi Publik</span>

          <button
            @click="goBackHome"
            class="ml-2 inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs transition-colors"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>Kembali ke Beranda</span>
          </button>
        </div>

        <div class="hidden lg:flex items-center gap-8 text-sm font-medium">
          <router-link to="/" class="text-slate-600 hover:text-slate-950 transition-colors">
            Panduan Sistem
          </router-link>
          <router-link to="/solutions" class="text-sky-600 font-semibold hover:text-sky-700 transition-colors relative py-1">
            Bank Solusi
            <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-sky-600 rounded-full"></span>
          </router-link>
          <a href="https://wa.me/6281910875514" target="_blank" class="text-slate-600 hover:text-slate-950 transition-colors">
            Hubungi Implementator
          </a>
        </div>
      </div>
    </nav>

    <!-- 3. Hero Section (Bank Solusi) -->
    <section class="pt-14 pb-12 px-6 sm:px-12 lg:px-20 xl:px-28 text-center bg-gradient-to-b from-white via-amber-50/20 to-transparent">
      <div class="max-w-4xl mx-auto">
        <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-md bg-amber-50 text-amber-800 text-xs font-semibold border border-amber-200/60 mb-3.5 shadow-2xs">
          <AlertCircle class="w-3.5 h-3.5" />
          <span>Knowledge Base & Troubleshooting</span>
        </div>

        <h1 class="text-3xl sm:text-4xl font-extrabold text-slate-900 tracking-tight">
          Bank Solusi Kendala & Troubleshooting
        </h1>
        <p class="mt-3 text-slate-500 text-sm sm:text-base leading-relaxed max-w-2xl mx-auto">
          Temukan solusi cepat untuk masalah printer kasir, cash drawer, jaringan LAN, dan error operasional sistem Probus.
        </p>

        <!-- Search Bar -->
        <div class="mt-7 max-w-3xl mx-auto relative shadow-2xs">
          <Search class="absolute left-4 w-4 h-4 text-slate-400 pointer-events-none top-4" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari solusi kendala (contoh: printer, cash drawer, paper jam, IP address)..."
            class="w-full pl-11 pr-10 py-3.5 rounded-xl text-sm bg-white border border-slate-300 text-slate-900 placeholder-slate-400 focus:outline-none focus:border-amber-500 focus:ring-4 focus:ring-amber-100 transition-all shadow-xs"
          />
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            class="absolute right-3.5 top-3.5 text-slate-400 hover:text-slate-600"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Filter Pills by System -->
        <div class="mt-5 flex flex-wrap items-center justify-center gap-2">
          <button
            v-for="opt in systemOptions"
            :key="opt.value"
            @click="selectedSystemFilter = opt.value"
            :class="[
              'px-3.5 py-1.5 rounded-lg text-xs font-semibold transition-all cursor-pointer',
              selectedSystemFilter === opt.value
                ? 'bg-slate-900 text-white shadow-xs'
                : 'bg-white text-slate-600 border border-slate-200 hover:bg-slate-100'
            ]"
          >
            {{ opt.label }}
          </button>
        </div>
      </div>
    </section>

    <!-- 4. Main Articles List -->
    <main class="flex-1 max-w-[1600px] w-full mx-auto px-6 sm:px-12 lg:px-20 xl:px-28 pb-20">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h2 class="text-xl font-bold text-slate-900">
            Daftar Artikel Solusi Masalah
          </h2>
          <p class="text-xs sm:text-sm text-slate-500 mt-0.5">
            Menampilkan {{ filteredArticles.length }} artikel pemecahan masalah teknis dari basis data
          </p>
        </div>

        <div v-if="isLoading" class="flex items-center gap-2 text-xs text-slate-400">
          <Loader2 class="w-4 h-4 animate-spin text-amber-500" />
          <span>Memuat artikel...</span>
        </div>
      </div>

      <!-- Articles Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="art in filteredArticles"
          :key="art.id"
          @click="openArticle(art)"
          class="bg-white rounded-2xl p-6 border border-slate-200 hover:border-amber-400 hover:shadow-md transition-all cursor-pointer flex flex-col justify-between group"
        >
          <div>
            <div class="flex items-center justify-between gap-2 mb-3.5">
              <span class="text-xs font-bold px-2.5 py-0.5 rounded-md bg-amber-50 text-amber-700 border border-amber-200/60">
                {{ art.system }} • {{ art.categoryName }}
              </span>
              <span class="text-xs text-slate-400 flex items-center gap-1">
                <Clock class="w-3 h-3" /> {{ art.readTime }}
              </span>
            </div>

            <h3 class="font-bold text-slate-900 text-base group-hover:text-amber-600 transition-colors leading-snug">
              {{ art.title }}
            </h3>

            <p class="text-xs sm:text-sm text-slate-500 mt-2.5 line-clamp-3 leading-relaxed">
              {{ art.summary }}
            </p>
          </div>

          <div class="pt-5 mt-5 border-t border-slate-100 flex items-center justify-between text-xs">
            <span class="text-slate-400">Terverifikasi Tim Teknis</span>
            <span class="text-amber-600 font-semibold flex items-center gap-1 group-hover:translate-x-0.5 transition-transform">
              <span>Buka Solusi</span>
              <ChevronRight class="w-3.5 h-3.5" />
            </span>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!isLoading && filteredArticles.length === 0" class="text-center py-16 bg-white rounded-2xl border border-slate-200">
        <AlertCircle class="w-12 h-12 text-slate-300 mx-auto mb-3" />
        <h3 class="text-base font-bold text-slate-800">Tidak ada artikel yang sesuai</h3>
        <p class="text-xs sm:text-sm text-slate-500 mt-1 max-w-md mx-auto">
          Coba ganti kata kunci pencarian atau hubungi teknisi implementator melalui WhatsApp.
        </p>
        <button
          @click="searchQuery = ''; selectedSystemFilter = 'all'"
          class="mt-4 px-4 py-2 rounded-lg bg-slate-900 text-white text-xs font-semibold hover:bg-slate-800 transition-colors"
        >
          Reset Filter Pencarian
        </button>
      </div>
    </main>

    <!-- 5. Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 px-6 sm:px-12 lg:px-20 xl:px-28 text-xs sm:text-sm text-slate-500">
      <div class="max-w-[1600px] mx-auto flex flex-col sm:flex-row items-center justify-between gap-3">
        <p>© 2026 PT Probus System Indonesia • Bank Solusi Resmi</p>
        <p>
          Masih terkendala? Hubungi WhatsApp Support:
          <a href="https://wa.me/6281910875514" target="_blank" class="font-semibold text-slate-700 hover:text-sky-600 transition-colors">
            +62 819-1087-5514
          </a>
        </p>
      </div>
    </footer>

    <!-- 6. Interactive Reader Modal -->
    <div
      v-if="readerModalOpen && activeArticle"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 backdrop-blur-xs p-4 sm:p-6"
    >
      <div class="bg-white w-full max-w-4xl max-h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-in fade-in duration-200">
        <div class="p-5 border-b border-slate-200 flex items-center justify-between bg-slate-50/50">
          <div class="flex items-center gap-2.5">
            <span class="text-[11px] font-bold px-2.5 py-0.5 rounded-full uppercase tracking-wider bg-amber-100 text-amber-800">
              Artikel Solusi
            </span>
            <span class="text-xs text-slate-500 font-medium">• {{ activeArticle.system }} ({{ activeArticle.categoryName }})</span>
          </div>
          <button
            @click="closeReader"
            class="w-8 h-8 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-500 hover:text-slate-800 flex items-center justify-center transition-colors cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-6 sm:p-8 overflow-y-auto space-y-4">
          <h2 class="text-xl sm:text-2xl font-black text-slate-900 leading-snug">
            {{ activeArticle.title }}
          </h2>

          <div
            class="manual-book-content text-slate-700 text-sm sm:text-base leading-relaxed pt-4 border-t border-slate-100"
            v-html="activeArticle.content"
          ></div>
        </div>

        <div class="p-4 sm:p-5 border-t border-slate-200 bg-slate-50 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs">
          <div class="text-slate-500">
            Apakah artikel troubleshooting ini membantu Anda?
          </div>
          <div class="flex items-center gap-2">
            <button
              v-if="!feedbackGiven"
              @click="giveFeedback"
              class="px-3.5 py-1.5 rounded-lg border border-slate-300 bg-white hover:bg-slate-100 text-slate-700 font-medium inline-flex items-center gap-1.5 transition-colors cursor-pointer"
            >
              <ThumbsUp class="w-3.5 h-3.5" />
              <span>Ya, Berhasil</span>
            </button>
            <div v-else class="text-emerald-600 font-bold inline-flex items-center gap-1">
              <CheckCircle2 class="w-4 h-4" />
              <span>Terima kasih! Masukan Anda tercatat.</span>
            </div>
            <button
              @click="closeReader"
              class="px-4 py-1.5 rounded-lg bg-slate-900 text-white font-semibold hover:bg-slate-800 transition-colors cursor-pointer"
            >
              Tutup
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
:deep(.manual-book-content table) {
  width: 100% !important;
  border-collapse: collapse !important;
  margin: 1.25rem 0 !important;
}

:deep(.manual-book-content th) {
  background-color: #0b1a36 !important;
  color: #ffffff !important;
  padding: 10px 14px !important;
  border: 1px solid #334155 !important;
  font-weight: 700 !important;
  font-size: 0.875rem !important;
}

:deep(.manual-book-content td) {
  padding: 10px 14px !important;
  border: 1px solid #cbd5e1 !important;
  font-size: 0.875rem !important;
  color: #334155 !important;
  line-height: 1.5 !important;
}

:deep(.manual-book-content h2),
:deep(.manual-book-content h3) {
  font-size: 1.125rem !important;
  font-weight: 700 !important;
  color: #0f172a !important;
  margin-top: 1.5rem !important;
  margin-bottom: 0.75rem !important;
}

:deep(.manual-book-content ol),
:deep(.manual-book-content ul) {
  padding-left: 1.5rem !important;
  margin-bottom: 1rem !important;
}

:deep(.manual-book-content li) {
  margin-bottom: 0.5rem !important;
  line-height: 1.625 !important;
}

:deep(.manual-book-content p) {
  font-size: 0.9375rem !important;
  color: #475569 !important;
  line-height: 1.625 !important;
  margin-bottom: 0.75rem !important;
}

:deep(.manual-book-content img) {
  max-width: 100% !important;
  height: auto !important;
  border-radius: 8px !important;
  margin: 1rem 0 !important;
  border: 1px solid #e2e8f0 !important;
}
</style>
