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
  ExternalLink
} from 'lucide-vue-next'
import apiClient from '../api/client'

const route = useRoute()
const router = useRouter()

// Selected System ID (default: 2 for PRESTO)
const selectedSystemId = ref(Number(route.params.systemId) || 2)

// Selected Module ID (default: 2 for Transaksi / Order Taker)
const selectedModuleId = ref(Number(route.params.moduleId) || 2)

// Systems list for the "PILIH SISTEM" pill switcher
const availableSystems = [
  { id: 2, name: 'PRESTO', title: 'PRESTO Resto POS (10 Modul)', shortName: 'PRESTO Resto POS' },
  { id: 1, name: 'PHIS', title: 'Probus Hotel PMS (10 Modul)', shortName: 'Probus Hotel PMS' },
  { id: 8, name: 'PPE', title: 'PPE 2.0 Management (10 Modul)', shortName: 'PPE 2.0 Management' },
  { id: 3, name: 'PPOS', title: 'PPOS Point of Sales (10 Modul)', shortName: 'PPOS Point of Sales' },
  { id: 5, name: 'PR-PO ONLINE', title: 'Back Office & Akuntansi (10 Modul)', shortName: 'Back Office & Akuntansi' },
]

const currentSystem = computed(() => {
  return availableSystems.find(s => s.id === selectedSystemId.value) || availableSystems[0]
})

