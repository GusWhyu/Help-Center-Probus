import apiClient from '../api/client'

/**
 * Utility functions for URL Slug generation and resolution
 * Maps systems, modules, and articles to human-readable & secure URL slugs without altering MySQL database schemas.
 */

// Mapping of System ID to canonical URL slug and vice-versa
export const systemSlugMap = {
  1: 'phis',
  2: 'presto',
  3: 'ppos',
  4: 'self-order',
  5: 'pr-po-online',
  6: 'order-taker',
  7: 'vrestropical',
  8: 'ppe',
  9: 'barong-pos',
  12: 'channel-manager-ezze',
}

// Aliases mapping for backward compatibility and friendly URLs
export const systemAliasToId = {
  // IDs as string
  '1': 1,
  '2': 2,
  '3': 3,
  '4': 4,
  '5': 5,
  '6': 6,
  '7': 7,
  '8': 8,
  '9': 9,
  // Canonical slugs
  'phis': 1,
  'presto': 2,
  'ppos': 3,
  'self-order': 4,
  'selforder': 4,
  'pr-po-online': 5,
  'pr-po': 5,
  'prpo': 5,
  'order-taker': 6,
  'ordertaker': 6,
  'vrestropical': 7,
  'ppe': 8,
  'barong-pos': 9,
  'barong': 9,
  // Full names
  'probus-hotel-pms': 1,
  'hotel-pms': 1,
  'presto-resto-pos': 2,
  'resto-pos': 2,
  'ppos-point-of-sales': 3,
  'point-of-sales': 3,
  'probus-self-order': 4,
  'order-taker-mobile': 6,
  'ppe-2-0-management': 8,
  'ppe-management': 8,
  'barong-pos-cloud': 9,
}

// In-memory registry for dynamic categories loaded from database
const dynamicCategoryMap = new Map()
const dynamicSlugToIdMap = new Map()

// Pre-seeded categories including new systems (e.g. 10: Project WEB, 11: CHANNEL MANAGER (eZee), 12: Channel Manager (Ezze))
const preseededCategories = [
  { id: 1, name: 'PHIS' },
  { id: 2, name: 'PRESTO' },
  { id: 3, name: 'PPOS' },
  { id: 4, name: 'SELF ORDER' },
  { id: 5, name: 'PR-PO ONLINE' },
  { id: 6, name: 'ORDER TAKER' },
  { id: 7, name: 'VRESTROPICAL' },
  { id: 8, name: 'PPE' },
  { id: 9, name: 'BARONG POS' },
  { id: 10, name: 'Project WEB' },
  { id: 11, name: 'CHANNEL MANAGER (eZee)' },
  { id: 12, name: 'Channel Manager (Ezze)' },
]

// In-flight categories fetch promise to avoid redundant API calls
let categoriesPromise = null

/**
 * Fetch all categories from backend API and register them dynamically.
 * Ensures unknown slugs can be resolved on initial refresh/F5 without falling back to PRESTO.
 */
export async function ensureCategoriesLoaded(forceRefresh = false) {
  if (categoriesPromise && !forceRefresh) {
    return categoriesPromise
  }
  try {
    categoriesPromise = apiClient.get('/categories')
      .then(res => {
        if (res.data && res.data.data && Array.isArray(res.data.data)) {
          registerCategories(res.data.data)
          return res.data.data
        }
        return []
      })
      .catch(err => {
        console.warn('[slug] Auto-fetch categories failed:', err)
        categoriesPromise = null
        return []
      })
    return await categoriesPromise
  } catch (err) {
    console.warn('[slug] Error in ensureCategoriesLoaded:', err)
    return []
  }
}

/**
 * Asynchronously resolve a system ID from route parameter, auto-fetching categories if unrecognised.
 */
export async function resolveSystemId(param, defaultId = 2) {
  if (param === undefined || param === null || param === '') {
    return defaultId
  }

  // 1. Try resolving with currently loaded memory
  const syncId = getSystemIdFromParam(param, null)
  if (syncId !== null) {
    return syncId
  }

  // 2. Fresh reload / F5 with new system slug: fetch latest categories from database
  await ensureCategoriesLoaded(true)

  // 3. Re-resolve with newly registered categories
  const reResolved = getSystemIdFromParam(param, null)
  if (reResolved !== null) {
    return reResolved
  }

  return defaultId
}

