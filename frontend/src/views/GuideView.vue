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
  HelpCircle,
  MessageSquare,
  Share2,
  Check
} from 'lucide-vue-next'
import apiClient from '../api/client'
import { detectSystemIcon } from '../utils/iconResolver'
import { slugify, getSystemSlug, getSystemIdFromParam, findModuleByParam, highlightKeyword, registerCategories } from '../utils/slug'

const route = useRoute()
const router = useRouter()
const mobileMenuOpen = ref(false)

// Selected System ID (supports slug e.g. 'presto' and legacy ID e.g. '2')
const selectedSystemId = ref(getSystemIdFromParam(route.params.systemId, 2))

// Selected Module Order / ID
const selectedModuleOrder = ref(route.params.moduleId || 2)

// Active selected manual book index inside the active module (0-based)
const activeManualIndex = ref(0)

// Share button state
const copiedShareLink = ref(false)
let copyTimeout = null

// Sidebar search and accordion state
const tocSearchQuery = ref('')
const expandedModules = ref({})

const isLoading = ref(false)
const availableSystems = ref([
  { id: 1, name: 'PHIS', title: 'PHIS (1 Modul)', modulCount: 1, shortName: 'PHIS' },
  { id: 2, name: 'PRESTO', title: 'PRESTO (4 Modul)', modulCount: 4, shortName: 'PRESTO' },
  { id: 3, name: 'PPOS', title: 'PPOS (1 Modul)', modulCount: 1, shortName: 'PPOS' },
  { id: 4, name: 'SELF ORDER', title: 'SELF ORDER (2 Modul)', modulCount: 2, shortName: 'SELF ORDER' },
  { id: 5, name: 'PR-PO ONLINE', title: 'PR-PO ONLINE (1 Modul)', modulCount: 1, shortName: 'PR-PO ONLINE' },
  { id: 6, name: 'ORDER TAKER', title: 'ORDER TAKER (1 Modul)', modulCount: 1, shortName: 'ORDER TAKER' },
  { id: 7, name: 'VRESTROPICAL', title: 'VRESTROPICAL (0 Modul)', modulCount: 0, shortName: 'VRESTROPICAL' },
  { id: 8, name: 'PPE', title: 'PPE (0 Modul)', modulCount: 0, shortName: 'PPE' },
  { id: 9, name: 'BARONG POS', title: 'BARONG POS (0 Modul)', modulCount: 0, shortName: 'BARONG POS' },
])

const currentSystem = computed(() => {
  const found = availableSystems.value.find(s => s.id === selectedSystemId.value)
  if (found) return found
  return {
    id: selectedSystemId.value,
    name: `Sistem ${selectedSystemId.value}`,
    title: `Sistem ${selectedSystemId.value}`,
    modulCount: 0,
    shortName: `Sistem ${selectedSystemId.value}`,
  }
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

// Helper untuk mencocokkan sub menu aktif dari query URL (?sub=... atau ?chapter=...)
function resolveActiveSubMenuFromQuery(mod) {
  if (!mod || !mod.manualBooks || mod.manualBooks.length === 0) {
    activeManualIndex.value = 0
    return
  }

  const querySub = route.query.sub || route.query.chapter
  if (querySub !== undefined && querySub !== null && querySub !== '') {
    // 1. Cocokkan dengan index urutan 1-based (e.g. ?sub=2 -> index 1)
    const num = parseInt(querySub, 10)
    if (!isNaN(num) && num >= 1 && num <= mod.manualBooks.length) {
      activeManualIndex.value = num - 1
      return
    }

    // 2. Cocokkan dengan nomor urut modul (order) jika ada
    const orderMatchIdx = mod.manualBooks.findIndex(mb => String(mb.order) === String(querySub))
    if (orderMatchIdx >= 0) {
      activeManualIndex.value = orderMatchIdx
      return
    }

    // 3. Cocokkan dengan id sub menu
    const idMatchIdx = mod.manualBooks.findIndex(mb => String(mb.id) === String(querySub))
    if (idMatchIdx >= 0) {
      activeManualIndex.value = idMatchIdx
      return
    }

    // 4. Cocokkan dengan slug judul
    const slugMatchIdx = mod.manualBooks.findIndex(mb => slugify(cleanChapterTitle(mb.title)) === slugify(String(querySub)))
    if (slugMatchIdx >= 0) {
      activeManualIndex.value = slugMatchIdx
      return
    }
  }

  activeManualIndex.value = 0
}

// Menghasilkan full direct URL untuk di-share
function getShareUrl() {
  if (!currentModule.value) return window.location.href
  const sysSlug = getSystemSlug(selectedSystemId.value, currentSystem.value)
  const modSlug = slugify(currentModule.value.name)
  const resolved = router.resolve({
    path: `/system/${sysSlug}/guide/${modSlug}`,
    query: { sub: activeManualIndex.value + 1 }
  })
  return `${window.location.origin}${resolved.href}`
}

// Fungsi tombol share
async function shareCurrentGuide() {
  const url = getShareUrl()
  const title = currentManualBook.value
    ? cleanChapterTitle(currentManualBook.value.title)
    : (currentModule.value?.title || 'Panduan Probus')

  // Coba Web Share API di perangkat mobile
  if (navigator.share && /mobile|android|iphone|ipad/i.test(navigator.userAgent || '')) {
    try {
      await navigator.share({
        title: `${title} - ${currentSystem.value?.name || 'Probus'}`,
        text: `Buku Panduan: ${title} (${currentSystem.value?.name})`,
        url: url
      })
      return
    } catch (err) {
      if (err.name === 'AbortError') return
    }
  }

  // Fallback: Salin ke clipboard
  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(url)
    } else {
      const textArea = document.createElement('textarea')
      textArea.value = url
      textArea.style.position = 'fixed'
      textArea.style.opacity = '0'
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
    }

    copiedShareLink.value = true
    if (copyTimeout) clearTimeout(copyTimeout)
    copyTimeout = setTimeout(() => {
      copiedShareLink.value = false
    }, 2500)
  } catch (err) {
    console.error('Failed to copy share link:', err)
  }
}

