<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  ArrowRight,
  Search,
  BookOpen,
  Lightbulb,
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
  Printer
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

// Get system ID from route parameter, default to 2 (PRESTO)
const systemId = computed(() => Number(route.params.id) || 2)

// Active Tab: 'manual' or 'solusi'
const activeTab = ref('manual')
const searchQuery = ref('')

// System metadata mapping
const systemsInfo = {
  1: { id: 1, name: 'PHIS', title: 'Panduan Pengguna Probus Hotel PMS', badge: 'Hotel & Resort' },
  2: { id: 2, name: 'PRESTO', title: 'Panduan Pengguna PRESTO Resto POS', badge: 'Restoran & POS' },
  3: { id: 3, name: 'PPOS', title: 'Panduan Pengguna PPOS Point of Sales', badge: 'Retail & Kasir' },
  4: { id: 4, name: 'SELF ORDER', title: 'Panduan Pengguna Probus Self Order', badge: 'Self Service' },
  5: { id: 5, name: 'PR-PO ONLINE', title: 'Panduan Pengguna PR-PO Online Logistik', badge: 'Logistik & Gudang' },
  6: { id: 6, name: 'ORDER TAKER', title: 'Panduan Pengguna Order Taker Mobile', badge: 'Waiter App' },
  7: { id: 7, name: 'VRESTROPICAL', title: 'Panduan Pengguna Vrestropical POS', badge: 'Club & Lounge' },
  8: { id: 8, name: 'PPE', title: 'Panduan Pengguna PPE 2.0 Management', badge: 'Facility & Asset' },
  9: { id: 9, name: 'BARONG POS', title: 'Panduan Pengguna Barong POS Cloud', badge: 'Cloud POS' },
}

const currentSystem = computed(() => systemsInfo[systemId.value] || systemsInfo[2])

