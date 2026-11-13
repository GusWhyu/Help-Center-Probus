<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  ArrowRight,
  Search,
  BookOpen,
  ChevronDown,
  ChevronRight,
  Folder,
  CheckCircle2,
  Phone,
  Menu,
  X,
  User,
  PenTool,
  HelpCircle,
  MessageSquare
} from 'lucide-vue-next'
import apiClient from '../api/client'
import { detectSystemIcon } from '../utils/iconResolver'
import { slugify, getSystemSlug, getSystemIdFromParam, findModuleByParam } from '../utils/slug'

const route = useRoute()
const router = useRouter()
const mobileMenuOpen = ref(false)

// Selected System ID (supports slug e.g. 'presto' and legacy ID e.g. '2')
const selectedSystemId = ref(getSystemIdFromParam(route.params.systemId, 2))

// Selected Module Order / ID
const selectedModuleOrder = ref(route.params.moduleId || 2)

// Active selected manual book index inside the active module (0-based)
const activeManualIndex = ref(0)

// Sidebar search and accordion state
const tocSearchQuery = ref('')
const expandedModules = ref({})

const isLoading = ref(false)
const availableSystems = ref([
  { id: 2, name: 'PRESTO', title: 'PRESTO (10 Modul)', modulCount: 10, shortName: 'PRESTO Resto POS' },
  { id: 1, name: 'PHIS', title: 'PHIS (10 Modul)', modulCount: 10, shortName: 'Probus Hotel PMS' },
  { id: 8, name: 'PPE', title: 'PPE (10 Modul)', modulCount: 10, shortName: 'PPE 2.0 Management' },
  { id: 3, name: 'PPOS', title: 'PPOS (10 Modul)', modulCount: 10, shortName: 'PPOS Point of Sales' },
  { id: 4, name: 'SELF ORDER', title: 'SELF ORDER (10 Modul)', modulCount: 10, shortName: 'Probus Self Order' },
  { id: 5, name: 'PR-PO ONLINE', title: 'PR-PO ONLINE (10 Modul)', modulCount: 10, shortName: 'PR-PO Online' },
  { id: 6, name: 'ORDER TAKER', title: 'ORDER TAKER (10 Modul)', modulCount: 10, shortName: 'Order Taker Mobile' },
  { id: 7, name: 'VRESTROPICAL', title: 'VRESTROPICAL (10 Modul)', modulCount: 10, shortName: 'Vrestropical POS' },
  { id: 9, name: 'BARONG POS', title: 'BARONG POS (8 Modul)', modulCount: 8, shortName: 'Barong POS Cloud' },
])

const currentSystem = computed(() => {
  return availableSystems.value.find(s => s.id === selectedSystemId.value) || availableSystems.value[0]
})

// Current Modules list for the selected system
const modulesList = ref([])

// All published manual books for the current system from database
const systemManualBooks = ref([])

// Helper: membersihkan nomor ganda pada judul bab (#1 01. Sales -> #1 Sales)
function cleanChapterTitle(title) {
  if (!title) return ''
  return title.replace(/^\d+[\.\s\-_]+/, '').trim() || title
}