// Fetch categories, system menus, and published manual books
async function loadSystemData() {
  isLoading.value = true
  try {
    // 1. Fetch Categories and all published Manual Books to determine real published module count
    const [catRes, allManualRes, detailRes, manualRes] = await Promise.all([
      apiClient.get('/categories').catch(err => {
        console.warn('Failed to fetch categories:', err)
        return { data: { data: [] } }
      }),
      apiClient.get('/manual-books').catch(err => {
        console.warn('Failed to fetch all manual books:', err)
        return { data: { data: [] } }
      }),
      apiClient.get(`/categories/${selectedSystemId.value}`).catch(err => {
        console.warn(`Failed to fetch category detail for ID ${selectedSystemId.value}:`, err)
        return { data: { data: null } }
      }),
      apiClient.get('/manual-books', { params: { kategori_id: selectedSystemId.value } }).catch(err => {
        console.warn(`Failed to fetch manual books for category ID ${selectedSystemId.value}:`, err)
        return { data: { data: [] } }
      })
    ])

    // Hitung jumlah modul yang benar-benar memiliki buku panduan terbit (published) per kategori
    const publishedModuleCounts = {}
    const allManuals = (allManualRes.data && allManualRes.data.data) ? allManualRes.data.data : []
    allManuals.forEach(mb => {
      if (mb.status === 'published' && Array.isArray(mb.menu_systems)) {
        mb.menu_systems.forEach(ms => {
          const kId = ms.kategori_id
          if (kId) {
            if (!publishedModuleCounts[kId]) {
              publishedModuleCounts[kId] = new Set()
            }
            publishedModuleCounts[kId].add(ms.id)
          }
        })
      }
    })

    if (catRes.data && catRes.data.data && Array.isArray(catRes.data.data)) {
      registerCategories(catRes.data.data)
      availableSystems.value = catRes.data.data.map(cat => {
        const count = publishedModuleCounts[cat.id] ? publishedModuleCounts[cat.id].size : 0
        return {
          id: cat.id,
          name: cat.name,
          title: `${cat.name} (${count} Modul)`,
          modulCount: count,
          shortName: cat.name,
        }
      })

      // Re-verify selectedSystemId with newly registered categories
      const reResolvedId = getSystemIdFromParam(route.params.systemId, selectedSystemId.value)
      if (reResolvedId && reResolvedId !== selectedSystemId.value) {
        selectedSystemId.value = reResolvedId
      }
    }

    // 2. Modul menus dari database untuk kategori yang sedang dipilih
    let menusFromDb = []
    if (detailRes.data && detailRes.data.data && detailRes.data.data.menus) {
      menusFromDb = detailRes.data.data.menus
    }

    // 3. Buku panduan terbit untuk kategori yang sedang dipilih (deduplikasi by id)
    const rawManuals = (manualRes.data && manualRes.data.data) ? manualRes.data.data : []
    const seenManualIds = new Set()
    const uniqueManuals = []
    const uniqueManualsMap = new Map()
    for (const mb of rawManuals) {
      if (mb && mb.id && !seenManualIds.has(mb.id)) {
        seenManualIds.add(mb.id)
        const formatted = {
          ...mb,
          description: isWordImportDesc(mb.description) ? '' : mb.description
        }
        uniqueManuals.push(formatted)
        uniqueManualsMap.set(mb.id, formatted)
      }
    }

    systemManualBooks.value = uniqueManuals

    // 4. Build module list combining database `menu_systems` and preloaded `manual_books`
    if (menusFromDb.length > 0) {
      const allModules = menusFromDb.map((menu, idx) => {
        const orderNum = menu.order || (idx + 1)
        const code = String(orderNum).padStart(2, '0')

        // Prioritaskan buku panduan yang tertaut langsung pada modul (menu.manual_books)
        // yang sudah diurutkan berdasarkan pivot order dari backend
        let uniqueLinked = []
        if (Array.isArray(menu.manual_books) && menu.manual_books.length > 0) {
          const seenIds = new Set()
          menu.manual_books.forEach(mb => {
            if (mb && mb.id && !seenIds.has(mb.id)) {
              seenIds.add(mb.id)
              const full = uniqueManualsMap.get(mb.id) || mb
              uniqueLinked.push({
                ...full,
                description: isWordImportDesc(full.description) ? '' : full.description
              })
            }
          })
        } else {
          // Fallback ke systemManualBooks jika menu.manual_books tidak tersedia
          const seenIds = new Set()
          systemManualBooks.value.forEach(mb => {
            if (mb && mb.id && !seenIds.has(mb.id) && Array.isArray(mb.menu_systems)) {
              if (mb.menu_systems.some(ms => ms.id === menu.id)) {
                seenIds.add(mb.id)
                uniqueLinked.push(mb)
              }
            }
          })
        }

        return {
          id: menu.id,
          order: orderNum,
          code: code,
          name: menu.name,
          title: `${code}. ${menu.name}`,
          desc: menu.description || 'Panduan operasional dan tata cara penggunaan fitur modul.',
          manualBooks: uniqueLinked,
          hasManualBook: uniqueLinked.length > 0,
        }
      })

      // Hanya tampilkan modul yang memiliki buku panduan di web client (modul kosong tidak ditampilkan)
      modulesList.value = allModules.filter(m => m.hasManualBook && m.manualBooks.length > 0)
    } else {
      modulesList.value = []
    }

    // Ensure active module is valid via slug or ID
    if (modulesList.value.length > 0) {
      const matched = findModuleByParam(modulesList.value, selectedModuleOrder.value || route.params.moduleId)
      if (matched) {
        selectedModuleOrder.value = matched.id
        expandedModules.value[matched.id] = true
        resolveActiveSubMenuFromQuery(matched)
      } else {
        selectedModuleOrder.value = modulesList.value[0].id
        expandedModules.value[modulesList.value[0].id] = true
        resolveActiveSubMenuFromQuery(modulesList.value[0])
      }

      // Synchronize canonical URL with sub-menu query parameter
      const active = modulesList.value.find(m => m.id === selectedModuleOrder.value) || modulesList.value[0]
      const sysSlug = getSystemSlug(selectedSystemId.value, currentSystem.value)
      const modSlug = slugify(active.name)
      const currentSub = active.manualBooks?.[activeManualIndex.value]
      const subParam = currentSub ? (currentSub.order || (activeManualIndex.value + 1)) : undefined

      if (route.params.systemId !== sysSlug || route.params.moduleId !== modSlug || (!route.query.sub && subParam)) {
        router.replace({
          path: `/system/${sysSlug}/guide/${modSlug}`,
          query: subParam ? { sub: subParam } : {}
        })
      }
    } else {
      selectedModuleOrder.value = null
    }
  } catch (err) {
    console.error('Failed to load system manual book data:', err)
  } finally {
    isLoading.value = false
  }
}

