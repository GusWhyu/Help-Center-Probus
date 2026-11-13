import {
  Building2,
  MonitorCheck,
  Printer,
  Smartphone,
  Package,
  UtensilsCrossed,
  Sparkles,
  Zap,
  Laptop,
  Settings,
  Store,
  PieChart,
  ShoppingBag,
  Receipt,
  Wallet,
  BarChart3,
  TrendingUp,
  Sliders,
  Monitor,
  FileText,
  Calendar,
  Calculator,
  Network,
  BookOpen,
  HelpCircle,
  Folder,
  Layers,
  Wrench,
  Search,
  CheckCircle2,
  Globe,
  Users,
  ShieldCheck,
  CreditCard,
  Coffee,
  Hotel
} from 'lucide-vue-next'

// Standard 10 Module Order Map (Fallback for Standard ERP Modules)
const standardOrderModuleIcons = {
  1: Settings,         // 01. Setup Master
  2: UtensilsCrossed,  // 02. Transaksi / Order Taker / Sales / Production
  3: Store,            // 03. Outlet Activity
  4: PieChart,         // 04. Cost Control
  5: ShoppingBag,      // 05. Purchasing AP
  6: Receipt,          // 06. AR (Account Receivable)
  7: Wallet,           // 07. General Cashier
  8: BarChart3,        // 08. Accounting System
  9: TrendingUp,       // 09. Report
  10: Sliders,         // 10. Tools
}

// Known Official Systems Mapping
const officialSystemIcons = {
  'PHIS': Building2,
  'PRESTO': MonitorCheck,
  'PPOS': Printer,
  'SELF ORDER': Smartphone,
  'PR-PO ONLINE': Package,
  'ORDER TAKER': UtensilsCrossed,
  'VRESTROPICAL': Sparkles,
  'PPE': Zap,
  'BARONG POS': Laptop,
}

/**
 * AUTO-DETECT SYSTEM ICON
 * Mendeteksi ikon yang paling sesuai untuk Card Sistem berdasarkan kata kunci pada nama sistem.
 * 
 * @param {string} systemName - Nama sistem / kategori dari database (e.g. 'PHIS', 'Probus Laundry', 'Villa PMS')
 * @returns Component
 */
export function detectSystemIcon(systemName) {
  if (!systemName) return MonitorCheck

  const key = systemName.trim().toUpperCase()
  if (officialSystemIcons[key]) {
    return officialSystemIcons[key]
  }

  const s = systemName.toLowerCase()

  // 1. Hotel & Akomodasi
  if (s.includes('hotel') || s.includes('pms') || s.includes('room') || s.includes('villa') || s.includes('resort') || s.includes('front desk') || s.includes('housekeeping')) {
    return Building2
  }

  // 2. Restoran, F&B & Dining
  if (s.includes('resto') || s.includes('food') || s.includes('cafe') || s.includes('coffee') || s.includes('kuliner') || s.includes('dining') || s.includes('bar') && !s.includes('barong')) {
    return UtensilsCrossed
  }

  // 3. Kasir, Thermal Printer & Retail
  if (s.includes('kasir') || s.includes('retail') || s.includes('pos') || s.includes('print') || s.includes('struk') || s.includes('minimarket') || s.includes('toko')) {
    return Printer
  }

  // 4. Mobile, Self Service & QR Code
  if (s.includes('self') || s.includes('qr') || s.includes('mobile') || s.includes('app') || s.includes('phone') || s.includes('tamu') || s.includes('order taker')) {
    return Smartphone
  }

  // 5. Gudang, Logistik, Purchasing & Stok
  if (s.includes('gudang') || s.includes('logistik') || s.includes('inventory') || s.includes('stok') || s.includes('purchasing') || s.includes('po ') || s.includes('supply')) {
    return Package
  }

  // 6. Entertainment, Club, Spa & Hiburan
  if (s.includes('club') || s.includes('lounge') || s.includes('spa') || s.includes('massage') || s.includes('wellness') || s.includes('tropical') || s.includes('karaoke') || s.includes('beach')) {
    return Sparkles
  }

  // 7. Maintenance, Facility & Asset
  if (s.includes('facility') || s.includes('asset') || s.includes('maintenance') || s.includes('pemeliharaan') || s.includes('teknisi') || s.includes('ppe') || s.includes('listrik')) {
    return Zap
  }

  // 8. Cloud & Web Portal
  if (s.includes('cloud') || s.includes('online') || s.includes('web') || s.includes('laptop') || s.includes('portal') || s.includes('barong')) {
    return Laptop
  }

  // 9. Finance & Akuntansi
  if (s.includes('akuntan') || s.includes('finance') || s.includes('keuangan') || s.includes('pajak')) {
    return BarChart3
  }

  // Default fallback untuk sistem baru
  return MonitorCheck
}

