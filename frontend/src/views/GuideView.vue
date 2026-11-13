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
  ImageIcon,
  CheckCircle2,
  Table,
  Info,
  Clock,
  Printer,
  Settings,
  UtensilsCrossed,
  Store,
  PieChart,
  ShoppingBag,
  Receipt,
  Wallet,
  BarChart3,
  TrendingUp,
  Sliders,
  ExternalLink,
  MessageSquare,
  HelpCircle,
  FileText
} from 'lucide-vue-next'
import apiClient from '../api/client'
import WhatsAppSupportDropdown from '../components/WhatsAppSupportDropdown.vue'

const route = useRoute()
const router = useRouter()

// Selected System ID (defaults to 2 for PRESTO)
const selectedSystemId = ref(Number(route.params.systemId) || 2)

// Selected Module Order / ID (defaults to 2 for Transaksi / Order Taker)
const selectedModuleOrder = ref(Number(route.params.moduleId) || 2)

// Active selected manual book inside the active module (if multiple exist)
const activeManualIndex = ref(0)

const isLoading = ref(false)
const availableSystems = ref([
  { id: 2, name: 'PRESTO', title: 'PRESTO Resto POS (10 Modul)', shortName: 'PRESTO Resto POS' },
  { id: 1, name: 'PHIS', title: 'Probus Hotel PMS (10 Modul)', shortName: 'Probus Hotel PMS' },
  { id: 8, name: 'PPE', title: 'PPE 2.0 Management (10 Modul)', shortName: 'PPE 2.0 Management' },
  { id: 3, name: 'PPOS', title: 'PPOS Point of Sales (10 Modul)', shortName: 'PPOS Point of Sales' },
  { id: 4, name: 'SELF ORDER', title: 'Probus Self Order (10 Modul)', shortName: 'Probus Self Order' },
  { id: 5, name: 'PR-PO ONLINE', title: 'PR-PO Online & Purchasing (10 Modul)', shortName: 'PR-PO Online' },
  { id: 6, name: 'ORDER TAKER', title: 'Order Taker Mobile (10 Modul)', shortName: 'Order Taker Mobile' },
  { id: 7, name: 'VRESTROPICAL', title: 'Vrestropical POS (10 Modul)', shortName: 'Vrestropical POS' },
  { id: 9, name: 'BARONG POS', title: 'Barong POS Cloud (8 Modul)', shortName: 'Barong POS Cloud' },
])

const currentSystem = computed(() => {
  return availableSystems.value.find(s => s.id === selectedSystemId.value) || availableSystems.value[0]
})

// Default 10 PRESTO modules metadata (fallback & icon mapping)
const defaultModuleMeta = {
  1: { code: '01', defaultTitle: '01. Setup Master', icon: Settings, lead: 'Konfigurasi master meja, denah resto, kategori menu, harga jual, dan tarif pajak/service.' },
  2: { code: '02', defaultTitle: '02. Order Taker & Transaksi Meja', icon: UtensilsCrossed, lead: 'Panduan kasir & pramusaji untuk membuka meja, input menu, catatan dapur, dan kirim order.' },
  3: { code: '03', defaultTitle: '03. Outlet Activity', icon: Store, lead: 'Monitoring operasional outlet: void pesanan, cancel item, move table, dan join bill.' },
  4: { code: '04', defaultTitle: '04. Cost Control', icon: PieChart, lead: 'Pengendalian HPP resep makanan, pemakaian bahan per porsi, dan waste management.' },
  5: { code: '05', defaultTitle: '05. Purcashing AP', icon: ShoppingBag, lead: 'Purchase Order (PO), penerimaan barang gudang (receiving), dan buku hutang supplier.' },
  6: { code: '06', defaultTitle: '06. AR (Account Receivable)', icon: Receipt, lead: 'Manajemen piutang tamu korporat, city ledger, faktur invoice, dan pelunasan piutang.' },
  7: { code: '07', defaultTitle: '07. General Cashier', icon: Wallet, lead: 'Sentralisasi setoran kasir outlet harian (drop cash), kas kecil petty cash, dan rekonsiliasi.' },
  8: { code: '08', defaultTitle: '08. Accounting System', icon: BarChart3, lead: 'Jurnal otomatis penjualan, buku besar, neraca saldo, dan laporan laba rugi restoran.' },
  9: { code: '09', defaultTitle: '09. Report', icon: TrendingUp, lead: 'Laporan analitik eksekutif: rekap penjualan kasir, ranking menu terlaris, dan pajak PB1.' },
  10: { code: '10', defaultTitle: '10. Tools', icon: Sliders, lead: 'Konfigurasi printer thermal LAN/USB, auto cash drawer, diagnostik perangkat keras.' },
}

