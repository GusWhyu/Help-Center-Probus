<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Search,
  Globe,
  Building2,
  Zap,
  BarChart3,
  Package,
  ArrowRight,
  Lightbulb,
  MessageSquare,
  MonitorCheck,
  BookOpen,
  AlertCircle,
  X,
  Printer,
  Smartphone,
  UtensilsCrossed,
  Sparkles,
  Laptop,
  Phone,
  Menu,
  CheckCircle2,
  Loader2
} from 'lucide-vue-next'
import apiClient from '../api/client'
import { getSystemSlug, slugify, highlightKeyword } from '../utils/slug'
import { detectSystemIcon } from '../utils/iconResolver'

const router = useRouter()
const searchQuery = ref('')
const isLoading = ref(false)
const isSearching = ref(false)
const liveSearchResults = ref([])
const mobileMenuOpen = ref(false)
let searchDebounceTimer = null

// Metadata mapping for the 9 database categories in dbprobuscs
const categoryMetaMap = {
  'PHIS': {
    badge: 'Hotel & Resort',
    desc: 'Property Management System komprehensif untuk front office, reservasi kamar, dan housekeeping.',
    icon: Building2,
    displayTitle: 'PHIS - Probus Hotel PMS',
  },
  'PRESTO': {
    badge: 'Restoran & POS',
    desc: 'Sistem point-of-sale restoran terpadu mulai dari order taker, kitchen display, hingga kasir & split bill.',
    icon: MonitorCheck,
    displayTitle: 'PRESTO Resto POS',
  },
  'PPOS': {
    badge: 'Retail & Kasir',
    desc: 'Sistem kasir cepat dan terintegrasi untuk outlet FnB, stand makanan, dan penjualan langsung.',
    icon: Printer,
    displayTitle: 'PPOS - Probus Point of Sales',
  },
  'SELF ORDER': {
    badge: 'Self Service',
    desc: 'Sistem pemesanan mandiri oleh tamu melalui scan QR code meja tanpa antre di kasir.',
    icon: Smartphone,
    displayTitle: 'Probus Self Order',
  },
  'PR-PO ONLINE': {
    badge: 'Logistik & Gudang',
    desc: 'Manajemen persediaan barang gudang, kartu stok, purchase request, dan penerimaan barang.',
    icon: Package,
    displayTitle: 'PR-PO Online & Purchasing',
  },
  'ORDER TAKER': {
    badge: 'Waiter App',
    desc: 'Aplikasi mobile waiter untuk input pesanan langsung dari meja tamu terkirim instan ke dapur.',
    icon: UtensilsCrossed,
    displayTitle: 'Order Taker Mobile',
  },
  'VRESTROPICAL': {
    badge: 'Club & Lounge',
    desc: 'Solusi kasir dan manajemen resto khusus konsep tropical bar, beach club, dan open-air dining.',
    icon: Sparkles,
    displayTitle: 'Vrestropical POS System',
  },
  'PPE': {
    badge: 'Facility & Asset',
    desc: 'Platform manajemen pemeliharaan fasilitas (maintenance), penugasan staf, dan kontrol aset.',
    icon: Zap,
    displayTitle: 'PPE 2.0 Management',
  },
  'BARONG POS': {
    badge: 'Cloud POS',
    desc: 'Sistem point of sales generasi baru dengan arsitektur modular dan sinkronisasi multi-outlet.',
    icon: Laptop,
    displayTitle: 'Barong POS Cloud',
  },
}

