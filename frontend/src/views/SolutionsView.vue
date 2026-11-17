<script setup>
import { ref, computed, onMounted, watch } from 'vue'
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
  Phone,
  Menu,
  Layers
} from 'lucide-vue-next'
import apiClient from '../api/client'
import { slugify, highlightKeyword } from '../utils/slug'
import { detectSystemIcon } from '../utils/iconResolver'

const router = useRouter()
const searchQuery = ref('')
const selectedSystemFilter = ref(2) // Default to PRESTO Resto POS (ID: 2)
const selectedModuleFilter = ref('all') // Default to 'all' (Semua Modul)
const isLoading = ref(false)
const mobileMenuOpen = ref(false)

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
  let formatted = html.replace(/https?:\/\/(127\.0\.0\.1|localhost)(:\d+)?\/uploads\//g, '/uploads/')

  // Auto-convert standalone YouTube links (anchor tag) ke responsive iframe
  formatted = formatted.replace(
    /<p>\s*<a\s+[^>]*href=["'](?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:watch\?.*v=|embed\/|shorts\/)|youtu\.be\/)([a-zA-Z0-9_-]{11})[^"']*["'][^>]*>.*?<\/a>\s*<\/p>/gi,
    '<div class="video-container my-4"><iframe src="https://www.youtube.com/embed/$1" width="100%" height="320" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen="allowfullscreen"></iframe></div>'
  )

  // Auto-convert standalone YouTube URL teks biasa ke responsive iframe
  formatted = formatted.replace(
    /<p>\s*(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:watch\?.*v=|embed\/|shorts\/)|youtu\.be\/)([a-zA-Z0-9_-]{11})(?:\S*)?\s*<\/p>/gi,
    '<div class="video-container my-4"><iframe src="https://www.youtube.com/embed/$1" width="100%" height="320" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen="allowfullscreen"></iframe></div>'
  )

  return formatted
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

// All 9 system categories directly from database `kategori` table + Semua Sistem di paling depan (Index 0)
const systemOptions = computed(() => {
  const result = []

  // 1. Opsi 'Semua Sistem' di urutan pertama (Index 0)
  result.push({
    value: 'all',
    code: 'ALL',
    label: 'Semua Sistem',
    articleCount: articles.value.length,
  })

  // 2. 9 Kategori Sistem Resmi Probus
  if (categories.value && categories.value.length > 0) {
    const sorted = [...categories.value].sort((a, b) => a.id - b.id)
    sorted.forEach(c => {
      const count = articles.value.filter(a => {
        const matchId = a.kategoriId === Number(c.id)
        const matchSys = a.system && a.system.toUpperCase() === c.name.toUpperCase()
        return matchId || matchSys
      }).length

      result.push({
        value: c.id,
        code: c.name,
        label: c.name, // PHIS, PRESTO, PPOS, SELF ORDER, etc.
        articleCount: count,
      })
    })
  } else {
    default9Categories.forEach(c => {
      const count = articles.value.filter(a => {
        const matchId = a.kategoriId === Number(c.value)
        const matchSys = a.system && a.system.toUpperCase() === c.label.toUpperCase()
        return matchId || matchSys
      }).length

      result.push({
        value: c.value,
        code: c.label,
        label: c.label,
        articleCount: count,
      })
    })
  }

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

// Hanya tampilkan modul yang memiliki artikel solusi (count > 0)
const modulesWithArticles = computed(() => {
  return availableModules.value.filter(m => m.count > 0)
})

// Sinkronisasi: jika modul yang dipilih tidak ada dalam daftar modul berartikel, kembalikan ke 'all'
watch([selectedSystemFilter, modulesWithArticles], () => {
  if (selectedModuleFilter.value !== 'all') {
    const exists = modulesWithArticles.value.some(m => m.id === selectedModuleFilter.value)
    if (!exists) {
      selectedModuleFilter.value = 'all'
    }
  }
})

// Total count of articles for current selected system
const totalArticlesCount = computed(() => {
  return articles.value.filter(a => {
    return selectedSystemFilter.value === 'all' || a.kategoriId === Number(selectedSystemFilter.value)
  }).length
})

// Filtered articles list
const filteredArticles = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  const queryWords = q ? q.split(/\s+/).filter(Boolean) : []

  return articles.value.filter(art => {
    // 1. System Filter
    const matchSystem = selectedSystemFilter.value === 'all' ||
      art.kategoriId === Number(selectedSystemFilter.value) ||
      art.system === selectedSystemFilter.value

    // 2. Module Filter
    const matchModule = selectedModuleFilter.value === 'all' ||
      art.menuSystemId === selectedModuleFilter.value ||
      (art.menuName && art.menuName.toLowerCase() === String(selectedModuleFilter.value).toLowerCase())

    // 3. Search Query (matches all keyword tokens across title, summary, category, system, and content)
    const matchQuery = queryWords.length === 0 || queryWords.every(word => {
      return (
        (art.title && art.title.toLowerCase().includes(word)) ||
        (art.summary && art.summary.toLowerCase().includes(word)) ||
        (art.categoryName && art.categoryName.toLowerCase().includes(word)) ||
        (art.system && art.system.toLowerCase().includes(word)) ||
        (art.content && art.content.toLowerCase().includes(word))
      )
    })

    return matchSystem && matchModule && matchQuery
  })
})