// Current Modules list for the selected system
const modulesList = ref([])

// All published manual books for the current system from database
const systemManualBooks = ref([])

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
        shortName: cat.name === 'PRESTO' ? 'PRESTO Resto POS' : (cat.name === 'PHIS' ? 'Probus Hotel PMS' : `${cat.name} System`),
      }))
    }

    // 2. Fetch specific Category details to get its `menu_systems`
    const detailRes = await apiClient.get(`/categories/${selectedSystemId.value}`)
    let menusFromDb = []
    if (detailRes.data && detailRes.data.data && detailRes.data.data.menus) {
      menusFromDb = detailRes.data.data.menus
    }

    // 3. Fetch published manual books for this category from database (menu_has_manual_books)
    const manualRes = await apiClient.get('/manual-books', {
      params: { kategori_id: selectedSystemId.value }
    })
    systemManualBooks.value = (manualRes.data && manualRes.data.data) ? manualRes.data.data : []

    // 4. Build module list combining database `menu_systems` and preloaded `manual_books`
    if (menusFromDb.length > 0) {
      modulesList.value = menusFromDb.map((menu, idx) => {
        const orderNum = menu.order || (idx + 1)
        const meta = defaultModuleMeta[orderNum] || {
          code: String(orderNum).padStart(2, '0'),
          defaultTitle: `${String(orderNum).padStart(2, '0')}. ${menu.name}`,
          icon: BookOpen,
          lead: menu.description || 'Panduan operasional dan tata cara penggunaan fitur modul.',
        }

        // Find manual books linked to this menu_system_id
        const linkedManuals = systemManualBooks.value.filter(mb => {
          if (!mb.menu_systems || !Array.isArray(mb.menu_systems)) return false
          return mb.menu_systems.some(ms => ms.id === menu.id)
        })

        return {
          id: menu.id,
          order: orderNum,
          code: meta.code,
          name: menu.name,
          title: `${meta.code}. ${menu.name}`,
          desc: menu.description || meta.lead,
          icon: meta.icon,
          manualBooks: linkedManuals,
          hasManualBook: linkedManuals.length > 0,
        }
      })
    } else {
      // Fallback to default 10 PRESTO modules if database returns empty
      modulesList.value = Object.keys(defaultModuleMeta).map(key => {
        const num = Number(key)
        const meta = defaultModuleMeta[num]
        return {
          id: num,
          order: num,
          code: meta.code,
          name: meta.defaultTitle.replace(/^\d+\.\s*/, ''),
          title: meta.defaultTitle,
          desc: meta.lead,
          icon: meta.icon,
          manualBooks: [],
          hasManualBook: false,
        }
      })
    }

    // Ensure active module is valid
    const exists = modulesList.value.some(m => m.order === selectedModuleOrder.value || m.id === selectedModuleOrder.value)
    if (!exists && modulesList.value.length > 0) {
      selectedModuleOrder.value = modulesList.value[0].order || 1
    }
    const current = modulesList.value.find(m => m.order === selectedModuleOrder.value || m.id === selectedModuleOrder.value)
    if (current && current.hasManualBook) {
      openDropdownId.value = current.id
    }
  } catch (err) {
    console.error('Failed to load system manual book data:', err)
  } finally {
    isLoading.value = false
  }
}