// Helper: normalisasi gambar TinyMCE & Word Importer ke /uploads/...
function formatArticleContent(html) {
  if (!html) return ''
  return html
    .replace(/https?:\/\/(127\.0\.0\.1|localhost)(:\d+)?\/uploads\//g, '/uploads/')
    .replace(/\.\.\/uploads\//g, '/uploads/')
}

// Helper tanggal
function formatDate(dateStr) {
  if (!dateStr) return ''
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch {
    return dateStr
  }
}

function formatDateRelative(dateStr) {
  if (!dateStr) return ''
  try {
    const d = new Date(dateStr)
    const now = new Date()
    const diffMs = now - d
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
    if (diffDays === 0) return 'Hari ini'
    if (diffDays === 1) return 'Kemarin'
    if (diffDays < 30) return `${diffDays} hari yang lalu`
    return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch {
    return dateStr
  }
}

function isWordImportDesc(desc) {
  if (!desc) return true
  const lower = desc.toLowerCase()
  return lower.includes('diimpor dari file word') || lower.includes('.docx')
}

// Fetch categories, system menus, and published manual books
async function loadSystemData() {
  isLoading.value = true
  try {
    // 1. Fetch Categories
    const catRes = await apiClient.get('/categories')
    if (catRes.data && catRes.data.data && Array.isArray(catRes.data.data)) {
      availableSystems.value = catRes.data.data.map(cat => ({
        id: cat.id,
        name: cat.name,
        title: `${cat.name} (${cat.total_menus || 10} Modul)`,
        modulCount: cat.total_menus || (cat.name === 'BARONG POS' ? 8 : 10),
        shortName: cat.name === 'PRESTO' ? 'PRESTO Resto POS' : (cat.name === 'PHIS' ? 'Probus Hotel PMS' : `${cat.name} System`),
      }))
    }

    // 2. Fetch specific Category details to get its `menu_systems`
    const detailRes = await apiClient.get(`/categories/${selectedSystemId.value}`)
    let menusFromDb = []
    if (detailRes.data && detailRes.data.data && detailRes.data.data.menus) {
      menusFromDb = detailRes.data.data.menus
    }

    // 3. Fetch published manual books for this category from database
    const manualRes = await apiClient.get('/manual-books', {
      params: { kategori_id: selectedSystemId.value }
    })
    const rawManuals = (manualRes.data && manualRes.data.data) ? manualRes.data.data : []
    systemManualBooks.value = rawManuals.map(mb => ({
      ...mb,
      description: isWordImportDesc(mb.description) ? '' : mb.description
    }))

    // 4. Build module list combining database `menu_systems` and preloaded `manual_books`
    if (menusFromDb.length > 0) {
      const allModules = menusFromDb.map((menu, idx) => {
        const orderNum = menu.order || (idx + 1)
        const code = String(orderNum).padStart(2, '0')

        // Find manual books linked to this menu_system_id
        const linkedManuals = systemManualBooks.value.filter(mb => {
          if (!mb.menu_systems || !Array.isArray(mb.menu_systems)) return false
          return mb.menu_systems.some(ms => ms.id === menu.id)
        })

        // Sort sub-menus ascending by order
        linkedManuals.sort((a, b) => (a.order || 0) - (b.order || 0))

        return {
          id: menu.id,
          order: orderNum,
          code: code,
          name: menu.name,
          title: `${code}. ${menu.name}`,
          desc: menu.description || 'Panduan operasional dan tata cara penggunaan fitur modul.',
          manualBooks: linkedManuals,
          hasManualBook: linkedManuals.length > 0,
        }
      })

      // Hanya tampilkan modul yang sudah memiliki buku panduan (published), fallback ke allModules jika belum ada
      const withDocs = allModules.filter(m => m.hasManualBook)
      modulesList.value = withDocs.length > 0 ? withDocs : allModules
    } else {
      modulesList.value = []
    }

    // Ensure active module is valid via slug or ID
    const matched = findModuleByParam(modulesList.value, selectedModuleOrder.value || route.params.moduleId)
    if (matched) {
      selectedModuleOrder.value = matched.id
      expandedModules.value[matched.id] = true
    } else if (modulesList.value.length > 0) {
      selectedModuleOrder.value = modulesList.value[0].id
      expandedModules.value[modulesList.value[0].id] = true
    }

    // Synchronize canonical URL
    if (modulesList.value.length > 0) {
      const active = modulesList.value.find(m => m.id === selectedModuleOrder.value) || modulesList.value[0]
      const sysSlug = getSystemSlug(selectedSystemId.value)
      const modSlug = slugify(active.name)
      if (route.params.systemId !== sysSlug || route.params.moduleId !== modSlug) {
        router.replace(`/system/${sysSlug}/guide/${modSlug}`)
      }
    }
  } catch (err) {
    console.error('Failed to load system manual book data:', err)
  } finally {
    isLoading.value = false
  }
}

// Filtered modules by sidebar TOC search query
const filteredModulesList = computed(() => {
  const q = (tocSearchQuery.value || '').trim().toLowerCase()
  if (!q) return modulesList.value
  return modulesList.value.filter(mod => {
    const titleMatch = (mod.title || '').toLowerCase().includes(q)
    const nameMatch = (mod.name || '').toLowerCase().includes(q)
    const subMatch = (mod.manualBooks || []).some(mb =>
      (mb.title || '').toLowerCase().includes(q) ||
      cleanChapterTitle(mb.title).toLowerCase().includes(q)
    )
    return titleMatch || nameMatch || subMatch
  })
})

// Current active module
const currentModule = computed(() => {
  return modulesList.value.find(m => m.id === selectedModuleOrder.value || m.order === selectedModuleOrder.value) || modulesList.value[0] || {
    id: 1,
    order: 1,
    code: '01',
    title: '01. Setup Master',
    desc: 'Panduan modul sistem.',
    manualBooks: [],
  }
})

// Current active manual book (sub menu)
const currentManualBook = computed(() => {
  if (currentModule.value && currentModule.value.manualBooks && currentModule.value.manualBooks.length > 0) {
    const idx = Math.min(activeManualIndex.value, currentModule.value.manualBooks.length - 1)
    return currentModule.value.manualBooks[idx]
  }
  return null
})

// Sub-Menu Navigation inside current module
const prevChapter = computed(() => {
  if (!currentModule.value || !currentModule.value.manualBooks) return null
  if (activeManualIndex.value > 0) {
    return currentModule.value.manualBooks[activeManualIndex.value - 1]
  }
  return null
})

const nextChapter = computed(() => {
  if (!currentModule.value || !currentModule.value.manualBooks) return null
  if (activeManualIndex.value < currentModule.value.manualBooks.length - 1) {
    return currentModule.value.manualBooks[activeManualIndex.value + 1]
  }
  return null
})

// Module Navigation in current system
const currentIndex = computed(() => {
  return modulesList.value.findIndex(m => m.id === currentModule.value.id)
})

const prevModule = computed(() => {
  if (currentIndex.value > 0) {
    return modulesList.value[currentIndex.value - 1]
  }
  return null
})

const nextModule = computed(() => {
  if (currentIndex.value >= 0 && currentIndex.value < modulesList.value.length - 1) {
    return modulesList.value[currentIndex.value + 1]
  }
  return null
})

// User interactions
function toggleModule(mod) {
  const id = mod.id
  expandedModules.value[id] = !expandedModules.value[id]
  if (expandedModules.value[id] && selectedModuleOrder.value !== mod.id) {
    selectModule(mod)
  }
}

function selectChapter(mod, chapterIdx) {
  selectedModuleOrder.value = mod.id
  activeManualIndex.value = chapterIdx
  expandedModules.value[mod.id] = true
  const sysSlug = getSystemSlug(selectedSystemId.value)
  const modSlug = slugify(mod.name)
  router.replace(`/system/${sysSlug}/guide/${modSlug}`)
  window.scrollTo({ top: 120, behavior: 'smooth' })
}

function selectModule(mod) {
  selectedModuleOrder.value = mod.id
  activeManualIndex.value = 0
  expandedModules.value[mod.id] = true
  const sysSlug = getSystemSlug(selectedSystemId.value)
  const modSlug = slugify(mod.name)
  router.replace(`/system/${sysSlug}/guide/${modSlug}`)
  window.scrollTo({ top: 120, behavior: 'smooth' })
}

function switchSystem(sys) {
  selectedSystemId.value = sys.id
  activeManualIndex.value = 0
  const sysSlug = getSystemSlug(sys.id)
  router.replace(`/system/${sysSlug}/guide/1`)
  loadSystemData()
}

function goBackToSystemDetail() {
  router.push('/system/' + getSystemSlug(selectedSystemId.value))
}

function goToPrevChapter() {
  if (prevChapter.value) {
    activeManualIndex.value--
    window.scrollTo({ top: 120, behavior: 'smooth' })
  }
}

function goToNextChapter() {
  if (nextChapter.value) {
    activeManualIndex.value++
    window.scrollTo({ top: 120, behavior: 'smooth' })
  }
}

function goToPrevModule() {
  if (prevModule.value) {
    selectModule(prevModule.value)
  }
}

function goToNextModule() {
  if (nextModule.value) {
    selectModule(nextModule.value)
  }
}

onMounted(() => {
  loadSystemData()
})

watch(
  () => route.params,
  newParams => {
    if (newParams.systemId) {
      const resolvedSysId = getSystemIdFromParam(newParams.systemId, selectedSystemId.value)
      if (resolvedSysId !== selectedSystemId.value) {
        selectedSystemId.value = resolvedSysId
        loadSystemData()
        return
      }
    }
    if (newParams.moduleId && modulesList.value.length > 0) {
      const matched = findModuleByParam(modulesList.value, newParams.moduleId)
      if (matched && matched.id !== selectedModuleOrder.value) {
        selectedModuleOrder.value = matched.id
        activeManualIndex.value = 0
        expandedModules.value[matched.id] = true
      }
    }
  }
)
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 text-slate-800 font-sans flex flex-col antialiased">
    
    <!-- 1. Utility Top Bar -->
    <header class="border-b border-slate-200/70 bg-white text-xs text-slate-500 py-2.5 px-4 sm:px-6 lg:px-12 xl:px-16">
      <div class="w-full flex items-center justify-between sm:justify-end">
        <span class="text-slate-400 text-[11px] sm:hidden">Pusat Bantuan Resmi</span>
        <div class="flex items-center gap-2 text-slate-600">
          <span class="text-slate-400 hidden sm:inline">Butuh bantuan cepat?</span>
          <a
            href="tel:+62361256500"
            class="inline-flex items-center gap-1.5 font-bold text-slate-800 hover:text-blue-600 transition-colors"
          >
            <Phone class="w-3.5 h-3.5 text-slate-500" />
            <span>Phone: +62 361 256 500</span>
          </a>
        </div>
      </div>
    </header>

    <!-- 2. Main Navigation Bar with "Kembali ke Pilihan Modul" Button -->
    <nav class="border-b border-slate-200/90 bg-white py-3.5 px-4 sm:px-6 lg:px-12 xl:px-16 sticky top-0 z-40 shadow-2xs">
      <div class="w-full flex items-center justify-between">
        <!-- Logo & Back Button -->
        <div class="flex items-center gap-2.5 sm:gap-3.5">
          <div class="flex items-center cursor-pointer" @click="router.push('/')">
            <img
              src="/probus-logo.png"
              alt="Probus System"
              class="h-7 sm:h-9 w-auto object-contain hover:opacity-90 transition-opacity"
            />
          </div>
          <span class="text-slate-300 font-light hidden sm:inline">|</span>
          <span class="text-slate-500 text-sm font-medium hidden md:inline">Pusat Bantuan Publik</span>

          <!-- Pill Button: "← Kembali ke Pilihan Modul" -->
          <button
            @click="goBackToSystemDetail"
            class="inline-flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-lg bg-blue-50 hover:bg-blue-100 text-blue-700 font-semibold text-xs border border-blue-100 transition-colors shadow-2xs cursor-pointer"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Kembali ke Pilihan Modul</span>
            <span class="sm:hidden">Pilihan Modul</span>
          </button>
        </div>

        <!-- Desktop Navigation Links -->
        <div class="hidden md:flex items-center gap-8 text-sm font-medium">
          <router-link to="/" class="text-blue-600 font-semibold hover:text-blue-700 transition-colors relative py-1">
            Panduan Sistem
            <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-600 rounded-full"></span>
          </router-link>
          <router-link to="/solutions" class="text-slate-600 hover:text-slate-950 transition-colors">
            Bank Solusi
          </router-link>
          <router-link to="/contact" class="text-slate-600 hover:text-slate-950 transition-colors">
            Hubungi Implementator
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
          class="block px-3 py-2 rounded-lg text-sm font-semibold bg-blue-50 text-blue-700"
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
          class="block px-3 py-2 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-slate-900"
        >
          Hubungi Implementator
        </router-link>
      </div>
    </nav>

    <!-- 3. Subheader: "Tata Cara Penggunaan Sistem" & System Pills Switcher -->
    <section class="border-b border-slate-200/80 bg-white py-5 sm:py-6 px-4 sm:px-6 lg:px-12 xl:px-16">
      <div class="w-full">
        <h1 class="text-xl sm:text-2xl lg:text-3xl font-extrabold text-slate-900 tracking-tight">
          Tata Cara Penggunaan Sistem
        </h1>
        <p class="text-xs sm:text-sm text-slate-500 mt-1">
          Pilih sistem di bawah, lalu pilih modul di menu kiri untuk membaca panduan tata cara operasional langkah demi langkah.
        </p>

        <!-- "PILIH SISTEM" Label & Pills Row -->
        <div class="mt-4">
          <div class="flex items-center justify-between mb-2.5">
            <div class="flex items-center gap-2">
              <span class="w-1.5 h-1.5 rounded-full bg-blue-600"></span>
              <span class="text-[11px] font-bold text-slate-500 uppercase tracking-wider">
                Pilih Sistem Probus
              </span>
            </div>
            <span class="text-[11px] text-slate-400 hidden sm:inline">
              Geser untuk melihat sistem lainnya →
            </span>
          </div>

          <!-- Horizontal Scroll Container -->
          <div class="custom-sky-scrollbar flex items-center gap-2.5 overflow-x-auto pb-3 pt-1">
            <button
              v-for="sys in availableSystems"
              :key="sys.id"
              @click="switchSystem(sys)"
              :class="[
                'group shrink-0 inline-flex items-center gap-2 text-xs px-3.5 sm:px-4 py-2.5 rounded-xl font-bold transition-all duration-200 cursor-pointer select-none',
                selectedSystemId === sys.id
                  ? 'bg-blue-600 text-white shadow-sm shadow-blue-500/25 ring-2 ring-blue-500/30'
                  : 'bg-white text-slate-700 border border-slate-200/90 hover:border-blue-300 hover:bg-blue-50/60 hover:text-blue-700 shadow-2xs'
              ]"
            >
              <!-- System Icon -->
              <div
                :class="[
                  'w-6 h-6 rounded-lg flex items-center justify-center shrink-0 transition-colors',
                  selectedSystemId === sys.id
                    ? 'bg-white/20 text-white'
                    : 'bg-slate-100 text-slate-500 group-hover:bg-blue-100 group-hover:text-blue-600'
                ]"
              >
                <component :is="detectSystemIcon(sys.name)" class="w-3.5 h-3.5" />
              </div>

              <!-- System Name -->
              <span class="tracking-tight">{{ sys.name }}</span>

              <!-- Modul Count Badge -->
              <span
                :class="[
                  'text-[10px] font-semibold px-2 py-0.5 rounded-full transition-colors',
                  selectedSystemId === sys.id
                    ? 'bg-white/25 text-white'
                    : 'bg-slate-100 text-slate-500 group-hover:bg-blue-100/80 group-hover:text-blue-700'
                ]"
              >
                {{ sys.modulCount || (sys.name === 'BARONG POS' ? 8 : 10) }} Modul
              </span>
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- 4. Two-Column Main Body: Left Sidebar + Right Reader Content Canvas -->
    <main class="flex-1 w-full px-4 sm:px-6 lg:px-12 xl:px-16 py-6 sm:py-8 lg:py-10">
      <div class="flex flex-col lg:flex-row gap-6 lg:gap-8 items-start">
        
        <!-- ======================================================== -->
        <!-- LEFT COLUMN: SIDEBAR DAFTAR MODUL & SUB MENU ACCORDION   -->
        <!-- Sesuai Desain Probus Implementator manualbook.blade.php  -->
        <!-- ======================================================== -->
        <aside class="w-full lg:w-80 xl:w-[340px] shrink-0">
          <div class="bg-white border border-slate-200/90 rounded-2xl p-4 sm:p-5 shadow-2xs sticky top-24">
            
            <!-- Sidebar Header -->
            <div class="flex items-center justify-between mb-3 px-0.5">
              <div>
                <div class="text-xs font-black text-slate-800 uppercase tracking-wider">
                  DAFTAR MODUL &amp; SUB MENU ({{ modulesList.length }})
                </div>
                <p class="text-[11px] text-slate-400 mt-0.5">
                  Pilih modul untuk melihat materi sub menu:
                </p>
              </div>
            </div>

            <!-- Mini Search Filter for Sub Menus -->
            <div class="relative mb-3.5">
              <Search class="w-3.5 h-3.5 absolute left-3 top-2.5 text-slate-400" />
              <input
                v-model="tocSearchQuery"
                type="text"
                placeholder="Cari dalam daftar sub menu..."
                class="w-full pl-8 pr-7 py-1.5 bg-slate-50/70 focus:bg-white border border-slate-200 rounded-xl text-xs text-slate-700 placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-blue-500 transition-colors"
              />
              <button
                v-if="tocSearchQuery"
                @click="tocSearchQuery = ''"
                class="absolute right-2.5 top-2 text-slate-400 hover:text-slate-600 text-xs font-bold"
              >
                ✕
              </button>
            </div>

            <!-- Accordion List of Modules and Sub-Menus -->
            <div class="space-y-2.5 max-h-[calc(100vh-280px)] overflow-y-auto pr-1 custom-sky-scrollbar">
              <div
                v-for="mod in filteredModulesList"
                :key="mod.id"
                :class="[
                  'rounded-2xl border transition-all duration-200 overflow-hidden',
                  currentModule.id === mod.id
                    ? 'border-blue-500 bg-white shadow-2xs ring-1 ring-blue-500/30'
                    : 'border-slate-200/90 bg-white hover:border-slate-300'
                ]"
              >
                <!-- Header Modul (Dapat diklik untuk expand / pilih modul) -->
                <div
                  class="p-3.5 flex items-center justify-between gap-2 cursor-pointer select-none group"
                  @click="toggleModule(mod)"
                >
                  <div class="flex items-center gap-2.5 min-w-0">
                    <div
                      :class="currentModule.id === mod.id ? 'bg-blue-600 text-white shadow-xs' : 'bg-slate-50 text-slate-500 border border-slate-200'"
                      class="w-7 h-7 rounded-xl flex items-center justify-center shrink-0 transition-colors"
                    >
                      <Folder class="w-3.5 h-3.5" />
                    </div>
                    <div class="min-w-0">
                      <h3
                        class="text-xs sm:text-[13px] font-bold text-slate-800 truncate group-hover:text-blue-600 transition-colors"
                      >
                        {{ mod.title }}
                      </h3>
                      <span
                        class="text-[10px] font-medium"
                        :class="currentModule.id === mod.id ? 'text-blue-600 font-semibold' : 'text-slate-400'"
                      >
                        {{ currentModule.id === mod.id ? 'Sedang Dibaca' : `${mod.manualBooks.length} Sub Menu` }}
                      </span>
                    </div>
                  </div>

                  <div class="flex items-center gap-1.5 shrink-0">
                    <span
                      class="text-[10px] font-bold px-2 py-0.5 rounded-full"
                      :class="mod.manualBooks.length > 0 ? 'bg-blue-50 text-blue-600' : 'bg-slate-100 text-slate-400'"
                    >
                      {{ mod.manualBooks.length }} Sub Menu
                    </span>
                    <ChevronDown
                      class="w-3.5 h-3.5 text-slate-400 transition-transform duration-200"
                      :class="expandedModules[mod.id] ? 'rotate-180 text-blue-500' : ''"
                    />
                  </div>
                </div>

                <!-- Daftar Sub Menu di bawah Modul (Inline Accordion) -->
                <div
                  v-show="expandedModules[mod.id] || tocSearchQuery.length > 0"
                  class="border-t border-slate-100 bg-slate-50/70 p-2.5 space-y-1.5"
                >
                  <div
                    v-for="(sub, sIdx) in mod.manualBooks"
                    :key="sub.id"
                    @click="selectChapter(mod, sIdx)"
                    :class="[
                      'p-2.5 rounded-xl transition-all duration-150 cursor-pointer flex items-center justify-between gap-2 select-none group',
                      currentModule.id === mod.id && activeManualIndex === sIdx
                        ? 'bg-blue-600 text-white font-semibold shadow-xs'
                        : 'bg-white text-slate-700 hover:bg-blue-50/80 border border-slate-200/80'
                    ]"
                  >
                    <div class="flex items-center gap-2 min-w-0">
                      <span
                        class="w-5 h-5 rounded-md flex items-center justify-center text-[10px] font-black shrink-0"
                        :class="currentModule.id === mod.id && activeManualIndex === sIdx ? 'bg-white/20 text-white' : 'bg-slate-100 text-slate-500'"
                      >
                        #{{ sub.order || (sIdx + 1) }}
                      </span>
                      <span class="text-xs truncate font-medium">
                        {{ cleanChapterTitle(sub.title) }}
                      </span>
                    </div>
                    <ChevronRight
                      v-if="currentModule.id === mod.id && activeManualIndex === sIdx"
                      class="w-3 h-3 text-white/80 shrink-0"
                    />
                  </div>

                  <div
                    v-if="!mod.manualBooks || mod.manualBooks.length === 0"
                    class="text-xs text-slate-400 italic py-3 text-center"
                  >
                    Belum ada sub menu panduan.
                  </div>
                </div>
              </div>

              <!-- Empty state pencarian sidebar -->
              <div
                v-if="filteredModulesList.length === 0"
                class="py-8 text-center text-xs text-slate-400 italic"
              >
                Tidak ada sub menu yang cocok.
              </div>
            </div>
          </div>
        </aside>

        <!-- ======================================================== -->
        <!-- RIGHT COLUMN: DETAILED READER CONTENT CANVAS             -->
        <!-- Sesuai Desain manual_book/index.blade.php baris 893-1004 -->
        <!-- ======================================================== -->
        <article class="flex-1 min-w-0 bg-white rounded-3xl border border-slate-200/90 p-6 sm:p-10 shadow-xs flex flex-col justify-between">
          
          <div>
            <!-- CASE 1: Real Published Manual Book Content from Database `manual_books` -->
            <template v-if="currentManualBook && currentManualBook.content">
              <div class="space-y-4">
                
                <!-- Breadcrumb Hierarki Sub Menu -->
                <div class="flex items-center flex-wrap gap-1.5 text-xs font-semibold text-slate-400 pb-1">
                  <span class="hover:text-slate-600 transition-colors uppercase tracking-wider font-bold text-sky-600">
                    {{ currentSystem.name }}
                  </span>
                  <ChevronRight class="w-3 h-3 text-slate-300" />
                  <span class="text-slate-600 font-medium">
                    {{ currentModule.title }}
                  </span>
                  <ChevronRight class="w-3 h-3 text-slate-300" />
                  <span class="text-blue-600 font-bold">
                    Sub Menu #{{ currentManualBook.order || (activeManualIndex + 1) }}
                  </span>
                </div>

                <!-- Judul Sub Menu Besar (Clean Title) -->
                <div>
                  <h1 class="text-xl sm:text-2xl md:text-3xl font-black text-slate-900 tracking-tight leading-snug">
                    {{ cleanChapterTitle(currentManualBook.title) }}
                  </h1>

                  <!-- Informasi Pembuat & Pengedit Sub Menu -->
                  <div class="flex items-center gap-2 sm:gap-3 text-[11px] sm:text-xs text-slate-500 mt-2.5 flex-wrap">
                    <!-- Pembuat Awal -->
                    <div class="inline-flex items-center gap-1.5 bg-slate-50 border border-slate-200/60 px-2.5 py-1 rounded-lg">
                      <User class="w-3 h-3 text-blue-500" />
                      <span>Dibuat: <strong class="text-slate-700">{{ currentManualBook.created_by_name || 'Superadministrator' }}</strong></span>
                      <span class="text-slate-400 text-[10px]" v-if="formatDate(currentManualBook.created_at)">
                        ({{ formatDate(currentManualBook.created_at) }})
                      </span>
                    </div>

                    <!-- Pengedit Terakhir -->
                    <div
                      v-if="currentManualBook.updated_at"
                      class="inline-flex items-center gap-1.5 bg-slate-50 border border-slate-200/60 px-2.5 py-1 rounded-lg"
                    >
                      <PenTool class="w-3 h-3 text-amber-500" />
                      <span>Diedit: <strong class="text-slate-700">{{ currentManualBook.updated_by_name || 'Superadministrator' }}</strong></span>
                      <span class="text-slate-400 text-[10px]">
                        ({{ formatDateRelative(currentManualBook.updated_at) }})
                      </span>
                    </div>
                  </div>

                  <!-- Deskripsi Sub Menu (Ringkasan Singkat) -->
                  <p
                    v-if="currentManualBook.description && !isWordImportDesc(currentManualBook.description)"
                    class="text-xs sm:text-sm text-slate-600 italic mt-3 bg-slate-50 p-3.5 rounded-2xl border border-slate-200/80 leading-relaxed"
                  >
                    {{ currentManualBook.description }}
                  </p>
                </div>

                <hr class="border-slate-100 my-4" />

                <!-- Isi Materi Artikel Panduan (TinyMCE Content) -->
                <div
                  class="article-content max-w-none text-slate-800 text-xs sm:text-sm leading-relaxed overflow-x-auto min-h-[220px]"
                  v-html="formatArticleContent(currentManualBook.content)"
                ></div>
              </div>
            </template>

            <!-- CASE 2: Draft / Not Yet Published in Database -->
            <template v-else>
              <div class="space-y-6">
                <!-- Notice Box -->
                <div class="bg-amber-50/70 border border-amber-200 rounded-2xl p-5 sm:p-6 flex items-start gap-4 text-amber-900">
                  <div class="w-10 h-10 rounded-xl bg-amber-100 text-amber-700 flex items-center justify-center shrink-0">
                    <HelpCircle class="w-5 h-5 stroke-[2]" />
                  </div>
                  <div class="flex-1">
                    <h4 class="text-sm font-bold text-amber-950">
                      Dokumen Panduan Modul Sedang Disusun
                    </h4>
                    <p class="text-xs sm:text-sm text-amber-800/90 mt-1 leading-relaxed">
                      Buku panduan interaktif resmi untuk modul <strong>{{ currentModule.name }}</strong> pada sistem <strong>{{ currentSystem.name }}</strong> sedang dalam proses kurasi dan publikasi oleh tim implementator Probus.
                    </p>
                    <div class="mt-4 flex flex-wrap items-center gap-3">
                      <a
                        href="https://wa.me/62811389211"
                        target="_blank"
                        class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-amber-600 hover:bg-amber-700 text-white font-semibold text-xs shadow-xs transition-colors"
                      >
                        <MessageSquare class="w-3.5 h-3.5 fill-white/20" />
                        <span>Tanyakan Penggunaan ke WhatsApp Support</span>
                      </a>
                      <button
                        @click="goBackToSystemDetail"
                        class="px-4 py-2 rounded-xl bg-white hover:bg-amber-100/50 text-amber-900 font-semibold text-xs border border-amber-300 transition-colors"
                      >
                        Kembali ke Daftar Modul
                      </button>
                    </div>
                  </div>
                </div>

                <!-- Standard Operational Procedure Outline -->
                <div class="bg-slate-50 border border-slate-200 rounded-2xl p-6">
                  <h3 class="text-base font-bold text-slate-900 mb-2">
                    Ringkasan Fungsi & Prosedur Modul
                  </h3>
                  <p class="text-sm text-slate-600 leading-relaxed">
                    {{ currentModule.desc }}
                  </p>
                  <div class="mt-4 p-4 bg-white rounded-xl border border-slate-200 text-xs text-slate-500 space-y-2">
                    <div class="flex items-center gap-2">
                      <CheckCircle2 class="w-4 h-4 text-emerald-500 shrink-0" />
                      <span>Pastikan akun Anda memiliki hak akses (role permission) untuk membuka modul ini.</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <CheckCircle2 class="w-4 h-4 text-emerald-500 shrink-0" />
                      <span>Untuk bantuan konfigurasi awal, hubungi tim implementator resmi Probus System.</span>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </div>

          <!-- ======================================================== -->
          <!-- BOTTOM SUB MENU NAVIGATION (SEBELUMNYA & SELANJUTNYA)     -->
          <!-- Sesuai Desain manual_book/index.blade.php baris 952-1004 -->
          <!-- ======================================================== -->
          <div class="mt-12 pt-6 border-t border-slate-200 shrink-0">
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
              
              <!-- Tombol Sub Menu Sebelumnya -->
              <button
                v-if="prevChapter"
                type="button"
                @click="goToPrevChapter"
                class="flex-1 group flex items-center gap-3 p-3.5 rounded-2xl border border-slate-200 hover:border-blue-500/50 bg-slate-50/60 hover:bg-blue-50/40 text-left transition-all cursor-pointer"
              >
                <div class="w-9 h-9 rounded-xl bg-white border border-slate-200 flex items-center justify-center text-slate-500 group-hover:text-blue-600 group-hover:border-blue-200 transition-colors shrink-0">
                  <ArrowLeft class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
                </div>
                <div class="truncate">
                  <div class="text-[10px] font-extrabold uppercase tracking-wider text-slate-400 group-hover:text-blue-600 transition-colors">
                    ← Sub Menu Sebelumnya
                  </div>
                  <div class="text-xs sm:text-sm font-bold text-slate-800 group-hover:text-blue-600 truncate mt-0.5">
                    {{ cleanChapterTitle(prevChapter.title) }}
                  </div>
                </div>
              </button>
              
              <!-- Jika bab pertama di modul ini, tampilkan tombol kembali ke modul sebelumnya jika ada -->
              <button
                v-else-if="prevModule"
                type="button"
                @click="goToPrevModule"
                class="flex-1 group flex items-center gap-2.5 p-3.5 rounded-2xl border border-slate-200 hover:border-slate-300 bg-slate-50/60 hover:bg-slate-100/60 text-left transition-all cursor-pointer"
              >
                <div class="w-8 h-8 rounded-xl bg-white border border-slate-200 flex items-center justify-center text-slate-500 shrink-0">
                  <ArrowLeft class="w-3.5 h-3.5 group-hover:-translate-x-0.5 transition-transform" />
                </div>
                <div class="truncate max-w-[200px]">
                  <div class="text-[10px] font-extrabold uppercase tracking-wider text-slate-400">
                    ← Modul Sebelumnya
                  </div>
                  <div class="text-xs font-bold text-slate-700 truncate mt-0.5">
                    {{ prevModule.title }}
                  </div>
                </div>
              </button>
              
              <div v-else class="flex-1 hidden sm:block"></div>

              <!-- Tombol Sub Menu Selanjutnya -->
              <button
                v-if="nextChapter"
                type="button"
                @click="goToNextChapter"
                class="flex-1 group flex items-center justify-end gap-3 p-3.5 rounded-2xl border border-slate-200 hover:border-blue-500/50 bg-slate-50/60 hover:bg-blue-50/40 text-right transition-all cursor-pointer"
              >
                <div class="truncate">
                  <div class="text-[10px] font-extrabold uppercase tracking-wider text-slate-400 group-hover:text-blue-600 transition-colors">
                    Sub Menu Selanjutnya →
                  </div>
                  <div class="text-xs sm:text-sm font-bold text-slate-800 group-hover:text-blue-600 truncate mt-0.5">
                    {{ cleanChapterTitle(nextChapter.title) }}
                  </div>
                </div>
                <div class="w-9 h-9 rounded-xl bg-white border border-slate-200 flex items-center justify-center text-slate-500 group-hover:text-blue-600 group-hover:border-blue-200 transition-colors shrink-0">
                  <ArrowRight class="w-4 h-4 group-hover:translate-x-0.5 transition-transform" />
                </div>
              </button>

              <!-- Jika sudah di akhir modul, tampilkan banner selesai & tombol lanjut ke modul berikutnya -->
              <div v-else-if="nextModule" class="flex-1 flex flex-col sm:flex-row items-stretch sm:items-center gap-2">
                <div class="flex-1 p-3 rounded-2xl border border-emerald-500/20 bg-emerald-50/50 text-center flex items-center justify-center gap-2">
                  <CheckCircle2 class="w-4 h-4 text-emerald-500 shrink-0" />
                  <span class="text-xs font-bold text-emerald-700">Selesai materi modul ini!</span>
                </div>
                <button
                  type="button"
                  @click="goToNextModule"
                  class="group flex items-center justify-end gap-2.5 p-3 rounded-2xl border border-blue-200 bg-blue-50/70 hover:bg-blue-100 text-right transition-all cursor-pointer shrink-0"
                >
                  <div class="truncate max-w-[150px]">
                    <div class="text-[10px] font-extrabold uppercase tracking-wider text-blue-500">
                      Modul Selanjutnya →
                    </div>
                    <div class="text-xs font-bold text-blue-700 truncate mt-0.5">
                      {{ nextModule.title }}
                    </div>
                  </div>
                  <div class="w-8 h-8 rounded-xl bg-blue-600 text-white flex items-center justify-center shrink-0 shadow-xs">
                    <ArrowRight class="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" />
                  </div>
                </button>
              </div>

              <!-- Selesai seluruh modul pada sistem ini -->
              <div v-else class="flex-1 p-3.5 rounded-2xl border border-emerald-500/20 bg-emerald-50/50 text-center flex items-center justify-center gap-2">
                <CheckCircle2 class="w-4 h-4 text-emerald-500" />
                <span class="text-xs font-bold text-emerald-700">Selesai membaca seluruh materi pada sistem ini!</span>
              </div>

            </div>
          </div>

        </article>
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

