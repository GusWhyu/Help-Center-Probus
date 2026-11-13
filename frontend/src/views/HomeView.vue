<script setup>
import { ref, computed, onMounted } from 'vue'
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
  Phone
} from 'lucide-vue-next'
import apiClient from '../api/client'

const router = useRouter()
const searchQuery = ref('')
const selectedSystem = ref(2) // default active: PRESTO (ID: 2 in DB)
const isLoading = ref(false)

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
          icon: MonitorCheck,
          displayTitle: `${cat.name} System`,
        }
        return {
          id: cat.id,
          name: cat.name,
          title: meta.displayTitle || cat.name,
          badge: meta.badge,
          desc: cat.description || meta.desc,
          modulCount: cat.modul_count || (cat.name === 'BARONG POS' ? 8 : 10),
          icon: meta.icon,
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
  { id: 1, type: 'article', title: 'Solusi Printer Kasir Thermal Tidak Mencetak / Error Paper Out', category: 'PRESTO', systemId: 2 },
  { id: 2, type: 'article', title: 'Cara Mengatasi Status Settlement EDC Pending / Gagal Kirim', category: 'PRESTO', systemId: 2 },
  { id: 3, type: 'article', title: 'Mengatasi Selisih Kasir saat Tutup Shift (Z-Report)', category: 'PRESTO', systemId: 2 },
  { id: 4, type: 'article', title: 'Panduan Split Bill & Merge Table Tamu Rombongan', category: 'PRESTO', systemId: 2 },
  { id: 5, type: 'article', title: 'Koneksi Kitchen Display (KDS) Terputus dari Server Kasir', category: 'PRESTO', systemId: 2 },
  { id: 6, type: 'manual', moduleId: 1, title: 'Tata Cara Setup Master Data Meja, Menu & Kategori PRESTO', category: 'PRESTO', systemId: 2 },
  { id: 7, type: 'manual', moduleId: 2, title: 'Alur Penerimaan Pesanan Meja & Pengiriman ke Dapur (Order Taker)', category: 'PRESTO', systemId: 2 },
  { id: 8, type: 'manual', moduleId: 10, title: 'Panduan Uji Coba Modul Hardware & Routing Printer', category: 'PRESTO', systemId: 2 },
]

const searchResults = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return []
  return quickSearchDatabase.filter(item =>
    item.title.toLowerCase().includes(q) || item.category.toLowerCase().includes(q)
  )
})

function selectCard(id) {
  selectedSystem.value = id
  router.push('/system/' + id)
}

function handleSearchItemClick(item) {
  if (item.type === 'manual') {
    router.push(`/system/${item.systemId || 2}/guide/${item.moduleId || 2}`)
  } else if (item.type === 'article') {
    router.push(`/solutions/${item.id}`)
  } else {
    router.push('/system/' + (item.systemId || 2))
  }
}

function clearSearch() {
  searchQuery.value = ''
}
</script>