// Current active module
const currentModule = computed(() => {
  return modulesList.value.find(m => m.order === selectedModuleOrder.value || m.id === selectedModuleOrder.value) || modulesList.value[0] || {
    id: 1,
    order: 1,
    code: '01',
    title: '01. Setup Master',
    desc: 'Panduan modul sistem.',
    manualBooks: [],
  }
})

// Current active manual book (if multiple exist for this module)
const currentManualBook = computed(() => {
  if (currentModule.value && currentModule.value.manualBooks && currentModule.value.manualBooks.length > 0) {
    const idx = Math.min(activeManualIndex.value, currentModule.value.manualBooks.length - 1)
    return currentModule.value.manualBooks[idx]
  }
  return null
})

// Navigation indices
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
  if (currentIndex.value < modulesList.value.length - 1) {
    return modulesList.value[currentIndex.value + 1]
  }
  return null
})

const openDropdownId = ref(null)

function toggleDropdown(modId) {
  if (openDropdownId.value === modId) {
    openDropdownId.value = null
  } else {
    openDropdownId.value = modId
  }
}

function selectModuleAndChapter(mod, mbIdx) {
  selectedModuleOrder.value = mod.order || mod.id
  activeManualIndex.value = mbIdx
  openDropdownId.value = mod.id
  router.replace({ params: { systemId: selectedSystemId.value, moduleId: mod.order || mod.id } })
  window.scrollTo({ top: 120, behavior: 'smooth' })
}

function selectModule(mod) {
  selectedModuleOrder.value = mod.order || mod.id
  activeManualIndex.value = 0
  if (mod.hasManualBook) {
    openDropdownId.value = mod.id
  } else {
    openDropdownId.value = null
  }
  router.replace({ params: { systemId: selectedSystemId.value, moduleId: mod.order || mod.id } })
  window.scrollTo({ top: 120, behavior: 'smooth' })
}

function switchSystem(sys) {
  selectedSystemId.value = sys.id
  selectedModuleOrder.value = 1
  activeManualIndex.value = 0
  router.replace({ params: { systemId: sys.id, moduleId: 1 } })
  loadSystemData()
}

function goBackToSystemDetail() {
  router.push('/system/' + selectedSystemId.value)
}

function goToPrev() {
  if (prevModule.value) {
    selectModule(prevModule.value)
  }
}

