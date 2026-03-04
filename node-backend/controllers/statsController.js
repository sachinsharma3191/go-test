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
    if (data === undefined || data === null) {
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
    const errMsg = error?.message ?? String(error ?? 'Failed to fetch stats');
    const statusCode = error?.statusCode || ((errMsg.includes('ECONNREFUSED') || errMsg.includes('unavailable')) ? 503 : 500);
    logger.logError(req.method, req.originalUrl, statusCode, 0, 'GetStatsError', errMsg);
    res.status(statusCode).json({
      error: errMsg,
      code: 'BACKEND_ERROR'
    });
  }
};

/**
 * Get cache statistics
 */
const getCacheStats = async (req, res) => {
  try {
    const raw = await makeRequest('/api/stats/cache');
    const cacheStats = raw?.data ?? raw;
    logger.logCacheStats(cacheStats?.cache ?? cacheStats);
    res.json(cacheStats);
  } catch (error) {
    const errMsg = error?.message ?? String(error ?? 'Failed to fetch cache stats');
    const statusCode = error?.statusCode || ((errMsg.includes('ECONNREFUSED') || errMsg.includes('unavailable')) ? 503 : 500);
    logger.logError(req.method, '/api/stats/cache', statusCode, 0, 'GetCacheStatsError', errMsg);
    res.status(statusCode).json({
      error: errMsg,
      code: 'BACKEND_ERROR'
    });
  }
};

module.exports = {
  getStats,
  getCacheStats
};