/**
 * Register categories dynamically into the in-memory lookup cache.
 * Automatically handles any new system added in the database without code changes.
 *
 * @param {Array<Object>} categories - List of category objects [{ id, name, ... }]
 */
export function registerCategories(categories) {
  if (!Array.isArray(categories)) return
  categories.forEach(cat => {
    if (!cat || cat.id === undefined || cat.id === null) return
    const id = Number(cat.id)
    if (isNaN(id) || id <= 0) return

    const name = cat.name || `System ${id}`
    const rawSlug = slugifyRaw(name)
    const cleanSlug = slugify(name)
    const canonicalSlug = systemSlugMap[id] || rawSlug || `system-${id}`

    dynamicCategoryMap.set(id, {
      id,
      name,
      slug: canonicalSlug,
    })

    // Register all slug variations to resolve back to ID
    dynamicSlugToIdMap.set(String(id), id)
    dynamicSlugToIdMap.set(`system-${id}`, id)
    if (canonicalSlug) dynamicSlugToIdMap.set(canonicalSlug.toLowerCase(), id)
    if (rawSlug) dynamicSlugToIdMap.set(rawSlug.toLowerCase(), id)
    if (cleanSlug) dynamicSlugToIdMap.set(cleanSlug.toLowerCase(), id)
  })
}

// Pre-register pre-seeded systems immediately
registerCategories(preseededCategories)

/**
 * Convert any string into a clean, URL-safe slug.
 * Example: "01. Setup Master & Meja" -> "setup-master-meja"
 */
