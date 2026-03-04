const {
  getStats,
  getCacheStats
} = require('./statsController');
const { makeRequest } = require('../utils/httpClient');
const logger = require('../utils/logger');

// Mock dependencies
jest.mock('../utils/httpClient');
jest.mock('../utils/logger');

describe('Stats Controller', () => {
  let mockReq, mockRes;

  beforeEach(() => {
    jest.clearAllMocks();

    mockReq = {
      method: 'GET',
      originalUrl: '/api/stats'
    };

    mockRes = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn().mockReturnThis()
    };
  });

  describe('getStats', () => {
    it('should return statistics from backend successfully', async () => {
      const backendStats = {
        users: { total: 10 },
        tasks: {
          total: 25,
          pending: 8,
          inProgress: 12,
          completed: 5
        }
      };

      makeRequest.mockResolvedValue({ data: backendStats });

      await getStats(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/stats');
      expect(mockRes.json).toHaveBeenCalledWith(backendStats);
    });

    it('should handle unwrapped response from backend', async () => {
      const backendStats = {
        users: { total: 5 },
        tasks: {
          total: 15,
          pending: 5,
          inProgress: 7,
          completed: 3
        }
      };

      makeRequest.mockResolvedValue(backendStats);

      await getStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(backendStats);
    });

    it('should handle backend connection error', async () => {
      const connectionError = new Error('ECONNREFUSED');
      makeRequest.mockRejectedValue(connectionError);

      await getStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(503);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'ECONNREFUSED',
        code: 'BACKEND_ERROR'
      });
      expect(logger.logError).toHaveBeenCalledWith(
        'GET',
        '/api/stats',
        503,
        0,
        'GetStatsError',
        'ECONNREFUSED'
      );
    });

    it('should handle backend server error', async () => {
      const serverError = new Error('Internal Server Error');
      serverError.statusCode = 500;
      makeRequest.mockRejectedValue(serverError);

      await getStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Internal Server Error',
        code: 'BACKEND_ERROR'
      });
    });

    it('should handle empty response from backend', async () => {
      makeRequest.mockResolvedValue({});

      await getStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith({});
    });

    it('should handle null response from backend', async () => {
      makeRequest.mockResolvedValue(null);

      await getStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith({
        users: { total: 0 },
        tasks: { total: 0, pending: 0, inProgress: 0, completed: 0 }
      });
    });

    it('should handle partial statistics data', async () => {
      const partialStats = {
        users: { total: 5 }
        // tasks data missing
      };

      makeRequest.mockResolvedValue({ data: partialStats });

      await getStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(partialStats);
    });

    it('should handle malformed response structure', async () => {
      const malformedResponse = {
        invalid: 'structure',
        not: 'what we expect'
      };

      makeRequest.mockResolvedValue(malformedResponse);

      await getStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(malformedResponse);
    });

    it('should use correct endpoint URL', async () => {
      const backendStats = { users: { total: 1 } };
      makeRequest.mockResolvedValue({ data: backendStats });

      await getStats(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/stats');
      expect(makeRequest).toHaveBeenCalledTimes(1);
    });
  });

  describe('getCacheStats', () => {
    it('should return cache statistics from backend successfully', async () => {
      const cacheStats = {
        hits: 150,
        misses: 25,
        evictions: 0,
        totalEntries: 10
      };

      makeRequest.mockResolvedValue({ data: cacheStats });

      await getCacheStats(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/stats/cache');
      expect(mockRes.json).toHaveBeenCalledWith(cacheStats);
    });

    it('should handle unwrapped cache stats response', async () => {
      const cacheStats = {
        hits: 200,
        misses: 30,
        evictions: 5,
        totalEntries: 15
      };

      makeRequest.mockResolvedValue(cacheStats);

      await getCacheStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(cacheStats);
    });

    it('should handle backend connection error for cache stats', async () => {
      const connectionError = new Error('ECONNREFUSED');
      makeRequest.mockRejectedValue(connectionError);

      await getCacheStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(503);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'ECONNREFUSED',
        code: 'BACKEND_ERROR'
      });
      expect(logger.logError).toHaveBeenCalledWith(
        'GET',
        '/api/stats/cache',
        503,
        0,
        'GetCacheStatsError',
        'ECONNREFUSED'
      );
    });

    it('should handle backend server error for cache stats', async () => {
      const serverError = new Error('Cache service unavailable');
      serverError.statusCode = 503;
      makeRequest.mockRejectedValue(serverError);

      await getCacheStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(503);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Cache service unavailable',
        code: 'BACKEND_ERROR'
      });
    });

    it('should handle empty cache stats response', async () => {
      makeRequest.mockResolvedValue({});

      await getCacheStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith({});
    });

    it('should handle zero cache statistics', async () => {
      const zeroStats = {
        hits: 0,
        misses: 0,
        evictions: 0,
        totalEntries: 0
      };

      makeRequest.mockResolvedValue({ data: zeroStats });

      await getCacheStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(zeroStats);
    });

    it('should handle high cache statistics', async () => {
      const highStats = {
        hits: 1000000,
        misses: 50000,
        evictions: 1000,
        totalEntries: 50000
      };

      makeRequest.mockResolvedValue({ data: highStats });

      await getCacheStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(highStats);
    });

    it('should use correct cache stats endpoint URL', async () => {
      const cacheStats = { hits: 100 };
      makeRequest.mockResolvedValue({ data: cacheStats });

      await getCacheStats(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/stats/cache');
      expect(makeRequest).toHaveBeenCalledTimes(1);
    });

    it('should update request URL for cache stats endpoint', async () => {
      mockReq.originalUrl = '/api/stats/cache';

      const cacheStats = { hits: 75 };
      makeRequest.mockResolvedValue({ data: cacheStats });

      await getCacheStats(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/stats/cache');
      expect(mockRes.json).toHaveBeenCalledWith(cacheStats);
    });
  });

  describe('Error Handling Edge Cases', () => {
    it('should handle makeRequest throwing non-Error objects in getStats', async () => {
      makeRequest.mockRejectedValue('String error');

      await getStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'String error',
        code: 'BACKEND_ERROR'
      });
    });

    it('should handle makeRequest throwing non-Error objects in getCacheStats', async () => {
      makeRequest.mockRejectedValue('String error');

      await getCacheStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'String error',
        code: 'BACKEND_ERROR'
      });
    });

    it('should handle makeRequest throwing undefined', async () => {
      makeRequest.mockRejectedValue(undefined);

      await getStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith(
        expect.objectContaining({
          code: 'BACKEND_ERROR'
        })
      );
    });

    it('should handle makeRequest throwing null', async () => {
      makeRequest.mockRejectedValue(null);

      await getCacheStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith(
        expect.objectContaining({
          code: 'BACKEND_ERROR'
        })
      );
    });

    it('should handle network timeout errors', async () => {
      const timeoutError = new Error('Network timeout');
      timeoutError.code = 'ETIMEDOUT';
      makeRequest.mockRejectedValue(timeoutError);

      await getStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Network timeout',
        code: 'BACKEND_ERROR'
      });
    });

    it('should handle HTTP 404 errors appropriately', async () => {
      const notFoundError = new Error('Stats endpoint not found');
      notFoundError.statusCode = 404;
      makeRequest.mockRejectedValue(notFoundError);

      await getStats(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(404);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Stats endpoint not found',
        code: 'BACKEND_ERROR'
      });
    });
  });

  describe('Response Format Validation', () => {
    it('should handle response with nested objects', async () => {
      const nestedResponse = {
        users: {
          total: 10,
          active: 8,
          inactive: 2
        },
        tasks: {
          total: 25,
          byStatus: {
            pending: 8,
            inProgress: 12,
            completed: 5
          },
          byPriority: {
            high: 5,
            medium: 15,
            low: 5
          }
        }
      };

      makeRequest.mockResolvedValue({ data: nestedResponse });

      await getStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(nestedResponse);
    });

    it('should handle response with arrays', async () => {
      const arrayResponse = {
        users: [
          { id: 1, name: 'John', taskCount: 5 },
          { id: 2, name: 'Jane', taskCount: 3 }
        ],
        tasks: [
          { status: 'pending', count: 8 },
          { status: 'completed', count: 5 }
        ]
      };

      makeRequest.mockResolvedValue({ data: arrayResponse });

      await getStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(arrayResponse);
    });

    it('should handle response with mixed data types', async () => {
      const mixedResponse = {
        users: { total: 10 },
        tasks: { total: 25 },
        uptime: '2h30m15s',
        version: '1.0.0',
        timestamp: '2026-03-04T12:00:00Z',
        healthy: true,
        memoryUsage: 85.5
      };

      makeRequest.mockResolvedValue({ data: mixedResponse });

      await getStats(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(mixedResponse);
    });
  });
});