<style scoped>
/* Styling Konten Artikel TinyMCE identik dengan Probus Implementator */
:deep(.article-content) {
  line-height: 1.8;
  color: #334155;
}

:deep(.article-content > :first-child) {
  margin-top: 0 !important;
  padding-top: 0 !important;
}

:deep(.article-content h1),
:deep(.article-content h2),
:deep(.article-content h3),
:deep(.article-content h4) {
  font-weight: 800;
  color: #0f172a;
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
  line-height: 1.3;
}

:deep(.article-content h1) { font-size: 1.5rem; }
:deep(.article-content h2) { font-size: 1.25rem; }
:deep(.article-content h3) { font-size: 1.125rem; }
:deep(.article-content h4) { font-size: 1rem; }

:deep(.article-content p) {
  margin-bottom: 1rem;
  color: #334155;
}

:deep(.article-content ul),
:deep(.article-content ol) {
  margin: 1rem 0;
  padding-left: 1.5rem;
}
:deep(.article-content ul) { list-style-type: disc; }
:deep(.article-content ol) { list-style-type: decimal; }
:deep(.article-content li) {
  margin-bottom: 0.35rem;
  color: #334155;
}

/* Tabel TinyMCE & Word Importer (Identik Probus Implementator) */
:deep(.article-content table) {
  width: 100% !important;
  border-collapse: collapse !important;
  margin: 1.25rem 0 !important;
  font-size: 0.825rem !important;
  border: 1px solid #cbd5e1 !important;
  border-radius: 0.5rem !important;
  overflow: hidden !important;
}