export function slugify(text) {
  if (!text) return ''
  return text
    .toString()
    .toLowerCase()
    .trim()
    .replace(/^\d+[\.\-\s]+/, '') // Strip leading module numbers like "01. " or "2. " for cleaner URLs
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

/**
 * Raw slugify without stripping leading numbers (e.g. "01-setup-master")
 */
export function slugifyRaw(text) {
  if (!text) return ''
  return text
    .toString()
    .toLowerCase()
    .trim()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

/**
 * Get canonical slug for a system ID, system object, or system name.
 * NEVER defaults to 'presto' for valid unknown systems.
 *
 * @param {number|string|object} idOrName - Numeric ID, slug/name string, or category object
 * @param {object} [fallbackCat] - Optional fallback category object { id, name }
 * @returns {string} Canonical URL slug
 */
export function getSystemSlug(idOrName, fallbackCat = null) {
  if (idOrName === undefined || idOrName === null || idOrName === '') {
    return 'presto'
  }

  // Handle object passed directly (e.g. sys = { id: 11, name: 'CHANNEL MANAGER (eZee)' })
  if (typeof idOrName === 'object' && idOrName !== null) {
    if (idOrName.name) {
      registerCategories([idOrName])
    }
    idOrName = idOrName.id
  }

  if (fallbackCat && fallbackCat.name && idOrName) {
    registerCategories([{ id: idOrName, name: fallbackCat.name }])
  }

  // Check if numeric ID (number or string containing positive integer)
  const numId = Number(idOrName)
  if (!isNaN(numId) && numId > 0 && Number.isInteger(numId)) {
    // 1. Static map (1..9)
    if (systemSlugMap[numId]) {
      return systemSlugMap[numId]
    }
    // 2. Dynamic registered map
    if (dynamicCategoryMap.has(numId)) {
      return dynamicCategoryMap.get(numId).slug
    }
    // 3. Fallback: dynamic system-{id} so it never redirects to presto
    return `system-${numId}`
  }

  // If string (slug, name, or alias)
  const clean = slugifyRaw(idOrName)
  if (systemAliasToId[clean] !== undefined) {
    const aliasId = systemAliasToId[clean]
    return systemSlugMap[aliasId] || clean
  }
  if (dynamicSlugToIdMap.has(clean)) {
    const dynId = dynamicSlugToIdMap.get(clean)
    if (systemSlugMap[dynId]) return systemSlugMap[dynId]
    if (dynamicCategoryMap.has(dynId)) return dynamicCategoryMap.get(dynId).slug
    return clean
  }

  return clean || 'presto'
}

/**
 * Resolve system ID from slug or numeric parameter.
 * Accepts any positive integer ID, system-XX pattern, or dynamic category slug.
 *
 * @param {string|number} param - Route parameter (slug or numeric ID)
 * @param {number} [defaultId=2] - Fallback ID if parameter cannot be resolved
 * @returns {number} Resolved system ID
 */
export function getSystemIdFromParam(param, defaultId = 2) {
  if (param === undefined || param === null || param === '') {
    return defaultId
  }

  const strParam = String(param).toLowerCase().trim()

  // 1. Static alias map
  if (systemAliasToId[strParam] !== undefined) {
    return systemAliasToId[strParam]
  }

  // 2. Dynamic slug map
  if (dynamicSlugToIdMap.has(strParam)) {
    return dynamicSlugToIdMap.get(strParam)
  }

  // 3. Pure numeric parameter (e.g. "10", "11", 11)
  const numeric = Number(strParam)
  if (!isNaN(numeric) && numeric > 0 && Number.isInteger(numeric)) {
    return numeric
  }

  // 4. system-XX pattern (e.g. "system-10", "system-11")
  const match = strParam.match(/^system-(\d+)$/)
  if (match) {
    const parsed = Number(match[1])
    if (parsed > 0) return parsed
  }

  // 5. Dynamic category scan
  for (const [id, cat] of dynamicCategoryMap.entries()) {
    if (
      cat.slug.toLowerCase() === strParam ||
      slugifyRaw(cat.name) === strParam ||
      slugify(cat.name) === strParam
    ) {
      return id
    }
  }

  return defaultId
}

/**
 * Match a module in modulesList by slug or ID
 */
export function findModuleByParam(modules, param) {
  if (!Array.isArray(modules) || modules.length === 0 || !param) return null
  const strParam = String(param).toLowerCase().trim()
  const numParam = Number(strParam)

  // 1. Direct ID / Order match if param is numeric
  if (!isNaN(numParam)) {
    const byId = modules.find(m => m.id === numParam || m.order === numParam)
    if (byId) return byId
  }

  // 2. Match clean slug (e.g. "setup-master")
  const byCleanSlug = modules.find(m => slugify(m.name) === strParam || slugify(m.title) === strParam)
  if (byCleanSlug) return byCleanSlug

  // 3. Match raw slug (e.g. "01-setup-master")
  const byRawSlug = modules.find(m => slugifyRaw(m.name) === strParam || slugifyRaw(m.title) === strParam)
  if (byRawSlug) return byRawSlug

  // 4. Fallback: fuzzy inclusion match
  const byFuzzy = modules.find(m => {
    const s1 = slugify(m.name)
    const s2 = slugify(m.title)
    return s1.includes(strParam) || strParam.includes(s1) || s2.includes(strParam)
  })
  if (byFuzzy) return byFuzzy

  return null
}

/**
 * Match an article from an articles list by slug or ID
 */
export function findArticleByParam(articles, param) {
  if (!Array.isArray(articles) || articles.length === 0 || !param) return null
  const strParam = String(param).toLowerCase().trim()
  const numParam = Number(strParam)

  // 1. Direct numeric ID match
  if (!isNaN(numParam)) {
    const byId = articles.find(a => Number(a.id) === numParam)
    if (byId) return byId
  }

  // 2. Match slugified title
  const bySlug = articles.find(a => slugify(a.title) === strParam || slugifyRaw(a.title) === strParam)
  if (bySlug) return bySlug

  // 3. Match fuzzy title
  const byFuzzy = articles.find(a => {
    const slug = slugify(a.title)
    return slug.includes(strParam) || strParam.includes(slug)
  })
  if (byFuzzy) return byFuzzy

  return null
}

/**
 * Safely highlight search keyword matches inside text using HTML <mark> tags.
 * Supports multi-word matching, escapes regex metacharacters, and escapes HTML
 * entities to prevent XSS.
 *
 * @param {string} text - Raw input text to search and highlight within
 * @param {string} keyword - Search term entered by user
 * @param {string} markClass - Tailwind CSS classes for the <mark> element
 * @returns {string} Text containing safely highlighted <mark> tags
 */
export function highlightKeyword(
  text,
  keyword,
  markClass = 'bg-amber-200 text-amber-950 font-bold px-1 py-0.5 rounded shadow-2xs'
) {
  if (!text) return ''
  if (!keyword || !String(keyword).trim()) {
    return String(text)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
  }

  const cleanKeyword = String(keyword).trim()
  const words = cleanKeyword.split(/\s+/).filter(w => w.length > 0)
  if (words.length === 0) {
    return String(text)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
  }

  const escapedWords = words.map(w => w.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'))
  const regex = new RegExp(`(${escapedWords.join('|')})`, 'gi')

  const safeText = String(text)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  return safeText.replace(regex, `<mark class="${markClass}">$1</mark>`)
}