function openArticle(art) {
  router.push(`/solutions/${slugify(art.title)}`)
}

function goBackHome() {
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-slate-50/40 text-slate-800 font-sans flex flex-col antialiased">
    <!-- Main Navigation Bar -->
    <nav class="border-b border-slate-200/90 bg-white py-3.5 px-4 sm:px-6 lg:px-12 xl:px-16 sticky top-0 z-40 shadow-xs">
      <div class="w-full flex items-center justify-between">
        <div class="flex items-center gap-2.5 sm:gap-3.5">
          <div class="flex items-center cursor-pointer" @click="goBackHome">
            <img
              src="/probus-logo.png"
              alt="Probus System"
              class="h-7 sm:h-9 w-auto object-contain hover:opacity-90 transition-opacity"
            />
          </div>
          <span class="text-slate-300 font-light hidden sm:inline">|</span>
          <span class="text-slate-500 text-sm font-medium hidden md:inline">Bank Solusi Publik</span>

          <button
            @click="goBackHome"
            class="inline-flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-lg bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs transition-colors"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Kembali ke Beranda</span>
            <span class="sm:hidden">Beranda</span>
          </button>
        </div>

        <!-- Desktop Navigation Links -->
        <div class="hidden md:flex items-center gap-8 text-sm font-medium">
          <router-link to="/" class="text-slate-600 hover:text-slate-950 transition-colors">
            Panduan Sistem
          </router-link>
          <router-link to="/solutions" class="text-emerald-600 font-semibold hover:text-emerald-700 transition-colors relative py-1">
            Bank Solusi
            <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-emerald-600 rounded-full"></span>
          </router-link>
          <a
            href="https://probussystem.com/contact"
            target="_blank"
            rel="noopener noreferrer"
            class="text-slate-600 hover:text-slate-950 transition-colors"
          >
            Hubungi Kami
          </a>
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
          class="block px-3 py-2 rounded-lg text-sm font-semibold bg-emerald-50 text-emerald-700"
        >
          Bank Solusi
        </router-link>
        <a
          href="https://probussystem.com/contact"
          target="_blank"
          rel="noopener noreferrer"
          @click="mobileMenuOpen = false"
          class="block px-3 py-2 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-slate-900"
        >
          Hubungi Kami
        </a>
        <div class="pt-2 border-t border-slate-100 px-3">
          <a
            href="tel:+62361256500"
            class="inline-flex items-center gap-2 text-xs font-bold text-slate-700 hover:text-emerald-600"
          >
            <Phone class="w-3.5 h-3.5 text-slate-500" />
            <span>Phone: +62 361 256 500</span>
          </a>
        </div>
      </div>
    </nav>

    <!-- 3. Hero Section (Bank Solusi) -->
    <section class="pt-10 pb-8 sm:pt-16 sm:pb-12 lg:pt-20 lg:pb-16 px-4 sm:px-6 lg:px-12 xl:px-16 text-center bg-gradient-to-b from-white via-emerald-50/20 to-transparent">
      <div class="w-full max-w-4xl lg:max-w-5xl xl:max-w-6xl mx-auto">
        <h1 class="text-2xl sm:text-3xl lg:text-4xl xl:text-5xl font-extrabold text-slate-900 tracking-tight leading-tight">
          Pusat Bantuan & Solusi Klien
        </h1>
        <p class="mt-3 sm:mt-4 text-slate-500 text-sm sm:text-base lg:text-lg leading-relaxed max-w-3xl mx-auto">
          Pilih sistem dan modul di bawah, atau ketik kendala Anda untuk menemukan panduan langkah demi langkah.
        </p>

        <!-- Search Omnibox -->
        <div class="mt-7 sm:mt-8 w-full max-w-2xl sm:max-w-3xl lg:max-w-4xl xl:max-w-5xl mx-auto relative group">
          <div class="absolute left-4 sm:left-5 top-1/2 -translate-y-1/2 pointer-events-none flex items-center justify-center transition-colors duration-200 text-slate-400 group-focus-within:text-emerald-600">
            <Search class="w-5 h-5" />
          </div>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari kendala (contoh: printer tidak cetak, laci kasir macet, gagal QRIS)..."
            class="w-full pl-12 sm:pl-14 pr-32 sm:pr-40 py-3.5 sm:py-4 rounded-2xl text-sm sm:text-base bg-white border border-slate-200/90 text-slate-800 placeholder:text-slate-400 font-normal shadow-xs hover:border-slate-300 hover:shadow-sm focus:outline-none focus:border-emerald-500 focus:ring-4 focus:ring-emerald-500/15 focus:shadow-md transition-all duration-200"
          />

          <!-- Right Side Indicators (Tanda Cocok & Clear Button) -->
          <div class="absolute right-3.5 sm:right-4 top-1/2 -translate-y-1/2 flex items-center gap-1.5">
            <!-- Match indicator badge inside search bar -->
            <span
              v-if="searchQuery.trim() && filteredArticles.length > 0"
              class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200 shadow-2xs pointer-events-none animate-in fade-in"
            >
              <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600" />
              <span>{{ filteredArticles.length }} cocok</span>
            </span>
            <span
              v-else-if="searchQuery.trim() && filteredArticles.length === 0"
              class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] sm:text-[11px] font-semibold bg-rose-50 text-rose-600 border border-rose-200 pointer-events-none"
            >
              <span class="hidden sm:inline">Tidak ada</span> cocok
            </span>

            <button
              v-if="searchQuery"
              @click="searchQuery = ''"
              type="button"
              class="p-1.5 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors focus:outline-none cursor-pointer"
              title="Hapus pencarian"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- 1. PILIH SISTEM ROW (Grid 5x2 Simetris Tanpa Scrollbar) -->
        <div class="mt-8 text-left">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-600"></span>
              <span class="text-[11px] font-bold text-slate-500 uppercase tracking-wider">
                Pilih Sistem Probus
              </span>
            </div>
            <span class="text-[11px] text-slate-400">
              10 Pilihan Sistem
            </span>
          </div>

          <!-- Grid 5 Kolom x 2 Baris Simetris -->
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2.5">
            <button
              v-for="opt in systemOptions"
              :key="opt.value"
              @click="selectSystem(opt.value)"
              :class="[
                'group w-full flex items-center justify-between gap-2 text-xs px-3 py-2.5 rounded-xl font-bold transition-all duration-200 cursor-pointer select-none text-left',
                selectedSystemFilter === opt.value
                  ? 'bg-emerald-600 text-white shadow-sm shadow-emerald-500/25 ring-2 ring-emerald-500/30'
                  : 'bg-white text-slate-700 border border-slate-200/90 hover:border-emerald-300 hover:bg-emerald-50/60 hover:text-emerald-700 shadow-2xs'
              ]"
            >
              <div class="flex items-center gap-2 min-w-0 truncate">
                <!-- System Icon -->
                <div
                  :class="[
                    'w-6 h-6 rounded-lg flex items-center justify-center shrink-0 transition-colors',
                    selectedSystemFilter === opt.value
                      ? 'bg-white/20 text-white'
                      : 'bg-slate-100 text-slate-500 group-hover:bg-emerald-100 group-hover:text-emerald-600'
                  ]"
                >
                  <component
                    :is="opt.value === 'all' ? Layers : detectSystemIcon(opt.code)"
                    class="w-3.5 h-3.5"
                  />
                </div>

                <!-- System Name -->
                <span class="tracking-tight truncate">{{ opt.label }}</span>
              </div>

              <!-- Solusi Count Badge -->
              <span
                :class="[
                  'text-[10px] font-semibold px-2 py-0.5 rounded-full shrink-0 transition-colors',
                  selectedSystemFilter === opt.value
                    ? 'bg-white/25 text-white'
                    : (opt.articleCount > 0
                        ? 'bg-emerald-50 text-emerald-700 border border-emerald-100/80 group-hover:bg-emerald-100'
                        : 'bg-slate-100 text-slate-400')
                ]"
              >
                {{ opt.articleCount }} <span class="hidden xl:inline">Solusi</span>
              </span>
            </button>
          </div>
        </div>

        <!-- 2. MODUL BAR (Hanya modul yang memiliki artikel solusi) -->
        <div
          v-if="modulesWithArticles.length > 0"
          class="mt-2.5 p-2.5 sm:p-3 bg-slate-50/90 border border-slate-200/90 rounded-2xl flex flex-wrap items-center gap-2 text-left shadow-2xs animate-in fade-in duration-200"
        >
          <span class="text-xs font-black text-slate-500 uppercase tracking-wider px-2 shrink-0">
            MODUL:
          </span>

          <!-- Semua Modul Pill (Paling Depan) -->
          <button
            @click="selectedModuleFilter = 'all'"
            :class="[
              'px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all cursor-pointer inline-flex items-center gap-1.5',
              selectedModuleFilter === 'all'
                ? 'bg-emerald-600 text-white shadow-xs'
                : 'bg-white text-slate-700 border border-slate-200/90 hover:bg-emerald-50 hover:text-emerald-700 hover:border-emerald-200 shadow-2xs'
            ]"
          >
            <span>Semua Modul</span>
            <span
              :class="[
                'text-[10px] px-1.5 py-0.5 rounded-full font-bold',
                selectedModuleFilter === 'all'
                  ? 'bg-emerald-700/80 text-white'
                  : 'bg-slate-100 text-slate-600'
              ]"
            >
              {{ totalArticlesCount }}
            </span>
          </button>

          <!-- Modul-Modul yang memiliki artikel -->
          <button
            v-for="mod in modulesWithArticles"
            :key="mod.id"
            @click="selectedModuleFilter = mod.id"
            :class="[
              'px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all cursor-pointer inline-flex items-center gap-1.5',
              selectedModuleFilter === mod.id
                ? 'bg-emerald-600 text-white shadow-xs'
                : 'bg-white text-slate-700 border border-slate-200/90 hover:bg-emerald-50 hover:text-emerald-700 hover:border-emerald-200 shadow-2xs'
            ]"
          >
            <span>{{ mod.displayName }}</span>
            <span
              :class="[
                'text-[10px] px-1.5 py-0.5 rounded-full font-bold',
                selectedModuleFilter === mod.id
                  ? 'bg-emerald-700/80 text-white'
                  : 'bg-emerald-50 text-emerald-700 border border-emerald-100'
              ]"
            >
              {{ mod.count }}
            </span>
          </button>
        </div>
      </div>
    </section>

    <!-- 4. Main Articles List -->
    <main class="flex-1 w-full px-4 sm:px-6 lg:px-12 xl:px-16 pb-20">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-6">
        <div>
          <div class="flex items-center gap-2.5">
            <h2 class="text-xl font-bold text-slate-900">
              Daftar Artikel Solusi Masalah
            </h2>
            <span
              v-if="searchQuery.trim() && filteredArticles.length > 0"
              class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-bold bg-emerald-100 text-emerald-900 border border-emerald-300 shadow-2xs"
            >
              <CheckCircle2 class="w-3.5 h-3.5 text-emerald-700" />
              {{ filteredArticles.length }} Ditemukan
            </span>
          </div>
          <p class="text-xs sm:text-sm text-slate-500 mt-1">
            <span v-if="searchQuery.trim()">
              Menampilkan {{ filteredArticles.length }} artikel yang cocok dengan kata kunci 
              <strong class="text-slate-900 bg-emerald-100/90 text-emerald-950 font-bold px-1.5 py-0.5 rounded border border-emerald-200/80">"{{ searchQuery.trim() }}"</strong>
            </span>
            <span v-else>
              Menampilkan {{ filteredArticles.length }} artikel pemecahan masalah teknis dari basis data
            </span>
          </p>
        </div>

        <div v-if="isLoading" class="flex items-center gap-2 text-xs text-slate-400">
          <Loader2 class="w-4 h-4 animate-spin text-emerald-500" />
          <span>Memuat artikel...</span>
        </div>
      </div>

      <!-- Articles Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5 lg:gap-6 xl:gap-8">
        <div
          v-for="art in filteredArticles"
          :key="art.id"
          @click="openArticle(art)"
          class="bg-white rounded-2xl p-6 border border-slate-200 hover:border-emerald-400 hover:shadow-md transition-all cursor-pointer flex flex-col justify-between group"
        >
          <div>
            <div class="flex items-center justify-between gap-2 mb-3.5">
              <span class="text-xs font-bold px-2.5 py-0.5 rounded-md bg-emerald-50 text-emerald-700 border border-emerald-200/60">
                {{ art.system }} • {{ art.categoryName }}
              </span>
              <span class="text-xs text-slate-400 flex items-center gap-1">
                <Clock class="w-3 h-3" /> {{ art.readTime }}
              </span>
            </div>

            <h3
              class="font-bold text-slate-900 text-base group-hover:text-emerald-600 transition-colors leading-snug"
              v-html="highlightKeyword(art.title, searchQuery)"
            ></h3>

            <p
              class="text-xs sm:text-sm text-slate-500 mt-2.5 line-clamp-3 leading-relaxed"
              v-html="highlightKeyword(art.summary, searchQuery)"
            ></p>
          </div>

          <div class="pt-5 mt-5 border-t border-slate-100 flex items-center justify-between text-xs">
            <span class="text-slate-400">Terverifikasi Tim Teknis</span>
            <span class="text-emerald-600 font-semibold flex items-center gap-1 group-hover:translate-x-0.5 transition-transform">
              <span>Buka Solusi</span>
              <ChevronRight class="w-3.5 h-3.5" />
            </span>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!isLoading && filteredArticles.length === 0" class="text-center py-16 bg-white rounded-2xl border border-slate-200">
        <AlertCircle class="w-12 h-12 text-slate-300 mx-auto mb-3" />
        <h3 class="text-base font-bold text-slate-800">
          <span v-if="searchQuery.trim()">Tidak ada artikel yang cocok dengan "{{ searchQuery.trim() }}"</span>
          <span v-else>Tidak ada artikel solusi yang sesuai</span>
        </h3>
        <p class="text-xs sm:text-sm text-slate-500 mt-1 max-w-md mx-auto">
          <span v-if="searchQuery.trim()">
            Coba gunakan kata kunci lain seperti nama perangkat (printer, drawer, lan) atau reset pencarian.
          </span>
          <span v-else>
            Belum ada artikel yang dipublikasikan untuk modul atau filter ini. Hubungi tim implementator Probus untuk bantuan langsung.
          </span>
        </p>
        <div class="mt-4 flex flex-wrap items-center justify-center gap-2.5">
          <button
            v-if="searchQuery.trim()"
            @click="searchQuery = ''"
            class="px-4 py-2 rounded-lg bg-emerald-600 text-white text-xs font-semibold hover:bg-emerald-700 transition-colors cursor-pointer"
          >
            Hapus Kata Kunci
          </button>
          <button
            @click="selectedModuleFilter = 'all'; selectedSystemFilter = 'all'; searchQuery = ''"
            class="px-4 py-2 rounded-lg bg-slate-900 text-white text-xs font-semibold hover:bg-slate-800 transition-colors cursor-pointer"
          >
            Tampilkan Semua Modul & Sistem
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
    <footer class="border-t border-slate-200 bg-white py-6 px-4 sm:px-6 lg:px-12 xl:px-16 text-xs sm:text-sm text-slate-500">
      <div class="w-full flex flex-col sm:flex-row items-center justify-between gap-3 text-center sm:text-left">
        <p>© 2026 CV Mitra Solusi Mandiri • Bank Solusi Resmi</p>
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

:deep(.manual-book-content iframe),
:deep(.manual-book-content video),
:deep(.video-container iframe) {
  width: 100% !important;
  max-width: 100% !important;
  aspect-ratio: 16 / 9 !important;
  height: auto !important;
  border-radius: 8px !important;
  margin: 1rem 0 !important;
  border: 1px solid #e2e8f0 !important;
  display: block !important;
}
</style>