// Filtered modules by sidebar TOC search query (strictly within selected system)
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

// Total matching items in current system
const totalMatchesCount = computed(() => {
  const q = (tocSearchQuery.value || '').trim().toLowerCase()
  if (!q) return 0
  let count = 0
  filteredModulesList.value.forEach(mod => {
    const titleMatch = (mod.title || '').toLowerCase().includes(q) || (mod.name || '').toLowerCase().includes(q)
    const matchingSubs = (mod.manualBooks || []).filter(mb =>
      (mb.title || '').toLowerCase().includes(q) ||
      cleanChapterTitle(mb.title).toLowerCase().includes(q)
    ).length
    count += matchingSubs > 0 ? matchingSubs : (titleMatch ? 1 : 0)
  })
  return count
})

// Current active module
const currentModule = computed(() => {
  if (modulesList.value.length === 0) return null
  return modulesList.value.find(m => m.id === selectedModuleOrder.value || m.order === selectedModuleOrder.value) || modulesList.value[0] || null
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
  if (!currentModule.value || !currentModule.value.manualBooks || currentModule.value.manualBooks.length === 0) return null
  if (activeManualIndex.value > 0) {
    return currentModule.value.manualBooks[activeManualIndex.value - 1]
  }
  return null
})

const nextChapter = computed(() => {
  if (!currentModule.value || !currentModule.value.manualBooks || currentModule.value.manualBooks.length === 0) return null
  if (activeManualIndex.value < currentModule.value.manualBooks.length - 1) {
    return currentModule.value.manualBooks[activeManualIndex.value + 1]
  }
  return null
})

