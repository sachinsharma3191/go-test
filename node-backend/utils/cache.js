/**
 * In-memory string cache for serialization/deserialization.
 * Stores values as strings; 5 min TTL; invalidate on mutations.
 */

const TTL_MS = 5 * 60 * 1000; // 5 minutes

const store = new Map(); // key -> { value: string, expiresAt: number }
const stats = { hits: 0, misses: 0, evictions: 0, totalEntries: 0 };

function cleanupExpired() {
  const now = Date.now();
  for (const [key, entry] of store.entries()) {
    if (entry.expiresAt && now > entry.expiresAt) {
      store.delete(key);
      stats.evictions++;
    }
  }
  stats.totalEntries = store.size;
}

// Run cleanup every minute (.unref() allows Jest/Node to exit when tests finish)
const cleanupInterval = setInterval(cleanupExpired, 60 * 1000);
if (cleanupInterval.unref) cleanupInterval.unref();

function get(key) {
  if (key == null) return null;
  const entry = store.get(key);
  if (!entry) {
    stats.misses++;
    return null;
  }
  if (entry.expiresAt && Date.now() > entry.expiresAt) {
    store.delete(key);
    stats.misses++;
    stats.evictions++;
    return null;
  }
  stats.hits++;
  return entry.value;
}

function set(key, value, ttlMs) {
  if (key == null) return;
  let expiresAt = null;
  if (ttlMs !== undefined) {
    if (ttlMs > 0) expiresAt = Date.now() + ttlMs;
    else if (ttlMs === 0) { store.delete(key); stats.totalEntries = store.size; return; }
  } else if (TTL_MS > 0) {
    expiresAt = Date.now() + TTL_MS;
  }
  store.set(key, { value, expiresAt });
  stats.totalEntries = store.size;
}

function invalidate(pattern) {
  if (pattern == null) return;
  for (const key of store.keys()) {
    if (key != null && pattern != null && String(key).includes(String(pattern))) {
      store.delete(key);
      stats.evictions++;
    }
  }
  stats.totalEntries = store.size;
}

function getStats() {
  stats.totalEntries = store.size;
  return { ...stats };
}

function clear() {
  store.clear();
  stats.hits = 0;
  stats.misses = 0;
  stats.evictions = 0;
  stats.totalEntries = 0;
}

module.exports = {
  get,
  set,
  invalidate,
  getStats,
  clear
};
