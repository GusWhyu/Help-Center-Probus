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
  Globe
} from 'lucide-vue-next'

// Registry of Lucide Icon Components
const iconRegistry = {
  'building-2': Building2,
  'building': Building2,
  'monitor-check': MonitorCheck,
  'printer': Printer,
  'smartphone': Smartphone,
  'package': Package,
  'utensils': UtensilsCrossed,
  'utensils-crossed': UtensilsCrossed,
  'sparkles': Sparkles,
  'zap': Zap,
  'laptop': Laptop,
  'settings': Settings,
  'store': Store,
  'pie-chart': PieChart,
  'shopping-bag': ShoppingBag,
  'receipt': Receipt,
  'wallet': Wallet,
  'bar-chart-3': BarChart3,
  'bar-chart': BarChart3,
  'trending-up': TrendingUp,
  'sliders': Sliders,
  'monitor': Monitor,
  'file-text': FileText,
  'calendar': Calendar,
  'calculator': Calculator,
  'network': Network,
  'book-open': BookOpen,
  'help-circle': HelpCircle,
  'folder': Folder,
  'layers': Layers,
  'wrench': Wrench,
  'search': Search,
  'check-circle': CheckCircle2,
  'globe': Globe
}

/**
 * Resolves a string icon name to a Lucide Vue Component.
 * Supports fallback detection based on item title/name if icon is missing.
 * 
 * @param {string} iconName - Stored icon identifier in database (e.g. 'printer', 'building-2')
 * @param {string} fallbackTitle - Title/Name of system or module for semantic fallback
 * @param {string} defaultType - 'system' or 'module'
 * @returns Component
 */
export function resolveIcon(iconName, fallbackTitle = '', defaultType = 'module') {
  if (iconName) {
    const normalized = iconName.trim().toLowerCase()
    if (iconRegistry[normalized]) {
      return iconRegistry[normalized]
    }
  }

  // Fallback 1: Semantic keyword detection based on title/name
  if (fallbackTitle) {
    const t = fallbackTitle.toLowerCase()
    if (t.includes('hotel') || t.includes('phis') || t.includes('room')) return Building2
    if (t.includes('resto') || t.includes('food') || t.includes('transaksi') || t.includes('order taker') || t.includes('kitchen')) return UtensilsCrossed
    if (t.includes('kasir') || t.includes('pos') || t.includes('print')) return Printer
    if (t.includes('qr') || t.includes('self') || t.includes('mobile')) return Smartphone
    if (t.includes('gudang') || t.includes('logistik') || t.includes('po ') || t.includes('stok') || t.includes('package')) return Package
    if (t.includes('setup') || t.includes('master') || t.includes('setting')) return Settings
    if (t.includes('outlet') || t.includes('activity') || t.includes('store')) return Store
    if (t.includes('cost') || t.includes('control')) return PieChart
    if (t.includes('purchasing') || t.includes('purcashing') || t.includes('ap')) return ShoppingBag
    if (t.includes('ar ') || t.includes('receivable') || t.includes('piutang') || t.includes('invoice')) return Receipt
    if (t.includes('general cashier') || t.includes('cashier') || t.includes('kas')) return Wallet
    if (t.includes('akuntan') || t.includes('accounting') || t.includes('finance') || t.includes('jurnal')) return BarChart3
    if (t.includes('lapor') || t.includes('report') || t.includes('grafik')) return TrendingUp
    if (t.includes('tool') || t.includes('hardware') || t.includes('alat')) return Sliders
    if (t.includes('calendar') || t.includes('agenda')) return Calendar
    if (t.includes('calculator') || t.includes('hitung')) return Calculator
    if (t.includes('kds') || t.includes('display')) return Monitor
    if (t.includes('shift')) return FileText
    if (t.includes('network') || t.includes('jaringan') || t.includes('reconnecting')) return Network
    if (t.includes('cloud') || t.includes('laptop') || t.includes('barong')) return Laptop
    if (t.includes('club') || t.includes('tropical') || t.includes('vrestropical')) return Sparkles
    if (t.includes('facility') || t.includes('asset') || t.includes('ppe')) return Zap
  }

  // Fallback 2: General default by type
  return defaultType === 'system' ? MonitorCheck : BookOpen
}

export default resolveIcon
