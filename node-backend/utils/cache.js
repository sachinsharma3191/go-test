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

// Run cleanup every minute
setInterval(cleanupExpired, 60 * 1000);

function get(key) {
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

function set(key, value) {
  const expiresAt = TTL_MS > 0 ? Date.now() + TTL_MS : null;
  store.set(key, { value, expiresAt });
  stats.totalEntries = store.size;
}

function invalidate(pattern) {
  for (const key of store.keys()) {
    if (key.includes(pattern)) {
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

module.exports = {
  get,
  set,
  invalidate,
  getStats
};