// 10 Real Modules of PRESTO matching dbprobuscs `menu_systems`
const presto10ModulesData = [
  {
    id: 1,
    order: 1,
    code: '01',
    title: '01. Setup Master',
    subtitle: 'Panduan Modul',
    icon: Settings,
    lead: 'Panduan konfigurasi data master awal restoran: master meja, kategori menu, item makanan & minuman, harga jual, serta penentuan pajak dan service charge.',
    intro: 'Panduan langkah demi langkah bagi tim manajerial atau admin kasir untuk mengelola data master pada sistem PRESTO Resto POS.',
    sections: [
      {
        heading: '1. Pengaturan Data Meja & Denah Restoran',
        desc: 'Buka menu Setup Master > Data Meja. Tentukan area (Indoor, Outdoor, Bar) dan kapasitas kursi per meja sebelum membuka sesi kasir.',
        diagram: {
          filename: 'setup_denah_meja.png (480 × 160 px)',
          items: [
            '[1] Form Tambah Meja: Nomor Meja & Area',
            '[2] Tombol Simpan Konfigurasi Denah Meja →',
          ],
        },
      },
      {
        heading: '2. Input Kategori & Master Menu Makanan / Minuman',
        desc: 'Tambahkan kategori menu utama (Appetizer, Main Course, Beverage), tentukan harga dasar, serta printer tujuan pesanan (Dapur Panas, Bar Minuman).',
        diagram: {
          filename: 'master_menu_item.png (480 × 180 px)',
          items: [
            '[1] Nama Item: Nasi Goreng Spesial — Kategori: Main Course',
            '[2] Harga: Rp 35.000 | Tax: 10% | Service: 5%',
            '[3] Target Cetak Pesanan: Printer Kitchen Dapur 01',
          ],
        },
      },
      {
        heading: '3. Aturan & Kebijakan Tarif Pajak / Service',
        desc: 'Berikut adalah acuan parameter tarif yang diterapkan otomatis pada saat pencetakan bill transaksi tamu:',
        table: {
          headers: ['No', 'Komponen Biaya', 'Keterangan Petunjuk Penggunaan'],
          rows: [
            ['1', 'PB1 / Pajak Resto (10%)', 'Dihitung otomatis dari total konsumsi makanan & minuman sebelum diskon.'],
            ['2', 'Service Charge (5%)', 'Biaya layanan operasional restoran, dihitung dari subtotal pemesanan.'],
            ['3', 'Diskon Promosi / Member', 'Dipotong sebelum pengenaan pajak sesuai kebijakan manajemen.'],
          ],
        },
      },
    ],
  },
  {
    id: 2,
    order: 2,
    code: '02',
    title: '02. Order Taker & Transaksi Meja',
    subtitle: 'Sedang Dibaca • Aktif',
    icon: UtensilsCrossed,
    lead: 'Panduan langkah demi langkah bagi kasir & pramusaji untuk membuka meja tamu, menginput pesanan menu beserta catatan dapur (notes), hingga memindahkan atau menggabungkan meja aktif.',
    intro: 'Panduan langkah demi langkah bagi kasir dan pramusaji untuk mengoperasikan modul Order Taker & Transaksi Meja pada sistem PRESTO Resto POS.',
    sections: [
      {
        heading: '1. Membuka Meja Baru (Open Table) & Mengisi Jumlah Tamu',
        desc: 'Pilih nomor meja yang berstatus Kosong (warna hijau) pada denah restoran, lalu tekan tombol "Buka Meja" untuk mulai mencatat pesanan tamu:',
        diagram: {
          filename: 'denah_meja_resto.png (480 × 160 px)',
          items: [
            '[1] Meja 04 • Status: Kosong (Hijau)',
            '[2] Tombol Buka Meja 04 (Input Pesanan) →',
          ],
        },
      },
      {
        heading: '2. Memilih Menu Pesanan, Catatan Dapur & Kirim Order',
        desc: 'Pilih kategori makanan atau minuman, tentukan jumlah porsi, lalu gunakan kolom Catatan Dapur (Notes) untuk instruksi khusus sebelum mengirim pesanan ke koki:',
        diagram: {
          filename: 'input_pesanan_menu.png (480 × 180 px)',
          items: [
            '[1] Nasi Goreng Spesial (x2) — Rp 70.000',
            '[2] Catatan Dapur: 1 Tidak Pedas, 1 Less Oil',
            '[3] Tombol Kirim ke Dapur (Send Order)',
          ],
        },
      },
      {
        heading: '3. Keterangan Status Meja & Panduan Pindah Meja (Move Table)',
        desc: 'Berikut adalah acuan status warna meja pada denah restoran dan tata cara pindah meja jika tamu berpindah tempat duduk:',
        table: {
          headers: ['No', 'Status / Aksi', 'Keterangan Petunjuk Penggunaan'],
          rows: [
            ['1', '🟢 Meja Kosong (Hijau)', 'Meja belum terisi tamu, siap dibuka untuk transaksi tamu baru.'],
            ['2', '🔴 Meja Terisi (Merah)', 'Meja aktif sedang makan dan pesanan telah terkirim ke dapur.'],
            ['3', '🔄 Pindah Meja (Move Table)', 'Buka meja aktif, klik Pindah / Gabung Meja di sudut kanan atas, lalu pilih nomor meja tujuan.'],
          ],
        },
      },
    ],
  },
  {
    id: 3,
    order: 3,
    code: '03',
    title: '03. Outlet Activity',
    subtitle: 'Panduan Modul',
    icon: Store,
    lead: 'Monitoring aktivitas operasional harian outlet: void pesanan, pembatalan item, pemindahan meja (move table), dan penggabungan tagihan (join bill).',
    intro: 'Panduan pengawasan aktivitas operasional outlet bagi supervisor dan head waiter PRESTO Resto POS.',
    sections: [
      {
        heading: '1. Otorisasi Void Item / Pembatalan Pesanan',
        desc: 'Item yang sudah terkirim ke dapur hanya dapat dibatalkan melalui persetujuan Supervisor dengan PIN otorisasi.',
        diagram: {
          filename: 'void_item_approval.png (480 × 160 px)',
          items: [
            '[1] Pilih Item Pesanan yang Dibatalkan',
            '[2] Masukkan Alasan Pembatalan (Salah Input / Tamu Cancel)',
            '[3] Input Supervisor PIN Otorisasi',
          ],
        },
      },
      {
        heading: '2. Penggabungan Bill (Join Bill) Antar Meja',
        desc: 'Jika tamu dua meja ingin membayar dalam satu tagihan bersama, gunakan fitur Join Bill pada meja utama.',
        diagram: {
          filename: 'join_table_bill.png (480 × 160 px)',
          items: [
            '[1] Meja Asal: Meja 02 + Meja Gabungan: Meja 03',
            '[2] Subtotal Gabungan Dihitung Otomatis',
            '[3] Konfirmasi Cetak Bill Gabungan',
          ],
        },
      },
      {
        heading: '3. Rekap Aktivitas Operasional Meja',
        desc: 'Tabel acuan penanganan khusus operasional meja di outlet resto:',
        table: {
          headers: ['No', 'Tipe Aktivitas', 'Keterangan Petunjuk Penggunaan'],
          rows: [
            ['1', 'Void Before Send', 'Membatalkan item sebelum tombol Kirim Dapur ditekan (tidak perlu PIN).'],
            ['2', 'Void After Send', 'Item sudah tercetak di dapur; wajib memasukkan PIN supervisor & alasan.'],
            ['3', 'Split Bill Item', 'Memecah pesanan menjadi beberapa tagihan terpisah per tamu.'],
          ],
        },
      },
    ],
  },
  {
    id: 4,
    order: 4,
    code: '04',
    title: '04. Cost Control',
    subtitle: 'Panduan Modul',
    icon: PieChart,
    lead: 'Pengendalian biaya bahan baku resep (HPP/COGS), pemakaian bahan per porsi, dan pemantauan selisih stok bahan dapur (waste management).',
    intro: 'Panduan tata cara input resep menu dan analisis biaya bahan baku pada modul Cost Control.',
    sections: [
      {
        heading: '1. Penyusunan Resep (Recipe / Bill of Materials)',
        desc: 'Tautkan setiap menu penjualan dengan bahan baku gudang untuk memotong stok bahan secara otomatis setiap porsi terjual.',
        diagram: {
          filename: 'recipe_cogs_setup.png (480 × 160 px)',
          items: [
            '[1] Menu: Nasi Goreng Spesial',
            '[2] Komposisi: Beras (150 gr), Telur (1 btr), Daging Ayam (50 gr)',
            '[3] Kalkulasi HPP Otomatis: Rp 14.200 per porsi',
          ],
        },
      },
      {
        heading: '2. Pencatatan Pembuangan Bahan (Spoilage & Waste)',
        desc: 'Catat bahan makanan yang rusak atau kedaluwarsa agar kartu stok bahan di dapur tetap akurat.',
        diagram: {
          filename: 'waste_logging_form.png (480 × 160 px)',
          items: [
            '[1] Form Spoilage / Waste Gudang Dapur',
            '[2] Pilih Alasan: Kedaluwarsa, Rusak Penyimpanan, Uji Rasa',
          ],
        },
      },
      {
        heading: '3. Standar Rasio Biaya Makanan (Food Cost Ratio)',
        desc: 'Panduan indikator kesehatan food cost restoran:',
        table: {
          headers: ['No', 'Kategori Menu', 'Target Standar Food Cost'],
          rows: [
            ['1', 'Makanan Utama (Food)', '28% - 35% dari harga jual menu.'],
            ['2', 'Minuman (Beverage)', '15% - 22% dari harga jual menu.'],
            ['3', 'Pastry & Dessert', '20% - 28% dari harga jual menu.'],
          ],
        },
      },
    ],
  },
  {
    id: 5,
    order: 5,
    code: '05',
    title: '05. Purcashing AP',
    subtitle: 'Panduan Modul',
    icon: ShoppingBag,
    lead: 'Alur pengadaan barang: pembuatan Purchase Order (PO), penerimaan barang gudang (receiving), verifikasi faktur supplier, dan buku hutang usaha.',
    intro: 'Panduan proses pengadaan bahan baku dari supplier hingga pencatatan hutang dagang pada PRESTO Resto POS.',
    sections: [
      {
        heading: '1. Pembuatan Pesanan Pembelian (Purchase Order - PO)',
        desc: 'Pilih supplier terdaftar, masukkan daftar bahan baku yang dibutuhkan outlet, dan tetapkan tanggal perkiraan pengiriman.',
        diagram: {
          filename: 'po_creation_slip.png (480 × 160 px)',
          items: [
            '[1] Supplier: PT Pangan Segar Nusantara',
            '[2] Daftar Barang: Daging Sapi 20 kg, Minyak Goreng 50 L',
            '[3] Cetak Dokumen PO & Kirim ke Supplier',
          ],
        },
      },
      {
        heading: '2. Penerimaan Barang di Gudang (Receiving Log)',
        desc: 'Cek fisik barang yang datang sesuai Surat Jalan supplier, sesuaikan jumlah aktual yang diterima sebelum menyimpan ke stok gudang.',
        diagram: {
          filename: 'receiving_item_checklist.png (480 × 160 px)',
          items: [
            '[1] Cek Berat / Jumlah Fisik vs Faktur',
            '[2] Validasi Tanggal Expired Bahan Masuk',
          ],
        },
      },
      {
        heading: '3. Verifikasi Status Tagihan Hutang (AP Aging)',
        desc: 'Pantau jatuh tempo faktur supplier untuk memastikan pembayaran tepat waktu:',
        table: {
          headers: ['No', 'Status Tempo', 'Tindakan Bagian Pembelian / Keuangan'],
          rows: [
            ['1', '0 - 15 Hari', 'Faktur baru, verifikasi kelengkapan nota dan surat jalan.'],
            ['2', '16 - 30 Hari', 'Jadwalkan batch pembayaran mingguan via General Cashier.'],
            ['3', '> 30 Hari (Overdue)', 'Prioritas pembayaran segera untuk menghindari penghentian suplai.'],
          ],
        },
      },
    ],
  },
  {
    id: 6,
    order: 6,
    code: '06',
    title: '06. AR (Account Receivable)',
    subtitle: 'Panduan Modul',
    icon: Receipt,
    lead: 'Manajemen piutang pelanggan korporat: penagihan faktur banquet, city ledger, voucher makan corporate, dan pencatatan penerimaan pelunasan piutang.',
    intro: 'Tata cara pencatatan transaksi kredit dan penagihan piutang tamu pada modul AR.',
    sections: [
      {
        heading: '1. Transaksi City Ledger / Tagihan Korporat',
        desc: 'Saat transaksi kasir dipilih metode bayar City Ledger, tagihan otomatis dialihkan ke buku besar AR atas nama instansi terkait.',
        diagram: {
          filename: 'city_ledger_posting.png (480 × 160 px)',
          items: [
            '[1] Kasir Memilih Payment Method: City Ledger',
            '[2] Pilih Akun Perusahaan: PT Bank Mega Perkasa',
            '[3] Validasi Voucher / Surat Jaminan Tamu',
          ],
        },
      },
      {
        heading: '2. Pencetakan Invoice Penagihan Resmi',
        desc: 'Kumpulkan seluruh bill konsumsi selama periode acara corporate, lalu cetak Faktur Gabungan (Official Invoice) lengkap dengan lampiran.',
        diagram: {
          filename: 'invoice_ar_statement.png (480 × 160 px)',
          items: [
            '[1] Generate Billing Statement Periode Acara',
            '[2] Tombol Cetak Invoice Resmi Penagihan →',
          ],
        },
      },
      {
        heading: '3. Penutupan Piutang & Penerimaan Transfer',
        desc: 'Tabel acuan status penagihan piutang pelanggan:',
        table: {
          headers: ['No', 'Status Piutang', 'Petunjuk Operasional'],
          rows: [
            ['1', 'Unbilled', 'Transaksi kasir sudah selesai, invoice belum dikirim ke klien.'],
            ['2', 'Billed / Sent', 'Faktur resmi sudah diterima klien, menunggu jatuh tempo pembayaran.'],
            ['3', 'Settled / Paid', 'Dana telah masuk ke rekening bank, lakukan posting pelunasan AR.'],
          ],
        },
      },
    ],
  },
  {
    id: 7,
    order: 7,
    code: '07',
    title: '07. General Cashier',
    subtitle: 'Panduan Modul',
    icon: Wallet,
    lead: 'Sentralisasi kasir utama: penerimaan setoran kasir harian, pengelolaan kas kecil (petty cash), pembayaran operasional outlet, dan rekonsiliasi bank.',
    intro: 'Panduan pengelolaan kas harian, setoran brankas, dan pencairan petty cash pada General Cashier.',
    sections: [
      {
        heading: '1. Penerimaan Uang Setoran Kasir Harian (Drop Cash)',
        desc: 'Verifikasi amplop setoran uang tunai dari kasir outlet yang telah ditutup shift, cocokkan dengan laporan Z-Report sistem.',
        diagram: {
          filename: 'drop_cash_verification.png (480 × 160 px)',
          items: [
            '[1] Terima Amplop Kasir Shift 1 & Shift 2',
            '[2] Input Nominal Uang Fisik vs Z-Report POS',
            '[3] Simpan ke Brankas Utama (Safe Deposit)',
          ],
        },
      },
      {
        heading: '2. Pengeluaran Kas Kecil Operasional (Petty Cash)',
        desc: 'Keluarkan dana kas kecil hanya untuk kebutuhan darurat outlet yang disertai bukti struk pengeluaran valid.',
        diagram: {
          filename: 'petty_cash_voucher.png (480 × 160 px)',
          items: [
            '[1] Form Voucher Kas Kecil (Biaya Gas Elpiji, Es Kristal Darurat)',
            '[2] Tanda Tangan Penerima & Approval Manager',
          ],
        },
      },
      {
        heading: '3. Rekonsiliasi Kas Harian',
        desc: 'Panduan batas toleransi selisih kas fisik:',
        table: {
          headers: ['No', 'Kondisi Kas', 'Prosedur Wajib'],
          rows: [
            ['1', 'Balance (Tepat)', 'Posting setoran ke rekening bank penampung.'],
            ['2', 'Shortage (Kurang)', 'Wajib dibuat berita acara selisih kas oleh kasir yang bertugas.'],
            ['3', 'Overage (Lebih)', 'Catat sebagai pendapatan lain-lain, telusuri kemungkinan salah kembalian.'],
          ],
        },
      },
    ],
  },
  {
    id: 8,
    order: 8,
    code: '08',
    title: '08. Accounting System',
    subtitle: 'Panduan Modul',
    icon: BarChart3,
    lead: 'Buku besar dan pembukuan otomatis: jurnal umum, neraca saldo, laporan laba rugi (P&L), dan neraca lajur keuangan restoran.',
    intro: 'Panduan akuntansi menyeluruh yang terintegrasi otomatis dari modul penjualan, persediaan, dan kas.',
    sections: [
      {
        heading: '1. Pemetaan Akun COA (Chart of Accounts)',
        desc: 'Seluruh penjualan menu kasir otomatis memicu jurnal debet Kas/Bank/Piutang dan kredit Pendapatan Restoran serta Hutang Pajak.',
        diagram: {
          filename: 'auto_journal_sales.png (480 × 160 px)',
          items: [
            '[1] Jurnal Otomatis Penjualan Resto Hari Ini',
            '[2] Dr. Kas Kasir (Rp 5.250.000) | Cr. Penjualan Makanan (Rp 4.500.000)',
            '[3] Cr. Hutang Pajak PB1 (Rp 500.000) | Cr. Service Charge (Rp 250.000)',
          ],
        },
      },
      {
        heading: '2. Tutup Buku Akhir Bulan (Period Closing)',
        desc: 'Lakukan posting jurnal depresiasi aset dan penyesuaian persediaan akhir sebelum mengunci periode buku.',
        diagram: {
          filename: 'month_end_closing.png (480 × 160 px)',
          items: [
            '[1] Validasi Seluruh Jurnal Transaksi Posted',
            '[2] Kunci Periode Akuntansi Bulan Berjalan',
          ],
        },
      },
      {
        heading: '3. Laporan Keuangan Inti',
        desc: 'Dokumen finansial yang dihasilkan otomatis oleh modul Accounting:',
        table: {
          headers: ['No', 'Jenis Laporan', 'Keterangan Analisis'],
          rows: [
            ['1', 'Laba Rugi (Income Statement)', 'Menampilkan pendapatan bersih, HPP, beban operasional, dan laba bersih.'],
            ['2', 'Neraca (Balance Sheet)', 'Posisi aset, liabilitas hutang supplier, dan ekuitas modal pemilik.'],
            ['3', 'Arus Kas (Cash Flow)', 'Mutasi penerimaan dan pengeluaran kas operasional dan investasi.'],
          ],
        },
      },
    ],
  },
  {
    id: 9,
    order: 9,
    code: '09',
    title: '09. Report',
    subtitle: 'Panduan Modul',
    icon: TrendingUp,
    lead: 'Pusat laporan analitik eksekutif: rekap penjualan kasir, item menu terlaris (top sales), analisis jam sibuk (hourly sales), dan laporan pajak daerah.',
    intro: 'Panduan penarikan data laporan penjualan, performa menu, dan rekapitulasi pajak pada PRESTO Resto POS.',
    sections: [
      {
        heading: '1. Laporan Penjualan Harian per Kasir (Daily Cashier Report)',
        desc: 'Tarik laporan ringkasan per kasir untuk memverifikasi metode pembayaran (Tunai, Kartu EDC, QRIS, Transfer) dan diskon yang diberikan.',
        diagram: {
          filename: 'daily_sales_summary.png (480 × 160 px)',
          items: [
            '[1] Filter: Tanggal Hari Ini | Outlet: Resto Utama',
            '[2] Ringkasan Penjualan Bersih: Rp 12.450.000 (142 Transaksi)',
            '[3] Export ke PDF / Excel untuk Laporan Owner',
          ],
        },
      },
      {
        heading: '2. Analisis Menu Terlaris (Menu Engineering)',
        desc: 'Pelajari peringkat menu makanan dan minuman berdasarkan volume penjualan dan profit margin untuk menyusun strategi promo.',
        diagram: {
          filename: 'top_menu_ranking.png (480 × 160 px)',
          items: [
            '[1] Top 1: Nasi Goreng Spesial (320 porsi)',
            '[2] Top 2: Es Kopi Susu Aren (280 gelas)',
            '[3] Top 3: Ayam Bakar Madu (195 porsi)',
          ],
        },
      },
      {
        heading: '3. Kategori Laporan Siap Ekspor',
        desc: 'Daftar jenis laporan yang dapat diunduh langsung:',
        table: {
          headers: ['No', 'Nama Laporan', 'Format Tersedia'],
          rows: [
            ['1', 'Laporan Omzet Kasir per Jam', 'PDF, Excel (.xlsx), Direct Thermal Print.'],
            ['2', 'Laporan Pembayaran EDC & QRIS', 'Excel untuk rekonsiliasi mutasi rekening bank.'],
            ['3', 'Laporan Pajak Daerah PB1', 'Format standar siap lampir Surat Pemberitahuan Pajak Daerah.'],
          ],
        },
      },
    ],
  },
  {
    id: 10,
    order: 10,
    code: '10',
    title: '10. Tools',
    subtitle: 'Panduan Modul',
    icon: Sliders,
    lead: 'Utilitas konfigurasi teknis: pengaturan printer kasir thermal LAN/USB, IP address hardware, auto cut, barcode scanner, dan backup database lokal.',
    intro: 'Panduan konfigurasi teknis perangkat keras kasir dan pemeliharaan database lokal.',
    sections: [
      {
        heading: '1. Konfigurasi Printer Kasir LAN / USB / Bluetooth',
        desc: 'Buka Tools > Printer Setup. Masukkan IP Address printer kasir dan printer dapur, lakukan uji cetak (Test Print) untuk memastikan koneksi stabil.',
        diagram: {
          filename: 'printer_hardware_setup.png (480 × 160 px)',
          items: [
            '[1] IP Address Printer Dapur: 192.168.1.200 (Port 9100)',
            '[2] IP Address Printer Kasir: 192.168.1.201 (Port 9100)',
            '[3] Tombol Test Print Thermal Receipt →',
          ],
        },
      },
      {
        heading: '2. Pengaturan Laci Uang (Cash Drawer Kick)',
        desc: 'Atur perintah trigger cash drawer agar laci terbuka otomatis saat transaksi pembayaran tunai diselesaikan kasir.',
        diagram: {
          filename: 'cash_drawer_command.png (480 × 160 px)',
          items: [
            '[1] Trigger Command: ESC p 0 25 250 (Standar ESC/POS)',
            '[2] Uji Buka Laci Kasir Secara Manual',
          ],
        },
      },
      {
        heading: '3. Troubleshooting Masalah Koneksi Printer',
        desc: 'Langkah darurat saat printer dapur tidak merespons:',
        table: {
          headers: ['No', 'Gejala Gangguan', 'Solusi Perbaikan Mandiri'],
          rows: [
            ['1', 'Lampu Error Merah Berkedip', 'Periksa kertas struk, pastikan kertas terpasang kencang dan tutup rapat.'],
            ['2', 'Kabel LAN Lepas / Timeout', 'Pastikan kabel LAN terpasang erat di router switch hub kasir.'],
            ['3', 'Printer Dapur Macet', 'Aktifkan fitur Emergency Kitchen Routing di menu Tools untuk alihkan cetak ke printer kasir.'],
          ],
        },
      },
    ],
  },
]