function goToNext() {
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
    if (newParams.systemId && Number(newParams.systemId) !== selectedSystemId.value) {
      selectedSystemId.value = Number(newParams.systemId)
      loadSystemData()
    }
    if (newParams.moduleId) {
      selectedModuleOrder.value = Number(newParams.moduleId)
      activeManualIndex.value = 0
    }
  }
)
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 text-slate-800 font-sans flex flex-col antialiased">
    
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

    <!-- 2. Main Navigation Bar with "Kembali ke Pilihan Modul" Button -->
    <nav class="border-b border-slate-200/80 bg-white/95 backdrop-blur-sm py-3.5 px-6 sm:px-12 lg:px-20 xl:px-28 sticky top-0 z-40 shadow-xs">
      <div class="max-w-[1600px] mx-auto flex items-center justify-between">
        <!-- Logo & Back Button -->
        <div class="flex items-center gap-3.5">
          <div class="flex items-center gap-2.5 cursor-pointer" @click="router.push('/')">
            <div class="w-8 h-8 rounded-lg bg-slate-950 flex items-center justify-center text-white font-extrabold text-base tracking-wider shadow-sm">
              P
            </div>
            <span class="font-black text-lg tracking-wider text-slate-950">PROBUS</span>
          </div>
          <span class="text-slate-300 font-light">|</span>
          <span class="text-slate-500 text-sm font-medium hidden sm:inline">Pusat Bantuan Publik</span>

          <!-- Pill Button: "← Kembali ke Pilihan Modul" -->
          <button
            @click="goBackToSystemDetail"
            class="ml-2 inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-sky-50 hover:bg-sky-100 text-sky-700 font-semibold text-xs border border-sky-100 transition-colors shadow-2xs cursor-pointer"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>Kembali ke Pilihan Modul</span>
          </button>
        </div>

        <!-- Navigation Links -->
        <div class="hidden lg:flex items-center gap-8 text-sm font-medium">
          <router-link to="/" class="text-sky-600 font-semibold hover:text-sky-700 transition-colors relative py-1">
            Panduan Sistem
            <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-sky-600 rounded-full"></span>
          </router-link>
          <router-link to="/solutions" class="text-slate-600 hover:text-slate-950 transition-colors">
            Bank Solusi
          </router-link>
          <a href="https://wa.me/6281910875514" target="_blank" class="text-slate-600 hover:text-slate-950 transition-colors">
            Hubungi Implementator
          </a>
        </div>
      </div>
    </nav>

    <!-- 3. Subheader: "Tata Cara Penggunaan Sistem" & System Pills Switcher -->
    <section class="border-b border-slate-200/80 bg-white py-6 sm:py-8 px-6 sm:px-12 lg:px-20 xl:px-28">
      <div class="max-w-[1600px] mx-auto">
        <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900 tracking-tight">
          Tata Cara Penggunaan Sistem
        </h1>
        <p class="text-xs sm:text-sm text-slate-500 mt-1">
          Pilih sistem di bawah, lalu pilih modul di menu kiri untuk mempelajari tata cara penggunaan langkah demi langkah.
        </p>

        <!-- "PILIH SISTEM" Label & Pills Row -->
        <div class="mt-4">
          <div class="text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2.5">
            PILIH SISTEM
          </div>
          <div class="flex items-center gap-2.5 overflow-x-auto pb-1 no-scrollbar">
            <button
              v-for="sys in availableSystems"
              :key="sys.id"
              @click="switchSystem(sys)"
              :class="[
                'shrink-0 text-xs px-4 py-2 rounded-lg font-semibold transition-all cursor-pointer',
                selectedSystemId === sys.id
                  ? 'bg-slate-900 text-white shadow-xs'
                  : 'bg-white text-slate-700 border border-slate-200 hover:bg-slate-50 hover:border-slate-300'
              ]"
            >
              {{ sys.title }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- 4. Two-Column Main Body: Left Sidebar (Daftar Modul) + Right Reader Content -->
    <main class="flex-1 max-w-[1600px] w-full mx-auto px-6 sm:px-12 lg:px-20 xl:px-28 py-8 sm:py-10">
      <div class="flex flex-col md:flex-row gap-8 items-start">
        
        <!-- LEFT COLUMN: Sidebar "DAFTAR MODUL" -->
        <aside class="w-full md:w-80 lg:w-[340px] shrink-0">
          <div class="bg-slate-100/70 border border-slate-200/80 rounded-2xl p-4 sm:p-5 sticky top-24">
            <div class="flex items-center justify-between text-xs font-bold text-slate-500 uppercase tracking-wider mb-3.5 px-1">
              <span>DAFTAR MODUL ({{ modulesList.length }})</span>
              <span class="text-[10px] font-normal text-slate-400 capitalize">
                {{ currentSystem.name }}
              </span>
            </div>

            <!-- List of Modules from Database `menu_systems` -->
            <div class="space-y-2.5 max-h-[calc(100vh-220px)] overflow-y-auto pr-1">
              <div
                v-for="mod in modulesList"
                :key="mod.id"
                @click="selectModule(mod)"
                :class="[
                  'rounded-xl p-4 border transition-all duration-150 cursor-pointer text-left',
                  currentModule.id === mod.id
                    ? 'bg-white border-2 border-sky-500 ring-2 ring-sky-100 shadow-sm'
                    : 'bg-white border-slate-200/90 hover:border-slate-300 hover:shadow-2xs'
                ]"
              >
                <div class="flex items-center justify-between gap-2">
                  <h3
                    :class="[
                      'text-sm font-bold tracking-tight',
                      currentModule.id === mod.id ? 'text-sky-700 font-extrabold' : 'text-slate-900'
                    ]"
                  >
                    {{ mod.title }}
                  </h3>

                  <!-- Indicator: Published Manual in DB -->
                  <span
                    v-if="mod.hasManualBook"
                    title="Buku panduan resmi telah dipublish"
                    class="w-2 h-2 rounded-full bg-emerald-500 shrink-0"
                  ></span>
                </div>

                <div class="flex items-center justify-between text-[11px] mt-1 font-medium">
                  <span
                    :class="currentModule.id === mod.id ? 'text-sky-600 font-semibold' : 'text-slate-400'"
                  >
                    <span v-if="currentModule.id === mod.id">Sedang Dibaca • Aktif</span>
                    <span v-else-if="mod.hasManualBook">
                      {{ mod.manualBooks.length > 1 ? `${mod.manualBooks.length} Bab Panduan` : 'Tersedia Panduan' }}
                    </span>
                    <span v-else>Panduan Modul</span>
                  </span>

                  <span v-if="mod.hasManualBook" class="text-[10px] text-emerald-600 font-bold px-1.5 py-0.5 rounded bg-emerald-50">
                    Published
                  </span>
                </div>

                <!-- Dropdown Bab Manual Book di Bagian Bawah Modul -->
                <div
                  v-if="mod.hasManualBook"
                  class="mt-2.5 pt-2 border-t border-slate-100"
                  @click.stop
                >
                  <button
                    type="button"
                    @click="toggleDropdown(mod.id)"
                    :class="[
                      'w-full flex items-center justify-between px-2.5 py-1.5 rounded-lg text-xs font-semibold transition-all cursor-pointer',
                      currentModule.id === mod.id
                        ? 'bg-sky-50 text-sky-900 border border-sky-200 hover:bg-sky-100/70'
                        : 'bg-slate-50 text-slate-700 border border-slate-200/80 hover:bg-slate-100'
                    ]"
                  >
                    <span class="flex items-center gap-1.5 truncate">
                      <BookOpen class="w-3.5 h-3.5 text-sky-600 shrink-0" />
                      <span class="truncate">
                        {{ currentModule.id === mod.id && mod.manualBooks[activeManualIndex] ? `Bab ${activeManualIndex + 1}: ${mod.manualBooks[activeManualIndex].title}` : `Pilih Bab (${mod.manualBooks.length} Bab)` }}
                      </span>
                    </span>
                    <ChevronDown
                      :class="[
                        'w-3.5 h-3.5 transition-transform duration-200 shrink-0 ml-1',
                        openDropdownId === mod.id ? 'rotate-180 text-sky-600' : 'text-slate-400'
                      ]"
                    />
                  </button>

                  <!-- Dropdown Menu / Daftar Bab -->
                  <div
                    v-if="openDropdownId === mod.id"
                    class="mt-1.5 p-1 bg-white rounded-xl shadow-md border border-slate-200 space-y-1 z-20 animate-in fade-in duration-150"
                  >
                    <button
                      v-for="(mb, mbIdx) in mod.manualBooks"
                      :key="mb.id"
                      type="button"
                      @click="selectModuleAndChapter(mod, mbIdx)"
                      :class="[
                        'w-full text-left px-2.5 py-1.5 rounded-lg text-xs font-medium transition-all flex items-start gap-2 cursor-pointer',
                        currentModule.id === mod.id && activeManualIndex === mbIdx
                          ? 'bg-sky-600 text-white font-bold shadow-xs'
                          : 'bg-white hover:bg-sky-50 text-slate-700 hover:text-slate-900'
                      ]"
                    >
                      <FileText
                        :class="[
                          'w-3.5 h-3.5 shrink-0 mt-0.5',
                          currentModule.id === mod.id && activeManualIndex === mbIdx ? 'text-white' : 'text-sky-600'
                        ]"
                      />
                      <span class="line-clamp-2 leading-snug">
                        Bab {{ mbIdx + 1 }}: {{ mb.title }}
                      </span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </aside>

        <!-- RIGHT COLUMN: Detailed Reader Content Canvas -->
        <article class="flex-1 min-w-0 bg-white rounded-2xl border border-slate-200/90 p-6 sm:p-10 lg:p-12 shadow-xs">
          
          <!-- Category Breadcrumb / Subtitle (Always Dynamically Matches Active System) -->
          <div class="text-xs font-bold text-sky-600 uppercase tracking-wider flex items-center gap-2">
            <span>{{ currentSystem.shortName }}</span>
            <span class="text-slate-300">•</span>
            <span>Modul {{ currentModule.code }} dari {{ String(modulesList.length).padStart(2, '0') }}</span>
          </div>

          <!-- Module Main Title -->
          <h2 class="text-2xl sm:text-3xl font-extrabold text-slate-900 mt-2 tracking-tight">
            {{ currentManualBook ? currentManualBook.title : currentModule.title }}
          </h2>

          <!-- Module Lead Text / Description -->
          <p class="mt-4 text-slate-600 text-sm sm:text-base leading-relaxed">
            {{ currentManualBook ? (currentManualBook.description || currentModule.desc) : currentModule.desc }}
          </p>

          <!-- Chapter / Section Tabs Switcher (If multiple manual books exist for this menu) -->
          <div
            v-if="currentModule.manualBooks && currentModule.manualBooks.length > 1"
            class="mt-6 flex items-center gap-2 p-1.5 bg-slate-50 border border-slate-200 rounded-xl overflow-x-auto"
          >
            <button
              v-for="(mb, mbIdx) in currentModule.manualBooks"
              :key="mb.id"
              @click="activeManualIndex = mbIdx"
              :class="[
                'px-4 py-2 rounded-lg text-xs font-bold whitespace-nowrap transition-all cursor-pointer flex items-center gap-1.5',
                activeManualIndex === mbIdx
                  ? 'bg-sky-600 text-white shadow-xs'
                  : 'bg-white text-slate-600 hover:text-slate-900 border border-slate-200/70'
              ]"
            >
              <FileText class="w-3.5 h-3.5" />
              <span>Bab {{ mbIdx + 1 }}: {{ mb.title }}</span>
            </button>
          </div>

          <hr class="my-8 border-slate-100" />

          <!-- CASE 1: Real Published Manual Book Content from Database `manual_books` -->
          <div v-if="currentManualBook && currentManualBook.content" class="manual-book-content">
            <!-- Render HTML Content directly from `dbprobuscs.manual_books.content` -->
            <div
              class="prose prose-slate max-w-none prose-headings:font-bold prose-headings:text-slate-900 prose-p:text-slate-600 prose-p:leading-relaxed prose-table:w-full prose-table:border-collapse prose-th:bg-slate-100 prose-th:p-3 prose-td:p-3 prose-img:rounded-xl prose-img:shadow-xs prose-img:border prose-img:border-slate-200"
              v-html="currentManualBook.content"
            ></div>
          </div>

          <!-- CASE 2: Draft / Not Yet Published in Database (Helpful Overview & Contact Support) -->
          <div v-else class="space-y-6">
            <!-- Notice Box -->
            <div class="bg-amber-50/70 border border-amber-200 rounded-xl p-5 sm:p-6 flex items-start gap-4 text-amber-900">
              <div class="w-10 h-10 rounded-lg bg-amber-100 text-amber-700 flex items-center justify-center shrink-0">
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
                    class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-amber-600 hover:bg-amber-700 text-white font-semibold text-xs shadow-xs transition-colors"
                  >
                    <MessageSquare class="w-3.5 h-3.5 fill-white/20" />
                    <span>Tanyakan Penggunaan ke WhatsApp Support</span>
                  </a>
                  <button
                    @click="goBackToSystemDetail"
                    class="px-4 py-2 rounded-lg bg-white hover:bg-amber-100/50 text-amber-900 font-semibold text-xs border border-amber-300 transition-colors"
                  >
                    Kembali ke Daftar Modul
                  </button>
                </div>
              </div>
            </div>

            <!-- Standard Operational Procedure Outline -->
            <div class="bg-slate-50 border border-slate-200 rounded-xl p-6">
              <h3 class="text-base font-bold text-slate-900 mb-2">
                Ringkasan Fungsi & Prosedur Modul
              </h3>
              <p class="text-sm text-slate-600 leading-relaxed">
                {{ currentModule.desc }}
              </p>
              <div class="mt-4 p-4 bg-white rounded-lg border border-slate-200 text-xs text-slate-500 space-y-2">
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

          <!-- Bottom Navigation Bar: "← Modul Sebelumnya" & "Modul Selanjutnya →" -->
          <div class="mt-12 pt-8 border-t border-slate-200 flex flex-col sm:flex-row items-center justify-between gap-4">
            <!-- Previous Button -->
            <button
              v-if="prevModule"
              @click="goToPrev"
              class="w-full sm:w-auto inline-flex items-center gap-2 text-xs sm:text-sm font-semibold text-slate-600 hover:text-slate-900 px-4 py-2.5 rounded-lg border border-slate-200 hover:bg-slate-50 transition-colors cursor-pointer"
            >
              <ArrowLeft class="w-4 h-4 text-slate-400" />
              <span>Modul Sebelumnya: {{ prevModule.title }}</span>
            </button>
            <div v-else class="hidden sm:block"></div>

            <!-- Next Button -->
            <button
              v-if="nextModule"
              @click="goToNext"
              class="w-full sm:w-auto inline-flex items-center gap-2 text-xs sm:text-sm font-bold text-sky-600 hover:text-sky-700 px-4 py-2.5 rounded-lg border border-sky-200 bg-sky-50/60 hover:bg-sky-50 transition-colors cursor-pointer shadow-2xs"
            >
              <span>Modul Selanjutnya: {{ nextModule.title }}</span>
              <ArrowRight class="w-4 h-4 text-sky-600" />
            </button>
            <div v-else class="hidden sm:block"></div>
          </div>
        </article>
      </div>
    </main>

    <!-- 5. Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 px-6 sm:px-12 lg:px-20 xl:px-28 text-xs sm:text-sm text-slate-500 mt-auto">
      <div class="max-w-[1600px] mx-auto flex flex-col sm:flex-row items-center justify-between gap-3">
        <p class="font-normal">© 2026 PT Probus System Indonesia</p>
        <p>
          Masih terkendala? Hubungi WhatsApp Support:
          <a href="https://wa.me/62811389211" target="_blank" class="font-semibold text-slate-700 hover:text-sky-600 transition-colors">
            +62 811-389-211
          </a>
        </p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* Ensure database-rendered HTML tables and images look clean */
