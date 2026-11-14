<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  CheckCircle2,
  ThumbsUp,
  Phone,
  AlertCircle,
  Loader2,
  ChevronRight,
  Printer,
  Globe,
  Lock,
  Calendar,
  Menu,
  X
} from 'lucide-vue-next'
import apiClient from '../api/client'
import { slugify, findArticleByParam } from '../utils/slug'

const route = useRoute()
const router = useRouter()
const mobileMenuOpen = ref(false)

const articleParam = computed(() => route.params.id)
const article = ref(null)
const relatedArticles = ref([])
const isLoading = ref(true)
const errorMessage = ref('')
const feedbackGiven = ref(false)

const systemNameMap = {
  1: 'PHIS',
  2: 'PRESTO',
  3: 'PPOS',
  4: 'SELF ORDER',
  5: 'PR-PO ONLINE',
  6: 'ORDER TAKER',
  7: 'VRESTROPICAL',
  8: 'PPE',
  9: 'BARONG POS',
}

const systemFullTitleMap = {
  1: 'PHIS Hotel PMS',
  2: 'PRESTO Resto POS',
  3: 'PPOS Point of Sales',
  4: 'Probus Self Order',
  5: 'PR-PO Online Logistik',
  6: 'Order Taker Mobile',
  7: 'Vrestropical POS',
  8: 'PPE 2.0 Management',
  9: 'Barong POS',
}