// Module Navigation in current system
const currentIndex = computed(() => {
  if (!currentModule.value) return -1
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
  const sysSlug = getSystemSlug(selectedSystemId.value, currentSystem.value)
  const modSlug = slugify(mod.name)
  router.replace({
    path: `/system/${sysSlug}/guide/${modSlug}`,
    query: { sub: chapterIdx + 1 }
  })
  window.scrollTo({ top: 120, behavior: 'smooth' })
}

function selectModule(mod) {
  selectedModuleOrder.value = mod.id
  activeManualIndex.value = 0
  expandedModules.value[mod.id] = true
  const sysSlug = getSystemSlug(selectedSystemId.value, currentSystem.value)
  const modSlug = slugify(mod.name)
  const firstSub = mod.manualBooks?.[0]
  const subParam = firstSub ? (firstSub.order || 1) : undefined
  router.replace({
    path: `/system/${sysSlug}/guide/${modSlug}`,
    query: subParam ? { sub: subParam } : {}
  })
  window.scrollTo({ top: 120, behavior: 'smooth' })
}

function switchSystem(sys) {
  selectedSystemId.value = sys.id
  tocSearchQuery.value = ''
  activeManualIndex.value = 0
  const sysSlug = getSystemSlug(sys.id, sys)
  router.replace(`/system/${sysSlug}/guide/1`)
  loadSystemData()
}

function goBackToSystemDetail() {
  router.push('/system/' + getSystemSlug(selectedSystemId.value, currentSystem.value))
}

function goToPrevChapter() {
  if (prevChapter.value && currentModule.value) {
    selectChapter(currentModule.value, activeManualIndex.value - 1)
  }
}

function goToNextChapter() {
  if (nextChapter.value && currentModule.value) {
    selectChapter(currentModule.value, activeManualIndex.value + 1)
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
        expandedModules.value[matched.id] = true
        resolveActiveSubMenuFromQuery(matched)
      }
    }
  }
)

