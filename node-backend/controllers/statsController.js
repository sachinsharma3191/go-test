const { makeRequest } = require('../utils/httpClient');
const logger = require('../utils/logger');
const cache = require('../utils/cache');

/**
 * Get statistics
 */
const getStats = async (req, res) => {
  try {
    const stats = await makeRequest('/api/stats');

    // Return shape frontend expects (unwrap if Go sent envelope)
    const data = stats?.data ?? stats;
    if (!data) {
      return res.json({
        users: { total: 0 },
        tasks: { total: 0, pending: 0, inProgress: 0, completed: 0 }
      });
    }

    // Refresh tasks cache when stats has more tasks than cache (keeps them in sync)
    const statsTaskTotal = data?.tasks?.total ?? 0;
    const cachedTasks = cache.get('tasks');
    if (cachedTasks) {
      try {
        const parsed = JSON.parse(cachedTasks);
        const cachedCount = parsed?.count ?? parsed?.tasks?.length ?? 0;
        if (statsTaskTotal > cachedCount) {
          cache.invalidate('tasks');
        }
      } catch {
        /* ignore */
      }
    }

    res.json(data);
  } catch (error) {
    const statusCode = error.statusCode || (error.message.includes('ECONNREFUSED') ? 503 : 500);
    res.status(statusCode).json({
      error: error.message || 'Failed to fetch stats',
      code: error.code || 'BACKEND_ERROR'
    });
  }
};

/**
 * Get cache statistics
 */
const getCacheStats = async (req, res) => {
  try {
    const raw = await makeRequest('/api/cache/stats');
    const cacheStats = raw?.data ?? raw;
    logger.logCacheStats(cacheStats?.cache ?? cacheStats);
    res.json(cacheStats);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

module.exports = {
  getStats,
  getCacheStats
};
