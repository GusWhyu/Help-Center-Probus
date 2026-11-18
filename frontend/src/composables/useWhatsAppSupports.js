import { ref, onMounted } from 'vue'
import {
  Headphones,
  ShoppingCart,
  GraduationCap,
  Wrench,
  Laptop,
  Code
} from 'lucide-vue-next'
import apiClient from '../api/client'

// Fallback initial contacts jika backend offline
const defaultContacts = [
  {
    id: 1,
    name: "CS Probussystem",
    role: "Customer Service & Bantuan Teknis",
    numberDisplay: "+62 819-1087-5514",
    waNumber: "6281910875514",
    badge: "CS Support",
    badgeClass: "bg-emerald-50 text-emerald-700 border-emerald-200",
    icon: Headphones,
    defaultMessage: "Halo CS Probussystem, saya butuh bantuan terkait kendala sistem Probus.",
  },
  {
    id: 2,
    name: "E Commerce Probussystem",
    role: "Online Booking, Channel & Integrasi",
    numberDisplay: "+62 878-6144-6124",
    waNumber: "6287861446124",
    badge: "E-Commerce",
    badgeClass: "bg-sky-50 text-sky-700 border-sky-200",
    icon: ShoppingCart,
    defaultMessage: "Halo E-Commerce Probussystem, saya ingin konsultasi terkait layanan integrasi & E-Commerce Probus.",
  },
  {
    id: 3,
    name: "Trainer Probussystem",
    role: "Panduan Operasional & Training Modul",
    numberDisplay: "+62 819-1087-5514",
    waNumber: "6281910875514",
    badge: "Trainer",
    badgeClass: "bg-indigo-50 text-indigo-700 border-indigo-200",
    icon: GraduationCap,
    defaultMessage: "Halo Trainer Probussystem, saya ingin panduan dan jadwal training modul sistem Probus.",
  },
]

// Global reactive state shared between FloatingWhatsAppSupport and WhatsAppSupportDropdown
const contacts = ref([...defaultContacts])
const isLoading = ref(false)
const hasFetched = ref(false)

function resolveContactIcon(iconName, badge, name) {
  const iconStr = (iconName || '').toLowerCase()
  const badgeStr = (badge || '').toLowerCase()
  const nameStr = (name || '').toLowerCase()

  if (iconStr.includes('headset') || badgeStr.includes('cs') || nameStr.includes('cs')) {
    return Headphones
  }
  if (iconStr.includes('cart') || badgeStr.includes('commerce') || nameStr.includes('commerce') || badgeStr.includes('shop')) {
    return ShoppingCart
  }
  if (iconStr.includes('graduation') || badgeStr.includes('train') || nameStr.includes('train')) {
    return GraduationCap
  }
  if (iconStr.includes('wrench') || iconStr.includes('tool') || badgeStr.includes('dev') || nameStr.includes('program') || badgeStr.includes('tech')) {
    return Wrench
  }
  if (iconStr.includes('code') || iconStr.includes('laptop')) {
    return Code
  }
  return Headphones
}

function resolveBadgeClass(color) {
  const c = (color || 'emerald').toLowerCase().trim()
  const colorMap = {
    emerald: 'bg-emerald-50 text-emerald-700 border-emerald-200',
    green: 'bg-emerald-50 text-emerald-700 border-emerald-200',
    sky: 'bg-sky-50 text-sky-700 border-sky-200',
    blue: 'bg-blue-50 text-blue-700 border-blue-200',
    indigo: 'bg-indigo-50 text-indigo-700 border-indigo-200',
    purple: 'bg-purple-50 text-purple-700 border-purple-200',
    violet: 'bg-purple-50 text-purple-700 border-purple-200',
    amber: 'bg-amber-50 text-amber-700 border-amber-200',
    yellow: 'bg-amber-50 text-amber-700 border-amber-200',
    rose: 'bg-rose-50 text-rose-700 border-rose-200',
    red: 'bg-rose-50 text-rose-700 border-rose-200',
    teal: 'bg-teal-50 text-teal-700 border-teal-200',
    slate: 'bg-slate-50 text-slate-700 border-slate-200',
  }
  return colorMap[c] || colorMap['emerald']
}

function sanitizeWaNumber(phone) {
  if (!phone) return ''
  const digits = String(phone).replace(/\D/g, '')
  if (digits.startsWith('0')) {
    return '62' + digits.slice(1)
  }
  if (digits.startsWith('62')) {
    return digits
  }
  return digits
}

export async function fetchWhatsAppSupports() {
  if (hasFetched.value && contacts.value.length > 0) return
  isLoading.value = true
  try {
    const res = await apiClient.get('/whatsapp-supports')
    if (res.data && res.data.data && Array.isArray(res.data.data) && res.data.data.length > 0) {
      contacts.value = res.data.data.map(item => {
        const rawPhone = item.phone_number || ''
        const waNum = sanitizeWaNumber(rawPhone)
        const role = item.description || (item.badge ? `Bantuan Layanan ${item.badge}` : 'Bantuan Operasional Probus')
        const badge = item.badge || 'Support'
        return {
          id: item.id,
          name: item.name,
          role: role,
          numberDisplay: rawPhone,
          waNumber: waNum,
          badge: badge,
          badgeClass: resolveBadgeClass(item.badge_color),
          icon: resolveContactIcon(item.icon, badge, item.name),
          defaultMessage: `Halo ${item.name}, saya butuh bantuan terkait kendala sistem Probus.`,
        }
      })
      hasFetched.value = true
    }
  } catch (err) {
    console.warn('Gagal memuat kontak WhatsApp dari server, menggunakan data fallback:', err)
  } finally {
    isLoading.value = false
  }
}

export function useWhatsAppSupports() {
  onMounted(() => {
    fetchWhatsAppSupports()
  })

  function openWhatsApp(contact) {
    if (!contact || !contact.waNumber) return
    const url = `https://wa.me/${contact.waNumber}?text=${encodeURIComponent(contact.defaultMessage)}`
    window.open(url, '_blank')
  }

  return {
    contacts,
    isLoading,
    openWhatsApp,
    fetchWhatsAppSupports,
  }
}