<template>
  <div class="min-h-screen bg-white text-slate-800 font-sans flex flex-col antialiased">
    
    <!-- 1. Utility Top Bar -->
    <header class="border-b border-slate-100 bg-white text-xs text-slate-500 py-2.5 px-6 sm:px-12 lg:px-20 xl:px-28">
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
    <nav class="border-b border-slate-200/80 bg-white/95 backdrop-blur-sm py-4 px-6 sm:px-12 lg:px-20 xl:px-28 sticky top-0 z-40 shadow-xs">
      <div class="max-w-[1600px] mx-auto flex items-center justify-between">
        <!-- Logo & Title -->
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2.5 cursor-pointer" @click="router.push('/')">
            <img
              src="/probus-logo.png"
              alt="Probus System"
              class="h-9 sm:h-10 w-auto object-contain hover:opacity-90 transition-opacity"
            />
          </div>
          <span class="text-slate-300 font-light">|</span>
          <span class="text-slate-500 text-sm sm:text-base font-medium">Pusat Bantuan Publik</span>
        </div>

        <!-- Navigation Links -->
        <div class="hidden lg:flex items-center gap-9 text-sm font-medium">
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

    <!-- 3. Clean Hero & Search Omnibox Section -->
    <section class="pt-16 pb-12 sm:pt-20 sm:pb-16 px-6 sm:px-12 lg:px-20 xl:px-28 text-center bg-gradient-to-b from-white via-sky-50/20 to-transparent">
      <div class="max-w-4xl mx-auto">
        <h1 class="text-3xl sm:text-4xl lg:text-[2.75rem] font-extrabold text-slate-900 tracking-tight leading-tight">
          Pusat Bantuan & Panduan Sistem Probus
        </h1>
        <p class="mt-4 text-slate-500 text-base sm:text-lg leading-relaxed max-w-2xl mx-auto">
          Pilih sistem atau produk yang Anda gunakan untuk melihat manual book dan solusi kendala teknis.
        </p>

        <!-- Search Omnibox -->
        <div class="mt-8 max-w-2xl mx-auto relative shadow-xs">
          <Search class="absolute left-4.5 w-5 h-5 text-slate-400 pointer-events-none top-4" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Ketik kendala, fitur, atau nama modul yang Anda cari..."
            class="w-full pl-12 pr-12 py-3.5 sm:py-4 rounded-xl text-base bg-white border border-slate-300 text-slate-900 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-100 transition-all shadow-xs"
          />
          <button
            v-if="searchQuery"
            @click="clearSearch"
            class="absolute right-4 top-4 text-slate-400 hover:text-slate-600"
          >
            <X class="w-4 h-4" />
          </button>

          <!-- Live Search Results Dropdown -->
          <div
            v-if="searchQuery && searchResults.length > 0"
            class="absolute left-0 right-0 top-full mt-2 bg-white rounded-xl border border-slate-200 shadow-xl overflow-hidden z-50 text-left animate-in fade-in slide-in-from-top-2 duration-150"
          >
            <div class="p-3 bg-slate-50 border-b border-slate-100 text-xs font-bold text-slate-500 uppercase tracking-wider">
              Hasil Pencarian ({{ searchResults.length }})
            </div>
            <div class="max-h-80 overflow-y-auto divide-y divide-slate-100">
              <div
                v-for="item in searchResults"
                :key="item.id"
                @click="handleSearchItemClick(item)"
                class="p-4 hover:bg-sky-50/60 cursor-pointer transition-colors flex items-start gap-3.5"
              >
                <div
                  :class="[
                    'w-8 h-8 rounded-lg flex items-center justify-center shrink-0 mt-0.5',
                    item.type === 'article' ? 'bg-amber-100 text-amber-700' : 'bg-sky-100 text-sky-700'
                  ]"
                >
                  <AlertCircle v-if="item.type === 'article'" class="w-4 h-4" />
                  <BookOpen v-else class="w-4 h-4" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <span
                      :class="[
                        'text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider',
                        item.type === 'article' ? 'bg-amber-50 text-amber-700 border border-amber-200' : 'bg-sky-50 text-sky-700 border border-sky-200'
                      ]"
                    >
                      {{ item.type === 'article' ? 'Artikel Solusi' : 'Buku Panduan' }}
                    </span>
                    <span class="text-xs text-slate-400">• {{ item.category }}</span>
                  </div>
                  <h4 class="text-sm font-semibold text-slate-800 mt-1 truncate hover:text-sky-600">
                    {{ item.title }}
                  </h4>
                </div>
                <ArrowRight class="w-4 h-4 text-slate-400 shrink-0 self-center" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 4. Section: Pilih Produk / Sistem Probus -->
    <main class="flex-1 max-w-[1600px] w-full mx-auto px-6 sm:px-12 lg:px-20 xl:px-28 pb-16">
      
      <!-- Section Header -->
      <div class="mb-7 flex flex-col sm:flex-row sm:items-end justify-between gap-3">
        <div>
          <h2 class="text-2xl sm:text-[1.65rem] font-bold text-slate-900 tracking-tight">
            Pilih Produk / Sistem Probus
          </h2>
          <p class="text-sm sm:text-base text-slate-500 mt-1.5">
            Pilih sistem yang digunakan oleh properti Anda untuk melihat modul penggunaan dan solusi kendala:
          </p>
        </div>
        <div class="text-xs font-semibold px-3 py-1.5 rounded-full bg-slate-50 border border-slate-200 text-slate-600 shadow-2xs w-fit">
          {{ systems.length }} Kategori Sistem Tersedia
        </div>
      </div>

      <!-- 9 Cards Grid (3 Columns on Large Screens, 2 on Tablet, 1 on Mobile) -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-7">
        <div
          v-for="sys in systems"
          :key="sys.id"
          @click="selectCard(sys.id)"
          :class="[
            'bg-white rounded-2xl p-6 sm:p-7 border transition-all duration-200 cursor-pointer flex flex-col justify-between group shadow-xs',
            selectedSystem === sys.id
              ? 'border-2 border-sky-500 shadow-md ring-4 ring-sky-500/10'
              : 'border-slate-200/90 hover:border-slate-300 hover:shadow-lg hover:-translate-y-1'
          ]"
        >
          <div>
            <!-- Card Top: Icon & Badge -->
            <div class="flex items-center justify-between">
              <div
                :class="[
                  'w-12 h-12 rounded-xl flex items-center justify-center transition-colors shadow-2xs',
                  selectedSystem === sys.id ? 'bg-sky-50 text-sky-600' : 'bg-slate-100 text-slate-700 group-hover:bg-sky-50 group-hover:text-sky-600'
                ]"
              >
                <component :is="sys.icon" class="w-6 h-6 stroke-[1.8]" />
              </div>
              <span
                :class="[
                  'text-xs font-semibold px-3 py-1.5 rounded-lg tracking-wide transition-colors',
                  selectedSystem === sys.id
                    ? 'bg-sky-600 text-white shadow-2xs'
                    : 'bg-slate-100 text-slate-600 group-hover:bg-slate-200'
                ]"
              >
                {{ sys.badge }}
              </span>
            </div>

            <!-- Card Middle: Title & Description -->
            <h3 class="font-bold text-slate-900 text-lg sm:text-xl mt-5 group-hover:text-sky-600 transition-colors">
              {{ sys.title }}
            </h3>
            <p class="text-sm text-slate-500 mt-2.5 leading-relaxed line-clamp-2">
              {{ sys.desc }}
            </p>
          </div>

          <!-- Card Bottom: Modul Count & Link -->
          <div class="pt-6 mt-6 border-t border-slate-100 flex items-center justify-between text-sm">
            <span
              :class="[
                'font-medium text-xs sm:text-sm',
                selectedSystem === sys.id ? 'text-sky-600 font-bold' : 'text-slate-500'
              ]"
            >
              {{ sys.modulCount }} Modul Penggunaan
            </span>
            <span class="text-sky-600 font-semibold flex items-center gap-1.5 hover:text-sky-700 text-xs sm:text-sm">
              <span>Pilih Sistem</span>
              <ArrowRight class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
            </span>
          </div>
        </div>
      </div>

      <!-- 5. Bottom Help / Notification Banner -->
      <div class="mt-12 bg-sky-50/70 border border-sky-200/90 rounded-2xl p-6 sm:p-7 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-5 shadow-xs">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-sky-100 text-sky-600 flex items-center justify-center shrink-0 shadow-2xs border border-sky-200/70">
            <Lightbulb class="w-6 h-6 stroke-[2]" />
          </div>
          <div>
            <h4 class="text-base sm:text-lg text-slate-900 font-bold">
              Tidak menemukan sistem yang Anda cari?
            </h4>
            <p class="text-xs sm:text-sm text-slate-600 mt-0.5">
              Konsultasikan langsung kendala Anda dengan tim support teknis Probus untuk panduan sistem lainnya.
            </p>
          </div>
        </div>

        <a
          href="https://wa.me/62811389211"
          target="_blank"
          class="shrink-0 bg-sky-600 hover:bg-sky-700 text-white font-semibold text-xs sm:text-sm px-5 py-3 rounded-xl flex items-center gap-2.5 shadow-xs transition-all hover:shadow-md"
        >
          <MessageSquare class="w-4 h-4 fill-white/20" />
          <span>Hubungi WhatsApp Support</span>
        </a>
      </div>
    </main>

    <!-- 6. Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 px-6 sm:px-12 lg:px-20 xl:px-28 text-xs sm:text-sm text-slate-500">
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