// 10 Real Modules of PRESTO matching dbprobuscs `menu_systems` (Order 1 to 10)
const presto10Modules = [
  {
    id: 1,
    dbMenuId: 34,
    order: 1,
    title: '01. Setup Master',
    countText: '1 Panduan',
    icon: Settings,
    desc: 'Modul setup master data, parameter sistem, konfigurasi outlet resto, denah meja, dan pengaturan menu makanan.',
    estTime: 'Est. 4 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Pengenalan Modul Setup Master</h3>
        <p class="text-slate-600">Modul Setup Master adalah pusat konfigurasi utama sistem PRESTO Resto POS untuk mendefinisikan seluruh parameter operasional sebelum transaksi kasir dijalankan.</p>
        
        <h3 class="text-lg font-bold text-slate-900 mt-6">2. Pengaturan yang Dikelola</h3>
        <ul class="list-disc pl-5 space-y-1.5 text-slate-600">
          <li><strong>Master Outlet & Pajak</strong>: Menentukan nama outlet, alamat, nomor NPWPD, dan persentase PB1 / Service Charge.</li>
          <li><strong>Master Kategori & Menu</strong>: Menginput daftar makanan, minuman, harga jual, barcode, dan pengelompokan departemen.</li>
          <li><strong>Master Denah Meja (Table Layout)</strong>: Mengatur denah meja visual per ruangan (Area Indoor AC, Outdoor, VIP Room, Bar).</li>
        </ul>
      </div>
    `,
  },
  {
    id: 2,
    dbMenuId: 35,
    order: 2,
    title: '02. Transaksi',
    countText: '3 Panduan',
    icon: UtensilsCrossed,
    desc: 'Modul transaksi penjualan operasional kasir POS, order taker meja, custom modifier pesanan, dan pembayaran split bill.',
    estTime: 'Est. 7 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Membuka Meja Baru (Open Table)</h3>
        <p>Pilih area denah restoran (misal: <em>Indoor AC</em>, <em>Outdoor Garden</em>, atau <em>Bar Area</em>). Klik pada nomor meja yang berstatus kosong (berwarna hijau). Masukkan jumlah tamu (Pax) lalu klik tombol <strong>Buka Meja</strong>.</p>
        
        <h3 class="text-lg font-bold text-slate-900 mt-6">2. Memilih & Menginput Menu Pesanan</h3>
        <p>Gunakan kategori menu di sebelah kiri (Makanan, Minuman, Dessert, atau Paket). Anda dapat mencari menu cepat dengan mengetik nama atau kode menu pada kolom pencarian.</p>
        <ul class="list-disc pl-5 space-y-1.5 text-slate-600">
          <li>Klik item menu untuk menambah jumlah porsi (+1).</li>
          <li>Klik ikon pensil untuk menambahkan <strong>Special Request / Modifier</strong> (contoh: <em>Less Sugar</em>, <em>Pedas Level 2</em>, <em>Tanpa Bawang</em>).</li>
        </ul>

        <h3 class="text-lg font-bold text-slate-900 mt-6">3. Mengirim Pesanan ke Dapur (Send Order)</h3>
        <p>Setelah pesanan selesai diinput, klik tombol <strong>Send Order / Simpan</strong>. Tiket pesanan otomatis tercetak di printer kitchen dan bar sesuai routing printer masing-masing departemen.</p>
      </div>
    `,
  },
  {
    id: 3,
    dbMenuId: 36,
    order: 3,
    title: '03. Outlet Activity',
    countText: '1 Panduan',
    icon: Store,
    desc: 'Modul aktivitas dan operasional outlet, pembagian shift kerja kasir, serah terima modal kas, dan pantauan kegiatan harian.',
    estTime: 'Est. 5 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Manajemen Aktivitas Harian Outlet</h3>
        <p class="text-slate-600">Modul ini digunakan oleh Supervisor Outlet untuk memantau aktivitas kasir yang sedang login, status meja yang sedang terisi, dan durasi rata-rata waktu makan tamu (table turnover rate).</p>
        
        <h3 class="text-lg font-bold text-slate-900 mt-6">2. Pembukaan & Penutupan Shift Kasir</h3>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li>Kasir melakukan <strong>Open Cash</strong> dengan menghitung uang kas awal di laci.</li>
          <li>Di akhir jam kerja, kasir melakukan <strong>Serah Terima Shift</strong> ke kasir shift berikutnya.</li>
          <li>Sistem memverifikasi selisih uang fisik dengan data transaksi sistem.</li>
        </ol>
      </div>
    `,
  },
  {
    id: 4,
    dbMenuId: 37,
    order: 4,
    title: '04. Cost Control',
    countText: '1 Panduan',
    icon: PieChart,
    desc: 'Modul pengendalian biaya, standar resep menu (recipe costing), food & beverage cost, dan laporan analisa selisih bahan.',
    estTime: 'Est. 6 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Fungsi Modul Cost Control</h3>
        <p class="text-slate-600">Memastikan biaya pokok penjualan (HPP / Cost of Goods Sold) F&B restoran selalu berada dalam target persentase yang ditetapkan manajemen.</p>

        <h3 class="text-lg font-bold text-slate-900 mt-6">2. Fitur Unggulan</h3>
        <ul class="list-disc pl-5 space-y-1.5 text-slate-600">
          <li><strong>Standar Resep Bahan Baku</strong>: Menghubungkan 1 porsi menu dengan gramasi bahan di gudang (otomatis memotong stok).</li>
          <li><strong>Analisa Variance (Selisih)</strong>: Menghitung selisih antara pemakaian teoritis sistem dengan stok opname fisik di dapur.</li>
        </ul>
      </div>
    `,
  },
  {
    id: 5,
    dbMenuId: 38,
    order: 5,
    title: '05. Purchasing AP',
    countText: '2 Panduan',
    icon: ShoppingBag,
    desc: 'Modul pengadaan barang, purchase request (PR), purchase order (PO), receiving gudang, dan hutang usaha (AP).',
    estTime: 'Est. 5 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Alur Pengadaan Barang Restoran</h3>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li><strong>Purchase Request (PR)</strong>: Kepala dapur/bar mengajukan permintaan belanja bahan baku yang menipis.</li>
          <li><strong>Purchase Order (PO)</strong>: Bagian purchasing menyetujui dan menerbitkan PO ke supplier rekanan.</li>
          <li><strong>Receiving</strong>: Bagian penerimaan memeriksa kecocokan fisik barang dan timbangan saat supplier tiba.</li>
          <li><strong>Accounts Payable (AP)</strong>: Tagihan supplier tercatat otomatis ke pos hutang usaha.</li>
        </ol>
      </div>
    `,
  },
  {
    id: 6,
    dbMenuId: 39,
    order: 6,
    title: '06. AR (Accounts Receivable)',
    countText: '1 Panduan',
    icon: Receipt,
    desc: 'Modul piutang usaha, penagihan invoice klien hotel/resto, city ledger voucher rekanan, dan rekonsiliasi pembayaran.',
    estTime: 'Est. 4 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Pengelolaan Piutang & City Ledger</h3>
        <p class="text-slate-600">Jika tamu resto melakukan pembayaran menggunakan sistem voucher perusahaan, kartu member kredit, atau tagihan kamar hotel (Room Charge):</p>
        <ul class="list-disc pl-5 space-y-1.5 text-slate-600">
          <li>Transaksi dialihkan ke akun City Ledger perusahaan terkait.</li>
          <li>Modul AR menghasilkan invoice resmi bertempo (due date) untuk ditagihkan ke pihak korporat/rekanan.</li>
        </ul>
      </div>
    `,
  },
  {
    id: 7,
    dbMenuId: 41,
    order: 7,
    title: '07. General Cashier',
    countText: '2 Panduan',
    icon: Wallet,
    desc: 'Modul kasir umum, pencatatan penerimaan uang kas harian, pengeluaran kas kecil (petty cash), dan setoran bank.',
    estTime: 'Est. 4 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Fungsi General Cashier</h3>
        <p class="text-slate-600">Merekonsiliasi seluruh uang tunai yang diserahkan oleh kasir-kasir operasional outlet di akhir hari, serta mengelola kas kecil untuk kebutuhan harian operasional (petty cash).</p>
        <h3 class="text-lg font-bold text-slate-900 mt-6">2. Prosedur Setoran Bank</h3>
        <p class="text-slate-600">Mencetak slip setoran fisik dan membukukan transfer setoran harian ke rekening giro operasional restoran.</p>
      </div>
    `,
  },
  {
    id: 8,
    dbMenuId: 42,
    order: 8,
    title: '08. Accounting System',
    countText: '2 Panduan',
    icon: BarChart3,
    desc: 'Modul sistem akuntansi, pembukuan jurnal otomatis dari transaksi POS, buku besar (GL), neraca, dan laporan keuangan.',
    estTime: 'Est. 6 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Integrasi Transaksi Kasir ke Jurnal Umum</h3>
        <p class="text-slate-600">Setiap penjualan makanan, minuman, diskon, dan pajak di kasir otomatis memicu jurnal debet/kredit ke bagan akun (Chart of Accounts):</p>
        <ul class="list-disc pl-5 space-y-1.5 text-slate-600">
          <li>Debet: Kas / Bank / Piutang EDC</li>
          <li>Kredit: Pendapatan Makanan & Minuman, Hutang Pajak Resto (PB1)</li>
        </ul>
      </div>
    `,
  },
  {
    id: 9,
    dbMenuId: 43,
    order: 9,
    title: '09. Report',
    countText: '3 Panduan',
    icon: TrendingUp,
    desc: 'Modul pusat laporan operasional penjualan harian/bulanan, analisa menu terlaris (menu engineering), dan statistik performa.',
    estTime: 'Est. 5 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Laporan Standar Operasional</h3>
        <ul class="list-disc pl-5 space-y-2 text-slate-600">
          <li><strong>Daily Sales Summary (Z-Report)</strong>: Rekap total omzet, breakdown tunai/non-tunai, dan tax/service harian.</li>
          <li><strong>Laporan Void & Cancel</strong>: Melacak pembatalan pesanan atau diskon manual oleh kasir untuk audit keamanan.</li>
          <li><strong>Menu Engineering (Star, Dog, Plowhorse, Puzzle)</strong>: Menganalisa menu mana yang paling laris dan menghasilkan margin tertinggi.</li>
        </ul>
      </div>
    `,
  },
  {
    id: 10,
    dbMenuId: 40,
    order: 10,
    title: '10. Tools',
    countText: '6 Panduan & Solusi',
    icon: Printer,
    desc: 'Modul utilitas sistem, konfigurasi hardware printer kasir & dapur LAN/Bluetooth, cash drawer laci uang, dan diagnostik.',
    estTime: 'Est. 6 menit',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Konfigurasi Hardware & Printer Thermal</h3>
        <p class="text-slate-600">Menu Tools digunakan untuk mengatur koneksi perangkat keras kasir:</p>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li><strong>Printer Kasir</strong>: Pengaturan IP address LAN Port 9100 atau port USB thermal printer.</li>
          <li><strong>Printer Dapur & Bar</strong>: Routing cetak otomatis tiket order pesanan berdasarkan kategori makanan/minuman.</li>
          <li><strong>Cash Drawer (Laci Uang)</strong>: Menghubungkan kabel RJ11 ke port DK printer dan mengaktifkan tendangan otomatis.</li>
          <li><strong>Customer Display & Barcode Scanner</strong>: Kalibrasi layar hadap pelanggan dan scanner QRIS.</li>
        </ol>
      </div>
    `,
  },
]