watch(
  () => route.query.sub,
  newSub => {
    if (newSub !== undefined && currentModule.value) {
      resolveActiveSubMenuFromQuery(currentModule.value)
    }
  }
)
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 text-slate-800 font-sans flex flex-col antialiased">
    <!-- Main Navigation Bar with "Kembali ke Pilihan Modul" Button -->
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
        <a
          href="https://probussystem.com/contact"
          target="_blank"
          rel="noopener noreferrer"
          @click="mobileMenuOpen = false"
          class="block px-3 py-2 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-slate-900"
        >
          Hubungi Kami
        </a>
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
                {{ sys.modulCount ?? 0 }} Modul
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
            
            <!-- Sidebar Header with Selected System Badge -->
            <div class="mb-3 px-0.5">
              <div class="flex items-center justify-between gap-2">
                <div class="text-xs font-black text-slate-800 uppercase tracking-wider truncate">
                  DAFTAR MODUL &amp; SUB MENU
                </div>
                <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10.5px] font-extrabold bg-blue-50 text-blue-700 border border-blue-200/80 shrink-0">
                  {{ currentSystem.name }} ({{ modulesList.length }})
                </span>
              </div>
              <p class="text-[11px] text-slate-400 mt-1 truncate">
                Panduan materi operasional {{ currentSystem.name }}:
              </p>
            </div>

            <!-- Option 1: Modern Clean Search Filter (Scoped strictly to Selected System) -->
            <div class="relative mb-3.5 group">
              <!-- Search Icon with focus transition -->
              <Search class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-blue-600 pointer-events-none transition-colors duration-200" />
              
              <!-- Input with rounded-xl, subtle contrast, and smooth focus ring -->
              <input
                v-model="tocSearchQuery"
                type="text"
                :placeholder="`Cari materi ${currentSystem.name}...`"
                class="w-full pl-9 pr-24 py-2.5 bg-slate-50/90 hover:bg-slate-100/70 focus:bg-white border border-slate-200 hover:border-slate-300 focus:border-blue-500 rounded-xl text-xs text-slate-800 placeholder:text-slate-400 font-medium focus:outline-none focus:ring-4 focus:ring-blue-500/15 focus:shadow-xs transition-all duration-200"
              />

              <!-- Right Actions: Matches Count & Clear Button -->
              <div class="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1">
                <!-- Matches count badge when query is active -->
                <span
                  v-if="tocSearchQuery.trim() && totalMatchesCount > 0"
                  class="inline-flex items-center gap-0.5 px-2 py-0.5 rounded-full text-[10px] font-bold bg-blue-100 text-blue-800 border border-blue-200 shadow-2xs pointer-events-none animate-in fade-in"
                >
                  <CheckCircle2 class="w-2.5 h-2.5 text-blue-600" />
                  <span>{{ totalMatchesCount }} cocok</span>
                </span>
                <span
                  v-else-if="tocSearchQuery.trim() && totalMatchesCount === 0"
                  class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[10px] font-semibold bg-rose-50 text-rose-600 border border-rose-200 pointer-events-none"
                >
                  0 cocok
                </span>

                <!-- Clear button -->
                <button
                  v-if="tocSearchQuery"
                  @click="tocSearchQuery = ''"
                  type="button"
                  title="Hapus pencarian"
                  class="p-1 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-200/70 transition-colors focus:outline-none cursor-pointer"
                >
                  <X class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>

            <!-- Accordion List of Modules and Sub-Menus (Single Unified Card Layout) -->
            <div class="divide-y divide-slate-100 max-h-[calc(100vh-280px)] overflow-y-auto pr-1 custom-sky-scrollbar">
              <div
                v-for="mod in filteredModulesList"
                :key="mod.id"
                class="py-2 first:pt-0 last:pb-0"
              >
                <!-- Header Modul (Dapat diklik untuk expand / pilih modul) -->
                <div
                  class="p-2 sm:p-2.5 rounded-xl flex items-center justify-between gap-2.5 cursor-pointer select-none group transition-all duration-150"
                  :class="currentModule?.id === mod.id ? 'bg-blue-50/70' : 'hover:bg-slate-50'"
                  @click="toggleModule(mod)"
                >
                  <div class="flex items-center gap-2.5 min-w-0">
                    <!-- Ikon folder rounded berwarna biru di sebelah kiri -->
                    <div
                      :class="currentModule?.id === mod.id
                        ? 'bg-blue-600 text-white shadow-xs shadow-blue-500/20'
                        : 'bg-blue-50 text-blue-600 border border-blue-100/90 group-hover:bg-blue-100/80'"
                      class="w-8 h-8 rounded-xl flex items-center justify-center shrink-0 transition-colors"
                    >
                      <Folder class="w-4 h-4" />
                    </div>
                    <div class="min-w-0">
                      <!-- Judul modul (misal: "01. Setup Master") dengan teks tebal font-bold -->
                      <h3
                        class="text-xs sm:text-[13px] font-bold text-slate-900 group-hover:text-blue-600 transition-colors truncate leading-tight"
                        v-html="highlightKeyword(mod.title, tocSearchQuery)"
                      ></h3>
                      <!-- Teks status di bawah judul modul -->
                      <p class="text-[11px] mt-0.5 truncate leading-tight">
                        <span v-if="currentModule?.id === mod.id" class="text-blue-600 font-bold inline-flex items-center gap-1">
                          <span class="w-1.5 h-1.5 rounded-full bg-blue-600 animate-pulse"></span>
                          Sedang Dibaca
                        </span>
                        <span v-else class="text-slate-400 font-medium">
                          {{ (mod.manualBooks || []).length }} Sub Menu
                        </span>
                      </p>
                    </div>
                  </div>

                  <div class="flex items-center gap-2 shrink-0">
                    <!-- Badge jumlah sub menu berbentuk kapsul/pill abu-abu / biru muda -->
                    <span
                      class="text-[10.5px] font-bold px-2.5 py-0.5 rounded-full border transition-colors whitespace-nowrap"
                      :class="currentModule?.id === mod.id
                        ? 'bg-blue-100/70 text-blue-700 border-blue-200/80'
                        : 'bg-slate-100 text-slate-600 border-slate-200/70 group-hover:bg-blue-50 group-hover:text-blue-600 group-hover:border-blue-100'"
                    >
                      {{ (mod.manualBooks || []).length }} Sub Menu
                    </span>
                    <!-- Ikon chevron (panah atas/bawah) untuk accordion buka-tutup -->
                    <ChevronDown
                      class="w-4 h-4 text-slate-400 group-hover:text-slate-600 transition-transform duration-200 shrink-0"
                      :class="expandedModules[mod.id] ? 'rotate-180 text-blue-600' : ''"
                    />
                  </div>
                </div>

                <!-- Tree Hierarchy for Sub Menus (Sesuai Referensi Gambar 1) -->
                <div
                  v-show="expandedModules[mod.id] || (tocSearchQuery && tocSearchQuery.length > 0)"
                  class="mt-1.5 ml-4 sm:ml-5 pl-3 sm:pl-3.5 border-l-2 space-y-1 py-1"
                  :class="currentModule?.id === mod.id ? 'border-blue-300' : 'border-slate-200'"
                >
                  <div
                    v-for="(sub, sIdx) in mod.manualBooks"
                    :key="sub.id"
                    @click="selectChapter(mod, sIdx)"
                    :class="[
                      'p-2 rounded-xl transition-all duration-150 cursor-pointer flex items-center justify-between gap-2.5 select-none group/item active:scale-[0.99]',
                      currentModule?.id === mod.id && activeManualIndex === sIdx
                        ? 'bg-blue-600 text-white font-semibold shadow-xs'
                        : 'hover:bg-blue-50/70 text-slate-700 hover:text-blue-900'
                    ]"
                  >
                    <div class="flex items-center gap-2 min-w-0">
                      <!-- Badge nomor urut berbentuk pill di kiri (#1, #2) -->
                      <span
                        class="px-2 py-0.5 rounded-full text-[10.5px] shrink-0 transition-colors"
                        :class="currentModule?.id === mod.id && activeManualIndex === sIdx
                          ? 'bg-white/20 text-white border border-white/25 font-black'
                          : 'bg-slate-100 text-slate-500 border border-slate-200/70 group-hover/item:bg-blue-100 group-hover/item:text-blue-700 group-hover/item:border-blue-200 font-bold'"
                      >
                        #{{ sIdx + 1 }}
                      </span>
                      <!-- Judul sub menu -->
                      <span
                        class="text-xs truncate leading-normal"
                        :class="currentModule?.id === mod.id && activeManualIndex === sIdx
                          ? 'text-white font-semibold'
                          : 'text-slate-700 group-hover/item:text-blue-900 font-medium'"
                        v-html="highlightKeyword(cleanChapterTitle(sub.title), tocSearchQuery)"
                      ></span>
                    </div>

                    <!-- Ikon chevron kanan (">") di pojok kanan -->
                    <ChevronRight
                      v-if="currentModule?.id === mod.id && activeManualIndex === sIdx"
                      class="w-3.5 h-3.5 text-white shrink-0 ml-auto"
                    />
                    <ChevronRight
                      v-else
                      class="w-3.5 h-3.5 text-slate-300 opacity-0 group-hover/item:opacity-100 group-hover/item:text-blue-500 transition-opacity shrink-0 ml-auto"
                    />
                  </div>

                  <!-- Empty state jika belum ada sub menu panduan (seperti gambar 1) -->
                  <div
                    v-if="!mod.manualBooks || mod.manualBooks.length === 0"
                    class="text-xs text-slate-400 italic py-2 pl-2"
                  >
                    Belum ada sub menu panduan.
                  </div>
                </div>
              </div>

              <!-- Empty state jika belum ada modul dengan panduan sama sekali -->
              <div
                v-if="modulesList.length === 0"
                class="py-10 px-2 text-center text-xs text-slate-400"
              >
                <div class="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center text-slate-400 mx-auto mb-2.5">
                  <BookOpen class="w-5 h-5" />
                </div>
                <p class="font-bold text-slate-700">Belum Ada Panduan</p>
                <p class="text-[11px] text-slate-400 mt-1 leading-relaxed">
                  Modul panduan untuk sistem ini sedang disiapkan.
                </p>
              </div>

              <!-- Empty state pencarian sidebar -->
              <div
                v-else-if="filteredModulesList.length === 0"
                class="py-8 px-2 text-center text-xs text-slate-500"
              >
                <div class="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-400 mx-auto mb-2">
                  <Search class="w-4 h-4" />
                </div>
                <p class="font-bold text-slate-700">Tidak ada materi yang cocok</p>
                <p class="text-[11px] text-slate-400 mt-0.5">di sistem <strong>{{ currentSystem.name }}</strong></p>
                <button
                  @click="tocSearchQuery = ''"
                  class="mt-2.5 px-3 py-1 rounded-lg bg-blue-50 text-blue-700 text-[11px] font-bold hover:bg-blue-100 border border-blue-200 transition-colors cursor-pointer"
                >
                  Hapus Pencarian
                </button>
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
                
                <!-- Breadcrumb Hierarki Sub Menu + Share Button -->
                <div class="flex items-center justify-between gap-3 flex-wrap sm:flex-nowrap pb-1">
                  <!-- Breadcrumb Hierarki -->
                  <div class="flex items-center flex-wrap gap-1.5 text-xs font-semibold text-slate-400">
                    <span class="hover:text-slate-600 transition-colors uppercase tracking-wider font-bold text-sky-600">
                      {{ currentSystem.name }}
                    </span>
                    <ChevronRight class="w-3 h-3 text-slate-300" />
                    <span class="text-slate-600 font-medium">
                      {{ currentModule.title }}
                    </span>
                    <ChevronRight class="w-3 h-3 text-slate-300" />
                    <span class="text-blue-600 font-bold">
                      Sub Menu #{{ activeManualIndex + 1 }}
                    </span>
                  </div>

                  <!-- Tombol Bagikan / Share Button -->
                  <button
                    @click="shareCurrentGuide"
                    type="button"
                    :title="copiedShareLink ? 'Link panduan tersalin!' : 'Bagikan panduan ini'"
                    :class="[
                      'inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-semibold transition-all duration-200 border select-none shadow-2xs active:scale-95 cursor-pointer shrink-0',
                      copiedShareLink
                        ? 'bg-emerald-50 text-emerald-700 border-emerald-300 ring-2 ring-emerald-500/20'
                        : 'bg-white hover:bg-slate-50 text-slate-700 hover:text-blue-600 border-slate-200 hover:border-blue-200'
                    ]"
                  >
                    <Check v-if="copiedShareLink" class="w-3.5 h-3.5 text-emerald-600 animate-in zoom-in" />
                    <Share2 v-else class="w-3.5 h-3.5 text-slate-500 group-hover:text-blue-600" />
                    <span>{{ copiedShareLink ? 'Link Tersalin!' : 'Bagikan' }}</span>
                  </button>
                </div>

                <!-- Judul Sub Menu Besar (Clean Title) -->
                <div>
                  <h1 class="text-xl sm:text-2xl md:text-3xl font-black text-slate-900 tracking-tight leading-snug">
                    {{ cleanChapterTitle(currentManualBook.title) }}
                  </h1>

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

            <!-- CASE 2: Modul Belum Ada Panduan atau Sistem Kosong -->
            <template v-else>
              <div class="py-16 px-4 text-center">
                <div class="w-14 h-14 bg-blue-50 text-blue-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
                  <BookOpen class="w-7 h-7" />
                </div>
                <h3 class="text-base sm:text-lg font-bold text-slate-900">
                  Buku Panduan Belum Tersedia
                </h3>
                <p class="text-xs sm:text-sm text-slate-500 mt-2 max-w-md mx-auto leading-relaxed">
                  Buku panduan resmi untuk <strong v-if="currentModule">{{ currentModule.name }}</strong><span v-else>sistem <strong>{{ currentSystem.name }}</strong></span> sedang dalam proses kurasi dan publikasi oleh tim implementator Probus.
                </p>
                <div class="mt-6 flex flex-wrap items-center justify-center gap-3">
                  <a
                    href="https://wa.me/62811389211"
                    target="_blank"
                    class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs shadow-xs transition-colors"
                  >
                    <MessageSquare class="w-3.5 h-3.5 fill-white/20" />
                    <span>Tanyakan Penggunaan ke WhatsApp Support</span>
                  </a>
                  <button
                    @click="goBackToSystemDetail"
                    class="px-4 py-2.5 rounded-xl bg-white hover:bg-slate-50 text-slate-700 font-semibold text-xs border border-slate-200 transition-colors shadow-2xs cursor-pointer"
                  >
                    Kembali ke Pilihan Modul
                  </button>
                </div>
              </div>
            </template>
          </div>

          <!-- ======================================================== -->
          <!-- BOTTOM SUB MENU NAVIGATION (SEBELUMNYA & SELANJUTNYA)     -->
          <!-- Sesuai Desain manual_book/index.blade.php baris 952-1004 -->
          <!-- ======================================================== -->
          <div v-if="modulesList.length > 0 && currentManualBook" class="mt-12 pt-6 border-t border-slate-200 shrink-0">
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

