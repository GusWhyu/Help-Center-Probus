<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSystemIdFromParam, getSystemSlug, slugify, highlightKeyword, registerCategories, ensureCategoriesLoaded, resolveSystemId } from '../utils/slug'
import apiClient from '../api/client'
import { detectModuleIcon, parseModuleDescription, resolveModuleIcon, isImageIcon } from '../utils/iconResolver'
import {
  ArrowLeft,
  ArrowRight,
  Search,
  BookOpen,
  Lightbulb,
  HelpCircle,
  X,
  ThumbsUp,
  CheckCircle2,
  Clock,
  Info,
  Settings,
  UtensilsCrossed,
  Store,
  PieChart,
  ShoppingBag,
  Receipt,
  Wallet,
  BarChart3,
  TrendingUp,
  Printer,
  Phone,
  Menu
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const mobileMenuOpen = ref(false)

// Get system ID from route parameter, supports both slugs ('self-order') and legacy IDs ('4')
const systemId = ref(getSystemIdFromParam(route.params.id, 2))

// Active Tab: 'manual' or 'solusi'
const activeTab = ref('manual')
const searchQuery = ref('')
const selectedCardId = ref(2) // default highlight: card 02 (Transaksi)

// System metadata mapping
const systemsInfo = {
  1: { id: 1, name: 'PHIS', title: 'Panduan Pengguna PHIS', badge: 'Hotel & Resort' },
  2: { id: 2, name: 'PRESTO', title: 'Panduan Pengguna PRESTO', badge: 'Restoran & POS' },
  3: { id: 3, name: 'PPOS', title: 'Panduan Pengguna PPOS', badge: 'Retail & Kasir' },
  4: { id: 4, name: 'SELF ORDER', title: 'Panduan Pengguna SELF ORDER', badge: 'Self Service' },
  5: { id: 5, name: 'PR-PO ONLINE', title: 'Panduan Pengguna PR-PO ONLINE', badge: 'Logistik & Gudang' },
  6: { id: 6, name: 'ORDER TAKER', title: 'Panduan Pengguna ORDER TAKER', badge: 'Waiter App' },
  7: { id: 7, name: 'VRESTROPICAL', title: 'Panduan Pengguna VRESTROPICAL', badge: 'Club & Lounge' },
  8: { id: 8, name: 'PPE', title: 'Panduan Pengguna PPE', badge: 'Facility & Asset' },
  9: { id: 9, name: 'BARONG POS', title: 'Panduan Pengguna BARONG POS', badge: 'Cloud POS' },
}

const categoryData = ref(null)

const currentSystem = computed(() => {
  if (categoryData.value && categoryData.value.name) {
    return {
      id: categoryData.value.id,
      name: categoryData.value.name,
      title: `Panduan Pengguna ${categoryData.value.name}`,
      badge: systemsInfo[systemId.value]?.badge || 'Sistem Probus',
    }
  }
  return systemsInfo[systemId.value] || {
    id: systemId.value,
    name: `Sistem ${systemId.value}`,
    title: `Panduan Pengguna Sistem ${systemId.value}`,
    badge: 'Sistem Probus',
  }
})

// Modul fallback PRESTO yang riil memiliki buku panduan di database dbprobuscs
const prestoFallbackModules = [
  {
    id: 34,
    dbMenuId: 34,
    order: 1,
    name: 'Setup Master',
    title: '01. Setup Master',
    count: 3,
    countText: '3 Panduan',
    icon: Settings,
    desc: 'Modul setup master data, parameter sistem, konfigurasi outlet resto, denah meja, dan pengaturan menu makanan.',
    estTime: 'Est. 4 menit',
  },
  {
    id: 35,
    dbMenuId: 35,
    order: 2,
    name: 'Transaksi',
    title: '02. Transaksi',
    count: 2,
    countText: '2 Panduan',
    icon: UtensilsCrossed,
    desc: 'Modul transaksi penjualan operasional kasir POS, order taker meja, custom modifier pesanan, dan pembayaran split bill.',
    estTime: 'Est. 7 menit',
  },
  {
    id: 36,
    dbMenuId: 36,
    order: 3,
    name: 'Outlet Activity',
    title: '03. Outlet Activity',
    count: 1,
    countText: '1 Panduan',
    icon: Store,
    desc: 'Modul aktivitas dan operasional outlet, pembagian shift kerja kasir, serah terima modal kas, dan pantauan kegiatan harian.',
    estTime: 'Est. 5 menit',
  },
  {
    id: 40,
    dbMenuId: 40,
    order: 10,
    name: 'Tools',
    title: '10. Tools',
    count: 1,
    countText: '1 Panduan',
    icon: Printer,
    desc: 'Modul utilitas sistem, konfigurasi hardware printer kasir & dapur LAN/Bluetooth, cash drawer laci uang, dan diagnostik.',
    estTime: 'Est. 6 menit',
  },
]

function formatArticleContent(html) {
  if (!html) return ''
  let formatted = html.replace(/https?:\/\/(127\.0\.0\.1|localhost)(:\d+)?\/uploads\//g, '/uploads/')

  // Auto-convert standalone YouTube links (anchor tag) ke responsive iframe
  formatted = formatted.replace(
    /<p>\s*<a\s+[^>]*href=["'](?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:watch\?.*v=|embed\/|shorts\/)|youtu\.be\/)([a-zA-Z0-9_-]{11})[^"']*["'][^>]*>.*?<\/a>\s*<\/p>/gi,
    '<div class="video-container my-4"><iframe src="https://www.youtube.com/embed/$1" title="YouTube video player" width="100%" height="320" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen="allowfullscreen"></iframe></div>'
  )

  // Auto-convert standalone YouTube URL teks biasa ke responsive iframe
  formatted = formatted.replace(
    /<p>\s*(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:watch\?.*v=|embed\/|shorts\/)|youtu\.be\/)([a-zA-Z0-9_-]{11})(?:\S*)?\s*<\/p>/gi,
    '<div class="video-container my-4"><iframe src="https://www.youtube.com/embed/$1" title="YouTube video player" width="100%" height="320" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen="allowfullscreen"></iframe></div>'
  )

  return formatted
}

function cleanSummary(text) {
  if (!text) return ''
  return text.replace(/&nbsp;/g, ' ').replace(/<[^>]*>?/gm, '').trim()
}

// Troubleshooting Articles for "Solusi" Tab loaded dynamically from database
const troubleshootingArticles = ref([])

async function loadSystemArticles() {
  try {
    const currentId = systemId.value
    if (!currentId) return
    const res = await apiClient.get('/articles', {
      params: { kategori_id: currentId, limit: 50 }
    })
    if (res.data && res.data.data && Array.isArray(res.data.data)) {
      troubleshootingArticles.value = res.data.data
        .filter(item => item.status === 'published' && item.visibility === 'publik')
        .map(item => {
        const primaryMenu = item.menu_systems && item.menu_systems.length > 0 ? item.menu_systems[0] : null
        const menuName = primaryMenu ? primaryMenu.name : 'Tools'
        const words = (item.content || '').split(/\s+/).length
        const minutes = Math.max(2, Math.min(10, Math.ceil(words / 150)))

        return {
          id: item.id,
          title: item.title,
          menuName: menuName,
          readTime: `${minutes} menit baca`,
          summary: cleanSummary(item.summary) || 'Panduan pemecahan masalah teknis dan operasional sistem Probus.',
          content: formatArticleContent(item.content),
        }
      })
    } else {
      troubleshootingArticles.value = []
    }
  } catch (err) {
    console.error('Failed to load system troubleshooting articles:', err)
    troubleshootingArticles.value = []
  }
}

// Dynamic modules from database
const dynamicModules = ref([])
const isModulesLoaded = ref(false)

async function loadSystemModules() {
  try {
    // 1. Pastikan seluruh kategori dari database teregistrasi ke registry slug
    await ensureCategoriesLoaded()

    // 2. Resolve systemId secara akurat berdasarkan parameter rute
    const resolvedId = await resolveSystemId(route.params.id, systemId.value || 2)
    systemId.value = resolvedId

    // 3. Ambil data kategori & modul untuk sistem yang benar
    const res = await apiClient.get(`/categories/${resolvedId}`)
    if (res.data && res.data.data) {
      registerCategories([res.data.data])
      categoryData.value = res.data.data
    }
    if (res.data && res.data.data && res.data.data.menus && res.data.data.menus.length > 0) {
      const activeOnly = []
      res.data.data.menus.forEach((menu, idx) => {
        const orderNum = menu.order || (idx + 1)
        const code = String(orderNum).padStart(2, '0')

        // Hitung dan deduplikasi manual books yang berstatus published secara riil
        const seenMbIds = new Set()
        const validManuals = (menu.manual_books || []).filter(mb => {
          if (!mb || !mb.id || seenMbIds.has(mb.id)) return false
          seenMbIds.add(mb.id)
          return mb.status === 'published'
        })
        const manualCount = validManuals.length

        // HANYA tampilkan modul yang memang memiliki buku panduan (count > 0)
        if (manualCount > 0) {
          const parsed = parseModuleDescription(menu.description)
          const customIcon = menu.icon || parsed.icon
          const cleanDesc = (parsed.text || menu.description || '').trim()
          const isImage = isImageIcon(customIcon)
          const moduleIcon = isImage ? customIcon : resolveModuleIcon(customIcon, menu.name, orderNum)
          activeOnly.push({
            id: menu.id,
            order: orderNum,
            name: menu.name,
            title: `${code}. ${menu.name}`,
            desc: cleanDesc || 'Panduan operasional dan tata cara penggunaan fitur modul.',
            icon: moduleIcon,
            isImage: isImage,
            count: manualCount,
            countText: `${manualCount} Panduan`,
            estTime: 'Est. 5 menit',
            manualBooks: validManuals,
          })
        }
      })
      dynamicModules.value = activeOnly
    } else {
      dynamicModules.value = []
    }
  } catch (err) {
    console.warn('API /categories/:id failed, using fallback modules:', err)
    if (systemId.value === 2) {
      dynamicModules.value = prestoFallbackModules
    } else {
      dynamicModules.value = []
    }
  } finally {
    isModulesLoaded.value = true
  }
}

onMounted(async () => {
  await loadSystemModules()
  await loadSystemArticles()
  const currentParam = route.params.id
  const canonicalSlug = getSystemSlug(systemId.value, currentSystem.value)
  if (currentParam && !isNaN(Number(currentParam))) {
    router.replace(`/system/${canonicalSlug}`)
  }
})

watch(
  () => route.params.id,
  async (newParam) => {
    if (newParam) {
      isModulesLoaded.value = false
      await loadSystemModules()
      await loadSystemArticles()
      const canonicalSlug = getSystemSlug(systemId.value, currentSystem.value)
      if (!isNaN(Number(newParam))) {
        router.replace(`/system/${canonicalSlug}`)
      }
    }
  }
)

const activeModuleCards = computed(() => {
  return dynamicModules.value
})

// Filtered lists
const filteredModuleCards = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  const base = activeModuleCards.value
  if (!q) return base
  return base.filter(c =>
    c.title.toLowerCase().includes(q) || c.desc.toLowerCase().includes(q)
  )
})