// Troubleshooting Articles for "Solusi" Tab
const troubleshootingArticles = [
  {
    id: 1,
    title: 'Cara Menangani Printer Kasir Thermal yang Tidak Mencetak Struk',
    menuName: '10. Tools',
    readTime: '4 menit baca',
    summary: 'Langkah cepat penanganan printer thermal saat struk tidak keluar pada jam operasional kasir.',
    content: `
      <div class="space-y-4">
        <div class="p-4 bg-amber-50 border border-amber-200 rounded-xl text-amber-900 text-sm">
          <strong>Peringatan Awal:</strong> Jangan mematikan saklar listrik utama kasir saat printer sedang memproses antrian cetak.
        </div>
        <h3 class="text-lg font-bold text-slate-900">1. Periksa Fisik Printer & Lampu Indikator (LED)</h3>
        <ul class="list-disc pl-5 space-y-1 text-slate-600">
          <li><strong>Lampu Merah (Error) Berkedip</strong>: Kemungkinan kertas struk habis atau penutup printer belum tertutup rapat.</li>
          <li><strong>Lampu Biru/Hijau Mati</strong>: Periksa kabel adaptor daya di bagian belakang printer apakah terlepas.</li>
        </ul>
        <h3 class="text-lg font-bold text-slate-900 mt-6">2. Periksa Posisi Kertas Thermal</h3>
        <p>Kertas thermal memiliki satu sisi sensitif panas. Pastikan gulungan kertas menghadap ke atas (posisi kertas keluar dari bawah gulungan).</p>
      </div>
    `,
  },
  {
    id: 2,
    title: 'Solusi Laci Uang (Cash Drawer) Tidak Terbuka Otomatis Saat Bayar',
    menuName: '10. Tools',
    readTime: '3 menit baca',
    summary: 'Pengecekan kabel RJ11 dan pengaturan printer driver agar cash drawer memicu tendangan pembuka otomatis.',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">Penyebab Umum & Solusi</h3>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li><strong>Kabel RJ11 Longgar</strong>: Kabel telepon (RJ11) dari laci uang harus terpasang kencang ke port DK.</li>
          <li><strong>Kunci Manual Terkunci Penuh</strong>: Pastikan lubang kunci laci berada pada posisi vertikal.</li>
        </ol>
      </div>
    `,
  },
  {
    id: 3,
    title: 'Mengatur IP Address Statis Printer Kitchen (LAN Port 9100)',
    menuName: '10. Tools',
    readTime: '5 menit baca',
    summary: 'Panduan konfigurasi IP Address static pada printer dapur LAN agar cetakan order makanan tidak terputus.',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">Langkah Menetapkan IP Statis via Web Browser</h3>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li>Cetak kertas Self-Test printer dengan menahan tombol FEED sambil menyalakan saklar daya.</li>
          <li>Buka browser di komputer kasir, ketik IP printer, lalu ubah mode dari DHCP ke Static IP.</li>
        </ol>
      </div>
    `,
  },
  {
    id: 4,
    title: 'Mengatasi Kertas Struk Macet (Paper Jam) pada Pemotong Otomatis',
    menuName: '10. Tools',
    readTime: '3 menit baca',
    summary: 'Cara mengembalikan pisau auto-cutter yang tersangkut tanpa merusak gear mesin printer.',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">Langkah Menyelamatkan Pisau Auto-Cutter</h3>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li>Matikan saklar printer kasir terlebih dahulu.</li>
          <li>Buka penutup plastik kecil di bagian depan pemotong (gear cover).</li>
          <li>Putar roda gerigi (gear) secara perlahan menggunakan jari hingga bilah pisau mundur ke posisi semula.</li>
        </ol>
      </div>
    `,
  },
  {
    id: 5,
    title: 'Mengalihkan Cetakan Order Dapur Sementara Saat Printer Utama Rusak',
    menuName: '02. Transaksi',
    readTime: '4 menit baca',
    summary: 'Prosedur darurat mengalihkan routing tiket order dapur ke printer kasir saat terjadi gangguan printer kitchen.',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">Fitur Emergency Kitchen Routing</h3>
        <p class="text-slate-600">Buka menu <strong>10. Tools > Hardware Setup</strong>, ubah target printer sementara ke Printer Kasir.</p>
      </div>
    `,
  },
]