// Fallback initial list containing all 9 categories from dbprobuscs
const systems = ref([
  { id: 2, name: 'PRESTO', title: 'PRESTO Resto POS', badge: 'Restoran & POS', desc: categoryMetaMap['PRESTO'].desc, modulCount: 10, icon: MonitorCheck },
  { id: 1, name: 'PHIS', title: 'PHIS - Probus Hotel PMS', badge: 'Hotel & Resort', desc: categoryMetaMap['PHIS'].desc, modulCount: 10, icon: Building2 },
  { id: 8, name: 'PPE', title: 'PPE 2.0 Management', badge: 'Facility & Asset', desc: categoryMetaMap['PPE'].desc, modulCount: 10, icon: Zap },
  { id: 3, name: 'PPOS', title: 'PPOS - Probus Point of Sales', badge: 'Retail & Kasir', desc: categoryMetaMap['PPOS'].desc, modulCount: 10, icon: Printer },
  { id: 5, name: 'PR-PO ONLINE', title: 'PR-PO Online & Purchasing', badge: 'Logistik & Gudang', desc: categoryMetaMap['PR-PO ONLINE'].desc, modulCount: 10, icon: Package },
  { id: 4, name: 'SELF ORDER', title: 'Probus Self Order', badge: 'Self Service', desc: categoryMetaMap['SELF ORDER'].desc, modulCount: 10, icon: Smartphone },
  { id: 6, name: 'ORDER TAKER', title: 'Order Taker Mobile', badge: 'Waiter App', desc: categoryMetaMap['ORDER TAKER'].desc, modulCount: 10, icon: UtensilsCrossed },
  { id: 7, name: 'VRESTROPICAL', title: 'Vrestropical POS System', badge: 'Club & Lounge', desc: categoryMetaMap['VRESTROPICAL'].desc, modulCount: 10, icon: Sparkles },
  { id: 9, name: 'BARONG POS', title: 'Barong POS Cloud', badge: 'Cloud POS', desc: categoryMetaMap['BARONG POS'].desc, modulCount: 8, icon: Laptop },
])

onMounted(async () => {
  try {
    isLoading.value = true
    const res = await apiClient.get('/categories')
    if (res.data && res.data.data && Array.isArray(res.data.data) && res.data.data.length > 0) {
      const dbCategories = res.data.data
      systems.value = dbCategories.map(cat => {
        const meta = categoryMetaMap[cat.name] || {
          badge: 'Sistem Probus',
          desc: cat.description || `Panduan operasional dan solusi kendala sistem ${cat.name}.`,
          displayTitle: `${cat.name} System`,
        }
        return {
          id: cat.id,
          name: cat.name,
          title: meta.displayTitle || cat.name,
          badge: meta.badge,
          desc: cat.description || meta.desc,
          modulCount: cat.total_menus || cat.modul_count || (cat.name === 'BARONG POS' ? 8 : 10),
          icon: detectSystemIcon(cat.name),
        }
      })
    }
  } catch (err) {
    console.warn('API /categories unavailable, falling back to local dataset:', err)
  } finally {
    isLoading.value = false
  }
})

// Quick search index
const quickSearchDatabase = [
  { id: 1, type: 'article', title: 'Solusi Printer Kasir Thermal Tidak Mencetak / Error Paper Out', category: 'PRESTO', systemId: 2, snippet: 'Pengecekan kabel printer, paper jam, setting baudrate, dan test print di kasir.' },
  { id: 2, type: 'article', title: 'Cara Mengatasi Status Settlement EDC Pending / Gagal Kirim', category: 'PRESTO', systemId: 2, snippet: 'Prosedur audit batch settlement mesin EDC merchant bank dan verifikasi koneksi LAN.' },
  { id: 3, type: 'article', title: 'Mengatasi Selisih Kasir saat Tutup Shift (Z-Report)', category: 'PRESTO', systemId: 2, snippet: 'Penyebab selisih kas fisik dan langkah rekonsiliasi pembayaran tunai & non-tunai.' },
  { id: 4, type: 'article', title: 'Panduan Split Bill & Merge Table Tamu Rombongan', category: 'PRESTO', systemId: 2, snippet: 'Tata cara pisah tagihan per tamu atau gabung beberapa meja menjadi satu bill transaksi.' },
  { id: 5, type: 'article', title: 'Koneksi Kitchen Display (KDS) Terputus dari Server Kasir', category: 'PRESTO', systemId: 2, snippet: 'Troubleshooting IP address statis, port WebSocket, dan restart service KDS dapur.' },
  { id: 6, type: 'manual', moduleId: 1, title: 'Tata Cara Setup Master Data Meja, Menu & Kategori PRESTO', category: 'PRESTO', systemId: 2, snippet: 'Konfigurasi layout denah meja restoran, kategori makanan minuman, dan harga jual.' },
  { id: 7, type: 'manual', moduleId: 2, title: 'Alur Penerimaan Pesanan Meja & Pengiriman ke Dapur (Order Taker)', category: 'PRESTO', systemId: 2, snippet: 'Langkah waiter input order makanan via tablet atau POS kasir langsung ke printer dapur.' },
  { id: 8, type: 'manual', moduleId: 10, title: 'Panduan Uji Coba Modul Hardware & Routing Printer', category: 'PRESTO', systemId: 2, snippet: 'Testing printer checker, printer kasir, cash drawer kick RJ-11, dan network printer.' },
]