/* Tabel TinyMCE & Word Importer (Border Penutup Keliling & Sel Penuh) */
:deep(.article-content table) {
  width: 100% !important;
  border-collapse: collapse !important;
  margin: 1.5rem 0 !important;
  font-size: 0.825rem !important;
  border: 1px solid #cbd5e1 !important;
}

:deep(.article-content table th),
:deep(.article-content table td) {
  border: 1px solid #cbd5e1 !important;
  padding: 0.625rem 0.75rem !important;
  text-align: left !important;
  vertical-align: middle !important;
  color: #334155 !important;
}

/* Garis Penutup Luar Tabel (Atas, Bawah, Kiri, Kanan Sesuai Kolom) */
:deep(.article-content table > thead > tr:first-child > th),
:deep(.article-content table > thead > tr:first-child > td),
:deep(.article-content table > tbody > tr:first-child > th),
:deep(.article-content table > tbody > tr:first-child > td),
:deep(.article-content table tr:first-child > th),
:deep(.article-content table tr:first-child > td) {
  border-top: 1px solid #cbd5e1 !important;
}

:deep(.article-content table > tbody > tr:last-child > th),
:deep(.article-content table > tbody > tr:last-child > td),
:deep(.article-content table tr:last-child > th),
:deep(.article-content table tr:last-child > td) {
  border-bottom: 1px solid #cbd5e1 !important;
}