// Filtered lists
const filteredModuleCards = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return presto10Modules
  return presto10Modules.filter(c =>
    c.title.toLowerCase().includes(q) || c.desc.toLowerCase().includes(q)
  )
})

const filteredTroubleshooting = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return troubleshootingArticles
  return troubleshootingArticles.filter(a =>
    a.title.toLowerCase().includes(q) || a.summary.toLowerCase().includes(q)
  )
})

// Reader Modal
const readerModalOpen = ref(false)
const activeReaderItem = ref(null)
const feedbackGiven = ref(false)

function openReader(item, type = 'manual') {
  if (type === 'manual') {
    router.push(`/system/${systemId.value}/guide/${item.id}`)
    return
  }
  activeReaderItem.value = { ...item, type }
  feedbackGiven.value = false
  readerModalOpen.value = true
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
    <nav class="border-b border-slate-200/80 bg-white/95 backdrop-blur-sm py-3.5 px-6 sm:px-12 lg:px-20 xl:px-28 sticky top-0 z-40 shadow-xs">
      <div class="max-w-[1600px] mx-auto flex items-center justify-between">
        <!-- Logo, Title & Back Button -->
        <div class="flex items-center gap-3.5">
          <div class="flex items-center cursor-pointer" @click="goBackHome">
            <img src="/img/probus-logo.png" alt="Probus System" class="h-8 sm:h-9 w-auto object-contain" />
          </div>
          <span class="text-slate-300 font-light">|</span>
          <span class="text-slate-500 text-sm font-medium hidden sm:inline">Pusat Bantuan Publik</span>

          <!-- Pill Button: "← Kembali ke Pilihan Sistem" (Matching Design) -->
          <button
            @click="goBackHome"
            class="ml-2 inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-sky-50 hover:bg-sky-100 text-sky-700 font-semibold text-xs border border-sky-100 transition-colors shadow-2xs"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>Kembali ke Pilihan Sistem</span>
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
          <a href="#" class="text-slate-600 hover:text-slate-950 transition-colors">
            Status Layanan
          </a>
          <a href="#" class="text-slate-600 hover:text-slate-950 transition-colors">
            Hubungi Implementator
          </a>
        </div>
      </div>
    </nav>

    <!-- 3. Header Hero Section (Centered with Pill Switcher & Note) -->
    <section class="pt-12 pb-10 px-6 sm:px-12 lg:px-20 xl:px-28 text-center bg-gradient-to-b from-white via-sky-50/20 to-transparent">
      <div class="max-w-4xl mx-auto">
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
        <div class="mt-6 max-w-2xl mx-auto relative group">
          <div class="absolute left-4 top-1/2 -translate-y-1/2 pointer-events-none flex items-center justify-center transition-colors duration-200 text-slate-400 group-focus-within:text-sky-600">
            <Search class="w-4.5 h-4.5" />
          </div>
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="`Ketik fitur, nomor modul, atau topik panduan ${currentSystem.name}...`"
            class="w-full pl-11 sm:pl-12 pr-11 sm:pr-12 py-3 sm:py-3.5 rounded-2xl text-sm bg-white border border-slate-200/90 text-slate-800 placeholder:text-slate-400 font-normal shadow-xs hover:border-slate-300 hover:shadow-sm focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/15 focus:shadow-md transition-all duration-200"
          />
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            type="button"
            class="absolute right-3.5 top-1/2 -translate-y-1/2 p-1.5 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors focus:outline-none"
            title="Hapus pencarian"
          >
            <X class="w-4 h-4" />
          </button>
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
              10 Modul
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
              12 Artikel
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

    <!-- 4. Main Section: 10 Modul Cards Grid (Setup Master hingga Tools) -->
    <main class="flex-1 max-w-[1600px] w-full mx-auto px-6 sm:px-12 lg:px-20 xl:px-28 pb-20">
      <!-- Section Header -->
      <div class="mb-6 flex flex-col sm:flex-row sm:items-end justify-between gap-2">
        <div>
          <h2 class="text-xl sm:text-2xl font-bold text-slate-900 tracking-tight">
            Pilih Modul Sesuai Kebutuhan Anda
          </h2>
          <p class="text-xs sm:text-sm text-slate-500 mt-1">
            Pelajari panduan tata cara penggunaan fitur {{ currentSystem.name }} langkah demi langkah per modul (10 Modul Sistem):
          </p>
        </div>
        <div class="text-xs font-semibold px-3 py-1 rounded-full bg-slate-100 text-slate-600 w-fit">
          {{ filteredModuleCards.length }} dari 10 Modul
        </div>
      </div>

      <!-- TAB 1: 10 Modul Cards Grid (Setup Master s/d Tools) -->
      <div v-if="activeTab === 'manual'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
        <div
          v-for="card in filteredModuleCards"
          :key="card.id"
          @click="openReader(card, 'manual')"
          class="bg-white rounded-2xl p-6 sm:p-7 border border-slate-200/90 hover:border-sky-500 hover:ring-4 hover:ring-sky-500/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200 cursor-pointer flex flex-col justify-between group shadow-xs"
        >
          <div>
            <!-- Card Top: Icon & Count Badge -->
            <div class="flex items-center justify-between">
              <div class="w-11 h-11 rounded-xl flex items-center justify-center bg-slate-100 text-slate-700 group-hover:bg-sky-50 group-hover:text-sky-600 transition-colors shadow-2xs">
                <component :is="card.icon" class="w-5 h-5 stroke-[1.8]" />
              </div>

              <!-- Top Right Badge (Pill) -->
              <span class="text-xs font-semibold px-2.5 py-1 rounded-md tracking-wide bg-slate-100 text-slate-600 group-hover:bg-sky-50 group-hover:text-sky-700 transition-colors">
                {{ card.countText }}
              </span>
            </div>

            <!-- Card Title & Description -->
            <h3 class="font-bold text-slate-900 text-base sm:text-lg mt-5 group-hover:text-sky-600 transition-colors">
              {{ card.title }}
            </h3>
            <p class="text-xs sm:text-sm text-slate-500 mt-2 leading-relaxed line-clamp-3">
              {{ card.desc }}
            </p>
          </div>

          <!-- Card Bottom: Action Link -->
          <div class="pt-5 mt-5 border-t border-slate-100 flex items-center justify-end text-xs sm:text-sm">
            <span class="text-sky-600 font-semibold flex items-center gap-1 group-hover:text-sky-700">
              <span>Buka Manual Book</span>
              <ArrowRight class="w-3.5 h-3.5 group-hover:translate-x-1.5 transition-transform" />
            </span>
          </div>
        </div>
      </div>

      <!-- TAB 2: Troubleshooting Articles List (Solusi Tab) -->
      <div v-else class="space-y-4">
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
              <h3 class="text-base sm:text-lg font-bold text-slate-900 mt-1.5 group-hover:text-emerald-600 transition-colors">
                {{ art.title }}
              </h3>
              <p class="text-xs sm:text-sm text-slate-500 mt-1">
                {{ art.summary }}
              </p>
            </div>
          </div>

          <div class="shrink-0 flex items-center gap-2 text-emerald-600 font-semibold text-xs sm:text-sm">
            <span>Baca Solusi Lengkap</span>
            <ArrowRight class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
          </div>
        </div>
      </div>

      <!-- 5. Bottom Green Help Banner (Hanya muncul di Tab Manual Book) -->
      <div
        v-if="activeTab === 'manual'"
        class="mt-10 sm:mt-12 bg-emerald-50/70 border border-emerald-200/90 rounded-2xl p-5 sm:p-6 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 shadow-xs"
      >
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
    <footer class="border-t border-slate-200 bg-white py-6 px-6 sm:px-12 lg:px-20 xl:px-28 text-xs sm:text-sm text-slate-500">
      <div class="max-w-[1600px] mx-auto flex flex-col sm:flex-row items-center justify-between gap-3">
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
            class="text-slate-700 text-sm sm:text-base leading-relaxed pt-4 border-t border-slate-100"
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