// Debounced live backend search
watch(searchQuery, (newVal) => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  const q = newVal ? newVal.trim() : ''
  if (!q) {
    liveSearchResults.value = []
    isSearching.value = false
    return
  }

  isSearching.value = true
  searchDebounceTimer = setTimeout(async () => {
    try {
      const res = await apiClient.get('/search', { params: { q, limit: 12 } })
      if (res.data && res.data.data && Array.isArray(res.data.data)) {
        liveSearchResults.value = res.data.data
      } else {
        liveSearchResults.value = []
      }
    } catch (err) {
      console.warn('Live search fallback to local index:', err)
      liveSearchResults.value = []
    } finally {
      isSearching.value = false
    }
  }, 200)
})

// Unified matching results across systems, backend database, and quick index
const searchResults = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return []
  const queryWords = q.split(/\s+/).filter(Boolean)

  const results = []
  const seenKeys = new Set()

  // 1. Check matching Systems from `systems.value`
  systems.value.forEach(sys => {
    const match = queryWords.every(w =>
      sys.name.toLowerCase().includes(w) ||
      (sys.title && sys.title.toLowerCase().includes(w)) ||
      (sys.badge && sys.badge.toLowerCase().includes(w)) ||
      (sys.desc && sys.desc.toLowerCase().includes(w))
    )
    if (match) {
      const key = 'system-' + sys.id
      if (!seenKeys.has(key)) {
        seenKeys.add(key)
        results.push({
          id: sys.id,
          type: 'system',
          title: `${sys.name} - ${sys.title}`,
          category: sys.badge || 'Sistem Probus',
          systemId: sys.id,
          snippet: sys.desc || `Panduan modul dan sistem ${sys.name}`,
        })
      }
    }
  })

  // 2. Add backend live search results
  if (liveSearchResults.value && liveSearchResults.value.length > 0) {
    liveSearchResults.value.forEach(item => {
      const itemType = item.type === 'manual_book' ? 'manual' : item.type
      const key = `${itemType}-${item.id}`
      if (!seenKeys.has(key)) {
        seenKeys.add(key)
        results.push({
          id: item.id,
          type: itemType,
          title: item.title,
          category: item.category || 'PRESTO',
          systemId: 2,
          moduleId: itemType === 'manual' ? item.id : null,
          snippet: item.snippet || '',
        })
      }
    })
  }

  // 3. Match supplementary quickSearchDatabase
  quickSearchDatabase.forEach(item => {
    const match = queryWords.every(w =>
      item.title.toLowerCase().includes(w) ||
      item.category.toLowerCase().includes(w) ||
      (item.snippet && item.snippet.toLowerCase().includes(w))
    )
    if (match) {
      const key = `${item.type}-${item.id}`
      if (!seenKeys.has(key)) {
        seenKeys.add(key)
        results.push(item)
      }
    }
  })

  return results
})

function selectCard(id) {
  router.push('/system/' + getSystemSlug(id))
}

function handleSearchItemClick(item) {
  if (item.type === 'system') {
    router.push('/system/' + getSystemSlug(item.systemId || item.id))
  } else if (item.type === 'manual' || item.type === 'manual_book') {
    router.push(`/system/${getSystemSlug(item.systemId || 2)}/guide/${item.moduleId || item.id || 2}`)
  } else if (item.type === 'article') {
    router.push(`/solutions/${slugify(item.title)}`)
  } else {
    router.push('/system/' + getSystemSlug(item.systemId || 2))
  }
}