const filteredTroubleshooting = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return troubleshootingArticles.value
  return troubleshootingArticles.value.filter(a =>
    a.title.toLowerCase().includes(q) || a.summary.toLowerCase().includes(q)
  )
})

const currentMatchesCount = computed(() => {
  return activeTab.value === 'manual' ? filteredModuleCards.value.length : filteredTroubleshooting.value.length
})

// Reader Modal
const readerModalOpen = ref(false)
const activeReaderItem = ref(null)
const feedbackGiven = ref(false)

function openReader(item, type = 'manual') {
  selectedCardId.value = item.id
  const sysSlug = getSystemSlug(systemId.value, currentSystem.value)
  if (type === 'manual') {
    const modSlug = slugify(item.title || item.name || String(item.id))
    router.push(`/system/${sysSlug}/guide/${modSlug}`)
    return
  }
  const artSlug = slugify(item.title || String(item.id))
  router.push(`/solutions/${artSlug}`)
}

function closeReader() {
  readerModalOpen.value = false
  activeReaderItem.value = null
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
    <!-- Main Navigation Bar with "Kembali ke Pilihan Sistem" Button -->
    <nav class="border-b border-slate-200/90 bg-white py-3.5 px-4 sm:px-6 lg:px-12 xl:px-16 sticky top-0 z-40 shadow-xs">
      <div class="w-full flex items-center justify-between">
        <!-- Logo, Title & Back Button -->
        <div class="flex items-center gap-2.5 sm:gap-3.5">
          <div class="flex items-center cursor-pointer" @click="goBackHome">
            <img
              src="/probus-logo.png"
              alt="Probus System"
              class="h-7 sm:h-9 w-auto object-contain hover:opacity-90 transition-opacity"
            />
          </div>
          <span class="text-slate-300 font-light hidden sm:inline">|</span>
          <span class="text-slate-500 text-sm font-medium hidden md:inline">Pusat Bantuan Publik</span>

          <!-- Pill Button: "← Kembali ke Pilihan Sistem" -->
          <button
            @click="goBackHome"
            class="inline-flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-lg bg-sky-50 hover:bg-sky-100 text-sky-700 font-semibold text-xs border border-sky-100 transition-colors shadow-2xs"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Kembali ke Pilihan Sistem</span>
            <span class="sm:hidden">Pilihan Sistem</span>
          </button>
        </div>

        <!-- Desktop Navigation Links -->
        <div class="hidden md:flex items-center gap-8 text-sm font-medium">
          <router-link to="/" class="text-sky-600 font-semibold hover:text-sky-700 transition-colors relative py-1">
            Panduan Sistem
            <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-sky-600 rounded-full"></span>
          </router-link>
          <router-link to="/solutions" class="text-slate-600 hover:text-slate-950 transition-colors">
            Bank Solusi
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
          class="block px-3 py-2 rounded-lg text-sm font-semibold bg-sky-50 text-sky-700"
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
            class="inline-flex items-center gap-2 text-xs font-bold text-slate-700 hover:text-sky-600"
          >
            <Phone class="w-3.5 h-3.5 text-slate-500" />
            <span>Phone: +62 361 256 500</span>
          </a>
        </div>
      </div>
    </nav>

    <!-- 3. Header Hero Section (Centered with Pill Switcher & Note) -->
    <section class="pt-10 pb-8 sm:pt-12 sm:pb-10 px-4 sm:px-6 lg:px-12 xl:px-16 text-center bg-gradient-to-b from-white via-sky-50/20 to-transparent">
      <div class="w-full max-w-4xl lg:max-w-5xl mx-auto">
        <!-- Category Badge -->
        <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-md bg-sky-50 text-sky-700 text-xs font-semibold border border-sky-100 mb-3.5 shadow-2xs">
          <span>🏛</span>
          <span>{{ currentSystem.badge }}</span>
        </div>

        <!-- Main Title -->
        <h1 class="text-2xl sm:text-3xl lg:text-4xl font-extrabold text-slate-900 tracking-tight">
          {{ currentSystem.title }}
        </h1>
        <p class="mt-2.5 text-slate-500 text-sm sm:text-base leading-relaxed max-w-2xl mx-auto">
          Pilih modul untuk mempelajari tata cara penggunaan sistem atau beralih ke solusi kendala teknis.
        </p>

        <!-- Search Bar -->
        <div class="mt-6 w-full max-w-2xl sm:max-w-3xl lg:max-w-4xl mx-auto relative group">
          <div class="absolute left-4 top-1/2 -translate-y-1/2 pointer-events-none flex items-center justify-center transition-colors duration-200 text-slate-400 group-focus-within:text-sky-600">
            <Search class="w-4.5 h-4.5" />
          </div>
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="`Ketik fitur, nomor modul, atau topik panduan ${currentSystem.name}...`"
            class="w-full pl-11 sm:pl-12 pr-32 sm:pr-40 py-3 sm:py-3.5 rounded-2xl text-sm bg-white border border-slate-200/90 text-slate-800 placeholder:text-slate-400 font-normal shadow-xs hover:border-slate-300 hover:shadow-sm focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/15 focus:shadow-md transition-all duration-200"
          />

          <!-- Right Side Indicators (Tanda Cocok & Clear Button) -->
          <div class="absolute right-3.5 top-1/2 -translate-y-1/2 flex items-center gap-1.5">
            <span
              v-if="searchQuery.trim() && currentMatchesCount > 0"
              class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-sky-100 text-sky-800 border border-sky-300 shadow-2xs pointer-events-none animate-in fade-in"
            >
              <CheckCircle2 class="w-3 h-3 text-sky-600" />
              <span>{{ currentMatchesCount }} cocok</span>
            </span>
            <span
              v-else-if="searchQuery.trim() && currentMatchesCount === 0"
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

        <!-- Segmented Tab Switcher (Matching Design) -->
        <div class="mt-6 inline-flex p-1 bg-white border border-slate-200 rounded-xl shadow-xs">
          <!-- Tab 1: Manual Book (Active) -->
          <button
            @click="activeTab = 'manual'"
            :class="[
              'inline-flex items-center gap-2 px-5 py-2.5 rounded-lg text-xs sm:text-sm font-bold transition-all',
              activeTab === 'manual'
                ? 'bg-sky-600 text-white shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            ]"
          >
            <BookOpen class="w-4 h-4" />
            <span>Manual Book</span>
            <span
              :class="[
                'text-[10px] font-bold px-2 py-0.5 rounded-md',
                activeTab === 'manual' ? 'bg-sky-700 text-white' : 'bg-slate-100 text-slate-600'
              ]"
            >
              {{ activeModuleCards.length }} Modul
            </span>
          </button>

          <!-- Tab 2: Solusi (Troubleshooting) -->
          <button
            @click="activeTab = 'solusi'"
            :class="[
              'inline-flex items-center gap-2 px-5 py-2.5 rounded-lg text-xs sm:text-sm font-bold transition-all',
              activeTab === 'solusi'
                ? 'bg-emerald-600 text-white shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            ]"
          >
            <Lightbulb class="w-4 h-4" />
            <span>Solusi</span>
            <span
              :class="[
                'text-[10px] font-bold px-2 py-0.5 rounded-md',
                activeTab === 'solusi' ? 'bg-emerald-700 text-white' : 'bg-slate-100 text-slate-600'
              ]"
            >
              {{ troubleshootingArticles.length }} Artikel
            </span>
          </button>
        </div>

        <!-- Explanatory Note below tab -->
        <p class="mt-2 text-xs text-slate-400 flex items-center justify-center gap-1">
          <Info class="w-3.5 h-3.5 inline" />
          <span>Tab 'Manual Book' aktif (tata cara per modul). Klik tab 'Solusi' untuk dialihkan ke Halaman Utama Solusi.</span>
        </p>
      </div>
    </section>

    <!-- 4. Main Section: Modul Cards Grid -->
    <main class="flex-1 w-full px-4 sm:px-6 lg:px-12 xl:px-16 pb-20">
      <!-- Section Header -->
      <div class="mb-6 flex flex-col sm:flex-row sm:items-end justify-between gap-2">
        <div>
          <h2 class="text-xl sm:text-2xl font-bold text-slate-900 tracking-tight">
            Pilih Modul Sesuai Kebutuhan Anda
          </h2>
          <p class="text-xs sm:text-sm text-slate-500 mt-1">
            Pelajari panduan tata cara penggunaan fitur {{ currentSystem.name }} langkah demi langkah per modul:
          </p>
        </div>
        <div v-if="isModulesLoaded" class="text-xs font-semibold px-3 py-1 rounded-full bg-slate-100 text-slate-600 w-fit">
          {{ filteredModuleCards.length }} Modul Tersedia
        </div>
        <div v-else class="h-6 w-28 bg-slate-100 rounded-full animate-pulse"></div>
      </div>

      <!-- TAB 1: Modul Cards Grid (Hanya yang memiliki buku panduan) -->
      <div v-if="activeTab === 'manual'">
        <!-- Loading Skeleton -->
        <div v-if="!isModulesLoaded" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5 lg:gap-6 xl:gap-8">
          <div v-for="i in 3" :key="i" class="bg-white rounded-2xl p-6 sm:p-7 border border-slate-200/90 flex flex-col justify-between h-[210px] animate-pulse shadow-xs">
            <div>
              <div class="flex items-center justify-between">
                <div class="w-12 h-12 rounded-xl bg-slate-100"></div>
                <div class="h-6 w-20 bg-slate-100 rounded-md"></div>
              </div>
              <div class="h-6 w-36 bg-slate-100 rounded mt-5"></div>
              <div class="h-4 w-full bg-slate-100 rounded mt-2.5"></div>
              <div class="h-4 w-2/3 bg-slate-100 rounded mt-2"></div>
            </div>
            <div class="pt-5 border-t border-slate-100 flex items-center justify-between">
              <div class="h-4 w-24 bg-slate-100 rounded"></div>
            </div>
          </div>
        </div>

        <div v-else-if="filteredModuleCards.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5 lg:gap-6 xl:gap-8">
          <div
            v-for="card in filteredModuleCards"
            :key="card.id"
            @click="openReader(card, 'manual')"
            class="bg-white rounded-2xl p-6 sm:p-7 border border-slate-200/90 hover:border-sky-500 hover:ring-4 hover:ring-sky-500/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200 cursor-pointer flex flex-col justify-between group shadow-xs"
          >
            <div>
              <!-- Card Top: Icon & Count Badge -->
              <div class="flex items-center justify-between">
                <div
                  class="rounded-xl flex items-center justify-center transition-all overflow-hidden"
                  :class="card.isImage ? 'w-12 h-12 bg-white border border-slate-200/90 p-0 shadow-2xs' : 'w-11 h-11 bg-slate-100 text-slate-700 group-hover:bg-sky-50 group-hover:text-sky-600'"
                >
                  <img
                    v-if="card.isImage"
                    :src="card.icon"
                    :alt="card.name"
                    class="w-full h-full object-cover rounded-xl transition-transform duration-200 group-hover:scale-105"
                  />
                  <component v-else :is="card.icon" class="w-5 h-5 stroke-[1.8]" />
                </div>

                <!-- Top Right Badge (Pill) - Menampilkan jumlah panduan riil -->
                <span
                  class="text-xs font-semibold px-2.5 py-1 rounded-md tracking-wide bg-slate-100 text-slate-600 group-hover:bg-sky-100/70 group-hover:text-sky-700 transition-colors"
                >
                  {{ card.countText }}
                </span>
              </div>

              <!-- Card Title & Description -->
              <h3
                class="font-bold text-slate-900 text-base sm:text-lg mt-5 group-hover:text-sky-600 transition-colors"
                v-html="highlightKeyword(card.title, searchQuery)"
              ></h3>
              <p
                class="text-xs sm:text-sm text-slate-500 mt-2 leading-relaxed line-clamp-3"
                v-html="highlightKeyword(card.desc, searchQuery)"
              ></p>
            </div>

            <!-- Card Bottom: Action Link -->
            <div class="pt-5 mt-5 border-t border-slate-100 flex items-center justify-end text-xs sm:text-sm">
              <span class="text-sky-600 font-semibold flex items-center gap-1 group-hover:text-sky-700">
                <span>Buka Manual Book</span>
                <ArrowRight class="w-3.5 h-3.5 group-hover:translate-x-1 transition-transform" />
              </span>
            </div>
          </div>
        </div>

        <!-- Empty state jika tidak ada modul atau tidak cocok dengan filter pencarian -->
        <div v-else class="py-16 px-4 text-center bg-white rounded-2xl border border-slate-200/80 shadow-xs">
          <div class="w-12 h-12 rounded-full bg-slate-100 flex items-center justify-center text-slate-400 mx-auto mb-3">
            <BookOpen class="w-6 h-6 stroke-[1.5]" />
          </div>
          <h4 class="text-base font-bold text-slate-800">
            {{ searchQuery ? 'Tidak Ada Modul yang Cocok' : 'Belum Ada Modul Panduan' }}
          </h4>
          <p class="text-xs sm:text-sm text-slate-500 mt-1 max-w-md mx-auto">
            {{ searchQuery 
              ? `Tidak ditemukan modul yang memuat kata kunci "${searchQuery}". Coba kata kunci lain.` 
              : `Modul buku panduan untuk sistem ${currentSystem.name} sedang dalam proses kurasi dan penyiapan.` 
            }}
          </p>
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            class="mt-4 px-4 py-1.5 rounded-lg bg-sky-50 text-sky-700 font-semibold text-xs hover:bg-sky-100 transition-colors cursor-pointer"
          >
            Reset Pencarian
          </button>
        </div>
      </div>

      <!-- TAB 2: Troubleshooting Articles List (Solusi Tab) -->
      <div v-else class="space-y-4">
        <!-- If articles exist -->
        <div v-if="filteredTroubleshooting.length > 0" class="space-y-4">
          <div
            v-for="art in filteredTroubleshooting"
            :key="art.id"
            @click="openReader(art, 'article')"
            class="bg-white p-6 rounded-2xl border border-slate-200 hover:border-emerald-400 hover:shadow-md transition-all cursor-pointer group flex flex-col sm:flex-row sm:items-center justify-between gap-4"
          >
            <div class="flex items-start gap-4">
              <div class="w-11 h-11 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0 mt-0.5 group-hover:bg-emerald-600 group-hover:text-white transition-colors">
                <Lightbulb class="w-5 h-5" />
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <span class="text-[11px] font-semibold px-2 py-0.5 rounded bg-emerald-50 text-emerald-700 border border-emerald-200">
                    {{ art.menuName }}
                  </span>
                  <span class="text-xs text-slate-400 flex items-center gap-1">
                    <Clock class="w-3 h-3" /> {{ art.readTime }}
                  </span>
                </div>
                <h3
                  class="text-base sm:text-lg font-bold text-slate-900 mt-1.5 group-hover:text-emerald-600 transition-colors"
                  v-html="highlightKeyword(art.title, searchQuery)"
                ></h3>
                <p
                  class="text-xs sm:text-sm text-slate-500 mt-1"
                  v-html="highlightKeyword(art.summary, searchQuery)"
                ></p>
              </div>
            </div>

            <div class="shrink-0 flex items-center gap-2 text-emerald-600 font-semibold text-xs sm:text-sm">
              <span>Baca Solusi Lengkap</span>
              <ArrowRight class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
            </div>
          </div>
        </div>

        <!-- Empty State for Solusi Tab (When system has no articles yet) -->
        <div v-else class="bg-white rounded-2xl border border-slate-200/90 p-8 sm:p-12 text-center shadow-xs">
          <div class="w-14 h-14 rounded-2xl bg-amber-50 text-amber-600 flex items-center justify-center mx-auto mb-4 border border-amber-200/60 shadow-2xs">
            <HelpCircle class="w-7 h-7 stroke-[1.8]" />
          </div>
          <h3 class="text-base sm:text-lg font-bold text-slate-900">
            Belum Ada Artikel Solusi untuk Sistem {{ currentSystem.name }}
          </h3>
          <p class="text-xs sm:text-sm text-slate-500 mt-1.5 max-w-lg mx-auto leading-relaxed">
            Artikel pemecahan masalah teknis untuk {{ currentSystem.title }} sedang dalam proses kurasi oleh tim teknis Probus. Anda dapat membaca Buku Panduan operasional modul.
          </p>
          <div class="mt-6 flex flex-wrap items-center justify-center gap-3">
            <button
              @click="activeTab = 'manual'"
              class="px-4 py-2.5 rounded-xl bg-sky-600 hover:bg-sky-700 text-white font-semibold text-xs sm:text-sm shadow-xs transition-colors flex items-center gap-2 cursor-pointer"
            >
              <BookOpen class="w-4 h-4" />
              <span>Buka Buku Panduan {{ currentSystem.name }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 5. Bottom Green Help Banner -->
      <div v-if="activeTab === 'manual'" class="mt-10 sm:mt-12 bg-emerald-50/70 border border-emerald-200/90 rounded-2xl p-5 sm:p-6 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 shadow-xs">
        <div class="flex items-start sm:items-center gap-4">
          <div class="w-10 h-10 rounded-xl bg-emerald-100 text-emerald-600 flex items-center justify-center shrink-0 shadow-2xs">
            <Lightbulb class="w-5 h-5 stroke-[2]" />
          </div>
          <div>
            <h4 class="text-sm sm:text-base text-slate-900 font-bold">
              Mencari Solusi Masalah Teknis atau Error {{ currentSystem.name }}?
            </h4>
            <p class="text-xs sm:text-sm text-slate-600 mt-0.5">
              Pilih modul di atas untuk tata cara penggunaan. Untuk artikel solusi & troubleshooting, buka Bank Solusi Klien.
            </p>
          </div>
        </div>

        <button
          @click="activeTab = 'solusi'"
          class="shrink-0 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs sm:text-sm px-5 py-2.5 rounded-xl flex items-center gap-2 shadow-xs transition-colors"
        >
          <span>Buka Halaman Solusi {{ currentSystem.name }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </div>
    </main>

    <!-- 6. Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 px-4 sm:px-6 lg:px-12 xl:px-16 text-xs sm:text-sm text-slate-500">
      <div class="w-full flex flex-col sm:flex-row items-center justify-between gap-3 text-center sm:text-left">
        <p class="font-normal">
          © 2026 CV Mitra Solusi Mandiri • Pusat Bantuan & Manual Book Resmi
        </p>
      </div>
    </footer>

    <!-- 7. Interactive Reader Modal -->
    <div
      v-if="readerModalOpen && activeReaderItem"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 backdrop-blur-xs p-4 sm:p-6"
    >
      <div class="bg-white w-full max-w-3xl max-h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-in fade-in duration-200">
        <div class="p-5 border-b border-slate-200 flex items-center justify-between bg-slate-50/50">
          <div class="flex items-center gap-2.5">
            <span
              :class="[
                'text-[11px] font-bold px-2.5 py-0.5 rounded-full uppercase tracking-wider',
                activeReaderItem.type === 'article'
                  ? 'bg-emerald-100 text-emerald-800'
                  : 'bg-sky-100 text-sky-800'
              ]"
            >
              {{ activeReaderItem.type === 'article' ? 'Artikel Solusi' : 'Buku Panduan' }}
            </span>
            <span class="text-xs text-slate-400">• {{ currentSystem.title }}</span>
          </div>
          <button
            @click="closeReader"
            class="w-8 h-8 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-500 hover:text-slate-800 flex items-center justify-center transition-colors"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-6 sm:p-8 overflow-y-auto space-y-4">
          <h2 class="text-xl sm:text-2xl font-black text-slate-900 leading-snug">
            {{ activeReaderItem.title }}
          </h2>

          <div
            class="manual-book-content text-slate-700 text-sm sm:text-base leading-relaxed pt-4 border-t border-slate-100"
            v-html="activeReaderItem.content"
          ></div>
        </div>

        <div class="p-4 sm:p-5 border-t border-slate-200 bg-slate-50 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs">
          <div class="text-slate-500">
            Apakah panduan langkah demi langkah ini membantu Anda?
          </div>
          <div class="flex items-center gap-2">
            <button
              v-if="!feedbackGiven"
              @click="giveFeedback"
              class="px-3.5 py-1.5 rounded-lg border border-slate-300 bg-white hover:bg-slate-100 text-slate-700 font-medium inline-flex items-center gap-1.5 transition-colors"
            >
              <ThumbsUp class="w-3.5 h-3.5" />
              <span>Ya, Membantu</span>
            </button>
            <div v-else class="text-emerald-600 font-bold inline-flex items-center gap-1">
              <CheckCircle2 class="w-4 h-4" />
              <span>Terima kasih atas masukan Anda!</span>
            </div>
            <button
              @click="closeReader"
              class="px-4 py-1.5 rounded-lg bg-sky-600 text-white font-semibold hover:bg-sky-700 transition-colors"
            >
              Tutup Panduan
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
