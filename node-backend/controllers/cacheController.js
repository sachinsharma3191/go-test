const cache = require('../utils/cache');

/**
 * Get cache statistics (hits, misses, evictions, total entries)
 */
const getCacheStats = (req, res) => {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'Method not allowed' });
  }
  const stats = cache.getStats();
  res.json(stats);
};

module.exports = {
  getCacheStats
};