function clearSearch() {
  searchQuery.value = ''
  liveSearchResults.value = []
}
</script>

<template>
  <div class="min-h-screen bg-white text-slate-800 font-sans flex flex-col antialiased">
    
    <!-- 1. Utility Top Bar -->
    <header class="border-b border-slate-100 bg-white text-xs text-slate-500 py-2.5 px-4 sm:px-6 lg:px-12 xl:px-16">
      <div class="w-full flex items-center justify-between sm:justify-end">
        <span class="text-slate-400 text-[11px] sm:hidden">Pusat Bantuan Resmi</span>
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
          <span class="hidden sm:inline text-slate-300 font-light">|</span>
          <span class="hidden sm:inline text-slate-500 text-xs sm:text-sm font-medium">Pusat Bantuan Publik</span>
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
        <router-link
          to="/contact"
          @click="mobileMenuOpen = false"
          class="block px-3 py-2 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-slate-900"
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

    <!-- 3. Clean Hero & Search Omnibox Section -->
    <section class="pt-10 pb-8 sm:pt-16 sm:pb-12 lg:pt-20 lg:pb-16 px-4 sm:px-6 lg:px-12 xl:px-16 text-center bg-gradient-to-b from-white via-sky-50/20 to-transparent">
      <div class="w-full max-w-4xl lg:max-w-5xl xl:max-w-6xl mx-auto">
        <h1 class="text-2xl sm:text-3xl lg:text-4xl xl:text-5xl font-extrabold text-slate-900 tracking-tight leading-tight">
          Pusat Bantuan & Panduan Sistem Probus
        </h1>
        <p class="mt-3 sm:mt-4 text-slate-500 text-sm sm:text-base lg:text-lg leading-relaxed max-w-3xl mx-auto">
          Pilih sistem atau produk yang Anda gunakan untuk melihat manual book dan solusi kendala teknis.
        </p>

        <!-- Search Omnibox -->
        <div class="mt-7 sm:mt-8 w-full max-w-2xl sm:max-w-3xl lg:max-w-4xl xl:max-w-5xl mx-auto relative group">
          <!-- Search Icon with Centering & Focus Effect -->
          <div class="absolute left-4 sm:left-5 top-1/2 -translate-y-1/2 pointer-events-none flex items-center justify-center transition-colors duration-200 text-slate-400 group-focus-within:text-sky-600">
            <Search class="w-5 h-5" />
          </div>

          <!-- Search Input Field -->
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Ketik kendala, fitur, atau nama modul yang Anda cari..."
            class="w-full pl-12 sm:pl-14 pr-32 sm:pr-40 py-3.5 sm:py-4 rounded-2xl text-sm sm:text-base bg-white border border-slate-200/90 text-slate-800 placeholder:text-slate-400 font-normal shadow-xs hover:border-slate-300 hover:shadow-sm focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/15 focus:shadow-md transition-all duration-200"
          />

          <!-- Right Side Indicators (Tanda Cocok & Clear Button) -->
          <div class="absolute right-3.5 sm:right-4 top-1/2 -translate-y-1/2 flex items-center gap-1.5">
            <!-- Match indicator badge -->
            <span
              v-if="searchQuery && isSearching"
              class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-sky-50 text-sky-700 border border-sky-200 animate-pulse pointer-events-none"
            >
              <Loader2 class="w-3 h-3 animate-spin" />
              <span class="hidden sm:inline">Mencari...</span>
            </span>
            <span
              v-else-if="searchQuery && searchResults.length > 0"
              class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-emerald-50 text-emerald-700 border border-emerald-300 shadow-2xs pointer-events-none animate-in fade-in"
            >
              <CheckCircle2 class="w-3 h-3 text-emerald-600" />
              <span>{{ searchResults.length }} cocok</span>
            </span>
            <span
              v-else-if="searchQuery && !isSearching && searchResults.length === 0"
              class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] sm:text-[11px] font-semibold bg-rose-50 text-rose-600 border border-rose-200 pointer-events-none"
            >
              <span class="hidden sm:inline">Tidak ada</span> cocok
            </span>

            <!-- Clear Search Button -->
            <button
              v-if="searchQuery"
              @click="clearSearch"
              type="button"
              class="p-1.5 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors focus:outline-none cursor-pointer"
              title="Hapus pencarian"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <!-- Live Search Results Dropdown -->
          <div
            v-if="searchQuery && searchResults.length > 0"
            class="absolute left-0 right-0 top-full mt-2.5 bg-white rounded-2xl border border-slate-200/90 shadow-xl overflow-hidden z-50 text-left animate-in fade-in slide-in-from-top-2 duration-150"
          >
            <div class="px-4 py-3 bg-slate-50/80 border-b border-slate-100 flex items-center justify-between text-xs font-semibold text-slate-600">
              <span class="flex items-center gap-1.5">
                <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600 shrink-0" />
                <span>Ditemukan <strong>{{ searchResults.length }}</strong> hasil cocok untuk "<strong class="text-sky-700">{{ searchQuery }}</strong>"</span>
              </span>
              <span class="bg-emerald-100 text-emerald-800 px-2 py-0.5 rounded-full font-bold text-[11px] shrink-0">{{ searchResults.length }} Cocok</span>
            </div>
            <div class="max-h-80 overflow-y-auto divide-y divide-slate-100 custom-sky-scrollbar">
              <div
                v-for="item in searchResults"
                :key="item.type + '-' + item.id"
                @click="handleSearchItemClick(item)"
                class="p-3.5 sm:p-4 hover:bg-sky-50/60 cursor-pointer transition-colors flex items-start gap-3.5 group"
              >
                <div
                  :class="[
                    'w-9 h-9 rounded-xl flex items-center justify-center shrink-0 mt-0.5',
                    item.type === 'system' ? 'bg-sky-100 text-sky-700 border border-sky-300' :
                    item.type === 'article' ? 'bg-amber-50 text-amber-600 border border-amber-200/60' : 'bg-emerald-50 text-emerald-600 border border-emerald-200/60'
                  ]"
                >
                  <Building2 v-if="item.type === 'system'" class="w-4 h-4" />
                  <AlertCircle v-else-if="item.type === 'article'" class="w-4 h-4" />
                  <BookOpen v-else class="w-4 h-4" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <span
                      :class="[
                        'text-[10px] font-bold px-2 py-0.5 rounded-md uppercase tracking-wider',
                        item.type === 'system' ? 'bg-sky-50 text-sky-700 border border-sky-200' :
                        item.type === 'article' ? 'bg-amber-50 text-amber-700 border border-amber-200' : 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                      ]"
                    >
                      {{ item.type === 'system' ? 'Sistem / Produk' : item.type === 'article' ? 'Artikel Solusi' : 'Buku Panduan' }}
                    </span>
                    <span class="text-xs text-slate-400">• {{ item.category }}</span>
                  </div>
                  <h4
                    class="text-sm font-semibold text-slate-800 mt-1 truncate group-hover:text-sky-600 transition-colors"
                    v-html="highlightKeyword(item.title, searchQuery)"
                  ></h4>
                  <p
                    v-if="item.snippet"
                    class="text-xs text-slate-500 mt-0.5 line-clamp-1"
                    v-html="highlightKeyword(item.snippet, searchQuery)"
                  ></p>
                </div>
                <ArrowRight class="w-4 h-4 text-slate-400 shrink-0 self-center group-hover:text-sky-600 group-hover:translate-x-0.5 transition-all" />
              </div>
            </div>
          </div>

          <!-- Empty State Dropdown -->
          <div
            v-else-if="searchQuery && searchResults.length === 0"
            class="absolute left-0 right-0 top-full mt-2.5 bg-white rounded-2xl border border-slate-200/90 shadow-xl p-6 z-50 text-center animate-in fade-in duration-150"
          >
            <p class="text-sm text-slate-600 font-medium">Tidak ada hasil untuk "<span class="text-slate-800 font-semibold">{{ searchQuery }}</span>"</p>
            <p class="text-xs text-slate-400 mt-1">Coba gunakan kata kunci lain seperti nama fitur, nomor modul, atau topik kendala.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 4. Section: Pilih Produk / Sistem Probus -->
    <main class="flex-1 w-full px-4 sm:px-6 lg:px-12 xl:px-16 pb-16">
      
      <!-- Section Header -->
      <div class="mb-6 sm:mb-7 flex flex-col sm:flex-row sm:items-end justify-between gap-2 sm:gap-3">
        <div>
          <h2 class="text-xl sm:text-2xl font-bold text-slate-900 tracking-tight">
            Pilih Produk / Sistem Probus
          </h2>
          <p class="text-xs sm:text-sm text-slate-500 mt-1">
            Pilih sistem yang digunakan oleh properti Anda untuk melihat modul penggunaan dan solusi kendala:
          </p>
        </div>
        <div class="text-xs font-semibold px-3 py-1.5 rounded-full bg-slate-50 border border-slate-200 text-slate-600 shadow-2xs w-fit">
          {{ systems.length }} Kategori Sistem Tersedia
        </div>
      </div>

      <!-- 9 Cards Grid (3 Columns on Large Screens, 2 on Tablet, 1 on Mobile) -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5 lg:gap-6 xl:gap-8">
        <div
          v-for="sys in systems"
          :key="sys.id"
          @click="selectCard(sys.id)"
          class="bg-white rounded-2xl p-5 sm:p-6 lg:p-7 xl:p-8 border border-slate-200/90 hover:border-sky-500 hover:ring-4 hover:ring-sky-500/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200 cursor-pointer flex flex-col justify-between group shadow-xs"
        >
          <div>
            <!-- Card Top: Icon & Badge -->
            <div class="flex items-center justify-between">
              <div
                class="w-11 sm:w-12 h-11 sm:h-12 rounded-xl flex items-center justify-center bg-slate-100 text-slate-700 group-hover:bg-sky-50 group-hover:text-sky-600 transition-colors shadow-2xs"
              >
                <component :is="sys.icon" class="w-5 sm:w-6 h-5 sm:h-6 stroke-[1.8]" />
              </div>
              <span
                class="text-xs font-semibold px-3 py-1.5 rounded-lg tracking-wide bg-slate-100 text-slate-600 group-hover:bg-sky-100/70 group-hover:text-sky-700 transition-colors"
              >
                {{ sys.badge }}
              </span>
            </div>

            <!-- Card Middle: Title & Description -->
            <h3 class="font-bold text-slate-900 text-base sm:text-lg lg:text-xl mt-4 sm:mt-5 group-hover:text-sky-600 transition-colors">
              {{ sys.title }}
            </h3>
            <p class="text-xs sm:text-sm text-slate-500 mt-2 sm:mt-2.5 leading-relaxed line-clamp-2">
              {{ sys.desc }}
            </p>
          </div>

          <!-- Card Bottom: Modul Count & Link -->
          <div class="pt-5 sm:pt-6 mt-5 sm:mt-6 border-t border-slate-100 flex items-center justify-between text-xs sm:text-sm">
            <span class="font-medium text-slate-500">
              {{ sys.modulCount }} Modul Penggunaan
            </span>
            <span class="text-sky-600 font-semibold flex items-center gap-1.5 group-hover:text-sky-700">
              <span>Pilih Sistem</span>
              <ArrowRight class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
            </span>
          </div>
        </div>
      </div>

    </main>

    <!-- 5. Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 px-4 sm:px-6 lg:px-12 xl:px-16 text-xs sm:text-sm text-slate-500">
      <div class="w-full flex flex-col sm:flex-row items-center justify-between gap-3 text-center sm:text-left">
        <p class="font-normal">© 2026 CV Mitra Solusi Mandiri</p>
      </div>
    </footer>
  </div>
</template>