:deep(.article-content table tr > th:first-child),
:deep(.article-content table tr > td:first-child) {
  border-left: 1px solid #cbd5e1 !important;
}

:deep(.article-content table tr > th:last-child),
:deep(.article-content table tr > td:last-child) {
  border-right: 1px solid #cbd5e1 !important;
}

:deep(.article-content table th),
:deep(.article-content table thead th),
:deep(.article-content table tbody tr:first-child th),
:deep(.article-content table tr:first-child th) {
  background-color: #f1f5f9 !important;
  color: #0f172a !important;
  font-weight: 700 !important;
}

/* Preserve dark navy header from TinyMCE Word importer if present */
:deep(.article-content table th[style*="background-color"]),
:deep(.article-content table td[style*="background-color"]) {
  color: inherit !important;
}

:deep(.article-content table th p),
:deep(.article-content table td p) {
  margin: 0 !important;
  padding: 0 !important;
}

/* Gambar Panduan Dalam Tabel (Icon / Tombol) */
:deep(.article-content table td img),
:deep(.article-content table th img) {
  max-height: 48px !important;
  width: auto !important;
  max-width: 100% !important;
  margin: 2px auto !important;
  display: inline-block !important;
  vertical-align: middle !important;
  border: none !important;
  box-shadow: none !important;
  border-radius: 4px !important;
}

/* Gambar Panduan Utama */
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
