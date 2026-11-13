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
 * Get canonical slug for a system ID or system name
 */
export function getSystemSlug(idOrName) {
  if (typeof idOrName === 'number' || (!isNaN(idOrName) && idOrName !== '')) {
    return systemSlugMap[Number(idOrName)] || 'presto'
  }
  const clean = slugifyRaw(idOrName)
  const id = systemAliasToId[clean]
  return id ? systemSlugMap[id] : clean
}

/**
 * Resolve system ID from slug or numeric parameter.
 * Defaults to 2 (PRESTO) if not found.
 */
export function getSystemIdFromParam(param, defaultId = 2) {
  if (!param) return defaultId
  const strParam = String(param).toLowerCase().trim()
  if (systemAliasToId[strParam] !== undefined) {
    return systemAliasToId[strParam]
  }
  const numeric = Number(strParam)
  if (!isNaN(numeric) && systemSlugMap[numeric]) {
    return numeric
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