function formatArticleContent(html) {
  if (!html) return ''
  return html.replace(/https?:\/\/(127\.0\.0\.1|localhost)(:\d+)?\/uploads\//g, '/uploads/')
}

function formatDate(dateStr) {
  if (!dateStr) return '09 Nov 2026'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '09 Nov 2026'
  return d.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

function formatDateTime(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function loadRelatedArticles(currentId, katId) {
  try {
    const listRes = await apiClient.get('/articles', { params: { limit: 100 } })
    if (listRes.data && listRes.data.data && Array.isArray(listRes.data.data)) {
      const publicArticles = listRes.data.data.filter(item => 
        item.status === 'published' && 
        item.id !== currentId
      )
      const sameCategory = publicArticles.filter(item => {
        const pMenu = item.menu_systems && item.menu_systems.length > 0 ? item.menu_systems[0] : null
        return pMenu && pMenu.kategori_id === katId
      })
      const candidates = sameCategory.length > 0 ? sameCategory : publicArticles
      relatedArticles.value = candidates.slice(0, 2).map(item => ({
        id: item.id,
        title: item.title,
        slug: slugify(item.title),
        visibility: item.visibility || 'publik',
      }))
    }
  } catch (err) {
    console.warn('Failed to load related articles:', err)
  }
}

async function loadArticle() {
  isLoading.value = true
  errorMessage.value = ''
  const param = articleParam.value
  const numId = Number(param)

  try {
    // 1. If param is numeric, try fetching directly by ID
    if (!isNaN(numId) && numId > 0) {
      try {
        const res = await apiClient.get(`/articles/${numId}`)
        if (res.data && res.data.data) {
          processArticleData(res.data.data)
          const canonicalSlug = slugify(res.data.data.title)
          if (canonicalSlug && route.params.id !== canonicalSlug) {
            router.replace(`/solutions/${canonicalSlug}`)
          }
          return
        }
      } catch (err) {
        console.warn('Direct numeric article fetch failed, falling back to slug lookup:', err)
      }
    }

    // 2. Fetch from articles list and resolve by slug or ID
    const listRes = await apiClient.get('/articles', { params: { limit: 100 } })
    if (listRes.data && listRes.data.data && Array.isArray(listRes.data.data)) {
      const publicArticles = listRes.data.data.filter(item => item.status === 'published')
      const matched = findArticleByParam(publicArticles, param) || findArticleByParam(listRes.data.data, param)
      if (matched) {
        try {
          const detailRes = await apiClient.get(`/articles/${matched.id}`)
          if (detailRes.data && detailRes.data.data) {
            processArticleData(detailRes.data.data)
          } else {
            processArticleData(matched)
          }
        } catch {
          processArticleData(matched)
        }

        const canonicalSlug = slugify(matched.title)
        if (canonicalSlug && route.params.id !== canonicalSlug) {
          router.replace(`/solutions/${canonicalSlug}`)
        }
        return
      }
    }
    errorMessage.value = 'Artikel solusi tidak ditemukan atau belum dipublikasikan.'
  } catch (err) {
    console.error('Failed to load article:', err)
    errorMessage.value = 'Gagal memuat artikel dari server. Pastikan koneksi internet atau server aktif.'
  } finally {
    isLoading.value = false
  }
}

function processArticleData(raw) {
  const primaryMenu = raw.menu_systems && raw.menu_systems.length > 0 ? raw.menu_systems[0] : null
  const katId = primaryMenu ? primaryMenu.kategori_id : 1
  const sysShort = systemNameMap[katId] || 'PHIS'
  const sysFull = systemFullTitleMap[katId] || sysShort
  const menuName = primaryMenu ? primaryMenu.name : 'Setup Master'

  const words = (raw.content || '').split(/\s+/).length
  const minutes = Math.max(2, Math.min(10, Math.ceil(words / 150)))

  article.value = {
    id: raw.id,
    title: raw.title,
    systemShort: sysShort,
    systemTitle: sysFull,
    moduleName: menuName,
    visibility: raw.visibility || 'publik',
    status: raw.status || 'published',
    authorName: (raw.user && raw.user.name) || 'Superadministrator',
    formattedDate: formatDate(raw.created_at),
    formattedUpdated: formatDateTime(raw.updated_at || raw.created_at),
    readTime: `${minutes} menit baca`,
    content: formatArticleContent(raw.content),
    summary: raw.summary || '',
    updatedAt: raw.updated_at || raw.created_at,
    raw
  }
  isLoading.value = false
  loadRelatedArticles(raw.id, katId)
}

function printArticle() {
  window.print()
}

onMounted(() => {
  window.scrollTo({ top: 0, behavior: 'instant' })
  loadArticle()
})

watch(
  () => route.params.id,
  (newId, oldId) => {
    if (newId && newId !== oldId) {
      loadArticle()
    }
  }
)
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 text-slate-800 font-sans flex flex-col antialiased">
    <!-- 1. Utility Top Bar -->
    <header class="border-b border-slate-200/70 bg-white text-xs text-slate-500 py-2.5 px-4 sm:px-6 lg:px-12 xl:px-16">
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

    <!-- 2. Main Navigation Bar with "Kembali ke Bank Solusi" Button -->
    <nav class="border-b border-slate-200/90 bg-white py-3.5 px-4 sm:px-6 lg:px-12 xl:px-16 sticky top-0 z-40 shadow-xs">
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

          <!-- Pill Button: "← Kembali ke Bank Solusi" -->
          <button
            @click="router.push('/solutions')"
            class="inline-flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-lg bg-sky-50 hover:bg-sky-100 text-sky-700 font-semibold text-xs border border-sky-100 transition-colors shadow-2xs cursor-pointer"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Kembali ke Bank Solusi</span>
            <span class="sm:hidden">Bank Solusi</span>
          </button>
        </div>

        <!-- Desktop Navigation Links -->
        <div class="hidden md:flex items-center gap-8 text-sm font-medium">
          <router-link to="/" class="text-slate-600 hover:text-slate-950 transition-colors">
            Panduan Sistem
          </router-link>
          <router-link to="/solutions" class="text-sky-600 font-semibold hover:text-sky-700 transition-colors relative py-1">
            Bank Solusi
            <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-sky-600 rounded-full"></span>
          </router-link>
          <router-link
            to="/contact"
            class="text-slate-600 hover:text-slate-950 transition-colors"
          >
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
          class="block px-3 py-2 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-slate-900"
        >
          Panduan Sistem
        </router-link>
        <router-link
          to="/solutions"
          @click="mobileMenuOpen = false"
          class="block px-3 py-2 rounded-lg text-sm font-semibold bg-sky-50 text-sky-700"
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

    <!-- 3. Article Content Section (Identik dengan Desain Bank Solusi Probus Implementator) -->
    <main class="flex-1 w-full px-4 sm:px-6 lg:px-8 py-6 sm:py-8 bg-slate-50/50">
      <!-- Loading State -->
      <div v-if="isLoading" class="w-full py-20 text-center">
        <Loader2 class="w-10 h-10 animate-spin text-sky-600 mx-auto mb-4" />
        <h3 class="text-base font-bold text-slate-800">Memuat artikel solusi...</h3>
        <p class="text-xs text-slate-500 mt-1">Mengambil data panduan pemecahan masalah dari basis data.</p>
      </div>

      <!-- Error State -->
      <div v-else-if="errorMessage || !article" class="w-full max-w-xl mx-auto py-20 text-center">
        <div class="w-14 h-14 rounded-full bg-rose-50 text-rose-500 flex items-center justify-center mx-auto mb-4 border border-rose-100">
          <AlertCircle class="w-7 h-7" />
        </div>
        <h3 class="text-lg font-bold text-slate-900 mb-2">Artikel Tidak Ditemukan</h3>
        <p class="text-sm text-slate-500 mb-6">{{ errorMessage || 'Artikel yang Anda cari tidak tersedia atau belum dipublikasikan.' }}</p>
        <button
          @click="router.push('/solutions')"
          class="px-5 py-2.5 rounded-xl bg-slate-900 text-white font-semibold text-xs hover:bg-slate-800 transition-colors cursor-pointer inline-flex items-center gap-2"
        >
          <ArrowLeft class="w-4 h-4" />
          <span>Kembali ke Bank Solusi</span>
        </button>
      </div>

      <!-- Article View (Identik Desain Probus Implementator) -->
      <div v-else class="w-full max-w-5xl mx-auto">
        <!-- TOP BAR: BREADCRUMB & ACTION BUTTONS -->
        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 mb-6 pb-4 border-b border-slate-200">
          <!-- Breadcrumb -->
          <div class="flex items-center gap-2 text-xs sm:text-sm font-semibold text-slate-500 overflow-hidden">
            <span class="hover:text-slate-700 transition-colors shrink-0">Knowledge Base</span>
            <ChevronRight class="w-3.5 h-3.5 text-slate-400 shrink-0" />
            <router-link to="/solutions" class="hover:text-sky-600 transition-colors shrink-0">Bank Solusi</router-link>
            <ChevronRight class="w-3.5 h-3.5 text-slate-400 shrink-0" />
            <span class="text-slate-800 font-bold truncate max-w-xs sm:max-w-md lg:max-w-lg">
              {{ article.title }}
            </span>
          </div>

          <!-- Action Buttons -->
          <div class="flex items-center gap-2.5 shrink-0">
            <router-link
              to="/solutions"
              class="inline-flex items-center gap-2 px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 border border-slate-200 rounded-xl text-xs sm:text-sm font-semibold transition-all shadow-2xs"
            >
              <ArrowLeft class="w-3.5 h-3.5" />
              <span>Kembali ke Bank Solusi</span>
            </router-link>

            <button
              @click="printArticle"
              class="inline-flex items-center gap-2 px-4 py-2 bg-sky-500 hover:bg-sky-600 text-white rounded-xl text-xs sm:text-sm font-bold shadow-md shadow-sky-500/20 transition-all hover:scale-[1.02] cursor-pointer"
            >
              <Printer class="w-3.5 h-3.5" />
              <span>Cetak Solusi</span>
            </button>
          </div>
        </div>

        <!-- MAIN ARTICLE CARD CONTAINER (Identik Bank Solusi Probus Implementator) -->
        <article class="bg-white border border-slate-200 rounded-3xl p-6 sm:p-10 shadow-sm space-y-8">
          <!-- Header Artikel -->
          <div class="space-y-4 pb-6 border-b border-slate-100">
            <!-- Kategori Sistem & Modul + Badge Visibilitas & Status -->
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div class="flex items-center gap-2 text-xs sm:text-sm font-extrabold uppercase tracking-wider text-sky-600">
                <span>{{ article.systemShort }}</span>
                <span class="text-slate-300">•</span>
                <span>{{ article.moduleName }}</span>
              </div>

              <div class="flex items-center gap-2">
                <!-- Badge Visibilitas -->
                <span
                  v-if="article.visibility === 'publik'"
                  class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-black uppercase tracking-wider bg-emerald-50 text-emerald-700 border border-emerald-200 shadow-2xs"
                >
                  <Globe class="w-3 h-3 text-emerald-500 shrink-0" />
                  <span>BLOG KLIEN</span>
                </span>
                <span
                  v-else
                  class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-black uppercase tracking-wider bg-amber-50 text-amber-700 border border-amber-200 shadow-2xs"
                >
                  <Lock class="w-3 h-3 text-amber-500 shrink-0" />
                  <span>INTERNAL</span>
                </span>

                <!-- Status Badge -->
                <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-md text-[10px] font-extrabold uppercase tracking-wider bg-blue-50 text-blue-700 border border-blue-200">
                  <CheckCircle2 class="w-3 h-3 text-blue-600" />
                  <span>Diterbitkan</span>
                </span>
              </div>
            </div>

            <!-- Judul Utama -->
            <h1 class="text-2xl sm:text-3xl md:text-4xl font-black text-slate-900 tracking-tight leading-tight">
              {{ article.title }}
            </h1>

            <!-- Metadata: Tanggal Saja -->
            <div class="flex items-center gap-1.5 text-xs font-semibold text-slate-500 pt-1">
              <Calendar class="w-3.5 h-3.5 text-slate-400" />
              <span>{{ article.formattedDate }}</span>
            </div>

            <!-- Ringkasan Pengantar (Callout Biru Garis Kiri Identik Probus Implementator) -->
            <div
              v-if="article.summary"
              class="pt-2 text-sm sm:text-base text-slate-600 font-medium leading-relaxed italic border-l-4 border-blue-500 pl-4 py-2 bg-blue-50/30 rounded-r-xl"
            >
              {{ article.summary }}
            </div>
          </div>

          <!-- ISI DOKUMEN / KONTEN ARTIKEL (TINYMCE RENDER) -->
          <div
            class="article-content max-w-none text-slate-800 text-sm sm:text-base leading-relaxed"
            v-html="article.content"
          ></div>

          <!-- Card Footer: Terakhir Diperbarui & Navigasi -->
          <div class="pt-6 border-t border-slate-100 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-slate-500">
            <div>
              Terakhir diperbarui: <span class="font-bold text-slate-700">{{ article.formattedUpdated }}</span>
            </div>
            <router-link
              to="/solutions"
              class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-xs font-bold transition-all"
            >
              Semua Solusi
            </router-link>
          </div>
        </article>

        <!-- Interactive Feedback Card -->
        <div class="mt-8 p-6 sm:p-7 rounded-3xl bg-white border border-slate-200 shadow-sm flex flex-col sm:flex-row items-center justify-between gap-4">
          <div>
            <h4 class="text-sm font-bold text-slate-900">Apakah artikel pemecahan masalah ini membantu Anda?</h4>
            <p class="text-xs text-slate-500 mt-0.5">Masukan Anda membantu kami meningkatkan kualitas dokumentasi sistem Probus.</p>
          </div>

          <div class="flex items-center gap-2 shrink-0">
            <button
              v-if="!feedbackGiven"
              @click="feedbackGiven = true"
              class="px-4 py-2 rounded-xl border border-slate-300 bg-white hover:bg-slate-100 text-slate-800 font-semibold text-xs inline-flex items-center gap-1.5 transition-colors shadow-2xs cursor-pointer"
            >
              <ThumbsUp class="w-3.5 h-3.5 text-sky-600" />
              <span>Ya, Berhasil</span>
            </button>
            <div
              v-else
              class="px-4 py-2 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-700 font-bold text-xs inline-flex items-center gap-1.5"
            >
              <CheckCircle2 class="w-4 h-4 text-emerald-600" />
              <span>Terima kasih atas masukannya!</span>
            </div>
          </div>
        </div>

        <!-- Related Articles Section (Solusi Terkait Lainnya di Modul ...) -->
        <div v-if="relatedArticles.length > 0" class="mt-8 space-y-4">
          <h3 class="text-base sm:text-lg font-black text-slate-800 tracking-tight">
            Solusi Terkait Lainnya di Modul <span class="text-sky-600">{{ article.moduleName }}</span>
          </h3>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <router-link
              v-for="rel in relatedArticles"
              :key="rel.id"
              :to="`/solutions/${rel.slug}`"
              class="p-5 bg-white border border-slate-200 rounded-2xl hover:border-sky-500 hover:shadow-md transition-all duration-200 group flex flex-col justify-between shadow-xs"
            >
              <div>
                <div class="text-[10px] font-black uppercase tracking-wider text-sky-600 mb-1.5">
                  {{ rel.visibility === 'publik' ? 'BLOG KLIEN' : 'INTERNAL' }}
                </div>
                <h4 class="text-sm font-bold text-slate-900 group-hover:text-sky-600 transition-colors line-clamp-2">
                  {{ rel.title }}
                </h4>
              </div>
              <div class="mt-4 pt-3 flex items-center justify-between border-t border-slate-100 text-xs font-bold text-sky-600">
                <span>Baca panduan</span>
                <ChevronRight class="w-3.5 h-3.5 group-hover:translate-x-1 transition-transform" />
              </div>
            </router-link>
          </div>
        </div>
      </div>
    </main>

    <!-- 4. Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 px-4 sm:px-6 lg:px-12 xl:px-16 text-xs text-slate-500">
      <div class="w-full flex flex-col sm:flex-row items-center justify-between gap-3 text-center sm:text-left">
        <p>© 2026 CV Mitra Solusi Mandiri • Bank Solusi Resmi</p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* Styling Khusus Format Konten TinyMCE (Identik Sama Persis dengan Probus Implementator) */
:deep(.article-content h1), 
:deep(.article-content h2), 
:deep(.article-content h3), 
:deep(.article-content h4) {
  color: inherit;
  font-weight: 800;
  letter-spacing: -0.025em;
  margin-top: 1.75em;
  margin-bottom: 0.75em;
  line-height: 1.3;
}
:deep(.article-content h1) { font-size: 1.75rem; }
:deep(.article-content h2) { font-size: 1.35rem; }
:deep(.article-content h3) { font-size: 1.15rem; }

:deep(.article-content p) {
  margin-bottom: 1.25em;
  line-height: 1.75;
  color: #334155;
}

:deep(.article-content ul), 
:deep(.article-content ol) {
  margin-top: 0.75em;
  margin-bottom: 1.25em;
  padding-left: 1.5em;
}
:deep(.article-content ul) { list-style-type: disc; }
:deep(.article-content ol) { list-style-type: decimal; }
:deep(.article-content li) {
  margin-bottom: 0.5em;
  line-height: 1.6;
}

/* Tabel Responsif & Indikator (Border Penutup Keliling & Sel Penuh) */
:deep(.article-content table) {
  width: 100% !important;
  border-collapse: collapse !important;
  margin: 1.5em 0 !important;
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
  font-weight: 700 !important;
  color: #0f172a !important;
}

/* Preserve dark navy header from TinyMCE Word importer (bg #0b1a36) */
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

/* Gambar Konten (Identik Probus Implementator) */
:deep(.article-content img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.75rem;
  margin: 1.5em auto;
  display: block;
  border: 1px solid rgba(148, 163, 184, 0.3);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

/* Callout Box Wireframe / Screenshot (Identik Probus Implementator) */
:deep(.article-content .mockup-box),
:deep(.article-content div[style*="border: 1px solid #ef4444"]),
:deep(.article-content div[style*="border: 1px solid red"]) {
  border-radius: 0.5rem;
  padding: 0.75rem 1rem;
  margin: 1.25em 0;
  background-color: rgba(239, 68, 68, 0.05);
}
</style>