:deep(.manual-book-content table) {
  width: 100% !important;
  border-collapse: collapse !important;
  margin: 1.25rem 0 !important;
  border: 1px solid #cbd5e1 !important;
  border-radius: 8px !important;
  overflow: hidden !important;
}

:deep(.manual-book-content th) {
  background-color: #f8fafc !important;
  color: #0f172a !important;
  font-weight: 700 !important;
  padding: 10px 14px !important;
  border: 1px solid #e2e8f0 !important;
  text-align: left !important;
  font-size: 0.875rem !important;
}

:deep(.manual-book-content td) {
  padding: 10px 14px !important;
  border: 1px solid #e2e8f0 !important;
  font-size: 0.875rem !important;
  color: #334155 !important;
  line-height: 1.5 !important;
}

:deep(.manual-book-content h3) {
  font-size: 1.125rem !important;
  font-weight: 700 !important;
  color: #0f172a !important;
  margin-top: 1.5rem !important;
  margin-bottom: 0.75rem !important;
}

:deep(.manual-book-content p) {
  font-size: 0.9375rem !important;
  color: #475569 !important;
  line-height: 1.625 !important;
  margin-bottom: 1rem !important;
}

:deep(.manual-book-content img) {
  max-width: 100% !important;
  height: auto !important;
  border-radius: 8px !important;
  margin: 1rem 0 !important;
  border: 1px solid #e2e8f0 !important;
}
</style>