/**
 * AUTO-DETECT MODULE ICON
 * Mendeteksi ikon yang paling sesuai untuk Card Modul berdasarkan nama modul dan nomor urutannya.
 * 
 * @param {string} moduleName - Nama modul dari database (e.g. 'Setup Master', 'Sales', 'Kitchen Display')
 * @param {number} orderNum - Nomor urut modul (1 s/d 10)
 * @returns Component
 */
export function detectModuleIcon(moduleName, orderNum = 0) {
  const m = (moduleName || '').toLowerCase()

  // 1. Setup & Master Data
  if (m.includes('setup') || m.includes('master') || m.includes('setting') || m.includes('konfigurasi') || m.includes('parameter')) {
    return Settings
  }

  // 2. Transaksi & Penjualan / Order
  if (m.includes('transaksi') || m.includes('sales') || m.includes('penjualan') || m.includes('order taker') || m.includes('order') || m.includes('production') || m.includes('dapur') && !m.includes('kds')) {
    return UtensilsCrossed
  }

  // 3. Outlet Activity
  if (m.includes('outlet') || m.includes('activity') || m.includes('toko') || m.includes('gerai') || m.includes('aktivitas')) {
    return Store
  }

  // 4. Cost Control & HPP
  if (m.includes('cost') || m.includes('control') || m.includes('hpp') || m.includes('biaya') || m.includes('waste')) {
    return PieChart
  }

  // 5. Purchasing & Penerimaan Barang
  if (m.includes('purchasing') || m.includes('purcashing') || m.includes('pembelian') || m.includes('pengadaan') || m.includes('po ') || m.includes('receiving') || m.includes('ap')) {
    return ShoppingBag
  }

  // 6. Accounts Receivable / Piutang
  if (m.includes('ar') || m.includes('receivable') || m.includes('piutang') || m.includes('invoice') || m.includes('billing') || m.includes('tagihan')) {
    return Receipt
  }

  // 7. General Cashier & Kas
  if (m.includes('general cashier') || m.includes('cashier') || m.includes('kasir') || m.includes('kas') || m.includes('petty cash') || m.includes('uang')) {
    return Wallet
  }

  // 8. Accounting & Keuangan
  if (m.includes('accounting') || m.includes('akuntansi') || m.includes('jurnal') || m.includes('buku besar') || m.includes('gl') || m.includes('neraca')) {
    return BarChart3
  }

  // 9. Report & Analitik
  if (m.includes('report') || m.includes('laporan') || m.includes('analitik') || m.includes('statistik') || m.includes('rekap') || m.includes('grafik')) {
    return TrendingUp
  }

  // 10. Tools & Hardware
  if (m.includes('tools') || m.includes('alat') || m.includes('hardware') || m.includes('perangkat') || m.includes('printer')) {
    return Sliders
  }

  // Additional Specific Modules
  if (m.includes('kds') || m.includes('kitchen display') || m.includes('display') || m.includes('screen')) {
    return Monitor
  }
  if (m.includes('shift') || m.includes('serah terima') || m.includes('closing')) {
    return FileText
  }
  if (m.includes('calendar') || m.includes('kalender') || m.includes('agenda') || m.includes('jadwal') || m.includes('booking')) {
    return Calendar
  }
  if (m.includes('calculator') || m.includes('kalkulator') || m.includes('hitung') || m.includes('simulasi')) {
    return Calculator
  }
  if (m.includes('network') || m.includes('jaringan') || m.includes('sync') || m.includes('reconnecting')) {
    return Network
  }
  if (m.includes('user') || m.includes('karyawan') || m.includes('staf') || m.includes('hak akses')) {
    return Users
  }

  // Fallback berdasarkan urutan modul 1-10 jika urutan valid
  if (orderNum && standardOrderModuleIcons[orderNum]) {
    return standardOrderModuleIcons[orderNum]
  }

  // Fallback umum untuk modul baru lainnya
  return BookOpen
}

export default {
  detectSystemIcon,
  detectModuleIcon,
}
