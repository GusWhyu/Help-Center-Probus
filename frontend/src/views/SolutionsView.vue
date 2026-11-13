<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Search,
  ArrowLeft,
  ChevronRight,
  AlertCircle,
  Clock,
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
  Loader2,
  Phone
} from 'lucide-vue-next'
import apiClient from '../api/client'

const router = useRouter()
const searchQuery = ref('')
const selectedSystemFilter = ref(2) // Default to PRESTO Resto POS (ID: 2)
const selectedModuleFilter = ref('all') // Default to 'all' (Semua Modul)
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

// Categories and menus from database
const categories = ref([])

// Articles list from database
const articles = ref([])

// Load categories with their menu systems
async function loadCategories() {
  try {
    const res = await apiClient.get('/categories')
    if (res.data && res.data.data && Array.isArray(res.data.data)) {
      categories.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load categories:', err)
  }
}

// Load articles from database (only published & public articles)
async function loadArticles() {
  isLoading.value = true
  try {
    const res = await apiClient.get('/articles', { params: { limit: 100 } })
    if (res.data && res.data.data && Array.isArray(res.data.data)) {
      articles.value = res.data.data
        .filter(item => item.status === 'published' && item.visibility === 'publik')
        .map(item => {
          const primaryMenu = item.menu_systems && item.menu_systems.length > 0 ? item.menu_systems[0] : null
          const katId = primaryMenu ? primaryMenu.kategori_id : 2
          const sysCode = systemNameMap[katId] || 'PRESTO'
          const menuName = primaryMenu ? primaryMenu.name : 'Umum'
          const menuId = primaryMenu ? primaryMenu.id : null

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
            menuSystemId: menuId,
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
  loadCategories()
  loadArticles()
})

// Exact database categories from `kategori` table
const default9Categories = [
  { value: 1, label: 'PHIS' },
  { value: 2, label: 'PRESTO' },
  { value: 3, label: 'PPOS' },
  { value: 4, label: 'SELF ORDER' },
  { value: 5, label: 'PR-PO ONLINE' },
  { value: 6, label: 'ORDER TAKER' },
  { value: 7, label: 'VRESTROPICAL' },
  { value: 8, label: 'PPE' },
  { value: 9, label: 'BARONG POS' },
]

// All 9 system categories directly from database `kategori` table + Semua Sistem
const systemOptions = computed(() => {
  const result = []

  if (categories.value && categories.value.length > 0) {
    // Sort strictly by id ASC (1 to 9) to match database table `kategori`
    const sorted = [...categories.value].sort((a, b) => a.id - b.id)

    sorted.forEach(c => {
      result.push({
        value: c.id,
        code: c.name,
        label: c.name, // Exact database name: PHIS, PRESTO, PPOS, SELF ORDER, etc.
      })
    })
  } else {
    default9Categories.forEach(c => {
      result.push({
        value: c.value,
        code: c.label,
        label: c.label,
      })
    })
  }

  // Append 'Semua Sistem' at the end
  result.push({
    value: 'all',
    code: 'ALL',
    label: 'Semua Sistem',
  })

  return result
})

function selectSystem(systemVal) {
  selectedSystemFilter.value = systemVal
  selectedModuleFilter.value = 'all'
}

// 10 Real Modules for the selected system directly from `menu_systems`
const availableModules = computed(() => {
  let menus = []
  if (selectedSystemFilter.value === 'all') {
    // If all systems selected, use PRESTO menus as baseline
    const presto = categories.value.find(c => c.id === 2 || c.name === 'PRESTO')
    if (presto && presto.menus && presto.menus.length > 0) {
      menus = presto.menus
    }
  } else {
    const cat = categories.value.find(c => c.id === Number(selectedSystemFilter.value) || c.name === selectedSystemFilter.value)
    if (cat && cat.menus && cat.menus.length > 0) {
      menus = cat.menus
    }
  }

  // Fallback to standard 10 modules if not yet loaded from DB
  if (!menus || menus.length === 0) {
    menus = [
      { id: 34, order: 1, name: 'Setup Master' },
      { id: 35, order: 2, name: 'Transaksi' },
      { id: 36, order: 3, name: 'Outlet Activity' },
      { id: 37, order: 4, name: 'Cost Control' },
      { id: 38, order: 5, name: 'Purcashing AP' },
      { id: 39, order: 6, name: 'AR' },
      { id: 41, order: 7, name: 'General Cashier' },
      { id: 42, order: 8, name: 'Accounting System' },
      { id: 43, order: 9, name: 'Report' },
      { id: 40, order: 10, name: 'Tools' },
    ]
  }

  // Sort menus strictly by order (1 to 10)
  const sortedMenus = [...menus].sort((a, b) => (a.order || 0) - (b.order || 0))

  return sortedMenus.map(m => {
    // Calculate count of articles for this module under currently selected system
    const count = articles.value.filter(a => {
      const matchSys = selectedSystemFilter.value === 'all' || a.kategoriId === Number(selectedSystemFilter.value)
      const matchMenu = a.menuSystemId === m.id || (a.menuName && a.menuName.toLowerCase() === m.name.toLowerCase())
      return matchSys && matchMenu
    }).length

    return {
      id: m.id,
      name: m.name,
      displayName: m.name, // Exact database name: Setup Master, Transaksi, Tools, etc.
      order: m.order,
      count,
    }
  })
})

// Total count of articles for current selected system
const totalArticlesCount = computed(() => {
  return articles.value.filter(a => {
    return selectedSystemFilter.value === 'all' || a.kategoriId === Number(selectedSystemFilter.value)
  }).length
})

// Filtered articles list
const filteredArticles = computed(() => {
  return articles.value.filter(art => {
    // 1. System Filter
    const matchSystem = selectedSystemFilter.value === 'all' ||
      art.kategoriId === Number(selectedSystemFilter.value) ||
      art.system === selectedSystemFilter.value

    // 2. Module Filter
    const matchModule = selectedModuleFilter.value === 'all' ||
      art.menuSystemId === selectedModuleFilter.value ||
      (art.menuName && art.menuName.toLowerCase() === String(selectedModuleFilter.value).toLowerCase())

    // 3. Search Query
    const q = searchQuery.value.trim().toLowerCase()
    const matchQuery = !q ||
      (art.title && art.title.toLowerCase().includes(q)) ||
      (art.summary && art.summary.toLowerCase().includes(q)) ||
      (art.categoryName && art.categoryName.toLowerCase().includes(q)) ||
      (art.system && art.system.toLowerCase().includes(q))

    return matchSystem && matchModule && matchQuery
  })
})

function openArticle(art) {
  router.push(`/solutions/${art.id}`)
}

function goBackHome() {
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-slate-50/40 text-slate-800 font-sans flex flex-col antialiased">
    <!-- 1. Utility Top Bar -->
    <header class="border-b border-slate-200/70 bg-white text-xs text-slate-500 py-2.5 px-6 sm:px-12 lg:px-20 xl:px-28">
      <div class="max-w-[1600px] mx-auto flex items-center justify-end">
        <div class="flex items-center gap-2 text-slate-600">
          <span class="text-slate-400 hidden sm:inline">Butuh bantuan cepat?</span>
          <a
            href="tel:+62361256500"
            class="inline-flex items-center gap-1.5 font-bold text-slate-800 hover:text-sky-600 transition-colors"
          >
            <Phone class="w-3.5 h-3.5 text-slate-500" />
            <span>Phone: +62 361 256 500</span>
          </a>
        </div>
      </div>
    </header>

    <!-- 2. Main Navigation Bar -->
    <nav class="border-b border-slate-200/80 bg-white/95 backdrop-blur-sm py-3.5 px-6 sm:px-12 lg:px-20 xl:px-28 sticky top-0 z-40 shadow-xs">
      <div class="max-w-[1600px] mx-auto flex items-center justify-between">
        <div class="flex items-center gap-3.5">
          <div class="flex items-center gap-2.5 cursor-pointer" @click="goBackHome">
            <img
              src="/probus-logo.png"
              alt="Probus System"
              class="h-8 sm:h-9 w-auto object-contain hover:opacity-90 transition-opacity"
            />
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
    <section class="pt-12 pb-10 px-6 sm:px-12 lg:px-20 xl:px-28 text-center bg-gradient-to-b from-white via-amber-50/20 to-transparent">
      <div class="max-w-5xl mx-auto">
        <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-md bg-amber-50 text-amber-800 text-xs font-semibold border border-amber-200/60 mb-3 shadow-2xs">
          <AlertCircle class="w-3.5 h-3.5" />
          <span>Knowledge Base & Troubleshooting</span>
        </div>

        <h1 class="text-3xl sm:text-4xl font-extrabold text-slate-900 tracking-tight">
          Bank Solusi Kendala & Troubleshooting
        </h1>
        <p class="mt-2.5 text-slate-500 text-sm sm:text-base leading-relaxed max-w-2xl mx-auto">
          Temukan solusi cepat untuk masalah printer kasir, cash drawer, jaringan LAN, dan error operasional sistem Probus.
        </p>

        <!-- Search Omnibox -->
        <div class="mt-6 max-w-3xl mx-auto relative shadow-2xs">
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

        <!-- 1. PILIH SISTEM ROW (9 Kategori Produk Probus + Semua Sistem) -->
        <div class="mt-7 text-left">
          <div class="text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2.5">
            PILIH SISTEM
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <button
              v-for="opt in systemOptions"
              :key="opt.value"
              @click="selectSystem(opt.value)"
              :class="[
                'px-3.5 py-2 rounded-xl text-xs sm:text-sm font-bold transition-all cursor-pointer',
                selectedSystemFilter === opt.value
                  ? 'bg-slate-950 text-white shadow-xs ring-2 ring-slate-900/10'
                  : 'bg-white text-slate-700 border border-slate-200 hover:bg-slate-50 hover:text-slate-950 shadow-2xs'
              ]"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>

        <!-- 2. MODUL BAR (10 Modul Sistem dari Database) -->
        <div class="mt-3.5 p-3 sm:p-3.5 bg-slate-50/90 border border-slate-200/90 rounded-2xl flex flex-wrap items-center gap-2 text-left shadow-2xs">
          <span class="text-xs font-black text-slate-500 uppercase tracking-wider px-2 shrink-0">
            MODUL:
          </span>

          <!-- Individual Module Pills (10 Modul) -->
          <button
            v-for="mod in availableModules"
            :key="mod.id"
            @click="selectedModuleFilter = mod.id"
            :class="[
              'px-3.5 py-1.5 rounded-lg text-xs font-semibold transition-all cursor-pointer inline-flex items-center gap-1.5',
              selectedModuleFilter === mod.id
                ? 'bg-sky-600 text-white shadow-xs'
                : 'bg-white text-slate-700 border border-slate-200/90 hover:bg-slate-100 hover:text-slate-900 shadow-2xs'
            ]"
          >
            <span>{{ mod.displayName }}</span>
            <span
              :class="[
                'text-[10px] px-1.5 py-0.2 rounded-full font-bold',
                selectedModuleFilter === mod.id
                  ? 'bg-sky-700/70 text-white'
                  : (mod.count > 0 ? 'bg-amber-100 text-amber-800' : 'bg-slate-100 text-slate-400')
              ]"
            >
              {{ mod.count }}
            </span>
          </button>

          <!-- Semua Modul Pill -->
          <button
            @click="selectedModuleFilter = 'all'"
            :class="[
              'px-3.5 py-1.5 rounded-lg text-xs font-semibold transition-all cursor-pointer inline-flex items-center gap-1.5',
              selectedModuleFilter === 'all'
                ? 'bg-sky-600 text-white shadow-xs'
                : 'bg-white text-slate-700 border border-slate-200/90 hover:bg-slate-100 hover:text-slate-900 shadow-2xs'
            ]"
          >
            <span>Semua Modul</span>
            <span
              :class="[
                'text-[10px] px-1.5 py-0.2 rounded-full font-bold',
                selectedModuleFilter === 'all'
                  ? 'bg-sky-700/70 text-white'
                  : 'bg-slate-100 text-slate-500'
              ]"
            >
              {{ totalArticlesCount }}
            </span>
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
        <h3 class="text-base font-bold text-slate-800">Tidak ada artikel solusi yang sesuai</h3>
        <p class="text-xs sm:text-sm text-slate-500 mt-1 max-w-md mx-auto">
          Belum ada artikel yang dipublikasikan untuk modul atau filter ini. Hubungi tim implementator Probus untuk bantuan langsung.
        </p>
        <div class="mt-4 flex flex-wrap items-center justify-center gap-2.5">
          <button
            @click="selectedModuleFilter = 'all'; searchQuery = ''"
            class="px-4 py-2 rounded-lg bg-slate-900 text-white text-xs font-semibold hover:bg-slate-800 transition-colors cursor-pointer"
          >
            Tampilkan Semua Modul
          </button>
          <a
            href="https://wa.me/6281910875514"
            target="_blank"
            class="px-4 py-2 rounded-lg bg-emerald-600 text-white text-xs font-semibold hover:bg-emerald-700 transition-colors inline-flex items-center gap-1.5"
          >
            <span>Hubungi WhatsApp Support</span>
          </a>
        </div>
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