// Current module list
const modulesList = ref(presto10ModulesData)

// Active Module object
const currentModule = computed(() => {
  return modulesList.value.find(m => m.id === selectedModuleId.value) || modulesList.value[1]
})

// Previous and Next module computed
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

function selectModule(mod) {
  selectedModuleId.value = mod.id
  router.replace({ params: { systemId: selectedSystemId.value, moduleId: mod.id } })
  window.scrollTo({ top: 120, behavior: 'smooth' })
}

function switchSystem(sys) {
  selectedSystemId.value = sys.id
  selectedModuleId.value = 1
  router.replace({ params: { systemId: sys.id, moduleId: 1 } })
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

// Sync route parameters
watch(
  () => route.params,
  newParams => {
    if (newParams.systemId) {
      selectedSystemId.value = Number(newParams.systemId)
    }
    if (newParams.moduleId) {
      selectedModuleId.value = Number(newParams.moduleId)
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
        <div class="flex items-center gap-1.5">
          <span>Butuh bantuan cepat? Hubungi WhatsApp Support:</span>
          <a
            href="https://wa.me/62811389211"
            target="_blank"
            class="text-slate-800 font-semibold hover:text-sky-600 transition-colors"
          >
            +62 811-389-211
          </a>
        </div>
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

          <!-- Pill Button: "← Kembali ke Pilihan Modul" (Exact Match to Design) -->
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
          <a href="#" class="text-slate-600 hover:text-slate-950 transition-colors">
            Status Layanan
          </a>
          <a href="#" class="text-slate-600 hover:text-slate-950 transition-colors">
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
        
        <!-- LEFT COLUMN: Sidebar "DAFTAR MODUL (10)" -->
        <aside class="w-full md:w-80 lg:w-[340px] shrink-0">
          <div class="bg-slate-100/70 border border-slate-200/80 rounded-2xl p-4 sm:p-5 sticky top-24">
            <div class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-3.5 px-1">
              DAFTAR MODUL ({{ modulesList.length }})
            </div>

            <!-- List of 10 Modules -->
            <div class="space-y-2.5 max-h-[calc(100vh-220px)] overflow-y-auto pr-1">
              <div
                v-for="mod in modulesList"
                :key="mod.id"
                @click="selectModule(mod)"
                :class="[
                  'rounded-xl p-4 border transition-all duration-150 cursor-pointer text-left',
                  selectedModuleId === mod.id
                    ? 'bg-white border-2 border-sky-500 ring-2 ring-sky-100 shadow-sm'
                    : 'bg-white border-slate-200/90 hover:border-slate-300 hover:shadow-2xs'
                ]"
              >
                <div class="flex items-center justify-between">
                  <h3
                    :class="[
                      'text-sm font-bold tracking-tight',
                      selectedModuleId === mod.id ? 'text-sky-700 font-extrabold' : 'text-slate-900'
                    ]"
                  >
                    {{ mod.title }}
                  </h3>
                </div>
                <div
                  :class="[
                    'text-[11px] mt-1 font-medium',
                    selectedModuleId === mod.id ? 'text-sky-600 font-semibold' : 'text-slate-400'
                  ]"
                >
                  <span v-if="selectedModuleId === mod.id">Sedang Dibaca • Aktif</span>
                  <span v-else>Panduan Modul</span>
                </div>
              </div>
            </div>
          </div>
        </aside>

        <!-- RIGHT COLUMN: Detailed Reader Content Canvas (Exact Match to Design) -->
        <article class="flex-1 min-w-0 bg-white rounded-2xl border border-slate-200/90 p-6 sm:p-10 lg:p-12 shadow-xs">
          
          <!-- Category Breadcrumb / Subtitle -->
          <div class="text-xs font-bold text-sky-600 uppercase tracking-wider">
            {{ currentSystem.shortName }} • Modul {{ currentModule.code }} dari {{ String(modulesList.length).padStart(2, '0') }}
          </div>

          <!-- Module Main Title -->
          <h2 class="text-2xl sm:text-3xl font-extrabold text-slate-900 mt-2 tracking-tight">
            {{ currentModule.title }}
          </h2>

          <!-- Module Lead Text / Summary -->
          <p class="mt-4 text-slate-600 text-sm sm:text-base leading-relaxed">
            {{ currentModule.lead }}
          </p>

          <p class="mt-4 text-slate-600 text-sm sm:text-base leading-relaxed">
            {{ currentModule.intro }}
          </p>

          <hr class="my-8 border-slate-100" />

          <!-- Dynamic Sections List (Steps 1, 2, 3...) -->
          <div class="space-y-10">
            <div
              v-for="(sec, idx) in currentModule.sections"
              :key="idx"
              class="space-y-4"
            >
              <!-- Section Heading -->
              <h3 class="text-lg sm:text-xl font-bold text-slate-900 tracking-tight">
                {{ sec.heading }}
              </h3>

              <!-- Section Description -->
              <p class="text-sm sm:text-base text-slate-600 leading-relaxed">
                {{ sec.desc }}
              </p>

              <!-- Visual Diagram / Mockup Step Box (Matching Screenshot denah_meja_resto.png & input_pesanan_menu.png) -->
              <div
                v-if="sec.diagram"
                class="bg-slate-50 border border-slate-200 rounded-xl p-4 sm:p-5 my-4"
              >
                <!-- Image Header Filename Pill -->
                <div class="flex items-center gap-2 text-xs font-mono text-slate-500 mb-3">
                  <ImageIcon class="w-4 h-4 text-slate-400 shrink-0" />
                  <span>{{ sec.diagram.filename }}</span>
                </div>

                <!-- Highlight Box Containers (Simulating Step Highlights from Screenshot) -->
                <div class="space-y-2.5">
                  <div
                    v-for="(item, itemIdx) in sec.diagram.items"
                    :key="itemIdx"
                    class="p-3 bg-white rounded-lg border border-red-400/80 text-xs sm:text-sm font-semibold text-slate-800 shadow-2xs flex items-center justify-between"
                  >
                    <span>{{ item }}</span>
                    <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded bg-red-50 text-red-600 border border-red-200">
                      Langkah {{ itemIdx + 1 }}
                    </span>
                  </div>
                </div>
              </div>

              <!-- Structured Table Guide (e.g. Status Meja / Aksi) -->
              <div
                v-if="sec.table"
                class="overflow-x-auto rounded-xl border border-slate-200 mt-4 shadow-2xs"
              >
                <table class="w-full text-left text-xs sm:text-sm border-collapse">
                  <thead>
                    <tr class="bg-slate-50/80 border-b border-slate-200 text-slate-700 font-bold">
                      <th
                        v-for="(th, thIdx) in sec.table.headers"
                        :key="thIdx"
                        :class="[
                          'py-3 px-4',
                          thIdx === 0 ? 'w-16 text-center' : (thIdx === 1 ? 'w-48 sm:w-64' : '')
                        ]"
                      >
                        {{ th }}
                      </th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100 text-slate-700">
                    <tr
                      v-for="(row, rowIdx) in sec.table.rows"
                      :key="rowIdx"
                      class="hover:bg-slate-50/60 transition-colors"
                    >
                      <td class="py-3.5 px-4 font-semibold text-center text-slate-500">
                        {{ row[0] }}
                      </td>
                      <td class="py-3.5 px-4 font-semibold text-slate-900">
                        {{ row[1] }}
                      </td>
                      <td class="py-3.5 px-4 leading-relaxed text-slate-600">
                        {{ row[2] }}
                      </td>
                    </tr>
                  </tbody>
                </table>
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

