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
  HelpCircle
} from 'lucide-vue-next'
import apiClient from '../api/client'
import WhatsAppSupportDropdown from '../components/WhatsAppSupportDropdown.vue'

const router = useRouter()
const searchQuery = ref('')
const selectedSystemFilter = ref('all')
const selectedTagFilter = ref('all')
const isLoading = ref(false)

// Real articles matching dbprobuscs articles table
const articles = ref([
  {
    id: 1,
    title: 'Cara Menangani Printer Kasir Thermal yang Tidak Mencetak Struk',
    system: 'PRESTO',
    categoryName: 'Hardware & Printer',
    readTime: '4 menit baca',
    summary: 'Langkah cepat penanganan printer thermal saat struk tidak keluar pada jam operasional kasir.',
    content: `
      <div class="space-y-4">
        <div class="p-4 bg-amber-50 border border-amber-200 rounded-xl text-amber-900 text-sm">
          <strong>Peringatan Awal:</strong> Jangan mematikan saklar listrik utama kasir saat printer sedang memproses antrian cetak.
        </div>
        <h3 class="text-lg font-bold text-slate-900">1. Periksa Fisik Printer & Lampu Indikator (LED)</h3>
        <ul class="list-disc pl-5 space-y-1.5 text-slate-600">
          <li><strong>Lampu Merah (Error) Berkedip</strong>: Kemungkinan kertas struk habis atau penutup printer belum tertutup rapat.</li>
          <li><strong>Lampu Biru/Hijau Mati</strong>: Periksa kabel adaptor daya di bagian belakang printer apakah terlepas.</li>
        </ul>
        <h3 class="text-lg font-bold text-slate-900 mt-6">2. Periksa Posisi Kertas Thermal</h3>
        <p>Kertas thermal memiliki satu sisi sensitif panas. Pastikan gulungan kertas menghadap ke atas (posisi kertas keluar dari bawah gulungan). Jika terbalik, printer tetap berjalan namun hasil cetak akan polos/kosong.</p>
        <h3 class="text-lg font-bold text-slate-900 mt-6">3. Restart Service Spooler & Test Ping IP</h3>
        <p>Jika printer menggunakan kabel LAN (Ethernet), buka CMD di komputer kasir dan ketik <code>ping 192.168.1.200</code> (sesuaikan dengan IP printer). Jika respons <em>Reply</em>, buka PRESTO POS dan lakukan test print ulang.</p>
      </div>
    `,
  },
  {
    id: 2,
    title: 'Mengatur IP Address Statis Printer Kitchen (LAN Port 9100)',
    system: 'PRESTO',
    categoryName: 'Jaringan & Network',
    readTime: '5 menit baca',
    summary: 'Panduan konfigurasi IP Address static pada printer dapur LAN agar cetakan order makanan tidak terputus.',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">1. Alasan Menggunakan IP Statis</h3>
        <p class="text-slate-600">Printer dapur harus memiliki alamat IP yang tetap (statis). Jika menggunakan DHCP otomatis dari router, alamat IP printer dapat berubah sewaktu-waktu sehingga pesanan dari kasir gagal tercetak di dapur.</p>
        
        <h3 class="text-lg font-bold text-slate-900 mt-6">2. Langkah Menetapkan IP Statis via Web Browser</h3>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li>Cetak kertas Self-Test printer dengan menahan tombol <em>FEED</em> sambil menyalakan saklar daya. Catat alamat IP default printer.</li>
          <li>Buka Google Chrome di komputer kasir yang berada di jaringan LAN yang sama, lalu ketik alamat IP printer tersebut pada address bar.</li>
          <li>Login ke menu pengaturan printer (Default user: <em>admin</em>, pass: <em>admin</em> atau <em>12345678</em>).</li>
          <li>Masuk ke menu <strong>Network Configuration > TCP/IP</strong>, ubah mode dari <em>DHCP</em> menjadi <em>Static IP</em>.</li>
          <li>Masukkan IP di luar rentang pool DHCP router (misal: <code>192.168.1.210</code>, Subnet: <code>255.255.255.0</code>, Gateway: <code>192.168.1.1</code>).</li>
          <li>Klik <strong>Save & Reboot</strong>. Update juga IP baru ini pada menu <em>Tools > Hardware Setup</em> di aplikasi PRESTO.</li>
        </ol>
      </div>
    `,
  },
  {
    id: 3,
    title: 'Solusi Laci Uang (Cash Drawer) Tidak Terbuka Otomatis Saat Bayar',
    system: 'PRESTO',
    categoryName: 'Hardware & Kasir',
    readTime: '3 menit baca',
    summary: 'Pengecekan kabel RJ11 dan pengaturan printer driver agar cash drawer memicu tendangan pembuka otomatis.',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">Penyebab Umum & Solusi</h3>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li><strong>Kabel RJ11 Longgar</strong>: Kabel telepon (RJ11) dari laci uang harus terpasang kencang ke port <em>DK (Drawer Kick)</em> di belakang printer thermal kasir.</li>
          <li><strong>Kunci Manual Terkunci Penuh</strong>: Pastikan lubang kunci di bagian depan laci uang berada pada posisi vertikal (otomatis). Jika berada di posisi horizontal, laci terkunci manual dan tidak dapat ditendang secara elektrik.</li>
          <li><strong>Driver Cash Drawer Kick</strong>: Buka <em>Device and Printers > Printer Properties > Device Settings</em>, pastikan opsi <em>Cash Drawer #1: Open Before Printing</em> diaktifkan.</li>
        </ol>
      </div>
    `,
  },
  {
    id: 4,
    title: 'Mengatasi Kertas Struk Macet (Paper Jam) pada Pemotong Otomatis',
    system: 'PRESTO',
    categoryName: 'Hardware & Printer',
    readTime: '3 menit baca',
    summary: 'Cara mengembalikan pisau auto-cutter yang tersangkut tanpa merusak gear mesin printer.',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">Langkah Menyelamatkan Pisau Auto-Cutter</h3>
        <p class="text-slate-600">Jika penutup printer tidak bisa dibuka karena pisau pemotong tersangkut di tengah:</p>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li>Matikan saklar printer kasir terlebih dahulu.</li>
          <li>Buka penutup plastik kecil di bagian depan pemotong (gear cover).</li>
          <li>Putar roda gerigi (gear) secara perlahan menggunakan jari hingga bilah pisau pemotong kembali mundur ke posisi semula.</li>
          <li>Buka tuas penutup printer, bersihkan sisa robekan kertas yang tersangkut, lalu tutup kembali dengan kencang hingga berbunyi 'klik'.</li>
        </ol>
      </div>
    `,
  },
  {
    id: 5,
    title: 'Mengalihkan Cetakan Order Dapur Sementara Saat Printer Utama Rusak',
    system: 'PRESTO',
    categoryName: 'Operasional & Transaksi',
    readTime: '4 menit baca',
    summary: 'Prosedur darurat mengalihkan routing tiket order dapur ke printer kasir saat terjadi gangguan printer kitchen.',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">Fitur Emergency Kitchen Routing</h3>
        <p class="text-slate-600">Saat jam sibuk makan siang atau makan malam dan printer dapur mendadak rusak atau habis kertas:</p>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li>Buka menu <strong>Tools > Hardware & Printer Setup</strong>.</li>
          <li>Pada kolom <em>Routing Kitchen Printer</em>, ubah target printer sementara ke <strong>Printer Kasir (Cashier Thermal)</strong>.</li>
          <li>Simpan pengaturan. Seluruh tiket order makanan akan tercetak ganda di meja kasir sehingga kasir dapat mengantarkan tiket fisik secara manual ke chef dapur.</li>
        </ol>
      </div>
    `,
  },
  {
    id: 6,
    title: 'Cara Menghubungkan Printer Thermal Bluetooth ke Tablet Kasir Android',
    system: 'PPOS',
    categoryName: 'Hardware & Kasir',
    readTime: '4 menit baca',
    summary: 'Panduan pairing Bluetooth dan konfigurasi nomor PIN 0000/1234 pada aplikasi PPOS Tablet.',
    content: `
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-slate-900">Langkah Pairing Bluetooth</h3>
        <ol class="list-decimal pl-5 space-y-2 text-slate-600">
          <li>Nyalakan printer thermal mobile bluetooth.</li>
          <li>Buka pengaturan <em>Settings > Bluetooth</em> pada tablet Android.</li>
          <li>Cari nama perangkat printer (contoh: <em>RPP02N</em>, <em>MPT-II</em>, atau <em>Bluetooth Printer</em>).</li>
          <li>Masukkan kode PIN pairing (biasanya <code>0000</code> atau <code>1234</code>).</li>
          <li>Buka aplikasi PPOS, masuk ke menu <strong>Pengaturan > Printer Bluetooth</strong>, lalu pilih printer yang telah terhubung. Lakukan test print.</li>
        </ol>
      </div>
    `,
  },
])

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
    const matchSystem = selectedSystemFilter.value === 'all' || art.system === selectedSystemFilter.value
    const matchQuery = !searchQuery.value.trim() ||
      art.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      art.summary.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      art.categoryName.toLowerCase().includes(searchQuery.value.toLowerCase())
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
          <a href="#" class="text-slate-600 hover:text-slate-950 transition-colors">
            Status Layanan
          </a>
          <a href="#" class="text-slate-600 hover:text-slate-950 transition-colors">
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
              'px-3.5 py-1.5 rounded-lg text-xs font-semibold transition-all',
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
            Menampilkan {{ filteredArticles.length }} artikel pemecahan masalah teknis
          </p>
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
      <div v-if="filteredArticles.length === 0" class="text-center py-16 bg-white rounded-2xl border border-slate-200">
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
        <p>© 2026 CV Mitra Solusi Mandiri • Bank Solusi Resmi</p>
      </div>
    </footer>

    <!-- 6. Interactive Reader Modal -->
    <div
      v-if="readerModalOpen && activeArticle"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 backdrop-blur-xs p-4 sm:p-6"
    >
      <div class="bg-white w-full max-w-3xl max-h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-in fade-in duration-200">
        <div class="p-5 border-b border-slate-200 flex items-center justify-between bg-slate-50/50">
          <div class="flex items-center gap-2.5">
            <span class="text-[11px] font-bold px-2.5 py-0.5 rounded-full uppercase tracking-wider bg-amber-100 text-amber-800">
              Artikel Solusi
            </span>
            <span class="text-xs text-slate-500 font-medium">• {{ activeArticle.system }}</span>
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
            {{ activeArticle.title }}
          </h2>

          <div
            class="text-slate-700 text-sm sm:text-base leading-relaxed pt-4 border-t border-slate-100"
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
              class="px-3.5 py-1.5 rounded-lg border border-slate-300 bg-white hover:bg-slate-100 text-slate-700 font-medium inline-flex items-center gap-1.5 transition-colors"
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
              class="px-4 py-1.5 rounded-lg bg-slate-900 text-white font-semibold hover:bg-slate-800 transition-colors"
            >
              Tutup
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