:deep(.article-content table th),
:deep(.article-content table td),
:deep(.article-content table tr th),
:deep(.article-content table tr td),
:deep(.article-content table tbody tr th),
:deep(.article-content table tbody tr td),
:deep(.article-content table thead tr th) {
  border: 1px solid #cbd5e1 !important;
  border-width: 1px !important;
  border-style: solid !important;
  border-color: #cbd5e1 !important;
  padding: 0.625rem 0.75rem !important;
  text-align: left !important;
  vertical-align: top !important;
}

:deep(.article-content table th),
:deep(.article-content table thead th),
:deep(.article-content table tbody tr:first-child th),
:deep(.article-content table tr:first-child th) {
  background-color: #f1f5f9 !important;
  color: #0f172a !important;
  font-weight: 700 !important;
}

/* Preserve dark navy header from TinyMCE Word importer (bg #0b1a36) */
:deep(.article-content table th[style*="background-color"]),
:deep(.article-content table td[style*="background-color"]) {
  color: inherit !important;
}

:deep(.article-content table th p),
:deep(.article-content table td p) {
  margin: 0 !important;
  padding: 0 !important;
}

:deep(.article-content table tbody tr) {
  border-bottom: 1px solid #cbd5e1 !important;
}

/* Gambar Panduan */
:deep(.article-content img) {
  max-width: 100% !important;
  height: auto !important;
  border-radius: 0.75rem !important;
  margin: 1.25rem auto !important;
  border: 1px solid #e2e8f0 !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05) !important;
  display: block !important;
}
</style>
