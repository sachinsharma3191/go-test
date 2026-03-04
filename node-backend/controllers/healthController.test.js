const { getHealth } = require('./healthController');
const { makeRequest } = require('../utils/httpClient');

// Mock the httpClient module
jest.mock('../utils/httpClient');

describe('Health Controller', () => {
  let mockReq, mockRes;

  beforeEach(() => {
    // Reset all mocks before each test
    jest.clearAllMocks();

    // Mock request object
    mockReq = {
      method: 'GET',
      originalUrl: '/health'
    };

    // Mock response object
    mockRes = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn().mockReturnThis()
    };
  });

  describe('getHealth', () => {
    it('should return healthy status when Go backend is available', async () => {
      // Mock successful response from Go backend
      makeRequest.mockResolvedValue({
        status: 'healthy',
        message: 'All systems operational'
      });

      await getHealth(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith({
        status: 'ok',
        message: 'Node.js backend is running',
        goBackend: {
          status: 'healthy',
          message: 'All systems operational'
        }
      });

      expect(makeRequest).toHaveBeenCalledWith('/health');
    });

    it('should return error status when Go backend is unavailable', async () => {
      // Mock error from Go backend
      const backendError = new Error('ECONNREFUSED');
      makeRequest.mockRejectedValue(backendError);

      await getHealth(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(503);
      expect(mockRes.json).toHaveBeenCalledWith({
        status: 'error',
        message: 'Node.js backend is running but Go backend is unavailable',
        error: 'ECONNREFUSED'
      });

      expect(makeRequest).toHaveBeenCalledWith('/health');
    });

    it('should handle network errors gracefully', async () => {
      // Mock network error
      const networkError = new Error('Network timeout');
      makeRequest.mockRejectedValue(networkError);

      await getHealth(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(503);
      expect(mockRes.json).toHaveBeenCalledWith({
        status: 'error',
        message: 'Node.js backend is running but Go backend is unavailable',
        error: 'Network timeout'
      });
    });

    it('should handle malformed Go backend response', async () => {
      // Mock malformed response
      makeRequest.mockResolvedValue(null);

      await getHealth(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith({
        status: 'ok',
        message: 'Node.js backend is running',
        goBackend: null
      });
    });

    it('should handle empty Go backend response', async () => {
      // Mock empty response
      makeRequest.mockResolvedValue({});

      await getHealth(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith({
        status: 'ok',
        message: 'Node.js backend is running',
        goBackend: {}
      });
    });

    it('should preserve Go backend response structure', async () => {
      // Mock complex response from Go backend
      const complexResponse = {
        status: 'degraded',
        message: 'Some services degraded',
        checks: {
          database: { status: 'healthy', latency: '5ms' },
          cache: { status: 'unhealthy', latency: '100ms' }
        },
        timestamp: '2026-03-04T12:00:00Z'
      };

      makeRequest.mockResolvedValue(complexResponse);

      await getHealth(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith({
        status: 'ok',
        message: 'Node.js backend is running',
        goBackend: complexResponse
      });
    });
  });

  describe('Error Handling', () => {
    it('should handle makeRequest throwing non-Error objects', async () => {
      // Mock non-Error throw
      makeRequest.mockRejectedValue('String error');

      await getHealth(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(503);
      expect(mockRes.json).toHaveBeenCalledWith({
        status: 'error',
        message: 'Node.js backend is running but Go backend is unavailable',
        error: 'String error'
      });
    });

    it('should handle makeRequest throwing undefined', async () => {
      // Mock undefined throw - should not crash
      makeRequest.mockRejectedValue(undefined);

      await getHealth(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(503);
      expect(mockRes.json).toHaveBeenCalledWith(
        expect.objectContaining({
          status: 'error',
          message: 'Node.js backend is running but Go backend is unavailable'
        })
      );
    });
  });

  describe('Request Logging', () => {
    it('should not log errors when request is successful', async () => {
      // Mock successful response
      makeRequest.mockResolvedValue({ status: 'healthy' });

      // Mock logger to ensure no error logging
      const logger = require('../utils/logger');
      const logErrorSpy = jest.spyOn(logger, 'logError');

      await getHealth(mockReq, mockRes);

      expect(logErrorSpy).not.toHaveBeenCalled();
    });

    it('should not log errors when Go backend is unavailable (expected behavior)', async () => {
      // Mock backend error
      makeRequest.mockRejectedValue(new Error('ECONNREFUSED'));

      // Mock logger
      const logger = require('../utils/logger');
      const logErrorSpy = jest.spyOn(logger, 'logError');

      await getHealth(mockReq, mockRes);

      // Should not log errors for expected backend unavailability
      expect(logErrorSpy).not.toHaveBeenCalled();
    });
  });
});
