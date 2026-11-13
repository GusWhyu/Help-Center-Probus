<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  Clock,
  CheckCircle2,
  ThumbsUp,
  Phone,
  AlertCircle,
  Loader2,
  ChevronRight,
  Share2,
  Printer
} from 'lucide-vue-next'
import apiClient from '../api/client'

const route = useRoute()
const router = useRouter()

const articleId = computed(() => Number(route.params.id))
const article = ref(null)
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
  return html.replace(/http:\/\/(127\.0\.0\.1|localhost):8000\/uploads\//g, '/uploads/')
}

async function loadArticle() {
  isLoading.value = true
  errorMessage.value = ''
  try {
    // 1. Try fetching specific article by ID
    const res = await apiClient.get(`/articles/${articleId.value}`)
    if (res.data && res.data.data) {
      processArticleData(res.data.data)
      return
    }
  } catch (err) {
    console.warn('Direct article fetch failed, falling back to articles list:', err)
  }

  // 2. Fallback: fetch from articles list
  try {
    const listRes = await apiClient.get('/articles', { params: { limit: 100 } })
    if (listRes.data && listRes.data.data && Array.isArray(listRes.data.data)) {
      const found = listRes.data.data.find(item => Number(item.id) === articleId.value)
      if (found) {
        processArticleData(found)
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
  const katId = primaryMenu ? primaryMenu.kategori_id : 2
  const sysShort = systemNameMap[katId] || 'PRESTO'
  const sysFull = systemFullTitleMap[katId] || sysShort
  const menuName = primaryMenu ? primaryMenu.name : 'Hardware & Printer'

  const words = (raw.content || '').split(/\s+/).length
  const minutes = Math.max(2, Math.min(10, Math.ceil(words / 150)))

  article.value = {
    id: raw.id,
    title: raw.title,
    systemShort: sysShort,
    systemTitle: sysFull,
    moduleName: menuName,
    readTime: `${minutes} menit baca`,
    content: formatArticleContent(raw.content),
    summary: raw.summary || '',
    updatedAt: raw.updated_at || raw.created_at,
    raw
  }
  isLoading.value = false
}

function printArticle() {
  window.print()
}

onMounted(() => {
  window.scrollTo({ top: 0, behavior: 'instant' })
  loadArticle()
})
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 text-slate-800 font-sans flex flex-col antialiased">
    <!-- 1. Utility Top Bar -->
    <header class="border-b border-slate-200/70 bg-white text-xs text-slate-500 py-2.5 px-6 sm:px-12 lg:px-20 xl:px-28">
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

    <!-- 2. Main Navigation Bar with "Kembali ke Bank Solusi" Button -->
    <nav class="border-b border-slate-200/80 bg-white/95 backdrop-blur-sm py-3.5 px-6 sm:px-12 lg:px-20 xl:px-28 sticky top-0 z-40 shadow-xs">
      <div class="max-w-[1600px] mx-auto flex items-center justify-between">
        <!-- Logo & Back Button -->
        <div class="flex items-center gap-3.5">
          <div class="flex items-center gap-2.5 cursor-pointer" @click="router.push('/')">
            <img
              src="/probus-logo.png"
              alt="Probus System"
              class="h-8 sm:h-9 w-auto object-contain hover:opacity-90 transition-opacity"
            />
          </div>
          <span class="text-slate-300 font-light">|</span>
          <span class="text-slate-500 text-sm font-medium hidden sm:inline">Pusat Bantuan Publik</span>

          <!-- Pill Button: "← Kembali ke Bank Solusi" (matches screenshot) -->
          <button
            @click="router.push('/solutions')"
            class="ml-2 inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-sky-50 hover:bg-sky-100 text-sky-700 font-semibold text-xs border border-sky-100 transition-colors shadow-2xs cursor-pointer"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>Kembali ke Bank Solusi</span>
          </button>
        </div>

        <!-- Navigation Links -->
        <div class="hidden lg:flex items-center gap-8 text-sm font-medium">
          <router-link to="/" class="text-slate-600 hover:text-slate-950 transition-colors">
            Panduan Sistem
          </router-link>
          <router-link to="/solutions" class="text-sky-600 font-semibold hover:text-sky-700 transition-colors relative py-1">
            Bank Solusi
            <span class="absolute bottom-0 left-0 right-0 h-0.5 bg-sky-600 rounded-full"></span>
          </router-link>
          <a
            href="https://wa.me/6281910875514"
            target="_blank"
            class="text-slate-600 hover:text-slate-950 transition-colors"
          >
            Hubungi Implementator
          </a>
        </div>
      </div>
    </nav>

    <!-- 3. Article Content Section -->
    <main class="flex-1 bg-white">
      <!-- Loading State -->
      <div v-if="isLoading" class="max-w-4xl mx-auto px-6 py-20 text-center">
        <Loader2 class="w-10 h-10 animate-spin text-sky-600 mx-auto mb-4" />
        <h3 class="text-base font-bold text-slate-800">Memuat artikel solusi...</h3>
        <p class="text-xs text-slate-500 mt-1">Mengambil data panduan pemecahan masalah dari basis data.</p>
      </div>

      <!-- Error State -->
      <div v-else-if="errorMessage || !article" class="max-w-xl mx-auto px-6 py-20 text-center">
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

      <!-- Article View (Matching Screenshot media_1794537273140_24a23ca4.png) -->
      <article v-else class="max-w-4xl mx-auto px-6 sm:px-8 py-10 sm:py-14">
        <!-- Breadcrumb / Category Metadata -->
        <div class="flex items-center gap-2 text-xs sm:text-sm font-semibold tracking-wide text-sky-600 mb-2.5">
          <span>{{ article.systemTitle }}</span>
          <span class="text-slate-300">•</span>
          <span>{{ article.moduleName }}</span>
        </div>

        <!-- Article Title -->
        <h1 class="text-2xl sm:text-3xl lg:text-4xl font-extrabold text-slate-900 leading-tight tracking-tight mb-4">
          {{ article.title }}
        </h1>

        <!-- Reading Metadata & Status Bar -->
        <div class="flex flex-wrap items-center justify-between gap-4 text-xs text-slate-500 pb-5 border-b border-slate-200/80 mb-8">
          <div class="flex items-center gap-4">
            <span class="flex items-center gap-1.5">
              <Clock class="w-3.5 h-3.5 text-slate-400" />
              <span>{{ article.readTime }}</span>
            </span>
            <span class="flex items-center gap-1.5 text-emerald-600 font-semibold">
              <CheckCircle2 class="w-3.5 h-3.5" />
              <span>Terverifikasi Tim Teknis Probus</span>
            </span>
          </div>

          <button
            @click="printArticle"
            class="text-slate-400 hover:text-slate-700 transition-colors inline-flex items-center gap-1 cursor-pointer"
            title="Cetak panduan ini"
          >
            <Printer class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Cetak Solusi</span>
          </button>
        </div>

        <!-- Article Rich HTML Content -->
        <div
          class="solution-article-content text-slate-700 text-sm sm:text-base leading-relaxed"
          v-html="article.content"
        ></div>

        <!-- Interactive Feedback Card -->
        <div class="mt-14 p-6 rounded-2xl bg-slate-50 border border-slate-200/90 flex flex-col sm:flex-row items-center justify-between gap-4">
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

        <!-- Back Link Navigation -->
        <div class="mt-10 pt-6 border-t border-slate-100 flex items-center justify-between text-xs text-slate-500">
          <button
            @click="router.push('/solutions')"
            class="inline-flex items-center gap-1.5 font-semibold text-slate-700 hover:text-sky-600 transition-colors cursor-pointer"
          >
            <ArrowLeft class="w-4 h-4" />
            <span>Kembali ke Daftar Bank Solusi</span>
          </button>

          <a
            href="https://wa.me/6281910875514"
            target="_blank"
            class="font-semibold text-emerald-600 hover:text-emerald-700 transition-colors inline-flex items-center gap-1"
          >
            <span>Masih ada kendala? Chat WhatsApp Tim Support</span>
            <ChevronRight class="w-3.5 h-3.5" />
          </a>
        </div>
      </article>
    </main>

    <!-- 4. Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 px-6 sm:px-12 lg:px-20 xl:px-28 text-xs text-slate-500">
      <div class="max-w-[1600px] mx-auto flex flex-col sm:flex-row items-center justify-between gap-3">
        <p>© 2026 CV Mitra Solusi Mandiri • Bank Solusi Resmi</p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
:deep(.solution-article-content) {
  color: #334155;
  font-size: 0.95rem;
  line-height: 1.75;
}

:deep(.solution-article-content p) {
  margin-bottom: 1.25rem;
  color: #334155;
}

:deep(.solution-article-content h1),
:deep(.solution-article-content h2),
:deep(.solution-article-content h3),
:deep(.solution-article-content h4) {
  color: #0f172a;
  font-weight: 700;
  margin-top: 2rem;
  margin-bottom: 0.75rem;
  line-height: 1.4;
}

:deep(.solution-article-content h2) {
  font-size: 1.15rem;
}

:deep(.solution-article-content h3) {
  font-size: 1.05rem;
}

:deep(.solution-article-content ul),
:deep(.solution-article-content ol) {
  padding-left: 1.5rem;
  margin-bottom: 1.25rem;
}

:deep(.solution-article-content ol) {
  list-style-type: decimal;
}

:deep(.solution-article-content ul) {
  list-style-type: disc;
}

:deep(.solution-article-content li) {
  margin-bottom: 0.5rem;
  line-height: 1.65;
}

:deep(.solution-article-content table) {
  width: 100% !important;
  border-collapse: collapse !important;
  margin: 1.5rem 0 !important;
  border: 1px solid #cbd5e1 !important;
  border-radius: 8px !important;
  overflow: hidden !important;
}

:deep(.solution-article-content th) {
  background-color: #f8fafc !important;
  color: #0f172a !important;
  font-weight: 700 !important;
  font-size: 0.875rem !important;
  text-align: left !important;
  padding: 10px 14px !important;
  border: 1px solid #cbd5e1 !important;
}

:deep(.solution-article-content td) {
  padding: 10px 14px !important;
  border: 1px solid #e2e8f0 !important;
  font-size: 0.875rem !important;
  color: #334155 !important;
  line-height: 1.5 !important;
}

:deep(.solution-article-content tr:nth-child(even) td) {
  background-color: #fafbfc !important;
}

:deep(.solution-article-content img) {
  max-width: 100% !important;
  height: auto !important;
  border-radius: 8px !important;
  margin: 1.25rem 0 !important;
  border: 1px solid #e2e8f0 !important;
}

:deep(.solution-article-content code) {
  background-color: #f1f5f9;
  color: #0f172a;
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 0.875rem;
  font-family: monospace;
}

:deep(.solution-article-content pre) {
  background-color: #0f172a;
  color: #f8fafc;
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 1.25rem 0;
  font-size: 0.875rem;
}
</style>

